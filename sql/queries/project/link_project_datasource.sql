-- name: LinkProjectDataSource :exec
INSERT INTO projects_datasources (
    project_id,
    datasource_id
) VALUES (
    $1,
    $2
);
