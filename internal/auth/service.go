package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"uuid"

	"github.com/prionkor/careermesh/internal/pending_user"
	"github.com/prionkor/careermesh/internal/utility"
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
	CreateWithHashedPassword(ctx context.Context, email, passwordHash string) (*models.User, error)
}

type Service struct {
	users       UserLookup
	jwtSecret   []byte
	PendingUser pending_user.Service
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

	ok, err := utility.VerifyPassword(usr.PasswordHash, password)
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

func (s *Service) Register(ctx context.Context, email, password string) (*models.PendingUser, error) {
	// Hash the password using the existing password-hashing mechanism
	passwordHash, err := utility.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// Generate a cryptographically secure verification token
	verificationToken, err := generateVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("generate verification token: %w", err)
	}

	// Hash the verification token for storage
	verificationTokenHash, err := utility.HashPassword(verificationToken)
	if err != nil {
		return nil, fmt.Errorf("hash verification token: %w", err)
	}

	// Create the pending user record
	expireAt := time.Now().Add(24 * time.Hour) // 24 hours from now
	user, err := s.PendingUser.Create(ctx, email, passwordHash, verificationTokenHash, expireAt)
	if err != nil {
		return nil, fmt.Errorf("create pending user: %w", err)
	}

	return user, nil
}

// Verify verifies a user's email using a verification token
func (s *Service) Verify(ctx context.Context, token string) error {
	// Find the pending user using the token
	pendingUser, err := s.PendingUser.FindByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("find pending user: %w", err)
	}

	// Check if the token has expired
	if time.Now().After(pendingUser.ExpiresAt) {
		return fmt.Errorf("token has expired")
	}

	// Verify the token matches the stored hash
	if ok, err := utility.VerifyPassword(pendingUser.TokenHash, token); err != nil || !ok {
		return fmt.Errorf("invalid token")
	}

	// Create the actual user using the existing user service/repository conventions
	_, err = s.users.GetByEmail(ctx, pendingUser.Email)
	if err == nil {
		// User already exists
		return fmt.Errorf("user already exists")
	}

	// Create the user using the existing user service/repository conventions
	_, err = s.users.CreateWithHashedPassword(ctx, pendingUser.Email, pendingUser.PasswordHash)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	// Delete the pending user record
	if err := s.PendingUser.Delete(ctx, pendingUser.ID); err != nil {
		return fmt.Errorf("delete pending user: %w", err)
	}

	// todo: send notificaiton email to user

	return nil
}

func generateVerificationToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}
