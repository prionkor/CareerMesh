package app_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prionkor/careermesh/internal/app"
	"github.com/prionkor/careermesh/tests/integration/testdb"
)

// testJWTSecret is the fixed signing secret used across integration tests.
const testJWTSecret = "test-secret-do-not-use-in-production"

// sharedPool is the single pool for the whole suite, created in TestMain
// against TEST_DATABASE_URL with migrations already applied.
var sharedPool *pgxpool.Pool

func TestMain(m *testing.M) {
	pool, err := testdb.Setup(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "integration test setup failed: %v\n", err)
		os.Exit(1)
	}
	sharedPool = pool

	code := m.Run()
	pool.Close()
	os.Exit(code)
}

// testApp builds a fully wired app on the shared test database pool.
func testApp(t *testing.T) *fiber.App {
	t.Helper()
	return app.New(sharedPool, []byte(testJWTSecret))
}

// testDB returns the shared pool for setup/assertions that bypass the HTTP API.
func testDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	return sharedPool
}

// grantAdminRole grants the 'admin' role directly via the database, bypassing the HTTP API.
func grantAdminRole(t *testing.T, pool *pgxpool.Pool, userID string) {
	t.Helper()

	_, err := pool.Exec(context.Background(), `
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, id FROM roles WHERE name = 'admin'
	`, userID)
	if err != nil {
		t.Fatalf("grant admin role: %v", err)
	}
}

func doRequest(t *testing.T, app *fiber.App, method, path, token string, body any) *http.Response {
	t.Helper()

	var reader *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(b)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := app.Test(req, fiber.TestConfig{Timeout: 10 * time.Second})
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}

	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, out any) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
}

// readBody reads and returns the raw response body as a string.
func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}
	return string(b)
}

// containsPasswordHash reports whether a raw JSON response leaks the password hash.
func containsPasswordHash(body string) bool {
	return strings.Contains(body, "password_hash") || strings.Contains(body, "argon2id")
}

// registerAndLogin creates a brand-new user through the public registration
// endpoint and logs in, returning the user id, access token and raw password.
func registerAndLogin(t *testing.T, app *fiber.App) (userID, token, email, password string) {
	t.Helper()

	email = fmt.Sprintf("test-%d-%d@example.com", time.Now().UnixNano(), os.Getpid())
	password = "correct horse battery staple"

	resp := doRequest(t, app, http.MethodPost, "/api/v1/users", "", map[string]string{
		"email":    email,
		"password": password,
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("register: expected 201, got %d", resp.StatusCode)
	}

	var created struct {
		ID string `json:"id"`
	}
	decodeJSON(t, resp, &created)

	loginResp := doRequest(t, app, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email":    email,
		"password": password,
	})
	if loginResp.StatusCode != http.StatusOK {
		t.Fatalf("login: expected 200, got %d", loginResp.StatusCode)
	}

	var loginBody struct {
		AccessToken string `json:"access_token"`
	}
	decodeJSON(t, loginResp, &loginBody)

	return created.ID, loginBody.AccessToken, email, password
}

// login exchanges an email/password for a fresh access token, reflecting
// any RBAC changes made since the last login (tokens are a permission snapshot).
func login(t *testing.T, app *fiber.App, email, password string) string {
	t.Helper()

	resp := doRequest(t, app, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email":    email,
		"password": password,
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login: expected 200, got %d", resp.StatusCode)
	}

	var body struct {
		AccessToken string `json:"access_token"`
	}
	decodeJSON(t, resp, &body)

	return body.AccessToken
}
