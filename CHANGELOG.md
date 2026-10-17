# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.0, 2026-10-17

- Added `latitude`, a REST client generated with oapi-codegen from the Latitude API specification version 1.0.0 (OpenAPI 3.1.0), published at https://raw.githubusercontent.com/latitude-dev/latitude-llm/main/apps/api/openapi.json. It covers all 48 operations across 35 paths: projects, traces and spans, scores, annotations, saved searches, issues, API keys, OAuth keys, account, members, health and OAuth resource metadata.
- Number fields use `float64`. Authentication uses organization-scoped Bearer API keys through a request editor.
- `make generate` downloads the spec into the gitignored `api/` directory, converts nullable types, nullable enum members and exclusive numeric bounds to OpenAPI 3.0 form, normalizes description punctuation, and regenerates the client.
- Added HTTP contract tests, the compiled README example, vendored dependencies, development tooling, CI workflows and the `go-latitude-llm` agent skill.
