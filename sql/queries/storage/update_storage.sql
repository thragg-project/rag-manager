-- name: UpdateStorage :one
UPDATE storages
SET
    collection = $2,
    embeddingModel = $3,
    dimensions = $4,
    metric = $5
WHERE id = $1
RETURNING
    id,
    collection,
    embeddingModel,
    dimensions,
    metric;
