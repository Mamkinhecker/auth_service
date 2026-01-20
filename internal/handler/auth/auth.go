package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"auth_service/internal/model/request"
	authService "auth_service/internal/service/auth"
)

type Auth_handler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
}

type AuthHandler struct {
	authService *authService.AuthService
}

func NewAuthHandler(authService *authService.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// SignUp
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя в системе
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body db.SignUpRequest true "Данные для регистрации"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/auth/signup [post]

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req request.SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, tokens, err := h.authService.SignUp(ctx, req)
	if err != nil {
		ErrorResponse(w, err.Error(), http.StatusBadRequest)
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

	JsonResponse(w, response, http.StatusCreated)
}

// SignIn
// @Summary Вход в систему
// @Description Аутентификация по номеру телефона и паролю
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body db.LoginRequest true "Учетные данные"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/signin [post]
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, tokens, err := h.authService.SignIn(ctx, req)
	if err != nil {
		ErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
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

	JsonResponse(w, response, http.StatusOK)
}

// Logout
// @Summary Выход из системы
// @Description Инвалидирует токены пользователя
// @Tags Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/profile/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(int64)
	/*if !ok {
		errorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}*/

	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	ctx := r.Context()
	err := h.authService.Logout(ctx, userID, token)
	if err != nil {
		ErrorResponse(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	JsonResponse(w, map[string]interface{}{
		"success": true,
		"message": "logged out successfully",
	}, http.StatusOK)
}

// Refresh
// @Summary Обновление токенов
// @Description Получение новой пары токенов
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body db.RefreshTokenRequest true "Refresh токен"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/refresh [post]
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
		ErrorResponse(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	JsonResponse(w, map[string]interface{}{
		"success": true,
		"data":    tokens,
		"message": "tokens refreshed",
	}, http.StatusOK)
}

func ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

func JsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
