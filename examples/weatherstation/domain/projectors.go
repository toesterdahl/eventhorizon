// Copyright (c) 2014 - The Event Horizon authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package domain

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/eventhandler/projector"
)

// Temperature is a read model object for an weather station.
type Temperature struct {
	ID      uuid.UUID `bson:"_id"`
	Version int
	Name string
	Temperature    float32
	History	[]float32
}

var _ = eh.Entity(&Temperature{})
var _ = eh.Versionable(&Temperature{})

// EntityID implements the EntityID method of the eventhorizon.Entity interface.
func (i *Temperature) EntityID() uuid.UUID {
	return i.ID
}

// AggregateVersion implements the AggregateVersion method of the
// eventhorizon.Versionable interface.
func (i *Temperature) AggregateVersion() int {
	return i.Version
}

// TemperatureReadProjector is a projector that show a weather station and it's temperature history.
type TemperatureReadProjector struct{
}

// NewTemperatureReadProjector creates a new TemperatureReadProjector.
func NewTemperatureReadProjector() *TemperatureReadProjector {
	return &TemperatureReadProjector{
	}
}

// ProjectorType implements the ProjectorType method of the Projector interface.
func (p *TemperatureReadProjector) ProjectorType() projector.Type {
	return projector.Type("TemperatureReadProjector")
}

// Project implements the Project method of the Projector interface.
func (p *TemperatureReadProjector) Project(ctx context.Context, event eh.Event, entity eh.Entity) (eh.Entity, error) {
	i, ok := entity.(*Temperature)
	if !ok {
		return nil, errors.New("model is of incorrect type")
	}

	// Apply the changes for the event.
	switch event.EventType() {

	case WeatherStationCreatedEvent:
		log.Printf("TemperatureReadProjector:Project: Received WeatherStationCreatedEvent")
		data, ok := event.Data().(*WeatherStationCreatedData)
		if !ok {
			return nil, fmt.Errorf("projector: invalid event data type: %v", event.Data())
		}
		i.ID = event.AggregateID()
		i.Name = data.Name
		i.History =  make([]float32,0)

	case TemperatureReportedEvent:
		log.Printf("TemperatureReadProjector:Project: Received TemperatureReportedEvent")
		data, ok := event.Data().(*TemperatureReportedData)
		if !ok {
			return nil, fmt.Errorf("projector: invalid event data type: %v", event.Data())
		}
		i.ID = event.AggregateID()
		i.Temperature = data.Temperature
		i.History = append(i.History, data.Temperature)

	default:
		return nil, errors.New("Could not handle event: " + event.String())
	}

	i.Version++
	return i, nil
}
