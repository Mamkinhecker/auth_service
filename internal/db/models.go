// internal/db/models.go
package db

import (
	"database/sql"
	"time"
)

type User struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name"`
	PhoneNumber string         `db:"phone_number"`
	Email       sql.NullString `db:"email"`
	Password    string         `db:"password"`
	PhotoObj    sql.NullString `db:"photo_object"`
	IsDeleted   bool           `db:"is_deleted"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type SignUpRequest struct {
	Name        string `json:"name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Email       string `json:"email"`
	Password    string `json:"password" validate:"required,min=6"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UploadPhotoResponse struct {
	PhotoURL string `json:"photo_url"`
}
