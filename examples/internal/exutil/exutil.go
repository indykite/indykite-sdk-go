// Copyright (c) 2026 IndyKite
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package exutil holds the few helpers shared by all examples.
package exutil

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/indykite/indykite-sdk-go/transport"
)

// Options resolves the transport options every example uses: an explicit
// INDYKITE_BASE_URL wins, otherwise the eu region.
func Options() []transport.Option {
	if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
		return []transport.Option{transport.WithBaseURL(base)}
	}
	return []transport.Option{transport.WithRegion("eu")}
}

// Print pretty-prints any API result as JSON.
func Print(v any) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		Fatal(err)
	}
	fmt.Println(string(out))
}

// Fatal reports an error (with platform details when it is an APIError) and exits.
func Fatal(err error) {
	if apiErr, ok := transport.AsAPIError(err); ok {
		fmt.Fprintf(os.Stderr, "API error (HTTP %d, request %s): %s\n",
			apiErr.StatusCode, apiErr.RequestID, apiErr.Message)
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(1)
}

// Usage prints the subcommand list and exits.
func Usage(name string, subcommands ...string) {
	fmt.Fprintf(os.Stderr, "usage: %s <subcommand> [flags]\nsubcommands:\n", name)
	for _, s := range subcommands {
		fmt.Fprintf(os.Stderr, "  %s\n", s)
	}
	os.Exit(2)
}
