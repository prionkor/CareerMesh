package user

import (
	"context"
	"fmt"

	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prionkor/careermesh/models"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	createdUser := &models.User{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash, created_at, updated_at
	`, user.Email, user.PasswordHash).Scan(
		&createdUser.ID,
		&createdUser.Email,
		&createdUser.PasswordHash,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return createdUser, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}

// Never call repository.Update directly; use service.Update instead.
// Calling the repository.Update directly can lead to security issues, such as updating passwords without proper validation.
func (r *Repository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	updatedUser := &models.User{}
	err := r.db.QueryRow(ctx, `
		UPDATE users
		SET email = $2, password_hash = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, email, password_hash, created_at, updated_at
	`, user.ID, user.Email, user.PasswordHash).Scan(
		&updatedUser.ID,
		&updatedUser.Email,
		&updatedUser.PasswordHash,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return updatedUser, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `
		DELETE FROM users
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("delete user: %w", pgx.ErrNoRows)
	}

	return nil
}

// defaultRoleName is granted to every newly created user.
const defaultRoleName = "user"

// AssignDefaultRole grants the standard 'user' role to a newly created user.
func (r *Repository) AssignDefaultRole(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, id FROM roles WHERE name = $2
	`, userID, defaultRoleName)
	if err != nil {
		return fmt.Errorf("assign default role: %w", err)
	}

	return nil
}

// GetPermissionsByUserID returns the flattened, deduplicated set of
// permission keys granted to a user through its roles.
func (r *Repository) GetPermissionsByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT p.key
		FROM user_roles ur
		JOIN role_permissions rp ON rp.role_id = ur.role_id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE ur.user_id = $1
		ORDER BY p.key
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("get permissions by user id: %w", err)
	}
	defer rows.Close()

	permissions := make([]string, 0)
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		permissions = append(permissions, key)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get permissions by user id: %w", err)
	}

	return permissions, nil
}
