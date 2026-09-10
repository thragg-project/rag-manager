package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/thragg-project/rag-manager/internal/core"
	apperrors "github.com/thragg-project/rag-manager/internal/errors"
	"github.com/thragg-project/rag-manager/internal/repo"
	"github.com/thragg-project/rag-manager/internal/repo/storage/sqlc"
	"github.com/thragg-project/rag-manager/pkg/database"
)

type StorageRepository struct {
	db database.Queryer
}

func NewStorageRepository(db database.Queryer) *StorageRepository {
	return &StorageRepository{db: db}
}

func (r *StorageRepository) Fetch(ctx context.Context) ([]core.Storage, error) {
	rows, err := r.queries(ctx).FetchStorages(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]core.Storage, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapStorage(row))
	}

	return result, nil
}

func (r *StorageRepository) FindByPK(ctx context.Context, id int64) (core.Storage, error) {
	row, err := r.queries(ctx).FindStorageByPK(ctx, id)
	if err != nil {
		return core.Storage{}, mapNoRows(err)
	}

	return mapStorage(row), nil
}

func (r *StorageRepository) Insert(
	ctx context.Context,
	storage core.Storage,
) (core.Storage, error) {
	row, err := r.queries(ctx).InsertStorage(ctx, sqlc.InsertStorageParams{
		Collection:     repo.NullString(storage.Collection),
		Embeddingmodel: repo.NullString(storage.EmbeddingModel),
		Dimensions:     repo.NullString(storage.Dimensions),
		Metric:         repo.NullString(storage.Metric),
	})
	if err != nil {
		return core.Storage{}, err
	}

	return mapStorage(row), nil
}

func (r *StorageRepository) Update(
	ctx context.Context,
	id int64,
	storage core.Storage,
) (core.Storage, error) {
	row, err := r.queries(ctx).UpdateStorage(ctx, sqlc.UpdateStorageParams{
		ID:             id,
		Collection:     repo.NullString(storage.Collection),
		Embeddingmodel: repo.NullString(storage.EmbeddingModel),
		Dimensions:     repo.NullString(storage.Dimensions),
		Metric:         repo.NullString(storage.Metric),
	})
	if err != nil {
		return core.Storage{}, err
	}

	return mapStorage(row), nil
}

func (r *StorageRepository) queries(ctx context.Context) *sqlc.Queries {
	return sqlc.New(database.DBTX(ctx, r.db))
}

func mapNoRows(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.ErrNotFound
	}

	return err
}

func mapStorage(row sqlc.Storage) core.Storage {
	return core.Storage{
		Id:             &row.ID,
		Collection:     row.Collection.String,
		EmbeddingModel: row.Embeddingmodel.String,
		Dimensions:     row.Dimensions.String,
		Metric:         row.Metric.String,
	}
}
