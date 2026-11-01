---
id: "261031-wjs3"
title: close rewrites the Markdown links its move breaks
status: raw
created: 2026-10-31
updated: 2026-10-31
kind: bug
summary: close moves an item into done/ and breaks relative Markdown links in and to it; it should rewrite them in the same write set
---

# close rewrites the Markdown links its move breaks

## Goal

`close` moves an item from `docs/work/` to `docs/work/done/` and edits nothing else, so every relative Markdown link into or out of the file breaks until someone fixes it by hand. Seen closing 261031-yw5v: its own links to `../adr/...` needed `../../adr/...`, and the sibling item 261031-j4bc linked to it by bare file name. `check` (rule L1) catches the breakage afterwards; the tool should not cause it. Records that use wikilinks (`reference_style: wikilink`) resolve by id and are not affected.

Alexander's decisions (2026-10-31): `close` rewrites the links itself, as part of the operation, not an agent afterwards; it covers records (work items and ADRs), `docs/ROADMAP.md` and other Markdown under the docs folder; this is a standalone item, not part of the body-schema work 261031-gtw6.

Design notes: the rewrite is a multi-file write, so it shares the problem of 261029-xtqj (no rollback when a multi-write stops half-way): do the plan first, write under the repository lock, and report every file changed in the response. 261029-m0vm (report inbound links `move_file` leaves broken) is the migrate-side twin; one link resolver should serve both. Links inside code fences must be left alone, which `internal/body` `Locate` can tell (see 261031-2j14, the first consumer of it), so this is best built after 2j14. Done when:

- [ ] After `close`, no markdown link in or to the moved item is broken: its own relative links are re-based, and inbound inline links and reference definitions in records, the roadmap and other docs are rewritten to the new path.
- [ ] Links inside code fences and links to anything other than the moved file are untouched; a test covers CRLF files and an anchor (`file.md#section`).
- [ ] The response and CLI output list every other file changed, and `commit.paths` includes them.
- [ ] A refused or failed write leaves the repository as it was, or the half-done state is reported (decided with 261029-xtqj).
- [ ] `docs/spec/akm-api-spec.md`, `.agents/skills/akm/references/close.md` and its embedded copy say what `close` now edits; fixtures cover it (both copies, see 261029-x4ws).
- [ ] The `akm-track-work` close steps no longer tell an agent to fix links by hand.

## Next

Not started. Read internal/api close (ops_write.go), 261029-xtqj and 261029-m0vm; decide the link forms to handle (inline, reference definition, autolink excluded); then make it ready. Build after 261031-2j14.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from the closing of 261031-yw5v, which broke three links.
