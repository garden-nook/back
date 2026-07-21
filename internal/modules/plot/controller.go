package plots

import (
	"garden-nook/internal/middleware"
	"garden-nook/internal/pkg/helpers"
	"garden-nook/internal/pkg/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	svc *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{svc: svc}
}

// ---------- PLOTS ----------

// ListPlots godoc
// @Summary      Список участков пользователя
// @Tags         plots
// @Produce      json
// @Security     UserAuth
// @Success      200 {object} response.Response{data=[]Plot}
// @Router       /api/v1/plots [get]
func (c *Controller) ListPlots(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.SubID(r.Context())
	plots, _, err := c.svc.ListPlots(r.Context(), ownerID)
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, plots)
}

// CreatePlot godoc
// @Summary      Создать участок
// @Tags         plots
// @Accept       json
// @Produce      json
// @Security     UserAuth
// @Param        body  body      CreatePlotRequest  true  "Данные участка"
// @Success      201 {object} response.Response{data=response.CreateUpdateUuidId}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/plots [post]
func (c *Controller) CreatePlot(w http.ResponseWriter, r *http.Request) {
	var req CreatePlotRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}

	ownerID := middleware.SubID(r.Context())
	plot, err := c.svc.CreatePlot(r.Context(), ownerID, req)
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, plot)
}

// UpdatePlot godoc
// @Summary      Обновить участок
// @Tags         plots
// @Accept       json
// @Produce      json
// @Security     UserAuth
// @Param        id    path      int                true  "ID участка"
// @Param        body  body      UpdatePlotRequest  true  "Данные участка"
// @Success      200 {object} response.Response{data=response.CreateUpdateUuidId}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/plots/{id} [put]
func (c *Controller) UpdatePlot(w http.ResponseWriter, r *http.Request) {
	var req UpdatePlotRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}

	plotID := chi.URLParam(r, "id")
	ownerID := middleware.SubID(r.Context())
	plot, err := c.svc.UpdatePlot(r.Context(), plotID, ownerID, req)
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, plot)
}

