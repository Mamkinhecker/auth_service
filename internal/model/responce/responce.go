package responce

import (
	"time"
)

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
