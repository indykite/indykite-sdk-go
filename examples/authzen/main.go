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

// Command authzen demonstrates the AuthZEN (authorization) API: single and
// batch decisions plus the three search directions.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/examples/internal/exutil"
)

func main() {
	if len(os.Args) < 2 {
		exutil.Usage("authzen", "evaluate", "batch", "search-action", "search-resource", "search-subject")
	}

	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx, exutil.Options()...)
	if err != nil {
		exutil.Fatal(err)
	}
	az := cli.AuthZEN()

	fs := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	subjectType := fs.String("subject-type", "Person", "subject node type")
	subjectID := fs.String("subject-id", "", "subject external id")
	action := fs.String("action", "", "action name, e.g. CAN_READ")
	resourceType := fs.String("resource-type", "Asset", "resource node type")
	resourceID := fs.String("resource-id", "", "resource external id")
	tag := fs.String("policy-tag", "", "optional policy tag filter")
	_ = fs.Parse(os.Args[2:])

	var opts []authzen.Option
	if *tag != "" {
		opts = append(opts, authzen.WithPolicyTags(*tag))
	}
	subject := authzen.NewNode(*subjectType, *subjectID)
	resource := authzen.NewNode(*resourceType, *resourceID)

	switch os.Args[1] {
	case "evaluate":
		// The one-shot boolean decision; use Evaluate for the full response.
		ok, err := az.Allowed(ctx, subject, *action, resource, opts...)
		if err != nil {
			exutil.Fatal(err)
		}
		fmt.Println("decision:", ok)

	case "batch":
		// One request, many decisions: the top-level subject is the default,
		// each entry overrides only what differs.
		resp, err := az.EvaluateBatch(ctx, authzen.EvaluationsRequest{
			Subject: &subject,
			Action:  &authzen.Action{Name: *action},
			Evaluations: []authzen.EvaluationItem{
				{Resource: &resource},
				{Resource: &authzen.Node{Type: *resourceType, ID: *resourceID + "-2"}},
			},
		})
		if err != nil {
			exutil.Fatal(err)
		}
		exutil.Print(resp.Evaluations)

	case "search-action":
		// What can this subject do with this resource?
		actions, err := az.SearchAction(ctx, authzen.SearchActionRequest{
			Subject: &subject, Resource: &resource,
		})
		if err != nil {
			exutil.Fatal(err)
		}
		exutil.Print(actions)

	case "search-resource":
		// Which resources of this type can the subject perform the action on?
		nodes, err := az.SearchResource(ctx, authzen.SearchResourceRequest{
			Subject:  &subject,
			Action:   &authzen.Action{Name: *action},
			Resource: &authzen.NodeType{Type: *resourceType},
		})
		if err != nil {
			exutil.Fatal(err)
		}
		exutil.Print(nodes)

	case "search-subject":
		// Who can perform the action on this resource?
		nodes, err := az.SearchSubject(ctx, authzen.SearchSubjectRequest{
			Subject:  &authzen.NodeType{Type: *subjectType},
			Action:   &authzen.Action{Name: *action},
			Resource: &resource,
		})
		if err != nil {
			exutil.Fatal(err)
		}
		exutil.Print(nodes)

	default:
		exutil.Usage("authzen", "evaluate", "batch", "search-action", "search-resource", "search-subject")
	}
}
