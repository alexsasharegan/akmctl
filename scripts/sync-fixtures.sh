#!/bin/sh
# Copy the API fixtures from their source, docs/spec/akm-api/.fixtures, to the copy the Go
# tests read, testdata/api and testdata/yaml-lite.json (work item 261029-x4ws).
# Edit the source, run this, commit both. TestFixtureCopiesAgree fails when they differ.
# make_fixtures.py stays in the source only.
set -eu
root=$(cd "$(dirname "$0")/.." && pwd)
src="$root/docs/spec/akm-api/.fixtures"
dst="$root/testdata"
[ -d "$src/cases" ] && [ -d "$src/repos" ] || { echo "sync-fixtures: $src has no cases/ and repos/" >&2; exit 1; }
rm -rf "$dst/api/cases" "$dst/api/repos"
mkdir -p "$dst/api"
cp -R "$src/cases" "$dst/api/cases"
cp -R "$src/repos" "$dst/api/repos"
cp "$src/yaml-lite.json" "$dst/yaml-lite.json"
echo "synced $src -> $dst"
