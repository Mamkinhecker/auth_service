package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"auth_service/internal/db"
	"auth_service/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req db.SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, tokens, err := h.authService.SignUp(ctx, req)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id":           user.ID,
				"name":         user.Name,
				"phone_number": user.PhoneNumber,
				"email":        user.Email.String,
				"photo_url":    user.PhotoURL.String,
				"created_at":   user.CreatedAt,
			},
			"tokens": tokens,
		},
		"message": "registration successful",
	}

	jsonResponse(w, response, http.StatusCreated)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req db.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, tokens, err := h.authService.SignIn(ctx, req)
	if err != nil {
		errorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id":           user.ID,
				"name":         user.Name,
				"phone_number": user.PhoneNumber,
				"email":        user.Email.String,
				"photo_url":    user.PhotoURL.String,
			},
			"tokens": tokens,
		},
		"message": "login successful",
	}

	jsonResponse(w, response, http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		errorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	ctx := r.Context()
	err := h.authService.Logout(ctx, userID, token)
	if err != nil {
		errorResponse(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"message": "logged out successfully",
	}, http.StatusOK)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	tokens, err := h.authService.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		errorResponse(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"data":    tokens,
		"message": "tokens refreshed",
	}, http.StatusOK)
}

func errorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
