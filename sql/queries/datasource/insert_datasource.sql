-- name: InsertDataSource :one
INSERT INTO datasources (
    provider,
    credentials,
    data_type
) VALUES (
    $1,
    $2,
    $3
)
RETURNING
    id,
    provider,
    credentials,
    data_type;
