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

// Package mongodb contains an example of a CQRS/ES app using the MongoDB adapter.
package weatherstation

import (
	"context"
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
	"github.com/looplab/eventhorizon/commandhandler/aggregate"
	"github.com/looplab/eventhorizon/commandhandler/bus"
	eventbus "github.com/looplab/eventhorizon/eventbus/local"
	"github.com/looplab/eventhorizon/eventhandler/projector"
	eventstore "github.com/looplab/eventhorizon/eventstore/mongodb"
	"log"
	"os"

	"github.com/looplab/eventhorizon/examples/weatherstation/domain"
)

const RepoHost = "localhost"
const RepoPort = "27017"
const DbPrefix = "weatherstation"
const DbNamespace = "nsweatherstation"
const LocationName = "Athens"

type Store struct {
	eventStore *eventstore.EventStore
	eventBus *eventbus.EventBus
	commandBus *bus.CommandHandler
	aggregateStore *events.AggregateStore
}

func NewStore() *Store {
	// Local Mongo testing with Docker
	url := os.Getenv("MONGO_HOST")
	if url == "" {
		// Default to localhost
		//url = "localhost:27017"
		url = RepoHost + ":" + RepoPort
	}

	var err error
	var eventStore *eventstore.EventStore
	var eventBus *eventbus.EventBus
	var commandBus *bus.CommandHandler
	var aggregateStore *events.AggregateStore
	var commandHandler *aggregate.CommandHandler

	// Create the event store. The collection used is currently always 'events'
	eventStore, err = eventstore.NewEventStore(url, DbPrefix)
	if err != nil {
		log.Fatalf("could not create event store: %s", err)
	}

	// Create the event bus that distributes events.
	eventBus = eventbus.NewEventBus(nil)
	go func() {
		for e := range eventBus.Errors() {
			log.Printf("eventbus: %s", e.Error())
		}
	}()

	// Create the command bus.
	commandBus = bus.NewCommandHandler()

	// Add a logger as an observer.
	eventBus.AddObserver(eh.MatchAny(), &domain.Logger{})

	// Create the aggregate repository.
	aggregateStore, err = events.NewAggregateStore(eventStore, eventBus)
	if err != nil {
		log.Fatalf("could not create aggregate store: %s", err)
	}

	// Create the aggregate command handler and register the commands it handles.
	commandHandler, err = aggregate.NewCommandHandler(domain.WeatherStationAggregateType, aggregateStore)
	if err != nil {
		log.Fatalf("could not create command handler: %s", err)
	}

	loggingCommandHandler := eh.UseCommandHandlerMiddleware(commandHandler, domain.LoggingMiddleware)
	err = commandBus.SetHandler(loggingCommandHandler, domain.CreateWeatherStationCommand)
	if err != nil {
		log.Fatalf("could not create command handler: %s", err)
	}
	err = commandBus.SetHandler(loggingCommandHandler, domain.ReportTemperatureCommand)
	if err != nil {
		log.Fatalf("could not create command handler: %s", err)
	}

	return &Store{
		eventStore,
		eventBus,
		commandBus,
		aggregateStore,
	}
}

func (ws *Store) Create(ctx context.Context, id uuid.UUID, name string) {
	if err := ws.commandBus.HandleCommand(ctx, &domain.Create{ID: id, Name: name}); err != nil {
		log.Println("error:", err)
	}
}

func (ws *Store) ReportTemperature(ctx context.Context, id uuid.UUID, temp float32) {
	if err := ws.commandBus.HandleCommand(ctx, &domain.ReportTemperature{ID: id, Temperature: temp}); err != nil {
		log.Println("error:", err)
	}
}

func (ws *Store) AddProjector(projectorEventHandler *projector.EventHandler) {
	ws.eventBus.AddHandler(eh.MatchAnyEventOf(
		domain.WeatherStationCreatedEvent,
		domain.TemperatureReportedEvent,
	), projectorEventHandler)
}

func (ws *Store) CleanEventStore(ctx context.Context) {
	err := ws.eventStore.Clear(ctx)
	if err != nil {
		log.Fatalf("could not create command handler: %s", err)
	}
}


