---
id: "261108-7szk"
title: Rename the init response value initialised
status: raw
kind: chore
created: 2026-11-08
updated: 2026-11-08
summary: Spell the git state of the init response American (initialized) in the schema, the code, the printed line and the fixtures, before 1.0.0 fixes the API
---

# Rename the init response value initialised

## Goal

The `init` response has a `git` field with the values `none`, `initialised`, `existing` and `failed`, and the command prints `initialised <dir> (records in <docs>/)`. The rest of the repository is American English (work item 261108-jvv4), so this one value is British. It was left alone there because it is part of the API: the schema (`docs/spec/akm-api/schemas/init.json`), the fixtures and their copy, the CLI golden files, `internal/api/ops_init.go`, `internal/cli/cli.go`, the README and the `akm` skill reference (`init.md`, two copies) all name it.

Changing a value of the API is a change a client can see. The tool is below 1.0.0, so now is the cheap time. Decide whether to rename it outright or to accept both for a release, say which in the spec and the changelog, then make the change everywhere in one commit and run the conformance fixtures (`scripts/sync-fixtures.sh` for the copy).

Done when:

- [ ] the value is `initialized` in the schema, the code, the printed line, the docs and the fixtures, or the decision to keep `initialised` is written down with its reason
- [ ] `go test ./...` and the schema check pass
- [ ] the changelog says what changed and who is affected

## Next

Read the spec's compatibility rules (section on MINOR and MAJOR changes), pick rename or accept-both, then grep for `initialised`.

## Log

- 2026-11-08: filed.
