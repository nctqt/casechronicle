-- name: AddMilestone :one
INSERT INTO milestones (
    id,
    case_id,
    title,
    description,
    event_date,
    created_at,
    updated_at
) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ListMilestonesByCase :many
SELECT * 
FROM milestones
WHERE case_id = $1
ORDER BY event_date ASC;

-- name: GetMilestoneByID :one
SELECT * 
FROM cases
WHERE id = $1 
LIMIT 1;

-- name: DeleteMilestone :exec
DELETE FROM milestones
WHERE id = $1;

-- name: UpdateTitle :exec
UPDATE milestones
SET title = $2,
    updated_at = $3
WHERE id = $1;

-- name: UpdateDescription :exec
UPDATE milestones
SET description = $2,
    updated_at = $3
WHERE id = $1;

-- name: UpdateEventDate :exec
UPDATE milestones
SET event_date = $2,
    updated_at = $3
WHERE id = $1;