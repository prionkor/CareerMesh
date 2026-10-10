package user_role

import (
	"errors"
	"log"

	"uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/prionkor/careermesh/internal/httpapi"
)

const (
	foreignKeyViolation = "23503"
	userFK              = "user_roles_user_id_fkey"
	roleFK              = "user_roles_role_id_fkey"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type UpdateRoleRequest struct {
	UserID string `json:"userId"`
	RoleID string `json:"roleId"`
}

func (h *Handler) Update(c fiber.Ctx) error {
	var req UpdateRoleRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid userId")
	}
	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid roleId")
	}

	if err := h.service.ReplaceAssignment(c.Context(), userID, roleID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) || isMissingReference(err) {
			return httpapi.WriteError(c, fiber.StatusNotFound, "user or role not found")
		}
		log.Printf("update user role: %v", err)
		return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// isMissingReference reports a foreign-key violation on the user or role reference.
func isMissingReference(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != foreignKeyViolation {
		return false
	}
	return pgErr.ConstraintName == userFK || pgErr.ConstraintName == roleFK
}
