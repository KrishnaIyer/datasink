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

// Package device provides functions to decode device data and creates database entries.
package device

import (
	"context"
	"fmt"

	"krishnaiyer.dev/golang/datasink/pkg/database/entry"
	"krishnaiyer.dev/golang/datasink/pkg/device/smartmeter"
)

// Config is the configuration for devices.
type Config struct {
	SmartMeter smartmeter.Config `name:"smart-meter" description:"smartmeter configuration"`
}

func (c Config) GetParser(ctx context.Context, key string) (Device, error) {
	if c.SmartMeter.SupportsKey(key) {
		return c.SmartMeter, nil
	}
	return nil, fmt.Errorf("no device found for key %s", key)
}

// Device is an IoT device.
type Device interface {
	// Parse parses device data on a particular topic.
	Parse(ctx context.Context, id, key string, value []byte) (*entry.Entry, error)
	// SupportsKey returns true if the device supports the given key.
	SupportsKey(key string) bool
}
