-- name: AddUser :one
INSERT INTO users (id, email, hashed_password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * 
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * 
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT * 
FROM users
ORDER BY created_at DESC;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdatePassword :exec
UPDATE users
SET hashed_password = $2,
    updated_at = $3
WHERE id = $1;
