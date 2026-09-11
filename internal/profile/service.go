package profile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/prionkor/careermesh/models"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, profile *models.Profile) (*models.Profile, error) {
	createdProfile, err := s.repository.Create(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("create profile: %w", err)
	}

	return createdProfile, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Profile, error) {
	profile, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get profile by id: %w", err)
	}

	return profile, nil
}

func (s *Service) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Profile, error) {
	profile, err := s.repository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get profile by user id: %w", err)
	}

	return profile, nil
}

func (s *Service) Update(ctx context.Context, profile *models.Profile) (*models.Profile, error) {
	updatedProfile, err := s.repository.Update(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}

	return updatedProfile, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}

	return nil
}
