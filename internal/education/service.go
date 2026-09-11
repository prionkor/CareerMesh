package education

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

func (s *Service) Create(ctx context.Context, education *models.Education) (*models.Education, error) {
	createdEducation, err := s.repository.Create(ctx, education)
	if err != nil {
		return nil, fmt.Errorf("create education: %w", err)
	}

	return createdEducation, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Education, error) {
	education, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get education by id: %w", err)
	}

	return education, nil
}

func (s *Service) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Education, error) {
	education, err := s.repository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list education by user id: %w", err)
	}

	return education, nil
}

func (s *Service) Update(ctx context.Context, education *models.Education) (*models.Education, error) {
	updatedEducation, err := s.repository.Update(ctx, education)
	if err != nil {
		return nil, fmt.Errorf("update education: %w", err)
	}

	return updatedEducation, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete education: %w", err)
	}

	return nil
}
