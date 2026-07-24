-- name: CreateSession :exec
INSERT INTO sessions (id, user_id, created_at, last_activity_at, expires_at, revoked_at, user_agent, ip_address, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetSessionByID :one
SELECT id, user_id, created_at, last_activity_at, expires_at, revoked_at, user_agent, ip_address, status
FROM sessions
WHERE id = $1;

-- name: UpdateSessionActivity :exec
UPDATE sessions SET last_activity_at = $1 WHERE id = $2;

-- name: RevokeSession :exec
UPDATE sessions SET status = $1, revoked_at = $2 WHERE id = $3;

-- name: RevokeAllSessionsByUser :exec
UPDATE sessions SET status = $1, revoked_at = $2 WHERE user_id = $3;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expires_at < $1;