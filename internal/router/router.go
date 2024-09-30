package router

import (
	"net/http"

	"github.com/pi-prakhar/go-gcp-auth/internal/handlers"
	"github.com/pi-prakhar/go-gcp-auth/internal/middleware"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/auth/google/home", handlers.HandleHome)
	mux.HandleFunc("/api/v1/auth/google/login", handlers.HandleGoogleLogin)
	mux.HandleFunc("/api/v1/auth/google/callback", handlers.HandleGoogleCallback)
	mux.HandleFunc("/api/v1/auth/protected", middleware.AuthMiddleware(handlers.HandleProtected))

	return mux
}
