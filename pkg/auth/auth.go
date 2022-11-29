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

// Package auth provides reusable auth function.
// Currently reading passwords from a htpasswd file is supported.
package auth

import (
	"fmt"

	"krishnaiyer.dev/golang/datasink/pkg/auth/htpasswd"
)

// Store is a generic auth store.
type Store interface {
	Verify(user, pass string) bool
}

// Config is the auth configuration.
type Config struct {
	Type         string `name:"type" description:"authentication file type. Supported values are 'htpasswd'"`
	HtpasswdFile string `name:"htpasswd-file" description:"location of the htpasswd file"`
}

// NewStore creates a new auth store.
func (c Config) NewStore() (Store, error) {
	switch c.Type {
	case "htpasswd":
		return htpasswd.NewStore(c.HtpasswdFile)
	default:
		return nil, fmt.Errorf("invalid auth type '%s'", c.Type)
	}
}
