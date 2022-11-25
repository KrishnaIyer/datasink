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
	logger "go.krishnaiyer.dev/dry/pkg/logger"
	"go.krishnaiyer.dev/mqtt-influx/pkg/http"
)

// Config contains the configuration.
type Config struct {
	HTTPAddress string `name:"http-address" description:"Address where the instrumentation HTTP server is served"`
}

var (
	config  = new(Config)
	manager *conf.Manager

	// Root is the root of the commands.
	Root = &cobra.Command{
		Use:           "mqtt-influx",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "mqtt-influx is tool that acts as an MQTT broker for incoming traffic and writes it to an Influx DB instance.",
		Long:          `mqtt-influx is tool that acts as an MQTT broker for incoming traffic and writes it to an Influx DB instance. More documentation at https://go.krishnaiyer.dev/mqtt-influx`,
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

			l, err := logger.New(baseCtx, false)
			if err != nil {
				panic(err)
			}

			ctx = logger.NewContextWithLogger(baseCtx, l)

			// Start the HTTP Server.
			s := http.New()
			go func() {
				err := s.Start(ctx, config.HTTPAddress)
				if err != nil {
					log.Fatal(err)
				}
			}()
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