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
	"fmt"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"krishnaiyer.dev/golang/datasink/pkg/database/entry"
	"krishnaiyer.dev/golang/dry/pkg/logger"
)

const (
	// DefaultWriteTimeout is the default write timeout.
	DefaultWriteTimeout = 5 * time.Second
)

// SetupOptions are used to setup the database.
type SetupOptions struct {
	Username           string `name:"username" description:"username"`
	Password           string `name:"password" description:"password"`
	RetentionPeriodHrs int    `name:"retention_period_hrs" description:"retention period in hours"`
}

// NonBlockingWrites uses the non-blocking write API.
// This scales well but is also more prone to error.
// In case of a crash, data may be lost.
// If set to false (default), blocking write API is used, which is more reliable.
type NonBlockingWrites struct {
	Enabled       bool `name:"enabled" description:"enable non-blocking writes"`
	BatchSize     int  `name:"batch_size" description:"batch size"`
	FlushInterval int  `name:"flush_interval" description:"flush interval"`
}

// Config configures the InfluxDB client.
type Config struct {
	NonBlockingWrites NonBlockingWrites `name:"non_blocking_writes"`
	Address           string            `name:"address" description:"server address"`
	Token             string            `name:"token" description:"auth token. Generate a random one using 'openssl rand -hex 32'"`
	Bucket            string            `name:"bucket" description:"data bucket"`
	Organization      string            `name:"organization" description:"organization"`
	WriteTimeout      time.Duration     `name:"write_timeout" description:"write timeout in seconds (for blocking writes)"`
	SetupOpts         SetupOptions      `name:"setup" description:"setup options"`
}

// Client is an InfluxDB client.
type Client struct {
	cfg Config
	cl  influxdb.Client
}

// Setup sets up the database with the provided config.
// Use this function only on the first run and for testing.
func (cfg Config) Setup(ctx context.Context) error {
	cl := influxdb.NewClient(cfg.Address, cfg.Token) // The token is technically not required for setup.
	defer cl.Close()
	_, err := cl.SetupWithToken(
		ctx,
		cfg.SetupOpts.Username,
		cfg.SetupOpts.Password,
		cfg.Organization,
		cfg.Bucket,
		cfg.SetupOpts.RetentionPeriodHrs,
		cfg.Token,
	)
	return err
}

// New returns a new client.
// Use Close() to close the client after done.
func (c *Config) NewClient(ctx context.Context) *Client {
	options := influxdb.DefaultOptions()
	if c.NonBlockingWrites.Enabled {
		options.SetBatchSize(uint(c.NonBlockingWrites.BatchSize))
		options.SetFlushInterval(uint(c.NonBlockingWrites.FlushInterval))
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = DefaultWriteTimeout
	}
	cl := influxdb.NewClientWithOptions(c.Address, c.Token,
		options,
	)
	return &Client{
		cfg: *c,
		cl:  cl,
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

func (c *Client) Query(ctx context.Context, query string) (map[time.Time]any, error) {
	queryAPI := c.cl.QueryAPI(c.cfg.Organization)

	logger := logger.LoggerFromContext(ctx).WithField("query", query)

	logger.Debug("Run query")

	ret := make(map[time.Time]any)

	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		if result.TableChanged() {
			logger.WithField("table", result.TableMetadata().String()).Debug("Table changed")
		}
		ret[result.Record().Time()] = result.Record().Value()
	}
	if result.Err() != nil {
		return nil, fmt.Errorf("query parsing error: %s\n", result.Err().Error())
	}
	return ret, nil
}
