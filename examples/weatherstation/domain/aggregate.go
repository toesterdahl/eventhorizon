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
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
)

func init() {
	eh.RegisterAggregate(func(id uuid.UUID) eh.Aggregate {
		return NewWeatherStationAggregate(id)
	})
}

// WeatherStationAggregateType is the type Name of the aggregate.
const WeatherStationAggregateType eh.AggregateType = "Temperature"

// WeatherStationAggregate is the root aggregate.
//
type WeatherStationAggregate struct {
	// AggregateBase implements most of the eventhorizon.Aggregate interface.
	*events.AggregateBase

	Name        string
	Temperature float32
}

var _ = eh.Aggregate(&WeatherStationAggregate{})

// NewWeatherStationAggregate creates a new WeatherStationAggregate with an ID.
func NewWeatherStationAggregate(id uuid.UUID) *WeatherStationAggregate {
	return &WeatherStationAggregate{
		AggregateBase: events.NewAggregateBase(WeatherStationAggregateType, id),
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *WeatherStationAggregate) HandleCommand(ctx context.Context, cmd eh.Command) error {
	switch cmd := cmd.(type) {
	case *Create:
		a.StoreEvent(WeatherStationCreatedEvent,
			&WeatherStationCreatedData{
				cmd.Name,
			},
			time.Now(),
		)
		return nil

	case *ReportTemperature:
		a.StoreEvent(TemperatureReportedEvent,
			&TemperatureReportedData{
			cmd.Temperature,
			},
			time.Now(),
		)
		return nil

	}
	return fmt.Errorf("couldn't handle command")
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *WeatherStationAggregate) ApplyEvent(ctx context.Context, event eh.Event) error {
	switch event.EventType() {
	case WeatherStationCreatedEvent:
		if data, ok := event.Data().(*WeatherStationCreatedData); ok {
			a.Name = data.Name
		} else {
			log.Println("invalid event data type:", event.Data())
		}
	case TemperatureReportedEvent:
		if data, ok := event.Data().(*TemperatureReportedData); ok {
			a.Temperature = data.Temperature
		} else {
			log.Println("invalid event data type:", event.Data())
		}
	}
	return nil
}

