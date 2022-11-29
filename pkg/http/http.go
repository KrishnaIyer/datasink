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

// Package http provides a simple HTTP Server for instrumentation.
package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"krishnaiyer.dev/dry/pkg/logger"
)

// Config is the configuration for the HTTP server.
type Config struct {
	Addr string `name:"address" description:"server address"`
}

// Server is an HTTP server.
type Server struct {
	s *http.Server
	c Config
}

// New creates a new Server.
func New(c Config) *Server {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	return &Server{
		c: c,
		s: &http.Server{
			Addr:           c.Addr,
			Handler:        r,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

// Start starts the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	logger.LoggerFromContext(ctx).WithField("address", s.c.Addr).Info("Start HTTP server")
	select {
	case <-ctx.Done():
		s.s.Shutdown(ctx)
		return ctx.Err()
	default:
		return s.s.ListenAndServe()
	}
}
