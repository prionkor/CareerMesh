package certification

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

func (r *Repository) Create(ctx context.Context, certification *models.Certification) (*models.Certification, error) {
	created := &models.Certification{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO certifications (user_id, name, issuer, issue_date, expiry_date, credential_id, url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, name, issuer, issue_date, expiry_date, credential_id, url, created_at, updated_at
	`, certification.UserID, certification.Name, certification.Issuer, certification.IssueDate, certification.ExpiryDate, certification.CredentialID, certification.URL).Scan(
		&created.ID, &created.UserID, &created.Name, &created.Issuer, &created.IssueDate,
		&created.ExpiryDate, &created.CredentialID, &created.URL, &created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create certification: %w", err)
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Certification, error) {
	certification := &models.Certification{}
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, name, issuer, issue_date, expiry_date, credential_id, url, created_at, updated_at
		FROM certifications
		WHERE id = $1
	`, id).Scan(
		&certification.ID, &certification.UserID, &certification.Name, &certification.Issuer,
		&certification.IssueDate, &certification.ExpiryDate, &certification.CredentialID,
		&certification.URL, &certification.CreatedAt, &certification.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get certification by id: %w", err)
	}

	return certification, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Certification, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, name, issuer, issue_date, expiry_date, credential_id, url, created_at, updated_at
		FROM certifications
		WHERE user_id = $1
		ORDER BY issue_date DESC NULLS LAST, id
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list certifications by user id: %w", err)
	}
	defer rows.Close()

	certifications := make([]*models.Certification, 0)
	for rows.Next() {
		certification := &models.Certification{}
		if err := rows.Scan(
			&certification.ID, &certification.UserID, &certification.Name, &certification.Issuer,
			&certification.IssueDate, &certification.ExpiryDate, &certification.CredentialID,
			&certification.URL, &certification.CreatedAt, &certification.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan certification: %w", err)
		}
		certifications = append(certifications, certification)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list certifications by user id: %w", err)
	}

	return certifications, nil
}

func (r *Repository) Update(ctx context.Context, certification *models.Certification) (*models.Certification, error) {
	updated := &models.Certification{}
	err := r.db.QueryRow(ctx, `
		UPDATE certifications
		SET name = $2, issuer = $3, issue_date = $4, expiry_date = $5, credential_id = $6, url = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, name, issuer, issue_date, expiry_date, credential_id, url, created_at, updated_at
	`, certification.ID, certification.Name, certification.Issuer, certification.IssueDate, certification.ExpiryDate, certification.CredentialID, certification.URL).Scan(
		&updated.ID, &updated.UserID, &updated.Name, &updated.Issuer, &updated.IssueDate,
		&updated.ExpiryDate, &updated.CredentialID, &updated.URL, &updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update certification: %w", err)
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM certifications WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete certification: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("delete certification: %w", pgx.ErrNoRows)
	}

	return nil
}
