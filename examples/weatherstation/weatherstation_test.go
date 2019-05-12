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

package weatherstation

import (
	"context"
	"github.com/brutella/hc/log"
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/examples/weatherstation/domain"
	"testing"
)

// The purpose of this test is demonstration three cases of using projections,
// 1. No projection at all. A persisted Aggregate may be recovered from the aggregate store.
// 2. Simple projection. The projection is created before creating any new aggregates.
// 3. The projection is created after the aggregate. The projection is restored from the persisted aggregate.
// 4. Same as (3.) but instead of relying of the aggregate's uuid we retrieve it based on a natural key.

var id = uuid.New()

// TestCreateWeatherStation will set up an eventstore WITHOUT a projection.
// It will issue commands,
// results will be verified against the aggregate store.
func TestCreateWeatherStation(t *testing.T) {
	// Set the namespace to use.
	ctx := eh.NewContextWithNamespace(context.Background(), DbNamespace)

	ws := NewStore()

	// FIXME: Clean fails if store does not exist!
	ws.CleanEventStore(ctx)

	ws.Create(ctx, id, LocationName)
	ws.ReportTemperature(ctx, id, 1.0)
	ws.ReportTemperature(ctx, id, 2.0)
	ws.ReportTemperature(ctx, id, 3.0)

	a, err := ws.aggregateStore.Load(ctx, domain.WeatherStationAggregateType, id)
	if err != nil {
		t.Errorf("Error loading aggregate %s %s", domain.WeatherStationAggregateType, id)
	}

	weatherStation := a.(*domain.WeatherStationAggregate)
	log.Info.Printf("Aggregate type: %s version: %d name: %s temperature (latest): %f", weatherStation.AggregateType(), weatherStation.Version(), weatherStation.Name, weatherStation.Temperature)

	if weatherStation.AggregateType() != "Temperature" {
		t.Errorf("Expected aggregate type '%s', got %s\n", "Temperature", weatherStation.AggregateType())
	}

	if weatherStation.AggregateType() != "Temperature" {
		t.Errorf("Expected aggregate type '%s', got %s\n", "Temperature", weatherStation.AggregateType())
	}

	if weatherStation.Name != LocationName {
		t.Errorf("Expected aggregate name '%s', got %s\n", LocationName, weatherStation.Name)
	}

	if weatherStation.Temperature != 3.0 {
		t.Errorf("Expected aggregate type '%s', got %s\n", 3.0, weatherStation.Temperature)
	}

}

// TestRestoreProjector will create a projector
// it will first initialize the same event eventstore as TestCreateWeatherStation
// (yes, tests are dependent, itch, but it is necessary to prove the point)
// The projector is 'restored' to catch up with the event store.
// An additional command is applied,
// results are checked against the projector.
func TestRestoreProjector(t *testing.T) {
	// Set the namespace to use.
	ctx := eh.NewContextWithNamespace(context.Background(), DbNamespace)

	// This is the same event store used previously.
	ws := NewStore()
	// The projector manager takes care of operations on the projection that require access to the event store.
	pm := SetupProjectorManager(ws.aggregateStore)
	//
	p := SetupProjector()
	ws.AddProjector(p.projectorEventHandler)

	// FIXME: Clean fails if repo does not exist!
	p.ClearRepo(ctx)
	pm.Restore(ctx, id)
	ws.ReportTemperature(ctx, id, 4.0)

	var resultName string
	var resultTemp float32
	var resultHistory []float32
	resultName, resultTemp, resultHistory = p.ListTemperatureHistory(ctx)

	if resultName != LocationName {
		t.Errorf("Expected Weather Station name %s got %s\n", LocationName, resultName)
	}

	if resultTemp != 4.0 {
		t.Errorf("Expected Weather Station temperature %f got %f\n", 4.0, resultTemp)
	}

	if len(resultHistory) != 4.0 {
		t.Errorf("Expected Weather Station history length %f got %f\n", 4.0, len(resultHistory))
	}

	if resultHistory[2] != 3.0 {
		t.Errorf("Expected Weather Station history last temperature %f got %f\n", 4.0, resultHistory[3])
	}

	if resultHistory[3] != 4.0 {
		t.Errorf("Expected Weather Station history last temperature %f got %f\n", 4.0, resultHistory[3])
	}
}
