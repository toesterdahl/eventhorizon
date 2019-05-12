package weatherstation

import (
	"context"
	"github.com/google/uuid"
	"github.com/looplab/eventhorizon/aggregatestore/events"
	"github.com/looplab/eventhorizon/eventhandler/projector"
	"github.com/looplab/eventhorizon/examples/weatherstation/domain"
	"log"
)

type ProjectorManager struct {
	projectorManager *projector.Manager
}

func SetupProjectorManager(aggregateStore *events.AggregateStore) *ProjectorManager {
	var err error
	var projectorManager *projector.Manager
	projectorManager, err = projector.NewManager(domain.WeatherStationAggregateType, aggregateStore)
	if err != nil {
		log.Fatalf("could not create projector manager: %s", err)
	}

	return &ProjectorManager {
		projectorManager,
	}
}

func (pm *ProjectorManager) Restore(ctx context.Context, id uuid.UUID) {
	aggregate, err := pm.projectorManager.Restore(ctx, id)
	if err != nil {
		log.Println("error:", err)
		return
	}
	log.Printf("restored: %#vs %#v", aggregate, aggregate.EntityID())
}
