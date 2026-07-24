package sessions

import (
	"context"
	"log"
	"time"
)

type Worker struct {
	repo     *PostgresRepo
	interval time.Duration
}

func NewWorker(repo *PostgresRepo, interval time.Duration) *Worker {
	return &Worker{
		repo:     repo,
		interval: interval,
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				err := w.repo.DeleteExpired(ctx, time.Now().UTC())
				if err != nil {
					log.Printf("error deleting expired sessions: %v", err)
				}
			}
		}
	}()
}
