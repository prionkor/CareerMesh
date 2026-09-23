package certification

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	certifications := v1.Group("/certifications")

	certifications.Get("/:id", h.GetByID)
	certifications.Post("", h.Create)
	certifications.Put("/:id", h.Update)
	certifications.Delete("/:id", h.Delete)
}
