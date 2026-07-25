package sessions

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(dsn string) (*PostgresRepo, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// basic ping to validate
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return &PostgresRepo{db: db}, nil
}

func NewPostgresRepoFromDB(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// DB exposes the underlying *sql.DB for integration tasks (migration/worker).
func (r *PostgresRepo) DB() *sql.DB { return r.db }

func (r *PostgresRepo) Create(s Session) error {
	q := `INSERT INTO sessions (id, user_id, created_at, last_activity_at, expires_at, revoked_at, user_agent, ip_address, status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.db.ExecContext(context.Background(), q, s.ID, s.UserID, s.CreatedAt, s.LastActivityAt, s.ExpiresAt, nullableTime(s.RevokedAt), s.UserAgent, s.IPAddress, string(s.Status))
	return err
}

func (r *PostgresRepo) Get(id string) (Session, error) {
	q := `SELECT id, user_id, created_at, last_activity_at, expires_at, revoked_at, user_agent, ip_address, status FROM sessions WHERE id = $1`
	var s Session
	var revoked sql.NullTime
	var status string
	row := r.db.QueryRowContext(context.Background(), q, id)
	err := row.Scan(&s.ID, &s.UserID, &s.CreatedAt, &s.LastActivityAt, &s.ExpiresAt, &revoked, &s.UserAgent, &s.IPAddress, &status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Session{}, errors.New("not found")
		}
		return Session{}, err
	}
	if revoked.Valid {
		s.RevokedAt = revoked.Time
	}
	s.Status = Status(status)
	if !s.ExpiresAt.IsZero() && time.Now().After(s.ExpiresAt) {
		s.Status = StatusExpired
	}
	return s, nil
}

func (r *PostgresRepo) Revoke(id string) error {
	q := `UPDATE sessions SET status = $1, revoked_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(context.Background(), q, string(StatusRevoked), time.Now().UTC(), id)
	return err
}

func (r *PostgresRepo) RevokeAllByUser(userID string) error {
	q := `UPDATE sessions SET status = $1, revoked_at = $2 WHERE user_id = $3`
	_, err := r.db.ExecContext(context.Background(), q, string(StatusRevoked), time.Now().UTC(), userID)
	return err
}

func (r *PostgresRepo) UpdateActivity(id string, lastActivity time.Time) error {
	q := `UPDATE sessions SET last_activity_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(context.Background(), q, lastActivity, id)
	return err
}

func nullableTime(t time.Time) interface{} {
	if t.IsZero() {
		return nil
	}
	return t
}

func (r *PostgresRepo) DeleteExpired(ctx context.Context, before time.Time) error {
	q := `DELETE FROM sessions WHERE expires_at < $1`
	_, err := r.db.ExecContext(ctx, q, before)
	return err
}
