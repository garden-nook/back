package middleware

import (
	"context"
	"net/http"
	"strings"

	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/jwt"
	"garden-nook/internal/pkg/response"
)

type contextKey string

const (
	ctxSubID   contextKey = "sub_id"
	ctxSubType contextKey = "sub_type"
)

type AuthMiddleware struct {
	jwtMgr *jwt.Manager
}

func NewAuth(jwtMgr *jwt.Manager) *AuthMiddleware {
	return &AuthMiddleware{jwtMgr: jwtMgr}
}

// RequireUser допускает только токены обычных пользователей.
func (m *AuthMiddleware) RequireUser(next http.Handler) http.Handler {
	return m.requireSubject(jwt.SubjectUser, next)
}

// RequireAdmin допускает только токены администраторов.
func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return m.requireSubject(jwt.SubjectAdmin, next)
}

// requireSubject — общая логика проверки токена + типа субъекта.
func (m *AuthMiddleware) requireSubject(expected jwt.SubjectType, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Error(w, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := m.jwtMgr.Parse(tokenStr)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// Ключевая проверка: тип субъекта должен совпадать
		if claims.SubType != expected {
			response.Error(w, http.StatusForbidden, "access denied: wrong token type")
			return
		}

		ctx := context.WithValue(r.Context(), ctxSubID, claims.SubID)
		ctx = context.WithValue(ctx, ctxSubType, claims.SubType)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Хелперы для извлечения из контекста
func SubID(ctx context.Context) string {
	v, _ := ctx.Value(ctxSubID).(string)
	return v
}

func SubType(ctx context.Context) jwt.SubjectType {
	v, _ := ctx.Value(ctxSubType).(jwt.SubjectType)
	return v
}
