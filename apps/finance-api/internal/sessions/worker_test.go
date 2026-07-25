package sessions

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSessionWorker(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := &PostgresRepo{db: db}

	w := NewWorker(repo, 5*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	
	// Expect it to run at least once
	mock.ExpectExec("DELETE FROM sessions").WillReturnResult(sqlmock.NewResult(1, 1))
	// Also test error path by returning an error on the second tick
	mock.ExpectExec("DELETE FROM sessions").WillReturnError(errors.New("db err"))
	
	w.Start(ctx)
	
	// Allow ticker to fire a few times
	time.Sleep(20 * time.Millisecond)
	cancel()
	// Wait a bit to ensure goroutine shuts down cleanly
	time.Sleep(10 * time.Millisecond)
}
