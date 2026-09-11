package project

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

func (s *Service) Create(ctx context.Context, project *models.Project) (*models.Project, error) {
	createdProject, err := s.repository.Create(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	return createdProject, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	project, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	return project, nil
}

func (s *Service) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Project, error) {
	projects, err := s.repository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list projects by user id: %w", err)
	}

	return projects, nil
}

func (s *Service) Update(ctx context.Context, project *models.Project) (*models.Project, error) {
	updatedProject, err := s.repository.Update(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}

	return updatedProject, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete project: %w", err)
	}

	return nil
}
