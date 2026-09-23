package skill

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/prionkor/careermesh/internal/httpapi"
	"github.com/prionkor/careermesh/models"
)

// Handler exposes the Skill service over HTTP.
type Handler struct {
	service *Service
}

// NewHandler creates a Skill HTTP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateSkillRequest is the client-provided body for creating a skill.
type CreateSkillRequest struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

// UpdateSkillRequest is the client-provided body for updating a skill.
type UpdateSkillRequest struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

// GetByID handles GET /api/v1/skills/:id.
func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid skill id")
	}

	sk, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get skill by id", err)
	}

	return c.Status(fiber.StatusOK).JSON(sk)
}

// Create handles POST /api/v1/skills.
func (h *Handler) Create(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Get(httpapi.DevUserIDHeader))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "missing or invalid "+httpapi.DevUserIDHeader+" header")
	}

	var req CreateSkillRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	sk := &models.Skill{
		UserID:   userID,
		Name:     req.Name,
		Category: req.Category,
	}

	created, err := h.service.Create(c.Context(), sk)
	if err != nil {
		return handleServiceError(c, "create skill", err)
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// Update handles PUT /api/v1/skills/:id.
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid skill id")
	}

	var req UpdateSkillRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	sk := &models.Skill{
		ID:       id,
		Name:     req.Name,
		Category: req.Category,
	}

	updated, err := h.service.Update(c.Context(), sk)
	if err != nil {
		return handleServiceError(c, "update skill", err)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

// Delete handles DELETE /api/v1/skills/:id.
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid skill id")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, "delete skill", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// handleServiceError logs the underlying error and writes a safe, generic
// JSON error response, mapping "not found" to 404.
func handleServiceError(c fiber.Ctx, action string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return httpapi.WriteError(c, fiber.StatusNotFound, "skill not found")
	}

	log.Printf("%s: %v", action, err)
	return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")
}
