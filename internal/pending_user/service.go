package pending_user

import (
	"context"
	"fmt"
	"time"

	"github.com/prionkor/careermesh/internal/utility"
	"github.com/prionkor/careermesh/models"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

// Create creates a new pending user
func (s *Service) Create(ctx context.Context, email, passwordHash, tokenHash string, expiresAt time.Time) (*models.PendingUser, error) {
	user, err := s.repository.Create(ctx, email, passwordHash, tokenHash, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("create pending user: %w", err)
	}
	return user, nil
}

// FindByEmail returns a pending user by email
func (s *Service) FindByEmail(ctx context.Context, email string) (*models.PendingUser, error) {
	user, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("find pending user by email: %w", err)
	}
	return user, nil
}

// FindByToken returns a pending user by verification token
func (s *Service) FindByToken(ctx context.Context, token string) (*models.PendingUser, error) {
	// Hash the input token before matching
	hashedToken, err := utility.HashPassword(token)
	if err != nil {
		return nil, fmt.Errorf("hash verification token: %w", err)
	}

	user, err := s.repository.FindByToken(ctx, hashedToken)
	if err != nil {
		return nil, fmt.Errorf("find pending user by token: %w", err)
	}
	return user, nil
}

// Update updates a pending user
func (s *Service) Update(ctx context.Context, id string, email, passwordHash, tokenHash string, expiresAt time.Time) error {
	if err := s.repository.Update(ctx, id, email, passwordHash, tokenHash, expiresAt); err != nil {
		return fmt.Errorf("update pending user: %w", err)
	}
	return nil
}

// Delete deletes a pending user
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete pending user: %w", err)
	}
	return nil
}
