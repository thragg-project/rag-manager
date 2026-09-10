package handlers

import (
	"context"
	"fmt"

	"github.com/thragg-project/rag-manager/internal/core"
	apperrors "github.com/thragg-project/rag-manager/internal/errors"
	"github.com/thragg-project/rag-manager/internal/server/gen"
)

type (
	ProjectService interface {
		Create(ctx context.Context, project core.Project) (core.Project, error)
		Get(ctx context.Context, id int64) (core.Project, error)
		Search(ctx context.Context) ([]core.Project, error)
	}
	ProjectHandler struct {
		service ProjectService
	}
)

func NewProjectHandler(service ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) SearchProjects(
	ctx context.Context,
	request gen.SearchProjectsRequestObject,
) (gen.SearchProjectsResponseObject, error) {
	prjsObjs := make([]gen.Project, 0)
	prjs, err := h.service.Search(ctx)
	if err != nil {
		return nil, err
	}
	for _, prj := range prjs {
		prjsObjs = append(prjsObjs, mapProject(prj))
	}
	return gen.SearchProjects200JSONResponse(prjsObjs), nil
}

func (h *ProjectHandler) CreateProject(
	ctx context.Context,
	request gen.CreateProjectRequestObject,
) (gen.CreateProjectResponseObject, error) {
	dataSources := make([]core.DataSource, 0)
	if request.Body == nil {
		return nil, fmt.Errorf("%w: request body is required", apperrors.ErrInvalidArgument)
	}
	for _, ds := range request.Body.DataSources {
		dataSources = append(dataSources, core.DataSource{
			Id:          ds.Id,
			DataType:    core.DataSourceType(ds.DataType),
			Provider:    core.DataSourceProvider(ds.Provider),
			Credentials: ds.Credentials,
		})
	}
	newPrj := core.Project{
		Id:          nil,
		Name:        request.Body.Name,
		DataSources: dataSources,
		Storage: core.Storage{
			Id:             request.Body.Storage.Id,
			Collection:     request.Body.Storage.Collection,
			Dimensions:     request.Body.Storage.Dimensions,
			EmbeddingModel: request.Body.Storage.EmbeddingModel,
			Metric:         request.Body.Storage.Metric,
		},
	}
	prj, err := h.service.Create(ctx, newPrj)
	if err != nil {
		return nil, err
	}
	if prj.Id == nil {
		return nil, fmt.Errorf("project id is nil")
	}
	return gen.CreateProject201JSONResponse(mapProject(prj)), nil
}

func (h *ProjectHandler) GetProject(
	ctx context.Context,
	request gen.GetProjectRequestObject,
) (gen.GetProjectResponseObject, error) {
	if request.ProjectId == 0 {
		return nil, fmt.Errorf("%w: project id is required", apperrors.ErrInvalidArgument)
	}
	prj, err := h.service.Get(ctx, request.ProjectId)
	if err != nil {
		return nil, err
	}
	if prj.Id == nil {
		return nil, fmt.Errorf("project id is nil")
	}
	return gen.GetProject200JSONResponse(mapProject(prj)), nil
}

func mapProject(prj core.Project) gen.Project {
	dss := make([]gen.DataSource, 0)
	for _, ds := range prj.DataSources {
		dss = append(dss, mapDataSource(ds))
	}
	return gen.Project{
		Id:          prj.Id,
		Name:        prj.Name,
		DataSources: dss,
		Storage:     mapStorage(prj.Storage),
	}
}
