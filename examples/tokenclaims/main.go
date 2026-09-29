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

// Command tokenclaims demonstrates policies that read token claims. On every
// AuthZEN decision and ContX IQ execution the platform exposes:
//
//   - $token: the claims of the end-user token sent as "Authorization: Bearer"
//   - $ik_token: the claims of the IndyKite delegated token sent in X-IK-Token,
//     including its RFC 8693 act delegation chain ($ik_token.act.sub is the
//     agent acting for the user)
//
// Both names are reserved: a policy never asks the caller for them, a value
// sent under them in input_params is replaced, and a policy reading a token
// that was not sent denies rather than fails.
//
// setup and teardown use the control plane (Service Account credential);
// evaluate and execute use the runtime plane (App Agent credential) plus the
// two tokens.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/ciq"
	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/examples/internal/exutil"
)

var subcommands = []string{"policies", "setup", "evaluate", "execute", "teardown"}

func main() {
	if len(os.Args) < 2 {
		exutil.Usage("tokenclaims", subcommands...)
	}

	fs := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	projectID := fs.String("project-id", os.Getenv("PROJECT_ID"), "project gid (setup)")
	name := fs.String("name", "sdk-example-claims", "name prefix for created resources (setup)")
	actor := fs.String("actor", "knightrider", "agent the delegated token must name in act.sub")
	subjectID := fs.String("subject-id", "karel", "Person external id; must equal the end-user token sub")
	resourceID := fs.String("resource-id", "docA", "Document external id the subject MANAGEs")
	endUserToken := fs.String("end-user-token", os.Getenv("END_USER_TOKEN"), "end-user token ($token)")
	delegatedToken := fs.String("delegated-token", os.Getenv("IK_TOKEN"), "IndyKite delegated token ($ik_token)")
	queryID := fs.String("query-id", os.Getenv("CLAIMS_QUERY_ID"), "knowledge query id (execute, teardown)")
	kbacPolicyID := fs.String("kbac-policy-id", os.Getenv("CLAIMS_KBAC_POLICY_ID"), "KBAC policy id (teardown)")
	ciqPolicyID := fs.String("ciq-policy-id", os.Getenv("CLAIMS_CIQ_POLICY_ID"), "CIQ policy id (teardown)")
	_ = fs.Parse(os.Args[2:])

	ctx := context.Background()
	switch os.Args[1] {
	case "policies":
		// Just print the two policies, nothing is sent.
		fmt.Println("KBAC policy:")
		fmt.Println(kbacPolicy(*actor))
		fmt.Println("CIQ policy:")
		fmt.Println(ciqPolicy(*actor))

	case "setup":
		setup(ctx, *projectID, *name, *actor)

	case "evaluate":
		evaluate(ctx, *subjectID, *resourceID, *endUserToken, *delegatedToken)

	case "execute":
		execute(ctx, *queryID, *endUserToken, *delegatedToken)

	case "teardown":
		teardown(ctx, *queryID, *ciqPolicyID, *kbacPolicyID)

	default:
		exutil.Usage("tokenclaims", subcommands...)
	}
}

// kbacPolicy lets a Person SHARE a Document they MANAGE, but only when the
// request carries the person's own end-user token and a delegated token whose
// acting agent is actor. The claims are read directly in the graph condition.
func kbacPolicy(actor string) string {
	return mustJSON(map[string]any{
		"meta":     map[string]any{"policy_version": "2.0-kbac"},
		"subject":  map[string]any{"type": "Person"},
		"actions":  []string{"CAN_SHARE"},
		"resource": map[string]any{"type": "Document"},
		"condition": map[string]any{
			"cypher": "MATCH (subject:Person)-[:MANAGE]->(resource:Document) " +
				"WHERE subject.external_id = $token.sub AND $ik_token.act.sub = " + cypherString(actor),
		},
	})
}

// ciqPolicy lets an agent acting for a Person read the external ids of the
// person's documents. $token.sub binds the subject in the cypher; the acting
// agent is checked in a filter whose attribute is the token-claim reference
// $ik_token.act.sub (a supported attribute form, distinct from a graph
// property), so the two ways of reading a claim are both shown.
func ciqPolicy(actor string) string {
	return mustJSON(map[string]any{
		"meta":    map[string]any{"policy_version": "1.0-ciq"},
		"subject": map[string]any{"type": "Person"},
		"condition": map[string]any{
			"cypher": "MATCH (subject:Person)-[:MANAGE]->(doc:Document) WHERE subject.external_id = $token.sub",
			"filter": []map[string]any{
				{"operator": "=", "attribute": "$ik_token.act.sub", "value": actor},
			},
		},
		"allowed_reads": map[string]any{"nodes": []string{"doc.external_id"}},
	})
}

func setup(ctx context.Context, projectID, name, actor string) {
	admin, err := indykite.NewAdminFromEnv(ctx, exutil.Options()...)
	if err != nil {
		exutil.Fatal(err)
	}
	pols := admin.AuthorizationPolicies()

	kbac, err := pols.Create(ctx, &config.CreateAuthorizationPolicy{
		ProjectID: projectID, Name: name + "-kbac", Policy: kbacPolicy(actor), Status: config.StatusActive,
	})
	if err != nil {
		exutil.Fatal(err)
	}
	ciqPol, err := pols.Create(ctx, &config.CreateAuthorizationPolicy{
		ProjectID: projectID, Name: name + "-ciq", Policy: ciqPolicy(actor), Status: config.StatusActive,
	})
	if err != nil {
		exutil.Fatal(err)
	}
	kq, err := admin.KnowledgeQueries().Create(ctx, &config.CreateKnowledgeQuery{
		ProjectID: projectID, Name: name + "-kq", Query: `{"nodes": ["doc.external_id"]}`,
		Status: config.StatusActive, PolicyID: ciqPol.ID,
	})
	if err != nil {
		exutil.Fatal(err)
	}

	fmt.Printf("export CLAIMS_KBAC_POLICY_ID=%s\n", kbac.ID)
	fmt.Printf("export CLAIMS_CIQ_POLICY_ID=%s\n", ciqPol.ID)
	fmt.Printf("export CLAIMS_QUERY_ID=%s\n", kq.ID)
}

