-- name: CreateUser :one
INSERT INTO users (full_name, role, login, password_hash, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, full_name, role, login, password_hash, status, created_at;

-- name: GetUserByLogin :one
SELECT id, full_name, role, login, password_hash, status, created_at
FROM users
WHERE login = $1;

-- name: GetUserByID :one
SELECT id, full_name, role, login, password_hash, status, created_at
FROM users
WHERE id = $1;
