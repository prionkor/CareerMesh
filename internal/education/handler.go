package education

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

// Handler exposes the Education service over HTTP.
type Handler struct {
	service *Service
}

// NewHandler creates an Education HTTP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateEducationRequest is the client-provided body for creating an education entry.
type CreateEducationRequest struct {
	Institution string     `json:"institution"`
	Degree      string     `json:"degree"`
	Field       string     `json:"field"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Description string     `json:"description"`
}

// UpdateEducationRequest is the client-provided body for updating an education entry.
type UpdateEducationRequest struct {
	Institution string     `json:"institution"`
	Degree      string     `json:"degree"`
	Field       string     `json:"field"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Description string     `json:"description"`
}

// GetByID handles GET /api/v1/education/:id.
func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid education id")
	}

	edu, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get education by id", err)
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		edu.UserID,
		middleware.Permissions(c),
		authorization.PermEducationReadOwn,
		authorization.PermEducationReadAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	return c.Status(fiber.StatusOK).JSON(edu)
}

// Create handles POST /api/v1/education.
func (h *Handler) Create(c fiber.Ctx) error {
	var req CreateEducationRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	edu := &models.Education{
		UserID:      middleware.UserID(c),
		Institution: req.Institution,
		Degree:      req.Degree,
		Field:       req.Field,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Description: req.Description,
	}

	created, err := h.service.Create(c.Context(), edu)
	if err != nil {
		return handleServiceError(c, "create education", err)
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// Update handles PUT /api/v1/education/:id.
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid education id")
	}

	var req UpdateEducationRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	saved, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get education for update", err)
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		saved.UserID,
		middleware.Permissions(c),
		authorization.PermEducationUpdateOwn,
		authorization.PermEducationUpdateAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	edu := &models.Education{
		ID:          id,
		UserID:      saved.UserID,
		Institution: req.Institution,
		Degree:      req.Degree,
		Field:       req.Field,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Description: req.Description,
	}

	updated, err := h.service.Update(c.Context(), edu)
	if err != nil {
		return handleServiceError(c, "update education", err)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

// Delete handles DELETE /api/v1/education/:id.
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid education id")
	}

	saved, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get education for delete", err)
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		saved.UserID,
		middleware.Permissions(c),
		authorization.PermEducationDeleteOwn,
		authorization.PermEducationDeleteAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, "delete education", err)
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
		return httpapi.WriteError(c, fiber.StatusNotFound, "education not found")
	}

	log.Printf("%s: %v", action, err)
	return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")
}
