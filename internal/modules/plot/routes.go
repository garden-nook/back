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

			// Состояние участка
			//r.Get("/{id}/state", ctrl.GetPlotState)

			// Грядки
			//r.Post("/{plot_id}/beds", ctrl.CreateBed)
			//r.Get("/{plot_id}/beds", ctrl.ListBeds)
			//r.Put("/{plot_id}/beds/{bed_id}", ctrl.UpdateBed)
			//r.Delete("/{plot_id}/beds/{bed_id}", ctrl.DeleteBed)

			// Объекты
			//r.Post("/{plot_id}/objects", ctrl.CreateObject)
			//r.Get("/{plot_id}/objects", ctrl.ListObjects)
			//r.Delete("/{plot_id}/objects/{object_id}", ctrl.DeleteObject)

			// Посадки
			//r.Post("/{plot_id}/plant", ctrl.PlantCrop)
			//r.Post("/{plot_id}/harvest", ctrl.HarvestCrop)

			// Timeline
			//r.Get("/{plot_id}/timeline", ctrl.GetTimeline)
		})
	})
}
