package auth

import (
	"context"
	"errors"
	"fmt"

	"uuid"

	"github.com/prionkor/careermesh/models"
)

// ErrInvalidCredentials is returned for any login failure. It intentionally
// does not distinguish between an unknown email and an incorrect password.
var ErrInvalidCredentials = errors.New("invalid email or password")

// UserLookup is the narrow slice of user.Service that auth needs. Defined
// here (rather than importing internal/user) to avoid an import cycle with
// internal/middleware, which both internal/auth and internal/user depend on.
type UserLookup interface {
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetPermissions(ctx context.Context, id uuid.UUID) ([]string, error)
}

type Service struct {
	users     UserLookup
	jwtSecret []byte
}

func NewService(users UserLookup, jwtSecret []byte) *Service {
	return &Service{users: users, jwtSecret: jwtSecret}
}

// LoginResult carries the outcome of a successful login.
type LoginResult struct {
	AccessToken string
	User        *models.User
}

// Login verifies the email/password pair and issues an access token
// containing the user's current effective permissions.
func (s *Service) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	usr, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	ok, err := VerifyPassword(usr.PasswordHash, password)
	if err != nil || !ok {
		return nil, ErrInvalidCredentials
	}

	permissions, err := s.users.GetPermissions(ctx, usr.ID)
	if err != nil {
		return nil, fmt.Errorf("load permissions: %w", err)
	}

	token, err := GenerateAccessToken(s.jwtSecret, usr.ID, permissions)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	return &LoginResult{AccessToken: token, User: usr}, nil
}
