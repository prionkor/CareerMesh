package app_test

import (
	"net/http"
	"testing"

	"uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/auth"
	"github.com/prionkor/careermesh/models"
)

func TestJobApplicationCRUDAndPagination(t *testing.T) {
	app := testApp(t)
	_, token := newJobDescriptionUser(t, app)
	jd1 := createJobDescription(t, app, token, jobDescriptionBody("First posting"))
	jd2 := createJobDescription(t, app, token, jobDescriptionBody("Second posting"))

	first := createJobApplication(t, app, token, map[string]any{"job_description_id": jd1.ID.String()})
	if first.Status != models.JobApplicationStatusDraft || first.AppliedAt != nil || first.Notes != "" {
		t.Fatalf("unexpected defaults: %+v", first)
	}

	second := createJobApplication(t, app, token, map[string]any{
		"job_description_id": jd1.ID.String(), "status": "applied", "applied_at": "2026-03-04", "notes": "via referral",
	})
	if second.Status != models.JobApplicationStatusApplied || second.AppliedAt == nil ||
		second.AppliedAt.Format("2006-01-02") != "2026-03-04" || second.Notes != "via referral" {
		t.Fatalf("unexpected application: %+v", second)
	}
	third := createJobApplication(t, app, token, map[string]any{"job_description_id": jd2.ID.String()})

	path := "/api/v1/job-applications/" + second.ID.String()
	updateResp := doRequest(t, app, http.MethodPut, path, token, map[string]any{
		"status": "interviewing", "applied_at": "2026-03-05", "notes": "phone screen done",
		// Ignored: ownership and reference are not editable.
		"job_description_id": jd2.ID.String(), "user_id": uuid.NewV7().String(),
	})
	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("update: expected 200, got %d", updateResp.StatusCode)
	}
	var updated models.JobApplication
	decodeJSON(t, updateResp, &updated)
	if updated.Status != models.JobApplicationStatusInterviewing || updated.Notes != "phone screen done" ||
		updated.AppliedAt == nil || updated.AppliedAt.Format("2006-01-02") != "2026-03-05" ||
		updated.JobDescriptionID != jd1.ID || updated.UserID != second.UserID {
		t.Fatalf("unexpected updated application: %+v", updated)
	}

	getResp := doRequest(t, app, http.MethodGet, path, token, nil)
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", getResp.StatusCode)
	}
	var fetched models.JobApplication
	decodeJSON(t, getResp, &fetched)
	if fetched.ID != second.ID || fetched.Status != models.JobApplicationStatusInterviewing {
		t.Fatalf("unexpected fetched application: %+v", fetched)
	}

	type listBody struct {
		Items  []models.JobApplication `json:"items"`
		Limit  int                     `json:"limit"`
		Offset int                     `json:"offset"`
	}
	list := func(query string) listBody {
		resp := doRequest(t, app, http.MethodGet, "/api/v1/job-applications"+query, token, nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("list %q: expected 200, got %d", query, resp.StatusCode)
		}
		var body listBody
		decodeJSON(t, resp, &body)
		return body
	}

	all := list("")
	if len(all.Items) != 3 || all.Limit != 20 || all.Offset != 0 {
		t.Fatalf("unexpected list: %+v", all)
	}
	if all.Items[0].ID != third.ID || all.Items[2].ID != first.ID {
		t.Fatalf("list not ordered newest first: %+v", all.Items)
	}

	page := list("?limit=1&offset=1")
	if len(page.Items) != 1 || page.Items[0].ID != second.ID || page.Limit != 1 || page.Offset != 1 {
		t.Fatalf("unexpected page: %+v", page)
	}

	filtered := list("?job_description_id=" + jd1.ID.String())
	if len(filtered.Items) != 2 {
		t.Fatalf("expected 2 applications for first job description, got %d", len(filtered.Items))
	}
	for _, item := range filtered.Items {
		if item.JobDescriptionID != jd1.ID {
			t.Fatalf("filter returned application for another job description: %+v", item)
		}
	}

	deleteResp := doRequest(t, app, http.MethodDelete, path, token, nil)
	if deleteResp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete: expected 204, got %d", deleteResp.StatusCode)
	}
	deleteResp.Body.Close()
	missingResp := doRequest(t, app, http.MethodGet, path, token, nil)
	if missingResp.StatusCode != http.StatusNotFound {
		t.Fatalf("get deleted: expected 404, got %d", missingResp.StatusCode)
	}
	missingResp.Body.Close()
}

