package project

import (
	"errors"
	"log"
	"time"

	"uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/httpapi"
	"github.com/prionkor/careermesh/internal/middleware"
	"github.com/prionkor/careermesh/models"
)

// Handler exposes the Project service over HTTP.
type Handler struct {
	service *Service
}

// NewHandler creates a Project HTTP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateProjectRequest is the client-provided body for creating a project.
type CreateProjectRequest struct {
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	URL           string     `json:"url"`
	RepositoryURL string     `json:"repository_url"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
}

// UpdateProjectRequest is the client-provided body for updating a project.
type UpdateProjectRequest struct {
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	URL           string     `json:"url"`
	RepositoryURL string     `json:"repository_url"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
}

// GetByID handles GET /api/v1/projects/:id.
func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid project id")
	}

	proj, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get project by id", err)
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		proj.UserID,
		middleware.Permissions(c),
		authorization.PermProjectsReadOwn,
		authorization.PermProjectsReadAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	return c.Status(fiber.StatusOK).JSON(proj)
}

// Create handles POST /api/v1/projects.
func (h *Handler) Create(c fiber.Ctx) error {
	var req CreateProjectRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	proj := &models.Project{
		UserID:        middleware.UserID(c),
		Name:          req.Name,
		Description:   req.Description,
		URL:           req.URL,
		RepositoryURL: req.RepositoryURL,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
	}

	created, err := h.service.Create(c.Context(), proj)
	if err != nil {
		return handleServiceError(c, "create project", err)
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// Update handles PUT /api/v1/projects/:id.
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid project id")
	}

	var req UpdateProjectRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	saved, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get project for update", err)
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		saved.UserID,
		middleware.Permissions(c),
		authorization.PermProjectsUpdateOwn,
		authorization.PermProjectsUpdateAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	proj := &models.Project{
		ID:            id,
		UserID:        saved.UserID,
		Name:          req.Name,
		Description:   req.Description,
		URL:           req.URL,
		RepositoryURL: req.RepositoryURL,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
	}

	updated, err := h.service.Update(c.Context(), proj)
	if err != nil {
		return handleServiceError(c, "update project", err)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

// Delete handles DELETE /api/v1/projects/:id.
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid project id")
	}

	saved, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get project for delete", err)
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		saved.UserID,
		middleware.Permissions(c),
		authorization.PermProjectsDeleteOwn,
		authorization.PermProjectsDeleteAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, "delete project", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// handleServiceError logs the underlying error and writes a safe, generic
// JSON error response, mapping "not found" to 404.
func handleServiceError(c fiber.Ctx, action string, err error) error {
	if errors.Is(err, authorization.ErrForbidden) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return httpapi.WriteError(c, fiber.StatusNotFound, "project not found")
	}

	log.Printf("%s: %v", action, err)
	return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")
}
