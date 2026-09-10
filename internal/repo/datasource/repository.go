package repo

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/thragg-project/rag-manager/internal/core"
	apperrors "github.com/thragg-project/rag-manager/internal/errors"
	"github.com/thragg-project/rag-manager/internal/repo"
	"github.com/thragg-project/rag-manager/internal/repo/datasource/sqlc"
	"github.com/thragg-project/rag-manager/pkg/database"
)

type DataSourceRepository struct {
	db database.Queryer
}

func NewDataSourceRepository(db database.Queryer) *DataSourceRepository {
	return &DataSourceRepository{db: db}
}

func (r *DataSourceRepository) Fetch(ctx context.Context) ([]core.DataSource, error) {
	rows, err := r.queries(ctx).FetchDataSources(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]core.DataSource, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapDataSource(row))
	}

	return result, nil
}

func (r *DataSourceRepository) FindByPK(ctx context.Context, id int64) (core.DataSource, error) {
	row, err := r.queries(ctx).FindDataSourceByPK(ctx, id)
	if err != nil {
		return core.DataSource{}, mapNoRows(err)
	}

	return mapDataSource(row), nil
}

func (r *DataSourceRepository) Insert(
	ctx context.Context,
	dataSource core.DataSource,
) (core.DataSource, error) {
	row, err := r.queries(ctx).InsertDataSource(ctx, sqlc.InsertDataSourceParams{
		Provider:    repo.NullString(string(dataSource.Provider)),
		Credentials: nullCredentials(dataSource.Credentials),
		DataType:    repo.NullString(string(dataSource.DataType)),
	})
	if err != nil {
		return core.DataSource{}, err
	}

	return mapDataSource(row), nil
}

func (r *DataSourceRepository) Update(
	ctx context.Context,
	id int64,
	dataSource core.DataSource,
) (core.DataSource, error) {
	row, err := r.queries(ctx).UpdateDataSource(ctx, sqlc.UpdateDataSourceParams{
		ID:          id,
		Provider:    repo.NullString(string(dataSource.Provider)),
		Credentials: nullCredentials(dataSource.Credentials),
		DataType:    repo.NullString(string(dataSource.DataType)),
	})
	if err != nil {
		return core.DataSource{}, err
	}

	return mapDataSource(row), nil
}

func (r *DataSourceRepository) queries(ctx context.Context) *sqlc.Queries {
	return sqlc.New(database.DBTX(ctx, r.db))
}

func mapNoRows(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.ErrNotFound
	}

	return err
}

func mapDataSource(row sqlc.Datasource) core.DataSource {
	credentials := map[string]string{}
	if len(row.Credentials) > 0 {
		_ = json.Unmarshal(row.Credentials, &credentials)
	}

	return core.DataSource{
		Id:          &row.ID,
		Provider:    core.DataSourceProvider(row.Provider.String),
		Credentials: credentials,
		DataType:    core.DataSourceType(row.DataType.String),
	}
}

func nullCredentials(value map[string]string) []byte {
	if len(value) == 0 {
		return nil
	}

	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}

	return raw
}
