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

// Package influxdb adds functions to interact with InfluxDB.
package influxdb

import (
	"context"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"go.krishnaiyer.dev/datasink/pkg/database/entry"
)

const (
	// DefaultWriteTimeout is the default write timeout.
	DefaultWriteTimeout = 5 * time.Second
)

// Config configures the InfluxDB client.
type Config struct {
	// NonBlockingWrites uses the non-blocking write API.
	// This scales well but is also more prone to error.
	// In case of a crash, data may be lost.
	// If set to false (default), blocking write API is used, which is more reliable.
	NonBlockingWrites struct {
		Enabled       bool `name:"enabled" description:"enable non-blocking writes"`
		BatchSize     int  `name:"address" description:"batch size"`
		FlushInterval int  `name:"address" description:"flush interval"`
	} `name:"non_blocking_writes"`
	Address      string        `name:"address" description:"server address"`
	Token        string        `name:"token" description:"auth token"`
	Bucket       string        `name:"bucket" description:"data bucket"`
	Organization string        `name:"organization" description:"organization"`
	WriteTimeout time.Duration `name:"write_timeout" description:"write timeout in seconds (for blocking writes)"`
}

// Client is an InfluxDB client.
type Client struct {
	cfg Config
	cl  influxdb.Client
}

// New returns a new client.
// Use Close() to close the client after done.
func (c *Config) NewClient() *Client {
	var options *influxdb.Options
	if c.NonBlockingWrites.Enabled {
		options = influxdb.DefaultOptions().
			SetBatchSize(uint(c.NonBlockingWrites.BatchSize)).
			SetFlushInterval(uint(c.NonBlockingWrites.FlushInterval))
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = DefaultWriteTimeout
	}
	return &Client{
		cfg: *c,
		cl: influxdb.NewClientWithOptions(c.Address, c.Token,
			options,
		),
	}
}

// Close closes the client.
func (c *Client) Close(ctx context.Context) {
	if c.cfg.NonBlockingWrites.Enabled {
		// Flush before closure.
		writeAPI := c.cl.WriteAPI(c.cfg.Organization, c.cfg.Bucket)
		writeAPI.Flush()
	}
	c.cl.Close()
}

// Record implements Database.
// We use the non-blocking write API. This scales well but is also more prone to error.
func (c *Client) Record(ctx context.Context, entry entry.Entry) error {
	point := influxdb.NewPoint(
		entry.Measurement,
		entry.Tags,
		entry.Fields,
		time.Now(),
	)
	if c.cfg.NonBlockingWrites.Enabled {
		writeAPI := c.cl.WriteAPI(c.cfg.Organization, c.cfg.Bucket)
		writeAPI.WritePoint(point)
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, c.cfg.WriteTimeout)
	defer cancel()
	writeAPI := c.cl.WriteAPIBlocking(c.cfg.Organization, c.cfg.Bucket)
	return writeAPI.WritePoint(ctx, point)
}
