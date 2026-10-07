package auth

import (
	"errors"
	"log"

	"uuid"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/httpapi"
)

var validate = validator.New()

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// LoginRequest is the client-provided body for logging in.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
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

	// Validate the request using validator/v10
	if err := validate.Struct(req); err != nil {
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

// RegisterRequest is the client-provided body for registration.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RegisterResponse is the response body for a successful registration.
type RegisterResponse struct {
	Success bool `json:"success"`
}

// Register handles POST /api/v1/auth/register.
func (h *Handler) Register(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	// Validate the request using validator/v10
	if err := validate.Struct(req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	// Register the user
	if _, err := h.service.Register(c.Context(), req.Email, req.Password); err != nil {
		log.Printf("register: %v", err)
		return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(RegisterResponse{Success: true})
}

// VerifyRequest is the client-provided query parameter for verification.
type VerifyRequest struct {
	Token string `query:"token" validate:"required"`
}

// VerifyResponse is the response body for a successful verification.
type VerifyResponse struct {
	Success bool `json:"success"`
}

// Verify handles GET /api/v1/auth/register/verify.
func (h *Handler) Verify(c fiber.Ctx) error {
	var req VerifyRequest
	if err := c.Bind().Query(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request query")
	}

	// Validate the request using validator/v10
	if err := validate.Struct(req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request query")
	}

	// Verify the token
	if err := h.service.Verify(c.Context(), req.Token); err != nil {
		// Don't expose the nature of the error
		return httpapi.WriteError(c, fiber.StatusUnauthorized, "invalid or expired token")
	}

	return c.Status(fiber.StatusOK).JSON(VerifyResponse{Success: true})
}
