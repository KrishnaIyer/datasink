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

package influxdb

import (
	"context"
	"os"
	"testing"

	"krishnaiyer.dev/golang/datasink/pkg/database/entry"
	"krishnaiyer.dev/golang/dry/pkg/logger"
)

func TestInfluxDB(t *testing.T) {
	address := os.Getenv("INFLUXDB_ADDRESS")
	if address == "" {
		t.Skip("INFLUXDB_ADDRESS not set")
	}
	ctx := context.Background()
	l, err := logger.New(ctx, false)
	if err != nil {
		panic(err)
	}
	ctx = logger.NewContextWithLogger(ctx, l)

	cfg := Config{
		Address:      address,
		Bucket:       "test",
		Organization: "test",
		Token:        "a817453653634fb34cf07ca8366c6e74",
		SetupOpts: SetupOptions{
			Username:           "testuser",
			Password:           "testpassword",
			RetentionPeriodHrs: 1,
		},
	}
	// Setup Database.
	err = cfg.Setup(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Create Client.
	client := cfg.NewClient(ctx)
	defer client.Close(ctx)

	// Write Data.
	err = client.Record(ctx, entry.Entry{
		Measurement: "test",
		Tags: map[string]string{
			"tag1": "value1",
		},
		Fields: map[string]interface{}{
			"field1": 1,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.Record(ctx, entry.Entry{
		Measurement: "test",
		Tags: map[string]string{
			"tag1": "value1",
		},
		Fields: map[string]interface{}{
			"field1": 2,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.Record(ctx, entry.Entry{
		Measurement: "test",
		Tags: map[string]string{
			"tag1": "value1",
		},
		Fields: map[string]interface{}{
			"field1": 3,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Query Data.
	val, err := client.Query(ctx, "from(bucket: \"test\") |> range(start: -1h) |> filter(fn: (r) => r._measurement == \"test\")")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}
