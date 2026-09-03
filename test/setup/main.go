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

// Command setup provisions (and removes) the dataset and config-API resources
// the SDK integration tests run against:
//
//	go run ./test/setup apply     # create config resources + ingest the dataset, print env
//	go run ./test/setup env       # print the fixture env vars for the existing resources
//	go run ./test/setup destroy   # delete the dataset and the config resources
//
// It needs INDYKITE_SERVICE_ACCOUNT_CREDENTIALS[_FILE] (control plane),
// INDYKITE_APPLICATION_CREDENTIALS[_FILE] (runtime plane) and PROJECT_ID.
// Apply is idempotent: config resources are found by name and only created when
// missing; graph upserts are idempotent by external id.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/capture"
	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/transport"
)

type dataset struct {
	Nodes         []capture.UpsertNode   `json:"nodes"`
	Relationships []capture.Relationship `json:"relationships"`
}

type policySpec struct {
	Name        string          `json:"name"`
	DisplayName string          `json:"display_name"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	Policy      json.RawMessage `json:"policy"`
}

type querySpec struct {
	Name        string          `json:"name"`
	DisplayName string          `json:"display_name"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	Query       json.RawMessage `json:"query"`
}

type pipelineSpec struct {
	Name                  string          `json:"name"`
	DisplayName           string          `json:"display_name"`
	Description           string          `json:"description"`
	NodeFilter            json.RawMessage `json:"node_filter"`
	SimilarityScoreCutoff float32         `json:"similarity_score_cutoff"`
}

// auditSigningSpec describes the audit-signing fixture. The fixture is
// PLATFORM_MANAGED so it carries no key resource, kid or auth params; the
// customer-managed providers would need real KMS material.
type auditSigningSpec struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Provider    string `json:"provider"`
}

type manifest struct {
	Env            map[string]string `json:"env"`
	ProjectID      string            `json:"-"`
	KQID           string            `json:"-"`
	PipelineID     string            `json:"-"`
	AuditSigningID string            `json:"-"`
	KBACPolicyID   string            `json:"-"`
	CIQPolicyID    string            `json:"-"`
	Query          querySpec         `json:"knowledge_query"`
	loadedDataset  dataset
	KBACPolicy     policySpec       `json:"kbac_policy"`
	CIQPolicy      policySpec       `json:"ciq_policy"`
	AuditSigning   auditSigningSpec `json:"audit_signing"`
	Pipeline       pipelineSpec     `json:"entity_matching_pipeline"`
}

// errUsage signals a usage error; main prints the usage line and exits 2.
var errUsage = errors.New(
	"usage: setup <apply|env|destroy> [-config file] [-dataset file] [-project-id gid]")

func main() {
	err := run(os.Args)
	switch {
	case err == nil:
	case errors.Is(err, errUsage):
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	default:
		log.Fatal(err)
	}
}

// run parses flags, loads the manifest, resolves the project and dispatches the
// subcommand. It returns errors instead of exiting so the logic stays testable.
func run(args []string) error {
	fs := flag.NewFlagSet("setup", flag.ExitOnError)
	configPath := fs.String("config", "test/fixtures/config.json", "config-resources manifest")
	datasetPath := fs.String("dataset", "test/fixtures/dataset.json", "graph dataset")
	projectID := fs.String("project-id", os.Getenv("PROJECT_ID"), "IndyKite project gid")
	if len(args) < 2 {
		return errUsage
	}
	_ = fs.Parse(args[2:])
	if *projectID == "" {
		return errors.New("a project is required: set PROJECT_ID or pass -project-id")
	}

	m, err := load(*configPath, *datasetPath)
	if err != nil {
		return err
	}
	m.ProjectID = *projectID

	ctx := context.Background()
	var opts []indykite.Option
	if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
		opts = append(opts, indykite.WithBaseURL(base))
	}
	admin, err := indykite.NewAdminFromEnv(ctx, opts...)
	if err != nil {
		return err
	}
	if err = m.resolveProject(ctx, admin); err != nil {
		return err
	}
	return m.dispatch(ctx, admin, args[1], opts)
}

