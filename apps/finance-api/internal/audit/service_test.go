package audit

import (
	"context"
	"database/sql"
	"testing"
)

type mockRepo struct{}

func (m *mockRepo) Insert(ctx context.Context, e AuditEvent) error { return nil }
func (m *mockRepo) InsertTx(ctx context.Context, tx *sql.Tx, e AuditEvent) error { return nil }

type mockRepoNoTx struct{}

func (m *mockRepoNoTx) Insert(ctx context.Context, e AuditEvent) error { return nil }

func TestRecordTx(t *testing.T) {
	// Test with TxInserter
	repo := &mockRepo{}
	svc := NewService(repo)

	err := svc.RecordTx(context.Background(), nil, AuditEvent{})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Test without TxInserter (fallback)
	repo2 := &mockRepoNoTx{}
	svc2 := NewService(repo2)

	err = svc2.RecordTx(context.Background(), nil, AuditEvent{})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
