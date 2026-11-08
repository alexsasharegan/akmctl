---
id: "261108-2z9m"
title: new files are lf in a crlf repository
status: raw
parent: 261108-2ncf
created: 2026-11-08
updated: 2026-11-08
---

# new files are lf in a crlf repository

## Goal

Found by the CRLF and path audit ([261108-kbve](done/261108-kbve-windows-crlf-and-path-audit.md)), part of [261108-2ncf](261108-2ncf-basic-windows-support.md).

Command: `akmctl new work <slug>`, `new adr <slug>` and `split <id> <slug>...` in a repository whose records are CRLF (a Windows checkout with `core.autocrlf=true`).

Actual: the new file is written with LF, so the repository holds CRLF and LF records side by side. Every edit of an existing record keeps its CRLF; only new files differ.

Expected: decide the rule. Either a new file follows the line ending of its neighbours (for example the first record in the folder), or LF is the documented rule because git normalizes on `add`. Whichever is chosen is written in the guide and tested.

Done when:

- [ ] The rule is decided and written down (an ADR draft if a real alternative is rejected).
- [ ] A test covers `new` and `split` in a CRLF repository.

## Next

Decide with the maintainer: follow the neighbours or stay LF. Low priority; git normalizes LF on add when autocrlf is on.

## Log

- 2026-11-08: filed.
- 2026-11-08 13:52-07:00: filed by the CRLF and path audit
