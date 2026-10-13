-- migrations/00018_create_job_description_versions.sql
-- +goose Up

CREATE TABLE job_description_versions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    job_description_id UUID NOT NULL
        REFERENCES job_descriptions(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL CHECK (version_number > 0),

    title TEXT NOT NULL,
    company TEXT NOT NULL,
    company_website TEXT NOT NULL DEFAULT '',
    source_url TEXT NOT NULL DEFAULT '',
    location TEXT NOT NULL DEFAULT '',
    employment_type employment_type,
    original_text TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_job_description_versions_number
        UNIQUE (job_description_id, version_number),

    CONSTRAINT uq_job_description_versions_parent_id
        UNIQUE (job_description_id, id)
);

CREATE INDEX idx_job_description_versions_job_description_id
    ON job_description_versions(job_description_id);

ALTER TABLE job_descriptions
    ADD CONSTRAINT fk_job_descriptions_current_version
    FOREIGN KEY (id, current_version_id)
    REFERENCES job_description_versions(job_description_id, id)
    DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE job_descriptions
    ALTER COLUMN current_version_id SET NOT NULL;

-- +goose Down

ALTER TABLE job_descriptions
    DROP CONSTRAINT fk_job_descriptions_current_version;

ALTER TABLE job_descriptions
    ALTER COLUMN current_version_id DROP NOT NULL;

DROP TABLE job_description_versions;