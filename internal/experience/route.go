package experience

import (
	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/middleware"
)

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	experiences := v1.Group("/experiences")

	experiences.Get("/:id", middleware.RequirePermission(authorization.PermExperiencesReadOwn), h.GetByID)
	experiences.Post("", middleware.RequirePermission(authorization.PermExperiencesCreateOwn), h.Create)
	experiences.Put("/:id", middleware.RequirePermission(authorization.PermExperiencesUpdateOwn), h.Update)
	experiences.Delete("/:id", middleware.RequirePermission(authorization.PermExperiencesDeleteOwn), h.Delete)
}
