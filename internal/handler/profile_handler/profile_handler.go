package profile_handler

import (
	"encoding/json"
	"net/http"

	"auth_service/internal/handler/auth"
	"auth_service/internal/handler/middleware"
	"auth_service/internal/model/request"
	profileService "auth_service/internal/service/profile"
)

type ProfileHandler struct {
	profileService *profileService.ProfileService
}

func NewProfileHandler(profileService *profileService.ProfileService) *ProfileHandler {
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
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		auth.ErrorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	user, err := h.profileService.GetProfile(ctx, userID)
	if err != nil {
		auth.ErrorResponse(w, "profile not found", http.StatusNotFound)
		return
	}

	auth.JsonResponse(w, map[string]interface{}{
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
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		auth.ErrorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		auth.ErrorResponse(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, err := h.profileService.UpdateProfile(ctx, userID, request.UpdateProfileRequest{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		auth.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	auth.JsonResponse(w, map[string]interface{}{
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
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		auth.ErrorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	err := h.profileService.DeleteProfile(ctx, userID)
	if err != nil {
		auth.ErrorResponse(w, "failed to delete profile", http.StatusInternalServerError)
		return
	}

	auth.JsonResponse(w, map[string]interface{}{
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
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		auth.ErrorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		auth.ErrorResponse(w, "file too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		auth.ErrorResponse(w, "no photo uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/jpg" {
		auth.ErrorResponse(w, "only JPEG and PNG images are allowed", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	photoURL, err := h.profileService.UploadPhoto(ctx, userID, file, header.Filename, header.Size)
	if err != nil {
		auth.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	auth.JsonResponse(w, map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"photo_url": photoURL,
		},
		"message": "photo uploaded successfully",
	}, http.StatusOK)
}
