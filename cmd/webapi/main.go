package main

import (
	"log"
	"net/http"

	"github.com/Frisbon/hungrymonke/service/api"
)

func main() {
	apirouter := api.New()

	srv := http.Server{
		Addr:    ":8082",
		Handler: apirouter.Handler(),
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to run server: %v", err)
	}
}
