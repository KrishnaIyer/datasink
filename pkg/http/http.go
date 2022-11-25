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

// Package HTTP provides a simple HTTP Server for instrumentation.
package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server is an HTTP server.
type Server struct {
	s *http.Server
}

// New creates a new Server.
func New() *Server {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	return &Server{
		s: &http.Server{
			Handler:        r,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

// Start starts the HTTP server.
func (s *Server) Start(ctx context.Context, addr string) error {
	s.s.Addr = addr
	select {
	case <-ctx.Done():
		s.s.Shutdown(ctx)
		return ctx.Err()
	default:
		return s.s.ListenAndServe()
	}
}
