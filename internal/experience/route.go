package experience

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	experiences := v1.Group("/experiences")

	experiences.Get("/:id", h.GetByID)
	experiences.Post("", h.Create)
	experiences.Put("/:id", h.Update)
	experiences.Delete("/:id", h.Delete)
}
