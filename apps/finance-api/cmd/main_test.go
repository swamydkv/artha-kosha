package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/sessions"
)

func TestMain_Run_NoDB(t *testing.T) {
	os.Setenv("PORT", "0")
	os.Setenv("DATABASE_URL", "")
	os.Setenv("TEST_MODE", "true")
	
	err := run()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestMain_Run_WithDB(t *testing.T) {
	os.Setenv("PORT", "0")
	os.Setenv("DATABASE_URL", "postgres://test")
	os.Setenv("TEST_MODE", "true")

	// Mock startServer
	origStart := startServer
	startServer = func(server *http.Server) error { return nil }
	defer func() { startServer = origStart }()

	// Mock newAuthProvider
	origNewAuth := newAuthProvider
	newAuthProvider = func(dsn string, ttl time.Duration) (*auth.LocalAuthProvider, *sessions.PostgresRepo, error) {
		db, _, _ := sqlmock.New()
		pgRepo := sessions.NewPostgresRepoFromDB(db)
		return auth.NewLocalAuthProviderFromDB(pgRepo, ttl), pgRepo, nil
	}
	defer func() { newAuthProvider = origNewAuth }()
	
	err := run()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
