package user

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/prionkor/careermesh/internal/certification"
	"github.com/prionkor/careermesh/internal/education"
	"github.com/prionkor/careermesh/internal/experience"
	"github.com/prionkor/careermesh/internal/httpapi"
	"github.com/prionkor/careermesh/internal/language"
	"github.com/prionkor/careermesh/internal/profile"
	"github.com/prionkor/careermesh/internal/project"
	"github.com/prionkor/careermesh/internal/skill"
)

type Handler struct {
	service              *Service
	profileService       *profile.Service
	experienceService    *experience.Service
	projectService       *project.Service
	skillService         *skill.Service
	educationService     *education.Service
	certificationService *certification.Service
	languageService      *language.Service
}

// NewHandler creates a User HTTP handler.
func NewHandler(
	service *Service,
	profileService *profile.Service,
	experienceService *experience.Service,
	projectService *project.Service,
	skillService *skill.Service,
	educationService *education.Service,
	certificationService *certification.Service,
	languageService *language.Service,
) *Handler {
	return &Handler{
		service:              service,
		profileService:       profileService,
		experienceService:    experienceService,
		projectService:       projectService,
		skillService:         skillService,
		educationService:     educationService,
		certificationService: certificationService,
		languageService:      languageService,
	}
}

// CreateUserRequest is the client-provided body for creating a user.
type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUserRequest is the client-provided body for updating account fields.
// It intentionally has no password field; password changes are a separate operation.
type UpdateUserRequest struct {
	Email string `json:"email"`
}

// GetByID handles GET /api/v1/users/:id.
func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	usr, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get user by id", err)
	}

	return c.Status(fiber.StatusOK).JSON(usr)
}

// Create handles POST /api/v1/users.
func (h *Handler) Create(c fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	created, err := h.service.Create(c.Context(), req.Email, req.Password)
	if err != nil {
		return handleServiceError(c, "create user", err)
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// Update handles PUT /api/v1/users/:id. It updates account fields only and
// does not change the password.
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	var req UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	updated, err := h.service.Update(c.Context(), id, UpdateUserInput{Email: req.Email})
	if err != nil {
		return handleServiceError(c, "update user", err)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

// Delete handles DELETE /api/v1/users/:id.
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, "delete user", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetProfile handles GET /api/v1/users/:id/profile.
func (h *Handler) GetProfile(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	prof, err := h.profileService.GetByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "get profile by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(prof)
}

// GetExperiences handles GET /api/v1/users/:id/experiences.
func (h *Handler) GetExperiences(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	experiences, err := h.experienceService.ListByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "list experiences by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(experiences)
}

// GetProjects handles GET /api/v1/users/:id/projects.
func (h *Handler) GetProjects(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	projects, err := h.projectService.ListByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "list projects by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(projects)
}

// GetSkills handles GET /api/v1/users/:id/skills.
func (h *Handler) GetSkills(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	skills, err := h.skillService.ListByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "list skills by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(skills)
}

// GetEducation handles GET /api/v1/users/:id/education.
func (h *Handler) GetEducation(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	entries, err := h.educationService.ListByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "list education by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(entries)
}

// GetCertifications handles GET /api/v1/users/:id/certifications.
func (h *Handler) GetCertifications(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	certifications, err := h.certificationService.ListByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "list certifications by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(certifications)
}

// GetLanguages handles GET /api/v1/users/:id/languages.
func (h *Handler) GetLanguages(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	languages, err := h.languageService.ListByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "list languages by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(languages)
}

// handleServiceError logs the underlying error and writes a safe, generic
// JSON error response, mapping "not found" to 404.
func handleServiceError(c fiber.Ctx, action string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return httpapi.WriteError(c, fiber.StatusNotFound, "not found")
	}

	log.Printf("%s: %v", action, err)
	return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")
}
