package helpers

import (
	"errors"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/response"
	"net/http"
)

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

func WriteErr(w http.ResponseWriter, err error) {
	status := mapError(err)
	msg := err.Error()
	if status == http.StatusInternalServerError {
		msg = "internal server error"
	}
	response.Error(w, status, msg)
}
