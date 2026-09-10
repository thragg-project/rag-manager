-- name: FetchStorages :many
SELECT
    id,
    collection,
    embeddingModel,
    dimensions,
    metric
FROM storages
ORDER BY id;
