package skill

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	skills := v1.Group("/skills")

	skills.Get("/:id", h.GetByID)
	skills.Post("", h.Create)
	skills.Put("/:id", h.Update)
	skills.Delete("/:id", h.Delete)
}
