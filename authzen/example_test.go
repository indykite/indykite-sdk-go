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

package authzen_test

import (
	"context"
	"fmt"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/authzen"
)

// One boolean decision, with policy input parameters.
func ExampleClient_Allowed() {
	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx)
	if err != nil {
		return
	}

	ok, err := cli.AuthZEN().Allowed(ctx,
		authzen.NewNode("Person", "ada"), "PROVISION", authzen.NewNode("Server", "gpu-7"),
		authzen.WithInputParams(map[string]any{"max_price": 500}))
	if err != nil {
		return
	}
	fmt.Println(ok)
}

// Many decisions in one round trip: the top-level fields are defaults, each
// entry overrides only what differs.
func ExampleClient_EvaluateBatch() {
	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx)
	if err != nil {
		return
	}

	subject := authzen.NewNode("Person", "ada")
	resp, err := cli.AuthZEN().EvaluateBatch(ctx, authzen.EvaluationsRequest{
		Subject: &subject,
		Action:  &authzen.Action{Name: "PROVISION"},
		Evaluations: []authzen.EvaluationItem{
			{Resource: &authzen.Node{Type: "Server", ID: "gpu-7"}},
			{Resource: &authzen.Node{Type: "Server", ID: "gpu-8"}},
		},
	})
	if err != nil {
		return
	}
	for _, e := range resp.Evaluations {
		fmt.Println(e.Decision)
	}
}

// Which servers may ada provision? The search endpoints enumerate instead of
// answering yes/no.
func ExampleClient_SearchResource() {
	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx)
	if err != nil {
		return
	}

	subject := authzen.NewNode("Person", "ada")
	servers, err := cli.AuthZEN().SearchResource(ctx, authzen.SearchResourceRequest{
		Subject:  &subject,
		Action:   &authzen.Action{Name: "PROVISION"},
		Resource: &authzen.NodeType{Type: "Server"},
	})
	if err != nil {
		return
	}
	for _, s := range servers {
		fmt.Println(s.ID)
	}
}
