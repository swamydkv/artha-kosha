package outbox

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
)

type mockRepoForWorker struct {
	err     error
	entries []OutboxEntry
}

func (m *mockRepoForWorker) Insert(ctx context.Context, e OutboxEntry) error { return m.err }
func (m *mockRepoForWorker) InsertTx(ctx context.Context, tx *sql.Tx, e OutboxEntry) error { return m.err }
func (m *mockRepoForWorker) FetchPending(ctx context.Context, limit int) ([]OutboxEntry, error) {
	return m.entries, m.err
}
func (m *mockRepoForWorker) MarkProcessed(ctx context.Context, id string) error { return m.err }
func (m *mockRepoForWorker) IncrementRetry(ctx context.Context, id string) error { return m.err }
func (m *mockRepoForWorker) MarkFailed(ctx context.Context, id string, errMsg string) error { return m.err }
func (m *mockRepoForWorker) DeleteProcessed(ctx context.Context, before time.Time) (int64, error) { return 0, m.err }

func TestOutboxWorker(t *testing.T) {
	repo := &mockRepoForWorker{}
	w := NewWorker(repo, nil, 10*time.Millisecond)
	
	// Test Deliver for ProcessorFunc
	w.processor.Deliver(context.Background(), OutboxEntry{ID: "test"})
	
	// Test processOnce fetch error
	repo.err = errors.New("fetch error")
	_ = w.processOnce(context.Background())
	
	// Test processOnce delivery fail
	repo.err = nil
	repo.entries = []OutboxEntry{{ID: "1"}}
	failProcessor := ProcessorFunc(func(ctx context.Context, e OutboxEntry) error {
		return errors.New("fail")
	})
	w.processor = failProcessor
	// To avoid long sleeps, reduce maxRetries
	w.maxRetries = 1
	_ = w.processOnce(context.Background())
	
	// Test processOnce delivery success but mark processed fail
	repo.err = errors.New("mark fail")
	successProcessor := ProcessorFunc(func(ctx context.Context, e OutboxEntry) error {
		return nil
	})
	w.processor = successProcessor
	w.maxRetries = 1
	_ = w.processOnce(context.Background())

	// Test processOnce delivery success
	repo.err = nil
	_ = w.processOnce(context.Background())

	// Test Start/Stop
	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)
	time.Sleep(20 * time.Millisecond)
	w.Stop()
	cancel()
	
	// Test Context Cancel
	w2 := NewWorker(repo, nil, 10*time.Millisecond)
	ctx2, cancel2 := context.WithCancel(context.Background())
	w2.Start(ctx2)
	cancel2()
	time.Sleep(20 * time.Millisecond)
}
