-- name: UpdateDataSource :one
UPDATE datasources
SET
    provider = $2,
    credentials = $3,
    data_type = $4
WHERE id = $1
RETURNING
    id,
    provider,
    credentials,
    data_type;
