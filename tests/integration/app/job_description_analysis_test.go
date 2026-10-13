package app_test

import (
	"context"
	"errors"
	"testing"

	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prionkor/careermesh/internal/job_description_analysis"
	"github.com/prionkor/careermesh/internal/utility"
	"github.com/prionkor/careermesh/models"
)

func TestJobDescriptionAnalysisRunLifecycle(t *testing.T) {
	ctx := context.Background()
	pool := testDB(t)
	service := job_description_analysis.NewService(job_description_analysis.NewRepository(pool))
	versionID := createJobDescriptionVersionForAnalysis(t, pool)
	otherVersionID := createJobDescriptionVersionForAnalysis(t, pool)

	first, err := service.CreateRun(ctx, versionID, "prompt-v1", "provider-a", "model-a")
	if err != nil {
		t.Fatalf("create first run: %v", err)
	}
	second, err := service.CreateRun(ctx, versionID, "prompt-v2", "provider-a", "model-a")
	if err != nil {
		t.Fatalf("create second run: %v", err)
	}
	pending, err := service.CreateRun(ctx, versionID, "prompt-v2", "provider-b", "model-b")
	if err != nil {
		t.Fatalf("create pending run: %v", err)
	}
	failed, err := service.CreateRun(ctx, otherVersionID, "prompt-v1", "provider-a", "model-a")
	if err != nil {
		t.Fatalf("create failure-test run: %v", err)
	}
	for _, run := range []*models.JobDescriptionAnalysis{first, second, pending, failed} {
		if run.Status != models.JobDescriptionAnalysisStatusPending || run.SchemaVersion != models.JobDescriptionAnalysisSchemaVersion {
			t.Fatalf("unexpected new run state: %+v", run)
		}
	}
	if first.ID == second.ID || first.ID == pending.ID {
		t.Fatal("new runs must have distinct IDs")
	}

	resultOne := analysisResult("first result")
	completedFirst, err := service.Complete(ctx, first.ID, &resultOne)
	if err != nil {
		t.Fatalf("complete first run: %v", err)
	}
	resultTwo := analysisResult("second result")
	completedSecond, err := service.Complete(ctx, second.ID, &resultTwo)
	if err != nil {
		t.Fatalf("complete second run: %v", err)
	}
	if completedFirst.Status != models.JobDescriptionAnalysisStatusCompleted || completedFirst.Result == nil ||
		completedFirst.Result.Requirements[0].Name != "first result" || completedFirst.Error != nil || completedFirst.CompletedAt == nil {
		t.Fatalf("unexpected completed run: %+v", completedFirst)
	}
	if completedSecond.Status != models.JobDescriptionAnalysisStatusCompleted || completedSecond.Result == nil ||
		completedSecond.Result.Requirements[0].Name != "second result" {
		t.Fatalf("unexpected second completed run: %+v", completedSecond)
	}

	_, err = pool.Exec(ctx, `
		UPDATE job_description_analyses
		SET created_at = CASE id
			WHEN $1 THEN NOW() - INTERVAL '2 minutes'
			WHEN $2 THEN NOW() - INTERVAL '3 minutes'
			WHEN $3 THEN NOW()
		END,
		completed_at = CASE id
			WHEN $1 THEN NOW() - INTERVAL '1 minute'
			WHEN $2 THEN NOW()
			ELSE completed_at
		END
		WHERE id IN ($1, $2, $3)
	`, first.ID, second.ID, pending.ID)
	if err != nil {
		t.Fatalf("set deterministic run timestamps: %v", err)
	}

	listed, err := service.ListByVersion(ctx, versionID)
	if err != nil {
		t.Fatalf("list runs: %v", err)
	}
	if len(listed) != 3 || listed[0].ID != pending.ID || listed[1].ID != first.ID || listed[2].ID != second.ID {
		t.Fatalf("unexpected newest-first run list: %+v", listed)
	}
	latest, err := service.GetLatestCompletedByVersion(ctx, versionID)
	if err != nil {
		t.Fatalf("get latest completed run: %v", err)
	}
	if latest.ID != second.ID || latest.Status != models.JobDescriptionAnalysisStatusCompleted {
		t.Fatalf("latest completed run = %s, want %s", latest.ID, second.ID)
	}

	changedResult := analysisResult("replacement must not persist")
	if _, err := service.Complete(ctx, first.ID, &changedResult); !errors.Is(err, utility.ErrConflict) {
		t.Fatalf("complete finalized run error = %v, want ErrConflict", err)
	}
	if _, err := service.Fail(ctx, first.ID, "replacement failure"); !errors.Is(err, utility.ErrConflict) {
		t.Fatalf("fail completed run error = %v, want ErrConflict", err)
	}
	unchanged, err := service.GetByID(ctx, first.ID)
	if err != nil {
		t.Fatalf("get completed run: %v", err)
	}
	if unchanged.Result == nil || unchanged.Result.Requirements[0].Name != "first result" {
		t.Fatalf("completed result was overwritten: %+v", unchanged.Result)
	}

	failedRun, err := service.Fail(ctx, failed.ID, "provider unavailable")
	if err != nil {
		t.Fatalf("fail run: %v", err)
	}
	if failedRun.Status != models.JobDescriptionAnalysisStatusFailed || failedRun.Error == nil ||
		*failedRun.Error != "provider unavailable" || failedRun.CompletedAt == nil || failedRun.Result != nil {
		t.Fatalf("unexpected failed run: %+v", failedRun)
	}
	if _, err := service.Complete(ctx, failed.ID, &resultOne); !errors.Is(err, utility.ErrConflict) {
		t.Fatalf("complete failed run error = %v, want ErrConflict", err)
	}
	if _, err := service.GetLatestCompletedByVersion(ctx, otherVersionID); !errors.Is(err, utility.ErrNotFound) {
		t.Fatalf("latest completed for failed-only version error = %v, want ErrNotFound", err)
	}

	retry, err := service.CreateRun(ctx, otherVersionID, "prompt-v2", "provider-a", "model-a")
	if err != nil {
		t.Fatalf("create retry run: %v", err)
	}
	if retry.ID == failed.ID || retry.Status != models.JobDescriptionAnalysisStatusPending {
		t.Fatalf("retry did not create a distinct pending run: %+v", retry)
	}
	otherRuns, err := service.ListByVersion(ctx, otherVersionID)
	if err != nil {
		t.Fatalf("list retry runs: %v", err)
	}
	if len(otherRuns) != 2 {
		t.Fatalf("got %d runs after retry, want 2", len(otherRuns))
	}

	missingID := uuid.NewV7()
	if _, err := service.GetByID(ctx, missingID); !errors.Is(err, utility.ErrNotFound) {
		t.Fatalf("get missing run error = %v, want ErrNotFound", err)
	}
	if _, err := service.Complete(ctx, missingID, &resultOne); !errors.Is(err, utility.ErrNotFound) {
		t.Fatalf("complete missing run error = %v, want ErrNotFound", err)
	}
	if _, err := service.Fail(ctx, missingID, "failure"); !errors.Is(err, utility.ErrNotFound) {
		t.Fatalf("fail missing run error = %v, want ErrNotFound", err)
	}
}

