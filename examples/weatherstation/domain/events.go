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
	eh "github.com/looplab/eventhorizon"
)

const (
	// WeatherStationCreatedEvent is when a new weather station is created.
	WeatherStationCreatedEvent eh.EventType = "WeatherStationCreated"

	// TemperatureReportedEvent is when a temperature is reported.
	TemperatureReportedEvent eh.EventType = "TemperatureReported"

)

func init() {
	// Only the event for creating a weather station has custom data.
	eh.RegisterEventData(WeatherStationCreatedEvent, func() eh.EventData {
		return &WeatherStationCreatedData{}
	})
	eh.RegisterEventData(TemperatureReportedEvent, func() eh.EventData {
		return &TemperatureReportedData{}
	})}

// WeatherStationCreatedData is the event data for when a new weather station is created.
type WeatherStationCreatedData struct {
	Name  string    `bson:"age"`
}

// WeatherStationCreatedData is the event data for when a temperature is reported.
type TemperatureReportedData struct {
	Temperature  float32    `bson:"age"`
}