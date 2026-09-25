package profile

import (
	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/middleware"
)

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	profile := v1.Group("/profiles")

	profile.Get("/:id", middleware.RequirePermission(authorization.PermProfilesReadOwn), h.GetByID)
	profile.Post("", middleware.RequirePermission(authorization.PermProfilesCreateOwn), h.Create)
	profile.Put("/:id", middleware.RequirePermission(authorization.PermProfilesUpdateOwn), h.Update)
	profile.Delete("/:id", middleware.RequirePermission(authorization.PermProfilesDeleteOwn), h.Delete)
}
