package crops

import (
	"errors"
	"net/http"
	"strconv"

	"garden-nook/internal/middleware"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/response"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	svc *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{svc: svc}
}

// mapError преобразует доменную ошибку в HTTP-статус.
func mapError(err error) int {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, apperrors.ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, apperrors.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, apperrors.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, apperrors.ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func writeErr(w http.ResponseWriter, err error) {
	status := mapError(err)
	msg := err.Error()
	if status == http.StatusInternalServerError {
		msg = "internal server error"
	}
	response.Error(w, status, msg)
}

func parseIntParam(r *http.Request, name string) (int32, error) {
	v := chi.URLParam(r, name)
	i, err := strconv.Atoi(v)
	if err != nil || i <= 0 {
		return 0, apperrors.ErrBadRequest
	}
	return int32(i), nil
}

// ---------- FAMILIES ----------

// ListFamilies godoc
// @Summary      Список семейств культур
// @Description  Возвращает все семейства культур из справочника
// @Tags         crops
// @Produce      json
// @Success      200 {object} response.Response{data=[]CropFamily}
// @Failure      500 {object} response.Response
// @Router       /api/v1/crop-families [get]
func (c *Controller) ListFamilies(w http.ResponseWriter, r *http.Request) {
	list, err := c.svc.ListFamilies(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, list)
}

// GetFamily godoc
// @Summary      Получить семейство по ID
// @Tags         crops
// @Produce      json
// @Param        id   path      int  true  "ID семейства"
// @Success      200 {object} response.Response{data=CropFamily}
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/crop-families/{id} [get]
func (c *Controller) GetFamily(w http.ResponseWriter, r *http.Request) {
	id, err := parseIntParam(r, "id")
	if err != nil {
		writeErr(w, err)
		return
	}
	f, err := c.svc.GetFamily(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, f)
}

// CreateFamily godoc
// @Summary      Создать семейство (admin)
// @Tags         crops-admin
// @Accept       json
// @Produce      json
// @Security     AdminAuth
// @Param        body  body      CreateFamilyRequest  true  "Данные семейства"
// @Success      201 {object} response.Response{data=CropFamily}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      409 {object} response.Response
// @Router       /api/v1/admin/crop-families [post]
func (c *Controller) CreateFamily(w http.ResponseWriter, r *http.Request) {
	var req CreateFamilyRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	f, err := c.svc.CreateFamily(r.Context(), req)
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, f)
}

// UpdateFamily godoc
// @Summary      Обновить семейство (admin)
// @Tags         crops-admin
// @Accept       json
// @Produce      json
// @Security     AdminAuth
// @Param        id    path      int                  true  "ID семейства"
// @Param        body  body      UpdateFamilyRequest  true  "Обновляемые поля"
// @Success      200 {object} response.Response{data=CropFamily}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/admin/crop-families/{id} [put]
func (c *Controller) UpdateFamily(w http.ResponseWriter, r *http.Request) {
	id, err := parseIntParam(r, "id")
	if err != nil {
		writeErr(w, err)
		return
	}
	var req UpdateFamilyRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	f, err := c.svc.UpdateFamily(r.Context(), id, req)
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, f)
}

