package skill

import (
	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/middleware"
)

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	skills := v1.Group("/skills")

	skills.Get("/:id", middleware.RequirePermission(authorization.PermSkillsReadOwn), h.GetByID)
	skills.Post("", middleware.RequirePermission(authorization.PermSkillsCreateOwn), h.Create)
	skills.Put("/:id", middleware.RequirePermission(authorization.PermSkillsUpdateOwn), h.Update)
	skills.Delete("/:id", middleware.RequirePermission(authorization.PermSkillsDeleteOwn), h.Delete)
}
