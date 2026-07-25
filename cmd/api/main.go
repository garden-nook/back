package main

import (
	"context"
	"errors"
	docs "garden-nook/docs/gen"
	"garden-nook/internal/config"
	"garden-nook/internal/middleware"
	"garden-nook/internal/migrator"
	"garden-nook/internal/modules/auth"
	"garden-nook/internal/modules/crops"
	cropRepos "garden-nook/internal/modules/crops/repository"
	cropSvcs "garden-nook/internal/modules/crops/service"
	"garden-nook/internal/modules/plot"
	plotRepos "garden-nook/internal/modules/plot/repository"
	plotSvcs "garden-nook/internal/modules/plot/service"
	"garden-nook/internal/pkg/database"
	"garden-nook/internal/pkg/helpers"
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
)

// @title           Garden Nook API
// @version         0.2.0
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
	seh := helpers.NewServiceErrorHandler(log)

	// Модуль auth
	authRepo := auth.NewRepository(pool)
	authSvc := auth.NewService(authRepo, jwtMgr, log)
	authCtrl := auth.NewController(authSvc)

	// Модуль crops
	soilTypeRepo := cropRepos.NewSoilTypeRepo(pool, errorMapper)
	cropFamilyRepo := cropRepos.NewCropFamilyRepo(pool, errorMapper)
	cropRepo := cropRepos.NewCropRepo(pool, errorMapper)
	cropRuleRepo := cropRepos.NewCropRuleRepo(pool, errorMapper)
	ruleCache, err := cropSvcs.NewRuleCache(cropRuleRepo)
	if err != nil {
		log.Error("crop rule cache refresh failed", "err", err)
		os.Exit(1)
	}
	soilTypeSvc := cropSvcs.NewSoilTypeService(soilTypeRepo, seh)
	cropFamilySvc := cropSvcs.NewCropFamilyService(cropFamilyRepo, seh)
	cropSvc := cropSvcs.NewCropService(cropRepo, cropFamilyRepo, cropRuleRepo, seh)
	cropRuleSvc := cropSvcs.NewCropRuleService(cropRuleRepo, ruleCache, seh)
	cropsCtrl := crops.NewController(soilTypeSvc, cropFamilySvc, cropSvc, cropRuleSvc)

	// Модуль plots
	plotRepo := plotRepos.NewPlotRepo(pool, errorMapper)
	bedRepo := plotRepos.NewBedRepo(pool, errorMapper)
	objectRepo := plotRepos.NewObjectRepo(pool, errorMapper)
	gridRepo := plotRepos.NewGridCellRepo(pool, errorMapper)
	eventRepo := plotRepos.NewEventStoreRepo(pool, errorMapper)
	historyRepo := plotRepos.NewHistoryRepo(pool, errorMapper)
	plotsSvc := plotSvcs.NewPlotService(pool, plotRepo, gridRepo, bedRepo, objectRepo, eventRepo, seh)
	eventSvc := plotSvcs.NewEventService(pool, plotRepo, bedRepo, objectRepo, eventRepo, historyRepo, seh)
	recommendationSvc := plotSvcs.NewRecommendationService(bedRepo, plotRepo, gridRepo, historyRepo, cropSvc, ruleCache, seh)
	plotsCtrl := plots.NewController(plotsSvc, eventSvc, recommendationSvc)

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
