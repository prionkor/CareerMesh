package testdb

import "testing"

func TestParseAndValidate(t *testing.T) {
	cases := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"test db", "postgres://u:p@localhost:5432/careermesh_test?sslmode=disable", false},
		{"dev db", "postgres://u:p@localhost:5432/careermesh?sslmode=disable", true},
		{"other db", "postgres://u:p@localhost:5432/other", true},
		{"suffix not at end", "postgres://u:p@localhost:5432/careermesh_test_old", true},
		{"dbname override in query", "postgres://u:p@localhost:5432/x_test?dbname=careermesh", true},
		{"malformed", "postgres://u:p@localhost:notaport/x_test", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseAndValidate(tc.url)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestPoolConfigFromEnv_Missing(t *testing.T) {
	t.Setenv(EnvAppEnv, requiredAppEnv)
	t.Setenv(EnvURL, "")
	if _, err := poolConfigFromEnv(); err == nil {
		t.Fatal("expected error when TEST_DATABASE_URL is unset")
	}
}
