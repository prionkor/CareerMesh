package language

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

func (r *Repository) Create(ctx context.Context, language *models.Language) (*models.Language, error) {
	created := &models.Language{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO languages (user_id, name, proficiency)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, name, proficiency, created_at, updated_at
	`, language.UserID, language.Name, language.Proficiency).Scan(
		&created.ID, &created.UserID, &created.Name, &created.Proficiency, &created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create language: %w", err)
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Language, error) {
	language := &models.Language{}
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, name, proficiency, created_at, updated_at
		FROM languages
		WHERE id = $1
	`, id).Scan(
		&language.ID, &language.UserID, &language.Name, &language.Proficiency, &language.CreatedAt, &language.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get language by id: %w", err)
	}

	return language, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Language, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, name, proficiency, created_at, updated_at
		FROM languages
		WHERE user_id = $1
		ORDER BY name, id
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list languages by user id: %w", err)
	}
	defer rows.Close()

	languages := make([]*models.Language, 0)
	for rows.Next() {
		language := &models.Language{}
		if err := rows.Scan(
			&language.ID, &language.UserID, &language.Name, &language.Proficiency, &language.CreatedAt, &language.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan language: %w", err)
		}
		languages = append(languages, language)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list languages by user id: %w", err)
	}

	return languages, nil
}

func (r *Repository) Update(ctx context.Context, language *models.Language) (*models.Language, error) {
	updated := &models.Language{}
	err := r.db.QueryRow(ctx, `
		UPDATE languages
		SET name = $2, proficiency = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, name, proficiency, created_at, updated_at
	`, language.ID, language.Name, language.Proficiency).Scan(
		&updated.ID, &updated.UserID, &updated.Name, &updated.Proficiency, &updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update language: %w", err)
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM languages WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete language: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("delete language: %w", pgx.ErrNoRows)
	}

	return nil
}
