package language

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	languages := v1.Group("/languages")

	languages.Get("/:id", h.GetByID)
	languages.Post("", h.Create)
	languages.Put("/:id", h.Update)
	languages.Delete("/:id", h.Delete)
}
