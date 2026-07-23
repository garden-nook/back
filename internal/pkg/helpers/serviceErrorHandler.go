package helpers

import (
	"errors"
	"garden-nook/internal/pkg/apperrors"
	"log/slog"
)

type ServiceErrorHandler struct {
	log *slog.Logger
}

func NewServiceErrorHandler(log *slog.Logger) *ServiceErrorHandler {
	return &ServiceErrorHandler{log: log}
}

func (h *ServiceErrorHandler) HandleError(err error, op string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, apperrors.ErrNotFound) ||
		errors.Is(err, apperrors.ErrConflict) ||
		errors.Is(err, apperrors.ErrBadRequest) {
		return err
	}
	h.log.Error("service error", slog.String("op", op), slog.Any("err", err), slog.String("source", GetLogSource()))
	return apperrors.ErrInternal
}
