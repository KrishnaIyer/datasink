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
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.krishnaiyer.dev/datasink/pkg/database"
	"go.krishnaiyer.dev/datasink/pkg/http"
	"go.krishnaiyer.dev/datasink/pkg/mqtt"
	conf "krishnaiyer.dev/dry/pkg/config"
	logger "krishnaiyer.dev/dry/pkg/logger"
)

// Config contains the configuration.
type Config struct {
	HTTP     http.Config     `name:"http"`
	MQTT     mqtt.Config     `name:"mqtt"`
	Database database.Config `name:"database"`
}

var (
	config  Config
	manager *conf.Manager

	// Root is the root of the commands.
	Root = &cobra.Command{
		Use:           "datasink",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "datasink is tool that acts as acts as a server with multiple protocols (ex: mqtt, websocket) for incoming traffic and writes to a time series database",
		Long:          `datasink is tool that acts as acts as a server with multiple protocols (ex: mqtt, websocket) for incoming traffic and writes to a time series database. More documentation at https://go.krishnaiyer.dev/datasink`,
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
			baseCtx := context.Background()
			ctx, cancel := context.WithCancel(baseCtx)
			defer cancel()

			l, err := logger.New(baseCtx, false)
			if err != nil {
				panic(err)
			}
			ctx = logger.NewContextWithLogger(baseCtx, l)

			errCh := make(chan error)
			defer close(errCh)

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
			go func() {
				s, err := mqtt.New(ctx, config.MQTT)
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
	manager.InitFlags(config)
	Root.PersistentFlags().AddFlagSet(manager.Flags())
	Root.AddCommand(VersionCommand(Root))
	manager.AddConfigFlag(Root.Flags())
}
