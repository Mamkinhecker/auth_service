package postgresql

import (
	"auth_service/internal/config"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitPostgres() {
	cfg := config.App.Postgres

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	fmt.Println(connStr)
	var err error
	DB, err = sqlx.Connect("pgx", connStr)

	if err != nil {
		log.Fatalf("fialed to connect to db: %v", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("connected to db")
}

func ClosePostgres() {
	if DB != nil {
		DB.Close()
		log.Println("closing db connection")
	}
}
