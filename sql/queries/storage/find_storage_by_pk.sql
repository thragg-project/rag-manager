-- name: FindStorageByPK :one
SELECT
    id,
    collection,
    embeddingModel,
    dimensions,
    metric
FROM storages
WHERE id = $1
LIMIT 1;
