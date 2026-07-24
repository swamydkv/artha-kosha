package outbox

import (
    "context"
    "database/sql"
    "log"
    "time"
)

// Worker polls pending outbox entries and marks them processed (placeholder processor)
type Worker struct {
    db *sql.DB
    interval time.Duration
    stop chan struct{}
}

func NewWorker(db *sql.DB, interval time.Duration) *Worker {
    return &Worker{db: db, interval: interval, stop: make(chan struct{})}
}

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
    rows, err := w.db.QueryContext(ctx, `SELECT id, payload FROM transactional_outbox WHERE processing_status = 'pending' ORDER BY created_at ASC LIMIT 10`)
    if err != nil {
        return err
    }
    defer rows.Close()
    var id string
    var payload []byte
    for rows.Next() {
        if err := rows.Scan(&id, &payload); err != nil {
            return err
        }
        // placeholder: pretend to process
        log.Printf("processing outbox id=%s payload_len=%d", id, len(payload))
        // mark processed
        if _, err := w.db.ExecContext(ctx, `UPDATE transactional_outbox SET processing_status = 'processed', processed_at = NOW() WHERE id = $1`, id); err != nil {
            return err
        }
    }
    return rows.Err()
}
