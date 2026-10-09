// Package testdb provides the shared database harness for integration tests.
package testdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	EnvURL    = "TEST_DATABASE_URL"
	EnvAppEnv = "APP_ENV"

	requiredSuffix = "_test"
	devDatabase    = "careermesh"
	requiredAppEnv = "test"
)

// Setup validates TEST_DATABASE_URL, connects, and applies migrations.
// The caller owns the returned pool and must Close it.
func Setup(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := poolConfigFromEnv()
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create test database pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("cannot connect to test database %q (is it running? see README \"Integration tests\"): %w", cfg.ConnConfig.Database, err)
	}

	if err := migrate(pool); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func poolConfigFromEnv() (*pgxpool.Config, error) {
	if appEnv := os.Getenv(EnvAppEnv); appEnv != requiredAppEnv {
		return nil, fmt.Errorf("refusing to run integration tests: %s must be %q (got %q)", EnvAppEnv, requiredAppEnv, appEnv)
	}

	url := os.Getenv(EnvURL)
	if url == "" {
		return nil, fmt.Errorf("%s is not set; integration tests require an explicit test database URL (e.g. postgres://careermesh:careermesh@localhost:5432/careermesh_test?sslmode=disable)", EnvURL)
	}
	return parseAndValidate(url)
}

// parseAndValidate checks the effective database name without connecting.
func parseAndValidate(url string) (*pgxpool.Config, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", EnvURL, err)
	}
	if err := ValidateDatabaseName(cfg.ConnConfig.Database); err != nil {
		return nil, err
	}
	return cfg, nil
}

// ValidateDatabaseName rejects anything that is not a dedicated test database.
func ValidateDatabaseName(name string) error {
	switch {
	case name == devDatabase:
		return fmt.Errorf("refusing to run integration tests against development database %q", name)
	case !strings.HasSuffix(name, requiredSuffix):
		return fmt.Errorf("refusing to run integration tests against database %q: name must end with %q", name, requiredSuffix)
	}
	return nil
}

func migrate(pool *pgxpool.Pool) error {
	dir, err := migrationsDir()
	if err != nil {
		return err
	}

	// The sql.DB wraps the shared pool; closing it does not close the pool.
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.Up((*sql.DB)(db), dir); err != nil {
		return fmt.Errorf("apply migrations from %s to test database: %w", dir, err)
	}
	return nil
}

// migrationsDir locates the repo's migrations directory by walking up to go.mod.
func migrationsDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}
	for dir := wd; ; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, "migrations"), nil
		}
		if dir == filepath.Dir(dir) {
			return "", errors.New("cannot locate repository root (go.mod) to find migrations")
		}
	}
}
