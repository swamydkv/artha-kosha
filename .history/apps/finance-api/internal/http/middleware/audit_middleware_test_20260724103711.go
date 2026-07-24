package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "artha-kosha/apps/finance-api/internal/audit"
    "context"
)

type fakeAuditRepo struct{
    ch chan audit.AuditEvent
}

func (f *fakeAuditRepo) Insert(ctx context.Context, e audit.AuditEvent) error {
    select {
    case f.ch <- e:
    default:
    }
    return nil
}

func TestAuditMiddleware_InsertsAuditEventAsync(t *testing.T) {
    ch := make(chan audit.AuditEvent, 1)
    repo := &fakeAuditRepo{ch: ch}

    next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    handler := AuditMiddleware(repo)(next)

    req := httptest.NewRequest("POST", "/resource/1", nil)
    req.Header.Set("X-Request-ID", "req-123")
    req.Header.Set("X-Session-ID", "sess-456")
    rr := httptest.NewRecorder()

    handler.ServeHTTP(rr, req)

    select {
    case ev := <-ch:
        if ev.RequestID != "req-123" {
            t.Fatalf("expected RequestID 'req-123' got '%s'", ev.RequestID)
        }
        if ev.SessionID != "sess-456" {
            t.Fatalf("expected SessionID 'sess-456' got '%s'", ev.SessionID)
        }
        if ev.Resource != "/resource/1" {
            t.Fatalf("expected Resource '/resource/1' got '%s'", ev.Resource)
        }
        if ev.Action != "POST" {
            t.Fatalf("expected Action 'POST' got '%s'", ev.Action)
        }
    case <-time.After(500 * time.Millisecond):
        t.Fatal("audit Insert was not called")
    }
}
