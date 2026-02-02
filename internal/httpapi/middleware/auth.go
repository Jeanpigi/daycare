package middleware

import (
	"context"
	"net/http"
	"strings"

	"daycare/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	ctxUserIDKey ctxKey = "user_id"
	ctxRoleKey   ctxKey = "role"
)

type AuthMiddleware struct {
	secret []byte
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{secret: []byte(secret)}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")

		tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, domain.ErrUnauthorized
			}
			return m.secret, nil
		})
		if err != nil || !tok.Valid {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		sub, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		role, _ := claims["role"].(string)

		ctx := context.WithValue(r.Context(), ctxUserIDKey, uint64(sub))
		ctx = context.WithValue(ctx, ctxRoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserID(ctx context.Context) uint64 {
	v := ctx.Value(ctxUserIDKey)
	if v == nil {
		return 0
	}
	return v.(uint64)
}

func Role(ctx context.Context) string {
	v := ctx.Value(ctxRoleKey)
	if v == nil {
		return ""
	}
	return v.(string)
}
