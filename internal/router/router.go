package router

import (
	"auth_service/internal/handler/auth"
	"auth_service/internal/handler/middleware"
	"auth_service/internal/handler/profile_handler"
	tokenrepo "auth_service/internal/repository/token"
	userrepo "auth_service/internal/repository/user"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func SetupRouter(
	authHandler *auth.AuthHandler,
	profileHandler *profile_handler.ProfileHandler,
	userRepo *userrepo.UserRepository,
	tokenRepo *tokenrepo.TokenRepository,
) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.AuthMiddleware(userRepo, tokenRepo))

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
