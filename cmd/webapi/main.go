package main

import (
	"log"
	"net/http"

	"github.com/Frisbon/hungrymonke/service/api/handlers"
)

func main() {
	r := handlers.NewRouter()

	srv := http.Server{
		Addr:    ":8082",
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to run server: %v", err)
	}
}
