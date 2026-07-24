package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestWithTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()

	err = WithTx(context.Background(), db, func(tx *sql.Tx) error {
		return nil
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestWithTx_BeginTxError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin().WillReturnError(errors.New("begin error"))

	err = WithTx(context.Background(), db, func(tx *sql.Tx) error {
		return nil
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "begin tx: begin error" {
		t.Errorf("expected begin tx: begin error, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestWithTx_FnErrorRollbacks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()

	err = WithTx(context.Background(), db, func(tx *sql.Tx) error {
		return errors.New("fn error")
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "fn error" {
		t.Errorf("expected fn error, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestWithTx_CommitError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(errors.New("commit error"))

	err = WithTx(context.Background(), db, func(tx *sql.Tx) error {
		return nil
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "commit tx: commit error" {
		t.Errorf("expected commit tx: commit error, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
