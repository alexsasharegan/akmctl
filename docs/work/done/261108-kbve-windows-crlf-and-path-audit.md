---
id: "261108-kbve"
title: windows crlf and path audit
status: done
kind: chore
parent: 261108-2ncf
priority: 2
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
---

# windows crlf and path audit

## Goal

An audit of CRLF and path handling. Part of [261108-2ncf](../261108-2ncf-basic-windows-support.md). Records use forward slashes, Windows users pass `C:\...` and git may check files out with CRLF. List every write operation and every place a path is built or compared, and check each one with CRLF input and Windows-style paths, simulating both on any platform where it can be done (tests with `\r\n` fixtures and a path-style shim). Do not change behavior in the audit item; file the faults.

Done when:

- [x] Every write operation (`new`, `set`, `claim`, `release`, `log`, `tick`, `section`, `remove`, `close`, `split`, `migrate` steps) is run on a CRLF record and keeps CRLF; the table of results is in the Log or a design note.
- [x] `--repo`, `find`, link rewriting and wikilink resolution are checked with backslash and mixed paths.
- [x] Each fault found is filed as its own item under this one.

## Next

## Outcome

Every write operation keeps CRLF on an existing record and no run left a mixed file, including the migrate steps. Path forms for --repo and find pass on the Windows runner. Faults: new files are LF in a CRLF repository (261108-2z9m), backslash link targets are neither checked nor rewritten (261108-sfj3), and index rewrites a CRLF ADR README as LF (added to 261108-ma45). The table is docs/design/windows-crlf-path-audit.md. Not covered: drive-letter and UNC paths, a console.

## Log

- 2026-11-08: filed.
- 2026-11-08 13:47-07:00: claimed (claude-sonnet-5-5-win).
- 2026-11-08 13:53-07:00: ticked "Every write operation (`new`, `set`, `claim`, `release`, `log`, `tick`, `section`, `remove`, `close`, `split`, `migrate` steps) is run on a CRLF record and keeps CRLF; the table of results is in the Log or a design note.": docs/design/windows-crlf-path-audit.md: table of every write op; none left mixed endings; two faults filed
- 2026-11-08 13:53-07:00: ticked "`--repo`, `find`, link rewriting and wikilink resolution are checked with backslash and mixed paths.": Windows runner cases in scripts/windows-smoke.ps1 (--repo forms, find, close link repair) and macOS probes, in the design note
- 2026-11-08 13:53-07:00: ticked "Each fault found is filed as its own item under this one.": 261108-2z9m and 261108-sfj3 filed under 261108-2ncf; 261108-ma45 (already filed) got the index finding
- 2026-11-08 13:53-07:00: audit done; results in docs/design/windows-crlf-path-audit.md; faults 261108-2z9m, 261108-sfj3, and 261108-ma45 (index)
