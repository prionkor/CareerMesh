package skill

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

func (s *Service) Create(ctx context.Context, skill *models.Skill) (*models.Skill, error) {
	createdSkill, err := s.repository.Create(ctx, skill)
	if err != nil {
		return nil, fmt.Errorf("create skill: %w", err)
	}

	return createdSkill, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Skill, error) {
	skill, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get skill by id: %w", err)
	}

	return skill, nil
}

func (s *Service) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Skill, error) {
	skills, err := s.repository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list skills by user id: %w", err)
	}

	return skills, nil
}

func (s *Service) Update(ctx context.Context, skill *models.Skill) (*models.Skill, error) {
	updatedSkill, err := s.repository.Update(ctx, skill)
	if err != nil {
		return nil, fmt.Errorf("update skill: %w", err)
	}

	return updatedSkill, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete skill: %w", err)
	}

	return nil
}
