---
id: "261108-jvv4"
title: American spelling pass
status: done
kind: chore
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
summary: Spell every word American in the repository's prose and code comments, and say so in AGENTS.md
---

# American spelling pass

## Goal

The repository's words are American English (color, behavior, license, neighbor, canceled, organize). Some files used British forms. Rewrite them, leave out what must not change, and add the rule to `AGENTS.md`.

Left alone on purpose: accepted ADRs (their bodies are never edited), the conformance fixtures and their copies (the expected outputs are data), the Unicode data files, file names and slugs, the API value `initialised` (the `init` response, its schema and its printed line), and `aria-labelledby`.

Done when:

- [x] no British spelling is left in prose, comments or identifiers outside the list above
- [x] `AGENTS.md` says to write American English
- [x] the tests, `akmctl check`, the format check and the schema check pass
- [x] what was left alone is filed as work (the `initialised` value)

## Next

## Outcome

- 68 files changed (150 lines): docs, work items, comments and the identifier `neighbour` (now `neighbor`, `js/moves.js` and its users). The error text `capitalised ... is not supported` and the migration hand-off prompt (`judgement`) changed wording.
- `AGENTS.md` has the rule.
- Checks run: `go vet ./...`, `go test -race ./...`, `akmctl check` (0 errors), `npm run format:check`, `npm run schemas`.
- Not changed, on purpose: accepted ADRs (the body of 261028-p4x9, 261028-zh0w, 261031-nvka, 261101-tr8e and others; none edited), the conformance fixtures and `testdata/api`, `internal/width/data`, file names and slugs, `aria-labelledby`, and the files of pull request 67 (the colour rename there; run the pass again after it merges). The API value `initialised` is work item 261108-7szk.
- Unverified: the pass is a script over whole words, so a British form inside a longer token (a hyphenated compound, a path) was skipped; a search for the common forms finds none outside the exclusions, but it is not a spell check.

Decision: none.

## Log

- 2026-11-08: filed.
- 2026-11-08 16:28-07:00: claimed (claude-sonnet-5-5-spelling).
- 2026-11-08 16:28-07:00: ticked "no British spelling is left in prose, comments or identifiers outside the list above": scripted pass over 582 tracked files, reviewed diff of code lines; a search for British forms finds only the exclusions
- 2026-11-08 16:28-07:00: ticked "`AGENTS.md` says to write American English": AGENTS.md rule added
- 2026-11-08 16:28-07:00: ticked "the tests, `akmctl check`, the format check and the schema check pass": go vet, go test -race ./..., akmctl check, format:check, schemas all pass
- 2026-11-08 16:28-07:00: ticked "what was left alone is filed as work (the `initialised` value)": work item 261108-7szk (initialised value of the init response)
