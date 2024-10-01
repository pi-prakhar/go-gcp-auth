package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pi-prakhar/go-gcp-auth/internal/handlers"
	"github.com/pi-prakhar/go-gcp-auth/internal/router"
	"github.com/pi-prakhar/go-gcp-auth/internal/services"
)

func main() {

	authService := services.NewGoogleAuthService()
	authHandler := handlers.NewAuthHandler(authService)
	router := router.NewRouter(authHandler)

	srv := http.Server{
		Addr:    ":8081",
		Handler: router.Mux,
	}

	fmt.Println("Server started at 8081")
	log.Fatal(srv.ListenAndServe())
}
