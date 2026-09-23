package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

	return createdUser, nil
}

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

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*models.User, error) {
	updatedUser, err := s.repository.Update(ctx, &models.User{
		ID:    id,
		Email: input.Email,
	})

	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return updatedUser, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
