package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prionkor/careermesh/internal/auth"
	"github.com/prionkor/careermesh/internal/certification"
	"github.com/prionkor/careermesh/internal/education"
	"github.com/prionkor/careermesh/internal/experience"
	"github.com/prionkor/careermesh/internal/language"
	"github.com/prionkor/careermesh/internal/middleware"
	"github.com/prionkor/careermesh/internal/profile"
	"github.com/prionkor/careermesh/internal/project"
	"github.com/prionkor/careermesh/internal/skill"
	"github.com/prionkor/careermesh/internal/user"
)

func New(db *pgxpool.Pool, jwtSecret []byte) *fiber.App {

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

	authService := auth.NewService(userService, jwtSecret)
	authHandler := auth.NewHandler(authService)

	app := fiber.New()

	v1 := app.Group("/api/v1")
	public := v1
	protected := v1.Group("", middleware.RequireAuth(jwtSecret))

	authHandler.RegisterRoutes(public)

	// protected routes
	profileHandler.RegisterRoutes(protected)
	userHandler.RegisterRoutes(protected)
	experienceHandler.RegisterRoutes(protected)
	projectHandler.RegisterRoutes(protected)
	skillHandler.RegisterRoutes(protected)
	educationHandler.RegisterRoutes(protected)
	certificationHandler.RegisterRoutes(protected)
	languageHandler.RegisterRoutes(protected)

	return app
}
