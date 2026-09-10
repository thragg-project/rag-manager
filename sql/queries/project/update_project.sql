-- name: UpdateProject :one
UPDATE projects
SET
    name = $2,
    storage_id = $3
WHERE id = $1
RETURNING
    id,
    name,
    storage_id;
