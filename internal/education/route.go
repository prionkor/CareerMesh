package education

import (
	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/middleware"
)

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	education := v1.Group("/education")

	education.Get("/:id", middleware.RequirePermission(authorization.PermEducationReadOwn), h.GetByID)
	education.Post("", middleware.RequirePermission(authorization.PermEducationCreateOwn), h.Create)
	education.Put("/:id", middleware.RequirePermission(authorization.PermEducationUpdateOwn), h.Update)
	education.Delete("/:id", middleware.RequirePermission(authorization.PermEducationDeleteOwn), h.Delete)
}
