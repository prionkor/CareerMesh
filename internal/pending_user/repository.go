package pending_user

import (
	"context"
	"fmt"
	"time"

	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prionkor/careermesh/models"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, email, passwordHash, tokenHash string, expiresAt time.Time) (*models.PendingUser, error) {
	// Check if email already exists
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pending_users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("check email existence: %w", err)
	}

	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// Generate UUIDv7
	id := uuid.New()

	// Insert new pending user and return the created record
	user := &models.PendingUser{}
	err = r.db.QueryRow(ctx, `
		INSERT INTO pending_users (
			id, 
			email, 
			password_hash,
			token_hash,
			expires_at
		) VALUES (
			$1, 
			$2, 
			$3, 
			$4, 
			$5
		)
		RETURNING id, email, password_hash, token_hash, expires_at, created_at
	`, id, email, passwordHash, tokenHash, expiresAt).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.TokenHash,
		&user.ExpiresAt,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create pending user: %w", err)
	}

	return user, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*models.PendingUser, error) {
	user := &models.PendingUser{}
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, token_hash, expires_at, created_at
		FROM pending_users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.TokenHash,
		&user.ExpiresAt,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("find pending user by email: %w", err)
	}

	return user, nil
}

func (r *Repository) FindByToken(ctx context.Context, tokenHash string) (*models.PendingUser, error) {
	user := &models.PendingUser{}
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, token_hash, expires_at, created_at
		FROM pending_users
		WHERE token_hash = $1
	`, tokenHash).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.TokenHash,
		&user.ExpiresAt,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("find pending user by token: %w", err)
	}

	return user, nil
}

func (r *Repository) Update(ctx context.Context, id string, email, passwordHash, tokenHash string, expiresAt time.Time) error {
	// Convert string ID to UUID
	pendingUserID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Update pending user
	_, err = r.db.Exec(ctx, `
		UPDATE pending_users
		SET email = $2, password_hash = $3, token_hash = $4, expires_at = $5, created_at = NOW()
		WHERE id = $1
	`, pendingUserID, email, passwordHash, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("update pending user: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	// Convert string ID to UUID
	pendingUserID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Delete pending user
	_, err = r.db.Exec(ctx, `
		DELETE FROM pending_users
		WHERE id = $1
	`, pendingUserID)
	if err != nil {
		return fmt.Errorf("delete pending user: %w", err)
	}

	return nil
}
