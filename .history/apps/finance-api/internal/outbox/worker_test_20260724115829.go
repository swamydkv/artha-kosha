package outbox

import (
    "context"
    "errors"
    "testing"
    "time"
)

type fakeRepo struct{
    entries map[string]OutboxEntry
    processed map[string]bool
    retries map[string]int
}

func newFakeRepo() *fakeRepo {
    return &fakeRepo{entries: make(map[string]OutboxEntry), processed: make(map[string]bool), retries: make(map[string]int)}
}

func (f *fakeRepo) Insert(ctx context.Context, e OutboxEntry) error { f.entries[e.ID]=e; return nil }
func (f *fakeRepo) FetchPending(ctx context.Context, limit int) ([]OutboxEntry, error) {
    var res []OutboxEntry
    for _, e := range f.entries {
        if e.ProcessingStatus=="pending" {
            res = append(res, e)
        }
    }
    return res, nil
}
func (f *fakeRepo) MarkProcessed(ctx context.Context, id string) error { f.processed[id]=true; e:=f.entries[id]; e.ProcessingStatus="processed"; f.entries[id]=e; return nil }
func (f *fakeRepo) IncrementRetry(ctx context.Context, id string) error { f.retries[id]++; return nil }
func (f *fakeRepo) MarkFailed(ctx context.Context, id string, reason string) error { e:=f.entries[id]; e.ProcessingStatus="failed"; f.entries[id]=e; return nil }

func TestWorker_DeliversAndRetries(t *testing.T) {
    repo := newFakeRepo()
    repo.Insert(context.Background(), OutboxEntry{ID:"e1", EventType:"EV", Payload:[]byte("p"), ProcessingStatus:"pending", CreatedAt: time.Now()})

    // processor fails first two times then succeeds
    attempts := 0
    proc := ProcessorFunc(func(ctx context.Context, e OutboxEntry) error {
        attempts++
        if attempts < 3 {
            return errors.New("transient")
        }
        return nil
    })

    w := NewWorker(repo, proc, 1*time.Hour)
    // run single processOnce
    if err := w.processOnce(context.Background()); err != nil {
        t.Fatalf("processOnce failed: %v", err)
    }
    if !repo.processed["e1"] {
        t.Fatalf("expected entry processed")
    }
    if repo.retries["e1"] < 2 {
        t.Fatalf("expected at least 2 retries, got %d", repo.retries["e1"])
    }
}
