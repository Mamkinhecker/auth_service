package storage

import (
	"auth_service/internal/config"
	"auth_service/internal/storage/minio"
	"auth_service/internal/storage/postgresql"
	"auth_service/internal/storage/redis"
	"log"
)

func BuildStorage() {
	config.Init()
	log.Println("Configuration loaded")

	postgresql.InitPostgres()
	defer postgresql.ClosePostgres()

	redis.InitRedis()
	defer redis.CloseRedis()

	minio.InitMinio()

}
