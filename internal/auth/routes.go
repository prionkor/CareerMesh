package auth

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	auth := v1.Group("/auth")

	auth.Post("/login", h.Login)
}
