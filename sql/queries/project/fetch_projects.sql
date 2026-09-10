-- name: FetchProjects :many
SELECT
    p.id,
    p.name,
    p.storage_id AS project_storage_id,
    s.id AS storage_id,
    s.collection AS storage_collection,
    s.embeddingModel AS storage_embedding_model,
    s.dimensions AS storage_dimensions,
    s.metric AS storage_metric,
    ds.id AS datasource_id,
    ds.provider AS datasource_provider,
    ds.credentials AS datasource_credentials,
    ds.data_type AS datasource_data_type
FROM projects p
LEFT JOIN storages s ON s.id = p.storage_id
LEFT JOIN projects_datasources pds ON pds.project_id = p.id
LEFT JOIN datasources ds ON ds.id = pds.datasource_id
ORDER BY p.id;
