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

	"krishnaiyer.dev/golang/datasink/pkg/database/entry"
)

// Device is an IoT device.
type Device interface {
	// Parse parses device data on a particular topic.
	Parse(ctx context.Context, id, dataType, key string, value []byte) (*entry.Entry, error)
}
