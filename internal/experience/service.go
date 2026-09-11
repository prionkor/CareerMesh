package experience

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

func (s *Service) Create(ctx context.Context, experience *models.Experience) (*models.Experience, error) {
	createdExperience, err := s.repository.Create(ctx, experience)
	if err != nil {
		return nil, fmt.Errorf("create experience: %w", err)
	}

	return createdExperience, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Experience, error) {
	experience, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get experience by id: %w", err)
	}

	return experience, nil
}

func (s *Service) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Experience, error) {
	experiences, err := s.repository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list experiences by user id: %w", err)
	}

	return experiences, nil
}

func (s *Service) Update(ctx context.Context, experience *models.Experience) (*models.Experience, error) {
	updatedExperience, err := s.repository.Update(ctx, experience)
	if err != nil {
		return nil, fmt.Errorf("update experience: %w", err)
	}

	return updatedExperience, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete experience: %w", err)
	}

	return nil
}
