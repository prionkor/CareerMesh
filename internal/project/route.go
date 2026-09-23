package project

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	projects := v1.Group("/projects")

	projects.Get("/:id", h.GetByID)
	projects.Post("", h.Create)
	projects.Put("/:id", h.Update)
	projects.Delete("/:id", h.Delete)
}
