package job_description

import (
	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/middleware"
)

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	jobDescriptions := v1.Group("/job-descriptions")

	jobDescriptions.Get("", middleware.RequirePermission(authorization.PermJobDescriptionsReadOwn), h.List)
	jobDescriptions.Get("/:id", middleware.RequirePermission(authorization.PermJobDescriptionsReadOwn), h.GetByID)
	jobDescriptions.Get("/:id/versions", middleware.RequirePermission(authorization.PermJobDescriptionsReadOwn), h.ListVersions)
	jobDescriptions.Get("/:id/versions/:versionID", middleware.RequirePermission(authorization.PermJobDescriptionsReadOwn), h.GetVersion)
	jobDescriptions.Post("", middleware.RequirePermission(authorization.PermJobDescriptionsCreateOwn), h.Create)
	jobDescriptions.Put("/:id", middleware.RequirePermission(authorization.PermJobDescriptionsUpdateOwn), h.Update)
	jobDescriptions.Delete("/:id", middleware.RequirePermission(authorization.PermJobDescriptionsDeleteOwn), h.Delete)
}
