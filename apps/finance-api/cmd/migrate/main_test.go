package main

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMigrate_NoDSN(t *testing.T) {
	os.Setenv("DATABASE_URL", "")
	err := run()
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMigrate_Success(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://test")
	os.Setenv("TEST_MODE", "") // turn off test mode skip

	db, mock, _ := sqlmock.New()
	
	origOpen := openDB
	openDB = func(driverName, dataSourceName string) (*sql.DB, error) { return db, nil }
	defer func() { openDB = origOpen }()

	// create a dummy migration file
	d, err := ioutil.TempDir("", "migrations")
	if err != nil { t.Fatalf("err: %v", err) }
	defer os.RemoveAll(d)

	err = ioutil.WriteFile(filepath.Join(d, "001.sql"), []byte("CREATE TABLE dummy (id int);"), 0644)
	if err != nil { t.Fatalf("err: %v", err) }

	os.Args = []string{"migrate", d}
	defer func() { os.Args = []string{"migrate"} }()

	mock.ExpectExec("CREATE TABLE dummy").WillReturnResult(sqlmock.NewResult(0, 0))

	err = run()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
