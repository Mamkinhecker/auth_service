package handlers

import (
	"auth_service/internal/repository"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func SetupRouter(
	authHandler *AuthHandler,
	profileHandler *ProfileHandler,
	userRepo *repository.UserRepository,
	tokenRepo *repository.TokenRepository,
) *mux.Router {
	router := mux.NewRouter()

	router.Use(CORSMiddleware)
	router.Use(LoggingMiddleware)
	router.Use(AuthMiddleware(userRepo, tokenRepo))

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "OK",
			"timestamp": time.Now().Unix(),
		})
	}).Methods("GET")

	api := router.PathPrefix("/api/v1").Subrouter()

	auth := api.PathPrefix("/auth").Subrouter()
	auth.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})
	auth.HandleFunc("/signup", authHandler.SignUp).Methods("POST")
	auth.HandleFunc("/signin", authHandler.SignIn).Methods("POST")
	auth.HandleFunc("/refresh", authHandler.Refresh).Methods("POST")

	profile := api.PathPrefix("/profile").Subrouter()
	profile.HandleFunc("", profileHandler.GetProfile).Methods("GET")
	profile.HandleFunc("", profileHandler.UpdateProfile).Methods("PUT")
	profile.HandleFunc("", profileHandler.DeleteProfile).Methods("DELETE")
	profile.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	profile.HandleFunc("/photo", profileHandler.UploadPhoto).Methods("POST")

	return router
}
