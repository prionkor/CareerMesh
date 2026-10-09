package user_role

import (
	"errors"
	"log"

	"uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/prionkor/careermesh/internal/httpapi"
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
		if errors.Is(err, pgx.ErrNoRows) {
			return httpapi.WriteError(c, fiber.StatusNotFound, "user or role not found")
		}
		log.Printf("update user role: %v", err)
		return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
