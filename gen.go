// Package latitude is a Go client for the Latitude LLM API,
// generated from the public OpenAPI document.
package latitude

//go:generate bash scripts/fetch-spec.sh
//go:generate bash scripts/patch-spec.sh
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.patched.json
