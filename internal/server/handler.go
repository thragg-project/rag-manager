package server

import (
	"context"
	"fmt"

	"github.com/thragg-project/rag-manager/internal/server/gen"
)

type (
	DataSourceHandler interface {
		SearchDataSources(
			ctx context.Context,
			request gen.SearchDataSourcesRequestObject,
		) (gen.SearchDataSourcesResponseObject, error)
		CreateDataSource(
			ctx context.Context,
			request gen.CreateDataSourceRequestObject,
		) (gen.CreateDataSourceResponseObject, error)
		GetDataSource(
			ctx context.Context,
			request gen.GetDataSourceRequestObject,
		) (gen.GetDataSourceResponseObject, error)
	}
	ProjectHandler interface {
		SearchProjects(
			ctx context.Context,
			request gen.SearchProjectsRequestObject,
		) (gen.SearchProjectsResponseObject, error)
		CreateProject(
			ctx context.Context,
			request gen.CreateProjectRequestObject,
		) (gen.CreateProjectResponseObject, error)
		GetProject(
			ctx context.Context,
			request gen.GetProjectRequestObject,
		) (gen.GetProjectResponseObject, error)
	}
	StorageHandler interface {
		SearchStorages(
			ctx context.Context,
			request gen.SearchStoragesRequestObject,
		) (gen.SearchStoragesResponseObject, error)
		CreateStorage(
			ctx context.Context,
			request gen.CreateStorageRequestObject,
		) (gen.CreateStorageResponseObject, error)
		GetStorage(
			ctx context.Context,
			request gen.GetStorageRequestObject,
		) (gen.GetStorageResponseObject, error)
	}
	APIHandler struct {
		dataSourceHandler DataSourceHandler
		projectHandler    ProjectHandler
		storageHandler    StorageHandler
	}
)

var _ gen.StrictServerInterface = (*APIHandler)(nil)

func NewAPIHandler(
	dataSourceHandler DataSourceHandler,
	projectHandler ProjectHandler,
	storageHandler StorageHandler,
) *APIHandler {
	return &APIHandler{
		dataSourceHandler: dataSourceHandler,
		projectHandler:    projectHandler,
		storageHandler:    storageHandler,
	}
}

func (h *APIHandler) SearchDataSources(
	ctx context.Context,
	request gen.SearchDataSourcesRequestObject,
) (gen.SearchDataSourcesResponseObject, error) {
	if h.dataSourceHandler == nil {
		return nil, notImplemented("datasource handler not implemented")
	}
	return h.dataSourceHandler.SearchDataSources(ctx, request)
}

func (h *APIHandler) CreateDataSource(
	ctx context.Context,
	request gen.CreateDataSourceRequestObject,
) (gen.CreateDataSourceResponseObject, error) {
	if h.dataSourceHandler == nil {
		return nil, notImplemented("datasource handler not implemented")
	}
	return h.dataSourceHandler.CreateDataSource(ctx, request)
}

func (h *APIHandler) GetDataSource(
	ctx context.Context,
	request gen.GetDataSourceRequestObject,
) (gen.GetDataSourceResponseObject, error) {
	if h.dataSourceHandler == nil {
		return nil, notImplemented("datasource handler not implemented")
	}
	return h.dataSourceHandler.GetDataSource(ctx, request)
}

func (h *APIHandler) SearchProjects(
	ctx context.Context,
	request gen.SearchProjectsRequestObject,
) (gen.SearchProjectsResponseObject, error) {
	if h.projectHandler == nil {
		return nil, notImplemented("project handler not implemented")
	}
	return h.projectHandler.SearchProjects(ctx, request)
}

func (h *APIHandler) CreateProject(
	ctx context.Context,
	request gen.CreateProjectRequestObject,
) (gen.CreateProjectResponseObject, error) {
	if h.projectHandler == nil {
		return nil, notImplemented("project handler not implemented")
	}
	return h.projectHandler.CreateProject(ctx, request)
}

func (h *APIHandler) GetProject(
	ctx context.Context,
	request gen.GetProjectRequestObject,
) (gen.GetProjectResponseObject, error) {
	if h.projectHandler == nil {
		return nil, notImplemented("project handler not implemented")
	}
	return h.projectHandler.GetProject(ctx, request)
}

func (h *APIHandler) SearchStorages(
	ctx context.Context,
	request gen.SearchStoragesRequestObject,
) (gen.SearchStoragesResponseObject, error) {
	if h.storageHandler == nil {
		return nil, notImplemented("storage handler not implemented")
	}
	return h.storageHandler.SearchStorages(ctx, request)
}

func (h *APIHandler) CreateStorage(
	ctx context.Context,
	request gen.CreateStorageRequestObject,
) (gen.CreateStorageResponseObject, error) {
	if h.storageHandler == nil {
		return nil, notImplemented("storage handler not implemented")
	}
	return h.storageHandler.CreateStorage(ctx, request)
}

func (h *APIHandler) GetStorage(
	ctx context.Context,
	request gen.GetStorageRequestObject,
) (gen.GetStorageResponseObject, error) {
	if h.storageHandler == nil {
		return nil, notImplemented("storage handler not implemented")
	}
	return h.storageHandler.GetStorage(ctx, request)
}

func notImplemented(name string) error {
	return fmt.Errorf("%s is not configured", name)
}
