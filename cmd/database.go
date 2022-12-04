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

package cmd

import (
	"context"

	"github.com/spf13/cobra"
	logger "krishnaiyer.dev/golang/dry/pkg/logger"
)

// InitDBCommand initializes the database.
func InitDBCommand(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "init-db",
		Short: "Initialize the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(baseCtx)
			defer cancel()

			l, err := logger.New(ctx, false)
			if err != nil {
				panic(err)
			}
			ctx = logger.NewContextWithLogger(ctx, l)

			l = l.WithField("type", config.Database.Type)

			l.Info("Initialize database")

			switch config.Database.Type {
			case "influxdb":
				err := config.Database.InfluxDB.Setup(ctx)
				if err != nil {
					return err
				}
			default:
				panic("invalid database type")
			}
			return nil
		},
	}
}
