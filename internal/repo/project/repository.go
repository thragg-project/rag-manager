package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/thragg-project/rag-manager/internal/core"
	apperrors "github.com/thragg-project/rag-manager/internal/errors"
	"github.com/thragg-project/rag-manager/internal/repo"
	"github.com/thragg-project/rag-manager/internal/repo/project/sqlc"
	"github.com/thragg-project/rag-manager/pkg/database"
)

type ProjectRepository struct {
	db database.Queryer
}

func NewProjectRepository(db database.Queryer) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Fetch(ctx context.Context) ([]core.Project, error) {
	rows, err := r.queries(ctx).FetchProjects(ctx)
	if err != nil {
		return nil, err
	}

	projects := make(map[int64]*core.Project)
	result := make([]core.Project, 0)
	for _, row := range rows {
		project, ok := projects[row.ID]
		if !ok {
			mapped := mapFetchedProject(row)
			projects[row.ID] = &mapped
			project = &mapped
			result = append(result, mapped)
		}
		addFetchedProjectDataSource(project, row)
	}

	for i := range result {
		result[i] = *projects[*result[i].Id]
	}

	return result, nil
}

func (r *ProjectRepository) FindByPK(ctx context.Context, id int64) (core.Project, error) {
	rows, err := r.queries(ctx).FindProjectByPK(ctx, id)
	if err != nil {
		return core.Project{}, err
	}
	if len(rows) == 0 {
		return core.Project{}, apperrors.ErrNotFound
	}

	project := mapFoundProject(rows[0])
	for _, row := range rows {
		addFoundProjectDataSource(&project, row)
	}

	return project, nil
}

func (r *ProjectRepository) Insert(
	ctx context.Context,
	project core.Project,
) (core.Project, error) {
	storageID := storageID(project.Storage)

	row, err := r.queries(ctx).InsertProject(ctx, sqlc.InsertProjectParams{
		Name:      repo.NullString(project.Name),
		StorageID: storageID,
	})
	if err != nil {
		return core.Project{}, err
	}

	if err := r.linkProjectDataSources(ctx, row.ID, project.DataSources); err != nil {
		return core.Project{}, err
	}

	return r.FindByPK(ctx, row.ID)
}

func (r *ProjectRepository) Update(
	ctx context.Context,
	id int64,
	project core.Project,
) (core.Project, error) {
	row, err := r.queries(ctx).UpdateProject(ctx, sqlc.UpdateProjectParams{
		ID:        id,
		Name:      repo.NullString(project.Name),
		StorageID: storageID(project.Storage),
	})
	if err != nil {
		return core.Project{}, err
	}

	return r.FindByPK(ctx, row.ID)
}

func mapFetchedProject(row sqlc.FetchProjectsRow) core.Project {
	project := core.Project{
		Id:          ptr(row.ID),
		Name:        row.Name.String,
		DataSources: []core.DataSource{},
	}

	if row.StorageID.Valid {
		project.Storage = mapProjectStorage(
			row.StorageID,
			row.StorageCollection,
			row.StorageEmbeddingModel,
			row.StorageDimensions,
			row.StorageMetric,
		)
	}

	return project
}

func mapFoundProject(row sqlc.FindProjectByPKRow) core.Project {
	project := core.Project{
		Id:          ptr(row.ID),
		Name:        row.Name.String,
		DataSources: []core.DataSource{},
	}

	if row.StorageID.Valid {
		project.Storage = mapProjectStorage(
			row.StorageID,
			row.StorageCollection,
			row.StorageEmbeddingModel,
			row.StorageDimensions,
			row.StorageMetric,
		)
	}

	return project
}

func (r *ProjectRepository) linkProjectDataSources(
	ctx context.Context,
	projectID int64,
	dataSources []core.DataSource,
) error {
	for _, dataSource := range dataSources {
		if dataSource.Id == nil {
			return fmt.Errorf("%w: project datasource id is required", apperrors.ErrInvalidArgument)
		}

		if err := r.queries(ctx).LinkProjectDataSource(ctx, sqlc.LinkProjectDataSourceParams{
			ProjectID:    projectID,
			DatasourceID: *dataSource.Id,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (r *ProjectRepository) queries(ctx context.Context) *sqlc.Queries {
	return sqlc.New(database.DBTX(ctx, r.db))
}

func storageID(storage core.Storage) pgtype.Int8 {
	if storage.Id == nil {
		return pgtype.Int8{}
	}

	return pgtype.Int8{
		Int64: *storage.Id,
		Valid: true,
	}
}

func mapProjectStorage(
	id pgtype.Int8,
	collection, embeddingModel, dimensions, metric pgtype.Text,
) core.Storage {
	return core.Storage{
		Id:             ptr(id.Int64),
		Collection:     collection.String,
		EmbeddingModel: embeddingModel.String,
		Dimensions:     dimensions.String,
		Metric:         metric.String,
	}
}

func addFetchedProjectDataSource(project *core.Project, row sqlc.FetchProjectsRow) {
	if !row.DatasourceID.Valid {
		return
	}

	project.DataSources = append(project.DataSources, mapProjectDataSource(
		row.DatasourceID,
		row.DatasourceProvider,
		row.DatasourceCredentials,
		row.DatasourceDataType,
	))
}

func addFoundProjectDataSource(project *core.Project, row sqlc.FindProjectByPKRow) {
	if !row.DatasourceID.Valid {
		return
	}

	project.DataSources = append(project.DataSources, mapProjectDataSource(
		row.DatasourceID,
		row.DatasourceProvider,
		row.DatasourceCredentials,
		row.DatasourceDataType,
	))
}

func mapProjectDataSource(
	id pgtype.Int8,
	provider pgtype.Text,
	credentialsRaw []byte,
	dataType pgtype.Text,
) core.DataSource {
	credentials := map[string]string{}
	if len(credentialsRaw) > 0 {
		_ = json.Unmarshal(credentialsRaw, &credentials)
	}

	return core.DataSource{
		Id:          ptr(id.Int64),
		Provider:    core.DataSourceProvider(provider.String),
		Credentials: credentials,
		DataType:    core.DataSourceType(dataType.String),
	}
}

func ptr[T any](value T) *T {
	return &value
}
