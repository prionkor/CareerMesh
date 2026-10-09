package app_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prionkor/careermesh/internal/auth"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/utility"
)

const authTestPassword = "correct horse battery staple"

type authIdentity struct {
	email    string
	password string
	token    string
	userID   string
}

func TestAuthLifecycle_Registration(t *testing.T) {
	application := testApp(t)
	identity := newAuthIdentity(t)

	invalid := doRequest(t, application, http.MethodPost, "/api/v1/auth/register", "", map[string]string{
		"email": "not-an-email", "password": authTestPassword,
	})
	if invalid.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid registration status = %d, want 400: %s", invalid.StatusCode, readBody(t, invalid))
	}
	invalid.Body.Close()

	body := registerAuthIdentity(t, application, &identity)
	assertSuccessResponse(t, body)
	assertNoAuthSecrets(t, body, identity.token)

	var passwordHash, tokenHash string
	if err := testDB(t).QueryRow(context.Background(), `
		SELECT password_hash, token_hash FROM pending_users WHERE email = $1
	`, identity.email).Scan(&passwordHash, &tokenHash); err != nil {
		t.Fatalf("load pending registration: %v", err)
	}
	if !strings.HasPrefix(passwordHash, "$argon2id$") || tokenHash != utility.HashToken(identity.token) {
		t.Fatal("pending registration did not store the expected password and verification token hashes")
	}
	if strings.Contains(body, passwordHash) || strings.Contains(body, tokenHash) || strings.Contains(body, authTestPassword) {
		t.Fatalf("registration response leaked stored credential material: %s", body)
	}

	duplicate := doRequest(t, application, http.MethodPost, "/api/v1/auth/register", "", map[string]string{
		"email": identity.email, "password": authTestPassword,
	})
	if duplicate.StatusCode != http.StatusConflict {
		t.Fatalf("duplicate registration status = %d, want 409: %s", duplicate.StatusCode, readBody(t, duplicate))
	}
	duplicate.Body.Close()
	assertAuthIdentityCounts(t, identity.email, 0, 1)
}

func TestAuthLifecycle_Verification(t *testing.T) {
	application := testApp(t)
	identity := newAuthIdentity(t)
	registerAuthIdentity(t, application, &identity)

	invalid := doRequest(t, application, http.MethodGet, "/api/v1/auth/register/verify?token=invalid-token", "", nil)
	if invalid.StatusCode != http.StatusUnauthorized {
		t.Fatalf("invalid verification token status = %d, want 401: %s", invalid.StatusCode, readBody(t, invalid))
	}
	invalid.Body.Close()

	resp := doRequest(t, application, http.MethodGet, verificationPath(identity.token), "", nil)
	body := readBody(t, resp)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("valid verification status = %d, want 200: %s", resp.StatusCode, body)
	}
	assertSuccessResponse(t, body)
	assertNoAuthSecrets(t, body, identity.token)
	assertAuthIdentityCounts(t, identity.email, 1, 0)

	reused := doRequest(t, application, http.MethodGet, verificationPath(identity.token), "", nil)
	if reused.StatusCode != http.StatusUnauthorized {
		t.Fatalf("consumed verification token status = %d, want 401: %s", reused.StatusCode, readBody(t, reused))
	}
	reused.Body.Close()

	duplicate := doRequest(t, application, http.MethodPost, "/api/v1/auth/register", "", map[string]string{
		"email": identity.email, "password": authTestPassword,
	})
	if duplicate.StatusCode != http.StatusConflict {
		t.Fatalf("verified email registration status = %d, want 409: %s", duplicate.StatusCode, readBody(t, duplicate))
	}
	duplicate.Body.Close()

	missing := doRequest(t, application, http.MethodGet, "/api/v1/auth/register/verify", "", nil)
	if missing.StatusCode != http.StatusBadRequest {
		t.Fatalf("missing verification token status = %d, want 400: %s", missing.StatusCode, readBody(t, missing))
	}
	missing.Body.Close()
}

func TestAuthLifecycle_ExpiredVerificationToken(t *testing.T) {
	application := testApp(t)
	identity := seedPendingAuthIdentity(t, time.Now().Add(-time.Hour))
	resp := doRequest(t, application, http.MethodGet, verificationPath(identity.token), "", nil)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expired verification token status = %d, want 401: %s", resp.StatusCode, readBody(t, resp))
	}
	resp.Body.Close()
	assertAuthIdentityCounts(t, identity.email, 0, 1)
}

