package user

import (
	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/middleware"
)

// RegisterRoutes registers user routes. Registration (POST) is intentionally
// public since there is no authenticated identity before an account exists;
// every other user route requires the caller to already be authenticated,
// so it is registered on protected.
func (h *Handler) RegisterRoutes(protected fiber.Router) {

	user := protected.Group("/users")
	user.Post("", h.Create)
	user.Get("/:id", middleware.RequirePermission(authorization.PermUsersReadOwn), h.GetByID)
	user.Put("/:id", middleware.RequirePermission(authorization.PermUsersUpdateOwn), h.Update)
	user.Delete("/:id", middleware.RequirePermission(authorization.PermUsersDeleteOwn), h.Delete)

	// get profile by user id
	user.Get("/:id/profile", middleware.RequirePermission(authorization.PermProfilesReadOwn), h.GetProfile)

	user.Get("/:id/experiences", middleware.RequirePermission(authorization.PermExperiencesReadOwn), h.GetExperiences)
	user.Get("/:id/projects", middleware.RequirePermission(authorization.PermProjectsReadOwn), h.GetProjects)
	user.Get("/:id/skills", middleware.RequirePermission(authorization.PermSkillsReadOwn), h.GetSkills)
	user.Get("/:id/education", middleware.RequirePermission(authorization.PermEducationReadOwn), h.GetEducation)
	user.Get("/:id/certifications", middleware.RequirePermission(authorization.PermCertificationsReadOwn), h.GetCertifications)
	user.Get("/:id/languages", middleware.RequirePermission(authorization.PermLanguagesReadOwn), h.GetLanguages)
}
