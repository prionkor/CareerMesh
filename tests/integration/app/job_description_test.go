package app_test

import (
	"net/http"
	"testing"

	"uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/prionkor/careermesh/internal/auth"
	"github.com/prionkor/careermesh/models"
)

func TestJobDescriptionCRUDAndPagination(t *testing.T) {
	app := testApp(t)
	_, token := newJobDescriptionUser(t, app)

	firstText := "  Build reliable systems.\nKeep this text verbatim.  "
	first := createJobDescription(t, app, token, jobDescriptionBody(firstText))
	if first.CurrentVersion == nil || first.CurrentVersion.VersionNumber != 1 || first.CurrentVersion.OriginalText != firstText {
		t.Fatalf("create did not preserve the initial version: %+v", first.CurrentVersion)
	}

	second := createJobDescription(t, app, token, jobDescriptionBody("Second posting"))
	updateBody := jobDescriptionBody("\nUpdated source text\n")
	updateBody["company_website"] = ""
	updateBody["source_url"] = ""
	updateBody["employment_type"] = nil
	updateResp := doRequest(t, app, http.MethodPut, "/api/v1/job-descriptions/"+first.ID.String(), token, updateBody)
	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("update job description: expected 200, got %d", updateResp.StatusCode)
	}
	var updated models.JobDescription
	decodeJSON(t, updateResp, &updated)
	if updated.CurrentVersion == nil || updated.CurrentVersion.VersionNumber != 2 ||
		updated.CurrentVersionID == first.CurrentVersionID ||
		updated.CurrentVersion.OriginalText != updateBody["original_text"].(string) {
		t.Fatalf("update did not advance to a new current version: %+v", updated)
	}
	if updated.CurrentVersion.CompanyWebsite != "" || updated.CurrentVersion.SourceURL != "" || updated.CurrentVersion.EmploymentType != nil {
		t.Fatalf("update did not clear optional values in the new version: %+v", updated.CurrentVersion)
	}

	versionsResp := doRequest(t, app, http.MethodGet, "/api/v1/job-descriptions/"+first.ID.String()+"/versions", token, nil)
	if versionsResp.StatusCode != http.StatusOK {
		t.Fatalf("list job description versions: expected 200, got %d", versionsResp.StatusCode)
	}
	var versionsPage struct {
		Items []models.JobDescriptionVersion `json:"items"`
	}
	decodeJSON(t, versionsResp, &versionsPage)
	if len(versionsPage.Items) != 2 || versionsPage.Items[0].ID != first.CurrentVersionID ||
		versionsPage.Items[0].VersionNumber != 1 || versionsPage.Items[0].OriginalText != firstText ||
		versionsPage.Items[1].ID != updated.CurrentVersionID || versionsPage.Items[1].VersionNumber != 2 {
		t.Fatalf("unexpected job description version history: %+v", versionsPage.Items)
	}

	versionResp := doRequest(t, app, http.MethodGet,
		"/api/v1/job-descriptions/"+first.ID.String()+"/versions/"+first.CurrentVersionID.String(), token, nil)
	if versionResp.StatusCode != http.StatusOK {
		t.Fatalf("get specific job description version: expected 200, got %d", versionResp.StatusCode)
	}
	var historicalVersion models.JobDescriptionVersion
	decodeJSON(t, versionResp, &historicalVersion)
	if historicalVersion.ID != first.CurrentVersionID || historicalVersion.OriginalText != firstText {
		t.Fatalf("specific version changed after update: %+v", historicalVersion)
	}

	wrongParentVersionResp := doRequest(t, app, http.MethodGet,
		"/api/v1/job-descriptions/"+first.ID.String()+"/versions/"+second.CurrentVersionID.String(), token, nil)
	if wrongParentVersionResp.StatusCode != http.StatusNotFound {
		t.Fatalf("version belonging to another job description: expected 404, got %d", wrongParentVersionResp.StatusCode)
	}

	getResp := doRequest(t, app, http.MethodGet, "/api/v1/job-descriptions/"+first.ID.String(), token, nil)
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("get job description: expected 200, got %d", getResp.StatusCode)
	}
	var fetched models.JobDescription
	decodeJSON(t, getResp, &fetched)
	if fetched.CurrentVersionID != updated.CurrentVersionID || fetched.CurrentVersion == nil ||
		fetched.CurrentVersion.OriginalText != updated.CurrentVersion.OriginalText {
		t.Fatalf("get did not return the current version: %+v", fetched)
	}

	listResp := doRequest(t, app, http.MethodGet, "/api/v1/job-descriptions?limit=1&offset=0", token, nil)
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("list job descriptions: expected 200, got %d", listResp.StatusCode)
	}
	var firstPage struct {
		Items  []models.JobDescription `json:"items"`
		Limit  int                     `json:"limit"`
		Offset int                     `json:"offset"`
	}
	decodeJSON(t, listResp, &firstPage)
	if len(firstPage.Items) != 1 || firstPage.Limit != 1 || firstPage.Offset != 0 {
		t.Fatalf("unexpected first page: %+v", firstPage)
	}

	secondPageResp := doRequest(t, app, http.MethodGet, "/api/v1/job-descriptions?limit=1&offset=1", token, nil)
	if secondPageResp.StatusCode != http.StatusOK {
		t.Fatalf("list second page: expected 200, got %d", secondPageResp.StatusCode)
	}
	var secondPage struct {
		Items []models.JobDescription `json:"items"`
	}
	decodeJSON(t, secondPageResp, &secondPage)
	if len(secondPage.Items) != 1 || secondPage.Items[0].ID == firstPage.Items[0].ID {
		t.Fatalf("unexpected second page: %+v", secondPage)
	}

	deleteResp := doRequest(t, app, http.MethodDelete, "/api/v1/job-descriptions/"+second.ID.String(), token, nil)
	if deleteResp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete job description: expected 204, got %d", deleteResp.StatusCode)
	}
	missingResp := doRequest(t, app, http.MethodGet, "/api/v1/job-descriptions/"+second.ID.String(), token, nil)
	if missingResp.StatusCode != http.StatusNotFound {
		t.Fatalf("get deleted job description: expected 404, got %d", missingResp.StatusCode)
	}
}

