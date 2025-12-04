package db

import (
	"auth_service/internal/config"
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	cfg := config.App.Redis

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("не получается подключится к редису: %v", err)
	}

	log.Println("подключился к редис успешно")
}

func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
		log.Println("отключился от редиса")
	}
}
