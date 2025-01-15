package main

import (
	"context"
	"log"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"

	"tcp-aws-crud/config"
	"tcp-aws-crud/internal/db"
	"tcp-aws-crud/internal/server"
)

func main() {
	ctx := context.Background()

	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	dynamoDB, err := db.New(ctx, cfg.AWS)
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB: %v", err)
	}

	srv, err := server.NewServer(ctx, cfg.Server, dynamoDB)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	if err = srv.Run(ctx); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
