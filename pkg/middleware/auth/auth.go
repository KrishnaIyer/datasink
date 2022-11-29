// Copyright © 2022 Krishna Iyer Easwaran
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package auth is a middleware for checking authentication.
// It supports
// - Basic authentication for HTTP.
package auth

import (
	"net/http"

	"krishnaiyer.dev/golang/datasink/pkg/auth"
)

// Auth abstracts basic authentication.
type Auth struct {
	Store auth.Store
}

// HTTP is a middleware that checks for HTTP Basic Authentication.
func (auth Auth) HTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Basic Authentication credentials
		user, pass, ok := r.BasicAuth()
		if !ok || !auth.Store.Verify(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Not authorized", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
