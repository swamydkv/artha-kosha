import os
import re

def generate_models_test():
    with open('apps/finance-api/internal/sqlc/models.go', 'r') as f:
        content = f.read()
    
    enums = re.findall(r'type (\w+) string', content)
    
    test_code = """package sqlc

import (
	"testing"
)

"""
    for enum in enums:
        test_code += f"""
func Test{enum}_ScanValue(t *testing.T) {{
    var e {enum}
    var ne Null{enum}
    
    // Test Scan on base type
    if err := e.Scan(string("valid")); err != nil {{
        t.Errorf("expected no error, got %v", err)
    }}
    if err := e.Scan([]byte("valid")); err != nil {{
        t.Errorf("expected no error, got %v", err)
    }}
    if err := e.Scan(123); err == nil {{
        t.Errorf("expected error")
    }}
    
    // Test Scan on Null type
    if err := ne.Scan(nil); err != nil {{
        t.Errorf("expected no error, got %v", err)
    }}
    if err := ne.Scan(string("valid")); err != nil {{
        t.Errorf("expected no error, got %v", err)
    }}
    
    // Test Value on Null type
    ne.Valid = false
    v1, _ := ne.Value()
    if v1 != nil {{
        t.Errorf("expected nil")
    }}
    
    ne.Valid = true
    ne.{enum} = "valid"
    v2, _ := ne.Value()
    if v2 != string("valid") {{
        t.Errorf("expected valid, got %v", v2)
    }}
}}
"""
    with open('apps/finance-api/internal/sqlc/models_test.go', 'w') as f:
        f.write(test_code)

def generate_db_test():
    test_code = """package sqlc

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
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

	err = queries.WithTx(context.Background(), tx, func(q *Queries) error {
		return nil
	})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

    mock.ExpectBegin()
    tx2, _ := db.Begin()
    mock.ExpectRollback()
    err = queries.WithTx(context.Background(), tx2, func(q *Queries) error {
		return sql.ErrNoRows
	})
	if err == nil {
		t.Errorf("expected error")
	}
}
"""
    with open('apps/finance-api/internal/sqlc/db_test.go', 'w') as f:
        f.write(test_code)

def main():
    generate_models_test()
    generate_db_test()

if __name__ == "__main__":
    main()