// dispatch runs the named subcommand.
func (m *manifest) dispatch(
	ctx context.Context,
	admin *config.AdminClient,
	cmd string,
	opts []indykite.Option,
) error {
	switch cmd {
	case "apply":
		cli, err := indykite.NewClientFromEnv(ctx, opts...)
		if err != nil {
			return err
		}
		if err = m.apply(ctx, admin, cli); err != nil {
			return err
		}
		m.printEnv()
	case "env":
		if err := m.resolve(ctx, admin); err != nil {
			return err
		}
		m.printEnv()
	case "destroy":
		cli, err := indykite.NewClientFromEnv(ctx, opts...)
		if err != nil {
			return err
		}
		return m.destroy(ctx, admin, cli)
	default:
		return errUsage
	}
	return nil
}

func load(configPath, datasetPath string) (*manifest, error) {
	var m manifest
	raw, err := os.ReadFile(configPath) //nolint:gosec // operator-supplied fixture path
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("%s: %w", configPath, err)
	}
	raw, err = os.ReadFile(datasetPath) //nolint:gosec // operator-supplied fixture path
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(raw, &m.loadedDataset); err != nil {
		return nil, fmt.Errorf("%s: %w", datasetPath, err)
	}
	return &m, nil
}

// resolveProject validates PROJECT_ID against the organization's projects,
// accepting either the project gid or its name. On no match it lists what is
// actually available so the operator can pick.
func (m *manifest) resolveProject(ctx context.Context, admin *config.AdminClient) error {
	org, err := admin.Organizations().ReadCurrent(ctx)
	if err != nil {
		return fmt.Errorf("read current organization: %w", err)
	}
	projects, err := admin.Projects().List(ctx, org.ID)
	if err != nil {
		return fmt.Errorf("list projects of %s: %w", org.ID, err)
	}
	for i := range projects {
		if projects[i].ID == m.ProjectID || projects[i].Name == m.ProjectID {
			if projects[i].Name == m.ProjectID {
				log.Printf("resolved project %q to %s", m.ProjectID, projects[i].ID)
			}
			m.ProjectID = projects[i].ID
			return nil
		}
	}
	log.Printf("PROJECT_ID %q matches no project in organization %s (%s); available:", m.ProjectID, org.Name, org.ID)
	for i := range projects {
		log.Printf("  %-30s %s", projects[i].Name, projects[i].ID)
	}
	return fmt.Errorf("set PROJECT_ID (or -project-id) to one of the gids or names above")
}

// ensurePolicy finds a policy by name or creates it, returning its gid.
func (m *manifest) ensurePolicy(ctx context.Context, admin *config.AdminClient, spec *policySpec) (string, error) {
	existing, err := admin.AuthorizationPolicies().Read(ctx, spec.Name, config.WithLocation(m.ProjectID))
	if err == nil {
		log.Printf("policy %q exists: %s", spec.Name, existing.ID)
		return existing.ID, nil
	}
	if apiErr, ok := transport.AsAPIError(err); !ok || !apiErr.IsNotFound() {
		return "", fmt.Errorf("read policy %q: %w", spec.Name, err)
	}
	created, err := admin.AuthorizationPolicies().Create(ctx, &config.CreateAuthorizationPolicy{
		ProjectID:   m.ProjectID,
		Name:        spec.Name,
		DisplayName: spec.DisplayName,
		Description: spec.Description,
		Policy:      string(spec.Policy),
		Status:      spec.Status,
	})
	if err != nil {
		return "", fmt.Errorf("create policy %q: %w", spec.Name, err)
	}
	log.Printf("created policy %q: %s", spec.Name, created.ID)
	return created.ID, nil
}

func (m *manifest) ensureQuery(ctx context.Context, admin *config.AdminClient, ciqPolicyID string) (string, error) {
	existing, err := admin.KnowledgeQueries().Read(ctx, m.Query.Name, config.WithLocation(m.ProjectID))
	if err == nil {
		log.Printf("knowledge query %q exists: %s", m.Query.Name, existing.ID)
		return existing.ID, nil
	}
	if apiErr, ok := transport.AsAPIError(err); !ok || !apiErr.IsNotFound() {
		return "", fmt.Errorf("read knowledge query %q: %w", m.Query.Name, err)
	}
	created, err := admin.KnowledgeQueries().Create(ctx, &config.CreateKnowledgeQuery{
		ProjectID:   m.ProjectID,
		Name:        m.Query.Name,
		DisplayName: m.Query.DisplayName,
		Description: m.Query.Description,
		Query:       string(m.Query.Query),
		Status:      m.Query.Status,
		PolicyID:    ciqPolicyID,
	})
	if err != nil {
		return "", fmt.Errorf("create knowledge query %q: %w", m.Query.Name, err)
	}
	log.Printf("created knowledge query %q: %s", m.Query.Name, created.ID)
	return created.ID, nil
}

