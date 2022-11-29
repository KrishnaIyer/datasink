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

// Package database abstracts database functions.
package database

import (
	"context"

	"go.krishnaiyer.dev/datasink/pkg/database/entry"
	"go.krishnaiyer.dev/datasink/pkg/database/influxdb"
)

// Config defines the database configuration.
type Config struct {
	Type     string          `name:"type" description:"The type of database to use. Ex: influxdb"`
	InfluxDB influxdb.Config `name:"influxdb"`
}

// Database is a database.
type Database interface {
	// Record records an entry.
	Record(ctx context.Context, entry entry.Entry) error
}
