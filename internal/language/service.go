package language

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

func (s *Service) Create(ctx context.Context, language *models.Language) (*models.Language, error) {
	createdLanguage, err := s.repository.Create(ctx, language)
	if err != nil {
		return nil, fmt.Errorf("create language: %w", err)
	}

	return createdLanguage, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Language, error) {
	language, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get language by id: %w", err)
	}

	return language, nil
}

func (s *Service) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Language, error) {
	languages, err := s.repository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list languages by user id: %w", err)
	}

	return languages, nil
}

func (s *Service) Update(ctx context.Context, language *models.Language) (*models.Language, error) {
	updatedLanguage, err := s.repository.Update(ctx, language)
	if err != nil {
		return nil, fmt.Errorf("update language: %w", err)
	}

	return updatedLanguage, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete language: %w", err)
	}

	return nil
}
