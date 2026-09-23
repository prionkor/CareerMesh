package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/prionkor/careermesh/internal/database"
	"github.com/prionkor/careermesh/internal/profile"
	"github.com/prionkor/careermesh/internal/user"
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

	profileRepository := profile.NewRepository(db)
	profileService := profile.NewService(profileRepository)
	profileHandler := profile.NewHandler(profileService)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService, profileService)

	app := fiber.New()

	v1 := app.Group("/api/v1")
	profileHandler.RegisterRoutes(v1)
	userHandler.RegisterRoutes(v1)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("server listening on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
