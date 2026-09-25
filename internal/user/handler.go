package user

import (
	"errors"
	"log"

	"uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/certification"
	"github.com/prionkor/careermesh/internal/education"
	"github.com/prionkor/careermesh/internal/experience"
	"github.com/prionkor/careermesh/internal/httpapi"
	"github.com/prionkor/careermesh/internal/language"
	"github.com/prionkor/careermesh/internal/middleware"
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

// GetByID handles GET /api/v1/users/.
func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	usr, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get user by id", err)
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		usr.ID,
		middleware.Permissions(c),
		authorization.PermUsersReadOwn,
		authorization.PermUsersReadAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
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

// Update handles PUT /api/v1/users/. It updates account fields only and
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

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		updated.ID,
		middleware.Permissions(c),
		authorization.PermUsersUpdateOwn,
		authorization.PermUsersUpdateAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	return c.Status(fiber.StatusOK).JSON(updated)

}

// Delete handles DELETE /api/v1/users/.
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	_, err = h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get user for delete", err)
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		id,
		middleware.Permissions(c),
		authorization.PermUsersDeleteOwn,
		authorization.PermUsersDeleteAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, "delete user", err)
	}

	return c.SendStatus(fiber.StatusNoContent)

}

// GetProfile handles GET /api/v1/users//profile.
func (h *Handler) GetProfile(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		userID,
		middleware.Permissions(c),
		authorization.PermProfilesReadOwn,
		authorization.PermProfilesReadAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	prof, err := h.profileService.GetByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "get profile by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(prof)

}

// GetExperiences handles GET /api/v1/users//experiences.
func (h *Handler) GetExperiences(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		userID,
		middleware.Permissions(c),
		authorization.PermExperiencesReadOwn,
		authorization.PermExperiencesReadAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	experiences, err := h.experienceService.ListByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "list experiences by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(experiences)

}

// GetProjects handles GET /api/v1/users//projects.
func (h *Handler) GetProjects(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		userID,
		middleware.Permissions(c),
		authorization.PermProjectsReadOwn,
		authorization.PermProjectsReadAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	projects, err := h.projectService.ListByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "list projects by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(projects)

}

// GetSkills handles GET /api/v1/users//skills.
func (h *Handler) GetSkills(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		userID,
		middleware.Permissions(c),
		authorization.PermSkillsReadOwn,
		authorization.PermSkillsReadAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	skills, err := h.skillService.ListByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "list skills by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(skills)

}

// GetEducation handles GET /api/v1/users//education.
func (h *Handler) GetEducation(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		userID,
		middleware.Permissions(c),
		authorization.PermEducationReadOwn,
		authorization.PermEducationReadAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	entries, err := h.educationService.ListByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "list education by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(entries)

}

// GetCertifications handles GET /api/v1/users//certifications.
func (h *Handler) GetCertifications(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		userID,
		middleware.Permissions(c),
		authorization.PermCertificationsReadOwn,
		authorization.PermCertificationsReadAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	certifications, err := h.certificationService.ListByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "list certifications by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(certifications)

}

// GetLanguages handles GET /api/v1/users//languages.
func (h *Handler) GetLanguages(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		userID,
		middleware.Permissions(c),
		authorization.PermLanguagesReadOwn,
		authorization.PermLanguagesReadAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	languages, err := h.languageService.ListByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "list languages by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(languages)

}

// handleServiceError logs the underlying error and writes a safe, generic
// JSON error response, mapping "not found" to 404 and ownership violations to 403.
func handleServiceError(c fiber.Ctx, action string, err error) error {
	if errors.Is(err, authorization.ErrForbidden) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return httpapi.WriteError(c, fiber.StatusNotFound, "not found")
	}

	log.Printf("%s: %v", action, err)
	return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")

}
