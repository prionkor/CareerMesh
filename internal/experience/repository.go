package experience

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

func (r *Repository) Create(ctx context.Context, experience *models.Experience) (*models.Experience, error) {
	created := &models.Experience{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO experiences (user_id, company, website, title, employment_type, location, start_date, end_date, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, company, website, title, employment_type, location, start_date, end_date, description, created_at, updated_at
	`, experience.UserID, experience.Company, experience.Website, experience.Title, experience.EmploymentType, experience.Location, experience.StartDate, experience.EndDate, experience.Description).Scan(
		&created.ID, &created.UserID, &created.Company, &created.Website, &created.Title,
		&created.EmploymentType, &created.Location, &created.StartDate, &created.EndDate,
		&created.Description, &created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create experience: %w", err)
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Experience, error) {
	experience := &models.Experience{}
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, company, website, title, employment_type, location, start_date, end_date, description, created_at, updated_at
		FROM experiences
		WHERE id = $1
	`, id).Scan(
		&experience.ID, &experience.UserID, &experience.Company, &experience.Website, &experience.Title,
		&experience.EmploymentType, &experience.Location, &experience.StartDate, &experience.EndDate,
		&experience.Description, &experience.CreatedAt, &experience.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get experience by id: %w", err)
	}

	return experience, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Experience, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, company, website, title, employment_type, location, start_date, end_date, description, created_at, updated_at
		FROM experiences
		WHERE user_id = $1
		ORDER BY start_date DESC, id
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list experiences by user id: %w", err)
	}
	defer rows.Close()

	experiences := make([]*models.Experience, 0)
	for rows.Next() {
		experience := &models.Experience{}
		if err := rows.Scan(
			&experience.ID, &experience.UserID, &experience.Company, &experience.Website, &experience.Title,
			&experience.EmploymentType, &experience.Location, &experience.StartDate, &experience.EndDate,
			&experience.Description, &experience.CreatedAt, &experience.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan experience: %w", err)
		}
		experiences = append(experiences, experience)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list experiences by user id: %w", err)
	}

	return experiences, nil
}

func (r *Repository) Update(ctx context.Context, experience *models.Experience) (*models.Experience, error) {
	updated := &models.Experience{}
	err := r.db.QueryRow(ctx, `
		UPDATE experiences
		SET company = $2, website = $3, title = $4, employment_type = $5, location = $6, start_date = $7, end_date = $8, description = $9, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, company, website, title, employment_type, location, start_date, end_date, description, created_at, updated_at
	`, experience.ID, experience.Company, experience.Website, experience.Title, experience.EmploymentType, experience.Location, experience.StartDate, experience.EndDate, experience.Description).Scan(
		&updated.ID, &updated.UserID, &updated.Company, &updated.Website, &updated.Title,
		&updated.EmploymentType, &updated.Location, &updated.StartDate, &updated.EndDate,
		&updated.Description, &updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update experience: %w", err)
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM experiences WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete experience: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("delete experience: %w", pgx.ErrNoRows)
	}

	return nil
}