func (m *manifest) ensurePipeline(ctx context.Context, admin *config.AdminClient) (string, error) {
	existing, err := admin.EntityMatchingPipelines().Read(ctx, m.Pipeline.Name, config.WithLocation(m.ProjectID))
	if err == nil {
		log.Printf("entity-matching pipeline %q exists: %s", m.Pipeline.Name, existing.ID)
		return existing.ID, nil
	}
	if apiErr, ok := transport.AsAPIError(err); !ok || !apiErr.IsNotFound() {
		return "", fmt.Errorf("read pipeline %q: %w", m.Pipeline.Name, err)
	}
	created, err := admin.EntityMatchingPipelines().Create(ctx, &config.CreateEntityMatchingPipeline{
		ProjectID:             m.ProjectID,
		Name:                  m.Pipeline.Name,
		DisplayName:           m.Pipeline.DisplayName,
		Description:           m.Pipeline.Description,
		NodeFilter:            m.Pipeline.NodeFilter,
		SimilarityScoreCutoff: m.Pipeline.SimilarityScoreCutoff,
	})
	if err != nil {
		return "", fmt.Errorf("create pipeline %q: %w", m.Pipeline.Name, err)
	}
	log.Printf("created entity-matching pipeline %q: %s", m.Pipeline.Name, created.ID)
	return created.ID, nil
}

func (m *manifest) ensureAuditSigning(ctx context.Context, admin *config.AdminClient) (string, error) {
	existing, err := admin.AuditSignings().Read(ctx, m.AuditSigning.Name, config.WithLocation(m.ProjectID))
	if err == nil {
		log.Printf("audit signing %q exists: %s", m.AuditSigning.Name, existing.ID)
		return existing.ID, nil
	}
	if apiErr, ok := transport.AsAPIError(err); !ok || !apiErr.IsNotFound() {
		return "", fmt.Errorf("read audit signing %q: %w", m.AuditSigning.Name, err)
	}
	created, err := admin.AuditSignings().Create(ctx, &config.CreateAuditSigning{
		ProjectID:   m.ProjectID,
		Name:        m.AuditSigning.Name,
		DisplayName: m.AuditSigning.DisplayName,
		Description: m.AuditSigning.Description,
		AuditSigningConfig: config.AuditSigningConfig{
			Provider: m.AuditSigning.Provider,
		},
	})
	if err != nil {
		return "", fmt.Errorf("create audit signing %q: %w", m.AuditSigning.Name, err)
	}
	log.Printf("created audit signing %q: %s", m.AuditSigning.Name, created.ID)
	return created.ID, nil
}

func (m *manifest) apply(ctx context.Context, admin *config.AdminClient, cli *indykite.Client) error {
	var err error
	if m.KBACPolicyID, err = m.ensurePolicy(ctx, admin, &m.KBACPolicy); err != nil {
		return err
	}
	if m.CIQPolicyID, err = m.ensurePolicy(ctx, admin, &m.CIQPolicy); err != nil {
		return err
	}
	if m.KQID, err = m.ensureQuery(ctx, admin, m.CIQPolicyID); err != nil {
		return err
	}
	if m.PipelineID, err = m.ensurePipeline(ctx, admin); err != nil {
		return err
	}
	if m.AuditSigningID, err = m.ensureAuditSigning(ctx, admin); err != nil {
		return err
	}

	if _, err = cli.Capture().UpsertNodesChunked(ctx, m.loadedDataset.Nodes, 0); err != nil {
		return fmt.Errorf("ingest nodes: %w", err)
	}
	if _, err = cli.Capture().UpsertRelationships(ctx, m.loadedDataset.Relationships...); err != nil {
		return fmt.Errorf("ingest relationships: %w", err)
	}
	log.Printf("ingested %d nodes, %d relationships",
		len(m.loadedDataset.Nodes), len(m.loadedDataset.Relationships))
	return nil
}

