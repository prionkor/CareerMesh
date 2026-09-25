package app

import (
	"net/http"
	"testing"
	"time"

	"uuid"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prionkor/careermesh/internal/auth"
)

func TestPermission_InsufficientPermissionsRejected(t *testing.T) {
	app := testApp(t)

	// A validly signed token with no permissions at all.
	token, err := auth.GenerateAccessToken([]byte(testJWTSecret), uuid.NewV7(), []string{})
	if err != nil {
		t.Fatalf("generate access token: %v", err)
	}

	resp := doRequest(t, app, http.MethodGet, "/api/v1/experiences/00000000-0000-0000-0000-000000000000", token, nil)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 for insufficient permissions, got %d", resp.StatusCode)
	}
}

func TestAuth_ExpiredTokenRejected(t *testing.T) {
	app := testApp(t)

	claims := auth.Claims{
		Permissions: []string{},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   uuid.NewV7().String(),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}

	expired, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testJWTSecret))
	if err != nil {
		t.Fatalf("sign expired token: %v", err)
	}

	resp := doRequest(t, app, http.MethodGet, "/api/v1/experiences/00000000-0000-0000-0000-000000000000", expired, nil)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d", resp.StatusCode)
	}
}
