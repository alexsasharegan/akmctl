---
id: "261031-wjs3"
title: close rewrites the Markdown links its move breaks
status: done
created: 2026-10-31
updated: 2026-10-31
closed: 2026-10-31
kind: bug
summary: close moves an item into done/ and breaks relative Markdown links in and to it; it should rewrite them in the same write set
---

# close rewrites the Markdown links its move breaks

## Goal

`close` moves an item from `docs/work/` to `docs/work/done/` and edits nothing else, so every relative Markdown link into or out of the file breaks until someone fixes it by hand. Seen closing 261031-yw5v: its own links to `../adr/...` needed `../../adr/...`, and the sibling item 261031-j4bc linked to it by bare file name. `check` (rule L1) catches the breakage afterwards; the tool should not cause it. Records that use wikilinks (`reference_style: wikilink`) resolve by id and are not affected.

Alexander's decisions (2026-10-31): `close` rewrites the links itself, as part of the operation, not an agent afterwards; it covers records (work items and ADRs), `docs/ROADMAP.md` and other Markdown under the docs folder; this is a standalone item, not part of the body-schema work 261031-gtw6.

Design notes: the rewrite is a multi-file write, so it shares the problem of 261029-xtqj (no rollback when a multi-write stops half-way): do the plan first, write under the repository lock, and report every file changed in the response. 261029-m0vm (report inbound links `move_file` leaves broken) is the migrate-side twin; one link resolver should serve both. Links inside code fences must be left alone, which `internal/body` `Locate` can tell (see 261031-2j14, the first consumer of it), so this is best built after 2j14. Done when:

- [x] After `close`, no markdown link in or to the moved item is broken: its own relative links are re-based, and inbound inline links and reference definitions in records, the roadmap and other docs are rewritten to the new path.
- [x] Links inside code fences and links to anything other than the moved file are untouched; a test covers CRLF files and an anchor (`file.md#section`).
- [x] The response and CLI output list every other file changed, and `commit.paths` includes them.
- [x] A refused or failed write leaves the repository as it was, or the half-done state is reported (decided with 261029-xtqj).
- [x] `docs/spec/akm-api-spec.md`, `.agents/skills/akm/references/close.md` and its embedded copy say what `close` now edits; fixtures cover it (both copies, see 261029-x4ws).
- [x] The `akm-track-work` close steps no longer tell an agent to fix links by hand.

## Next

Done and verified; close with the Outcome.

## Outcome

- Links after a close: `TestCloseRewritesLinks` (`internal/api/closelinks_test.go`) closes an item that links to a sibling, an ADR, itself, a missing file and an anchor, with a note, a work file and the roadmap linking back; afterwards the item's links are re-based, every inbound link points into `done/`, and `check` reports only the link that was already broken. Fixture `close-8` (both copies) covers the same through the API.
- Code, other targets, CRLF, anchors: the same test and `internal/mdlinks/mdlinks_test.go` (`TestScanSkips`, `TestRewriteKeepsTheRest`) check that links in fences, code spans and frontmatter and links to other files stay as written, that a CRLF line and a `#fragment`, an angle-bracket form with a title and a `./` prefix survive, and that escaped and spaced paths are escaped.
- Reporting: `links_rewritten` in the response and CLI output, and `commit.paths` lists every file written (test and `close-8`; schema `close.json`).
- Failure: `TestCloseReportsLinkWriteFailure` makes one file unwritable; the move stands, the other file is rewritten and the failure is in `link_errors`. Decision with 261029-xtqj: report, do not roll back; a journal stays that item's question.
- Docs: `docs/spec/akm-api-spec.md` (4.12 and atomicity), `close.md` and the `akm-track-work` close step in `.agents/skills` and the embedded copy (a test keeps them equal), glossary, architecture, changelog.

`gofmt -l .` (empty), `go vet ./...`, `go test -race ./...`, `npm run format:check` and `npm run schemas` (79 cases, 0 failures) pass, and `go list -deps ./cmd/akmctl | grep -c goldmark` prints 0. This item was closed with the new build, so the close itself was the first live run.

Not done: wikilink repair (they resolve by name); reference-style link definitions (`[id]: path`), which `check` does not read either; links in files outside AGENTS.md, the docs folders and the roadmap; the journal of 261029-xtqj; `m0vm`'s migrate-side report, which can reuse `mdlinks` (noted there).

Decision: none (the policy was set by Alexander on 2026-10-31 and is in this item's Goal; no hard-to-reverse choice beyond it).

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from the closing of 261031-yw5v, which broke three links.
- 2026-10-31: claimed (claude).
- 2026-10-31: built internal/mdlinks and the close link repair; check L1 now reads links through mdlinks; spec, skills, schema, fixture close-8, glossary, architecture and changelog updated.
