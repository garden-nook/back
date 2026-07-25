package plots

import (
	"garden-nook/internal/middleware"
	"garden-nook/internal/modules/plot/dto"
	_ "garden-nook/internal/modules/plot/model"
	"garden-nook/internal/modules/plot/service"
	"garden-nook/internal/pkg/helpers"
	"garden-nook/internal/pkg/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	plotSvc           *service.PlotService
	eventSvc          *service.EventService
	recommendationSvc *service.RecommendationService
}

func NewController(
	plotSvc *service.PlotService,
	eventSvc *service.EventService,
	recommendationSvc *service.RecommendationService,
) *Controller {
	return &Controller{
		plotSvc:           plotSvc,
		eventSvc:          eventSvc,
		recommendationSvc: recommendationSvc,
	}
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
// @Description  В качестве payload могут выступать BedCreatedRequest, BedUpdatedRequest, BedDeletedRequest, CropPlantedRequest, CropRemovedRequest
// @Tags         plots
// @Accept       json
// @Produce      json
// @Security     UserAuth
// @Param        id    path      string           true  "Plot ID"
// @Param        body  body      dto.PlotEvents   true  "События"
// @Success      204 {object} dto.BedCreatedRequest
// @Success      204 {object} dto.BedUpdatedRequest
// @Success      204 {object} dto.BedDeletedRequest
// @Success      204 {object} dto.CropPlantedRequest
// @Success      204 {object} dto.CropRemovedRequest
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

// GetBedRecommendations godoc
// @Summary      Получить рекомендации для посадки
// @Tags         plots
// @Produce      json
// @Security     UserAuth
// @Param        id  path      string  true  "Bed ID"
// @Param        search   		query     string  false  "Поиск по названию"
// @Param        limit    		query     int     false  "Максимум рекомендаций (default 10)"
// @Param        searchLimit    query     int     false  "Максимум культур при поиске (default 10)"
// @Param        disableFilters query     string  false  "Использовать ли фильтры по почве и освещённости для глобального поиска (default false)"
// @Success      200 {object} response.Response{data=dto.BedRecommendationsResponse}
// @Router       /api/v1/plots/bed/{id}/recommendation [get]
func (c *Controller) GetBedRecommendations(w http.ResponseWriter, r *http.Request) {
	bedID := chi.URLParam(r, "id")
	ownerID := middleware.SubID(r.Context())

	q := r.URL.Query()

	search := q.Get("search")

	limit := 10
	if limitStr := q.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	searchLimit := 10
	if limitStr := q.Get("searchLimit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			searchLimit = l
		}
	}
	disableFilters := false
	if disableFiltersStr := q.Get("disableFilters"); disableFiltersStr == "true" {
		disableFilters = true
	}

	recommendations, err := c.recommendationSvc.GetBedRecommendations(r.Context(), bedID, ownerID, search, limit, searchLimit, disableFilters)
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, recommendations)
}

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
