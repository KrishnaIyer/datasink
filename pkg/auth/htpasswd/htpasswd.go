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

// Package htpasswd reads htpasswd files.
package htpasswd

import (
	htpasswd "github.com/tg123/go-htpasswd"
)

// Store is a htpasswd store.
type Store struct {
	store *htpasswd.File
}

// NewStore returns a new auth Store.
func NewStore(file string) (*Store, error) {
	st, err := htpasswd.New(file, htpasswd.DefaultSystems, nil)
	if err != nil {
		return nil, err
	}
	return &Store{
		store: st,
	}, nil
}

func (st *Store) Verify(user, pass string) bool {
	return st.store.Match(user, pass)
}
