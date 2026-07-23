package plots

import (
	"garden-nook/internal/middleware"
	"garden-nook/internal/modules/plot/dto"
	_ "garden-nook/internal/modules/plot/model"
	"garden-nook/internal/modules/plot/service"
	"garden-nook/internal/pkg/helpers"
	"garden-nook/internal/pkg/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	plotSvc  *service.PlotService
	eventSvc *service.EventService
}

func NewController(plotSvc *service.PlotService, eventSvc *service.EventService) *Controller {
	return &Controller{plotSvc: plotSvc, eventSvc: eventSvc}
}

// ListPlots godoc
// @Summary      Список участков пользователя
// @Tags         plots
// @Produce      json
// @Security     UserAuth
// @Success      200 {object} response.Response{data=[]model.Plot}
// @Router       /api/v1/plots [get]
func (c *Controller) ListPlots(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.SubID(r.Context())
	plots, _, err := c.plotSvc.ListPlots(r.Context(), ownerID, nil)
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
// @Param        body  body   dto.CreatePlotRequest  true  "Данные участка"
// @Success      201 {object} response.Response{data=response.CreateUpdateUuidId}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/plots [post]
func (c *Controller) CreatePlot(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePlotRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}

	ownerID := middleware.SubID(r.Context())
	plot, err := c.plotSvc.CreatePlot(r.Context(), ownerID, req)
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
// @Param        id    path      string                true  "ID участка"
// @Param        body  body      dto.UpdatePlotRequest  true  "Данные участка"
// @Success      200 {object} response.Response{data=response.CreateUpdateUuidId}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/plots/{id} [put]
func (c *Controller) UpdatePlot(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdatePlotRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}

	plotID := chi.URLParam(r, "id")
	ownerID := middleware.SubID(r.Context())
	plot, err := c.plotSvc.UpdatePlot(r.Context(), plotID, ownerID, req)
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

	if err := c.plotSvc.DeletePlot(r.Context(), plotID, ownerID); err != nil {
		helpers.WriteErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetPlotStructure godoc
// @Summary      Получить информацию о участке и его структуре
// @Tags         plots
// @Produce      json
// @Security     UserAuth
// @Param        id  path      string  true  "Plot ID"
// @Success      200 {object} response.Response{data=model.PlotStructure}
// @Router       /api/v1/plots/{id}/structure [get]
func (c *Controller) GetPlotStructure(w http.ResponseWriter, r *http.Request) {
	plotID := chi.URLParam(r, "id")
	ownerID := middleware.SubID(r.Context())

	structure, err := c.plotSvc.GetPlotStructure(r.Context(), plotID, ownerID)
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, structure)
}

// HandleEvents godoc
// @Summary      Отправить события действий над участком
// @Description  В качестве payload могут выступать BedCreatedRequest, BedDeletedRequest
// @Tags         plots
// @Accept       json
// @Produce      json
// @Security     UserAuth
// @Param        id    path      string           true  "Plot ID"
// @Param        body  body      dto.PlotEvents   true  "События"
// @Success      204 {object} dto.BedCreatedRequest
// @Success      204 {object} dto.BedDeletedRequest
// @Success      204 "No Content"
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/plots/{id}/events [post]
func (c *Controller) HandleEvents(w http.ResponseWriter, r *http.Request) {
	plotID := chi.URLParam(r, "id")
	ownerID := middleware.SubID(r.Context())

	var req dto.PlotEvents
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}

	if err := c.eventSvc.HandleEvents(r.Context(), plotID, ownerID, req.Events); err != nil {
		helpers.WriteErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

//
//// ListBeds godoc
//// @Summary      Список грядок участка
//// @Tags         beds
//// @Produce      json
//// @Security     UserAuth
//// @Param        plot_id  path      string  true  "Plot ID"
//// @Success      200 {object} response.Response{data=[]Bed}
//// @Router       /api/v1/plots/{plot_id}/beds [get]
//func (c *Controller) ListBeds(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "plot_id")
//	ownerID := middleware.SubID(r.Context())
//
//	beds, err := c.plotSvc.ListBeds(r.Context(), plotID, ownerID)
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
//// @Security     UserAuth
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
//	bed, err := c.plotSvc.UpdateBed(r.Context(), plotID, bedID, ownerID, req)
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
//// @Security     UserAuth
//// @Param        plot_id  path      string  true  "Plot ID"
//// @Param        bed_id   path      string  true  "Bed ID"
//// @Success      204 "No Content"
//// @Router       /api/v1/plots/{plot_id}/beds/{bed_id} [delete]
//func (c *Controller) DeleteBed(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "plot_id")
//	bedID := chi.URLParam(r, "bed_id")
//	ownerID := middleware.SubID(r.Context())
//
//	if err := c.plotSvc.DeleteBed(r.Context(), plotID, bedID, ownerID); err != nil {
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
//// @Security     UserAuth
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
//	obj, err := c.plotSvc.CreateObject(r.Context(), plotID, ownerID, req)
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
//// @Security     UserAuth
//// @Param        plot_id  path      string  true  "Plot ID"
//// @Success      200 {object} response.Response{data=[]UIObject}
//// @Router       /api/v1/plots/{plot_id}/objects [get]
//func (c *Controller) ListObjects(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "plot_id")
//	ownerID := middleware.SubID(r.Context())
//
//	objects, err := c.plotSvc.ListObjects(r.Context(), plotID, ownerID)
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
//// @Security     UserAuth
//// @Param        plot_id    path      string  true  "Plot ID"
//// @Param        object_id  path      string  true  "Object ID"
//// @Success      204 "No Content"
//// @Router       /api/v1/plots/{plot_id}/objects/{object_id} [delete]
//func (c *Controller) DeleteObject(w http.ResponseWriter, r *http.Request) {
//	plotID := chi.URLParam(r, "plot_id")
//	objectID := chi.URLParam(r, "object_id")
//	ownerID := middleware.SubID(r.Context())
//
//	if err := c.plotSvc.DeleteObject(r.Context(), plotID, objectID, ownerID); err != nil {
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
//// @Security     UserAuth
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
//	if err := c.plotSvc.PlantCrop(r.Context(), plotID, ownerID, req); err != nil {
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
//// @Security     UserAuth
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
//	if err := c.plotSvc.HarvestCrop(r.Context(), plotID, ownerID, req); err != nil {
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
//// @Security     UserAuth
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
//	events, err := c.plotSvc.GetTimeline(r.Context(), plotID, ownerID, filter)
//	if err != nil {
//		helpers.WriteErr(w, err)
//		return
//	}
//	response.JSON(w, http.StatusOK, events)
//}
