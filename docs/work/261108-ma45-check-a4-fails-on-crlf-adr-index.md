---
id: "261108-ma45"
title: check a4 fails on crlf adr index
status: raw
parent: 261108-2ncf
created: 2026-11-08
updated: 2026-11-08
---

# check a4 fails on crlf adr index

## Goal

Found by the Windows smoke harness ([261108-hwcq](done/261108-hwcq-windows-powershell-harness.md)), part of [261108-2ncf](261108-2ncf-basic-windows-support.md).

Command: `akmctl check` in a repository checked out on Windows with `core.autocrlf=true`, so `docs/adr/README.md` has CRLF line endings.

Actual: exit 1, `docs/adr/README.md:1: [E] A4 the ADR index differs from what akmctl index generates; run akmctl index`. Running `index` would rewrite the whole file to LF, which git then shows as a change.

Expected: `check` accepts an index that differs from the generated one only in line endings (or says so and `index` keeps the file's endings), because the records keep their line endings byte for byte.

Done when:

- [ ] A test with a CRLF `README.md` passes `check` on any platform.
- [ ] The Windows smoke harness drops its known-fault exception for this item.

## Next

Find rule A4 in internal/check, compare after normalizing CRLF, add the test; remove the known-fault exception in scripts/windows-smoke.ps1.

## Log

- 2026-11-08: filed.
- 2026-11-08 13:16-07:00: filed by the Windows smoke harness run (PR 51)
- 2026-11-08 13:52-07:00: the CRLF audit (261108-kbve) found the same cause in index: it rewrites a CRLF docs/adr/README.md to LF, so the fix should keep the file's endings
