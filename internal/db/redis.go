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
		log.Fatalf("failed ro connect redis: %v", err)
	}

	log.Println("connected with redis successfuly")
}

func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
		log.Println("closing redis")
	}
}
