package language

import (
	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/middleware"
)

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	languages := v1.Group("/languages")

	languages.Get("/:id", middleware.RequirePermission(authorization.PermLanguagesReadOwn), h.GetByID)
	languages.Post("", middleware.RequirePermission(authorization.PermLanguagesCreateOwn), h.Create)
	languages.Put("/:id", middleware.RequirePermission(authorization.PermLanguagesUpdateOwn), h.Update)
	languages.Delete("/:id", middleware.RequirePermission(authorization.PermLanguagesDeleteOwn), h.Delete)
}
