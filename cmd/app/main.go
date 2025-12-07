package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth_service/internal/config"
	"auth_service/internal/db"
	"auth_service/internal/handlers"
	"auth_service/internal/repository"
	"auth_service/internal/services"
)

func main() {
	config.Init()
	log.Println("Configuration loaded")

	db.InitPostgres()
	defer db.ClosePostgres()

	db.InitRedis()
	defer db.CloseRedis()

	db.InitMinio()

	userRepo := repository.NewUserRepository(db.DB)
	tokenRepo := repository.NewTokenRepository(db.RedisClient)

	authService := services.NewAuthService(userRepo, tokenRepo)
	profileService := services.NewProfileService(userRepo, tokenRepo)

	authHandler := handlers.NewAuthHandler(authService)
	profileHandler := handlers.NewProfileHandler(profileService)

	router := handlers.SetupRouter(authHandler, profileHandler, userRepo, tokenRepo)

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

	log.Println("Server exited properly")
}
