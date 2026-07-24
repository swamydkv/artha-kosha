-- name: InsertAuditEvent :exec
INSERT INTO audit_events (id, request_id, user_id, session_id, resource, resource_id, action, result, timestamp, user_agent, client_ip)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: GetAuditEventsByRequestID :many
SELECT id, request_id, user_id, session_id, resource, resource_id, action, result, timestamp, user_agent, client_ip
FROM audit_events
WHERE request_id = $1
ORDER BY timestamp DESC;