func TestJobApplicationValidation(t *testing.T) {
	app := testApp(t)
	_, token := newJobDescriptionUser(t, app)
	jd := createJobDescription(t, app, token, jobDescriptionBody("Posting"))
	created := createJobApplication(t, app, token, map[string]any{"job_description_id": jd.ID.String()})
	jdID := jd.ID.String()

	createCases := map[string]map[string]any{
		"missing job description id": {},
		"invalid job description id": {"job_description_id": "not-a-uuid"},
		"invalid status":             {"job_description_id": jdID, "status": "ghosted"},
		"invalid date":               {"job_description_id": jdID, "applied_at": "03/04/2026"},
		"non-existent date":          {"job_description_id": jdID, "applied_at": "2026-02-30"},
	}
	for name, body := range createCases {
		resp := doRequest(t, app, http.MethodPost, "/api/v1/job-applications", token, body)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("create %s: expected 400, got %d", name, resp.StatusCode)
		}
		if got := readBody(t, resp); got != `{"error":"invalid request body"}` {
			t.Fatalf("create %s: unexpected body %s", name, got)
		}
	}

	updatePath := "/api/v1/job-applications/" + created.ID.String()
	updateCases := map[string]map[string]any{
		"missing status": {"notes": "x"},
		"invalid status": {"status": "ghosted"},
		"invalid date":   {"status": "applied", "applied_at": "yesterday"},
	}
	for name, body := range updateCases {
		resp := doRequest(t, app, http.MethodPut, updatePath, token, body)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("update %s: expected 400, got %d", name, resp.StatusCode)
		}
		resp.Body.Close()
	}

	for _, request := range []struct {
		name, method, path string
		body               any
	}{
		{"get malformed id", http.MethodGet, "/api/v1/job-applications/not-a-uuid", nil},
		{"update malformed id", http.MethodPut, "/api/v1/job-applications/not-a-uuid", map[string]any{"status": "applied"}},
		{"delete malformed id", http.MethodDelete, "/api/v1/job-applications/not-a-uuid", nil},
		{"limit zero", http.MethodGet, "/api/v1/job-applications?limit=0", nil},
		{"negative offset", http.MethodGet, "/api/v1/job-applications?offset=-1", nil},
		{"malformed filter", http.MethodGet, "/api/v1/job-applications?job_description_id=nope", nil},
	} {
		resp := doRequest(t, app, request.method, request.path, token, request.body)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("%s: expected 400, got %d", request.name, resp.StatusCode)
		}
		resp.Body.Close()
	}

	unknown := doRequest(t, app, http.MethodGet, "/api/v1/job-applications/"+uuid.NewV7().String(), token, nil)
	if unknown.StatusCode != http.StatusNotFound {
		t.Fatalf("unknown application: expected 404, got %d", unknown.StatusCode)
	}
	unknown.Body.Close()

	missingJD := doRequest(t, app, http.MethodPost, "/api/v1/job-applications", token,
		map[string]any{"job_description_id": uuid.NewV7().String()})
	if missingJD.StatusCode != http.StatusNotFound {
		t.Fatalf("missing job description: expected 404, got %d", missingJD.StatusCode)
	}
	missingJD.Body.Close()
}

