package job_description

import (
	"context"
	"errors"
	"fmt"

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

const jobDescriptionColumns = `id, user_id, title, company, company_website, source_url, location, employment_type, original_text, created_at, updated_at`

const (
	restrictViolation               = "23001"
	foreignKeyViolation             = "23503"
	jobApplicationsJobDescriptionFK = "job_applications_job_description_id_fkey"
)

func (r *Repository) Create(ctx context.Context, jobDescription *models.JobDescription) (*models.JobDescription, error) {
	created, err := scanJobDescription(r.db.QueryRow(ctx, `
		INSERT INTO job_descriptions (user_id, title, company, company_website, source_url, location, employment_type, original_text)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING `+jobDescriptionColumns,
		jobDescription.UserID, jobDescription.Title, jobDescription.Company, jobDescription.CompanyWebsite,
		jobDescription.SourceURL, jobDescription.Location, employmentTypeValue(jobDescription.EmploymentType), jobDescription.OriginalText,
	))
	if err != nil {
		return nil, fmt.Errorf("create job description: %w", err)
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.JobDescription, error) {
	jobDescription, err := scanJobDescription(r.db.QueryRow(ctx, `
		SELECT `+jobDescriptionColumns+`
		FROM job_descriptions
		WHERE id = $1
	`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get job description by id: %w", utility.ErrNotFound)
		}
		return nil, fmt.Errorf("get job description by id: %w", err)
	}

	return jobDescription, nil
}

func (r *Repository) List(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*models.JobDescription, error) {
	query := `SELECT ` + jobDescriptionColumns + ` FROM job_descriptions`
	args := []any{limit, offset}
	if userID != nil {
		query += ` WHERE user_id = $1`
		args = []any{*userID, limit, offset}
		query += ` ORDER BY created_at DESC, id DESC LIMIT $2 OFFSET $3`
	} else {
		query += ` ORDER BY created_at DESC, id DESC LIMIT $1 OFFSET $2`
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list job descriptions: %w", err)
	}
	defer rows.Close()

	jobDescriptions := make([]*models.JobDescription, 0)
	for rows.Next() {
		jobDescription, err := scanJobDescription(rows)
		if err != nil {
			return nil, fmt.Errorf("scan job description: %w", err)
		}
		jobDescriptions = append(jobDescriptions, jobDescription)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list job descriptions: %w", err)
	}

	return jobDescriptions, nil
}

func (r *Repository) Update(ctx context.Context, jobDescription *models.JobDescription) (*models.JobDescription, error) {
	updated, err := scanJobDescription(r.db.QueryRow(ctx, `
		UPDATE job_descriptions
		SET title = $2, company = $3, company_website = $4, source_url = $5, location = $6,
			employment_type = $7, original_text = $8, updated_at = NOW()
		WHERE id = $1
		RETURNING `+jobDescriptionColumns,
		jobDescription.ID, jobDescription.Title, jobDescription.Company, jobDescription.CompanyWebsite,
		jobDescription.SourceURL, jobDescription.Location, employmentTypeValue(jobDescription.EmploymentType), jobDescription.OriginalText,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("update job description: %w", utility.ErrNotFound)
		}
		return nil, fmt.Errorf("update job description: %w", err)
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM job_descriptions WHERE id = $1`, id)
	if err != nil {
		if hasJobApplications(err) {
			return fmt.Errorf("delete job description: %w", utility.ErrConflict)
		}
		return fmt.Errorf("delete job description: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("delete job description: %w", utility.ErrNotFound)
	}

	return nil
}

func employmentTypeValue(employmentType *models.EmploymentType) any {
	if employmentType == nil {
		return nil
	}
	return *employmentType
}

func scanJobDescription(row pgx.Row) (*models.JobDescription, error) {
	jobDescription := &models.JobDescription{}
	var employmentType pgtype.Text
	if err := row.Scan(
		&jobDescription.ID, &jobDescription.UserID, &jobDescription.Title, &jobDescription.Company,
		&jobDescription.CompanyWebsite, &jobDescription.SourceURL, &jobDescription.Location,
		&employmentType, &jobDescription.OriginalText, &jobDescription.CreatedAt, &jobDescription.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if employmentType.Valid {
		value := models.EmploymentType(employmentType.String)
		jobDescription.EmploymentType = &value
	}

	return jobDescription, nil
}

func hasJobApplications(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.ConstraintName != jobApplicationsJobDescriptionFK {
		return false
	}
	return pgErr.Code == restrictViolation || pgErr.Code == foreignKeyViolation
}
