package profile

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	profile := v1.Group("/profiles")

	profile.Get("/:id", h.GetByID)
	profile.Post("", h.Create)
	profile.Put("/:id", h.Update)
	profile.Delete("/:id", h.Delete)
}
