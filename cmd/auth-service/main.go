package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pi-prakhar/go-gcp-auth/internal/router"
)

func main() {
	router := router.NewRouter()

	srv := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server started at 8081")
	log.Fatal(srv.ListenAndServe())
}
