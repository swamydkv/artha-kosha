-- name: InsertAccount :exec
INSERT INTO accounts (id, owner_id, name) VALUES ($1,$2,$3);
