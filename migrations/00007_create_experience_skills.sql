-- migrations/00007_create_experience_skills.sql

-- +goose Up

CREATE TABLE experience_skills (
    experience_id UUID NOT NULL REFERENCES experiences(id) ON DELETE CASCADE,
    skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    PRIMARY KEY (experience_id, skill_id)
);

CREATE INDEX idx_experience_skills_skill_id
    ON experience_skills(skill_id);

-- +goose Down

DROP TABLE experience_skills;