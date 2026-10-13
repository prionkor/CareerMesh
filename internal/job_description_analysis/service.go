package job_description_analysis

import (
	"context"
	"fmt"
	"strings"

	"uuid"

	"github.com/prionkor/careermesh/internal/utility"
	"github.com/prionkor/careermesh/models"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateRun(ctx context.Context, versionID uuid.UUID, promptVersion, provider, modelName string) (*models.JobDescriptionAnalysis, error) {
	if versionID == (uuid.UUID{}) || strings.TrimSpace(promptVersion) == "" ||
		strings.TrimSpace(provider) == "" || strings.TrimSpace(modelName) == "" {
		return nil, fmt.Errorf("create job description analysis: %w", utility.ErrInvalidInput)
	}
	analysis, err := s.repository.Create(ctx, &models.JobDescriptionAnalysis{
		JobDescriptionVersionID: versionID,
		SchemaVersion:           models.JobDescriptionAnalysisSchemaVersion,
		PromptVersion:           promptVersion,
		Provider:                provider,
		Model:                   modelName,
		Status:                  models.JobDescriptionAnalysisStatusPending,
	})
	if err != nil {
		return nil, fmt.Errorf("create job description analysis: %w", err)
	}
	return analysis, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.JobDescriptionAnalysis, error) {
	analysis, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get job description analysis by id: %w", err)
	}
	return analysis, nil
}

func (s *Service) ListByVersion(ctx context.Context, versionID uuid.UUID) ([]*models.JobDescriptionAnalysis, error) {
	analyses, err := s.repository.ListByVersion(ctx, versionID)
	if err != nil {
		return nil, fmt.Errorf("list job description analyses: %w", err)
	}
	return analyses, nil
}

func (s *Service) GetLatestCompletedByVersion(ctx context.Context, versionID uuid.UUID) (*models.JobDescriptionAnalysis, error) {
	analysis, err := s.repository.GetLatestCompletedByVersion(ctx, versionID)
	if err != nil {
		return nil, fmt.Errorf("get latest completed job description analysis: %w", err)
	}
	return analysis, nil
}

func (s *Service) Complete(ctx context.Context, id uuid.UUID, result *models.JobDescriptionAnalysisResult) (*models.JobDescriptionAnalysis, error) {
	if result == nil {
		return nil, fmt.Errorf("complete job description analysis: %w", utility.ErrInvalidInput)
	}
	completedResult := *result
	completedResult.SchemaVersion = models.JobDescriptionAnalysisSchemaVersion
	if err := completedResult.Validate(); err != nil {
		return nil, fmt.Errorf("complete job description analysis: %w: %v", utility.ErrInvalidInput, err)
	}
	analysis, err := s.repository.Complete(ctx, id, &completedResult)
	if err != nil {
		return nil, fmt.Errorf("complete job description analysis: %w", err)
	}
	return analysis, nil
}

func (s *Service) Fail(ctx context.Context, id uuid.UUID, failure string) (*models.JobDescriptionAnalysis, error) {
	if strings.TrimSpace(failure) == "" {
		return nil, fmt.Errorf("fail job description analysis: %w", utility.ErrInvalidInput)
	}
	analysis, err := s.repository.Fail(ctx, id, failure)
	if err != nil {
		return nil, fmt.Errorf("fail job description analysis: %w", err)
	}
	return analysis, nil
}
