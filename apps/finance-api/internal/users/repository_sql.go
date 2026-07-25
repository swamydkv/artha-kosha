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
		"user_id":   row.UserID.String(),
		"username":  row.Username,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
	deID := uuid.New()
	err = qtx.InsertDomainEvent(ctx, sqlc.InsertDomainEventParams{
		ID:               deID,
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

	err = qtx.InsertOutboxEntry(ctx, sqlc.InsertOutboxEntryParams{
		ID:               uuid.New(),
		DomainEventID:    deID,
		EventType:        string(sqlc.DomainEventTypeUSERREGISTERED),
		Payload:          dePayload,
		CreatedAt:        time.Now().UTC(),
		ProcessingStatus: sqlc.ProcessingStatusPending,
		RetryCount:       0,
	})
	if err != nil {
		return nil, fmt.Errorf("create outbox entry: %w", err)
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

func (r *SQLRepository) DeleteUser(ctx context.Context, userID string, archiveRetentionDays int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}
	now := time.Now().UTC()
	retentionExpiresAt := now.AddDate(0, 0, archiveRetentionDays)

	var u User
	err = tx.QueryRowContext(ctx, "SELECT username, email, full_name, mobile_number, password_hash FROM users WHERE user_id = $1 AND is_archived = false", uid).Scan(&u.Username, &u.Email, &u.FullName, &u.MobileNumber, &u.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found or already archived")
		}
		return fmt.Errorf("fetch user: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO archived_users (original_user_id, username, email, name, mobile, password_hash, retention_expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, uid, u.Username, u.Email, u.FullName, u.MobileNumber, u.PasswordHash, retentionExpiresAt)
	if err != nil {
		return fmt.Errorf("insert archived user: %w", err)
	}

	dePayload, _ := json.Marshal(map[string]interface{}{
		"user_id":              uid.String(),
		"username":             u.Username,
		"email":                u.Email,
		"deleted_at":           now.Format(time.RFC3339),
		"retention_expires_at": retentionExpiresAt.Format(time.RFC3339),
	})
	deID := uuid.New()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO domain_events (id, event_type, aggregate_id, aggregate_type, event_data, timestamp, processing_status, retry_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, deID, sqlc.DomainEventTypeUSERDELETED, uid.String(), "user", dePayload, now, sqlc.ProcessingStatusPending, 0)
	if err != nil {
		return fmt.Errorf("create domain event: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO transactional_outbox (id, domain_event_id, event_type, payload, created_at, processing_status, retry_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, uuid.New(), deID, string(sqlc.DomainEventTypeUSERDELETED), dePayload, now, sqlc.ProcessingStatusPending, 0)
	if err != nil {
		return fmt.Errorf("create outbox entry: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO audit_events (id, request_id, user_id, session_id, resource, resource_id, action, result, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, uuid.New(), "system", uuid.NullUUID{UUID: uid, Valid: true}, uuid.NullUUID{}, "user", sql.NullString{String: uid.String(), Valid: true}, "USER_DELETED", sqlc.AuditResultSuccess, now)
	if err != nil {
		return fmt.Errorf("create audit event: %w", err)
	}

	scrambled := fmt.Sprintf("deleted_user_%s", uid.String()[:8])
	_, err = tx.ExecContext(ctx, `
		UPDATE users
		SET username = $2, email = $2, full_name = $2, mobile_number = $2, password_hash = '', is_archived = true, deleted_at = NOW()
		WHERE user_id = $1
	`, uid, scrambled)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *SQLRepository) PruneArchivedUsers(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM archived_users WHERE retention_expires_at < NOW()`)
	return err
}
