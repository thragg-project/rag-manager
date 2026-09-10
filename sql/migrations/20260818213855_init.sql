-- +goose Up
create table storages (
    id bigserial primary key,
    collection varchar,
    embeddingModel varchar,
    dimensions varchar,
    metric varchar
);

create table projects (
    id bigserial primary key,
    name varchar,
    storage_id bigint references storages (id)
);

create table datasources (
    id bigserial primary key,
    provider varchar,
    credentials json,
    data_type varchar
);

create table projects_datasources (
    id bigserial primary key,
    project_id bigint not null references projects (id) on delete cascade,
    datasource_id bigint not null references datasources (id) on delete cascade,
    unique (project_id, datasource_id)
);

-- +goose Down
drop table projects_datasources;
drop table datasources;
drop table projects;
drop table storages;
