package db

import (
	"auth_service/internal/config"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitPostgres() {
	cfg := config.App.Postgres

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	var err error
	DB, err = sqlx.Connect("postgres", connStr)

	if err != nil {
		log.Fatalf("не удалось подключиться к бд: %v", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("подключен к постгре")
}

func ClosePostgres() {
	if DB != nil {
		DB.Close()
		log.Println("отключение от постгры")
	}
}
