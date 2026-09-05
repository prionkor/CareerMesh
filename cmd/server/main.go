package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/prionkor/careermesh/internal/database"
)

func main() {
	// Load .env if it exists. Environment variables take precedence.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("load .env: %v", err)
	}

	ctx := context.Background()

	db, err := database.NewPool(ctx, database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	log.Println("database connected")
}
