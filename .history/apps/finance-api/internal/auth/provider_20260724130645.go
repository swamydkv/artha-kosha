package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"artha-kosha/apps/finance-api/internal/accounts"
	"artha-kosha/apps/finance-api/internal/audit"
	"artha-kosha/apps/finance-api/internal/budgets"
	"artha-kosha/apps/finance-api/internal/domain"
	"artha-kosha/apps/finance-api/internal/sessions"
	"artha-kosha/apps/finance-api/internal/transactions"
)

type AuthProvider interface {
	Register(RegisterUserRequest) (RegisterUserResponse, error)
	Login(LoginRequest) (LoginResponse, error)
	Logout(string) error
	GetSession(string) (sessions.Session, error)
	RevokeAll(string) error
}

type RegisterUserRequest struct {
	FullName        string `json:"full_name"`
	DateOfBirth     string `json:"date_of_birth"`
	MobileNumber    string `json:"mobile_number"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type RegisterUserResponse struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	FirstName    string `json:"first_name"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID         string `json:"user_id"`
	SessionID      string `json:"session_id"`
	WelcomeMessage string `json:"welcome_message"`
}

type localUser struct {
	ID           string
	Username     string
	Email        string
	MobileNumber string
	FullName     string
	FirstName    string
	PasswordHash string
	CreatedAt    time.Time
}

type localSession struct {
	SessionID string
	UserID    string
	CreatedAt time.Time
}

type LocalAuthProvider struct {
	mu          sync.RWMutex
	users       map[string]*localUser
	sessSvc     *sessions.Service
	domainSvc   *domain.Service
	auditSvc    *audit.Service
	accountsSvc *accounts.Service
	txSvc       *transactions.Service
	budgetsSvc  *budgets.Service
}

func NewLocalAuthProvider() *LocalAuthProvider {
	repo := sessions.NewInMemoryRepo()
	svc := sessions.NewService(repo, 2*time.Hour)
	return &LocalAuthProvider{
		users:   make(map[string]*localUser),
		sessSvc: svc,
	}
}

// NewLocalAuthProviderWithRepo constructs a LocalAuthProvider using a custom sessions repository.
func NewLocalAuthProviderWithRepo(repo sessions.Repository, ttl time.Duration) *LocalAuthProvider {
	svc := sessions.NewService(repo, ttl)
	return &LocalAuthProvider{
		users:   make(map[string]*localUser),
		sessSvc: svc,
	}
}

// NewLocalAuthProviderFromDSN constructs a LocalAuthProvider backed by Postgres using the provided DSN.
// Returns the provider and the underlying *sql.DB for integration use (workers/migrations).
func NewLocalAuthProviderFromDSN(dsn string, ttl time.Duration) (*LocalAuthProvider, *sessions.PostgresRepo, error) {
	pg, err := sessions.NewPostgresRepo(dsn)
	if err != nil {
		return nil, nil, err
	}
	svc := sessions.NewService(pg, ttl)
	return &LocalAuthProvider{users: make(map[string]*localUser), sessSvc: svc}, pg, nil
}

// SetDomainService attaches a domain service to emit domain events for actions like register/login
func (p *LocalAuthProvider) SetDomainService(svc *domain.Service) { p.domainSvc = svc }

// SetAuditService attaches an audit service for recording audit events.
func (p *LocalAuthProvider) SetAuditService(svc *audit.Service) { p.auditSvc = svc }

// SetAccountsService attaches accounts service for routing.
func (p *LocalAuthProvider) SetAccountsService(svc *accounts.Service) { p.accountsSvc = svc }
func (p *LocalAuthProvider) GetAccountsService() *accounts.Service    { return p.accountsSvc }

// SetTransactionsService attaches transactions service for routing.
func (p *LocalAuthProvider) SetTransactionsService(svc *transactions.Service) { p.txSvc = svc }
func (p *LocalAuthProvider) GetTransactionsService() *transactions.Service    { return p.txSvc }

// SetBudgetsService attaches budgets service for routing.
func (p *LocalAuthProvider) SetBudgetsService(svc *budgets.Service) { p.budgetsSvc = svc }
func (p *LocalAuthProvider) GetBudgetsService() *budgets.Service    { return p.budgetsSvc }

func (p *LocalAuthProvider) Register(req RegisterUserRequest) (RegisterUserResponse, error) {
	if err := validateRegistrationRequest(req); err != nil {
		return RegisterUserResponse{}, err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.users[strings.ToLower(req.Username)]; exists {
		return RegisterUserResponse{}, errors.New("duplicate username")
	}
	if _, exists := p.users[strings.ToLower(req.Email)]; exists {
		return RegisterUserResponse{}, errors.New("duplicate email")
	}
	if _, exists := p.users[strings.ToLower(req.MobileNumber)]; exists {
		return RegisterUserResponse{}, errors.New("duplicate mobile number")
	}

	userID := generateID("user")
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return RegisterUserResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}
	firstName := firstNameFromFullName(req.FullName)
	user := &localUser{
		ID:           userID,
		Username:     strings.ToLower(req.Username),
		Email:        strings.ToLower(req.Email),
		MobileNumber: strings.TrimSpace(req.MobileNumber),
		FullName:     strings.TrimSpace(req.FullName),
		FirstName:    firstName,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
	}

	p.users[strings.ToLower(user.Username)] = user
	p.users[strings.ToLower(user.Email)] = user
	p.users[strings.ToLower(user.MobileNumber)] = user

	// emit domain event if service available (best-effort)
	if p.domainSvc != nil {
		payload := fmt.Sprintf(`{"user_id":"%s","username":"%s","first_name":"%s","timestamp":"%s"}`,
			user.ID, user.Username, user.FirstName, time.Now().UTC().Format(time.RFC3339))
		evt := domain.DomainEvent{
			ID:            generateID("de"),
			EventType:     "USER_REGISTERED",
			AggregateID:   user.ID,
			AggregateType: "user",
			EventData:     []byte(payload),
			Timestamp:     time.Now().UTC(),
		}
		_ = p.domainSvc.Emit(context.Background(), evt)
	}

	return RegisterUserResponse{
		UserID:       user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
	}, nil
}

