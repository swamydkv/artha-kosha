package outbox

import (
	"context"
	"log"
	"math"
	"time"
)

// Processor delivers outbox entries to their destinations.
type Processor interface {
	Deliver(ctx context.Context, e OutboxEntry) error
}

// Worker polls pending outbox entries and delivers them using a Processor with retries.
type Worker struct {
	repo       ReadWriter
	processor  Processor
	interval   time.Duration
	maxRetries int
	stop       chan struct{}
}

func NewWorker(repo ReadWriter, processor Processor, interval time.Duration) *Worker {
	if processor == nil {
		// default processor just logs and marks processed via repo (no-op delivery)
		processor = ProcessorFunc(func(ctx context.Context, e OutboxEntry) error {
			log.Printf("default deliver: id=%s event=%s len=%d", e.ID, e.EventType, len(e.Payload))
			return nil
		})
	}
	return &Worker{repo: repo, processor: processor, interval: interval, maxRetries: 3, stop: make(chan struct{})}
}

// ProcessorFunc is a helper to implement Processor from a function.
type ProcessorFunc func(ctx context.Context, e OutboxEntry) error

func (f ProcessorFunc) Deliver(ctx context.Context, e OutboxEntry) error { return f(ctx, e) }

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := w.processOnce(ctx); err != nil {
					log.Printf("outbox worker: %v", err)
				}
			case <-w.stop:
				ticker.Stop()
				return
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (w *Worker) Stop() { close(w.stop) }

func (w *Worker) processOnce(ctx context.Context) error {
	entries, err := w.repo.FetchPending(ctx, 10)
	if err != nil {
		return err
	}
	for _, e := range entries {
		// attempt delivery with exponential backoff
		var attempt int
		var lastErr error
		for attempt = 0; attempt <= w.maxRetries; attempt++ {
			if attempt > 0 {
				backoff := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
				time.Sleep(backoff)
			}
			if err := w.processor.Deliver(ctx, e); err != nil {
				lastErr = err
				_ = w.repo.IncrementRetry(ctx, e.ID)
				continue
			}
			// success
			if err := w.repo.MarkProcessed(ctx, e.ID); err != nil {
				log.Printf("failed to mark processed id=%s: %v", e.ID, err)
			}
			lastErr = nil
			break
		}
		if lastErr != nil {
			// mark failed after retries
			_ = w.repo.MarkFailed(ctx, e.ID, lastErr.Error())
			log.Printf("outbox entry id=%s failed after %d attempts: %v", e.ID, attempt, lastErr)
		}
	}
	return nil
}