func TestJobDescriptionValidationAndOwnership(t *testing.T) {
	app := testApp(t)
	_, ownerToken := newJobDescriptionUser(t, app)
	_, otherToken := newJobDescriptionUser(t, app)

	validBody := jobDescriptionBody("Original posting")
	created := createJobDescription(t, app, ownerToken, validBody)

	requiredFields := []string{"title", "company", "original_text"}
	for _, field := range requiredFields {
		for _, test := range []struct {
			name    string
			value   string
			present bool
		}{
			{name: "missing"},
			{name: "empty", present: true},
			{name: "whitespace", value: " \t\n", present: true},
		} {
			t.Run(field+"/"+test.name, func(t *testing.T) {
				body := jobDescriptionBody("Original posting")
				if test.present {
					body[field] = test.value
				} else {
					delete(body, field)
				}
				resp := doRequest(t, app, http.MethodPost, "/api/v1/job-descriptions", ownerToken, body)
				if resp.StatusCode != http.StatusBadRequest {
					t.Fatalf("invalid %s: expected 400, got %d", field, resp.StatusCode)
				}
				if field == "title" && test.name == "missing" {
					if body := readBody(t, resp); body != `{"error":"invalid request body"}` {
						t.Fatalf("unexpected validation response: %s", body)
					}
				} else {
					resp.Body.Close()
				}
			})
		}
	}

	optionalURLCases := []struct {
		name       string
		companyURL string
		sourceURL  string
		omitURLs   bool
	}{
		{name: "http and https", companyURL: "http://company.example", sourceURL: "https://jobs.example/role"},
		{name: "https and http", companyURL: "https://company.example", sourceURL: "http://jobs.example/role"},
		{name: "omitted", omitURLs: true},
		{name: "empty"},
	}
	for _, test := range optionalURLCases {
		t.Run("optional URLs/"+test.name, func(t *testing.T) {
			body := jobDescriptionBody("Posting with optional URLs")
			if test.omitURLs {
				delete(body, "company_website")
				delete(body, "source_url")
			} else {
				body["company_website"] = test.companyURL
				body["source_url"] = test.sourceURL
			}
			created := createJobDescription(t, app, ownerToken, body)
			if created.CurrentVersion == nil || created.CurrentVersion.CompanyWebsite != test.companyURL ||
				created.CurrentVersion.SourceURL != test.sourceURL {
				t.Fatalf("optional URLs = (%+v), want (%q, %q)", created.CurrentVersion, test.companyURL, test.sourceURL)
			}
		})
	}

	invalidBodies := []map[string]any{
		{"title": "Engineer", "company": "Example", "original_text": "text", "employment_type": "seasonal"},
		{"title": "Engineer", "company": "Example", "original_text": "text", "source_url": "/jobs/1"},
		{"title": "Engineer", "company": "Example", "original_text": "text", "company_website": "ftp://example.com"},
		{"title": "Engineer", "company": "Example", "original_text": "text", "source_url": "https://:443"},
	}
	for index, body := range invalidBodies {
		resp := doRequest(t, app, http.MethodPost, "/api/v1/job-descriptions", ownerToken, body)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("invalid body %d: expected 400, got %d", index, resp.StatusCode)
		}
		resp.Body.Close()
	}

	getResp := doRequest(t, app, http.MethodGet, "/api/v1/job-descriptions/"+created.ID.String(), otherToken, nil)
	if getResp.StatusCode != http.StatusForbidden {
		t.Fatalf("other user get: expected 403, got %d", getResp.StatusCode)
	}
	getResp.Body.Close()

	versionsResp := doRequest(t, app, http.MethodGet,
		"/api/v1/job-descriptions/"+created.ID.String()+"/versions", otherToken, nil)
	if versionsResp.StatusCode != http.StatusForbidden {
		t.Fatalf("other user list versions: expected 403, got %d", versionsResp.StatusCode)
	}
	versionsResp.Body.Close()

	versionResp := doRequest(t, app, http.MethodGet,
		"/api/v1/job-descriptions/"+created.ID.String()+"/versions/"+created.CurrentVersionID.String(), otherToken, nil)
	if versionResp.StatusCode != http.StatusForbidden {
		t.Fatalf("other user get version: expected 403, got %d", versionResp.StatusCode)
	}
	versionResp.Body.Close()

	updateResp := doRequest(t, app, http.MethodPut, "/api/v1/job-descriptions/"+created.ID.String(), otherToken, validBody)
	if updateResp.StatusCode != http.StatusForbidden {
		t.Fatalf("other user update: expected 403, got %d", updateResp.StatusCode)
	}
	updateResp.Body.Close()

	deleteResp := doRequest(t, app, http.MethodDelete, "/api/v1/job-descriptions/"+created.ID.String(), otherToken, nil)
	if deleteResp.StatusCode != http.StatusForbidden {
		t.Fatalf("other user delete: expected 403, got %d", deleteResp.StatusCode)
	}
	deleteResp.Body.Close()

	listResp := doRequest(t, app, http.MethodGet, "/api/v1/job-descriptions", otherToken, nil)
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("other user list: expected 200, got %d", listResp.StatusCode)
	}
	var list struct {
		Items []models.JobDescription `json:"items"`
	}
	decodeJSON(t, listResp, &list)
	if len(list.Items) != 0 {
		t.Fatalf("other user list returned another user's job descriptions: %+v", list.Items)
	}

	invalidIDResp := doRequest(t, app, http.MethodGet, "/api/v1/job-descriptions/not-a-uuid", ownerToken, nil)
	if invalidIDResp.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid ID: expected 400, got %d", invalidIDResp.StatusCode)
	}
	invalidIDResp.Body.Close()

	invalidVersionIDResp := doRequest(t, app, http.MethodGet,
		"/api/v1/job-descriptions/"+created.ID.String()+"/versions/not-a-uuid", ownerToken, nil)
	if invalidVersionIDResp.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid version ID: expected 400, got %d", invalidVersionIDResp.StatusCode)
	}
	invalidVersionIDResp.Body.Close()

	invalidLimitResp := doRequest(t, app, http.MethodGet, "/api/v1/job-descriptions?limit=0", ownerToken, nil)
	if invalidLimitResp.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid limit: expected 400, got %d", invalidLimitResp.StatusCode)
	}
	invalidLimitResp.Body.Close()

	unauthorizedResp := doRequest(t, app, http.MethodGet, "/api/v1/job-descriptions", "", nil)
	if unauthorizedResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthenticated list: expected 401, got %d", unauthorizedResp.StatusCode)
	}
	unauthorizedResp.Body.Close()

	noPermissionsToken, err := auth.GenerateAccessToken([]byte(testJWTSecret), uuid.NewV7(), []string{})
	if err != nil {
		t.Fatalf("generate token without permissions: %v", err)
	}
	forbiddenResp := doRequest(t, app, http.MethodGet, "/api/v1/job-descriptions", noPermissionsToken, nil)
	if forbiddenResp.StatusCode != http.StatusForbidden {
		t.Fatalf("token without permissions: expected 403, got %d", forbiddenResp.StatusCode)
	}
	forbiddenResp.Body.Close()
}