func (p *LocalAuthProvider) Login(req LoginRequest) (LoginResponse, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	user, found := p.users[strings.ToLower(req.Username)]
	if !found {
		return LoginResponse{}, errors.New("invalid credentials")
	}
	match, err := passwordMatches(req.Password, user.PasswordHash)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("password comparison error: %w", err)
	}
	if !match {
		return LoginResponse{}, errors.New("invalid credentials")
	}

	sessionID := generateID("session")
	// create session via session service
	_, err = p.sessSvc.CreateSession(sessionID, user.ID, "", "")
	if err != nil {
		return LoginResponse{}, err
	}

	if p.domainSvc != nil {
		payload := fmt.Sprintf(`{"user_id":"%s","username":"%s","session_id":"%s","timestamp":"%s"}`,
			user.ID, user.Username, sessionID, time.Now().UTC().Format(time.RFC3339))
		evt := domain.DomainEvent{
			ID:            generateID("de"),
			EventType:     "USER_LOGGED_IN",
			AggregateID:   user.ID,
			AggregateType: "user",
			EventData:     []byte(payload),
			Timestamp:     time.Now().UTC(),
		}
		_ = p.domainSvc.Emit(context.Background(), evt)
	}

	return LoginResponse{
		UserID:         user.ID,
		SessionID:      sessionID,
		WelcomeMessage: fmt.Sprintf("Welcome, %s!", user.FirstName),
	}, nil
}

func (p *LocalAuthProvider) Logout(sessionID string) error {
	// delegate to session service
	return p.sessSvc.RevokeSession(sessionID)
}

// GetSession retrieves session details
func (p *LocalAuthProvider) GetSession(id string) (sessions.Session, error) {
	return p.sessSvc.GetSession(id)
}

// RevokeAll revokes all sessions for a user
func (p *LocalAuthProvider) RevokeAll(userID string) error {
	return p.sessSvc.RevokeAll(userID)
}

func validateRegistrationRequest(req RegisterUserRequest) error {
	if strings.TrimSpace(req.FullName) == "" ||
		strings.TrimSpace(req.DateOfBirth) == "" ||
		strings.TrimSpace(req.MobileNumber) == "" ||
		strings.TrimSpace(req.Email) == "" ||
		strings.TrimSpace(req.Username) == "" ||
		strings.TrimSpace(req.Password) == "" ||
		strings.TrimSpace(req.ConfirmPassword) == "" {
		return errors.New("all fields are required")
	}
	if req.Password != req.ConfirmPassword {
		return errors.New("password and confirm password must match")
	}
	if len(req.Password) < 12 {
		return errors.New("password must be at least 12 characters")
	}
	if !containsUpper(req.Password) || !containsLower(req.Password) || !containsDigit(req.Password) || !containsSpecial(req.Password) {
		return errors.New("password must meet complexity requirements")
	}
	if !regexp.MustCompile(`^[A-Za-z0-9_.]{4,30}$`).MatchString(req.Username) {
		return errors.New("username format is invalid")
	}
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return errors.New("date of birth must be a valid date")
	}
	if dob.After(time.Now()) {
		return errors.New("date of birth cannot be in the future")
	}
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		return errors.New("email format is invalid")
	}
	// Basic validation for format: local@domain.tld
	parts := strings.Split(req.Email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return errors.New("email format is invalid")
	}
	return nil
}

func firstNameFromFullName(name string) string {
	parts := strings.Fields(strings.TrimSpace(name))
	if len(parts) == 0 {
		return "User"
	}
	return parts[0]
}

func generateID(prefix string) string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(buf))
}

func containsUpper(value string) bool {
	for _, ch := range value {
		if ch >= 'A' && ch <= 'Z' {
			return true
		}
	}
	return false
}

func containsLower(value string) bool {
	for _, ch := range value {
		if ch >= 'a' && ch <= 'z' {
			return true
		}
	}
	return false
}

func containsDigit(value string) bool {
	for _, ch := range value {
		if ch >= '0' && ch <= '9' {
			return true
		}
	}
	return false
}

func containsSpecial(value string) bool {
	for _, ch := range value {
		switch ch {
		case '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '_', '+', '=', '{', '}', '[', ']', '|', ';', ':', '<', '>', ',', '.', '?', '/':
			return true
		}
	}
	return false
}
