package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	pkgauth "SyncNote/pkg/auth"
)

const CtxUserIDKey = "currentUserID"

type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
			writeUnauthorized(w, "missing or invalid authorization header")
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
		if token == "" {
			writeUnauthorized(w, "empty token")
			return
		}

		claims, err := pkgauth.ParseToken(token)
		if err != nil {
			writeUnauthorized(w, "invalid token")
			return
		}
		if claims.UserID == "" {
			writeUnauthorized(w, "token user_id is empty")
			return
		}

		ctx := context.WithValue(r.Context(), CtxUserIDKey, claims.UserID)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code":    http.StatusUnauthorized,
		"message": msg,
	})
}
