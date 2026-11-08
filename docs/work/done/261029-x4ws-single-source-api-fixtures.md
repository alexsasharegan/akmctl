---
id: "261029-x4ws"
title: Keep one copy of the API fixtures
summary: testdata/api is a hand-kept copy of docs/spec/akm-api/.fixtures; make one the source and generate or check the other.
status: done
kind: chore
priority: 3
created: 2026-10-29
updated: 2026-11-08
closed: 2026-11-08
---

# Keep one copy of the API fixtures

## Goal

The API fixtures exist twice: `docs/spec/akm-api/.fixtures` (the spec's, with `make_fixtures.py`) and `testdata/api` plus `testdata/yaml-lite.json` (the tests'). Since 2026-10-29 both are in this repository and are kept in step by hand ([261029-j0d0](../../adr/261029-j0d0-one-tool-spec-and-fixtures-in-this-repository.md)). Make one the source: either `testdata` reads the spec directory, or a script or test copies it and a check fails when they differ. `testdata/` stays byte-exact (`.gitattributes`).

Done when:

- [x] One directory is the source; the other is generated or checked, and the choice is written in `docs/spec/akm-api/README.md` and `docs/architecture.md`.
- [x] A change to one copy without the other fails `go test` or CI.
- [x] `make_fixtures.py` is documented as the way to regenerate, or replaced.

## Next

## Outcome

docs/spec/akm-api/.fixtures is the one source. scripts/sync-fixtures.sh writes the copy in testdata, and TestFixtureCopiesAgree fails on any file or byte that differs. make_fixtures.py is documented as a frozen reference, since it no longer reproduces the cases.

## Log

- 2026-10-29: filed.
- 2026-11-08 00:23-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08 00:24-07:00: ticked "One directory is the source; the other is generated or checked, and the choice is written in `docs/spec/akm-api/README.md` and `docs/architecture.md`.": spec .fixtures is the source; scripts/sync-fixtures.sh copies it; written in akm-api/README.md, architecture.md, AGENTS.md
- 2026-11-08 00:24-07:00: ticked "A change to one copy without the other fails `go test` or CI.": TestFixtureCopiesAgree fails on a one-byte change and on a missing or extra file; passes after the sync script
- 2026-11-08 00:24-07:00: ticked "`make_fixtures.py` is documented as the way to regenerate, or replaced.": make_fixtures.py had drifted (18 differing files in a scratch run); marked a frozen reference in its docstring, the README and spec section 9
- 2026-11-08 00:24-07:00: Chose the check-and-copy form over pointing 100 test paths at a hidden spec directory: tests and testdata stay as they are. make_fixtures.py cannot be the regeneration route (a scratch run differed in 18 files: it would overwrite hand-kept cases and lacks bases, show-3 and ls-7), so it is documented as a frozen reference. No CHANGELOG entry: nothing a user sees changed.
