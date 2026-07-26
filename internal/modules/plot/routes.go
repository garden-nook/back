package plots

import (
	"garden-nook/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, ctrl *Controller, auth *middleware.AuthMiddleware) {
	r.Route("/api/v1/plots", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireUser)

			// CRUD участков
			r.Get("/", ctrl.ListPlots)
			r.Post("/", ctrl.CreatePlot)
			r.Put("/{id}", ctrl.UpdatePlot)
			r.Delete("/{id}", ctrl.DeletePlot)

			// Конструктор участка
			r.Get("/{id}/structure", ctrl.GetPlotStructure)
			r.Post("/{id}/events", ctrl.HandleEvents)

			r.Get("/bed/{id}/recommendation", ctrl.GetBedRecommendations)
			r.Get("/bed/{id}/history", ctrl.GetBedCropHistory)

			// Timeline
			//r.Get("/{plot_id}/timeline", ctrl.GetTimeline)
		})
	})
}
