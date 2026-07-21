# Integration-test fixtures

Everything the integration tests need in a real project, as data:

- [`fixtures/dataset.json`](fixtures/dataset.json) — the seed graph (Capture API
  payload shapes): three `Person` identities, two `Customer` records sharing
  emails with the persons (entity-matching material), three `Server` nodes with
  `region`/`price` properties, and the `OWNS` relationships. All external ids
  are prefixed `sdk-it-fix-`.
- [`fixtures/config.json`](fixtures/config.json) — the Config-API resources:
  a KBAC policy (*Person may `SDK_IT_CAN_USE` a Server they OWN*), a CIQ read
  policy + knowledge query (*Servers by `$region`*, `_Application` subject so
  execution needs only the App Agent key), and an entity-matching pipeline
  (Person → Customer). Plus the static env values the tests read.
- [`setup/`](setup/) — the SDK-backed tool that applies them.

## Usage

```sh
export INDYKITE_SERVICE_ACCOUNT_CREDENTIALS_FILE=~/sa.json
export INDYKITE_APPLICATION_CREDENTIALS_FILE=~/app-agent.json
export PROJECT_ID=gid:...

go run ./test/setup apply     # idempotent: create-if-missing + upsert dataset
go run ./test/setup env       # print the export block for existing fixtures
go run ./test/setup destroy   # remove the dataset and config resources
```

`apply`/`env` end with an export block — evaluate it before running the tests:

```sh
eval "$(go run ./test/setup env)"
make integration
```

It sets `CIQ_QUERY_ID` and `EM_PIPELINE_ID` (the created gids) plus the static
`AUTHZEN_{SUBJECT_TYPE,SUBJECT_ID,ACTION,RESOURCE_TYPE,RESOURCE_ID}` and
`CIQ_INPUT_PARAMS` values from the manifest.

CI needs no manual copies: both the Integration job (go-tests.yaml) and the
go-sdk-tests pipeline image (docker/infra/startscript.sh) run `setup apply`
themselves before the tests and adopt the env it prints.

Idempotency: config resources are looked up by their stable names
(`sdk-it-fix-*`) and only created when absent; graph upserts are keyed by
external id. Re-running `apply` after a partial failure is always safe.
