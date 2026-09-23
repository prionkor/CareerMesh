package education

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	education := v1.Group("/education")

	education.Get("/:id", h.GetByID)
	education.Post("", h.Create)
	education.Put("/:id", h.Update)
	education.Delete("/:id", h.Delete)
}
