package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"auth_service/internal/repository"
	"auth_service/internal/utils"
)

type contextKey string

const userIDKey contextKey = "user_id"

func AuthMiddleware(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.URL.Path == "/api/v1/auth/signup" ||
				r.URL.Path == "/api/v1/auth/signin" ||
				r.URL.Path == "/api/v1/auth/refresh" ||
				r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				errorResponse(w, "authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				errorResponse(w, "invalid authorization format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			claims, err := utils.ValidateAccessToken(token)
			if err != nil {
				errorResponse(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			blacklisted, err := tokenRepo.IsTokenBlacklisted(r.Context(), token)
			if err != nil {
				errorResponse(w, "internal server error", http.StatusInternalServerError)
				return
			}
			if blacklisted {
				errorResponse(w, "token is blacklisted", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf("[%s] %s %s %d %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			rw.status,
			duration,
		)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
