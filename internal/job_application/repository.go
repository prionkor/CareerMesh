package job_application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prionkor/careermesh/internal/utility"
	"github.com/prionkor/careermesh/models"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

const jobApplicationColumns = `id, user_id, job_description_id, status, applied_at, notes, created_at, updated_at`

const (
	foreignKeyViolation = "23503"
	jobDescriptionFK    = "job_applications_job_description_id_fkey"
)

func (r *Repository) Create(ctx context.Context, application *models.JobApplication) (*models.JobApplication, error) {
	created, err := scanJobApplication(r.db.QueryRow(ctx, `
		INSERT INTO job_applications (user_id, job_description_id, status, applied_at, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING `+jobApplicationColumns,
		application.UserID, application.JobDescriptionID, application.Status,
		appliedAtValue(application.AppliedAt), application.Notes,
	))
	if err != nil {
		if isMissingJobDescription(err) {
			return nil, fmt.Errorf("create job application: %w", utility.ErrNotFound)
		}
		return nil, fmt.Errorf("create job application: %w", err)
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.JobApplication, error) {
	application, err := scanJobApplication(r.db.QueryRow(ctx, `
		SELECT `+jobApplicationColumns+`
		FROM job_applications
		WHERE id = $1
	`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get job application by id: %w", utility.ErrNotFound)
		}
		return nil, fmt.Errorf("get job application by id: %w", err)
	}

	return application, nil
}

// List filters by owner and/or job description when the respective argument is non-nil.
func (r *Repository) List(ctx context.Context, userID, jobDescriptionID *uuid.UUID, limit, offset int) ([]*models.JobApplication, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+jobApplicationColumns+`
		FROM job_applications
		WHERE ($1::uuid IS NULL OR user_id = $1)
			AND ($2::uuid IS NULL OR job_description_id = $2)
		ORDER BY created_at DESC, id DESC
		LIMIT $3 OFFSET $4
	`, userID, jobDescriptionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list job applications: %w", err)
	}
	defer rows.Close()

	applications := make([]*models.JobApplication, 0)
	for rows.Next() {
		application, err := scanJobApplication(rows)
		if err != nil {
			return nil, fmt.Errorf("scan job application: %w", err)
		}
		applications = append(applications, application)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list job applications: %w", err)
	}

	return applications, nil
}

func (r *Repository) Update(ctx context.Context, application *models.JobApplication) (*models.JobApplication, error) {
	updated, err := scanJobApplication(r.db.QueryRow(ctx, `
		UPDATE job_applications
		SET status = $2, applied_at = $3, notes = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING `+jobApplicationColumns,
		application.ID, application.Status, appliedAtValue(application.AppliedAt), application.Notes,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("update job application: %w", utility.ErrNotFound)
		}
		return nil, fmt.Errorf("update job application: %w", err)
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM job_applications WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete job application: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("delete job application: %w", utility.ErrNotFound)
	}

	return nil
}

func appliedAtValue(appliedAt *models.Date) any {
	if appliedAt == nil {
		return nil
	}
	return appliedAt.Time
}

func scanJobApplication(row pgx.Row) (*models.JobApplication, error) {
	application := &models.JobApplication{}
	var appliedAt pgtype.Date
	if err := row.Scan(
		&application.ID, &application.UserID, &application.JobDescriptionID, &application.Status,
		&appliedAt, &application.Notes, &application.CreatedAt, &application.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if appliedAt.Valid {
		application.AppliedAt = &models.Date{Time: appliedAt.Time.UTC().Truncate(24 * time.Hour)}
	}

	return application, nil
}

func isMissingJobDescription(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) &&
		pgErr.Code == foreignKeyViolation && pgErr.ConstraintName == jobDescriptionFK
}
