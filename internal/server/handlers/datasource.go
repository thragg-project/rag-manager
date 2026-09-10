package handlers

import (
	"context"
	"fmt"

	"github.com/thragg-project/rag-manager/internal/core"
	apperrors "github.com/thragg-project/rag-manager/internal/errors"
	"github.com/thragg-project/rag-manager/internal/server/gen"
)

type (
	DataSourceService interface {
		Create(ctx context.Context, ds core.DataSource) (core.DataSource, error)
		Get(ctx context.Context, id int64) (core.DataSource, error)
		Search(ctx context.Context) ([]core.DataSource, error)
	}
	DataSourceHandler struct {
		service DataSourceService
	}
)

func NewDataSourceHandler(service DataSourceService) *DataSourceHandler {
	return &DataSourceHandler{service: service}
}

func (h *DataSourceHandler) SearchDataSources(
	ctx context.Context,
	request gen.SearchDataSourcesRequestObject,
) (gen.SearchDataSourcesResponseObject, error) {
	var dssObjs []gen.DataSource
	dss, err := h.service.Search(ctx)
	if err != nil {
		return nil, err
	}
	for _, ds := range dss {
		dssObjs = append(dssObjs, mapDataSource(ds))
	}
	return gen.SearchDataSources200JSONResponse(dssObjs), nil
}

func (h *DataSourceHandler) CreateDataSource(
	ctx context.Context,
	request gen.CreateDataSourceRequestObject,
) (gen.CreateDataSourceResponseObject, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("%w: request body is required", apperrors.ErrInvalidArgument)
	}
	newDS := core.DataSource{
		Id:          nil,
		DataType:    core.DataSourceType(request.Body.DataType),
		Provider:    core.DataSourceProvider(request.Body.Provider),
		Credentials: request.Body.Credentials,
	}
	ds, err := h.service.Create(ctx, newDS)
	if err != nil {
		return nil, err
	}
	if ds.Id == nil {
		return nil, fmt.Errorf("datasource id is nil")
	}
	return gen.CreateDataSource201JSONResponse(mapDataSource(ds)), nil
}

func (h *DataSourceHandler) GetDataSource(
	ctx context.Context,
	request gen.GetDataSourceRequestObject,
) (gen.GetDataSourceResponseObject, error) {
	if request.DataSourceId == 0 {
		return nil, fmt.Errorf("%w: datasource id is required", apperrors.ErrInvalidArgument)
	}
	ds, err := h.service.Get(ctx, request.DataSourceId)
	if err != nil {
		return nil, err
	}
	if ds.Id == nil {
		return nil, fmt.Errorf("datasource id is nil")
	}
	return gen.GetDataSource200JSONResponse(mapDataSource(ds)), nil
}

func mapDataSource(ds core.DataSource) gen.DataSource {
	return gen.DataSource{
		DataType:    string(ds.DataType),
		Id:          ds.Id,
		Provider:    string(ds.Provider),
		Credentials: ds.Credentials,
	}
}
