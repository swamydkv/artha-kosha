package sqlc

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDBMethods(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	queries := New(db)

	mock.ExpectBegin()
	tx, _ := db.Begin()
	mock.ExpectCommit()

	qTx := queries.WithTx(tx)
	if qTx == nil {
		t.Errorf("expected *Queries, got nil")
	}
	tx.Commit()
}
