package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"artha-kosha/apps/finance-api/internal/accounts"
	"artha-kosha/apps/finance-api/internal/audit"
	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/budgets"
	"artha-kosha/apps/finance-api/internal/config"
	"artha-kosha/apps/finance-api/internal/domain"
	internalhttp "artha-kosha/apps/finance-api/internal/http"
	"artha-kosha/apps/finance-api/internal/outbox"
	"artha-kosha/apps/finance-api/internal/transactions"
)

func main() {
	cfg := config.Load()

	// If DATABASE_URL provided, use Postgres-backed sessions
	var provider *auth.LocalAuthProvider
	var auditRepo audit.Repository
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

		// instantiate and wire business services
		accountsRepo := accounts.NewSQLRepository(db)
		accountsSvc := accounts.NewServiceWithDB(accountsRepo, db)
		accountsSvc.SetAuditService(audit.NewService(audit.NewSQLRepository(db)))
		accountsSvc.SetDomainService(domainSvc)

		transactionsRepo := transactions.NewSQLRepository(db)
		transactionsSvc := transactions.NewServiceWithDB(transactionsRepo, db)
		transactionsSvc.SetAuditService(audit.NewService(audit.NewSQLRepository(db)))
		transactionsSvc.SetDomainService(domainSvc)

		budgetsRepo := budgets.NewSQLRepository(db)
		budgetsSvc := budgets.NewServiceWithDB(budgetsRepo, db)
		budgetsSvc.SetAuditService(audit.NewService(audit.NewSQLRepository(db)))
		budgetsSvc.SetDomainService(domainSvc)

		// expose via provider so router can register handlers
		provider.SetAccountsService(accountsSvc)
		provider.SetTransactionsService(transactionsSvc)
		provider.SetBudgetsService(budgetsSvc)
		auditRepo = audit.NewSQLRepository(db)
		// create audit service and attach to provider
		auditSvc := audit.NewService(auditRepo)
		provider.SetAuditService(auditSvc)
		// start outbox worker (use SQL repo and default processor)
		w := outbox.NewWorker(outboxRepo, nil, 5*time.Second)
		w.Start(context.Background())
		// pass auditRepo to router below
	} else {
		provider = auth.NewLocalAuthProvider()
	}

	mux := internalhttp.NewRouter(provider, auditRepo)

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
