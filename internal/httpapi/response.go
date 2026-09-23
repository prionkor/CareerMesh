// Package httpapi contains small, shared HTTP response helpers used by
// handler packages. It has no dependency on any domain package.
package httpapi

import "github.com/gofiber/fiber/v3"

// ErrorResponse is the standard error body returned by the API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteError writes a simple {"error": message} JSON body with the given status code.
func WriteError(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(ErrorResponse{Error: message})
}