func TestJobApplicationOwnership(t *testing.T) {
	app := testApp(t)
	_, ownerToken := newJobDescriptionUser(t, app)
	_, otherToken := newJobDescriptionUser(t, app)
	jd := createJobDescription(t, app, ownerToken, jobDescriptionBody("Private posting"))
	created := createJobApplication(t, app, ownerToken, map[string]any{"job_description_id": jd.ID.String()})
	path := "/api/v1/job-applications/" + created.ID.String()

	crossCreate := doRequest(t, app, http.MethodPost, "/api/v1/job-applications", otherToken,
		map[string]any{"job_description_id": jd.ID.String()})
	if crossCreate.StatusCode != http.StatusNotFound {
		t.Fatalf("create for another user's job description: expected 404, got %d", crossCreate.StatusCode)
	}
	crossCreate.Body.Close()

	for _, request := range []struct {
		method string
		body   any
	}{
		{http.MethodGet, nil},
		{http.MethodPut, map[string]any{"status": "applied"}},
		{http.MethodDelete, nil},
	} {
		resp := doRequest(t, app, request.method, path, otherToken, request.body)
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("other user %s: expected 403, got %d", request.method, resp.StatusCode)
		}
		resp.Body.Close()
	}

	listResp := doRequest(t, app, http.MethodGet, "/api/v1/job-applications", otherToken, nil)
	var list struct {
		Items []models.JobApplication `json:"items"`
	}
	decodeJSON(t, listResp, &list)
	if len(list.Items) != 0 {
		t.Fatalf("other user list returned another user's applications: %+v", list.Items)
	}

	unauthenticated := doRequest(t, app, http.MethodGet, "/api/v1/job-applications", "", nil)
	if unauthenticated.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthenticated: expected 401, got %d", unauthenticated.StatusCode)
	}
	unauthenticated.Body.Close()

	noPermissionsToken, err := auth.GenerateAccessToken([]byte(testJWTSecret), uuid.NewV7(), []string{})
	if err != nil {
		t.Fatalf("generate token without permissions: %v", err)
	}
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		resp := doRequest(t, app, method, "/api/v1/job-applications", noPermissionsToken,
			map[string]any{"job_description_id": jd.ID.String()})
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("no permissions %s: expected 403, got %d", method, resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func TestJobApplicationAdminAccess(t *testing.T) {
	app := testApp(t)
	_, ownerToken := newJobDescriptionUser(t, app)
	adminIdentity := createVerifiedAuthIdentity(t, app)
	jd := createJobDescription(t, app, ownerToken, jobDescriptionBody("Private posting"))
	created := createJobApplication(t, app, ownerToken, map[string]any{"job_description_id": jd.ID.String()})

	setAuthTestRole(t, testDB(t), adminIdentity.userID, "admin")
	adminToken := loginAuthIdentity(t, app, adminIdentity)
	path := "/api/v1/job-applications/" + created.ID.String()

	getResp := doRequest(t, app, http.MethodGet, path, adminToken, nil)
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("admin get: expected 200, got %d", getResp.StatusCode)
	}
	getResp.Body.Close()

	listResp := doRequest(t, app, http.MethodGet, "/api/v1/job-applications?job_description_id="+jd.ID.String(), adminToken, nil)
	var list struct {
		Items []models.JobApplication `json:"items"`
	}
	decodeJSON(t, listResp, &list)
	if len(list.Items) != 1 || list.Items[0].ID != created.ID {
		t.Fatalf("admin list omitted another user's application: %+v", list.Items)
	}

	updateResp := doRequest(t, app, http.MethodPut, path, adminToken, map[string]any{"status": "rejected", "notes": "admin edit"})
	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("admin update: expected 200, got %d", updateResp.StatusCode)
	}
	var updated models.JobApplication
	decodeJSON(t, updateResp, &updated)
	if updated.UserID != created.UserID || updated.Status != models.JobApplicationStatusRejected {
		t.Fatalf("unexpected admin update result: %+v", updated)
	}

	// Admins cannot create an application under another user's job description.
	crossCreate := doRequest(t, app, http.MethodPost, "/api/v1/job-applications", adminToken,
		map[string]any{"job_description_id": jd.ID.String()})
	if crossCreate.StatusCode != http.StatusNotFound {
		t.Fatalf("admin create on another user's job description: expected 404, got %d", crossCreate.StatusCode)
	}
	crossCreate.Body.Close()

	deleteResp := doRequest(t, app, http.MethodDelete, path, adminToken, nil)
	if deleteResp.StatusCode != http.StatusNoContent {
		t.Fatalf("admin delete: expected 204, got %d", deleteResp.StatusCode)
	}
	deleteResp.Body.Close()
}

func TestJobDescriptionDeleteBlockedByApplications(t *testing.T) {
	app := testApp(t)
	_, token := newJobDescriptionUser(t, app)
	jd := createJobDescription(t, app, token, jobDescriptionBody("Posting"))
	application := createJobApplication(t, app, token, map[string]any{"job_description_id": jd.ID.String()})
	jdPath := "/api/v1/job-descriptions/" + jd.ID.String()

	blocked := doRequest(t, app, http.MethodDelete, jdPath, token, nil)
	if blocked.StatusCode != http.StatusConflict {
		t.Fatalf("delete job description with applications: expected 409, got %d", blocked.StatusCode)
	}
	blocked.Body.Close()

	deleteApp := doRequest(t, app, http.MethodDelete, "/api/v1/job-applications/"+application.ID.String(), token, nil)
	if deleteApp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete application: expected 204, got %d", deleteApp.StatusCode)
	}
	deleteApp.Body.Close()

	allowed := doRequest(t, app, http.MethodDelete, jdPath, token, nil)
	if allowed.StatusCode != http.StatusNoContent {
		t.Fatalf("delete job description without applications: expected 204, got %d", allowed.StatusCode)
	}
	allowed.Body.Close()
}

func createJobApplication(t *testing.T, app *fiber.App, token string, body map[string]any) models.JobApplication {
	t.Helper()
	resp := doRequest(t, app, http.MethodPost, "/api/v1/job-applications", token, body)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create job application: expected 201, got %d", resp.StatusCode)
	}
	var created models.JobApplication
	decodeJSON(t, resp, &created)
	return created
}
