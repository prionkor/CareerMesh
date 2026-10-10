package user_role

import (
	"context"
	"errors"
	"fmt"

	"uuid"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prionkor/careermesh/internal/utility"
)

type Repository struct {
	db *pgxpool.Pool
}

const (
	foreignKeyViolation = "23503"
	userFK              = "user_roles_user_id_fkey"
	roleFK              = "user_roles_role_id_fkey"
)

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
		if isMissingReference(err) {
			return fmt.Errorf("replace role assignment: %w", utility.ErrNotFound)
		}
		return fmt.Errorf("replace role assignment: %w", err)
	}

	return nil
}

func isMissingReference(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != foreignKeyViolation {
		return false
	}
	return pgErr.ConstraintName == userFK || pgErr.ConstraintName == roleFK
}
