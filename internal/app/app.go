package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prionkor/careermesh/internal/certification"
	"github.com/prionkor/careermesh/internal/education"
	"github.com/prionkor/careermesh/internal/experience"
	"github.com/prionkor/careermesh/internal/language"
	"github.com/prionkor/careermesh/internal/profile"
	"github.com/prionkor/careermesh/internal/project"
	"github.com/prionkor/careermesh/internal/skill"
	"github.com/prionkor/careermesh/internal/user"
)

func New(db *pgxpool.Pool) *fiber.App {

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

	return app
}