// DeletePlot godoc
// @Summary      Удалить участок
// @Tags         plots
// @Produce      json
// @Security     UserAuth
// @Param        id  path      string  true  "Plot ID"
// @Success      204 "No Content"
// @Router       /api/v1/plots/{id} [delete]
func (c *Controller) DeletePlot(w http.ResponseWriter, r *http.Request) {
	plotID := chi.URLParam(r, "id")
	ownerID := middleware.SubID(r.Context())

	if err := c.svc.DeletePlot(r.Context(), plotID, ownerID); err != nil {
		helpers.WriteErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

//// GetPlotState godoc
//// @Summary      Получить полное состояние участка (грядки, объекты, сетка)
//// @Tags         plots
//// @Produce      json
//// @Security     BearerAuth
//// @Param        id  path      string  true  "Plot ID"
//// @Success      200 {object} response.Response{data=PlotState}
//// @Router       /api/v1/plots/{id}/state [get]
//func (c *Controller) GetPlotState(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "id")
//	ownerID := middleware.SubID(r.Context())
//
//	state, err := c.svc.GetPlotState(r.Context(), plotID, ownerID)
//	if err != nil {
//		helpers.WriteErr(w, err)
//		return
//	}
//	response.JSON(w, http.StatusOK, state)
//}
//
//// ---------- BEDS ----------
//
//// CreateBed godoc
//// @Summary      Создать грядку
//// @Tags         beds
//// @Accept       json
//// @Produce      json
//// @Security     BearerAuth
//// @Param        plot_id  path      string             true  "Plot ID"
//// @Param        body     body      CreateBedRequest   true  "Данные грядки"
//// @Success      201 {object} response.Response{data=Bed}
//// @Router       /api/v1/plots/{plot_id}/beds [post]
//func (c *Controller) CreateBed(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "plot_id")
//	ownerID := middleware.SubID(r.Context())
//
//	var req CreateBedRequest
//	if err := middleware.DecodeAndValidate(r, &req); err != nil {
//		middleware.WriteValidationError(w, err)
//		return
//	}
//
//	bed, err := c.svc.CreateBed(r.Context(), plotID, ownerID, req)
//	if err != nil {
//		helpers.WriteErr(w, err)
//		return
//	}
//	response.JSON(w, http.StatusCreated, bed)
//}
//
//// ListBeds godoc
//// @Summary      Список грядок участка
//// @Tags         beds
//// @Produce      json
//// @Security     BearerAuth
//// @Param        plot_id  path      string  true  "Plot ID"
//// @Success      200 {object} response.Response{data=[]Bed}
//// @Router       /api/v1/plots/{plot_id}/beds [get]
//func (c *Controller) ListBeds(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "plot_id")
//	ownerID := middleware.SubID(r.Context())
//
//	beds, err := c.svc.ListBeds(r.Context(), plotID, ownerID)
//	if err != nil {
//		helpers.WriteErr(w, err)
//		return
//	}
//	response.JSON(w, http.StatusOK, beds)
//}
//
//// UpdateBed godoc
//// @Summary      Обновить грядку
//// @Tags         beds
//// @Accept       json
//// @Produce      json
//// @Security     BearerAuth
//// @Param        plot_id  path      string             true  "Plot ID"
//// @Param        bed_id   path      string             true  "Bed ID"
//// @Param        body     body      UpdateBedRequest   true  "Обновления"
//// @Success      200 {object} response.Response{data=Bed}
//// @Router       /api/v1/plots/{plot_id}/beds/{bed_id} [put]
//func (c *Controller) UpdateBed(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "plot_id")
//	bedID := chi.URLParam(r, "bed_id")
//	ownerID := middleware.SubID(r.Context())
//
//	var req UpdateBedRequest
//	if err := middleware.DecodeAndValidate(r, &req); err != nil {
//		middleware.WriteValidationError(w, err)
//		return
//	}
//
//	bed, err := c.svc.UpdateBed(r.Context(), plotID, bedID, ownerID, req)
//	if err != nil {
//		helpers.WriteErr(w, err)
//		return
//	}
//	response.JSON(w, http.StatusOK, bed)
//}
//
//// DeleteBed godoc
//// @Summary      Удалить грядку
//// @Tags         beds
//// @Produce      json
//// @Security     BearerAuth
//// @Param        plot_id  path      string  true  "Plot ID"
//// @Param        bed_id   path      string  true  "Bed ID"
//// @Success      204 "No Content"
//// @Router       /api/v1/plots/{plot_id}/beds/{bed_id} [delete]
//func (c *Controller) DeleteBed(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "plot_id")
//	bedID := chi.URLParam(r, "bed_id")
//	ownerID := middleware.SubID(r.Context())
//
//	if err := c.svc.DeleteBed(r.Context(), plotID, bedID, ownerID); err != nil {
//		helpers.WriteErr(w, err)
//		return
//	}
//	w.WriteHeader(http.StatusNoContent)
//}
//
//// ---------- OBJECTS ----------
//
//// CreateObject godoc
//// @Summary      Создать объект (дерево, постройку)
//// @Tags         objects
//// @Accept       json
//// @Produce      json
//// @Security     BearerAuth
//// @Param        plot_id  path      string               true  "Plot ID"
//// @Param        body     body      CreateObjectRequest  true  "Данные объекта"
//// @Success      201 {object} response.Response{data=UIObject}
//// @Router       /api/v1/plots/{plot_id}/objects [post]
//func (c *Controller) CreateObject(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "plot_id")
//	ownerID := middleware.SubID(r.Context())
//
//	var req CreateObjectRequest
//	if err := middleware.DecodeAndValidate(r, &req); err != nil {
//		middleware.WriteValidationError(w, err)
//		return
//	}
//
//	obj, err := c.svc.CreateObject(r.Context(), plotID, ownerID, req)
//	if err != nil {
//		helpers.WriteErr(w, err)
//		return
//	}
//	response.JSON(w, http.StatusCreated, obj)
//}
//
//// ListObjects godoc
//// @Summary      Список объектов участка
//// @Tags         objects
//// @Produce      json
//// @Security     BearerAuth
//// @Param        plot_id  path      string  true  "Plot ID"
//// @Success      200 {object} response.Response{data=[]UIObject}
//// @Router       /api/v1/plots/{plot_id}/objects [get]
//func (c *Controller) ListObjects(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "plot_id")
//	ownerID := middleware.SubID(r.Context())
//
//	objects, err := c.svc.ListObjects(r.Context(), plotID, ownerID)
//	if err != nil {
//		helpers.WriteErr(w, err)
//		return
//	}
//	response.JSON(w, http.StatusOK, objects)
//}
//
//// DeleteObject godoc
//// @Summary      Удалить объект
//// @Tags         objects
//// @Produce      json
//// @Security     BearerAuth
//// @Param        plot_id    path      string  true  "Plot ID"
//// @Param        object_id  path      string  true  "Object ID"
//// @Success      204 "No Content"
//// @Router       /api/v1/plots/{plot_id}/objects/{object_id} [delete]
//func (c *Controller) DeleteObject(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "plot_id")
//	objectID := chi.URLParam(r, "object_id")
//	ownerID := middleware.SubID(r.Context())
//
//	if err := c.svc.DeleteObject(r.Context(), plotID, objectID, ownerID); err != nil {
//		helpers.WriteErr(w, err)
//		return
//	}
//	w.WriteHeader(http.StatusNoContent)
//}
//
//// ---------- PLANTINGS ----------
//
//// PlantCrop godoc
//// @Summary      Посадить культуру на грядку
//// @Tags         plantings
//// @Accept       json
//// @Produce      json
//// @Security     BearerAuth
//// @Param        plot_id  path      string              true  "Plot ID"
//// @Param        body     body      PlantCropRequest    true  "Данные посадки"
//// @Success      200 {object} response.Response{data=string}
//// @Router       /api/v1/plots/{plot_id}/plant [post]
//func (c *Controller) PlantCrop(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "plot_id")
//	ownerID := middleware.SubID(r.Context())
//
//	var req PlantCropRequest
//	if err := middleware.DecodeAndValidate(r, &req); err != nil {
//		middleware.WriteValidationError(w, err)
//		return
//	}
//
//	if err := c.svc.PlantCrop(r.Context(), plotID, ownerID, req); err != nil {
//		helpers.WriteErr(w, err)
//		return
//	}
//	response.JSON(w, http.StatusOK, "crop planted successfully")
//}
//
//// HarvestCrop godoc
//// @Summary      Собрать урожай с грядки
//// @Tags         plantings
//// @Accept       json
//// @Produce      json
//// @Security     BearerAuth
//// @Param        plot_id  path      string               true  "Plot ID"
//// @Param        body     body      HarvestCropRequest   true  "Данные сбора"
//// @Success      200 {object} response.Response{data=string}
//// @Router       /api/v1/plots/{plot_id}/harvest [post]
//func (c *Controller) HarvestCrop(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "plot_id")
//	ownerID := middleware.SubID(r.Context())
//
//	var req HarvestCropRequest
//	if err := middleware.DecodeAndValidate(r, &req); err != nil {
//		middleware.WriteValidationError(w, err)
//		return
//	}
//
//	if err := c.svc.HarvestCrop(r.Context(), plotID, ownerID, req); err != nil {
//		helpers.WriteErr(w, err)
//		return
//	}
//	response.JSON(w, http.StatusOK, "harvest completed successfully")
//}
//
//// ---------- TIMELINE ----------
//
//// GetTimeline godoc
//// @Summary      Получить историю изменений участка
//// @Tags         timeline
//// @Produce      json
//// @Security     BearerAuth
//// @Param        plot_id  path      string  true   "Plot ID"
//// @Param        from     query     string  false  "Начальная дата (YYYY-MM-DD)"
//// @Param        to       query     string  false  "Конечная дата (YYYY-MM-DD)"
//// @Param        limit    query     int     false  "Максимум событий (default 100)"
//// @Success      200 {object} response.Response{data=[]TimelineEvent}
//// @Router       /api/v1/plots/{plot_id}/timeline [get]
//func (c *Controller) GetTimeline(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "plot_id")
//	ownerID := middleware.SubID(r.Context())
//
//	q := r.URL.Query()
//	filter := TimelineFilter{}
//
//	if fromStr := q.Get("from"); fromStr != "" {
//		if t, err := time.Parse("2006-01-02", fromStr); err == nil {
//			filter.From = &t
//		}
//	}
//	if toStr := q.Get("to"); toStr != "" {
//		if t, err := time.Parse("2006-01-02", toStr); err == nil {
//			filter.To = &t
//		}
//	}
//	if limitStr := q.Get("limit"); limitStr != "" {
//		if l, err := strconv.Atoi(limitStr); err == nil {
//			filter.Limit = l
//		}
//	}
//
//	events, err := c.svc.GetTimeline(r.Context(), plotID, ownerID, filter)
//	if err != nil {
//		helpers.WriteErr(w, err)
//		return
//	}
//	response.JSON(w, http.StatusOK, events)
//}
