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
	"log"

	"github.com/spf13/cobra"
	conf "go.krishnaiyer.dev/dry/pkg/config"
)

// Config contains the configuration.
type Config struct {
}

var (
	config  = new(Config)
	manager *conf.Manager

	// Root is the root of the commands.
	Root = &cobra.Command{
		Use:           "mqtt-influx",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "mqtt-influx is a simple command line tool to parse CSV files and convert them to JSON",
		Long:          `mqtt-influx is a simple command line tool to parse CSV files and convert them to JSON. More documentation at https://go.krishnaiyer.dev/mqtt-influx`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := manager.Unmarshal(config)
			if err != nil {
				panic(err)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			baseCtx := context.Background()
			ctx, cancel := context.WithCancel(baseCtx)
			defer cancel()

			_ = ctx

			// logger

			// logger, err := zephyrus.New(context.Background(), config.Debug)
			// if err != nil {
			// 	log.Fatal(err.Error())
			// }
			// defer logger.Clean()
			// ctx := zephyrus.NewContextWithLogger(context.Background(), logger)

		},
	}
)

// Execute ...
func Execute() {
	if err := Root.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	manager = conf.New("config")
	manager.InitFlags(*config)
	Root.PersistentFlags().AddFlagSet(manager.Flags())
	Root.AddCommand(VersionCommand(Root))
}
