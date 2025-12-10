package storage

import (
	"auth_service/internal/config"
	"auth_service/internal/handler/auth"
	"auth_service/internal/handler/profile_handler"
	tokenrepo "auth_service/internal/repository/token"
	userrepo "auth_service/internal/repository/user"
	"auth_service/internal/router"
	authService "auth_service/internal/service/auth"
	profileService "auth_service/internal/service/profile"
	"auth_service/internal/storage/minio"
	"auth_service/internal/storage/postgresql"
	"auth_service/internal/storage/redis"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

func RunAll() {
	config.Init()
	log.Println("Configuration loaded")

	postgresql.InitPostgres()
	defer postgresql.ClosePostgres()

	redis.InitRedis()
	defer redis.CloseRedis()

	minio.InitMinio()

	userRepo := userrepo.NewUserRepository(postgresql.DB)
	tokenRepo := tokenrepo.NewTokenRepository(redis.RedisClient)

	authService := authService.NewAuthService(userRepo, tokenRepo)
	profileService := profileService.NewProfileService(userRepo, tokenRepo)

	authHandler := auth.NewAuthHandler(authService)
	profileHandler := profile_handler.NewProfileHandler(profileService)

	router := router.SetupRouter(authHandler, profileHandler, userRepo, tokenRepo)
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:         ":" + config.App.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", config.App.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
