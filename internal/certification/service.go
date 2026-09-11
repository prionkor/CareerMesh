package certification

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

func (s *Service) Create(ctx context.Context, certification *models.Certification) (*models.Certification, error) {
	createdCertification, err := s.repository.Create(ctx, certification)
	if err != nil {
		return nil, fmt.Errorf("create certification: %w", err)
	}

	return createdCertification, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Certification, error) {
	certification, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get certification by id: %w", err)
	}

	return certification, nil
}

func (s *Service) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Certification, error) {
	certifications, err := s.repository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list certifications by user id: %w", err)
	}

	return certifications, nil
}

func (s *Service) Update(ctx context.Context, certification *models.Certification) (*models.Certification, error) {
	updatedCertification, err := s.repository.Update(ctx, certification)
	if err != nil {
		return nil, fmt.Errorf("update certification: %w", err)
	}

	return updatedCertification, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete certification: %w", err)
	}

	return nil
}
