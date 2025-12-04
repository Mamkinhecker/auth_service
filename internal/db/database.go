package db

import (
	"auth_service/internal/config"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitPostgres() {
	cfg := config.App.Postgres

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	db, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		log.Fatalf("не удалось подключиться к бд: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("подключен к постгре")
}

func ClosePostgres() {
	if db != nil {
		db.Close()
		log.Println("отключение от постгры")
	}
}
