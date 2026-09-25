package certification

import (
	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/middleware"
)

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	certifications := v1.Group("/certifications")

	certifications.Get("/:id", middleware.RequirePermission(authorization.PermCertificationsReadOwn), h.GetByID)
	certifications.Post("", middleware.RequirePermission(authorization.PermCertificationsCreateOwn), h.Create)
	certifications.Put("/:id", middleware.RequirePermission(authorization.PermCertificationsUpdateOwn), h.Update)
	certifications.Delete("/:id", middleware.RequirePermission(authorization.PermCertificationsDeleteOwn), h.Delete)
}
