package projector

import (
	"context"
	"errors"
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	)

// ErrNilAggregateStore is when a dispatcher is created with a nil aggregate store.
var ErrNilAggregateStore = errors.New("aggregate store is nil")

// projector.Manager control the projector.
// it implement the following use-cases
//
// 1. Restore an aggregate based on uuid
// 2. Search an aggregate, get uuid
type Manager struct {
	t     eh.AggregateType
	store eh.AggregateStore
}

// NewCommandHandler creates a new CommandHandler for an aggregate type.
func NewManager(t eh.AggregateType, store eh.AggregateStore) (*Manager, error) {
	if store == nil {
		return nil, ErrNilAggregateStore
	}

	h := &Manager{
		t:     t,
		store: store,
	}
	return h, nil
}

func (m *Manager) Restore(ctx context.Context, uuid uuid.UUID) (eh.Aggregate, error) {
	return m.store.Restore(ctx, m.t, uuid)
}