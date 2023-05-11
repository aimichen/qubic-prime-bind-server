// Copyright © 2021 AMIS Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"net/http"
	"time"
)

var (
	DefaultClient *http.Client
)

func init() {
	// Copy from http default transport
	defaultTransport := http.DefaultTransport.(*http.Transport)
	newTransport := &http.Transport{}
	//nolint:govet
	*newTransport = *defaultTransport

	// Set idle conn timeout to 55 seconds to avoid reset by peer error
	// The default idle timoeut in:
	// Azure: 240s
	// Aws: 60s
	// Gcp: 240s
	newTransport.IdleConnTimeout = 55 * time.Second
	newTransport.DisableKeepAlives = true
	DefaultClient = &http.Client{
		Transport: newTransport,
	}
}
