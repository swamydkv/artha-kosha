package audit

import (
	"context"
	"testing"
)

type mockRepo struct {
	insertErr error
}

func (m *mockRepo) Insert(ctx context.Context, e AuditEvent) error {
	return m.insertErr
}
func (m *mockRepo) Update(ctx context.Context, e AuditEvent) error {
	return nil
}
func (m *mockRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func TestAuditService(t *testing.T) {
	repo := &mockRepo{}
	svc := NewService(repo)

	err := svc.Record(context.Background(), AuditEvent{})
	if err != nil {
		t.Error("expected no error")
	}
}
