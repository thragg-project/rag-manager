-- name: FindDataSourceByPK :one
SELECT
    id,
    provider,
    credentials,
    data_type
FROM datasources
WHERE id = $1
LIMIT 1;
