---
id: "261101-nhj4"
title: scrub private names from fixtures and tests
status: done
kind: chore
summary: Replace the maintainer's vault, client and note names in fixtures, tests, the spec and the records with neutral ones, keeping every test
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
---

# scrub private names from fixtures and tests

## Goal

The repository is public and carried the maintainer's personal vault path, a client name, a note name and a local folder name in API fixtures (folder names too), Go and JS tests, the API spec, the find skill reference and old records. They were sample strings; the tests behind them (custom layouts, paths with spaces, wikilink and project handling) stay. Replace the names with neutral ones (`Notes/tasks`, `Acme Fabrication`, `Project Notes`, `notes.vaults`) and move the fixture folders to match. Git history is not rewritten.

The API fixtures exist twice (`docs/spec/akm-api/.fixtures`, the source, and `testdata/api`, the copy; 261029-x4ws): change both and diff the trees. If an expected output depends on a path or a hash, regenerate through `make_fixtures.py`.

Done when:

- [x] no file or folder outside `.git` carries the old vault, client, note or folder names
- [x] both fixture trees are identical and every test passes
- [x] `go vet`, `go test -race`, the records check and the format check pass

## Next

Scripted replace across the tree, then diff the two fixture trees and run the full check.

## Outcome

The vault path, client, note and folder names are replaced with neutral ones across fixtures (both trees, including folder names), tests, the spec, the find skill reference and old records; every test kept. Git history is not rewritten. ADR check: none.

## Log

- 2026-11-01: filed.
- 2026-11-01 15:42-07:00: claimed (claude).
- 2026-11-01 15:44-07:00: ticked "no file or folder outside `.git` carries the old vault, client, note or folder names": Scripted replace and fixture folder moves; both fixture trees diff clean; go vet, go test -race, records check, format check and schema validation pass
- 2026-11-01 15:44-07:00: ticked "both fixture trees are identical and every test passes": Scripted replace and fixture folder moves; both fixture trees diff clean; go vet, go test -race, records check, format check and schema validation pass
- 2026-11-01 15:44-07:00: ticked "`go vet`, `go test -race`, the records check and the format check pass": Scripted replace and fixture folder moves; both fixture trees diff clean; go vet, go test -race, records check, format check and schema validation pass
- 2026-11-01 15:44-07:00: Replaced the names in 47 files, moved the four fixture folders to Notes/, updated the minted ids that depend on the path (migrate-2), and the find test. Old names remain in git history only.
