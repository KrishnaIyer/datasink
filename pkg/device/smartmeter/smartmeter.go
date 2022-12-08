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
	"strings"

	"krishnaiyer.dev/golang/datasink/pkg/database/entry"
	"krishnaiyer.dev/golang/dry/pkg/logger"
)

const (
	rootPrefix  = "dsmr"
	measurement = "smartmeter"
)

// Config is the configuration for the smart meter.
type Config struct {
	Values map[string]string `name:"values" description:"Values to record and the corresponding data type"`
}

// SupportsKey implements device.Device.
func (c Config) SupportsKey(key string) bool {
	k := strings.Split(key, "/")
	if len(k) > 0 && k[0] == rootPrefix {
		return true
	}
	return false
}

// Parse implements device.Device.
// The value returned could be nil without error. Callers must skip these.
// This function does not error on unknown message types to prevent a rogue device from crashing the server.
func (c Config) Parse(ctx context.Context, id, key string, value []byte) (*entry.Entry, error) {
	logger := logger.LoggerFromContext(ctx).WithField("id", id)

	// Split the key get the last part.
	k := strings.Split(key, "/")
	if len(k) == 0 {
		logger.WithField("key", key).Warn("invalid key")
		return nil, nil
	}
	dbKey := k[len(k)-1]
	typ, ok := c.Values[dbKey]
	if !ok {
		logger.WithField("key", key).Info("Key not configured for logging, skip")
		return nil, nil
	}
	fields := make(map[string]any)
	switch typ {
	case "float":
		if v, err := strconv.ParseFloat(string(value), 64); err == nil {
			fields[dbKey] = v
		}
	case "int":
		if v, err := strconv.Atoi(string(value)); err == nil {
			fields[dbKey] = v
		}
	case "string":
		fields[dbKey] = string(value)
	default:
		logger.WithField("type", typ).Warn("unknown type")
		return nil, nil
	}
	return &entry.Entry{
		Measurement: measurement,
		Tags: map[string]string{
			"id": id,
		},
		Fields: fields,
	}, nil
}
