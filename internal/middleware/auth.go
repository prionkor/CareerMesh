// Package middleware provides Fiber middleware for authenticating requests
// and enforcing RBAC permission checks.
package middleware

import (
	"strings"

	"uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/auth"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/httpapi"
)

const (
	localsUserID      = "auth_user_id"
	localsPermissions = "auth_permissions"
)

// RequireAuth validates the JWT bearer token and stores the authenticated
// user id and permissions in request locals for downstream handlers.
func RequireAuth(jwtSecret []byte) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")

		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || token == "" {
			return httpapi.WriteError(c, fiber.StatusUnauthorized, "missing or invalid authorization header")
		}

		claims, err := auth.ParseAccessToken(jwtSecret, token)
		if err != nil {
			return httpapi.WriteError(c, fiber.StatusUnauthorized, "invalid or expired token")
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return httpapi.WriteError(c, fiber.StatusUnauthorized, "invalid token subject")
		}

		c.Locals(localsUserID, userID)
		c.Locals(localsPermissions, claims.Permissions)

		return c.Next()
	}
}

// RequirePermission rejects the request unless the authenticated permissions
// include the required permission. A frontend hiding UI is not a substitute
// for this check.
func RequirePermission(permission string) fiber.Handler {
	return func(c fiber.Ctx) error {
		if !authorization.HasPermission(Permissions(c), permission) {
			return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
		}

		return c.Next()
	}
}

// UserID returns the authenticated user id stored by RequireAuth.
func UserID(c fiber.Ctx) uuid.UUID {
	id, _ := c.Locals(localsUserID).(uuid.UUID)
	return id
}

// Permissions returns the authenticated permissions stored by RequireAuth.
func Permissions(c fiber.Ctx) []string {
	permissions, _ := c.Locals(localsPermissions).([]string)
	return permissions
}
