-- migrations/00019_create_job_description_analyses.sql
-- +goose Up
CREATE TABLE job_description_analyses (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    job_description_version_id UUID NOT NULL,
    schema_version TEXT NOT NULL,
    prompt_version TEXT NOT NULL,
    provider TEXT NOT NULL,
    model TEXT NOT NULL,
    status TEXT NOT NULL,
    result JSONB,
    error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

-- Supports listing analysis runs for a job description version, newest first.
CREATE INDEX idx_job_description_analyses_version_created_at
    ON job_description_analyses (job_description_version_id, created_at DESC);

-- Supports finding the latest completed analysis for a job description version.
CREATE INDEX idx_job_description_analyses_version_completed_at
    ON job_description_analyses (job_description_version_id, completed_at DESC)
    WHERE status = 'completed';

-- +goose Down
DROP TABLE job_description_analyses;