---
id: "261029-x4ws"
title: Keep one copy of the API fixtures
summary: testdata/api is a hand-kept copy of docs/spec/akm-api/.fixtures; make one the source and generate or check the other.
status: ready
kind: chore
priority: 3
created: 2026-10-29
updated: 2026-10-29
---

# Keep one copy of the API fixtures

## Goal

The API fixtures exist twice: `docs/spec/akm-api/.fixtures` (the spec's, with `make_fixtures.py`) and `testdata/api` plus `testdata/yaml-lite.json` (the tests'). Since 2026-10-29 both are in this repository and are kept in step by hand ([261029-j0d0](../adr/261029-j0d0-one-tool-spec-and-fixtures-in-this-repository.md)). Make one the source: either `testdata` reads the spec directory, or a script or test copies it and a check fails when they differ. `testdata/` stays byte-exact (`.gitattributes`).

Done when:

- [ ] One directory is the source; the other is generated or checked, and the choice is written in `docs/spec/akm-api/README.md` and `docs/architecture.md`.
- [ ] A change to one copy without the other fails `go test` or CI.
- [ ] `make_fixtures.py` is documented as the way to regenerate, or replaced.

## Log

- 2026-10-29: filed.