// resolve looks up the gids of already-provisioned resources.
func (m *manifest) resolve(ctx context.Context, admin *config.AdminClient) error {
	kq, err := admin.KnowledgeQueries().Read(ctx, m.Query.Name, config.WithLocation(m.ProjectID))
	if err != nil {
		return fmt.Errorf("knowledge query %q not found; run apply first: %w", m.Query.Name, err)
	}
	m.KQID = kq.ID
	pipe, err := admin.EntityMatchingPipelines().Read(ctx, m.Pipeline.Name, config.WithLocation(m.ProjectID))
	if err != nil {
		return fmt.Errorf("pipeline %q not found; run apply first: %w", m.Pipeline.Name, err)
	}
	m.PipelineID = pipe.ID
	signing, err := admin.AuditSignings().Read(ctx, m.AuditSigning.Name, config.WithLocation(m.ProjectID))
	if err != nil {
		return fmt.Errorf("audit signing %q not found; run apply first: %w", m.AuditSigning.Name, err)
	}
	m.AuditSigningID = signing.ID
	return nil
}

// printEnv emits the export block the integration tests consume.
func (m *manifest) printEnv() {
	fmt.Println("# fixture environment for the SDK integration tests")
	for _, kv := range [][2]string{
		{"CIQ_QUERY_ID", m.KQID},
		{"EM_PIPELINE_ID", m.PipelineID},
		{"AUDIT_SIGNING_ID", m.AuditSigningID},
	} {
		fmt.Printf("export %s=%q\n", kv[0], kv[1])
	}
	for k, v := range m.Env {
		fmt.Printf("export %s=%q\n", k, v)
	}
}

func (m *manifest) destroy(ctx context.Context, admin *config.AdminClient, cli *indykite.Client) error {
	if len(m.loadedDataset.Relationships) > 0 {
		if _, err := cli.Capture().DeleteRelationships(ctx, m.loadedDataset.Relationships...); err != nil {
			log.Printf("delete relationships: %v", err)
		}
	}
	nodes := make([]capture.Node, 0, len(m.loadedDataset.Nodes))
	for _, n := range m.loadedDataset.Nodes {
		nodes = append(nodes, n.Node)
	}
	if len(nodes) > 0 {
		if _, err := cli.Capture().DeleteNodes(ctx, nodes...); err != nil {
			log.Printf("delete nodes: %v", err)
		}
	}

	deleteByName := []struct {
		del  func(ctx context.Context) error
		name string
	}{
		{name: m.Query.Name, del: func(ctx context.Context) error {
			kq, err := admin.KnowledgeQueries().Read(ctx, m.Query.Name, config.WithLocation(m.ProjectID))
			if err != nil {
				return err
			}
			return admin.KnowledgeQueries().Delete(ctx, kq.ID, "")
		}},
		{name: m.CIQPolicy.Name, del: func(ctx context.Context) error {
			pol, err := admin.AuthorizationPolicies().Read(ctx, m.CIQPolicy.Name, config.WithLocation(m.ProjectID))
			if err != nil {
				return err
			}
			return admin.AuthorizationPolicies().Delete(ctx, pol.ID, "")
		}},
		{name: m.KBACPolicy.Name, del: func(ctx context.Context) error {
			pol, err := admin.AuthorizationPolicies().Read(ctx, m.KBACPolicy.Name, config.WithLocation(m.ProjectID))
			if err != nil {
				return err
			}
			return admin.AuthorizationPolicies().Delete(ctx, pol.ID, "")
		}},
		{name: m.Pipeline.Name, del: func(ctx context.Context) error {
			pipe, err := admin.EntityMatchingPipelines().Read(ctx, m.Pipeline.Name, config.WithLocation(m.ProjectID))
			if err != nil {
				return err
			}
			return admin.EntityMatchingPipelines().Delete(ctx, pipe.ID, "")
		}},
		{name: m.AuditSigning.Name, del: func(ctx context.Context) error {
			signing, err := admin.AuditSignings().Read(ctx, m.AuditSigning.Name, config.WithLocation(m.ProjectID))
			if err != nil {
				return err
			}
			return admin.AuditSignings().Delete(ctx, signing.ID, "")
		}},
	}
	for _, d := range deleteByName {
		if err := d.del(ctx); err != nil {
			if apiErr, ok := transport.AsAPIError(err); ok && apiErr.IsNotFound() {
				log.Printf("%q already absent", d.name)
				continue
			}
			return fmt.Errorf("delete %q: %w", d.name, err)
		}
		log.Printf("deleted %q", d.name)
	}
	return nil
}
