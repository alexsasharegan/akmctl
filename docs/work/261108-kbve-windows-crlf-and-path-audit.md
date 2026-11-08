---
id: "261108-kbve"
title: windows crlf and path audit
status: active
kind: chore
parent: 261108-2ncf
priority: 2
created: 2026-11-08
updated: 2026-11-08
claimed_by: claude-sonnet-5-5-win
claimed_at: 2026-11-08T20:47
---

# windows crlf and path audit

## Goal

An audit of CRLF and path handling. Part of [261108-2ncf](261108-2ncf-basic-windows-support.md). Records use forward slashes, Windows users pass `C:\...` and git may check files out with CRLF. List every write operation and every place a path is built or compared, and check each one with CRLF input and Windows-style paths, simulating both on any platform where it can be done (tests with `\r\n` fixtures and a path-style shim). Do not change behavior in the audit item; file the faults.

Done when:

- [ ] Every write operation (`new`, `set`, `claim`, `release`, `log`, `tick`, `section`, `remove`, `close`, `split`, `migrate` steps) is run on a CRLF record and keeps CRLF; the table of results is in the Log or a design note.
- [ ] `--repo`, `find`, link rewriting and wikilink resolution are checked with backslash and mixed paths.
- [ ] Each fault found is filed as its own item under this one.

## Next

Drive every write op on a CRLF repository with the built binary (scratch script), check for kept endings; check --repo, find, link rewriting and wikilinks with backslash and mixed paths; put the result table in the Log; file one item per fault.

## Log

- 2026-11-08: filed.
- 2026-11-08 13:47-07:00: claimed (claude-sonnet-5-5-win).
