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

// Package mqtt provides an MQTT broker.
package mqtt

import (
	"context"

	"go.krishnaiyer.dev/dry/pkg/logger"

	mqttnet "github.com/TheThingsIndustries/mystique/pkg/net"
	mqtt "github.com/TheThingsIndustries/mystique/pkg/server"
)

// Config is the configuration for the MQTT server.
type Config struct {
	Addr string `name:"address" description:"server address"`
}

// Server is an MQTT server.
type Server struct {
	srv mqtt.Server
	c   Config
}

// New creates a new Server.
func New(ctx context.Context, c Config) *Server {
	return &Server{
		srv: mqtt.New(ctx),
		c:   c,
	}
}

// Start starts the MQTT server.
func (s *Server) Start(ctx context.Context) error {
	logger := logger.LoggerFromContext(ctx)

	// Start a TCP listener at the given address.
	lis, err := mqttnet.Listen("tcp", s.c.Addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	go func() {
		for {
			// Each connection here is equivalent to an MQTT `CONNECT`.
			conn, err := lis.Accept()
			if err != nil {
				// TODO: Add WithError() method to logger.
				logger.Error(err.Error())
				return
			}
			go s.srv.Handle(conn)
		}
	}()
	logger.WithField("address", s.c.Addr).Info("Start MQTT server")

	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}
