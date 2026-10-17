# go-latitude-llm

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-latitude-llm.svg)](https://pkg.go.dev/github.com/psyb0t/go-latitude-llm)
[![CI](https://github.com/psyb0t/go-latitude-llm/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-latitude-llm/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-latitude-llm/badges/version.svg)](https://github.com/psyb0t/go-latitude-llm/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-latitude-llm/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-latitude-llm/badges/importers.svg)](https://github.com/psyb0t/go-latitude-llm/blob/badges/importers.md)

Unofficial Go client for the [Latitude LLM API](https://latitude.so), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from Latitude's public OpenAPI document. I am not affiliated with Latitude.

The package is `latitude`, in the repository root. Everything in `client.gen.go` is generated. Don't edit it by hand, regenerate it. This client covers the ops platform at latitude.so.

## Install

```sh
go get github.com/psyb0t/go-latitude-llm
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`: projects, traces and spans, annotations and scores, saved searches, issue management, API keys, OAuth keys, account and members. Responses decode into `JSON200`, `JSON201` and the other status fields.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-latitude-llm"
)

func main() {
	auth := latitude.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(
				"Authorization",
				"Bearer "+os.Getenv("LATITUDE_API_KEY"),
			)

			return nil
		},
	)

	client, err := latitude.NewClientWithResponses(
		"https://api.latitude.so",
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.ListProjectsWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

The same code is in `example_test.go`, so `make test` checks that it compiles. It needs a real key to run. Tests use local HTTP servers and never call Latitude.

The full API, with the spec's descriptions as doc comments, is in the [Go reference](https://pkg.go.dev/github.com/psyb0t/go-latitude-llm) and in `go doc github.com/psyb0t/go-latitude-llm`.

Check `StatusCode()` and the relevant JSON field before using a response. A non-success HTTP status does not by itself produce a Go error. Transport and JSON decoding failures do. Nullable fields are pointers; API-key account responses have null `User` and `Role`.

Method names come from operation IDs, such as `ListProjectsWithResponse`, `CreateProjectWithResponse`, `ListTracesWithResponse` and `GetAccountWithResponse`. The spec says the project list fits in one page, and `ListProjects` takes no cursor parameter. Other paginated collections expose `HasMore` and `NextCursor`; pass the returned cursor using that operation's params or body until `HasMore` is false.

## Auth

The spec's `ApiKeyAuth` scheme is HTTP Bearer authentication with an organization-scoped API key. Set `LATITUDE_API_KEY` and attach `Authorization: Bearer <key>` through `WithRequestEditorFn`, as shown above. The client does not add credentials automatically. Keep keys out of source and logs.

Use `https://api.latitude.so` as the base URL. The generated paths already contain `/v1`; do not append it to the base URL. For self-hosted installations, pass that installation's trusted API URL. `WithHTTPClient` accepts a custom HTTP client for timeouts and transport settings. The generated client does not retry requests automatically.

The spec's descriptions require OAuth for creating saved searches and inviting, updating or removing members, even though their security declarations also name `ApiKeyAuth`. Those operations need a user OAuth Bearer token. API keys do not represent a user. OAuth token acquisition is outside this client.

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` downloads `https://raw.githubusercontent.com/latitude-dev/latitude-llm/main/apps/api/openapi.json` to `api/openapi.json`.
- `scripts/patch-spec.sh` writes `api/openapi.patched.json`. It converts OpenAPI 3.1 nullable type arrays into scalar types with `nullable: true`, removes null enum entries while preserving nullability, and converts numeric exclusive bounds to 3.0 bounds and flags. It marks the generator input as OpenAPI 3.0.3. No operations are removed.
- oapi-codegen runs with `config.yaml`, which maps `number` to `float64` instead of the default `float32`.

`.github/workflows/regen.yml` does the same every day in CI. When the spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude can repair the code and repeat the checks up to three times, followed by review. The workflow protects tooling and rule paths from these repairs. A passing change gets a CHANGELOG entry and a release tag when the workflow is installed in a published repository.

## Agent integrations

The [skill](.agents/skills/go-latitude-llm) teaches a coding agent how to use go-latitude-llm when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-latitude-llm@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-latitude-llm@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The release workflow publishes the skill to ClawHub on tag pushes:

```bash
openclaw skills install @psyb0t/go-latitude-llm
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). The generated code comes from Latitude's published API specification. The specification itself is not redistributed here. The spec does not declare a license or terms-of-service URL; this repository's license does not grant rights to the provider's specification or service.

Vendored dependencies retain their own licenses. The development tools include GPL-licensed code, including golangci-lint; these tools are separate from the client runtime.
