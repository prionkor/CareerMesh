package education

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

func (r *Repository) Create(ctx context.Context, education *models.Education) (*models.Education, error) {
	created := &models.Education{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO education (user_id, institution, degree, field, start_date, end_date, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, institution, degree, field, start_date, end_date, description, created_at, updated_at
	`, education.UserID, education.Institution, education.Degree, education.Field, education.StartDate, education.EndDate, education.Description).Scan(
		&created.ID, &created.UserID, &created.Institution, &created.Degree, &created.Field,
		&created.StartDate, &created.EndDate, &created.Description, &created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create education: %w", err)
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Education, error) {
	education := &models.Education{}
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, institution, degree, field, start_date, end_date, description, created_at, updated_at
		FROM education
		WHERE id = $1
	`, id).Scan(
		&education.ID, &education.UserID, &education.Institution, &education.Degree, &education.Field,
		&education.StartDate, &education.EndDate, &education.Description, &education.CreatedAt, &education.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get education by id: %w", err)
	}

	return education, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Education, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, institution, degree, field, start_date, end_date, description, created_at, updated_at
		FROM education
		WHERE user_id = $1
		ORDER BY start_date DESC NULLS LAST, id
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list education by user id: %w", err)
	}
	defer rows.Close()

	entries := make([]*models.Education, 0)
	for rows.Next() {
		entry := &models.Education{}
		if err := rows.Scan(
			&entry.ID, &entry.UserID, &entry.Institution, &entry.Degree, &entry.Field,
			&entry.StartDate, &entry.EndDate, &entry.Description, &entry.CreatedAt, &entry.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan education: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list education by user id: %w", err)
	}

	return entries, nil
}

func (r *Repository) Update(ctx context.Context, education *models.Education) (*models.Education, error) {
	updated := &models.Education{}
	err := r.db.QueryRow(ctx, `
		UPDATE education
		SET institution = $2, degree = $3, field = $4, start_date = $5, end_date = $6, description = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, institution, degree, field, start_date, end_date, description, created_at, updated_at
	`, education.ID, education.Institution, education.Degree, education.Field, education.StartDate, education.EndDate, education.Description).Scan(
		&updated.ID, &updated.UserID, &updated.Institution, &updated.Degree, &updated.Field,
		&updated.StartDate, &updated.EndDate, &updated.Description, &updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update education: %w", err)
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM education WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete education: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("delete education: %w", pgx.ErrNoRows)
	}

	return nil
}
