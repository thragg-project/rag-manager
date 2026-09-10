-- name: FetchDataSources :many
SELECT
    id,
    provider,
    credentials,
    data_type
FROM datasources
ORDER BY id;
