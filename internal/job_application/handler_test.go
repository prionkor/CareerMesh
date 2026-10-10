package job_application

import (
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsMissingJobDescription(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"job_description_id fk", &pgconn.PgError{Code: "23503", ConstraintName: jobDescriptionFK}, true},
		{"wrapped", fmt.Errorf("create: %w", &pgconn.PgError{Code: "23503", ConstraintName: jobDescriptionFK}), true},
		{"user_id fk", &pgconn.PgError{Code: "23503", ConstraintName: "job_applications_user_id_fkey"}, false},
		{"other code", &pgconn.PgError{Code: "23505", ConstraintName: jobDescriptionFK}, false},
		{"non-database error", fmt.Errorf("boom"), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isMissingJobDescription(test.err); got != test.want {
				t.Fatalf("isMissingJobDescription() = %v, want %v", got, test.want)
			}
		})
	}
}
