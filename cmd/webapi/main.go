package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/Frisbon/hungrymonke/service/api/handlers"
)

func main() {
	r := handlers.NewRouter()

	// Read PORT from env; default 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("API listening on :%s", port)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to run server: %v", err)
	}
}
