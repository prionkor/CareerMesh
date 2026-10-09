package user_role

import (
	"context"
	"fmt"

	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ReplaceAssignment makes the requested role the user's only assigned role.
func (r *Repository) ReplaceAssignment(ctx context.Context, userID, roleID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO user_roles (user_id, role_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id)
        DO UPDATE SET role_id = EXCLUDED.role_id
    `, userID, roleID)
	if err != nil {
		return fmt.Errorf("replace role assignment: %w", err)
	}

	return nil
}
