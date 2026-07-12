package crops

import (
	"garden-nook/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes монтирует маршруты модуля crops в основной роутер.
func RegisterRoutes(r chi.Router, ctrl *Controller, auth *middleware.AuthMiddleware) {
	// ---- Публичные эндпоинты (только чтение справочников) ----
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/crop-families", ctrl.ListFamilies)
		r.Get("/crop-families/{id}", ctrl.GetFamily)
		r.Get("/crops", ctrl.ListCrops)
		r.Get("/crops/{id}", ctrl.GetCrop)
	})

	// ---- Админские эндпоинты ----
	r.Route("/api/v1/admin", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAdmin)

			r.Post("/crop-families", ctrl.CreateFamily)
			r.Put("/crop-families/{id}", ctrl.UpdateFamily)
			r.Delete("/crop-families/{id}", ctrl.DeleteFamily)

			r.Post("/crops", ctrl.CreateCrop)
			r.Put("/crops/{id}", ctrl.UpdateCrop)
			r.Delete("/crops/{id}", ctrl.DeleteCrop)

			r.Get("/crop-rules", ctrl.ListRules)
			r.Post("/crop-rules", ctrl.CreateRule)
			r.Delete("/crop-rules/{id}", ctrl.DeleteRule)
		})
	})
}
