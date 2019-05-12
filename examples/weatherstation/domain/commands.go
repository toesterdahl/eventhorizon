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
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
)

func init() {
	eh.RegisterCommand(func() eh.Command { return &Create{} })
	eh.RegisterCommand(func() eh.Command { return &ReportTemperature{} })
}

const (
	CreateWeatherStationCommand eh.CommandType = "CreateWeatherStation"
	ReportTemperatureCommand    eh.CommandType = "ReportTemperature"
)

// Create is a command for creating a weather station.
type Create struct {
	ID   uuid.UUID
	Name string
}

func (c Create) AggregateID() uuid.UUID          { return c.ID }
func (c Create) AggregateType() eh.AggregateType { return WeatherStationAggregateType }
func (c Create) CommandType() eh.CommandType     { return CreateWeatherStationCommand }

// ReportTemperature is a command for reporting temperature.
type ReportTemperature struct {
	ID          uuid.UUID
	Temperature float32
}

func (c ReportTemperature) AggregateID() uuid.UUID          { return c.ID }
func (c ReportTemperature) AggregateType() eh.AggregateType { return WeatherStationAggregateType }
func (c ReportTemperature) CommandType() eh.CommandType     { return ReportTemperatureCommand }

