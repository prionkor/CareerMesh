package user

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	user := v1.Group("/users")

	user.Get("/:id", h.GetByID)
	user.Post("", h.Create)
	user.Put("/:id", h.Update)
	user.Delete("/:id", h.Delete)

	// get profile by user id
	user.Get("/:id/profile", h.GetProfile)
}
