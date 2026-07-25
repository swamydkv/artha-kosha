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

func TestCreateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)

	mock.ExpectBegin()

	// CreateUser query
	rows := sqlmock.NewRows([]string{"user_id", "username", "email", "full_name", "created_at"}).
		AddRow(uuid.New(), "testuser", "test@example.com", "Test User", time.Now())
	mock.ExpectQuery("INSERT INTO users").WillReturnRows(rows)

	// DomainEvent
	mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))

	// OutboxEntry
	mock.ExpectExec("INSERT INTO transactional_outbox").WillReturnResult(sqlmock.NewResult(1, 1))

	// AuditEvent
	mock.ExpectExec("INSERT INTO audit_events").WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	_, err = repo.CreateUser(context.Background(), CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
	})

	if err != nil {
		t.Errorf("error was not expected while creating user: %s", err)
	}
}

func TestCreateUser_Errors(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewSQLRepository(db)
	defer db.Close()

	// 1. BeginTx error
	mock.ExpectBegin().WillReturnError(errors.New("begin error"))
	_, _ = repo.CreateUser(context.Background(), CreateUserRequest{})

	// 2. CreateUser error (Duplicate)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO users").WillReturnError(&pq.Error{Code: "23505", Message: "users_username_key"})
	mock.ExpectRollback()
	_, _ = repo.CreateUser(context.Background(), CreateUserRequest{})

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO users").WillReturnError(&pq.Error{Code: "23505", Message: "users_email_key"})
	mock.ExpectRollback()
	_, _ = repo.CreateUser(context.Background(), CreateUserRequest{})

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO users").WillReturnError(&pq.Error{Code: "23505", Message: "users_mobile_number_key"})
	mock.ExpectRollback()
	_, _ = repo.CreateUser(context.Background(), CreateUserRequest{})

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO users").WillReturnError(&pq.Error{Code: "23505", Message: "other_key"})
	mock.ExpectRollback()
	_, _ = repo.CreateUser(context.Background(), CreateUserRequest{})

	// 3. Domain Event error
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"user_id", "username", "email", "full_name", "created_at"}).
		AddRow(uuid.New(), "testuser", "test@example.com", "Test User", time.Now())
	mock.ExpectQuery("INSERT INTO users").WillReturnRows(rows)
	mock.ExpectExec("INSERT INTO domain_events").WillReturnError(errors.New("de err"))
	mock.ExpectRollback()
	_, _ = repo.CreateUser(context.Background(), CreateUserRequest{})

	// 4. Outbox error
	mock.ExpectBegin()
	rows2 := sqlmock.NewRows([]string{"user_id", "username", "email", "full_name", "created_at"}).
		AddRow(uuid.New(), "testuser", "test@example.com", "Test User", time.Now())
	mock.ExpectQuery("INSERT INTO users").WillReturnRows(rows2)
	mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO transactional_outbox").WillReturnError(errors.New("outbox err"))
	mock.ExpectRollback()
	_, _ = repo.CreateUser(context.Background(), CreateUserRequest{})

	// 5. Audit error
	mock.ExpectBegin()
	rows3 := sqlmock.NewRows([]string{"user_id", "username", "email", "full_name", "created_at"}).
		AddRow(uuid.New(), "testuser", "test@example.com", "Test User", time.Now())
	mock.ExpectQuery("INSERT INTO users").WillReturnRows(rows3)
	mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO transactional_outbox").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO audit_events").WillReturnError(errors.New("audit err"))
	mock.ExpectRollback()
	_, _ = repo.CreateUser(context.Background(), CreateUserRequest{})

	// 6. Commit error
	mock.ExpectBegin()
	rows4 := sqlmock.NewRows([]string{"user_id", "username", "email", "full_name", "created_at"}).
		AddRow(uuid.New(), "testuser", "test@example.com", "Test User", time.Now())
	mock.ExpectQuery("INSERT INTO users").WillReturnRows(rows4)
	mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO transactional_outbox").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO audit_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(errors.New("commit err"))
	_, _ = repo.CreateUser(context.Background(), CreateUserRequest{})
}

func TestGetUserBy(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewSQLRepository(db)
	defer db.Close()

	_, _ = repo.GetUserByID(context.Background(), "1")
	_, _ = repo.GetUserByEmail(context.Background(), "1")
	_, _ = repo.GetUserByMobileNumber(context.Background(), "1")
	_, _ = repo.CheckUserExists(context.Background(), "1", "1", "1")

	// GetUserByUsername
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
	_, _ = repo.GetUserByUsername(context.Background(), "test")

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("db err"))
	_, _ = repo.GetUserByUsername(context.Background(), "test")

	rows := sqlmock.NewRows([]string{"user_id", "username", "email", "full_name", "mobile_number", "password_hash", "created_at", "updated_at", "deleted_at", "is_archived"}).
		AddRow(uuid.New(), "test", "a@a.com", "t", "1", "hash", time.Now(), time.Now(), nil, false)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	_, _ = repo.GetUserByUsername(context.Background(), "test")
}

