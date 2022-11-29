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
	"errors"
	"fmt"
	"io"

	"go.krishnaiyer.dev/datasink/pkg/auth"
	"krishnaiyer.dev/dry/pkg/logger"

	"github.com/TheThingsIndustries/mystique/pkg/apex"
	mqttnet "github.com/TheThingsIndustries/mystique/pkg/net"
	"github.com/TheThingsIndustries/mystique/pkg/packet"
	mqtt "github.com/TheThingsIndustries/mystique/pkg/server"
	"github.com/TheThingsIndustries/mystique/pkg/session"
)

// Config is the configuration for the MQTT server.
type Config struct {
	Addr  string      `name:"address" description:"server address"`
	Debug bool        `name:"debug" description:"enable debug mode"`
	Auth  auth.Config `name:"auth" description:"authentication configuration"`
}

// Server is an MQTT server.
type Server struct {
	srv  mqtt.Server
	c    Config
	auth auth.Store
}

// New creates a new Server.
func New(ctx context.Context, c Config) (*Server, error) {
	if c.Debug {
		apex.SetLevelFromString("debug")
	}

	auth, err := c.Auth.NewStore()
	if err != nil {
		return nil, err
	}
	return &Server{
		srv:  mqtt.New(ctx),
		c:    c,
		auth: auth,
	}, nil
}

// Start starts the MQTT server.
func (s *Server) Start(ctx context.Context) error {
	logger := logger.LoggerFromContext(ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start a TCP listener at the given address.
	lis, err := mqttnet.Listen("tcp", s.c.Addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	logger.WithField("address", s.c.Addr).Info("Start MQTT server")

	// Loop incoming connections and handle them.
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			conn, err := lis.Accept()
			if err != nil {
				return err
			}
			// handleConnection closes the connection when done so we don't need to do it here.
			go s.handleConnection(ctx, conn)
		}
	}
}

// handleConnection handles a single connection.
func (s *Server) handleConnection(ctx context.Context, conn mqttnet.Conn) {
	logger := logger.LoggerFromContext(ctx).WithField("remote_addr", conn.RemoteAddr().String())
	logger.Info("Connect")
	defer func() {
		logger.Info("Disconnect")
		defer conn.Close()
	}()

	session := session.New(ctx, conn, s.deliver)

	// Handle the `CONNECT` packet. This method sends back `CONNACK` packet.
	if err := session.ReadConnect(); err != nil {
		// TODO: Add WithError semantics.
		logger.Error(fmt.Sprintf("Read connect packet: %s", err))
		return
	}
	defer session.Close()

	// Check auth and allowed topic access from the incoming connection.
	if s.auth != nil && !s.auth.Verify(session.AuthInfo().Username, string(session.AuthInfo().Password)) {
		logger.Error("Invalid credentials for user")
		return
	}

	controlCh := make(chan packet.ControlPacket)
	errCh := make(chan error, 1)
	// Read packets from the connection.
	go func() {
		for {
			response, err := session.ReadPacket()
			if err != nil {
				errCh <- err
				close(errCh)
				return
			}
			if response != nil {
				controlCh <- response
			}
		}
	}()

	for {
		var (
			pkt packet.ControlPacket
		)
		select {
		case err := <-errCh:
			if !errors.Is(err, io.EOF) {
				logger.Error(fmt.Sprintf("Read packet: %s", err))
			}
			return
		case pkt = <-controlCh:
			err := conn.Send(pkt)
			if err != nil {
				logger.Error(fmt.Sprintf("Publish packet: %s", err))
				return
			}
		case pkt = <-session.PublishChan():
			// PublishChan intercepts publish packets.
			// We can use this branch to observe latency or check rate limits.
			// Once done, we can use conn.Send().
			// In this case since we are only listening for incoming connections
			// and writing them to a Database, we don't actually publish anything to subscribers.
		}
	}
}

// deliver is a callback attached to the initial session to read all submitted packets.
func (s *Server) deliver(pkt *packet.PublishPacket) {

	// Only store required topics.
	fmt.Println("Received packet:", pkt)
}
