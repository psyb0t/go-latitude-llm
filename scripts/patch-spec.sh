#!/bin/bash
# Downconverts Latitude's OpenAPI 3.1 schemas for oapi-codegen's 3.0 reader.
# Nullable type lists become scalar types plus nullable: true. Null enum
# entries become nullable instead of generating an invalid Go enum value.
# Numeric exclusive bounds become minimum/maximum plus exclusive flags.
# No operations or schema properties are removed.
set -euo pipefail
trap 'on_error "$?" "$LINENO"' ERR

LOG_FILE="${LOG_FILE:-/tmp/$(basename "$0" .sh).log}"
exec > >(tee -a "$LOG_FILE") 2>&1

SPEC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/api"
readonly SPEC_DIR
readonly SPEC_IN="${SPEC_DIR}/openapi.json"
readonly SPEC_OUT="${SPEC_DIR}/openapi.patched.json"

on_error() {
	local exit_code="$1"
	local line="$2"

	printf '{"level":"ERROR","file":"%s","line":%d,"msg":"command failed","exit_code":%d}\n' \
		"${BASH_SOURCE[1]}" "$line" "$exit_code" >&2
}

log() {
	local level="$1"
	shift

	jq -nc \
		--arg time "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--arg level "$level" \
		--arg file "${BASH_SOURCE[1]##*/}" \
		--arg line "${BASH_LINENO[0]}" \
		--arg func "${FUNCNAME[1]:-main}" \
		--arg msg "$*" \
		'{time:$time,level:$level,file:$file,line:($line|tonumber),func:$func,msg:$msg}' >&2
}

main() {
	local tmp
	tmp="$(mktemp "${SPEC_DIR}/.openapi.patched.XXXXXX")"
	# shellcheck disable=SC2064 # expand tmp now, the local is gone when EXIT fires
	trap "rm -f -- '${tmp}'" EXIT

	jq '
	  walk(if type == "object" then
	    (if (.type | type) == "array" then
	      if (.type | length) == 2 and (.type | index("null")) != null then
	        .type = (.type | map(select(. != "null")) | .[0])
	        | .nullable = true
	      else error("unsupported schema type union") end
	    else . end)
	    | (if has("enum") and (.enum | index(null)) != null then
	        .enum |= map(select(. != null)) | .nullable = true
	      else . end)
	    | (if (.exclusiveMinimum | type) == "number" then
	        .minimum = .exclusiveMinimum | .exclusiveMinimum = true
	      else . end)
	    | (if (.exclusiveMaximum | type) == "number" then
	        .maximum = .exclusiveMaximum | .exclusiveMaximum = true
	      else . end)
	  else . end)
	  | .openapi = "3.0.3"
	' "$SPEC_IN" >"$tmp"

	chmod 0644 "$tmp"
	mv -- "$tmp" "$SPEC_OUT"
	log INFO "downconverted OpenAPI 3.1 schemas out=${SPEC_OUT##*/}"
}

main "$@"
