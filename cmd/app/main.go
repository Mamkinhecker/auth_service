package main

import (
	"log"

	_ "auth_service/docs"
	"auth_service/internal/storage"
)

// @title Auth Service API
// @version 1.0.0
// @description Микросервис аутентификации и управления пользователями с JWT токенами
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите "Bearer {token}" для авторизации
func main() {

	storage.RunAll()
	log.Println("Server exited properly")
}