func TestAuthLifecycle_Login(t *testing.T) {
	application := testApp(t)
	identity := newAuthIdentity(t)
	registerAuthIdentity(t, application, &identity)

	pendingLogin := doRequest(t, application, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email": identity.email, "password": authTestPassword,
	})
	if pendingLogin.StatusCode != http.StatusUnauthorized {
		t.Fatalf("pending user login status = %d, want 401: %s", pendingLogin.StatusCode, readBody(t, pendingLogin))
	}
	pendingLogin.Body.Close()

	verifyResp := doRequest(t, application, http.MethodGet, verificationPath(identity.token), "", nil)
	if verifyResp.StatusCode != http.StatusOK {
		t.Fatalf("verify user before login: expected 200, got %d: %s", verifyResp.StatusCode, readBody(t, verifyResp))
	}
	verifyResp.Body.Close()
	if err := testDB(t).QueryRow(context.Background(), `SELECT id::text FROM users WHERE email = $1`, identity.email).Scan(&identity.userID); err != nil {
		t.Fatalf("load verified user id: %v", err)
	}

	loginResp := doRequest(t, application, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email": identity.email, "password": authTestPassword,
	})
	loginBody := readBody(t, loginResp)
	if loginResp.StatusCode != http.StatusOK {
		t.Fatalf("valid login status = %d, want 200: %s", loginResp.StatusCode, loginBody)
	}
	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		User        struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}
	if err := json.Unmarshal([]byte(loginBody), &result); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if result.AccessToken == "" || result.TokenType != "Bearer" || result.ExpiresIn != int(auth.AccessTokenTTL.Seconds()) || result.User.ID != identity.userID || result.User.Email != identity.email {
		t.Fatalf("login response missing expected fields: %s", loginBody)
	}
	assertNoAuthSecrets(t, loginBody, identity.token)
	var passwordHash string
	if err := testDB(t).QueryRow(context.Background(), `SELECT password_hash FROM users WHERE email = $1`, identity.email).Scan(&passwordHash); err != nil {
		t.Fatalf("load password hash: %v", err)
	}
	if strings.Contains(loginBody, passwordHash) {
		t.Fatalf("login response leaked password hash: %s", loginBody)
	}

	wrongPassword := doRequest(t, application, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email": identity.email, "password": "incorrect-password",
	})
	wrongBody := readBody(t, wrongPassword)
	unknownEmail := doRequest(t, application, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email": newAuthEmail(), "password": authTestPassword,
	})
	unknownBody := readBody(t, unknownEmail)
	if wrongPassword.StatusCode != http.StatusUnauthorized || unknownEmail.StatusCode != http.StatusUnauthorized || wrongBody != unknownBody {
		t.Fatalf("login failures were distinguishable: wrong=(%d, %s), unknown=(%d, %s)", wrongPassword.StatusCode, wrongBody, unknownEmail.StatusCode, unknownBody)
	}

	invalid := doRequest(t, application, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email": "invalid-email", "password": authTestPassword,
	})
	if invalid.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid login input status = %d, want 400: %s", invalid.StatusCode, readBody(t, invalid))
	}
	invalid.Body.Close()
}

func TestAuthLifecycle_AccessTokenValidation(t *testing.T) {
	application := testApp(t)
	identity := createVerifiedAuthIdentity(t, application)
	token := loginAuthIdentity(t, application, identity)
	path := "/api/v1/users/" + identity.userID

	valid := doRequest(t, application, http.MethodGet, path, token, nil)
	if valid.StatusCode != http.StatusOK {
		t.Fatalf("valid access token status = %d, want 200: %s", valid.StatusCode, readBody(t, valid))
	}
	valid.Body.Close()

	claims := auth.Claims{
		Permissions: authorization.PermissionsForRoles([]authorization.Role{authorization.RoleUser}),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   identity.userID,
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
		},
	}
	expired := signTestAccessToken(t, jwt.SigningMethodHS256, []byte(testJWTSecret), claims)
	wrongSignatureClaims := auth.Claims{
		Permissions: claims.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: identity.userID, ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	wrongSignature := signTestAccessToken(t, jwt.SigningMethodHS256, []byte("different-test-secret"), wrongSignatureClaims)
	wrongAlgorithm := signTestAccessToken(t, jwt.SigningMethodHS384, []byte(testJWTSecret), wrongSignatureClaims)

	tests := []struct {
		name  string
		token string
	}{
		{name: "missing", token: ""},
		{name: "malformed", token: "not-a-jwt"},
		{name: "tampered", token: tamperAccessToken(token)},
		{name: "expired", token: expired},
		{name: "wrong signature", token: wrongSignature},
		{name: "wrong algorithm", token: wrongAlgorithm},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := doRequest(t, application, http.MethodGet, path, test.token, nil)
			if resp.StatusCode != http.StatusUnauthorized {
				t.Fatalf("protected request status = %d, want 401: %s", resp.StatusCode, readBody(t, resp))
			}
			resp.Body.Close()
		})
	}
}

