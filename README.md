# IndyKite Client Library for Go

<div align="left">
<a href="https://indykite.com">
<img src="https://raw.githubusercontent.com/indykite/.github/master/assets/squareformatlogo.png" alt="IndyKite Logo" width="100px" height="100px" align="right">
</a>
</div>

As AI agents and users execute work, IndyKite continuously captures context, relationships, trust signals, and actions.
This becomes the point of decision, pulling signals from multiple systems and applying them at runtime, with full traceability.
The result is agentic AI that can operate across platforms with precision and deliver an entirely new class of intelligent services.

This repository contains the Go client SDK for the [IndyKite Platform](https://indykite.com) **REST APIs**
([OpenAPI reference](https://openapi.indykite.com)).

[![codecov](https://codecov.io/gh/indykite/indykite-sdk-go/branch/master/graph/badge.svg?token=TFCDLXbnsh)](https://codecov.io/gh/indykite/indykite-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/indykite/indykite-sdk-go)](https://goreportcard.com/report/github.com/indykite/indykite-sdk-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/indykite/indykite-sdk-go.svg)](https://pkg.go.dev/github.com/indykite/indykite-sdk-go)

## Requirements

- Go 1.26+
- An **App Agent** credential for the runtime APIs (AuthZEN, ContX IQ, Capture, Entity Matching)
- A **Service Account** credential for the Config API

Credentials can be obtained from the Hub ([EU](https://eu.hub.indykite.com/) or [US](https://us.hub.indykite.com/))
or from your IndyKite contact.

## Two planes, two clients

| Plane | Credential | Entry point | Services |
| --- | --- | --- | --- |
| Runtime / data | App Agent | `indykite.NewClient[FromEnv]` | AuthZEN, ContX IQ, Capture, Entity Matching |
| Control | Service Account | `indykite.NewAdmin[FromEnv]` | Config management (`/configs/v1`) |

The SDK picks the right auth header per plane automatically (`X-IK-ClientKey` for runtime,
`Authorization: Bearer` for control). No TLS or connection setup is needed.

The two credential artifacts differ: `INDYKITE_APPLICATION_CREDENTIALS[_FILE]` holds the
**App Agent credential token itself** (the raw `X-IK-ClientKey` value — always pass
`WithRegion`/`WithBaseURL`, the token carries no URL); `INDYKITE_SERVICE_ACCOUNT_CREDENTIALS[_FILE]`
holds the **Service Account JSON artifact** (`serviceAccountId`, `baseUrl`, a pre-issued `token`,
private key). The SDK sends the pre-issued token while valid and mints a fresh JWT from the key
when it expires.

## Quick start

```go
import (
    indykite "github.com/indykite/indykite-sdk-go"
    "github.com/indykite/indykite-sdk-go/authzen"
    "github.com/indykite/indykite-sdk-go/ciq"
    "github.com/indykite/indykite-sdk-go/config"
)

// Runtime plane — reads INDYKITE_APPLICATION_CREDENTIALS[_FILE].
cli, err := indykite.NewClientFromEnv(ctx, indykite.WithRegion("eu"))

ok, err := cli.AuthZEN().Allowed(ctx,
    authzen.NewNode("Person", "ada"), "PROVISION", authzen.NewNode("Server", "gpu-7"))

rows, err := cli.CIQ().All(ctx, ciq.ExecuteRequest{
    ID: "get-servers", InputParams: map[string]any{"region": "eu"},
})

// Control plane — reads INDYKITE_SERVICE_ACCOUNT_CREDENTIALS[_FILE].
admin, err := indykite.NewAdminFromEnv(ctx, indykite.WithRegion("eu"))
pol, err := admin.AuthorizationPolicies().Create(ctx, &config.CreateAuthorizationPolicy{
    ProjectID: projectID, Name: "can-provision", Policy: policyJSON, Status: config.StatusActive,
})
```

Config mutations are guarded by ETags for optimistic concurrency: `Read` captures the ETag,
`Update`/`Delete` send it as `If-Match`.

Errors are `*transport.APIError` (HTTP status, platform code/message, request id):

```go
if apiErr, ok := transport.AsAPIError(err); ok && apiErr.IsNotFound() { /* ... */ }
```

Options: `WithRegion`, `WithBaseURL`, `WithRetry`, `WithTracing` (OpenTelemetry),
`WithHTTPClient`, `WithUserAgent`.

## Examples

Runnable per-service CLIs live in [examples/](examples/) — one subcommand per operation.

## Tests

```sh
make test          # unit tests (offline, httptest-based)
make integration   # integration tests against the platform
```

Integration tests cover the Config, Capture, AuthZEN, ContX IQ, and Entity Matching APIs.
`TestIntegrationIKGEndToEnd` additionally exercises an actual IKG with zero pre-provisioned
fixtures: it seeds a KBAC policy plus a CIQ read policy/knowledge query, ingests a
`Person -[:OWNS]-> Server` graph, asserts the AuthZEN decision (positive and negative) and
the knowledge-query read over that live data, then removes everything it created.
They authenticate from `INDYKITE_APPLICATION_CREDENTIALS[_FILE]` /
`INDYKITE_SERVICE_ACCOUNT_CREDENTIALS[_FILE]` and scope to `PROJECT_ID` (and optionally
`ORGANIZATION_ID`). Fixture-dependent tests read `CIQ_QUERY_ID`, `EM_PIPELINE_ID`, and
`AUTHZEN_{SUBJECT_TYPE,SUBJECT_ID,ACTION,RESOURCE_TYPE,RESOURCE_ID}`; each test skips
cleanly when its inputs are unset. `INDYKITE_BASE_URL` overrides the default region URL.
The fixture dataset and config resources behind those env values live in [test/fixtures/](test/fixtures/)
and are provisioned idempotently with `go run ./test/setup apply` (see [test/README.md](test/README.md)).
The data-schema API is feature-flagged and intentionally has no integration tests.

When `SDK_AUDIT_TABLE_NAME` is set, the AuthZEN and ContX IQ integration tests additionally
verify that the platform's audit events reached the BigQuery audit-log table (authenticating
with Application Default Credentials; `BQ_PROJECT_ID` selects the hosting GCP project).

## SDK Documentation

- [IndyKite documentation](https://docs.indykite.com)
- [REST API reference](https://openapi.indykite.com)
- [Go package reference](https://pkg.go.dev/github.com/indykite/indykite-sdk-go)
