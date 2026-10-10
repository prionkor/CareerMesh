-- +goose Up

CREATE TYPE job_application_status AS ENUM (
    'draft',
    'applied',
    'screening',
    'interviewing',
    'offer',
    'accepted',
    'rejected',
    'withdrawn'
);

CREATE TABLE job_applications (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    job_description_id UUID NOT NULL REFERENCES job_descriptions(id) ON DELETE RESTRICT,
    status job_application_status NOT NULL DEFAULT 'draft',
    applied_at DATE NULL,
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_job_applications_user_id ON job_applications(user_id);
CREATE INDEX idx_job_applications_job_description_id ON job_applications(job_description_id);

-- +goose Down

DROP TABLE job_applications;

DROP TYPE job_application_status;