func TestAuthLifecycle_PermissionIntegration(t *testing.T) {
	application := testApp(t)
	roles := []struct {
		name             string
		role             authorization.Role
		wantRoleUpdate   bool
		wantUpdateStatus int
	}{
		{name: "user", role: authorization.RoleUser, wantUpdateStatus: http.StatusForbidden},
		{name: "admin", role: authorization.RoleAdmin, wantUpdateStatus: http.StatusForbidden},
		{name: "superadmin", role: authorization.RoleSuperadmin, wantRoleUpdate: true, wantUpdateStatus: http.StatusNoContent},
	}
	adminRoleID := lookupRoleID(t, testDB(t), "admin")
	for _, test := range roles {
		t.Run(test.name, func(t *testing.T) {
			identity := createVerifiedAuthIdentity(t, application)
			setAuthTestRole(t, testDB(t), identity.userID, test.name)
			token := loginAuthIdentity(t, application, identity)
			claims, err := auth.ParseAccessToken([]byte(testJWTSecret), token)
			if err != nil {
				t.Fatalf("parse login token: %v", err)
			}
			gotPermission := authorization.HasPermission(claims.Permissions, authorization.PermRoleUpdateAll)
			if gotPermission != test.wantRoleUpdate {
				t.Fatalf("role:update:all = %t, want %t", gotPermission, test.wantRoleUpdate)
			}
			if test.role == authorization.RoleUser && !authorization.HasPermission(claims.Permissions, authorization.PermProfilesUpdateOwn) {
				t.Fatal("regular user did not receive profiles:update:own")
			}
			if test.role == authorization.RoleAdmin && !authorization.HasPermission(claims.Permissions, authorization.PermProfilesUpdateAll) {
				t.Fatal("admin did not receive profiles:update:all")
			}

			resp := doRequest(t, application, http.MethodPut, "/api/v1/user_roles", token, map[string]string{
				"userId": identity.userID, "roleId": adminRoleID,
			})
			if resp.StatusCode != test.wantUpdateStatus {
				t.Fatalf("role update status = %d, want %d: %s", resp.StatusCode, test.wantUpdateStatus, readBody(t, resp))
			}
			resp.Body.Close()
		})
	}
}

func newAuthIdentity(t *testing.T) authIdentity {
	t.Helper()
	identity := authIdentity{
		email:    "auth-" + uuid.NewV7().String() + "@example.com",
		password: authTestPassword,
	}
	pool := testDB(t)
	t.Cleanup(func() {
		if _, err := pool.Exec(context.Background(), `DELETE FROM pending_users WHERE email = $1`, identity.email); err != nil {
			t.Errorf("delete pending auth test user: %v", err)
		}
		if _, err := pool.Exec(context.Background(), `DELETE FROM users WHERE email = $1`, identity.email); err != nil {
			t.Errorf("delete auth test user: %v", err)
		}
	})
	return identity
}

func newAuthEmail() string {
	return "auth-" + uuid.NewV7().String() + "@example.com"
}

func registerAuthIdentity(t *testing.T, application *fiber.App, identity *authIdentity) string {
	t.Helper()
	resp := doRequest(t, application, http.MethodPost, "/api/v1/auth/register", "", map[string]string{
		"email": identity.email, "password": identity.password,
	})
	body := readBody(t, resp)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("registration status = %d, want 200: %s", resp.StatusCode, body)
	}
	token, err := utility.GenerateVerificationToken()
	if err != nil {
		t.Fatalf("generate verification token: %v", err)
	}
	identity.token = token
	if _, err := testDB(t).Exec(context.Background(), `
		UPDATE pending_users SET token_hash = $2 WHERE email = $1
	`, identity.email, utility.HashToken(token)); err != nil {
		t.Fatalf("store verification token hash: %v", err)
	}
	if err := testDB(t).QueryRow(context.Background(), `
		SELECT id::text FROM pending_users WHERE email = $1
	`, identity.email).Scan(&identity.userID); err != nil {
		t.Fatalf("load pending auth user id: %v", err)
	}
	return body
}

