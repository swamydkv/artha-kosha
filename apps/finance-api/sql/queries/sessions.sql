-- name: CreateSession :one
INSERT INTO sessions (user_id, expires_at)
VALUES ($1, $2)
RETURNING session_id, user_id, created_at, expires_at;

-- name: GetSessionByID :one
SELECT session_id, user_id, created_at, expires_at, revoked_at
FROM sessions
WHERE session_id = $1;

-- name: GetActiveSessionByUserID :one
SELECT session_id, user_id, created_at, expires_at, revoked_at
FROM sessions
WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > NOW()
ORDER BY created_at DESC
LIMIT 1;

-- name: RevokeSession :one
UPDATE sessions
SET revoked_at = NOW()
WHERE session_id = $1
RETURNING session_id, user_id, created_at, expires_at, revoked_at;

-- name: RevokeAllUserSessions :exec
UPDATE sessions
SET revoked_at = NOW()
WHERE user_id = $1 AND revoked_at IS NULL;