-- name: InsertStorage :one
INSERT INTO storages (
    collection,
    embeddingModel,
    dimensions,
    metric
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING
    id,
    collection,
    embeddingModel,
    dimensions,
    metric;
