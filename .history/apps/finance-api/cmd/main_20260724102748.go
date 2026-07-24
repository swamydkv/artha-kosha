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
		// start outbox worker using DB and wire domain + audit services
		db := pgRepo.DB()
		outboxRepo := outbox.NewSQLRepository(db)
		outboxSvc := outbox.NewService(outboxRepo)
		domainRepo := domain.NewSQLRepository(db)
		domainSvc := domain.NewService(domainRepo, outboxSvc)
		provider.SetDomainService(domainSvc)
		auditRepo := audit.NewSQLRepository(db)
		// start outbox worker
		w := outbox.NewWorker(db, 5*time.Second)
		w.Start(context.Background())
		// pass auditRepo to router below
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
