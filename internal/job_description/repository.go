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

const versionColumns = `v.id, v.job_description_id, v.version_number, v.title, v.company, v.company_website, v.source_url, v.location, v.employment_type, v.original_text, v.created_at`

const jobDescriptionSelect = `
	SELECT jd.id, jd.user_id, jd.current_version_id, jd.created_at, jd.updated_at, ` + versionColumns + `
	FROM job_descriptions jd
	JOIN job_description_versions v ON v.id = jd.current_version_id AND v.job_description_id = jd.id`

const insertVersion = `
	INSERT INTO job_description_versions (id, job_description_id, version_number, title, company, company_website, source_url, location, employment_type, original_text)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING created_at`

const (
	restrictViolation               = "23001"
	foreignKeyViolation             = "23503"
	jobApplicationsJobDescriptionFK = "job_applications_job_description_id_fkey"
)

// Create inserts the logical record and its version 1 in one transaction; the current-version FK is deferred.
func (r *Repository) Create(ctx context.Context, jobDescription *models.JobDescription, version *models.JobDescriptionVersion) (*models.JobDescription, error) {
	jobDescriptionID := uuid.NewV7()
	versionID := uuid.NewV7()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("create job description: begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	created := &models.JobDescription{ID: jobDescriptionID, UserID: jobDescription.UserID, CurrentVersionID: versionID}
	if err := tx.QueryRow(ctx, `
		INSERT INTO job_descriptions (id, user_id, current_version_id)
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at`,
		jobDescriptionID, jobDescription.UserID, versionID,
	).Scan(&created.CreatedAt, &created.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create job description: %w", err)
	}

	created.CurrentVersion, err = insertJobDescriptionVersion(ctx, tx, versionID, jobDescriptionID, 1, version)
	if err != nil {
		return nil, fmt.Errorf("create job description: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("create job description: commit transaction: %w", err)
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.JobDescription, error) {
	jobDescription, err := scanJobDescription(r.db.QueryRow(ctx, jobDescriptionSelect+`
		WHERE jd.id = $1
	`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get job description by id: %w", utility.ErrNotFound)
		}
		return nil, fmt.Errorf("get job description by id: %w", err)
	}

	return jobDescription, nil
}

func (r *Repository) GetVersion(ctx context.Context, jobDescriptionID, versionID uuid.UUID) (*models.JobDescriptionVersion, error) {
	version, err := scanJobDescriptionVersion(r.db.QueryRow(ctx, `
		SELECT `+versionColumns+`
		FROM job_description_versions v
		WHERE v.job_description_id = $1 AND v.id = $2
	`, jobDescriptionID, versionID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get job description version: %w", utility.ErrNotFound)
		}
		return nil, fmt.Errorf("get job description version: %w", err)
	}

	return version, nil
}

func (r *Repository) ListVersions(ctx context.Context, jobDescriptionID uuid.UUID) ([]*models.JobDescriptionVersion, error) {
	var exists bool
	if err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM job_descriptions WHERE id = $1)`, jobDescriptionID).Scan(&exists); err != nil {
		return nil, fmt.Errorf("list job description versions: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("list job description versions: %w", utility.ErrNotFound)
	}

	rows, err := r.db.Query(ctx, `
		SELECT `+versionColumns+`
		FROM job_description_versions v
		WHERE v.job_description_id = $1
		ORDER BY v.version_number ASC
	`, jobDescriptionID)
	if err != nil {
		return nil, fmt.Errorf("list job description versions: %w", err)
	}
	defer rows.Close()

	versions := make([]*models.JobDescriptionVersion, 0)
	for rows.Next() {
		version, err := scanJobDescriptionVersion(rows)
		if err != nil {
			return nil, fmt.Errorf("scan job description version: %w", err)
		}
		versions = append(versions, version)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list job description versions: %w", err)
	}

	return versions, nil
}

func (r *Repository) List(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*models.JobDescription, error) {
	query := jobDescriptionSelect
	args := []any{limit, offset}
	if userID != nil {
		query += ` WHERE jd.user_id = $1`
		args = []any{*userID, limit, offset}
		query += ` ORDER BY jd.created_at DESC, jd.id DESC LIMIT $2 OFFSET $3`
	} else {
		query += ` ORDER BY jd.created_at DESC, jd.id DESC LIMIT $1 OFFSET $2`
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

// CreateVersion appends an immutable version and advances the current-version pointer.
// The parent row lock serializes concurrent writers so version numbers cannot collide.
func (r *Repository) CreateVersion(ctx context.Context, jobDescriptionID uuid.UUID, version *models.JobDescriptionVersion) (*models.JobDescription, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("create job description version: begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var userID uuid.UUID
	if err := tx.QueryRow(ctx, `SELECT user_id FROM job_descriptions WHERE id = $1 FOR UPDATE`, jobDescriptionID).Scan(&userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("create job description version: %w", utility.ErrNotFound)
		}
		return nil, fmt.Errorf("create job description version: lock job description: %w", err)
	}

	var versionNumber int
	if err := tx.QueryRow(ctx, `
		SELECT COALESCE(MAX(version_number), 0) + 1
		FROM job_description_versions
		WHERE job_description_id = $1`, jobDescriptionID).Scan(&versionNumber); err != nil {
		return nil, fmt.Errorf("create job description version: next version number: %w", err)
	}

	versionID := uuid.NewV7()
	created, err := insertJobDescriptionVersion(ctx, tx, versionID, jobDescriptionID, versionNumber, version)
	if err != nil {
		return nil, fmt.Errorf("create job description version: %w", err)
	}

	jobDescription := &models.JobDescription{ID: jobDescriptionID, UserID: userID, CurrentVersionID: versionID, CurrentVersion: created}
	if err := tx.QueryRow(ctx, `
		UPDATE job_descriptions
		SET current_version_id = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING created_at, updated_at`, jobDescriptionID, versionID,
	).Scan(&jobDescription.CreatedAt, &jobDescription.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create job description version: advance current version: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("create job description version: commit transaction: %w", err)
	}

	return jobDescription, nil
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

func insertJobDescriptionVersion(ctx context.Context, tx pgx.Tx, id, jobDescriptionID uuid.UUID, versionNumber int, version *models.JobDescriptionVersion) (*models.JobDescriptionVersion, error) {
	created := &models.JobDescriptionVersion{
		ID:               id,
		JobDescriptionID: jobDescriptionID,
		VersionNumber:    versionNumber,
		Title:            version.Title,
		Company:          version.Company,
		CompanyWebsite:   version.CompanyWebsite,
		SourceURL:        version.SourceURL,
		Location:         version.Location,
		EmploymentType:   version.EmploymentType,
		OriginalText:     version.OriginalText,
	}
	if err := tx.QueryRow(ctx, insertVersion,
		id, jobDescriptionID, versionNumber, version.Title, version.Company, version.CompanyWebsite,
		version.SourceURL, version.Location, employmentTypeValue(version.EmploymentType), version.OriginalText,
	).Scan(&created.CreatedAt); err != nil {
		return nil, fmt.Errorf("insert version: %w", err)
	}

	return created, nil
}

func employmentTypeValue(employmentType *models.EmploymentType) any {
	if employmentType == nil {
		return nil
	}
	return *employmentType
}

func scanJobDescription(row pgx.Row) (*models.JobDescription, error) {
	jobDescription := &models.JobDescription{CurrentVersion: &models.JobDescriptionVersion{}}
	version := jobDescription.CurrentVersion
	var employmentType pgtype.Text
	if err := row.Scan(
		&jobDescription.ID, &jobDescription.UserID, &jobDescription.CurrentVersionID, &jobDescription.CreatedAt, &jobDescription.UpdatedAt,
		&version.ID, &version.JobDescriptionID, &version.VersionNumber, &version.Title, &version.Company,
		&version.CompanyWebsite, &version.SourceURL, &version.Location, &employmentType, &version.OriginalText, &version.CreatedAt,
	); err != nil {
		return nil, err
	}
	version.EmploymentType = employmentTypePtr(employmentType)

	return jobDescription, nil
}

func scanJobDescriptionVersion(row pgx.Row) (*models.JobDescriptionVersion, error) {
	version := &models.JobDescriptionVersion{}
	var employmentType pgtype.Text
	if err := row.Scan(
		&version.ID, &version.JobDescriptionID, &version.VersionNumber, &version.Title, &version.Company,
		&version.CompanyWebsite, &version.SourceURL, &version.Location, &employmentType, &version.OriginalText, &version.CreatedAt,
	); err != nil {
		return nil, err
	}
	version.EmploymentType = employmentTypePtr(employmentType)

	return version, nil
}

func employmentTypePtr(text pgtype.Text) *models.EmploymentType {
	if !text.Valid {
		return nil
	}
	value := models.EmploymentType(text.String)
	return &value
}

func hasJobApplications(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.ConstraintName != jobApplicationsJobDescriptionFK {
		return false
	}
	return pgErr.Code == restrictViolation || pgErr.Code == foreignKeyViolation
}
