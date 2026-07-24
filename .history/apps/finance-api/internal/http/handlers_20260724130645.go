package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"artha-kosha/apps/finance-api/internal/accounts"
	"artha-kosha/apps/finance-api/internal/budgets"
	"artha-kosha/apps/finance-api/internal/transactions"
)

func RegisterAccountsHandlers(r chi.Router, svc *accounts.Service) {
	if svc == nil {
		return
	}
	r.Post("/accounts", func(w http.ResponseWriter, r *http.Request) {
		var req accounts.CreateAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		acc, err := svc.CreateAccount(r.Context(), req)
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(acc)
	})
}

func RegisterTransactionsHandlers(r chi.Router, svc *transactions.Service) {
	if svc == nil {
		return
	}
	r.Post("/transactions", func(w http.ResponseWriter, r *http.Request) {
		var req transactions.CreateTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		tx, err := svc.CreateTransaction(r.Context(), req)
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(tx)
	})
}

func RegisterBudgetsHandlers(r chi.Router, svc *budgets.Service) {
	if svc == nil {
		return
	}
	r.Post("/budgets", func(w http.ResponseWriter, r *http.Request) {
		var req budgets.CreateBudgetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		b, err := svc.CreateBudget(r.Context(), req)
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(b)
	})
}
