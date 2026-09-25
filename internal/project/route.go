package project

import (
	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/middleware"
)

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	projects := v1.Group("/projects")

	projects.Get("/:id", middleware.RequirePermission(authorization.PermProjectsReadOwn), h.GetByID)
	projects.Post("", middleware.RequirePermission(authorization.PermProjectsCreateOwn), h.Create)
	projects.Put("/:id", middleware.RequirePermission(authorization.PermProjectsUpdateOwn), h.Update)
	projects.Delete("/:id", middleware.RequirePermission(authorization.PermProjectsDeleteOwn), h.Delete)
}
