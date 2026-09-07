package project

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

func (r *Repository) Create(ctx context.Context, project *models.Project) (*models.Project, error) {
	created := &models.Project{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO projects (user_id, name, description, url, repository_url, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, name, description, url, repository_url, start_date, end_date, created_at, updated_at
	`, project.UserID, project.Name, project.Description, project.URL, project.RepositoryURL, project.StartDate, project.EndDate).Scan(
		&created.ID, &created.UserID, &created.Name, &created.Description, &created.URL,
		&created.RepositoryURL, &created.StartDate, &created.EndDate, &created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	project := &models.Project{}
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, name, description, url, repository_url, start_date, end_date, created_at, updated_at
		FROM projects
		WHERE id = $1
	`, id).Scan(
		&project.ID, &project.UserID, &project.Name, &project.Description, &project.URL,
		&project.RepositoryURL, &project.StartDate, &project.EndDate, &project.CreatedAt, &project.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	return project, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Project, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, name, description, url, repository_url, start_date, end_date, created_at, updated_at
		FROM projects
		WHERE user_id = $1
		ORDER BY start_date DESC NULLS LAST, id
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list projects by user id: %w", err)
	}
	defer rows.Close()

	projects := make([]*models.Project, 0)
	for rows.Next() {
		project := &models.Project{}
		if err := rows.Scan(
			&project.ID, &project.UserID, &project.Name, &project.Description, &project.URL,
			&project.RepositoryURL, &project.StartDate, &project.EndDate, &project.CreatedAt, &project.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list projects by user id: %w", err)
	}

	return projects, nil
}

func (r *Repository) Update(ctx context.Context, project *models.Project) (*models.Project, error) {
	updated := &models.Project{}
	err := r.db.QueryRow(ctx, `
		UPDATE projects
		SET name = $2, description = $3, url = $4, repository_url = $5, start_date = $6, end_date = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, name, description, url, repository_url, start_date, end_date, created_at, updated_at
	`, project.ID, project.Name, project.Description, project.URL, project.RepositoryURL, project.StartDate, project.EndDate).Scan(
		&updated.ID, &updated.UserID, &updated.Name, &updated.Description, &updated.URL,
		&updated.RepositoryURL, &updated.StartDate, &updated.EndDate, &updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("delete project: %w", pgx.ErrNoRows)
	}

	return nil
}
