package http

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http/httptest"
    "testing"

    "github.com/go-chi/chi/v5"

    "artha-kosha/apps/finance-api/internal/accounts"
    "artha-kosha/apps/finance-api/internal/budgets"
    "artha-kosha/apps/finance-api/internal/transactions"
)

type mockAccountsRepo struct{}
func (m *mockAccountsRepo) CreateAccount(ctx context.Context, req accounts.CreateAccountRequest) (*accounts.Account, error) {
    return &accounts.Account{AccountID: req.OwnerID+"-acct", OwnerID: req.OwnerID, Name: req.Name}, nil
}

type mockTransactionsRepo struct{}
func (m *mockTransactionsRepo) CreateTransaction(ctx context.Context, req transactions.CreateTransactionRequest) (*transactions.Transaction, error) {
    return &transactions.Transaction{TransactionID: req.AccountID+"-tx", AccountID: req.AccountID, Amount: req.Amount}, nil
}

type mockBudgetsRepo struct{}
func (m *mockBudgetsRepo) CreateBudget(ctx context.Context, req budgets.CreateBudgetRequest) (*budgets.Budget, error) {
    return &budgets.Budget{BudgetID: req.OwnerID+"-bgt", OwnerID: req.OwnerID, Amount: req.Amount}, nil
}

func TestRegisterAccountsHandler(t *testing.T) {
    r := chi.NewRouter()
    svc := accounts.NewService(&mockAccountsRepo{})
    RegisterAccountsHandlers(r, svc)

    body := map[string]string{"OwnerID":"owner1","Name":"My Account"}
    b, _ := json.Marshal(body)
    req := httptest.NewRequest("POST", "/accounts", bytes.NewReader(b))
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != 201 {
        t.Fatalf("expected 201 got %d body=%s", w.Code, w.Body.String())
    }
}

func TestRegisterTransactionsHandler(t *testing.T) {
    r := chi.NewRouter()
    svc := transactions.NewService(&mockTransactionsRepo{})
    RegisterTransactionsHandlers(r, svc)

    body := map[string]interface{}{"AccountID":"acc1","Amount":100.0,"Memo":"note"}
    b, _ := json.Marshal(body)
    req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(b))
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != 201 {
        t.Fatalf("expected 201 got %d body=%s", w.Code, w.Body.String())
    }
}

func TestRegisterBudgetsHandler(t *testing.T) {
    r := chi.NewRouter()
    svc := budgets.NewService(&mockBudgetsRepo{})
    RegisterBudgetsHandlers(r, svc)

    body := map[string]interface{}{"OwnerID":"owner1","Amount":500.0,"Name":"monthly"}
    b, _ := json.Marshal(body)
    req := httptest.NewRequest("POST", "/budgets", bytes.NewReader(b))
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != 201 {
        t.Fatalf("expected 201 got %d body=%s", w.Code, w.Body.String())
    }
}
