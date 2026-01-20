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

	redis.InitRedis()

	minio.InitMinio()

}
