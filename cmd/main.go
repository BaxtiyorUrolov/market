package main

import (
	"context"
	"fmt"
	"log"
	"market/api"
	"market/config"
	"market/service"
	"market/storage/postgres"
)

func main() {
	cfg := config.Load()

	store, err := postgres.New(context.Background(), cfg)
	if err != nil {
		log.Fatalf("error while connecting to db: %v", err)
	}
	defer store.Close()

	services := service.New(store)

	server := api.New(services, store)

	if err := server.Run("localhost:8080"); err != nil {
		fmt.Printf("error while running server: %v\n", err)
	}
}
