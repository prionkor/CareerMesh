package profile

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/prionkor/careermesh/internal/httpapi"
	"github.com/prionkor/careermesh/models"
)

// DevUserIDHeader is a temporary, explicit development-only way to identify
// the acting user for POST /api/v1/profiles until real authentication
// exists. There is currently no authenticated session or request context
// carrying a user ID anywhere in the project, so this header must be sent
// by the client. Replace this with the authenticated user ID once
// authentication is implemented; do not build a fake authentication layer
// around it.
const DevUserIDHeader = "X-Debug-User-ID"

// Handler exposes the Profile service over HTTP.
type Handler struct {
	service *Service
}

// NewHandler creates a Profile HTTP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateProfileRequest is the client-provided body for creating a profile.
type CreateProfileRequest struct {
	Name     string `json:"name"`
	Headline string `json:"headline"`
	Location string `json:"location"`
	Phone    string `json:"phone"`
	Website  string `json:"website"`
	GitHub   string `json:"github"`
	LinkedIn string `json:"linkedin"`
}

// UpdateProfileRequest is the client-provided body for updating a profile.
type UpdateProfileRequest struct {
	Name     string `json:"name"`
	Headline string `json:"headline"`
	Location string `json:"location"`
	Phone    string `json:"phone"`
	Website  string `json:"website"`
	GitHub   string `json:"github"`
	LinkedIn string `json:"linkedin"`
}

// GetByID handles GET /api/v1/profiles/:id.
func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid profile id")
	}

	prof, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get profile by id", err)
	}

	return c.Status(fiber.StatusOK).JSON(prof)
}

// GetByUserID handles GET /api/v1/users/:userId/profile.
func (h *Handler) GetByUserID(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid user id")
	}

	prof, err := h.service.GetByUserID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, "get profile by user id", err)
	}

	return c.Status(fiber.StatusOK).JSON(prof)
}

// Create handles POST /api/v1/profiles.
func (h *Handler) Create(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Get(DevUserIDHeader))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "missing or invalid "+DevUserIDHeader+" header")
	}

	var req CreateProfileRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	prof := &models.Profile{
		UserID:   userID,
		Name:     req.Name,
		Headline: req.Headline,
		Location: req.Location,
		Phone:    req.Phone,
		Website:  req.Website,
		GitHub:   req.GitHub,
		LinkedIn: req.LinkedIn,
	}

	created, err := h.service.Create(c.Context(), prof)
	if err != nil {
		return handleServiceError(c, "create profile", err)
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// Update handles PUT /api/v1/profiles/:id.
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid profile id")
	}

	var req UpdateProfileRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	prof := &models.Profile{
		ID:       id,
		Name:     req.Name,
		Headline: req.Headline,
		Location: req.Location,
		Phone:    req.Phone,
		Website:  req.Website,
		GitHub:   req.GitHub,
		LinkedIn: req.LinkedIn,
	}

	updated, err := h.service.Update(c.Context(), prof)
	if err != nil {
		return handleServiceError(c, "update profile", err)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

// Delete handles DELETE /api/v1/profiles/:id.
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid profile id")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, "delete profile", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// handleServiceError logs the underlying error and writes a safe, generic
// JSON error response, mapping "not found" to 404.
func handleServiceError(c fiber.Ctx, action string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return httpapi.WriteError(c, fiber.StatusNotFound, "profile not found")
	}

	log.Printf("%s: %v", action, err)
	return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")
}
