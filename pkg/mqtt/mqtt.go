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

	"krishnaiyer.dev/golang/datasink/pkg/auth"
	"krishnaiyer.dev/golang/dry/pkg/logger"

	"github.com/TheThingsIndustries/mystique/pkg/apex"
	mqttnet "github.com/TheThingsIndustries/mystique/pkg/net"
	"github.com/TheThingsIndustries/mystique/pkg/packet"
	mqtt "github.com/TheThingsIndustries/mystique/pkg/server"
	"github.com/TheThingsIndustries/mystique/pkg/session"
)

// Config is the configuration for the MQTT server.
type Config struct {
	Addr               string            `name:"address" description:"server address"`
	Debug              bool              `name:"debug" description:"enable debug mode"`
	Auth               auth.Config       `name:"auth" description:"authentication configuration"`
	AllowedTopicPrefix map[string]string `name:"allowed-topic-prefix" description:"allowed topic prefix per username"`
}

// Server is an MQTT server.
type Server struct {
	srv   mqtt.Server
	c     Config
	auth  auth.Store
	msgCh chan *Message
}

// Message is a message received on the MQTT server.
type Message struct {
	Username string
	Topic    string
	Payload  []byte
}

type userSession struct {
	ctx           context.Context
	username      string
	allowedPrefix string
	srv           *Server
}

// New creates a new Server.
func New(ctx context.Context, c Config, messagesCh chan *Message) (*Server, error) {
	if c.Debug {
		apex.SetLevelFromString("debug")
	}
	auth, err := c.Auth.NewStore()
	if err != nil {
		return nil, err
	}
	return &Server{
		srv:   mqtt.New(ctx),
		c:     c,
		auth:  auth,
		msgCh: messagesCh,
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
	defer func() {
		lis.Close()
		logger.Info("Stop MQTT server")
	}()

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

	userSession := &userSession{
		ctx: ctx,
		srv: s,
	}
	session := session.New(ctx, conn, userSession.deliver)

	// Handle the `CONNECT` packet. This method sends back `CONNACK` packet.
	if err := session.ReadConnect(); err != nil {
		logger.WithError(err).Error("Read connect packet: %s")
		return
	}
	defer session.Close()

	// Check auth and allowed topic access from the incoming connection.
	if s.auth != nil && !s.auth.Verify(session.AuthInfo().Username, string(session.AuthInfo().Password)) {
		logger.Error("Invalid credentials for user")
		return
	}

	userSession.allowedPrefix = s.c.AllowedTopicPrefix[session.AuthInfo().Username]
	userSession.username = session.AuthInfo().Username

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

// TODO: Fix this function to allow wildcards.
// Currently it doesn't take different sizes of topics into account.
func isTopicAllowed(requested, allowed []string) bool {
	for i, part := range requested {
		if part == allowed[i] || allowed[i] == "#" {
			continue
		}
		return false
	}
	return true
}

// deliver is a callback attached to the initial session to read all submitted packets.
func (session *userSession) deliver(pkt *packet.PublishPacket) {
	logger := logger.LoggerFromContext(session.ctx).WithField("username", session.username)

	logger.Info("Message received from client")

	if len(pkt.TopicParts) == 0 || pkt.TopicParts[0] != session.allowedPrefix {
		logger.WithField("topic", pkt.TopicName).Error("User not allowed to publish to topic")
	}
	select {
	case <-session.ctx.Done():
		return
	default:
		session.srv.msgCh <- &Message{
			Username: session.username,
			Topic:    pkt.TopicName,
			Payload:  pkt.Message,
		}
	}
}
