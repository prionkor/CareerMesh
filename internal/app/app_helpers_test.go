package app

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
	"github.com/joho/godotenv"
	"github.com/prionkor/careermesh/internal/database"
)

// testJWTSecret is the fixed signing secret used across integration tests.
const testJWTSecret = "test-secret-do-not-use-in-production"

// testApp builds a fully wired app against the real, migrated PostgreSQL
// instance used for local development (see docker-compose.yaml). Tests are
// skipped if the database is not reachable.
func testApp(t *testing.T) *fiber.App {
	t.Helper()

	_ = godotenv.Load(mustFindEnvFile(t))

	db, err := database.NewPool(context.Background(), database.Config{
		Host:     envOr("DB_HOST", "localhost"),
		Port:     envOr("DB_PORT", "5432"),
		User:     envOr("DB_USER", "careermesh"),
		Password: envOr("DB_PASS", "careermesh"),
		Name:     envOr("DB_NAME", "careermesh"),
		SSLMode:  envOr("DB_SSLMODE", "disable"),
	})
	if err != nil {
		t.Skipf("database not available, skipping integration test: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return New(db, []byte(testJWTSecret))
}

// testDB opens a direct pool for test setup/assertions that bypass the HTTP API.
func testDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	pool, err := database.NewPool(context.Background(), database.Config{
		Host:     envOr("DB_HOST", "localhost"),
		Port:     envOr("DB_PORT", "5432"),
		User:     envOr("DB_USER", "careermesh"),
		Password: envOr("DB_PASS", "careermesh"),
		Name:     envOr("DB_NAME", "careermesh"),
		SSLMode:  envOr("DB_SSLMODE", "disable"),
	})
	if err != nil {
		t.Skipf("database not available, skipping integration test: %v", err)
	}
	t.Cleanup(func() { pool.Close() })

	return pool
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

// mustFindEnvFile lets godotenv.Load() find the repo-root .env regardless of
// the package's working directory during `go test ./...`.
func mustFindEnvFile(t *testing.T) string {
	t.Helper()
	for _, candidate := range []string{".env", "../../.env"} {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ".env"
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
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
