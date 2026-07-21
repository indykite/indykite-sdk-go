# IndyKite REST SDK examples

Small CLIs demonstrating every service of the REST SDK. Each example is a
standalone `main` package with one subcommand per operation, mirroring the
structure of the platform APIs.

| Example | Plane | Credential env var | Demonstrates |
| --- | --- | --- | --- |
| [`authzen`](authzen/) | runtime | `INDYKITE_APPLICATION_CREDENTIALS[_FILE]` | evaluate, batch evaluate, search action/resource/subject |
| [`capture`](capture/) | runtime | `INDYKITE_APPLICATION_CREDENTIALS[_FILE]` | upsert/delete nodes & relationships, property deletes, chunked batches |
| [`ciq`](ciq/) | runtime | `INDYKITE_APPLICATION_CREDENTIALS[_FILE]` | execute a ContX IQ query, paginate all records |
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
