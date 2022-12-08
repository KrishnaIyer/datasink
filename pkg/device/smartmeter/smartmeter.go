// Copyright © 2022 Krishna Iyer Easwaran
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

// Package smartmeter parses data from the Smart Gateways smart meter (https://smartgateways.nl/product/slimme-meter-wifi-gateway/).
// Reference: https://smartgateways.nl/slimme-meter-p1-dsmr-uitlezen/
package smartmeter

import (
	"context"
	"strconv"

	"krishnaiyer.dev/golang/datasink/pkg/database/entry"
	"krishnaiyer.dev/golang/dry/pkg/logger"
)

const (
	// Identifier identifies the device type.
	Identifier = "dsmr"

	rootPrefix  = "dsmr"
	measurement = "smartmeter"

	// Message types.
	messageTypeReading     = "reading"
	messageTypeInfo        = "smart_gateways"
	messageTypeConsumption = "consumption"

	// Reading Subtypes.
	readingTypeElectricityEquipmentID      = "electricity_equipment_id"
	readingElectricityHourlyUsage          = "electricity_hourly_usage"
	readingTypeElectricityDelivered1       = "electricity_delivered_1"
	readingTypeElectricityReturned1        = "electricity_returned_1"
	readingTypeElectricityDelivered2       = "electricity_delivered_2"
	readingTypeElectricityReturned2        = "electricity_returned_2"
	readingTypeElectricityDeliveredCurrent = "electricity_currently_delivered"
	readingTypeElectricityReturnedCurrent  = "electricity_currently_returned"
	// Info Subtypes.
	infoWiFiRSSI = "wifi_rssi"
)

// SmartMeter is a smart meter.
type Meter struct {
}

// Parse implements device.Device.
// The value returned could be nil without error. Callers must skip these.
// This function does not error on unknown message types to prevent a rogue device from crashing the server.
func (m Meter) Parse(ctx context.Context, id, dataType, key string, value []byte) (*entry.Entry, error) {
	logger := logger.LoggerFromContext(ctx).WithField("id", id)
	switch dataType {
	case messageTypeReading:
		return parseReading(ctx, id, key, value)
	case messageTypeInfo:
		return parseInfo(ctx, id, key, value)
	case messageTypeConsumption:
		return parseConsumption(ctx, id, key, value)
	default:
		logger.WithField("type", dataType).Warn("unknown message type")
		return nil, nil
	}
}

// parseReadings readings `dsmr/reading/*` messages.
func parseReading(ctx context.Context, id string, reading string, value []byte) (*entry.Entry, error) {
	var (
		ret    *entry.Entry
		fields = make(map[string]any)
	)
	switch reading {
	case readingTypeElectricityEquipmentID:
		fields[readingTypeElectricityEquipmentID] = string(value)
	case readingElectricityHourlyUsage:
		if v, err := strconv.ParseFloat(string(value), 64); err == nil {
			fields[readingElectricityHourlyUsage] = v
		}
	case readingTypeElectricityDelivered1:
		if v, err := strconv.ParseFloat(string(value), 64); err == nil {
			fields[readingTypeElectricityDelivered1] = v
		}
	case readingTypeElectricityReturned1:
		if v, err := strconv.ParseFloat(string(value), 64); err == nil {
			fields[readingTypeElectricityReturned1] = v
		}
	case readingTypeElectricityDelivered2:
		if v, err := strconv.ParseFloat(string(value), 64); err == nil {
			fields[readingTypeElectricityDelivered2] = v
		}
	case readingTypeElectricityReturned2:
		if v, err := strconv.ParseFloat(string(value), 64); err == nil {
			fields[readingTypeElectricityReturned2] = v
		}
	case readingTypeElectricityDeliveredCurrent:
		if v, err := strconv.ParseFloat(string(value), 64); err == nil {
			fields[readingTypeElectricityDeliveredCurrent] = v
		}
	case readingTypeElectricityReturnedCurrent:
		if v, err := strconv.ParseFloat(string(value), 64); err == nil {
			fields[readingTypeElectricityReturnedCurrent] = v
		}
	}
	if len(fields) != 0 {
		ret = &entry.Entry{
			Measurement: measurement,
			Tags: map[string]string{
				"id": id,
			},
			Fields: fields,
		}
	}
	return ret, nil
}

// parseInfo parses `dsmr/smart_gateways/*` messages.
func parseInfo(ctx context.Context, id string, info string, value []byte) (*entry.Entry, error) {
	var (
		ret    *entry.Entry
		fields = make(map[string]any)
	)
	switch info {
	case infoWiFiRSSI:
		if v, err := strconv.Atoi(string(value)); err == nil {
			fields[infoWiFiRSSI] = v
		}
	}
	if len(fields) != 0 {
		ret = &entry.Entry{
			Measurement: measurement,
			Tags: map[string]string{
				"id": id,
			},
			Fields: fields,
		}
	}
	return ret, nil
}

// parseConsumption parses `dsmr/consumption/*` messages.
func parseConsumption(ctx context.Context, id string, reading string, data []byte) (*entry.Entry, error) {
	return nil, nil
}
