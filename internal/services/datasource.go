package services

import (
	"context"
	"fmt"

	"github.com/thragg-project/rag-manager/internal/core"
)

type (
	DataSourceRepository interface {
		Fetch(ctx context.Context) ([]core.DataSource, error)
		FindByPK(ctx context.Context, id int64) (core.DataSource, error)
		Insert(ctx context.Context, dataSource core.DataSource) (core.DataSource, error)
		Update(ctx context.Context, id int64, dataSource core.DataSource) (core.DataSource, error)
	}
	DataSourceService struct {
		dataSources DataSourceRepository
		events      EventPublisher
		tx          Transactor
	}
)

func NewDataSourceService(
	dataSources DataSourceRepository,
	events EventPublisher,
	tx Transactor,
) *DataSourceService {
	return &DataSourceService{
		dataSources: dataSources,
		events:      events,
		tx:          tx,
	}
}

func (s *DataSourceService) Create(
	ctx context.Context,
	ds core.DataSource,
) (core.DataSource, error) {
	var created core.DataSource
	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		created, err = s.dataSources.Insert(ctx, ds)
		if err != nil {
			return err
		}

		if created.Id == nil {
			return fmt.Errorf("datasource id is nil after insert")
		}

		return s.events.Publish(ctx, Event{
			AggregateType: core.DatasourceAggregateType,
			AggregateID:   *created.Id,
			EventType:     core.DatasourceEventCreated,
			Payload:       created,
		})
	})
	if err != nil {
		return core.DataSource{}, err
	}

	return created, nil
}

func (s *DataSourceService) Get(ctx context.Context, id int64) (core.DataSource, error) {
	return s.dataSources.FindByPK(ctx, id)
}

func (s *DataSourceService) Search(ctx context.Context) ([]core.DataSource, error) {
	return s.dataSources.Fetch(ctx)
}