func createJobDescriptionVersionForAnalysis(t *testing.T, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	userID := uuid.NewV7()
	jobDescriptionID := uuid.NewV7()
	versionID := uuid.NewV7()

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin analysis fixture transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3)`,
		userID, "analysis-"+userID.String()+"@example.test", "test-hash"); err != nil {
		t.Fatalf("insert analysis fixture user: %v", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO job_descriptions (id, user_id, current_version_id)
		VALUES ($1, $2, $3)
	`, jobDescriptionID, userID, versionID); err != nil {
		t.Fatalf("insert analysis fixture job description: %v", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO job_description_versions (id, job_description_id, version_number, title, company, original_text)
		VALUES ($1, $2, 1, 'Engineer', 'Example', 'Original posting')
	`, versionID, jobDescriptionID); err != nil {
		t.Fatalf("insert analysis fixture version: %v", err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatalf("commit analysis fixture: %v", err)
	}

	t.Cleanup(func() {
		if _, err := pool.Exec(ctx, `DELETE FROM job_description_analyses WHERE job_description_version_id = $1`, versionID); err != nil {
			t.Errorf("delete analysis fixture runs: %v", err)
		}
		if _, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID); err != nil {
			t.Errorf("delete analysis fixture user: %v", err)
		}
	})
	return versionID
}

func analysisResult(name string) models.JobDescriptionAnalysisResult {
	return models.JobDescriptionAnalysisResult{
		Requirements: []models.JobDescriptionAnalysisRequirement{{
			Name:      name,
			Category:  models.JobDescriptionAnalysisRequirementCategorySkill,
			Priority:  models.JobDescriptionAnalysisRequirementPriorityRequired,
			Assertion: models.JobDescriptionAnalysisAssertionExplicit,
			Weight:    0.95,
			Evidence:  []models.JobDescriptionAnalysisEvidence{{Text: "Experience operating Kubernetes in production"}},
		}},
		Responsibilities: []models.JobDescriptionAnalysisResponsibility{},
	}
}
