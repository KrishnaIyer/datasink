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
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"krishnaiyer.dev/golang/datasink/pkg/database"
	"krishnaiyer.dev/golang/datasink/pkg/http"
	"krishnaiyer.dev/golang/datasink/pkg/mqtt"
	conf "krishnaiyer.dev/golang/dry/pkg/config"
	logger "krishnaiyer.dev/golang/dry/pkg/logger"
)

const (
	defaultBufferSize = 64
)

// Config contains the configuration.
type Config struct {
	HTTP     http.Config     `name:"http"`
	MQTT     mqtt.Config     `name:"mqtt"`
	Database database.Config `name:"database"`
}

var (
	config  = &Config{}
	manager *conf.Manager
	baseCtx = context.Background()

	// Root is the root of the commands.
	Root = &cobra.Command{
		Use:           "datasink",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "datasink is tool that acts as acts as a server with multiple protocols (ex: mqtt, websocket) for incoming traffic and writes to a time series database",
		Long:          `datasink is tool that acts as acts as a server with multiple protocols (ex: mqtt, websocket) for incoming traffic and writes to a time series database. More documentation at https://krishnaiyer.dev/golang/datasink`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := manager.ReadFromFile(cmd.Flags())
			if err != nil {
				panic(err)
			}
			err = manager.Unmarshal(&config)
			if err != nil {
				panic(err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(baseCtx)
			defer cancel()

			l, err := logger.New(ctx, false)
			if err != nil {
				panic(err)
			}
			ctx = logger.NewContextWithLogger(ctx, l)

			errCh := make(chan error)
			defer close(errCh)

			var database database.Database
			switch config.Database.Type {
			case "influxdb":
				// Create Client.
				client := config.Database.InfluxDB.NewClient(ctx)
				database = client
			}
			defer database.Close(ctx)

			// Start the HTTP Server.
			go func() {
				s := http.New(config.HTTP)
				err = s.Start(ctx)
				if err != nil {
					errCh <- err
					return
				}
			}()

			// Start the MQTT Server.
			messageCh := make(chan mqtt.Message, defaultBufferSize)
			defer close(messageCh)
			go func() {
				s, err := mqtt.New(ctx, config.MQTT, messageCh)
				if err != nil {
					errCh <- err
					return
				}
				err = s.Start(ctx)
				if err != nil {
					errCh <- err
					return
				}
			}()

			// Listen for messages and write to database.
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case msg := <-messageCh:
						fmt.Println(msg)
						// err := database.Write(ctx, msg)
						// if err != nil {
						// 	l.WithError(err).Error("Error writing to database")
						// }
					}
				}
			}()

			select {
			case err := <-errCh:
				return err
			case <-ctx.Done():
				return ctx.Err()
			default:
				// Wait for a signal to stop the server.
				sigChan := make(chan os.Signal, 1)
				signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
				signal := (<-sigChan).String()
				l.WithField("signal", signal).Info("Signal received. Shut down server")
				return nil
			}
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
	// This line is needed to persist the config file to subcommands.
	manager.AddConfigFlag(manager.Flags())
	Root.PersistentFlags().AddFlagSet(manager.Flags())
	Root.AddCommand(VersionCommand(Root))
	Root.AddCommand(InitDBCommand(Root))
	Root.AddCommand(ConfigCommand(Root))
}
