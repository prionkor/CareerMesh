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
	if first.OriginalText != firstText {
		t.Fatalf("create changed original text: got %q, want %q", first.OriginalText, firstText)
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
	if updated.OriginalText != updateBody["original_text"].(string) {
		t.Fatalf("update did not replace original text verbatim: got %q", updated.OriginalText)
	}
	if updated.CompanyWebsite != "" || updated.SourceURL != "" || updated.EmploymentType != nil {
		t.Fatalf("update did not clear optional values: %+v", updated)
	}

	getResp := doRequest(t, app, http.MethodGet, "/api/v1/job-descriptions/"+first.ID.String(), token, nil)
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("get job description: expected 200, got %d", getResp.StatusCode)
	}
	var fetched models.JobDescription
	decodeJSON(t, getResp, &fetched)
	if fetched.OriginalText != updated.OriginalText {
		t.Fatalf("get original_text = %q, want %q", fetched.OriginalText, updated.OriginalText)
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
			if created.CompanyWebsite != test.companyURL || created.SourceURL != test.sourceURL {
				t.Fatalf("optional URLs = (%q, %q), want (%q, %q)", created.CompanyWebsite, created.SourceURL, test.companyURL, test.sourceURL)
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
