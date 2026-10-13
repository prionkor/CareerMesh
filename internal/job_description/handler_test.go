package job_description

import (
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/prionkor/careermesh/models"
)

func TestJobDescriptionVersionFromRequest(t *testing.T) {
	employmentType := models.EmploymentTypeFullTime
	version := jobDescriptionVersionFromRequest(JobDescriptionRequest{
		Title: "Senior Engineer", Company: "Example Corp", CompanyWebsite: "https://example.com",
		SourceURL: "https://jobs.example.com/role", Location: "Remote",
		EmploymentType: &employmentType, OriginalText: "Original posting",
	})

	if version.Title != "Senior Engineer" || version.Company != "Example Corp" ||
		version.CompanyWebsite != "https://example.com" || version.SourceURL != "https://jobs.example.com/role" ||
		version.Location != "Remote" || version.EmploymentType == nil || *version.EmploymentType != employmentType ||
		version.OriginalText != "Original posting" {
		t.Fatalf("jobDescriptionVersionFromRequest() = %+v", version)
	}
}

func TestHasJobApplications(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"restrict on job_applications fk", &pgconn.PgError{Code: "23001", ConstraintName: jobApplicationsJobDescriptionFK}, true},
		{"fk violation on job_applications fk", &pgconn.PgError{Code: "23503", ConstraintName: jobApplicationsJobDescriptionFK}, true},
		{"wrapped", fmt.Errorf("delete: %w", &pgconn.PgError{Code: "23001", ConstraintName: jobApplicationsJobDescriptionFK}), true},
		{"unrelated constraint", &pgconn.PgError{Code: "23503", ConstraintName: "job_descriptions_user_id_fkey"}, false},
		{"unrelated code", &pgconn.PgError{Code: "23505", ConstraintName: jobApplicationsJobDescriptionFK}, false},
		{"non-database error", fmt.Errorf("boom"), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := hasJobApplications(test.err); got != test.want {
				t.Fatalf("hasJobApplications() = %v, want %v", got, test.want)
			}
		})
	}
}
