package user

import (
	"context"
	"fmt"

	"uuid"

	"github.com/prionkor/careermesh/models"
)

type Service struct {
	repository *Repository
}

type UpdateUserInput struct {
	Email string
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, email, password string) (*models.User, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{Email: email, PasswordHash: passwordHash}

	createdUser, err := s.repository.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := s.repository.AssignDefaultRole(ctx, createdUser.ID); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return createdUser, nil
}

// GetByID returns the user identified by id. Only the user themselves or an
// admin may fetch the account.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}

// GetPermissions returns the effective, flattened permission set for a user.
func (s *Service) GetPermissions(ctx context.Context, id uuid.UUID) ([]string, error) {
	permissions, err := s.repository.GetPermissionsByUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get permissions: %w", err)
	}

	return permissions, nil
}

// Update updates account fields for id. Only the user themselves or an admin may update the account.
func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*models.User, error) {
	_, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	updatedUser, err := s.repository.Update(ctx, &models.User{
		ID:    id,
		Email: input.Email,
	})

	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return updatedUser, nil
}

// Delete deletes the account for id. Only the user themselves or an admin may delete the account.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
