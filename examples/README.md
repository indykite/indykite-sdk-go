# IndyKite REST SDK examples

Small CLIs demonstrating every service of the REST SDK. Each example is a
standalone `main` package with one subcommand per operation, mirroring the
structure of the platform APIs.

| Example | Plane | Credential env var | Demonstrates |
| --- | --- | --- | --- |
| [`authzen`](authzen/) | runtime | `INDYKITE_APPLICATION_CREDENTIALS[_FILE]` | evaluate (optionally with end-user/delegated tokens), batch evaluate, search action/resource/subject, list policies |
| [`capture`](capture/) | runtime | `INDYKITE_APPLICATION_CREDENTIALS[_FILE]` | upsert/delete nodes & relationships, property deletes, chunked batches |
| [`ciq`](ciq/) | runtime | `INDYKITE_APPLICATION_CREDENTIALS[_FILE]` | execute a ContX IQ query, paginate all records |
| [`tokenclaims`](tokenclaims/) | both | both | KBAC and CIQ policies reading `$token` / `$ik_token` claims: setup, decisions with and without the delegated token, teardown |
| [`audit`](audit/) | runtime | `INDYKITE_APPLICATION_CREDENTIALS[_FILE]` | page through signed audit logs, manifests and checkpoints, fetch the verification keys |
| [`entitymatching`](entitymatching/) | runtime | `INDYKITE_APPLICATION_CREDENTIALS[_FILE]` | run pipeline, read status, suggested property mappings |
| [`config`](config/) | control | `INDYKITE_SERVICE_ACCOUNT_CREDENTIALS[_FILE]` | organization, projects, app agents, credentials, policies, knowledge queries, event sinks, and other config resources |

Run any example with no arguments to see its subcommands, e.g.:

```sh
export INDYKITE_APPLICATION_CREDENTIALS_FILE=~/app-agent.json
go run ./examples/authzen evaluate -subject-type Person -subject-id ada \
    -action PROVISION -resource-type Server -resource-id gpu-node-7
```

All examples honor `INDYKITE_BASE_URL` (explicit base URL, e.g. a staging
gateway) and default to the `eu` region otherwise.

The `tokenclaims` example needs both credentials plus an end-user token and an
IndyKite delegated token (from the Token Service) for the same user:

```sh
eval "$(go run ./examples/tokenclaims setup -project-id "$PROJECT_ID")"   # exports CLAIMS_*
export END_USER_TOKEN=... IK_TOKEN=...
go run ./examples/tokenclaims evaluate -subject-id karel -resource-id docA
go run ./examples/tokenclaims execute
go run ./examples/tokenclaims teardown
```
