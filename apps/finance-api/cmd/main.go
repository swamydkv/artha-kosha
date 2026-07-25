package main

import (
	"context"
	"fmt"
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
	"artha-kosha/apps/finance-api/internal/sessions"
	"artha-kosha/apps/finance-api/internal/transactions"
	"artha-kosha/apps/finance-api/internal/users"
)

var (
	newAuthProvider = auth.NewLocalAuthProviderFromDSN
	startServer     = func(server *http.Server) error {
		// We check for a specific env var to avoid actually blocking during tests
		if os.Getenv("TEST_MODE") == "true" {
			return nil
		}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	}
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.Load()

	// If DATABASE_URL provided, use Postgres-backed sessions
	var provider *auth.LocalAuthProvider
	var auditRepo audit.Repository
	var mux http.Handler

	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		p, pgRepo, err := newAuthProvider(dsn, cfg.SessionTTL)
		if err != nil {
			return fmt.Errorf("init provider from dsn: %w", err)
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
		w := outbox.NewWorker(outboxRepo, nil, cfg.WorkerInterval)
		w.Start(context.Background())
		// start session cleanup worker
		sessionWorker := sessions.NewWorker(pgRepo, cfg.PruneInterval)
		sessionWorker.Start(context.Background())
		
		// start sweepers for outbox and archived users
		uRepo := users.NewSQLRepository(db)
		go func() {
			ticker := time.NewTicker(cfg.PruneInterval)
			for {
				select {
				case <-ticker.C:
					_ = outboxRepo.DeleteProcessed(context.Background(), cfg.OutboxRetentionDays)
					_ = uRepo.PruneArchivedUsers(context.Background())
				}
			}
		}()
		
		loginSvc := auth.NewLoginServiceWithDB(provider, nil, db)
		loginSvc.SetAuditService(auditSvc)
		loginSvc.SetDomainService(domainSvc)
		mux = internalhttp.NewRouter(loginSvc, auditRepo)
	} else {
		provider = auth.NewLocalAuthProvider()
		mux = internalhttp.NewRouter(provider, nil)
	}

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: cfg.RequestTimeout,
	}

	log.Printf("starting auth API on %s", server.Addr)
	return startServer(server)
}