func seedPendingAuthIdentity(t *testing.T, expiresAt time.Time) authIdentity {
	t.Helper()
	identity := newAuthIdentity(t)
	token, err := utility.GenerateVerificationToken()
	if err != nil {
		t.Fatalf("generate verification token: %v", err)
	}
	passwordHash, err := utility.HashPassword(identity.password)
	if err != nil {
		t.Fatalf("hash test password: %v", err)
	}
	if _, err := testDB(t).Exec(context.Background(), `
		INSERT INTO pending_users (email, password_hash, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
	`, identity.email, passwordHash, utility.HashToken(token), expiresAt); err != nil {
		t.Fatalf("insert pending auth test user: %v", err)
	}
	identity.token = token
	return identity
}

func createVerifiedAuthIdentity(t *testing.T, application *fiber.App) authIdentity {
	t.Helper()
	identity := newAuthIdentity(t)
	registerAuthIdentity(t, application, &identity)
	resp := doRequest(t, application, http.MethodGet, verificationPath(identity.token), "", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("verify auth test user: expected 200, got %d: %s", resp.StatusCode, readBody(t, resp))
	}
	resp.Body.Close()
	if err := testDB(t).QueryRow(context.Background(), `
		SELECT id::text FROM users WHERE email = $1
	`, identity.email).Scan(&identity.userID); err != nil {
		t.Fatalf("load verified auth user id: %v", err)
	}
	return identity
}

func loginAuthIdentity(t *testing.T, application *fiber.App, identity authIdentity) string {
	t.Helper()
	resp := doRequest(t, application, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email": identity.email, "password": identity.password,
	})
	body := readBody(t, resp)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login auth test user: expected 200, got %d: %s", resp.StatusCode, body)
	}
	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		t.Fatalf("decode auth test login: %v", err)
	}
	return result.AccessToken
}

func setAuthTestRole(t *testing.T, pool *pgxpool.Pool, userID, roleName string) {
	t.Helper()
	result, err := pool.Exec(context.Background(), `
		UPDATE user_roles
		SET role_id = (SELECT id FROM roles WHERE name = $2)
		WHERE user_id = $1
	`, userID, roleName)
	if err != nil {
		t.Fatalf("assign %s role: %v", roleName, err)
	}
	if result.RowsAffected() != 1 {
		t.Fatalf("assign %s role updated %d rows, want 1", roleName, result.RowsAffected())
	}
}

func assertAuthIdentityCounts(t *testing.T, email string, users, pending int) {
	t.Helper()
	var userCount, pendingCount int
	if err := testDB(t).QueryRow(context.Background(), `SELECT COUNT(*) FROM users WHERE email = $1`, email).Scan(&userCount); err != nil {
		t.Fatalf("count auth test users: %v", err)
	}
	if err := testDB(t).QueryRow(context.Background(), `SELECT COUNT(*) FROM pending_users WHERE email = $1`, email).Scan(&pendingCount); err != nil {
		t.Fatalf("count pending auth test users: %v", err)
	}
	if userCount != users || pendingCount != pending {
		t.Fatalf("auth database state = users:%d pending:%d, want users:%d pending:%d", userCount, pendingCount, users, pending)
	}
}

func assertSuccessResponse(t *testing.T, body string) {
	t.Helper()
	var response struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		t.Fatalf("decode success response: %v", err)
	}
	if !response.Success {
		t.Fatalf("response success = false: %s", body)
	}
}

func assertNoAuthSecrets(t *testing.T, body, verificationToken string) {
	t.Helper()
	for _, secret := range []string{"password_hash", "token_hash", "argon2id", verificationToken} {
		if secret != "" && strings.Contains(body, secret) {
			t.Fatalf("auth response exposed %q: %s", secret, body)
		}
	}
}

func verificationPath(token string) string {
	return "/api/v1/auth/register/verify?token=" + url.QueryEscape(token)
}

func signTestAccessToken(t *testing.T, method jwt.SigningMethod, secret []byte, claims auth.Claims) string {
	t.Helper()
	token, err := jwt.NewWithClaims(method, claims).SignedString(secret)
	if err != nil {
		t.Fatalf("sign %s token: %v", method.Alg(), err)
	}
	return token
}

func tamperAccessToken(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "tampered-token"
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || len(signature) == 0 {
		return "tampered-token"
	}
	signature[0] ^= 0xff
	parts[2] = base64.RawURLEncoding.EncodeToString(signature)
	return strings.Join(parts, ".")
}
