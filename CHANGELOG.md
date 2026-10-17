# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.2, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.
- The skill's `homepage` and `permissions` sit at the top level of the frontmatter, where agents read them, and the skill is user-invocable.

## v0.1.1, 2026-10-17

- Generated doc comments keep the API's own description text. The spec patch no longer rewrites dashes or wording in descriptions.
- `make sec` uses Semgrep 1.181.0. Its pinned requirements now hash every published file, so the development image also builds on arm64.

## v0.1.0, 2026-10-17

- Added `latitude`, a REST client generated with oapi-codegen from the Latitude API specification version 1.0.0 (OpenAPI 3.1.0), published at https://raw.githubusercontent.com/latitude-dev/latitude-llm/main/apps/api/openapi.json. It covers all 48 operations across 35 paths: projects, traces and spans, scores, annotations, saved searches, issues, API keys, OAuth keys, account, members, health and OAuth resource metadata.
- Number fields use `float64`. Authentication uses organization-scoped Bearer API keys through a request editor.
- `make generate` downloads the spec into the gitignored `api/` directory, converts nullable types, nullable enum members and exclusive numeric bounds to OpenAPI 3.0 form, normalizes description punctuation, and regenerates the client.
- Added HTTP contract tests, the compiled README example, vendored dependencies, development tooling, CI workflows and the `go-latitude-llm` agent skill.
