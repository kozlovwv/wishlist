package middleware

import (
	"context"
	"net/http"
	"strings"

	"wishlist/internal/application"
	"wishlist/internal/infrastructure/in/response"
)

const UserIDKey string = "user_id"

type AuthMiddleware struct {
	tokenManager application.TokenManager
}

func NewAuthMiddleware(tokenManager application.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokenManager: tokenManager}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			response.Unauthorized(w, "missing authorization header")
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(w, "invalid authorization header format")
			return
		}

		userID, err := m.tokenManager.Parse(parts[1])
		if err != nil {
			response.Unauthorized(w, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) int64 {
	userID, _ := r.Context().Value(UserIDKey).(int64)
	return userID
}