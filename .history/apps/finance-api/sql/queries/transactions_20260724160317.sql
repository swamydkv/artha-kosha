-- name: InsertTransaction :exec
INSERT INTO transactions (id, account_id, amount, memo) VALUES ($1,$2,$3,$4);
