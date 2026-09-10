package heartbeat

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/thragg-project/rag-manager/pkg/probes"
)

func NewProbes(pool *pgxpool.Pool) *probes.Probes {
	appProbes := probes.NewProbes()
	appProbes.SetReadinessHeartbeat(NewReadinessHeartbeat(pool), time.Second*2, time.Second*3)
	return appProbes
}
