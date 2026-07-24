package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sqlc-dev/pqtype"

	"artha-kosha/apps/finance-api/internal/sqlc"
)

type SQLRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *SQLRepository) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)

	row, err := qtx.CreateUser(ctx, sqlc.CreateUserParams{
		FullName:     req.FullName,
		DateOfBirth:  req.DateOfBirth,
		MobileNumber: req.MobileNumber,
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: req.PasswordHash,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if strings.Contains(pqErr.Message, "users_username_key") {
				return nil, errors.New("username already exists")
			}
			if strings.Contains(pqErr.Message, "users_email_key") {
				return nil, errors.New("email already exists")
			}
			if strings.Contains(pqErr.Message, "users_mobile_number_key") {
				return nil, errors.New("mobile number already exists")
			}
			return nil, errors.New("duplicate user entry")
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Insert domain event
	dePayload, _ := json.Marshal(map[string]interface{}{
		"user_id":    row.UserID.String(),
		"username":   row.Username,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	})
	err = qtx.InsertDomainEvent(ctx, sqlc.InsertDomainEventParams{
		ID:               uuid.New(),
		EventType:        sqlc.DomainEventTypeUSERREGISTERED,
		AggregateID:      row.UserID.String(),
		AggregateType:    "user",
		EventData:        dePayload,
		Timestamp:        time.Now().UTC(),
		ProcessingStatus: sqlc.ProcessingStatusPending,
		RetryCount:       0,
	})
	if err != nil {
		return nil, fmt.Errorf("create domain event: %w", err)
	}

	// Insert audit event
	err = qtx.InsertAuditEvent(ctx, sqlc.InsertAuditEventParams{
		ID:         uuid.New(),
		RequestID:  "system", // Ideally from context
		UserID:     uuid.NullUUID{UUID: row.UserID, Valid: true},
		Action:     "USER_REGISTERED",
		Resource:   "user",
		ResourceID: sql.NullString{String: row.UserID.String(), Valid: true},
		Result:     sqlc.AuditResultSuccess,
		Timestamp:  time.Now().UTC(),
	})
	if err != nil {
		return nil, fmt.Errorf("create audit event: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &User{
		UserID:       row.UserID.String(),
		Username:     row.Username,
		Email:        row.Email,
		FullName:     row.FullName,
		MobileNumber: req.MobileNumber,
		PasswordHash: req.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.CreatedAt,
	}, nil
}

func (r *SQLRepository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return nil, errors.New("not implemented")
}

func (r *SQLRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	row, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return &User{
		UserID:       row.UserID.String(),
		Username:     row.Username,
		Email:        row.Email,
		FullName:     row.FullName,
		MobileNumber: row.MobileNumber,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}, nil
}

func (r *SQLRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return nil, errors.New("not implemented")
}

func (r *SQLRepository) GetUserByMobileNumber(ctx context.Context, mobileNumber string) (*User, error) {
	return nil, errors.New("not implemented")
}

func (r *SQLRepository) CheckUserExists(ctx context.Context, username, email, mobileNumber string) (*UserExistsCheck, error) {
	return nil, errors.New("not implemented")
}

func (r *SQLRepository) CreateSession(ctx context.Context, sessionID, userID, ipAddress, userAgent string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)

	sid, _ := uuid.Parse(sessionID)
	uid, _ := uuid.Parse(userID)
	now := time.Now().UTC()

	var ip pqtype.Inet
	if ipAddress != "" {
		_ = ip.Scan(ipAddress)
	}
	
	err = qtx.CreateSession(ctx, sqlc.CreateSessionParams{
		ID:             sid,
		UserID:         uid,
		CreatedAt:      now,
		LastActivityAt: now,
		ExpiresAt:      now.Add(24 * time.Hour),
		RevokedAt:      sql.NullTime{},
		UserAgent:      sql.NullString{String: userAgent, Valid: userAgent != ""},
		IpAddress:      ip,
		Status:         sqlc.SessionStatusActive,
	})
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	dePayload, _ := json.Marshal(map[string]interface{}{
		"user_id":    userID,
		"session_id": sessionID,
		"timestamp":  now.Format(time.RFC3339),
	})
	err = qtx.InsertDomainEvent(ctx, sqlc.InsertDomainEventParams{
		ID:               uuid.New(),
		EventType:        sqlc.DomainEventTypeUSERLOGGEDIN,
		AggregateID:      userID,
		AggregateType:    "user",
		EventData:        dePayload,
		Timestamp:        now,
		ProcessingStatus: sqlc.ProcessingStatusPending,
		RetryCount:       0,
	})
	if err != nil {
		return fmt.Errorf("create domain event: %w", err)
	}

	err = qtx.InsertAuditEvent(ctx, sqlc.InsertAuditEventParams{
		ID:         uuid.New(),
		RequestID:  "system",
		UserID:     uuid.NullUUID{UUID: uid, Valid: true},
		SessionID:  uuid.NullUUID{UUID: sid, Valid: true},
		Action:     "USER_LOGGED_IN",
		Resource:   "session",
		ResourceID: sql.NullString{String: sessionID, Valid: true},
		Result:     sqlc.AuditResultSuccess,
		Timestamp:  now,
	})
	if err != nil {
		return fmt.Errorf("create audit event: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
