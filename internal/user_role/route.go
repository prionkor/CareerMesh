package user_role

import (
	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/middleware"
)

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	userRoles := v1.Group("/user_roles")
	userRoles.Put("", middleware.RequirePermission(authorization.PermRoleUpdateAll), h.Update)
}
