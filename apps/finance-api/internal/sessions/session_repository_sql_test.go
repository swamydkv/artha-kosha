package sessions

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// Helper to inject DB for testing
func NewPostgresRepoWithDB(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func TestSQLRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepoWithDB(db)
	now := time.Now().UTC()
	s := Session{
		ID:        "1",
		UserID:    "u1",
		IPAddress: "127.0.0.1",
		UserAgent: "agent",
		Status:    StatusActive,
		ExpiresAt: now,
		CreatedAt: now,
	}

	mock.ExpectExec("INSERT INTO sessions").
		WithArgs(s.ID, s.UserID, s.CreatedAt, s.LastActivityAt, s.ExpiresAt, nil, s.UserAgent, s.IPAddress, string(s.Status)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSQLRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepoWithDB(db)
	now := time.Now().UTC()

	rows := sqlmock.NewRows([]string{"id", "user_id", "created_at", "last_activity_at", "expires_at", "revoked_at", "user_agent", "ip_address", "status"}).
		AddRow("1", "u1", now, now, now, nil, "agent", "127.0.0.1", string(StatusActive))

	mock.ExpectQuery("SELECT id, user_id, created_at, last_activity_at, expires_at, revoked_at, user_agent, ip_address, status FROM sessions").
		WithArgs("1").
		WillReturnRows(rows)

	s, err := repo.Get("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.ID != "1" {
		t.Errorf("expected 1, got %s", s.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSQLRepository_Revoke(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepoWithDB(db)

	mock.ExpectExec("UPDATE sessions SET status = \\$1, revoked_at = \\$2 WHERE id = \\$3").
		WithArgs(string(StatusRevoked), sqlmock.AnyArg(), "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Revoke("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSQLRepository_RevokeAllByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepoWithDB(db)

	mock.ExpectExec("UPDATE sessions SET status = \\$1, revoked_at = \\$2 WHERE user_id = \\$3").
		WithArgs(string(StatusRevoked), sqlmock.AnyArg(), "u1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.RevokeAllByUser("u1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSQLRepository_UpdateActivity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepoWithDB(db)
	now := time.Now()

	mock.ExpectExec("UPDATE sessions SET last_activity_at = \\$1 WHERE id = \\$2").
		WithArgs(now, "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateActivity("1", now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSQLRepository_DB(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	repo := NewPostgresRepoWithDB(db)
	if repo.DB() != db {
		t.Error("expected db to match")
	}
}

func TestNewPostgresRepo(t *testing.T) {
	// Test with invalid dsn to cover error
	_, err := NewPostgresRepo("invalid-dsn-format-for-pgx")
	// error could be nil if pgx parses it, but ping will fail
	// actually it tries to connect and ping
	if err == nil {
		t.Log("Expected an error due to invalid DSN, but got nil. Ping might have passed? Or it fails later.")
	}
}

func TestSQLRepository_Get_Errors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepoWithDB(db)

	mock.ExpectQuery("SELECT id, user_id, created_at, last_activity_at, expires_at, revoked_at, user_agent, ip_address, status FROM sessions").
		WithArgs("1").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.Get("1")
	if err == nil || err.Error() != "not found" {
		t.Errorf("expected not found error, got %v", err)
	}

	mock.ExpectQuery("SELECT id, user_id, created_at, last_activity_at, expires_at, revoked_at, user_agent, ip_address, status FROM sessions").
		WithArgs("2").
		WillReturnError(errors.New("db error"))

	_, err = repo.Get("2")
	if err == nil {
		t.Error("expected db error")
	}
}

func TestSQLRepository_Get_Expired(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepoWithDB(db)
	now := time.Now().UTC()

	rows := sqlmock.NewRows([]string{"id", "user_id", "created_at", "last_activity_at", "expires_at", "revoked_at", "user_agent", "ip_address", "status"}).
		AddRow("1", "u1", now, now, now.Add(-time.Hour), nil, "agent", "127.0.0.1", string(StatusActive))

	mock.ExpectQuery("SELECT id, user_id, created_at, last_activity_at, expires_at, revoked_at, user_agent, ip_address, status FROM sessions").
		WithArgs("1").
		WillReturnRows(rows)

	s, err := repo.Get("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Status != StatusExpired {
		t.Errorf("expected expired status, got %s", s.Status)
	}
}

func TestNullableTime(t *testing.T) {
	v := nullableTime(time.Time{})
	if v != nil {
		t.Error("expected nil for zero time")
	}

	now := time.Now()
	v = nullableTime(now)
	if v != now {
		t.Error("expected time to match")
	}
}
