package certification

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

// Handler exposes the Certification service over HTTP.
type Handler struct {
	service *Service
}

// NewHandler creates a Certification HTTP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateCertificationRequest is the client-provided body for creating a certification.
type CreateCertificationRequest struct {
	Name         string     `json:"name"`
	Issuer       string     `json:"issuer"`
	IssueDate    *time.Time `json:"issue_date"`
	ExpiryDate   *time.Time `json:"expiry_date"`
	CredentialID string     `json:"credential_id"`
	URL          string     `json:"url"`
}

// UpdateCertificationRequest is the client-provided body for updating a certification.
type UpdateCertificationRequest struct {
	Name         string     `json:"name"`
	Issuer       string     `json:"issuer"`
	IssueDate    *time.Time `json:"issue_date"`
	ExpiryDate   *time.Time `json:"expiry_date"`
	CredentialID string     `json:"credential_id"`
	URL          string     `json:"url"`
}

// GetByID handles GET /api/v1/certifications/:id.
func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid certification id")
	}

	cert, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get certification by id", err)
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		cert.UserID,
		middleware.Permissions(c),
		authorization.PermCertificationsReadOwn,
		authorization.PermCertificationsReadAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	return c.Status(fiber.StatusOK).JSON(cert)
}

// Create handles POST /api/v1/certifications.
func (h *Handler) Create(c fiber.Ctx) error {
	var req CreateCertificationRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	cert := &models.Certification{
		UserID:       middleware.UserID(c),
		Name:         req.Name,
		Issuer:       req.Issuer,
		IssueDate:    req.IssueDate,
		ExpiryDate:   req.ExpiryDate,
		CredentialID: req.CredentialID,
		URL:          req.URL,
	}

	created, err := h.service.Create(c.Context(), cert)
	if err != nil {
		return handleServiceError(c, "create certification", err)
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// Update handles PUT /api/v1/certifications/:id.
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid certification id")
	}

	var req UpdateCertificationRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	saved, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get certification for update", err)
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		saved.UserID,
		middleware.Permissions(c),
		authorization.PermCertificationsUpdateOwn,
		authorization.PermCertificationsUpdateAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	cert := &models.Certification{
		ID:           id,
		UserID:       saved.UserID,
		Name:         req.Name,
		Issuer:       req.Issuer,
		IssueDate:    req.IssueDate,
		ExpiryDate:   req.ExpiryDate,
		CredentialID: req.CredentialID,
		URL:          req.URL,
	}

	updated, err := h.service.Update(c.Context(), cert)
	if err != nil {
		return handleServiceError(c, "update certification", err)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

// Delete handles DELETE /api/v1/certifications/:id.
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid certification id")
	}

	saved, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get certification for delete", err)
	}

	if !authorization.CanAccessResource(
		middleware.UserID(c),
		saved.UserID,
		middleware.Permissions(c),
		authorization.PermCertificationsDeleteOwn,
		authorization.PermCertificationsDeleteAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, "delete certification", err)
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
		return httpapi.WriteError(c, fiber.StatusNotFound, "certification not found")
	}

	log.Printf("%s: %v", action, err)
	return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")
}
