package app

import (
	"net/http"
	"testing"
)

func TestLogin_Success(t *testing.T) {
	app := testApp(t)

	_, token, _, _ := registerAndLogin(t, app)
	if token == "" {
		t.Fatal("expected a non-empty access token")
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	app := testApp(t)

	_, _, email, _ := registerAndLogin(t, app)

	resp := doRequest(t, app, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email":    email,
		"password": "definitely-the-wrong-password",
	})
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	app := testApp(t)

	resp := doRequest(t, app, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email":    "no-such-user@example.com",
		"password": "whatever",
	})
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestLogin_ResponseNeverContainsPasswordHash(t *testing.T) {
	app := testApp(t)

	_, token, _, _ := registerAndLogin(t, app)

	// Also check the registration response and the authenticated GetByID response.
	userID, _, _, _ := registerAndLogin(t, app)

	resp := doRequest(t, app, http.MethodGet, "/api/v1/users/"+userID, token, nil)
	body := readBody(t, resp)
	if containsPasswordHash(body) {
		t.Fatalf("response leaked password_hash: %s", body)
	}
}

func TestProtectedRoute_MissingToken(t *testing.T) {
	app := testApp(t)

	_, _, _, _ = registerAndLogin(t, app)

	resp := doRequest(t, app, http.MethodGet, "/api/v1/experiences/00000000-0000-0000-0000-000000000000", "", nil)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing token, got %d", resp.StatusCode)
	}
}

func TestProtectedRoute_MalformedToken(t *testing.T) {
	app := testApp(t)

	resp := doRequest(t, app, http.MethodGet, "/api/v1/experiences/00000000-0000-0000-0000-000000000000", "not-a-real-jwt", nil)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for malformed token, got %d", resp.StatusCode)
	}
}

func TestProtectedRoute_ValidTokenGrantsAccess(t *testing.T) {
	app := testApp(t)

	userID, token, _, _ := registerAndLogin(t, app)

	resp := doRequest(t, app, http.MethodGet, "/api/v1/users/"+userID, token, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
