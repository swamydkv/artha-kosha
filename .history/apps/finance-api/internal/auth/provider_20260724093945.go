package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"artha-kosha/apps/finance-api/internal/sessions"
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
	mu      sync.RWMutex
	users   map[string]*localUser
	sessSvc *sessions.Service
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
func NewLocalAuthProviderFromDSN(dsn string, ttl time.Duration) (*LocalAuthProvider, error) {
	pg, err := NewPostgresRepo(dsn)
	if err != nil {
		return nil, err
	}
	svc := sessions.NewService(pg, ttl)
	return &LocalAuthProvider{users: make(map[string]*localUser), sessSvc: svc}, nil
}

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
	passwordHash := hashPassword(req.Password)
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
	if !passwordMatches(req.Password, user.PasswordHash) {
		return LoginResponse{}, errors.New("invalid credentials")
	}

	sessionID := generateID("session")
	// create session via session service
	_, err := p.sessSvc.CreateSession(sessionID, user.ID, "", "")
	if err != nil {
		return LoginResponse{}, err
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
