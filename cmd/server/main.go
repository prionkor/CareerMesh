package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/prionkor/careermesh/internal/certification"
	"github.com/prionkor/careermesh/internal/database"
	"github.com/prionkor/careermesh/internal/education"
	"github.com/prionkor/careermesh/internal/experience"
	"github.com/prionkor/careermesh/internal/language"
	"github.com/prionkor/careermesh/internal/profile"
	"github.com/prionkor/careermesh/internal/project"
	"github.com/prionkor/careermesh/internal/skill"
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

	experienceRepository := experience.NewRepository(db)
	experienceService := experience.NewService(experienceRepository)
	experienceHandler := experience.NewHandler(experienceService)

	projectRepository := project.NewRepository(db)
	projectService := project.NewService(projectRepository)
	projectHandler := project.NewHandler(projectService)

	skillRepository := skill.NewRepository(db)
	skillService := skill.NewService(skillRepository)
	skillHandler := skill.NewHandler(skillService)

	educationRepository := education.NewRepository(db)
	educationService := education.NewService(educationRepository)
	educationHandler := education.NewHandler(educationService)

	certificationRepository := certification.NewRepository(db)
	certificationService := certification.NewService(certificationRepository)
	certificationHandler := certification.NewHandler(certificationService)

	languageRepository := language.NewRepository(db)
	languageService := language.NewService(languageRepository)
	languageHandler := language.NewHandler(languageService)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(
		userService,
		profileService,
		experienceService,
		projectService,
		skillService,
		educationService,
		certificationService,
		languageService,
	)

	app := fiber.New()

	v1 := app.Group("/api/v1")
	profileHandler.RegisterRoutes(v1)
	userHandler.RegisterRoutes(v1)
	experienceHandler.RegisterRoutes(v1)
	projectHandler.RegisterRoutes(v1)
	skillHandler.RegisterRoutes(v1)
	educationHandler.RegisterRoutes(v1)
	certificationHandler.RegisterRoutes(v1)
	languageHandler.RegisterRoutes(v1)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("server listening on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