func TestCreateSession(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewSQLRepository(db)
	defer db.Close()

	mock.ExpectBegin().WillReturnError(errors.New("begin error"))
	_ = repo.CreateSession(context.Background(), uuid.NewString(), uuid.NewString(), "1.1.1.1", "agent")

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO sessions").WillReturnError(errors.New("insert error"))
	mock.ExpectRollback()
	_ = repo.CreateSession(context.Background(), uuid.NewString(), uuid.NewString(), "1.1.1.1", "agent")

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO sessions").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO domain_events").WillReturnError(errors.New("de error"))
	mock.ExpectRollback()
	_ = repo.CreateSession(context.Background(), uuid.NewString(), uuid.NewString(), "1.1.1.1", "agent")

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO sessions").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO audit_events").WillReturnError(errors.New("audit error"))
	mock.ExpectRollback()
	_ = repo.CreateSession(context.Background(), uuid.NewString(), uuid.NewString(), "1.1.1.1", "agent")

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO sessions").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO audit_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(errors.New("commit error"))
	_ = repo.CreateSession(context.Background(), uuid.NewString(), uuid.NewString(), "1.1.1.1", "agent")

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO sessions").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO audit_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	_ = repo.CreateSession(context.Background(), uuid.NewString(), uuid.NewString(), "1.1.1.1", "agent")
}

func TestDeleteUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewSQLRepository(db)
	defer db.Close()

	// 1. Begin err
	mock.ExpectBegin().WillReturnError(errors.New("begin error"))
	_ = repo.DeleteUser(context.Background(), uuid.NewString(), 30)

	// 2. UUID parse err
	mock.ExpectBegin()
	_ = repo.DeleteUser(context.Background(), "invalid", 30)

	// 3. Query err
	uid := uuid.NewString()
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()
	_ = repo.DeleteUser(context.Background(), uid, 30)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnError(errors.New("db error"))
	mock.ExpectRollback()
	_ = repo.DeleteUser(context.Background(), uid, 30)

	// 4. Insert archived error
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"username", "email", "full_name", "mobile_number", "password_hash"}).
		AddRow("u", "e", "f", "m", "p")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectExec("INSERT INTO archived_users").WillReturnError(errors.New("ins err"))
	mock.ExpectRollback()
	_ = repo.DeleteUser(context.Background(), uid, 30)

	// 5. Domain event error
	mock.ExpectBegin()
	rows2 := sqlmock.NewRows([]string{"username", "email", "full_name", "mobile_number", "password_hash"}).
		AddRow("u", "e", "f", "m", "p")
	mock.ExpectQuery("SELECT").WillReturnRows(rows2)
	mock.ExpectExec("INSERT INTO archived_users").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO domain_events").WillReturnError(errors.New("de err"))
	mock.ExpectRollback()
	_ = repo.DeleteUser(context.Background(), uid, 30)

	// 6. Outbox error
	mock.ExpectBegin()
	rows3 := sqlmock.NewRows([]string{"username", "email", "full_name", "mobile_number", "password_hash"}).
		AddRow("u", "e", "f", "m", "p")
	mock.ExpectQuery("SELECT").WillReturnRows(rows3)
	mock.ExpectExec("INSERT INTO archived_users").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO transactional_outbox").WillReturnError(errors.New("outbox err"))
	mock.ExpectRollback()
	_ = repo.DeleteUser(context.Background(), uid, 30)

	// 7. Audit error
	mock.ExpectBegin()
	rows4 := sqlmock.NewRows([]string{"username", "email", "full_name", "mobile_number", "password_hash"}).
		AddRow("u", "e", "f", "m", "p")
	mock.ExpectQuery("SELECT").WillReturnRows(rows4)
	mock.ExpectExec("INSERT INTO archived_users").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO transactional_outbox").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO audit_events").WillReturnError(errors.New("audit err"))
	mock.ExpectRollback()
	_ = repo.DeleteUser(context.Background(), uid, 30)

	// 8. Commit error
	mock.ExpectBegin()
	rows5 := sqlmock.NewRows([]string{"username", "email", "full_name", "mobile_number", "password_hash"}).
		AddRow("u", "e", "f", "m", "p")
	mock.ExpectQuery("SELECT").WillReturnRows(rows5)
	mock.ExpectExec("INSERT INTO archived_users").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO transactional_outbox").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO audit_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(errors.New("com err"))
	_ = repo.DeleteUser(context.Background(), uid, 30)

	// 9. Success
	mock.ExpectBegin()
	rows6 := sqlmock.NewRows([]string{"username", "email", "full_name", "mobile_number", "password_hash"}).
		AddRow("u", "e", "f", "m", "p")
	mock.ExpectQuery("SELECT").WillReturnRows(rows6)
	mock.ExpectExec("INSERT INTO archived_users").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO transactional_outbox").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO audit_events").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	_ = repo.DeleteUser(context.Background(), uid, 30)
}

func TestPruneArchivedUsers(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewSQLRepository(db)
	defer db.Close()

	mock.ExpectExec("DELETE FROM archived_users").WillReturnResult(sqlmock.NewResult(1, 1))
	_ = repo.PruneArchivedUsers(context.Background())
}
