-- name: InsertOutboxEntry :exec
INSERT INTO transactional_outbox (id, domain_event_id, event_type, payload, created_at, processing_status, retry_count)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: FetchPendingOutbox :many
SELECT id, domain_event_id, event_type, payload, created_at, processing_status, retry_count
FROM transactional_outbox
WHERE processing_status = 'pending'
ORDER BY created_at ASC
LIMIT $1;

-- name: MarkOutboxProcessed :exec
UPDATE transactional_outbox SET processing_status = 'processed', processed_at = NOW() WHERE id = $1;

-- name: IncrementOutboxRetry :exec
UPDATE transactional_outbox SET retry_count = COALESCE(retry_count,0)+1, last_retry_at = NOW() WHERE id = $1;

-- name: MarkOutboxFailed :exec
UPDATE transactional_outbox SET processing_status = 'failed', last_error = $2, processed_at = NOW() WHERE id = $1;
