package user_role

import (
	"context"
	"fmt"

	"uuid"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) ReplaceAssignment(ctx context.Context, userID, roleID uuid.UUID) error {
	if err := s.repository.ReplaceAssignment(ctx, userID, roleID); err != nil {
		return fmt.Errorf("update user role: %w", err)
	}
	return nil
}
