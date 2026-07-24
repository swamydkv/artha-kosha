package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/config"
	internalhttp "artha-kosha/apps/finance-api/internal/http"
	"artha-kosha/apps/finance-api/internal/outbox"
)

func main() {
	cfg := config.Load()

	// If DATABASE_URL provided, use Postgres-backed sessions
	var provider *auth.LocalAuthProvider
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		p, pgRepo, err := auth.NewLocalAuthProviderFromDSN(dsn, 2*time.Hour)
		if err != nil {
			log.Fatalf("init provider from dsn: %v", err)
		}
		provider = p
		// start outbox worker using DB
		db := pgRepo.DB()
		w := outbox.NewWorker(db, 5*time.Second)
		w.Start(context.Background())
	} else {
		provider = auth.NewLocalAuthProvider()
	}

	mux := internalhttp.NewRouter(provider)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("starting auth API on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
