-- name: AddCase :one
INSERT INTO cases (
    id, 
    title, 
    description, 
    created_at, 
    updated_at
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCaseByID :one
SELECT * 
FROM cases
WHERE id = $1 
LIMIT 1;

-- name: ListCases :many
SELECT * 
FROM cases
ORDER BY created_at DESC;

-- name: DeleteCase :exec
DELETE FROM cases
WHERE id = $1;

-- name: UpdateTitle :exec
UPDATE cases
SET title = $2,
    updated_at = $3
WHERE id = $1;

-- name: UpdateDescription :exec
UPDATE cases
SET description = $2,
    updated_at = $3
WHERE id = $1;