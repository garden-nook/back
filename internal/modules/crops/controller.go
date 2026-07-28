package crops

import (
	"garden-nook/internal/modules/crops/dto"
	_ "garden-nook/internal/modules/crops/model"
	"garden-nook/internal/modules/crops/service"
	"garden-nook/internal/pkg/helpers"
	"net/http"
	"strconv"

	"garden-nook/internal/middleware"
	"garden-nook/internal/pkg/response"
)

type Controller struct {
	soilTypeSvc   *service.SoilTypeService
	cropFamilySvc *service.CropFamilyService
	cropSvc       *service.CropService
	cropRuleSvc   *service.CropRuleService
}

func NewController(
	soilTypeSvc *service.SoilTypeService,
	cropFamilySvc *service.CropFamilyService,
	cropSvc *service.CropService,
	cropRuleSvc *service.CropRuleService,
) *Controller {
	return &Controller{
		soilTypeSvc:   soilTypeSvc,
		cropFamilySvc: cropFamilySvc,
		cropSvc:       cropSvc,
		cropRuleSvc:   cropRuleSvc,
	}
}

// ---------- SOIL TYPES ----------

// ListSoilTypes godoc
// @Summary      Список типов почвы
// @Description  Возвращает все типы почв из справочника
// @Tags         soils
// @Produce      json
// @Success      200 {object} response.Response{data=[]model.SoilType}
// @Failure      500 {object} response.Response
// @Router       /api/v1/soil-types [get]
func (c *Controller) ListSoilTypes(w http.ResponseWriter, r *http.Request) {
	list, _, err := c.soilTypeSvc.ListSoilTypes(r.Context())
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
// @Success      200 {object} response.Response{data=model.SoilType}
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/soil-types/{id} [get]
func (c *Controller) GetSoilType(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ParseIntParam(r, "id")
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	f, err := c.soilTypeSvc.GetSoilType(r.Context(), id)
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
// @Param        body  body      dto.CreateSoilTypeRequest  true  "Данные типа почвы"
// @Success      201 {object} response.Response{data=response.CreateUpdateIntId}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      409 {object} response.Response
// @Router       /api/v1/admin/soil-types [post]
func (c *Controller) CreateSoilType(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSoilTypeRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	f, err := c.soilTypeSvc.CreateSoilType(r.Context(), req)
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
// @Param        body  body      dto.UpdateSoilTypeRequest  true  "Обновляемые поля"
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
	var req dto.UpdateSoilTypeRequest
	if err = middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	f, err := c.soilTypeSvc.UpdateSoilType(r.Context(), id, req)
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
	if err = c.soilTypeSvc.DeleteSoilType(r.Context(), id); err != nil {
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
// @Success      200 {object} response.Response{data=[]model.CropFamily}
// @Failure      500 {object} response.Response
// @Router       /api/v1/crop-families [get]
func (c *Controller) ListFamilies(w http.ResponseWriter, r *http.Request) {
	list, _, err := c.cropFamilySvc.ListFamilies(r.Context())
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
// @Success      200 {object} response.Response{data=model.CropFamily}
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/crop-families/{id} [get]
func (c *Controller) GetFamily(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ParseIntParam(r, "id")
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	f, err := c.cropFamilySvc.GetFamily(r.Context(), id)
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
// @Param        body  body      dto.CreateFamilyRequest  true  "Данные семейства"
// @Success      201 {object} response.Response{data=response.CreateUpdateIntId}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      409 {object} response.Response
// @Router       /api/v1/admin/crop-families [post]
func (c *Controller) CreateFamily(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFamilyRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	f, err := c.cropFamilySvc.CreateFamily(r.Context(), req)
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
// @Param        body  body      dto.UpdateFamilyRequest  true  "Обновляемые поля"
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
	var req dto.UpdateFamilyRequest
	if err = middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	f, err := c.cropFamilySvc.UpdateFamily(r.Context(), id, req)
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
	if err = c.cropFamilySvc.DeleteFamily(r.Context(), id); err != nil {
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
// @Success      200 {object} response.Response{data=[]model.Crop}
// @Router       /api/v1/crops [get]
func (c *Controller) ListCrops(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := dto.ListCropsFilter{
		Search: q.Get("search"),
	}
	if fid := q.Get("family_id"); fid != "" {
		if i, err := strconv.Atoi(fid); err == nil {
			v := int32(i)
			filter.FamilyID = &v
		}
	}

	list, _, err := c.cropSvc.ListCrops(r.Context(), filter)
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
// @Success      200 {object} response.Response{data=dto.CropExtended}
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/crops/{id} [get]
func (c *Controller) GetCrop(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ParseIntParam(r, "id")
	if err != nil {
		helpers.WriteErr(w, err)
		return
	}
	cr, err := c.cropSvc.GetCrop(r.Context(), id)
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
// @Param        body  body      dto.CreateCropRequest  true  "Данные культуры"
// @Success      201 {object} response.Response{data=response.CreateUpdateIntId}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Router       /api/v1/admin/crops [post]
func (c *Controller) CreateCrop(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCropRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	cr, err := c.cropSvc.CreateCrop(r.Context(), req)
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
// @Param        body  body      dto.UpdateCropRequest   true  "Обновляемые поля"
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
	var req dto.UpdateCropRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	cr, err := c.cropSvc.UpdateCrop(r.Context(), id, req)
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
	if err := c.cropSvc.DeleteCrop(r.Context(), id); err != nil {
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
// @Success      200 {object} response.Response{data=[]model.CropRule}
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Router       /api/v1/admin/crop-rules [get]
func (c *Controller) ListRules(w http.ResponseWriter, r *http.Request) {
	list, _, err := c.cropRuleSvc.ListRules(r.Context())
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
// @Param        body  body      dto.CreateRuleRequest  true  "Правило"
// @Success      201 {object} response.Response{data=response.CreateUpdateIntId}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Router       /api/v1/admin/crop-rules [post]
func (c *Controller) CreateRule(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRuleRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	ru, err := c.cropRuleSvc.CreateRule(r.Context(), req)
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
	if err = c.cropRuleSvc.DeleteRule(r.Context(), id); err != nil {
		helpers.WriteErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
