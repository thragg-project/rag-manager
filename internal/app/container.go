package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/thragg-project/rag-manager/internal/heartbeat"
	"github.com/thragg-project/rag-manager/internal/outbox"
	datasourcerepo "github.com/thragg-project/rag-manager/internal/repo/datasource"
	projectrepo "github.com/thragg-project/rag-manager/internal/repo/project"
	storagerepo "github.com/thragg-project/rag-manager/internal/repo/storage"
	"github.com/thragg-project/rag-manager/internal/services"
	"github.com/thragg-project/rag-manager/pkg/database"
	"github.com/thragg-project/rag-manager/pkg/probes"
)

type (
	Container struct {
		// DB
		DBPool *pgxpool.Pool

		// Probes
		Probes *probes.Probes

		// Services
		DataSourceService *services.DataSourceService
		ProjectService    *services.ProjectService
		StorageService    *services.StorageService
	}
)

var _ services.Transactor = (*database.Transactor)(nil)

func NewContainer(ctx context.Context, cfg Config) (Container, error) {
	// Setup Database Connection
	pool, err := database.GetConnectionPool(ctx, cfg.Database.DSN())
	if err != nil {
		return Container{}, err
	}

	// Setup Probes
	appProbes := heartbeat.NewProbes(pool)

	// Setup Database Adapters
	dataSourceRepository := datasourcerepo.NewDataSourceRepository(pool)
	projectRepository := projectrepo.NewProjectRepository(pool)
	storageRepository := storagerepo.NewStorageRepository(pool)
	eventPublisher := outbox.NewPublisher()
	transactor := database.NewTransactor(pool)

	return Container{
		DBPool: pool,
		Probes: appProbes,
		DataSourceService: services.NewDataSourceService(
			dataSourceRepository,
			eventPublisher,
			transactor,
		),
		ProjectService: services.NewProjectService(
			projectRepository,
			eventPublisher,
			transactor,
		),
		StorageService: services.NewStorageService(
			storageRepository,
			eventPublisher,
			transactor,
		),
	}, nil
}

func (c Container) Close() {
	if c.DBPool != nil {
		c.DBPool.Close()
	}
}
