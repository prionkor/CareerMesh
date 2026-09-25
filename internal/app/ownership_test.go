package app

import (
	"net/http"
	"testing"
)

func TestOwnership_UserCanAccessAndUpdateOwnResource(t *testing.T) {
	app := testApp(t)

	_, token, _, _ := registerAndLogin(t, app)

	createResp := doRequest(t, app, http.MethodPost, "/api/v1/experiences", token, map[string]any{
		"company":         "Acme",
		"title":           "Engineer",
		"employment_type": "full_time",
		"start_date":      "2020-01-01T00:00:00Z",
	})
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("create experience: expected 201, got %d", createResp.StatusCode)
	}

	var created struct {
		ID string `json:"id"`
	}
	decodeJSON(t, createResp, &created)

	getResp := doRequest(t, app, http.MethodGet, "/api/v1/experiences/"+created.ID, token, nil)
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("get own experience: expected 200, got %d", getResp.StatusCode)
	}

	updateResp := doRequest(t, app, http.MethodPut, "/api/v1/experiences/"+created.ID, token, map[string]any{
		"company":         "Acme Corp",
		"title":           "Senior Engineer",
		"employment_type": "full_time",
		"start_date":      "2020-01-01T00:00:00Z",
	})
	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("update own experience: expected 200, got %d", updateResp.StatusCode)
	}
}

func TestOwnership_UserCannotModifyAnotherUsersResource(t *testing.T) {
	app := testApp(t)

	_, tokenA, _, _ := registerAndLogin(t, app)
	_, tokenB, _, _ := registerAndLogin(t, app)

	createResp := doRequest(t, app, http.MethodPost, "/api/v1/experiences", tokenA, map[string]any{
		"company":         "Acme",
		"title":           "Engineer",
		"employment_type": "full_time",
		"start_date":      "2020-01-01T00:00:00Z",
	})
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("create experience: expected 201, got %d", createResp.StatusCode)
	}

	var created struct {
		ID string `json:"id"`
	}
	decodeJSON(t, createResp, &created)

	updateResp := doRequest(t, app, http.MethodPut, "/api/v1/experiences/"+created.ID, tokenB, map[string]any{
		"company":         "Hijacked",
		"title":           "Hijacked",
		"employment_type": "full_time",
		"start_date":      "2020-01-01T00:00:00Z",
	})
	if updateResp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 updating another user's experience, got %d", updateResp.StatusCode)
	}

	deleteResp := doRequest(t, app, http.MethodDelete, "/api/v1/experiences/"+created.ID, tokenB, nil)
	if deleteResp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 deleting another user's experience, got %d", deleteResp.StatusCode)
	}
}

func TestOwnership_AdminCanBypassOwnership(t *testing.T) {
	app := testApp(t)
	pool := testDB(t)

	_, tokenA, _, _ := registerAndLogin(t, app)
	adminID, _, adminEmail, adminPassword := registerAndLogin(t, app)

	createResp := doRequest(t, app, http.MethodPost, "/api/v1/experiences", tokenA, map[string]any{
		"company":         "Acme",
		"title":           "Engineer",
		"employment_type": "full_time",
		"start_date":      "2020-01-01T00:00:00Z",
	})
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("create experience: expected 201, got %d", createResp.StatusCode)
	}

	var created struct {
		ID string `json:"id"`
	}
	decodeJSON(t, createResp, &created)

	grantAdminRole(t, pool, adminID)
	adminToken := login(t, app, adminEmail, adminPassword) // refresh token to pick up the new permission

	updateResp := doRequest(t, app, http.MethodPut, "/api/v1/experiences/"+created.ID, adminToken, map[string]any{
		"company":         "Updated By Admin",
		"title":           "Senior Engineer",
		"employment_type": "full_time",
		"start_date":      "2020-01-01T00:00:00Z",
	})
	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("expected admin to bypass ownership with 200, got %d", updateResp.StatusCode)
	}
}
