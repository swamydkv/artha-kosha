package integration

import (
	"context"
	"errors"
	"testing"
	"time"

	"artha-kosha/apps/finance-api/internal/audit"
)

type failingAuditRepo struct{}

func (f *failingAuditRepo) Insert(ctx context.Context, e audit.AuditEvent) error {
	return errors.New("audit database unavailable")
}

func TestAuditEventStorageFailureHandling(t *testing.T) {
	repo := &failingAuditRepo{}
	svc := audit.NewService(repo)

	// In this system, audit failures should not crash the business logic,
	// but might just be logged, or if strict mode is enabled, it should return an error.

	evt := audit.AuditEvent{
		ID:        "evt-test",
		Timestamp: time.Now(),
		Action:    "USER_LOGIN",
		UserID:    "user-123",
	}

	err := svc.Record(context.Background(), evt)
	if err == nil {
		t.Error("Expected error when audit storage fails, but got nil")
	}
}
