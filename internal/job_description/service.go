package job_description

import (
	"context"
	"fmt"

	"uuid"

	"github.com/prionkor/careermesh/models"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, jobDescription *models.JobDescription) (*models.JobDescription, error) {
	created, err := s.repository.Create(ctx, jobDescription)
	if err != nil {
		return nil, fmt.Errorf("create job description: %w", err)
	}
	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.JobDescription, error) {
	jobDescription, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get job description by id: %w", err)
	}
	return jobDescription, nil
}

func (s *Service) List(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*models.JobDescription, error) {
	jobDescriptions, err := s.repository.List(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list job descriptions: %w", err)
	}
	return jobDescriptions, nil
}

func (s *Service) Update(ctx context.Context, jobDescription *models.JobDescription) (*models.JobDescription, error) {
	updated, err := s.repository.Update(ctx, jobDescription)
	if err != nil {
		return nil, fmt.Errorf("update job description: %w", err)
	}
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete job description: %w", err)
	}
	return nil
}
