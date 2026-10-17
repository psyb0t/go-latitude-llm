---
name: go-latitude-llm
description: "Use the unofficial Latitude LLM Go client when Go code needs the latitude.so ops platform API or already imports github.com/psyb0t/go-latitude-llm. Covers projects, traces, scores, annotations, saved searches, issues, API keys and organization management."
metadata:
  homepage: https://github.com/psyb0t/go-latitude-llm
  permissions:
    filesystem: "Guides edits to Go source in the current project. The client performs no filesystem operations."
    shell: "The Go toolchain to add the module, compile examples and run tests."
    network: "Module downloads through the Go proxy. Runtime requests to https://api.latitude.so, or a trusted self-hosted API URL explicitly selected by the user."
  openclaw:
    requires:
      bins:
        - go
      env:
        - LATITUDE_API_KEY
---

# go-latitude-llm

Unofficial, typed Go client for the Latitude LLM ops platform at latitude.so. The package is `latitude` at the module root. This is not a client for Latitude.sh.

Read [references/setup.md](references/setup.md) before writing calls. It covers installation, auth, the compiled usage example, methods, pagination and response handling.

## Security and safety

- Read an organization API key from `LATITUDE_API_KEY`. Use `WithRequestEditorFn` to set `Authorization: Bearer <key>`. Never hardcode, print or log credentials or request headers.
- Use `https://api.latitude.so` as the base URL. Generated paths already contain `/v1`. Set a custom HTTP timeout through `WithHTTPClient` when needed.
- The operation descriptions require user OAuth tokens for `CreateSavedSearch`, `InviteMember`, `UpdateMember` and `RemoveMember`, despite their `ApiKeyAuth` security declarations. Token acquisition is outside this client. API-key account responses have null user and role fields.
- Creating scores, annotations, searches and projects, changing issues, revoking keys, inviting members and deleting resources affect the organization. Perform those actions only within the user's request. Exports can send data by email.
- Keep trace contents and member information out of logs. Only send credentials to the user's trusted API URL.

## When to use

- Go code needs to inspect Latitude projects, traces, spans, scores or issues.
- Go code needs an authorized write or organization-management operation.
- The project imports `github.com/psyb0t/go-latitude-llm`.

## When not to use

- The code is not Go, or the service is Latitude.sh.
- The user needs a client maintained by Latitude. This package is unofficial.
- The task changes the generated client itself. Follow the repository's `AGENTS.md` instead.

## Working with the client

Use `NewClientWithResponses` and the operation's `...WithResponse` method. Check the Go error, HTTP status and expected success field before reading the result. `JSON200` and `JSON201` are pointers; non-success statuses do not automatically produce Go errors. Nullable model fields are pointers too.

Method names come from the spec's operation IDs. Scores use `float64`, and union schemas expose generated `From...` and `As...` helpers. Inspect the generated signature before constructing a union request. `ListProjects` takes no params and currently returns a single page. Other collections carry operation-specific cursor fields.

Finish by compiling the calling code and testing its requests against a local HTTP server. Do not call the live API merely to validate an example.
