package crops

import (
	"garden-nook/internal/pkg/helpers"
	"net/http"
	"strconv"

	"garden-nook/internal/middleware"
	"garden-nook/internal/pkg/response"
)

type Controller struct {
	svc *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{svc: svc}
}

// ---------- SOIL TYPES ----------

// ListSoilTypes godoc
// @Summary      Список типов почвы
// @Description  Возвращает все типы почв из справочника
// @Tags         soils
// @Produce      json
// @Success      200 {object} response.Response{data=[]SoilType}
// @Failure      500 {object} response.Response
// @Router       /api/v1/soil-types [get]
func (c *Controller) ListSoilTypes(w http.ResponseWriter, r *http.Request) {
	list, _, err := c.svc.ListSoilTypes(r.Context())
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, list)
}

// GetSoilType godoc
// @Summary      Получить тип почвы по ID
// @Tags         soils
// @Produce      json
// @Param        id   path      int  true  "ID типа почвы"
// @Success      200 {object} response.Response{data=SoilType}
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/soil-types/{id} [get]
func (c *Controller) GetSoilType(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ParseIntParam(r, "id")
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	f, err := c.svc.GetSoilType(r.Context(), id)
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, f)
}

// CreateSoilType godoc
// @Summary      Создать тип почвы (admin)
// @Tags         soils-admin
// @Accept       json
// @Produce      json
// @Security     AdminAuth
// @Param        body  body      CreateSoilTypeRequest  true  "Данные типа почвы"
// @Success      201 {object} response.Response{data=response.CreateUpdateIntId}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      409 {object} response.Response
// @Router       /api/v1/admin/soil-types [post]
func (c *Controller) CreateSoilType(w http.ResponseWriter, r *http.Request) {
	var req CreateSoilTypeRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	f, err := c.svc.CreateSoilType(r.Context(), req)
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, f)
}

// UpdateSoilType godoc
// @Summary      Обновить тип почвы (admin)
// @Tags         soils-admin
// @Accept       json
// @Produce      json
// @Security     AdminAuth
// @Param        id    path      int                  	true  "ID типа почвы"
// @Param        body  body      UpdateSoilTypeRequest  true  "Обновляемые поля"
// @Success      200 {object} response.Response{data=response.CreateUpdateIntId}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/admin/soil-types/{id} [put]
func (c *Controller) UpdateSoilType(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ParseIntParam(r, "id")
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	var req UpdateSoilTypeRequest
	if err = middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	f, err := c.svc.UpdateSoilType(r.Context(), id, req)
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, f)
}

// DeleteSoilType godoc
// @Summary      Удалить тип почвы (admin)
// @Tags         soils-admin
// @Produce      json
// @Security     AdminAuth
// @Param        id  path      int  true  "ID типа почвы"
// @Success      204 "No Content"
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      409 {object} response.Response "Есть связанные культуры"
// @Router       /api/v1/admin/soil-types/{id} [delete]
func (c *Controller) DeleteSoilType(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ParseIntParam(r, "id")
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	if err = c.svc.DeleteSoilType(r.Context(), id); err != nil {
		helpers.WriteErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
	list, _, err := c.svc.ListFamilies(r.Context())
	if err != nil {
		helpers.WriteErr(w, err)
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
	id, err := helpers.ParseIntParam(r, "id")
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	f, err := c.svc.GetFamily(r.Context(), id)
	if err != nil {
		helpers.WriteErr(w, err)
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
// @Success      201 {object} response.Response{data=response.CreateUpdateIntId}
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
		helpers.WriteErr(w, err)
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
// @Success      200 {object} response.Response{data=response.CreateUpdateIntId}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/admin/crop-families/{id} [put]
func (c *Controller) UpdateFamily(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ParseIntParam(r, "id")
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	var req UpdateFamilyRequest
	if err = middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	f, err := c.svc.UpdateFamily(r.Context(), id, req)
	if err != nil {
		helpers.WriteErr(w, err)
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
	id, err := helpers.ParseIntParam(r, "id")
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	if err := c.svc.DeleteFamily(r.Context(), id); err != nil {
		helpers.WriteErr(w, err)
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
// @Success      200 {object} response.Response{data=[]Crop}
// @Router       /api/v1/crops [get]
func (c *Controller) ListCrops(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := ListCropsFilter{
		Search: q.Get("search"),
	}
	if fid := q.Get("family_id"); fid != "" {
		if i, err := strconv.Atoi(fid); err == nil {
			v := int32(i)
			filter.FamilyID = &v
		}
	}

	list, _, err := c.svc.ListCrops(r.Context(), filter)
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, list)
}

// GetCrop godoc
// @Summary      Получить культуру по ID и связи с другими культурами
// @Tags         crops
// @Produce      json
// @Param        id   path      int  true  "ID культуры"
// @Success      200 {object} response.Response{data=CropExtended}
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/crops/{id} [get]
func (c *Controller) GetCrop(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ParseIntParam(r, "id")
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	cr, err := c.svc.GetCrop(r.Context(), id)
	if err != nil {
		helpers.WriteErr(w, err)
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
// @Success      201 {object} response.Response{data=response.CreateUpdateIntId}
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
		helpers.WriteErr(w, err)
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
// @Success      200 {object} response.Response{data=response.CreateUpdateIntId}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/admin/crops/{id} [put]
func (c *Controller) UpdateCrop(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ParseIntParam(r, "id")
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	var req UpdateCropRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	cr, err := c.svc.UpdateCrop(r.Context(), id, req)
	if err != nil {
		helpers.WriteErr(w, err)
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
	id, err := helpers.ParseIntParam(r, "id")
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	if err := c.svc.DeleteCrop(r.Context(), id); err != nil {
		helpers.WriteErr(w, err)
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
	list, _, err := c.svc.ListRules(r.Context())
	if err != nil {
		helpers.WriteErr(w, err)
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
// @Success      201 {object} response.Response{data=response.CreateUpdateIntId}
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
		helpers.WriteErr(w, err)
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
	id, err := helpers.ParseIntParam(r, "id")
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	if err := c.svc.DeleteRule(r.Context(), id); err != nil {
		helpers.WriteErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
