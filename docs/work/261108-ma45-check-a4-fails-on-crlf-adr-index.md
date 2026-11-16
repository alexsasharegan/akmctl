---
id: "261108-ma45"
title: check a4 fails on crlf adr index
status: active
parent: 261108-2ncf
created: 2026-11-08
updated: 2026-11-16
claimed_by: claude-fable-5-2
claimed_at: 2026-11-16T14:07
---

# check a4 fails on crlf adr index

## Goal

Found by the Windows smoke harness ([261108-hwcq](done/261108-hwcq-windows-powershell-harness.md)), part of [261108-2ncf](261108-2ncf-basic-windows-support.md).

Command: `akmctl check` in a repository checked out on Windows with `core.autocrlf=true`, so `docs/adr/README.md` has CRLF line endings.

Actual: exit 1, `docs/adr/README.md:1: [E] A4 the ADR index differs from what akmctl index generates; run akmctl index`. Running `index` would rewrite the whole file to LF, which git then shows as a change.

Expected: `check` accepts an index that differs from the generated one only in line endings (or says so and `index` keeps the file's endings), because the records keep their line endings byte for byte.

Done when:

- [x] A test with a CRLF `README.md` passes `check` on any platform.
- [x] The Windows smoke harness drops its known-fault exception for this item.

## Next

Find rule A4 in internal/check, compare after normalizing CRLF, add the test; remove the known-fault exception in scripts/windows-smoke.ps1.

## Outcome

- A CRLF `docs/adr/README.md` passes `check` on any platform: A4 compares after reading CRLF as LF (`index.LF`); `TestIndexOnACRLFCheckout` and a near-miss case prove it.
- `index` does not call such a file changed, and writes CRLF when it rewrites a CRLF file (`index.WithEndingsOf`), as the audit asked.
- The smoke harness no longer excuses this fault. It has not run here (no `pwsh`): the Windows CI job is the check.

The standard and the API spec say that line endings are not compared. Decision: none.

## Log

- 2026-11-08: filed.
- 2026-11-08 13:16-07:00: filed by the Windows smoke harness run (PR 51)
- 2026-11-08 13:52-07:00: the CRLF audit (261108-kbve) found the same cause in index: it rewrites a CRLF docs/adr/README.md to LF, so the fix should keep the file's endings
- 2026-11-16 07:07-07:00: claimed (claude-fable-5-2).
- 2026-11-16 07:08-07:00: ticked "A test with a CRLF `README.md` passes `check` on any platform.": TestIndexOnACRLFCheckout (internal/api) and the near miss 'A4 index checked out with CRLF' (internal/check)
- 2026-11-16 07:08-07:00: ticked "The Windows smoke harness drops its known-fault exception for this item.": scripts/windows-smoke.ps1: the known-fault branch is gone, check must exit 0 with 0 errors on the CRLF run
- 2026-11-16 07:08-07:00: not verified: the smoke harness on Windows. No pwsh here; the Windows job of the pull request runs it. index keeps CRLF when the file has any CRLF (a file with mixed endings becomes all CRLF)
