// Package httpapi contains small, shared HTTP response helpers used by
// handler packages. It has no dependency on any domain package.
package httpapi

import "github.com/gofiber/fiber/v3"

// DevUserIDHeader is a temporary, explicit development-only way to identify
// the acting user for POST requests on user-owned resources until real
// authentication exists. There is currently no authenticated session or
// request context carrying a user ID anywhere in the project, so this
// header must be sent by the client. Replace this with the authenticated
// user ID once authentication is implemented; do not build a fake
// authentication layer around it.
const DevUserIDHeader = "X-Debug-User-ID"

// ErrorResponse is the standard error body returned by the API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteError writes a simple {"error": message} JSON body with the given status code.
func WriteError(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(ErrorResponse{Error: message})
}
