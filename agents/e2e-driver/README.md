# e2e-driver (SDK-backed end-to-end driver)

A Go program that drives a full **seed → ingest → assert → cleanup** scenario
against the IndyKite platform using the SDK. It is the REST-native counterpart to
the `agent-gateway-e2e` Bruno collection's platform-side setup and verification —
a typed Go alternative to the curl/Bruno steps.

It exercises **both planes**:

| Stage | Plane | SDK |
| --- | --- | --- |
| seed a KBAC authorization policy | control | `config.AdminClient.AuthorizationPolicies().Create` |
| ingest the subject/resource graph | runtime | `cli.Capture().UpsertNodes` / `UpsertRelationships` |
| assert the AuthZEN decision | runtime | `cli.AuthZEN().Allowed` |
| assert a ContX IQ read | runtime | `cli.CIQ().All` |
| cleanup (always runs) | control | `config.AdminClient.AuthorizationPolicies().Delete` |

`Driver.Run` returns one `Step` per stage and **always attempts cleanup** (even
when an assertion fails), so a failed run still deletes the policy it created.

## Scope

This drives the platform directly over REST. It does **not** exercise the
gateway's agent-to-agent (A2A) hops — that's the axis the agent-gateway owns and
which the Bruno `12 IAG Tests` collection covers by calling through the gateway.
This driver is the platform-side setup + verification the SDK can own.

## Run

```bash
export INDYKITE_SERVICE_ACCOUNT_CREDENTIALS='{"serviceAccountId":"...","token":"...","privateKeyJWK":{...}}'  # control plane (JSON artifact)
export INDYKITE_APPLICATION_CREDENTIALS='<app-agent-credential-token>'                                        # runtime plane (the raw token)
export INDYKITE_REGION=eu                 # or INDYKITE_BASE_URL=https://api.eu.indykite.com
export INDYKITE_PROJECT_ID=gid:...        # where the policy is created
export CIQ_QUERY_ID=get-doc               # optional: a knowledge query to read

go run ./agents/e2e-driver
```

It prints a PASS/FAIL line per step and exits non-zero if any step failed —
suitable as a CI smoke test alongside the gateway e2e suite.