func evaluate(ctx context.Context, subjectID, resourceID, endUserToken, delegatedToken string) {
	requireTokens(endUserToken, delegatedToken)
	az := runtime(ctx).AuthZEN()
	subject := authzen.NewNode("Person", subjectID)
	resource := authzen.NewNode("Document", resourceID)

	// Both tokens: $token.sub and $ik_token.act.sub both resolve.
	ok, err := az.Allowed(ctx, subject, "CAN_SHARE", resource,
		authzen.WithEndUserToken(endUserToken),
		authzen.WithDelegatedToken(delegatedToken))
	if err != nil {
		exutil.Fatal(err)
	}
	fmt.Println("both tokens:                       ", ok)

	// No delegated token: $ik_token is an empty claim set, the condition
	// compares against null and the policy denies.
	ok, err = az.Allowed(ctx, subject, "CAN_SHARE", resource,
		authzen.WithEndUserToken(endUserToken))
	if err != nil {
		exutil.Fatal(err)
	}
	fmt.Println("end-user token only (expect false):", ok)

	// A spoofed ik_token input param is replaced by the real claims, so
	// without the header the decision stays false.
	ok, err = az.Allowed(ctx, subject, "CAN_SHARE", resource,
		authzen.WithEndUserToken(endUserToken),
		authzen.WithInputParams(map[string]any{
			"ik_token": map[string]any{"act": map[string]any{"sub": "knightrider"}},
		}))
	if err != nil {
		exutil.Fatal(err)
	}
	fmt.Println("spoofed input param (expect false):", ok)

	// Search works the same way: the resources the subject may share when
	// acting through the delegated agent.
	docs, err := az.SearchResource(ctx, authzen.SearchResourceRequest{
		Subject:  &subject,
		Action:   &authzen.Action{Name: "CAN_SHARE"},
		Resource: &authzen.NodeType{Type: "Document"},
		Context:  &authzen.Context{EndUserToken: endUserToken, DelegatedToken: delegatedToken},
	})
	if err != nil {
		exutil.Fatal(err)
	}
	fmt.Println("shareable documents:")
	exutil.Print(docs)
}

func execute(ctx context.Context, queryID, endUserToken, delegatedToken string) {
	requireTokens(endUserToken, delegatedToken)
	if queryID == "" {
		exutil.Fatal(errors.New("-query-id (or CLAIMS_QUERY_ID) is required; run setup first"))
	}
	q := runtime(ctx).CIQ()

	rows, err := q.All(ctx, ciq.ExecuteRequest{
		ID: queryID, EndUserToken: endUserToken, DelegatedToken: delegatedToken,
	})
	if err != nil {
		exutil.Fatal(err)
	}
	fmt.Printf("both tokens: %d records\n", len(rows))
	exutil.Print(rows)

	// Without the delegated token the $ik_token filter matches nothing.
	rows, err = q.All(ctx, ciq.ExecuteRequest{ID: queryID, EndUserToken: endUserToken})
	if err != nil {
		exutil.Fatal(err)
	}
	fmt.Printf("end-user token only (expect 0): %d records\n", len(rows))
}

func teardown(ctx context.Context, queryID, ciqPolicyID, kbacPolicyID string) {
	admin, err := indykite.NewAdminFromEnv(ctx, exutil.Options()...)
	if err != nil {
		exutil.Fatal(err)
	}
	// The knowledge query references the CIQ policy, so it goes first.
	if queryID != "" {
		if err = admin.KnowledgeQueries().Delete(ctx, queryID, ""); err != nil {
			exutil.Fatal(err)
		}
	}
	for _, id := range []string{ciqPolicyID, kbacPolicyID} {
		if id == "" {
			continue
		}
		if err = admin.AuthorizationPolicies().Delete(ctx, id, ""); err != nil {
			exutil.Fatal(err)
		}
	}
	fmt.Println("cleaned up")
}

func runtime(ctx context.Context) *indykite.Client {
	cli, err := indykite.NewClientFromEnv(ctx, exutil.Options()...)
	if err != nil {
		exutil.Fatal(err)
	}
	return cli
}

func requireTokens(endUserToken, delegatedToken string) {
	if endUserToken == "" || delegatedToken == "" {
		exutil.Fatal(errors.New("-end-user-token and -delegated-token (or END_USER_TOKEN and IK_TOKEN) are required"))
	}
}

// cypherString quotes s as a single-quoted cypher string literal.
func cypherString(s string) string {
	return "'" + strings.NewReplacer(`\`, `\\`, "'", `\'`).Replace(s) + "'"
}

// mustJSON renders v as indented JSON, leaving "->" in the cypher unescaped.
func mustJSON(v any) string {
	var b strings.Builder
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		exutil.Fatal(err)
	}
	return strings.TrimSpace(b.String())
}
