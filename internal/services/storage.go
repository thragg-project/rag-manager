package services

import (
	"context"
	"fmt"

	"github.com/thragg-project/rag-manager/internal/core"
)

type (
	StorageRepository interface {
		Fetch(ctx context.Context) ([]core.Storage, error)
		FindByPK(ctx context.Context, id int64) (core.Storage, error)
		Insert(ctx context.Context, storage core.Storage) (core.Storage, error)
		Update(ctx context.Context, id int64, storage core.Storage) (core.Storage, error)
	}
	StorageService struct {
		storages StorageRepository
		events   EventPublisher
		tx       Transactor
	}
)

func NewStorageService(
	storages StorageRepository,
	events EventPublisher,
	tx Transactor,
) *StorageService {
	return &StorageService{
		storages: storages,
		events:   events,
		tx:       tx,
	}
}

func (s *StorageService) Create(ctx context.Context, storage core.Storage) (core.Storage, error) {
	var created core.Storage
	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		created, err = s.storages.Insert(ctx, storage)
		if err != nil {
			return err
		}

		if created.Id == nil {
			return fmt.Errorf("project id is nil after insert")
		}

		return s.events.Publish(ctx, Event{
			AggregateType: core.StorageAggregateType,
			AggregateID:   *created.Id,
			EventType:     core.StorageEventCreated,
			Payload:       created,
		})
	})
	if err != nil {
		return core.Storage{}, err
	}

	return created, nil
}

func (s *StorageService) Get(ctx context.Context, id int64) (core.Storage, error) {
	return s.storages.FindByPK(ctx, id)
}

func (s *StorageService) Search(ctx context.Context) ([]core.Storage, error) {
	return s.storages.Fetch(ctx)
}