// DeleteFamily godoc
// @Summary      Удалить семейство (admin)
// @Tags         crops-admin
// @Produce      json
// @Security     AdminAuth
// @Param        id  path      int  true  "ID семейства"
// @Success      204 "No Content"
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      409 {object} response.Response "Есть связанные культуры"
// @Router       /api/v1/admin/crop-families/{id} [delete]
func (c *Controller) DeleteFamily(w http.ResponseWriter, r *http.Request) {
	id, err := parseIntParam(r, "id")
	if err != nil {
		writeErr(w, err)
		return
	}
	if err := c.svc.DeleteFamily(r.Context(), id); err != nil {
		writeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---------- CROPS ----------

// ListCrops godoc
// @Summary      Список культур с пагинацией и фильтрами
// @Tags         crops
// @Produce      json
// @Param        family_id  query     int     false  "Фильтр по семейству"
// @Param        search     query     string  false  "Поиск по названию"
// @Param        page       query     int     false  "Номер страницы (default 1)"
// @Param        limit      query     int     false  "Элементов на странице (default 50, max 100)"
// @Success      200 {object} response.Response{data=[]Crop}
// @Router       /api/v1/crops [get]
func (c *Controller) ListCrops(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := ListCropsFilter{
		Search: q.Get("search"),
		Page:   atoiDefault(q.Get("page"), 1),
		Limit:  atoiDefault(q.Get("limit"), 50),
	}
	if fid := q.Get("family_id"); fid != "" {
		if i, err := strconv.Atoi(fid); err == nil {
			v := int32(i)
			filter.FamilyID = &v
		}
	}

	list, total, err := c.svc.ListCrops(r.Context(), filter)
	if err != nil {
		writeErr(w, err)
		return
	}
	response.Paginated(w, list, total, filter.Page, filter.Limit)
}

// GetCrop godoc
// @Summary      Получить культуру по ID
// @Tags         crops
// @Produce      json
// @Param        id   path      int  true  "ID культуры"
// @Success      200 {object} response.Response{data=Crop}
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/crops/{id} [get]
func (c *Controller) GetCrop(w http.ResponseWriter, r *http.Request) {
	id, err := parseIntParam(r, "id")
	if err != nil {
		writeErr(w, err)
		return
	}
	cr, err := c.svc.GetCrop(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, cr)
}

// CreateCrop godoc
// @Summary      Создать культуру (admin)
// @Tags         crops-admin
// @Accept       json
// @Produce      json
// @Security     AdminAuth
// @Param        body  body      CreateCropRequest  true  "Данные культуры"
// @Success      201 {object} response.Response{data=Crop}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Router       /api/v1/admin/crops [post]
func (c *Controller) CreateCrop(w http.ResponseWriter, r *http.Request) {
	var req CreateCropRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	cr, err := c.svc.CreateCrop(r.Context(), req)
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, cr)
}

// UpdateCrop godoc
// @Summary      Обновить культуру (admin)
// @Tags         crops-admin
// @Accept       json
// @Produce      json
// @Security     AdminAuth
// @Param        id    path      int                 true  "ID культуры"
// @Param        body  body      UpdateCropRequest   true  "Обновляемые поля"
// @Success      200 {object} response.Response{data=Crop}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/admin/crops/{id} [put]
func (c *Controller) UpdateCrop(w http.ResponseWriter, r *http.Request) {
	id, err := parseIntParam(r, "id")
	if err != nil {
		writeErr(w, err)
		return
	}
	var req UpdateCropRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	cr, err := c.svc.UpdateCrop(r.Context(), id, req)
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, cr)
}

// DeleteCrop godoc
// @Summary      Удалить культуру (soft delete, admin)
// @Tags         crops-admin
// @Produce      json
// @Security     AdminAuth
// @Param        id  path      int  true  "ID культуры"
// @Success      204 "No Content"
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/admin/crops/{id} [delete]
func (c *Controller) DeleteCrop(w http.ResponseWriter, r *http.Request) {
	id, err := parseIntParam(r, "id")
	if err != nil {
		writeErr(w, err)
		return
	}
	if err := c.svc.DeleteCrop(r.Context(), id); err != nil {
		writeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---------- RULES ----------

// ListRules godoc
// @Summary      Список правил совместимости (admin)
// @Tags         crops-admin
// @Produce      json
// @Security     AdminAuth
// @Success      200 {object} response.Response{data=[]CropRule}
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Router       /api/v1/admin/crop-rules [get]
func (c *Controller) ListRules(w http.ResponseWriter, r *http.Request) {
	list, err := c.svc.ListRules(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, list)
}

// CreateRule godoc
// @Summary      Создать правило (admin)
// @Tags         crops-admin
// @Accept       json
// @Produce      json
// @Security     AdminAuth
// @Param        body  body      CreateRuleRequest  true  "Правило"
// @Success      201 {object} response.Response{data=CropRule}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Router       /api/v1/admin/crop-rules [post]
func (c *Controller) CreateRule(w http.ResponseWriter, r *http.Request) {
	var req CreateRuleRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	ru, err := c.svc.CreateRule(r.Context(), req)
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, ru)
}

// DeleteRule godoc
// @Summary      Удалить правило (admin)
// @Tags         crops-admin
// @Produce      json
// @Security     AdminAuth
// @Param        id  path      int  true  "ID правила"
// @Success      204 "No Content"
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/admin/crop-rules/{id} [delete]
func (c *Controller) DeleteRule(w http.ResponseWriter, r *http.Request) {
	id, err := parseIntParam(r, "id")
	if err != nil {
		writeErr(w, err)
		return
	}
	if err := c.svc.DeleteRule(r.Context(), id); err != nil {
		writeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return def
}
