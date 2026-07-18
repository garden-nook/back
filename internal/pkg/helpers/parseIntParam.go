package helpers

import (
	"garden-nook/internal/pkg/apperrors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ParseIntParam(r *http.Request, name string) (int32, error) {
	v := chi.URLParam(r, name)
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, apperrors.ErrBadRequest
	}
	return int32(i), nil
}
