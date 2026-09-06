package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	"github.com/prionkor/careermesh/internal/database"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("load .env: %v", err)
	}

	cfg := database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./cmd/migrate <up|down|status>")
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("set goose dialect: %v", err)
	}

	switch os.Args[1] {
	case "up":
		err = goose.Up(db, "migrations")
	case "down":
		err = goose.Down(db, "migrations")
	case "status":
		err = goose.Status(db, "migrations")
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}

	if err != nil {
		log.Fatalf("migration %s failed: %v", os.Args[1], err)
	}

	log.Printf("migration %s completed", os.Args[1])
}
