-- name: InsertOutboxEntry :exec
INSERT INTO transactional_outbox (id, domain_event_id, event_type, payload, created_at, processing_status, retry_count)
VALUES (:id, :domain_event_id, :event_type, :payload, :created_at, :processing_status, :retry_count);

-- name: FetchPendingOutbox :many
SELECT id, domain_event_id, event_type, payload, created_at, processing_status, retry_count
FROM transactional_outbox
WHERE processing_status = 'pending'
ORDER BY created_at ASC
LIMIT $1;

-- name: MarkOutboxProcessed :exec
UPDATE transactional_outbox SET processing_status = 'processed', processed_at = NOW() WHERE id = :id;

-- name: IncrementOutboxRetry :exec
UPDATE transactional_outbox SET retry_count = COALESCE(retry_count,0)+1, last_retry_at = NOW() WHERE id = :id;

-- name: MarkOutboxFailed :exec
UPDATE transactional_outbox SET processing_status = 'failed', last_error = :reason, processed_at = NOW() WHERE id = :id;
