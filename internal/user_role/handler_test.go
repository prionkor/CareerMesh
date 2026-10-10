package user_role

import (
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsMissingReference(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"user fk", &pgconn.PgError{Code: "23503", ConstraintName: userFK}, true},
		{"role fk", &pgconn.PgError{Code: "23503", ConstraintName: roleFK}, true},
		{"wrapped", fmt.Errorf("replace: %w", &pgconn.PgError{Code: "23503", ConstraintName: userFK}), true},
		{"unrelated constraint", &pgconn.PgError{Code: "23503", ConstraintName: "other_fkey"}, false},
		{"other code", &pgconn.PgError{Code: "23505", ConstraintName: userFK}, false},
		{"non-database error", fmt.Errorf("boom"), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isMissingReference(test.err); got != test.want {
				t.Fatalf("isMissingReference() = %v, want %v", got, test.want)
			}
		})
	}
}
