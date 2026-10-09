package app_test

import (
	"context"
	"net/http"
	"testing"

	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prionkor/careermesh/internal/auth"
	"github.com/prionkor/careermesh/internal/authorization"
)

func TestUserRoles_PutReplacesAssignment(t *testing.T) {
	app := testApp(t)
	pool := testDB(t)

	adminID := createRoleTestUser(t, pool)
	targetID := createRoleTestUser(t, pool)
	token := roleTestToken(t, adminID, authorization.RoleSuperadmin)

	for _, roleName := range []string{"admin", "admin", "user"} {
		roleID := lookupRoleID(t, pool, roleName)
		resp := doRequest(t, app, http.MethodPut, "/api/v1/user_roles", token, map[string]string{
			"userId": targetID,
			"roleId": roleID,
		})
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("assign %s role: expected 204, got %d: %s", roleName, resp.StatusCode, readBody(t, resp))
		}
		resp.Body.Close()

		var count int
		var assignedRoleID string
		if err := pool.QueryRow(context.Background(), `
			SELECT COUNT(*), COALESCE(MAX(role_id::text), '')
			FROM user_roles WHERE user_id = $1
		`, targetID).Scan(&count, &assignedRoleID); err != nil {
			t.Fatalf("read role assignment: %v", err)
		}
		if count != 1 || assignedRoleID != roleID {
			t.Fatalf("assignment = (%d rows, %s), want exactly role %s", count, assignedRoleID, roleID)
		}
	}
}

func TestUserRoles_PutRequiresSuperadmin(t *testing.T) {
	app := testApp(t)
	pool := testDB(t)

	regularID := createRoleTestUser(t, pool)
	regularToken := roleTestToken(t, regularID, authorization.RoleUser)
	adminID := createRoleTestUser(t, pool)
	grantUserRole(t, pool, adminID, "admin")
	adminToken := roleTestToken(t, adminID, authorization.RoleAdmin)
	roleID := lookupRoleID(t, pool, "admin")

	for _, test := range []struct {
		name  string
		token string
		user  string
	}{
		{name: "regular user", token: regularToken, user: regularID},
		{name: "admin", token: adminToken, user: regularID},
	} {
		t.Run(test.name, func(t *testing.T) {
			resp := doRequest(t, app, http.MethodPut, "/api/v1/user_roles", test.token, map[string]string{
				"userId": test.user,
				"roleId": roleID,
			})
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusForbidden {
				t.Fatalf("expected 403, got %d", resp.StatusCode)
			}
		})
	}
}

func TestUserRoles_PutValidatesIDs(t *testing.T) {
	app := testApp(t)
	pool := testDB(t)

	actorID := createRoleTestUser(t, pool)
	token := roleTestToken(t, actorID, authorization.RoleSuperadmin)
	targetID := createRoleTestUser(t, pool)
	roleID := lookupRoleID(t, pool, "admin")

	tests := []struct {
		name string
		body map[string]string
		want int
	}{
		{
			name: "missing user",
			body: map[string]string{"userId": uuid.NewV7().String(), "roleId": roleID},
			want: http.StatusNotFound,
		},
		{
			name: "missing role",
			body: map[string]string{"userId": targetID, "roleId": uuid.NewV7().String()},
			want: http.StatusNotFound,
		},
		{
			name: "invalid user id",
			body: map[string]string{"userId": "not-a-uuid", "roleId": roleID},
			want: http.StatusBadRequest,
		},
		{
			name: "missing role id",
			body: map[string]string{"userId": targetID},
			want: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := doRequest(t, app, http.MethodPut, "/api/v1/user_roles", token, test.body)
			defer resp.Body.Close()
			if resp.StatusCode != test.want {
				t.Fatalf("expected %d, got %d", test.want, resp.StatusCode)
			}
		})
	}
}

func grantUserRole(t *testing.T, pool *pgxpool.Pool, userID, roleName string) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, id FROM roles WHERE name = $2
	`, userID, roleName)
	if err != nil {
		t.Fatalf("grant %s role: %v", roleName, err)
	}
}

func lookupRoleID(t *testing.T, pool *pgxpool.Pool, roleName string) string {
	t.Helper()
	var roleID string
	if err := pool.QueryRow(context.Background(), `SELECT id::text FROM roles WHERE name = $1`, roleName).Scan(&roleID); err != nil {
		t.Fatalf("lookup %s role: %v", roleName, err)
	}
	return roleID
}

func createRoleTestUser(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()

	email := uuid.NewV7().String() + "@example.com"
	var userID string
	if err := pool.QueryRow(context.Background(), `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id::text
	`, email, "test-password-hash").Scan(&userID); err != nil {
		t.Fatalf("create test user: %v", err)
	}
	grantUserRole(t, pool, userID, "user")
	return userID
}

func roleTestToken(t *testing.T, userID string, role authorization.Role) string {
	t.Helper()
	id, err := uuid.Parse(userID)
	if err != nil {
		t.Fatalf("parse test user id: %v", err)
	}
	token, err := auth.GenerateAccessToken([]byte(testJWTSecret), id, authorization.PermissionsForRoles([]authorization.Role{role}))
	if err != nil {
		t.Fatalf("generate %s token: %v", role, err)
	}
	return token
}
