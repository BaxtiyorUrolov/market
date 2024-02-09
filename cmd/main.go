package main

import (
	"context"
	"fmt"
	"log"
	"market/api"
	"market/config"
	"market/storage/postgres"
)

func main() {
	cfg := config.Load()

	store, err := postgres.New(context.Background(), cfg)
	if err != nil {
		log.Fatalf("error while connecting to db: %v", err)
	}
	defer store.Close()

	server := api.New(store)

	if err := server.Run("localhost:8080"); err != nil {
		fmt.Printf("error while running server: %v\n", err)
	}
}
