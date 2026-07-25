package http

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"artha-kosha/apps/finance-api/internal/accounts"
	"artha-kosha/apps/finance-api/internal/budgets"
	"artha-kosha/apps/finance-api/internal/constants"
	"artha-kosha/apps/finance-api/internal/transactions"
)

type mockAccRepo struct {
	err error
}

func (m mockAccRepo) CreateAccount(ctx context.Context, req accounts.CreateAccountRequest) (*accounts.Account, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &accounts.Account{AccountID: "acc-1", Name: req.Name}, nil
}
func (m mockAccRepo) GetAccountsByUserID(ctx context.Context, userID string) ([]accounts.Account, error) {
	return nil, nil
}

type mockTxRepo struct {
	err error
}

func (m mockTxRepo) CreateTransaction(ctx context.Context, req transactions.CreateTransactionRequest) (*transactions.Transaction, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &transactions.Transaction{TransactionID: "tx-1", Amount: req.Amount}, nil
}
func (m mockTxRepo) GetTransactions(ctx context.Context, userID, accountID string) ([]transactions.Transaction, error) {
	return nil, nil
}

type mockBdgRepo struct {
	err error
}

func (m mockBdgRepo) CreateBudget(ctx context.Context, req budgets.CreateBudgetRequest) (*budgets.Budget, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &budgets.Budget{BudgetID: "bdg-1"}, nil
}
func (m mockBdgRepo) GetBudgets(ctx context.Context, userID string) ([]budgets.Budget, error) {
	return nil, nil
}

func TestAccountsHandlers(t *testing.T) {
	tests := []struct {
		name         string
		repoErr      error
		body         string
		expectedCode int
	}{
		{"Invalid Body", nil, "{invalid", http.StatusBadRequest},
		{"Service Error", errors.New("db error"), `{"user_id":"u","name":"test"}`, http.StatusInternalServerError},
		{"Success", nil, `{"user_id":"u","name":"test","type":"SAVINGS","currency":"USD","initial_balance":0}`, http.StatusCreated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			svc := accounts.NewService(mockAccRepo{err: tt.repoErr})
			RegisterAccountsHandlers(r, svc)

			req := httptest.NewRequest("POST", "/accounts", bytes.NewBufferString(tt.body))
			req.Header.Set(constants.HeaderUserID, "usr-1")
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected %v, got %v", tt.expectedCode, rr.Code)
			}
		})
	}

	// Test nil service
	r := chi.NewRouter()
	RegisterAccountsHandlers(r, nil)
}

func TestBudgetsHandlers(t *testing.T) {
	tests := []struct {
		name         string
		repoErr      error
		body         string
		expectedCode int
	}{
		{"Invalid Body", nil, "{invalid", http.StatusBadRequest},
		{"Service Error", errors.New("db error"), `{"user_id":"u","category":"FOOD"}`, http.StatusInternalServerError},
		{"Success", nil, `{"user_id":"u","category":"FOOD","limit_amount":100,"period":"MONTHLY","start_date":"2023-01-01","end_date":"2023-01-31"}`, http.StatusCreated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			svc := budgets.NewService(mockBdgRepo{err: tt.repoErr})
			RegisterBudgetsHandlers(r, svc)

			req := httptest.NewRequest("POST", "/budgets", bytes.NewBufferString(tt.body))
			req.Header.Set(constants.HeaderUserID, "usr-1")
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected %v, got %v", tt.expectedCode, rr.Code)
			}
		})
	}

	// Test nil service
	r := chi.NewRouter()
	RegisterBudgetsHandlers(r, nil)
}

func TestTransactionsHandlers(t *testing.T) {
	tests := []struct {
		name         string
		repoErr      error
		body         string
		expectedCode int
	}{
		{"Invalid Body", nil, "{invalid", http.StatusBadRequest},
		{"Service Error", errors.New("db error"), `{"user_id":"u","amount":10}`, http.StatusInternalServerError},
		{"Success", nil, `{"user_id":"u","account_id":"a","amount":10,"type":"EXPENSE","category":"FOOD","date":"2023-01-01"}`, http.StatusCreated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			svc := transactions.NewService(mockTxRepo{err: tt.repoErr})
			RegisterTransactionsHandlers(r, svc)

			req := httptest.NewRequest("POST", "/transactions", bytes.NewBufferString(tt.body))
			req.Header.Set(constants.HeaderUserID, "usr-1")
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected %v, got %v", tt.expectedCode, rr.Code)
			}
		})
	}

	// Test nil service
	r := chi.NewRouter()
	RegisterTransactionsHandlers(r, nil)
}
