-- name: CreateUser :one
INSERT INTO users (full_name, date_of_birth, mobile_number, email, username, password_hash)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING user_id, username, email, full_name, created_at;

-- name: GetUserByUsername :one
SELECT user_id, username, email, full_name, mobile_number, password_hash, created_at, updated_at
FROM users
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT user_id, username, email, full_name, mobile_number, password_hash, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByMobileNumber :one
SELECT user_id, username, email, full_name, mobile_number, password_hash, created_at, updated_at
FROM users
WHERE mobile_number = $1;

-- name: GetUserByID :one
SELECT user_id, username, email, full_name, mobile_number, password_hash, created_at, updated_at
FROM users
WHERE user_id = $1;

-- name: CheckUsernameExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: CheckMobileNumberExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE mobile_number = $1);