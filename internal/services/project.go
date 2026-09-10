package services

import (
	"context"
	"fmt"

	"github.com/thragg-project/rag-manager/internal/core"
	apperrors "github.com/thragg-project/rag-manager/internal/errors"
)

type (
	ProjectRepository interface {
		Fetch(ctx context.Context) ([]core.Project, error)
		FindByPK(ctx context.Context, id int64) (core.Project, error)
		Insert(ctx context.Context, project core.Project) (core.Project, error)
		Update(ctx context.Context, id int64, project core.Project) (core.Project, error)
	}
	ProjectService struct {
		projects ProjectRepository
		events   EventPublisher
		tx       Transactor
	}
)

func NewProjectService(
	projects ProjectRepository,
	events EventPublisher,
	tx Transactor,
) *ProjectService {
	return &ProjectService{
		projects: projects,
		events:   events,
		tx:       tx,
	}
}

func (s *ProjectService) Create(ctx context.Context, project core.Project) (core.Project, error) {
	if err := validateProjectReferences(project); err != nil {
		return core.Project{}, err
	}

	var created core.Project
	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		created, err = s.projects.Insert(ctx, project)
		if err != nil {
			return err
		}

		if created.Id == nil {
			return fmt.Errorf("project id is nil after insert")
		}

		return s.events.Publish(ctx, Event{
			AggregateType: core.ProjectAggregateType,
			AggregateID:   *created.Id,
			EventType:     core.ProjectEventCreated,
			Payload:       created,
		})
	})
	if err != nil {
		return core.Project{}, err
	}

	return created, nil
}

func (s *ProjectService) Get(ctx context.Context, id int64) (core.Project, error) {
	return s.projects.FindByPK(ctx, id)
}

func (s *ProjectService) Search(ctx context.Context) ([]core.Project, error) {
	return s.projects.Fetch(ctx)
}

func validateProjectReferences(project core.Project) error {
	if project.Storage.Id == nil {
		return fmt.Errorf("%w: project storage id is required", apperrors.ErrInvalidArgument)
	}

	for _, dataSource := range project.DataSources {
		if dataSource.Id == nil {
			return fmt.Errorf("%w: project datasource id is required", apperrors.ErrInvalidArgument)
		}
	}

	return nil
}
