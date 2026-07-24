-- SQL queries for transactional_outbox
-- name: InsertOutboxEntry :exec
INSERT INTO transactional_outbox (id, domain_event_id, event_type, payload, created_at, processing_status, retry_count)
VALUES ($1,$2,$3,$4,$5,$6,$7);

-- name: GetPendingOutboxEntries :many
SELECT id, domain_event_id, event_type, payload, created_at, processing_status, retry_count, processed_at, error_message
FROM transactional_outbox WHERE processing_status = 'pending' ORDER BY created_at ASC LIMIT $1;
