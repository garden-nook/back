package middleware

import (
	"encoding/json"
	"errors"
	"garden-nook/internal/pkg/apperrors"
	"net/http"
	"strings"

	"garden-nook/internal/pkg/response"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// DecodeAndValidate читает JSON-body, декодирует в dst и валидирует по тегам.
// Используется внутри контроллеров (не как http-middleware) для гибкости.
func DecodeAndValidate(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return &ValidationError{Message: "empty request body"}
	}
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return &ValidationError{Message: "invalid json: " + err.Error()}
	}

	if err := validate.Struct(dst); err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return apperrors.ErrBadRequest
		}
		return &ValidationError{Message: formatValidationErrors(err)}
	}
	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

// WriteValidationError отвечает клиенту 400 с описанием.
func WriteValidationError(w http.ResponseWriter, err error) {
	msg := "validation failed"
	var ve *ValidationError
	if errors.As(err, &ve) {
		msg = ve.Message
	}
	response.Error(w, http.StatusBadRequest, msg)
}

func formatValidationErrors(err error) string {
	var parts []string
	for _, e := range err.(validator.ValidationErrors) {
		parts = append(parts, e.Field()+" "+ruleToMessage(e))
	}
	return strings.Join(parts, "; ")
}

func ruleToMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "min":
		return "min length " + e.Param()
	case "max":
		return "max length " + e.Param()
	case "email":
		return "must be a valid email"
	case "gte":
		return "must be >= " + e.Param()
	case "lte":
		return "must be <= " + e.Param()
	case "oneof":
		return "must be one of: " + e.Param()
	default:
		return "is invalid (" + e.Tag() + ")"
	}
}
