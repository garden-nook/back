package main

import (
	"context"
	"errors"
	"garden-nook/docs"
	"garden-nook/internal/config"
	"garden-nook/internal/middleware"
	"garden-nook/internal/migrator"
	"garden-nook/internal/modules/auth"
	"garden-nook/internal/modules/crops"
	"garden-nook/internal/modules/plot"
	"garden-nook/internal/modules/plot/repositories"
	"garden-nook/internal/modules/plot/services"
	"garden-nook/internal/pkg/database"
	"garden-nook/internal/pkg/jwt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "garden-nook/docs"
)

// @title           Garden Nook API
// @version         0.0.7
// @host            localhost:8000
// @BasePath        /
// @securityDefinitions.apikey  UserAuth
// @in                          header
// @name                        Authorization
// @description                 Access JWT token в формате: Bearer {token}
// @securityDefinitions.apikey  AdminAuth
// @in                          header
// @name                        Authorization
// @description                 Access JWT token в формате: Bearer {token}
func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg := config.Load()

	if err := migrator.Run(cfg.Database.DSN(), log); err != nil {
		log.Error("migration failed, aborting startup", "err", err)
		os.Exit(1)
	}

	pool, err := database.NewPool(cfg.Database.DSN())
	if err != nil {
		log.Error("db connection failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	jwtMgr := jwt.NewManager(cfg.JWT.AccessSecret, cfg.JWT.UserAccessTTL, cfg.JWT.AdminAccessTTL)
	authMW := middleware.NewAuth(jwtMgr)

	errorMapper := database.NewErrorMapper(map[string]error{}, log)

	// Модуль auth
	authRepo := auth.NewRepository(pool)
	authSvc := auth.NewService(authRepo, jwtMgr, log)
	authCtrl := auth.NewController(authSvc)

	// Модуль crops
	cropsRepo := crops.NewRepository(pool, errorMapper)
	cropsSvc := crops.NewService(cropsRepo, log)
	cropsCtrl := crops.NewController(cropsSvc)

	// Модуль plots
	plotsRepo := repositories.NewRepository(pool, errorMapper)
	//plotsProjector := plots.NewProjector(plotsRepo, log)
	plotsSvc := services.NewService(plotsRepo /*plotsProjector,*/, log)
	plotsCtrl := plots.NewController(plotsSvc)

	r := chi.NewRouter()

	// Базовые middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.ClientIPFromHeader("X-Real-IP"))
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))

	// CORS (упрощённо)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Swagger UI
	docs.SwaggerInfo.Host = cfg.Docs.Host
	docs.SwaggerInfo.Schemes = []string{cfg.Docs.Schema}
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Регистрация модулей
	auth.RegisterRoutes(r, authCtrl, authMW)
	crops.RegisterRoutes(r, cropsCtrl, authMW)
	plots.RegisterRoutes(r, plotsCtrl, authMW)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Info("server starting", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Info("server stopped")
}
