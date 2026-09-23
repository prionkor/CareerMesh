package language

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/prionkor/careermesh/internal/httpapi"
	"github.com/prionkor/careermesh/models"
)

// Handler exposes the Language service over HTTP.
type Handler struct {
	service *Service
}

// NewHandler creates a Language HTTP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateLanguageRequest is the client-provided body for creating a language.
type CreateLanguageRequest struct {
	Name        string                     `json:"name"`
	Proficiency models.LanguageProficiency `json:"proficiency"`
}

// UpdateLanguageRequest is the client-provided body for updating a language.
type UpdateLanguageRequest struct {
	Name        string                     `json:"name"`
	Proficiency models.LanguageProficiency `json:"proficiency"`
}

// GetByID handles GET /api/v1/languages/:id.
func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid language id")
	}

	lang, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get language by id", err)
	}

	return c.Status(fiber.StatusOK).JSON(lang)
}

// Create handles POST /api/v1/languages.
func (h *Handler) Create(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Get(httpapi.DevUserIDHeader))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "missing or invalid "+httpapi.DevUserIDHeader+" header")
	}

	var req CreateLanguageRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	lang := &models.Language{
		UserID:      userID,
		Name:        req.Name,
		Proficiency: req.Proficiency,
	}

	created, err := h.service.Create(c.Context(), lang)
	if err != nil {
		return handleServiceError(c, "create language", err)
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// Update handles PUT /api/v1/languages/:id.
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid language id")
	}

	var req UpdateLanguageRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	lang := &models.Language{
		ID:          id,
		Name:        req.Name,
		Proficiency: req.Proficiency,
	}

	updated, err := h.service.Update(c.Context(), lang)
	if err != nil {
		return handleServiceError(c, "update language", err)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

// Delete handles DELETE /api/v1/languages/:id.
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid language id")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, "delete language", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// handleServiceError logs the underlying error and writes a safe, generic
// JSON error response, mapping "not found" to 404.
func handleServiceError(c fiber.Ctx, action string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return httpapi.WriteError(c, fiber.StatusNotFound, "language not found")
	}

	log.Printf("%s: %v", action, err)
	return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")
}
