#!/usr/bin/env bash
# Updates the Go versions this module supports.
#
# Usage: update-go-versions.sh <latest> <penultimate>
#
# Called by the shared sdk-go-versions workflow, which passes versions it has already
# validated against a strict version pattern.

set -euo pipefail

latest="${1:?latest version required}"
penultimate="${2:?penultimate version required}"

versions_file=./.github/variables/go-versions.env

# Read the outgoing penultimate before rewriting it; the README references it by value.
previous_penultimate=$(grep -E '^penultimate=' "$versions_file" | head -1 | cut -d= -f2- | tr -d '[:space:]')

sed -i -e "s#^latest=.*#latest=$latest#" \
       -e "s#^penultimate=.*#penultimate=$penultimate#" \
       "$versions_file"

# The README states the minimum supported Go version, which tracks the penultimate release.
if [ -n "$previous_penultimate" ]; then
  sed -i "s/Go version $previous_penultimate/Go version $penultimate/g" README.md
fi

go mod edit -go="$penultimate"
