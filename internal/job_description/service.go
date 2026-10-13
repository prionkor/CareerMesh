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

func (s *Service) Create(ctx context.Context, jobDescription *models.JobDescription, version *models.JobDescriptionVersion) (*models.JobDescription, error) {
	created, err := s.repository.Create(ctx, jobDescription, version)
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

func (s *Service) GetVersion(ctx context.Context, jobDescriptionID, versionID uuid.UUID) (*models.JobDescriptionVersion, error) {
	version, err := s.repository.GetVersion(ctx, jobDescriptionID, versionID)
	if err != nil {
		return nil, fmt.Errorf("get job description version: %w", err)
	}
	return version, nil
}

func (s *Service) ListVersions(ctx context.Context, jobDescriptionID uuid.UUID) ([]*models.JobDescriptionVersion, error) {
	versions, err := s.repository.ListVersions(ctx, jobDescriptionID)
	if err != nil {
		return nil, fmt.Errorf("list job description versions: %w", err)
	}
	return versions, nil
}

func (s *Service) List(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*models.JobDescription, error) {
	jobDescriptions, err := s.repository.List(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list job descriptions: %w", err)
	}
	return jobDescriptions, nil
}

func (s *Service) CreateVersion(ctx context.Context, jobDescriptionID uuid.UUID, version *models.JobDescriptionVersion) (*models.JobDescription, error) {
	created, err := s.repository.CreateVersion(ctx, jobDescriptionID, version)
	if err != nil {
		return nil, fmt.Errorf("create job description version: %w", err)
	}
	return created, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete job description: %w", err)
	}
	return nil
}
