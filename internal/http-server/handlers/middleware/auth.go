package middleware

import (
	"WarehouseControl/internal/lib/api/response"
	"WarehouseControl/internal/models"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(secret string, log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Warn("authorization header is missing")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response.Error("missing authorization header"))
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				log.Warn("bearer token is missing")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response.Error("invalid authorization header format"))
				return
			}

			claims := &models.JWTClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				log.Warn("invalid token", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response.Error("invalid token"))
				return
			}

			// Сохраняем в контекст с правильным ключом
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) (*models.JWTClaims, bool) {
	user, ok := ctx.Value(UserContextKey).(*models.JWTClaims)
	return user, ok
}
