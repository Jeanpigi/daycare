package middleware

import (
	"net/http"

	"daycare/internal/domain"
)

type RequireAdminMiddleware struct{}

func NewRequireAdminMiddleware() *RequireAdminMiddleware {
	return &RequireAdminMiddleware{}
}

func (m *RequireAdminMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if Role(r.Context()) != "ADMIN" {
			_ = domain.ErrForbidden
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
