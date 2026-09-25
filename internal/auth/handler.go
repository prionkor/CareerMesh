package auth

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/httpapi"
	"uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// LoginRequest is the client-provided body for logging in.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUserResponse is the safe user projection returned on login.
// PasswordHash is intentionally omitted.
type LoginUserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

// LoginResponse is the response body for a successful login.
type LoginResponse struct {
	AccessToken string            `json:"access_token"`
	TokenType   string            `json:"token_type"`
	ExpiresIn   int               `json:"expires_in"`
	User        LoginUserResponse `json:"user"`
}

// Login handles POST /api/v1/auth/login.
func (h *Handler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	result, err := h.service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return httpapi.WriteError(c, fiber.StatusUnauthorized, "invalid email or password")
		}

		log.Printf("login: %v", err)
		return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(LoginResponse{
		AccessToken: result.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(AccessTokenTTL.Seconds()),
		User: LoginUserResponse{
			ID:    result.User.ID,
			Email: result.User.Email,
		},
	})
}
