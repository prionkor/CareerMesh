package job_description_analysis

import (
	"context"
	"errors"
	"testing"

	"uuid"

	"github.com/prionkor/careermesh/internal/utility"
	"github.com/prionkor/careermesh/models"
)

func TestCreateRunRejectsInvalidMetadata(t *testing.T) {
	service := NewService(nil)
	tests := []struct {
		name          string
		versionID     uuid.UUID
		promptVersion string
		provider      string
		modelName     string
	}{
		{name: "zero version id", promptVersion: "p1", provider: "provider", modelName: "model"},
		{name: "missing prompt version", versionID: uuid.NewV7(), provider: "provider", modelName: "model"},
		{name: "missing provider", versionID: uuid.NewV7(), promptVersion: "p1", modelName: "model"},
		{name: "missing model", versionID: uuid.NewV7(), promptVersion: "p1", provider: "provider"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := service.CreateRun(context.Background(), test.versionID, test.promptVersion, test.provider, test.modelName)
			if !errors.Is(err, utility.ErrInvalidInput) {
				t.Fatalf("CreateRun() error = %v, want ErrInvalidInput", err)
			}
		})
	}
}

func TestCompleteRejectsInvalidResult(t *testing.T) {
	service := NewService(nil)
	tests := []struct {
		name   string
		result *models.JobDescriptionAnalysisResult
	}{
		{name: "nil result"},
		{name: "missing arrays", result: &models.JobDescriptionAnalysisResult{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := service.Complete(context.Background(), uuid.NewV7(), test.result)
			if !errors.Is(err, utility.ErrInvalidInput) {
				t.Fatalf("Complete() error = %v, want ErrInvalidInput", err)
			}
		})
	}
}

func TestFailRejectsEmptyFailure(t *testing.T) {
	service := NewService(nil)
	if _, err := service.Fail(context.Background(), uuid.NewV7(), " \t "); !errors.Is(err, utility.ErrInvalidInput) {
		t.Fatalf("Fail() error = %v, want ErrInvalidInput", err)
	}
}
