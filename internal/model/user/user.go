package user

import (
	"auth_service/internal/model/responce"
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

func (u *User) ToResponse() responce.UserResponse {
	return responce.UserResponse{
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
