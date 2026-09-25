package experience

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

// Handler exposes the Experience service over HTTP.
type Handler struct {
	service *Service
}

// NewHandler creates an Experience HTTP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateExperienceRequest is the client-provided body for creating an experience.
type CreateExperienceRequest struct {
	Company        string                `json:"company"`
	Website        string                `json:"website"`
	Title          string                `json:"title"`
	EmploymentType models.EmploymentType `json:"employment_type"`
	Location       string                `json:"location"`
	StartDate      time.Time             `json:"start_date"`
	EndDate        *time.Time            `json:"end_date"`
	Description    string                `json:"description"`
}

// UpdateExperienceRequest is the client-provided body for updating an experience.
type UpdateExperienceRequest struct {
	Company        string                `json:"company"`
	Website        string                `json:"website"`
	Title          string                `json:"title"`
	EmploymentType models.EmploymentType `json:"employment_type"`
	Location       string                `json:"location"`
	StartDate      time.Time             `json:"start_date"`
	EndDate        *time.Time            `json:"end_date"`
	Description    string                `json:"description"`
}

// GetByID handles GET /api/v1/experiences/:id.
func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid experience id")
	}

	exp, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get experience by id", err)
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		exp.UserID,
		middleware.Permissions(c),
		authorization.PermExperiencesReadOwn,
		authorization.PermExperiencesReadAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	return c.Status(fiber.StatusOK).JSON(exp)
}

// Create handles POST /api/v1/experiences.
func (h *Handler) Create(c fiber.Ctx) error {
	var req CreateExperienceRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	exp := &models.Experience{
		UserID:         middleware.UserID(c),
		Company:        req.Company,
		Website:        req.Website,
		Title:          req.Title,
		EmploymentType: req.EmploymentType,
		Location:       req.Location,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Description:    req.Description,
	}

	created, err := h.service.Create(c.Context(), exp)
	if err != nil {
		return handleServiceError(c, "create experience", err)
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// Update handles PUT /api/v1/experiences/:id.
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid experience id")
	}

	var req UpdateExperienceRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	saved, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get experience for update", err)
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		saved.UserID,
		middleware.Permissions(c),
		authorization.PermExperiencesUpdateOwn,
		authorization.PermExperiencesUpdateAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	exp := &models.Experience{
		ID:             id,
		UserID:         saved.UserID,
		Company:        req.Company,
		Website:        req.Website,
		Title:          req.Title,
		EmploymentType: req.EmploymentType,
		Location:       req.Location,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Description:    req.Description,
	}

	updated, err := h.service.Update(c.Context(), exp)
	if err != nil {
		return handleServiceError(c, "update experience", err)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

// Delete handles DELETE /api/v1/experiences/:id.
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid experience id")
	}

	saved, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get experience for delete", err)
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		saved.UserID,
		middleware.Permissions(c),
		authorization.PermExperiencesDeleteOwn,
		authorization.PermExperiencesDeleteAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, "delete experience", err)
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
		return httpapi.WriteError(c, fiber.StatusNotFound, "experience not found")
	}

	log.Printf("%s: %v", action, err)
	return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")
}
