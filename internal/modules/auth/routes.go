package auth

import (
	"garden-nook/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, ctrl *Controller, auth *middleware.AuthMiddleware) {
	// ---- Публичные эндпоинты (без токена) ----
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", ctrl.RegisterUser)
		r.Post("/login", ctrl.LoginUser)

		// ---- User-эндпоинты (требуют user-токен) ----
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireUser)
			r.Get("/me", ctrl.MeUser)
		})
	})

	// ---- Admin auth (отдельный эндпоинт входа, без токена) ----
	r.Post("/api/v1/admin/auth/login", ctrl.LoginAdmin)

	// ---- Admin-эндпоинты (требуют admin-токен) ----
	r.Route("/api/v1/admin/auth", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAdmin)
			r.Get("/me", ctrl.MeAdmin)
		})
	})
}
