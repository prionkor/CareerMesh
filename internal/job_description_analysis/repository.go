package job_description_analysis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"uuid"

	"github.com/jackc/pgx/v5"
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

const analysisColumns = `id, job_description_version_id, schema_version, prompt_version, provider, model, status, result, error, created_at, completed_at`

func (r *Repository) Create(ctx context.Context, analysis *models.JobDescriptionAnalysis) (*models.JobDescriptionAnalysis, error) {
	created, err := scanAnalysis(r.db.QueryRow(ctx, `
		INSERT INTO job_description_analyses (job_description_version_id, schema_version, prompt_version, provider, model, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING `+analysisColumns,
		analysis.JobDescriptionVersionID, analysis.SchemaVersion, analysis.PromptVersion,
		analysis.Provider, analysis.Model, analysis.Status,
	))
	if err != nil {
		return nil, fmt.Errorf("create job description analysis: %w", err)
	}
	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.JobDescriptionAnalysis, error) {
	analysis, err := scanAnalysis(r.db.QueryRow(ctx, `
		SELECT `+analysisColumns+`
		FROM job_description_analyses
		WHERE id = $1
	`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get job description analysis by id: %w", utility.ErrNotFound)
		}
		return nil, fmt.Errorf("get job description analysis by id: %w", err)
	}
	return analysis, nil
}

func (r *Repository) ListByVersion(ctx context.Context, versionID uuid.UUID) ([]*models.JobDescriptionAnalysis, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+analysisColumns+`
		FROM job_description_analyses
		WHERE job_description_version_id = $1
		ORDER BY created_at DESC, id DESC
	`, versionID)
	if err != nil {
		return nil, fmt.Errorf("list job description analyses: %w", err)
	}
	defer rows.Close()

	analyses := make([]*models.JobDescriptionAnalysis, 0)
	for rows.Next() {
		analysis, err := scanAnalysis(rows)
		if err != nil {
			return nil, fmt.Errorf("scan job description analysis: %w", err)
		}
		analyses = append(analyses, analysis)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list job description analyses: %w", err)
	}
	return analyses, nil
}

func (r *Repository) GetLatestCompletedByVersion(ctx context.Context, versionID uuid.UUID) (*models.JobDescriptionAnalysis, error) {
	analysis, err := scanAnalysis(r.db.QueryRow(ctx, `
		SELECT `+analysisColumns+`
		FROM job_description_analyses
		WHERE job_description_version_id = $1 AND status = $2
		ORDER BY completed_at DESC NULLS LAST, created_at DESC, id DESC
		LIMIT 1
	`, versionID, models.JobDescriptionAnalysisStatusCompleted))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get latest completed job description analysis: %w", utility.ErrNotFound)
		}
		return nil, fmt.Errorf("get latest completed job description analysis: %w", err)
	}
	return analysis, nil
}

func (r *Repository) Complete(ctx context.Context, id uuid.UUID, result *models.JobDescriptionAnalysisResult) (*models.JobDescriptionAnalysis, error) {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("complete job description analysis: marshal result: %w", err)
	}
	analysis, err := scanAnalysis(r.db.QueryRow(ctx, `
		UPDATE job_description_analyses
		SET status = $2, result = $3::jsonb, error = NULL, completed_at = NOW()
		WHERE id = $1 AND status = $4
		RETURNING `+analysisColumns,
		id, models.JobDescriptionAnalysisStatusCompleted, string(resultJSON), models.JobDescriptionAnalysisStatusPending,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, r.transitionMiss(ctx, id, "complete job description analysis")
		}
		return nil, fmt.Errorf("complete job description analysis: %w", err)
	}
	return analysis, nil
}

func (r *Repository) Fail(ctx context.Context, id uuid.UUID, failure string) (*models.JobDescriptionAnalysis, error) {
	analysis, err := scanAnalysis(r.db.QueryRow(ctx, `
		UPDATE job_description_analyses
		SET status = $2, result = NULL, error = $3, completed_at = NOW()
		WHERE id = $1 AND status = $4
		RETURNING `+analysisColumns,
		id, models.JobDescriptionAnalysisStatusFailed, failure, models.JobDescriptionAnalysisStatusPending,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, r.transitionMiss(ctx, id, "fail job description analysis")
		}
		return nil, fmt.Errorf("fail job description analysis: %w", err)
	}
	return analysis, nil
}

func (r *Repository) transitionMiss(ctx context.Context, id uuid.UUID, operation string) error {
	var status models.JobDescriptionAnalysisStatus
	if err := r.db.QueryRow(ctx, `SELECT status FROM job_description_analyses WHERE id = $1`, id).Scan(&status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", operation, utility.ErrNotFound)
		}
		return fmt.Errorf("%s: inspect current status: %w", operation, err)
	}
	return fmt.Errorf("%s: analysis is already finalized with status %q: %w", operation, status, utility.ErrConflict)
}

func scanAnalysis(row pgx.Row) (*models.JobDescriptionAnalysis, error) {
	analysis := &models.JobDescriptionAnalysis{}
	var resultJSON []byte
	var errorText pgtype.Text
	var completedAt pgtype.Timestamptz
	if err := row.Scan(
		&analysis.ID, &analysis.JobDescriptionVersionID, &analysis.SchemaVersion,
		&analysis.PromptVersion, &analysis.Provider, &analysis.Model, &analysis.Status,
		&resultJSON, &errorText, &analysis.CreatedAt, &completedAt,
	); err != nil {
		return nil, err
	}
	if resultJSON != nil {
		var result models.JobDescriptionAnalysisResult
		if err := json.Unmarshal(resultJSON, &result); err != nil {
			return nil, fmt.Errorf("decode job description analysis result: %w", err)
		}
		analysis.Result = &result
	}
	if errorText.Valid {
		analysis.Error = &errorText.String
	}
	if completedAt.Valid {
		analysis.CompletedAt = &completedAt.Time
	}
	return analysis, nil
}
