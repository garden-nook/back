package auth

import (
	"errors"
	"garden-nook/internal/middleware"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/response"
	"net/http"
)

type Controller struct {
	svc *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{svc: svc}
}

func writeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, apperrors.ErrUnauthorized):
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, apperrors.ErrForbidden):
		response.Error(w, http.StatusForbidden, err.Error())
	case errors.Is(err, apperrors.ErrConflict):
		response.Error(w, http.StatusConflict, "email already registered")
	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}

// RegisterUser godoc
// @Summary      Регистрация пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterRequest  true  "Данные регистрации"
// @Success      201 {object} response.Response{data=User}
// @Failure      400 {object} response.Response
// @Failure      409 {object} response.Response
// @Router       /api/v1/auth/register [post]
func (c *Controller) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	u, err := c.svc.RegisterUser(r.Context(), req)
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, u)
}

// LoginUser godoc
// @Summary      Вход пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      LoginRequest  true  "Email и пароль"
// @Success      200 {object} response.Response{data=TokenResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Router       /api/v1/auth/login [post]
func (c *Controller) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	tok, err := c.svc.LoginUser(r.Context(), req)
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, tok)
}

// MeUser godoc
// @Summary      Профиль текущего пользователя
// @Tags         auth
// @Produce      json
// @Security     UserAuth
// @Success      200 {object} response.Response{data=MeResponse}
// @Failure      401 {object} response.Response
// @Router       /api/v1/auth/me [get]
func (c *Controller) MeUser(w http.ResponseWriter, r *http.Request) {
	userID := middleware.SubID(r.Context())
	me, err := c.svc.GetUserProfile(r.Context(), userID)
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, me)
}

// LoginAdmin godoc
// @Summary      Вход администратора
// @Description  Отдельный эндпоинт для получения admin-токена
// @Tags         auth-admin
// @Accept       json
// @Produce      json
// @Param        body  body      LoginRequest  true  "Email и пароль админа"
// @Success      200 {object} response.Response{data=TokenResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Router       /api/v1/admin/auth/login [post]
func (c *Controller) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := middleware.DecodeAndValidate(r, &req); err != nil {
		middleware.WriteValidationError(w, err)
		return
	}
	tok, err := c.svc.LoginAdmin(r.Context(), req)
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, tok)
}

// MeAdmin godoc
// @Summary      Профиль текущего администратора
// @Tags         auth-admin
// @Produce      json
// @Security     AdminAuth
// @Success      200 {object} response.Response{data=MeResponse}
// @Failure      401 {object} response.Response
// @Router       /api/v1/admin/auth/me [get]
func (c *Controller) MeAdmin(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.SubID(r.Context())
	me, err := c.svc.GetAdminProfile(r.Context(), adminID)
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, me)
}
