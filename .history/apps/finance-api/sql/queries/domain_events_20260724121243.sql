-- name: InsertDomainEvent :exec
INSERT INTO domain_events (id, event_type, aggregate_id, aggregate_type, event_data, timestamp, processing_status, retry_count)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: FetchPendingDomainEvents :many
SELECT id, event_type, aggregate_id, aggregate_type, event_data, timestamp, processing_status, retry_count
FROM domain_events
WHERE processing_status = 'pending'
ORDER BY timestamp ASC
LIMIT $1;

-- name: MarkDomainEventProcessed :exec
UPDATE domain_events SET processing_status = 'processed', processed_at = NOW() WHERE id = $1;
