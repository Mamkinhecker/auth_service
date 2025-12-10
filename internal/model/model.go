// internal/model/models.go
package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID          int64          `db:"id" json:"id"`
	Name        string         `db:"name" json:"name" validate:"required,min=2,max=100"`
	PhoneNumber string         `db:"phone_number" json:"phone_number" validate:"required,startswith=+,min=11,max=15"`
	Email       sql.NullString `db:"email" json:"email,omitempty" validate:"omitempty,email,max=255"`
	Password    string         `db:"password" json:"-" validate:"required,min=6,max=100"`
	PhotoURL    sql.NullString `db:"photo_object" json:"photo_url,omitempty"`
	IsDeleted   bool           `db:"is_deleted" json:"-"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}

// UserResponse представляет ответ с информацией о пользователе
// @Description Информация о пользователе
type UserResponse struct {
	// Уникальный идентификатор
	// @Example 1
	ID int64 `json:"id"`

	// Имя пользователя
	// @Example Иван Иванов
	Name string `json:"name"`

	// Номер телефона
	// @Example +79161234567
	PhoneNumber string `json:"phone_number"`

	// Email адрес
	// @Example user@example.com
	Email string `json:"email,omitempty"`

	// URL фотографии профиля
	// @Example http://localhost:9000/user-photos/users/1/profile.jpg
	PhotoURL string `json:"photo_url,omitempty"`

	// Дата создания
	// @Example 2024-12-09T01:00:00Z`
	CreatedAt time.Time `json:"created_at"`

	// Дата последнего обновления
	// @Example 2024-12-09T01:30:00Z
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Name:        u.Name,
		PhoneNumber: u.PhoneNumber,
		Email:       u.Email.String,
		PhotoURL:    u.PhotoURL.String,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// LoginRequest для входа
type LoginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,startswith=+,min=11,max=15"`
	Password    string `json:"password" validate:"required,min=6,max=100"`
}

// SignUpRequest представляет запрос на регистрацию
// @Description Запрос на создание нового пользователя
type SignUpRequest struct {
	// Имя пользователя
	// @Example Иван Иванов
	Name string `json:"name" validate:"required,min=2,max=100"`

	// Номер телефона в международном формате
	// @Example +79161234567
	PhoneNumber string `json:"phone_number" validate:"required,startswith=+,min=11,max=15"`

	// Email адрес (опционально)
	// @Example user@example.com
	Email string `json:"email" validate:"omitempty,email,max=255"`

	// Пароль пользователя
	// @Example SecurePass123!
	Password string `json:"password" validate:"required,min=6,max=100"`
}

// UpdateProfileRequest для обновления профиля
type UpdateProfileRequest struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=100"`
	Email string `json:"email" validate:"omitempty,email,max=255"`
}

// @Param request body db.RefreshTokenRequest true "Refresh токен"
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UploadPhotoResponse для ответа с фото
type UploadPhotoResponse struct {
	PhotoURL string `json:"photo_url"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
}
