package request

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
