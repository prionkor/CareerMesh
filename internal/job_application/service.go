package job_application

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

func (s *Service) Create(ctx context.Context, application *models.JobApplication) (*models.JobApplication, error) {
	created, err := s.repository.Create(ctx, application)
	if err != nil {
		return nil, fmt.Errorf("create job application: %w", err)
	}
	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.JobApplication, error) {
	application, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get job application by id: %w", err)
	}
	return application, nil
}

func (s *Service) List(ctx context.Context, userID, jobDescriptionID *uuid.UUID, limit, offset int) ([]*models.JobApplication, error) {
	applications, err := s.repository.List(ctx, userID, jobDescriptionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list job applications: %w", err)
	}
	return applications, nil
}

func (s *Service) Update(ctx context.Context, application *models.JobApplication) (*models.JobApplication, error) {
	updated, err := s.repository.Update(ctx, application)
	if err != nil {
		return nil, fmt.Errorf("update job application: %w", err)
	}
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete job application: %w", err)
	}
	return nil
}
