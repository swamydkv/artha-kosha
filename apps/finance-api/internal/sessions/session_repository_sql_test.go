package sessions

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestSessionPostgresRepo(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := &PostgresRepo{db: db}
	defer db.Close()

	// 1. Create - error
	mock.ExpectExec("INSERT INTO sessions").WillReturnError(errors.New("db err"))
	_ = repo.Create(Session{})

	// 2. Create - success
	mock.ExpectExec("INSERT INTO sessions").WillReturnResult(sqlmock.NewResult(1, 1))
	_ = repo.Create(Session{})

	// 3. Get - error
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
	_, _ = repo.Get("id")

	// 4. Get - success
	rows := sqlmock.NewRows([]string{"id", "user_id", "created_at", "last_activity_at", "expires_at", "revoked_at", "user_agent", "ip_address", "status"}).
		AddRow(uuid.NewString(), uuid.NewString(), time.Now(), time.Now(), time.Now(), nil, nil, "127.0.0.1", "active")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	_, _ = repo.Get("id")

	// 5. Revoke - error
	mock.ExpectExec("UPDATE sessions").WillReturnError(errors.New("db err"))
	_ = repo.Revoke("id")

	// 6. Revoke - success
	mock.ExpectExec("UPDATE sessions").WillReturnResult(sqlmock.NewResult(1, 1))
	_ = repo.Revoke("id")

	// 7. RevokeAllByUser - error
	mock.ExpectExec("UPDATE sessions").WillReturnError(errors.New("db err"))
	_ = repo.RevokeAllByUser("id")

	// 8. RevokeAllByUser - success
	mock.ExpectExec("UPDATE sessions").WillReturnResult(sqlmock.NewResult(1, 1))
	_ = repo.RevokeAllByUser("id")

	// 9. UpdateActivity - error
	mock.ExpectExec("UPDATE sessions").WillReturnError(errors.New("db err"))
	_ = repo.UpdateActivity("id", time.Now())

	// 10. UpdateActivity - success
	mock.ExpectExec("UPDATE sessions").WillReturnResult(sqlmock.NewResult(1, 1))
	_ = repo.UpdateActivity("id", time.Now())

	// 11. DeleteExpired - error
	mock.ExpectExec("DELETE FROM sessions").WillReturnError(errors.New("db err"))
	_ = repo.DeleteExpired(context.Background(), time.Now())

	// 12. DeleteExpired - success
	mock.ExpectExec("DELETE FROM sessions").WillReturnResult(sqlmock.NewResult(5, 5))
	_ = repo.DeleteExpired(context.Background(), time.Now())
}
