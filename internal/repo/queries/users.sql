-- name: CreateUser :one
INSERT INTO app_user (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, username, email, password_hash, created_at;

-- name: GetUserByUsername :one
SELECT id, username, email, password_hash, created_at
FROM app_user
WHERE username = $1;
