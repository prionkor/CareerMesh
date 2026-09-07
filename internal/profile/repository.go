package profile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

func (r *Repository) Create(ctx context.Context, profile *models.Profile) (*models.Profile, error) {
	created := &models.Profile{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO profiles (user_id, name, headline, location, phone, website, github, linkedin)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, name, headline, location, phone, website, github, linkedin, created_at, updated_at
	`, profile.UserID, profile.Name, profile.Headline, profile.Location, profile.Phone, profile.Website, profile.GitHub, profile.LinkedIn).Scan(
		&created.ID, &created.UserID, &created.Name, &created.Headline, &created.Location,
		&created.Phone, &created.Website, &created.GitHub, &created.LinkedIn,
		&created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create profile: %w", err)
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Profile, error) {
	profile := &models.Profile{}
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, name, headline, location, phone, website, github, linkedin, created_at, updated_at
		FROM profiles
		WHERE id = $1
	`, id).Scan(
		&profile.ID, &profile.UserID, &profile.Name, &profile.Headline, &profile.Location,
		&profile.Phone, &profile.Website, &profile.GitHub, &profile.LinkedIn,
		&profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get profile by id: %w", err)
	}

	return profile, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Profile, error) {
	profile := &models.Profile{}
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, name, headline, location, phone, website, github, linkedin, created_at, updated_at
		FROM profiles
		WHERE user_id = $1
	`, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Name, &profile.Headline, &profile.Location,
		&profile.Phone, &profile.Website, &profile.GitHub, &profile.LinkedIn,
		&profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get profile by user id: %w", err)
	}

	return profile, nil
}

func (r *Repository) Update(ctx context.Context, profile *models.Profile) (*models.Profile, error) {
	updated := &models.Profile{}
	err := r.db.QueryRow(ctx, `
		UPDATE profiles
		SET name = $2, headline = $3, location = $4, phone = $5, website = $6, github = $7, linkedin = $8, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, name, headline, location, phone, website, github, linkedin, created_at, updated_at
	`, profile.ID, profile.Name, profile.Headline, profile.Location, profile.Phone, profile.Website, profile.GitHub, profile.LinkedIn).Scan(
		&updated.ID, &updated.UserID, &updated.Name, &updated.Headline, &updated.Location,
		&updated.Phone, &updated.Website, &updated.GitHub, &updated.LinkedIn,
		&updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM profiles WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("delete profile: %w", pgx.ErrNoRows)
	}

	return nil
}
