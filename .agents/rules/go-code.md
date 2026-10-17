# Go code

- **GO1** `client.gen.go` is generated. Never edit it by hand. Change `config.yaml` or the spec patch in `scripts/`, then run `make generate` and commit the result.
- **GO2** Spec fixes live in `scripts/`, applied before generation. `scripts/patch-spec.sh` downconverts nullable type arrays, null enum members and exclusive numeric bounds from OpenAPI 3.1 to 3.0 form. Each patch explains its purpose in a comment.
- **GO3** The downloaded spec in `api/` is gitignored and never committed.
- **GO4** Number fields stay `float64` through the `type-mapping` in `config.yaml`. Scores lose precision in `float32`.
- **GO5** Hand-written Go (`gen.go` and the tests) wraps errors with `ctxerrors.Wrap` against `ctxerrors/commerr` sentinels. No `fmt.Errorf`, no `errors.New` outside a sentinel declaration.
- **GO6** The exported API is whatever oapi-codegen emits. A config change that renames or removes an exported name is a breaking release (see `releases.md`).
- **GO7** `make lint` passes with zero issues. Hand-written lines are at most 80 columns, the `.golangci.yml` limit.
