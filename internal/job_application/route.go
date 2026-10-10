package job_application

import (
	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/middleware"
)

func (h *Handler) RegisterRoutes(v1 fiber.Router) {
	jobApplications := v1.Group("/job-applications")

	jobApplications.Get("", middleware.RequirePermission(authorization.PermJobApplicationsReadOwn), h.List)
	jobApplications.Get("/:id", middleware.RequirePermission(authorization.PermJobApplicationsReadOwn), h.GetByID)
	jobApplications.Post("", middleware.RequirePermission(authorization.PermJobApplicationsCreateOwn), h.Create)
	jobApplications.Put("/:id", middleware.RequirePermission(authorization.PermJobApplicationsUpdateOwn), h.Update)
	jobApplications.Delete("/:id", middleware.RequirePermission(authorization.PermJobApplicationsDeleteOwn), h.Delete)
}