func TestJobDescriptionAdminCanAccessAnotherUsersResource(t *testing.T) {
	app := testApp(t)
	_, ownerToken := newJobDescriptionUser(t, app)
	adminIdentity := createVerifiedAuthIdentity(t, app)
	created := createJobDescription(t, app, ownerToken, jobDescriptionBody("Private posting"))

	setAuthTestRole(t, testDB(t), adminIdentity.userID, "admin")
	adminToken := loginAuthIdentity(t, app, adminIdentity)

	getResp := doRequest(t, app, http.MethodGet, "/api/v1/job-descriptions/"+created.ID.String(), adminToken, nil)
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("admin get: expected 200, got %d", getResp.StatusCode)
	}
	getResp.Body.Close()

	versionsResp := doRequest(t, app, http.MethodGet,
		"/api/v1/job-descriptions/"+created.ID.String()+"/versions", adminToken, nil)
	if versionsResp.StatusCode != http.StatusOK {
		t.Fatalf("admin list versions: expected 200, got %d", versionsResp.StatusCode)
	}
	versionsResp.Body.Close()

	versionResp := doRequest(t, app, http.MethodGet,
		"/api/v1/job-descriptions/"+created.ID.String()+"/versions/"+created.CurrentVersionID.String(), adminToken, nil)
	if versionResp.StatusCode != http.StatusOK {
		t.Fatalf("admin get version: expected 200, got %d", versionResp.StatusCode)
	}
	versionResp.Body.Close()

	listResp := doRequest(t, app, http.MethodGet, "/api/v1/job-descriptions", adminToken, nil)
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("admin list: expected 200, got %d", listResp.StatusCode)
	}
	var list struct {
		Items []models.JobDescription `json:"items"`
	}
	decodeJSON(t, listResp, &list)
	found := false
	for _, item := range list.Items {
		if item.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("admin list omitted resource owned by another user: %s", created.ID)
	}
}

func createJobDescription(t *testing.T, app *fiber.App, token string, body map[string]any) models.JobDescription {
	t.Helper()
	resp := doRequest(t, app, http.MethodPost, "/api/v1/job-descriptions", token, body)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create job description: expected 201, got %d", resp.StatusCode)
	}
	var created models.JobDescription
	decodeJSON(t, resp, &created)
	return created
}

func newJobDescriptionUser(t *testing.T, app *fiber.App) (string, string) {
	t.Helper()
	identity := createVerifiedAuthIdentity(t, app)
	return identity.userID, loginAuthIdentity(t, app, identity)
}

func jobDescriptionBody(originalText string) map[string]any {
	return map[string]any{
		"title":           "Senior Engineer",
		"company":         "Example Corp",
		"company_website": "https://example.com",
		"source_url":      "https://jobs.example.com/role",
		"location":        "Remote",
		"employment_type": "full_time",
		"original_text":   originalText,
	}
}
