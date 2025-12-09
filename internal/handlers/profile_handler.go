package handlers

import (
	"encoding/json"
	"net/http"

	"auth_service/internal/db"
	"auth_service/internal/services"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

// GetProfile
// @Summary Получение профиля пользователя
// @Description Возвращает информацию о пользователе
// @Tags Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/profile [get]
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		errorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	user, err := h.profileService.GetProfile(ctx, userID)
	if err != nil {
		errorResponse(w, "profile not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"data":    user.ToResponse(),
	}, http.StatusOK)
}

// UpdateProfile
// @Summary Обновление профиля
// @Description Обновляет данные пользователя
// @Tags Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body db.UpdateProfileRequest true "Новые данные"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/profile [put]
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		errorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, err := h.profileService.UpdateProfile(ctx, userID, db.UpdateProfileRequest{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"data":    user.ToResponse(),
		"message": "profile updated successfully",
	}, http.StatusOK)
}

// DeleteProfile
// @Summary Удаление профиля
// @Description Мягкое удаление профиля
// @Tags Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/profile [delete]
func (h *ProfileHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		errorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	err := h.profileService.DeleteProfile(ctx, userID)
	if err != nil {
		errorResponse(w, "failed to delete profile", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"message": "profile deleted successfully",
	}, http.StatusOK)
}

// UploadPhoto
// @Summary Загрузка фото профиля
// @Description Загружает изображение для профиля
// @Tags Profile
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param photo formData file true "Файл изображения"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/profile/photo [post]
func (h *ProfileHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		errorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		errorResponse(w, "file too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		errorResponse(w, "no photo uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/jpg" {
		errorResponse(w, "only JPEG and PNG images are allowed", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	photoURL, err := h.profileService.UploadPhoto(ctx, userID, file, header.Filename, header.Size)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"photo_url": photoURL,
		},
		"message": "photo uploaded successfully",
	}, http.StatusOK)
}
