package handlers

import (
	"context"
	"fmt"

	"github.com/thragg-project/rag-manager/internal/core"
	apperrors "github.com/thragg-project/rag-manager/internal/errors"
	"github.com/thragg-project/rag-manager/internal/server/gen"
)

type (
	StorageService interface {
		Create(ctx context.Context, storage core.Storage) (core.Storage, error)
		Get(ctx context.Context, id int64) (core.Storage, error)
		Search(ctx context.Context) ([]core.Storage, error)
	}
	StorageHandler struct {
		service StorageService
	}
)

func NewStorageHandler(service StorageService) *StorageHandler {
	return &StorageHandler{service: service}
}

func (h *StorageHandler) SearchStorages(
	ctx context.Context,
	request gen.SearchStoragesRequestObject,
) (gen.SearchStoragesResponseObject, error) {
	storageObjs := make([]gen.Storage, 0)
	storages, err := h.service.Search(ctx)
	if err != nil {
		return nil, err
	}
	for _, storage := range storages {
		storageObjs = append(storageObjs, mapStorage(storage))
	}
	return gen.SearchStorages200JSONResponse(storageObjs), nil
}

func (h *StorageHandler) CreateStorage(
	ctx context.Context,
	request gen.CreateStorageRequestObject,
) (gen.CreateStorageResponseObject, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("%w: request body is required", apperrors.ErrInvalidArgument)
	}
	newStorage := core.Storage{
		Id:             nil,
		Collection:     request.Body.Collection,
		Dimensions:     request.Body.Dimensions,
		Metric:         request.Body.Metric,
		EmbeddingModel: request.Body.EmbeddingModel,
	}
	storage, err := h.service.Create(ctx, newStorage)
	if err != nil {
		return nil, err
	}
	if storage.Id == nil {
		return nil, fmt.Errorf("storage id is nil")
	}
	return gen.CreateStorage201JSONResponse(mapStorage(storage)), nil
}

func (h *StorageHandler) GetStorage(
	ctx context.Context,
	request gen.GetStorageRequestObject,
) (gen.GetStorageResponseObject, error) {
	if request.StorageId == 0 {
		return nil, fmt.Errorf("%w: storage id is required", apperrors.ErrInvalidArgument)
	}
	storage, err := h.service.Get(ctx, request.StorageId)
	if err != nil {
		return nil, err
	}
	if storage.Id == nil {
		return nil, fmt.Errorf("storage id is nil")
	}
	return gen.GetStorage200JSONResponse(mapStorage(storage)), nil
}

func mapStorage(storage core.Storage) gen.Storage {
	return gen.Storage{
		Collection:     storage.Collection,
		Dimensions:     storage.Dimensions,
		EmbeddingModel: storage.EmbeddingModel,
		Id:             storage.Id,
		Metric:         storage.Metric,
	}
}
