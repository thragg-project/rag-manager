-- name: InsertProject :one
INSERT INTO projects (
    name,
    storage_id
) VALUES (
    $1,
    $2
)
RETURNING
    id,
    name,
    storage_id;
