package main

import (
	"log"
	"net/http"

	"auth_service/internal/config"
	"auth_service/internal/db"
	"auth_service/internal/handlers"
)

func main() {
	// Инициализация конфигурации
	config.Init()

	// Инициализация БД
	db.InitPostgres()
	defer db.ClosePostgres()

	db.InitRedis()
	defer db.CloseRedis()

	// Инициализация MinIO
	db.InitMinio()

	// Создание роутера
	router := handlers.SetupRouter()

	// Запуск сервера
	port := config.App.Server.Port
	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
