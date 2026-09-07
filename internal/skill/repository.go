package skill

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

func (r *Repository) Create(ctx context.Context, skill *models.Skill) (*models.Skill, error) {
	created := &models.Skill{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO skills (user_id, name, category)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, name, category, created_at, updated_at
	`, skill.UserID, skill.Name, skill.Category).Scan(
		&created.ID, &created.UserID, &created.Name, &created.Category, &created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create skill: %w", err)
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Skill, error) {
	skill := &models.Skill{}
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, name, category, created_at, updated_at
		FROM skills
		WHERE id = $1
	`, id).Scan(
		&skill.ID, &skill.UserID, &skill.Name, &skill.Category, &skill.CreatedAt, &skill.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get skill by id: %w", err)
	}

	return skill, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Skill, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, name, category, created_at, updated_at
		FROM skills
		WHERE user_id = $1
		ORDER BY name, id
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list skills by user id: %w", err)
	}
	defer rows.Close()

	skills := make([]*models.Skill, 0)
	for rows.Next() {
		skill := &models.Skill{}
		if err := rows.Scan(
			&skill.ID, &skill.UserID, &skill.Name, &skill.Category, &skill.CreatedAt, &skill.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan skill: %w", err)
		}
		skills = append(skills, skill)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list skills by user id: %w", err)
	}

	return skills, nil
}

func (r *Repository) Update(ctx context.Context, skill *models.Skill) (*models.Skill, error) {
	updated := &models.Skill{}
	err := r.db.QueryRow(ctx, `
		UPDATE skills
		SET name = $2, category = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, name, category, created_at, updated_at
	`, skill.ID, skill.Name, skill.Category).Scan(
		&updated.ID, &updated.UserID, &updated.Name, &updated.Category, &updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update skill: %w", err)
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM skills WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete skill: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("delete skill: %w", pgx.ErrNoRows)
	}

	return nil
}
