package users

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func TestSQLRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)

	req := CreateUserRequest{
		FullName:     "John Doe",
		DateOfBirth:  time.Now(),
		MobileNumber: "1234567890",
		Email:        "john@example.com",
		Username:     "johndoe",
		PasswordHash: "hashed_password",
	}

	// TDD: Set up mock expectations but they will fail since method is not implemented
	mock.ExpectBegin()
	mock.ExpectQuery(`(?is).*INSERT INTO users.*`).
		WithArgs(req.FullName, req.DateOfBirth, req.MobileNumber, req.Email, req.Username, req.PasswordHash).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "email", "full_name", "created_at"}).
			AddRow(uuid.New(), req.Username, req.Email, req.FullName, time.Now()))
	
	// Expect inserts into domain_events and audit_events
	mock.ExpectExec(`(?is).*INSERT INTO domain_events.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`(?is).*INSERT INTO audit_events.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user, err := repo.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Username != req.Username {
		t.Fatalf("expected username %s, got %s", req.Username, user.Username)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSQLRepository_GetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)

	mock.ExpectQuery(`(?is).*SELECT .* FROM users WHERE username = \$1.*`).
		WithArgs("johndoe").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "email", "full_name", "mobile_number", "password_hash", "created_at", "updated_at"}).
			AddRow(uuid.New(), "johndoe", "john@example.com", "John Doe", "1234567890", "hash", time.Now(), time.Now()))

	user, err := repo.GetUserByUsername(context.Background(), "johndoe")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Username != "johndoe" {
		t.Fatalf("expected username johndoe, got %s", user.Username)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSQLRepository_CreateUser_Duplicate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)

	req := CreateUserRequest{
		FullName:     "Duplicate User",
		DateOfBirth:  time.Now(),
		MobileNumber: "0987654321",
		Email:        "duplicate@example.com",
		Username:     "duplicate",
		PasswordHash: "hash",
	}

	mock.ExpectBegin()
	// TDD: we simulate a Postgres unique constraint violation
	mock.ExpectQuery(`(?is).*INSERT INTO users.*`).
		WithArgs(req.FullName, req.DateOfBirth, req.MobileNumber, req.Email, req.Username, req.PasswordHash).
		WillReturnError(&pq.Error{
			Code:    "23505",
			Message: "duplicate key value violates unique constraint \"users_username_key\"",
		})
	mock.ExpectRollback()

	_, err = repo.CreateUser(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	
	// Test email violation
	mock.ExpectBegin()
	mock.ExpectQuery(`(?is).*INSERT INTO users.*`).WillReturnError(&pq.Error{Code: "23505", Message: "duplicate key value violates unique constraint \"users_email_key\""})
	mock.ExpectRollback()
	_, err = repo.CreateUser(context.Background(), req)
	if err == nil || err.Error() != "email already exists" {
		t.Fatalf("expected 'email already exists', got %v", err)
	}

	// Test mobile number violation
	mock.ExpectBegin()
	mock.ExpectQuery(`(?is).*INSERT INTO users.*`).WillReturnError(&pq.Error{Code: "23505", Message: "duplicate key value violates unique constraint \"users_mobile_number_key\""})
	mock.ExpectRollback()
	_, err = repo.CreateUser(context.Background(), req)
	if err == nil || err.Error() != "mobile number already exists" {
		t.Fatalf("expected 'mobile number already exists', got %v", err)
	}

	// Test generic constraint violation
	mock.ExpectBegin()
	mock.ExpectQuery(`(?is).*INSERT INTO users.*`).WillReturnError(&pq.Error{Code: "23505", Message: "duplicate key value violates unique constraint \"some_other_key\""})
	mock.ExpectRollback()
	_, err = repo.CreateUser(context.Background(), req)
	if err == nil || err.Error() != "duplicate user entry" {
		t.Fatalf("expected 'duplicate user entry', got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSQLRepository_CreateUser_Errors(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewSQLRepository(db)
	req := CreateUserRequest{Username: "test"}

	// Test BeginTx failure
	mock.ExpectBegin().WillReturnError(errors.New("begin err"))
	_, err := repo.CreateUser(context.Background(), req)
	if err == nil { t.Error("expected err") }

	// Test generic CreateUser failure
	mock.ExpectBegin()
	mock.ExpectQuery(`(?is).*INSERT INTO users.*`).WillReturnError(errors.New("insert err"))
	mock.ExpectRollback()
	_, err = repo.CreateUser(context.Background(), req)
	if err == nil { t.Error("expected err") }

	// Test InsertDomainEvent failure
	mock.ExpectBegin()
	mock.ExpectQuery(`(?is).*INSERT INTO users.*`).WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(uuid.New()))
	mock.ExpectExec(`(?is).*INSERT INTO domain_events.*`).WillReturnError(errors.New("domain err"))
	mock.ExpectRollback()
	_, err = repo.CreateUser(context.Background(), req)
	if err == nil { t.Error("expected err") }

	// Test InsertAuditEvent failure
	mock.ExpectBegin()
	mock.ExpectQuery(`(?is).*INSERT INTO users.*`).WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(uuid.New()))
	mock.ExpectExec(`(?is).*INSERT INTO domain_events.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`(?is).*INSERT INTO audit_events.*`).WillReturnError(errors.New("audit err"))
	mock.ExpectRollback()
	_, err = repo.CreateUser(context.Background(), req)
	if err == nil { t.Error("expected err") }

	// Test Commit failure
	mock.ExpectBegin()
	mock.ExpectQuery(`(?is).*INSERT INTO users.*`).WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(uuid.New()))
	mock.ExpectExec(`(?is).*INSERT INTO domain_events.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`(?is).*INSERT INTO audit_events.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(errors.New("commit err"))
	_, err = repo.CreateUser(context.Background(), req)
	if err == nil { t.Error("expected err") }
}

func TestSQLRepository_GetUserByUsername_Errors(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewSQLRepository(db)

	mock.ExpectQuery(`(?is).*SELECT .* FROM users.*`).WillReturnError(sql.ErrNoRows)
	_, err := repo.GetUserByUsername(context.Background(), "test")
	if err != sql.ErrNoRows { t.Error("expected sql.ErrNoRows") }

	mock.ExpectQuery(`(?is).*SELECT .* FROM users.*`).WillReturnError(errors.New("other err"))
	_, err = repo.GetUserByUsername(context.Background(), "test")
	if err == nil { t.Error("expected err") }
}

func TestSQLRepository_Unimplemented(t *testing.T) {
	repo := NewSQLRepository(nil)
	ctx := context.Background()
	_, err := repo.GetUserByID(ctx, "1")
	if err == nil { t.Error("expected err") }
	_, err = repo.GetUserByEmail(ctx, "a@b.com")
	if err == nil { t.Error("expected err") }
	_, err = repo.GetUserByMobileNumber(ctx, "123")
	if err == nil { t.Error("expected err") }
	_, err = repo.CheckUserExists(ctx, "1", "2", "3")
	if err == nil { t.Error("expected err") }
}

func TestSQLRepository_CreateSession(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewSQLRepository(db)
	ctx := context.Background()

	sid := uuid.New().String()
	uid := uuid.New().String()

	// Success
	mock.ExpectBegin()
	mock.ExpectExec(`(?is).*INSERT INTO sessions.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`(?is).*INSERT INTO domain_events.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`(?is).*INSERT INTO audit_events.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreateSession(ctx, sid, uid, "127.0.0.1", "agent")
	if err != nil { t.Errorf("expected nil err, got %v", err) }

	// Errors
	mock.ExpectBegin().WillReturnError(errors.New("begin err"))
	_ = repo.CreateSession(ctx, sid, uid, "", "")

	mock.ExpectBegin()
	mock.ExpectExec(`(?is).*INSERT INTO sessions.*`).WillReturnError(errors.New("session err"))
	mock.ExpectRollback()
	_ = repo.CreateSession(ctx, sid, uid, "", "")

	mock.ExpectBegin()
	mock.ExpectExec(`(?is).*INSERT INTO sessions.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`(?is).*INSERT INTO domain_events.*`).WillReturnError(errors.New("domain err"))
	mock.ExpectRollback()
	_ = repo.CreateSession(ctx, sid, uid, "", "")

	mock.ExpectBegin()
	mock.ExpectExec(`(?is).*INSERT INTO sessions.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`(?is).*INSERT INTO domain_events.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`(?is).*INSERT INTO audit_events.*`).WillReturnError(errors.New("audit err"))
	mock.ExpectRollback()
	_ = repo.CreateSession(ctx, sid, uid, "", "")

	mock.ExpectBegin()
	mock.ExpectExec(`(?is).*INSERT INTO sessions.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`(?is).*INSERT INTO domain_events.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`(?is).*INSERT INTO audit_events.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(errors.New("commit err"))
	_ = repo.CreateSession(ctx, sid, uid, "", "")
}
