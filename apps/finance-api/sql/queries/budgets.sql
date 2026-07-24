-- name: InsertBudget :exec
INSERT INTO budgets (id, owner_id, amount, name) VALUES ($1,$2,$3,$4);
