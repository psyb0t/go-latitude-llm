# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.4, 2026-10-28

- Updated Semgrep in the development image to 1.182.1 and refreshed its locked dependencies. certifi 2026.10.23 drops the emSign root CAs whose intermediate keys were compromised (GHSA-6q48-6c7w-9gw2).
- Dependabot no longer opens pull requests for `mcp` or the OpenTelemetry packages in `requirements-semgrep.txt`. Semgrep pins their versions itself, so bumping one of them on its own cannot install.
- Updated kin-openapi, which the code generator uses, from 0.142.0 to 0.149.0. The newer version fixes several panics on malformed specs reported in GitHub security advisories.
- The security scan merges its two gitleaks runs into one SARIF run. GitHub code scanning rejects a report with two runs from the same tool, so no scan results had been uploaded.

## v0.1.3, 2026-10-27

- The development image is based on Debian (`golang:1.26.10-trixie`) instead of Alpine. Alpine keeps only the latest build of each package, so the exact pins broke the image build whenever a package was updated. System packages now come from the snapshot.debian.org date the base image was built from, which keeps them fixed for a given base image digest. jq 1.8.2, shellcheck 0.11.0 and shfmt 3.13.1 keep their versions, and Semgrep runs in its own virtual environment.
- Raised the minimum Go version in `go.mod` from 1.26.9 to 1.26.10, which fixes the `crypto/x509` vulnerability GO-2026-6809 found by govulncheck.
- The issue-pull workflow runs at hours derived from the repository name instead of the same four hours as every other repository.

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
