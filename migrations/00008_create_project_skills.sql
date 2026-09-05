-- migrations/00008_create_project_skills.sql

-- +goose Up

CREATE TABLE project_skills (
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, skill_id)
);

CREATE INDEX idx_project_skills_skill_id
    ON project_skills(skill_id);

-- +goose Down

DROP TABLE project_skills;