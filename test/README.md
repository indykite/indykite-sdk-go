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
  execution needs only the App Agent key), an entity-matching pipeline
  (Person → Customer), and a `PLATFORM_MANAGED` audit-signing config (no
  customer key material needed). Plus the static env values the tests read.
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

It sets `CIQ_QUERY_ID`, `EM_PIPELINE_ID` and `AUDIT_SIGNING_ID` (the created
gids) plus the static
`AUTHZEN_{SUBJECT_TYPE,SUBJECT_ID,ACTION,RESOURCE_TYPE,RESOURCE_ID}` and
`CIQ_INPUT_PARAMS` values from the manifest.

The fixture audit-signing config is `PLATFORM_MANAGED` on purpose: a
customer-managed config pointing at fake key material would be accepted by the
config API but could break audit signing for the project. To exercise the
customer-managed path set `AUDIT_SIGNING_KEY_RESOURCE` and `AUDIT_SIGNING_KID`
(plus optional `AUDIT_SIGNING_PROVIDER`, default `CUSTOMER_GCP_KMS`, and
`AUDIT_SIGNING_AUTH_PARAMS` as a JSON object) to real material; the
`TestIntegrationConfigAuditSigningCustomerManaged` test skips otherwise.

CI needs no manual copies: both the Integration job (go-tests.yaml) and the
go-sdk-tests pipeline image (docker/infra/startscript.sh) run `setup apply`
themselves before the tests and adopt the env it prints.

Idempotency: config resources are looked up by their stable names
(`sdk-it-fix-*`) and only created when absent; graph upserts are keyed by
external id. Re-running `apply` after a partial failure is always safe.
