---
id: "261108-mwtr"
title: windows write fails while record is open
status: raw
parent: 261108-2ncf
created: 2026-11-08
updated: 2026-11-08
---

# windows write fails while record is open

## Goal

Found by the Windows smoke harness ([261108-hwcq](261108-hwcq-windows-powershell-harness.md)), part of [261108-2ncf](261108-2ncf-basic-windows-support.md).

Command: `akmctl set <id> priority=3` on `windows-latest` while another process holds the record open for reading with `FileShare.ReadWrite | FileShare.Delete` (an editor, a sync client or a virus scanner does this).

Actual: exit 2, `akmctl: unreadable: rename ...\docs\work\.akm-NNN.tmp ...\261108-xxxx-slug.md: Access is denied.` The record is unchanged. The code `unreadable` is wrong for a failed write, and the message does not say a program holds the file or to retry.

Expected: the write retries for a short while (as the lock does) and succeeds once the reader lets go, or the refusal says that another program holds the record, and no `.akm-*.tmp` file is left behind.

Done when:

- [ ] The command, the error code and the message are decided, with a test or a harness case that fails before the change and passes after.
- [ ] No `.akm-*.tmp` file remains after a refused write.

## Next

Read how the atomic write renames (internal/repo or the writer), check whether MoveFileEx with a retry loop is enough, then fix; the harness case is in scripts/windows-smoke.ps1 (known fault 261108-mwtr).

## Log

- 2026-11-08: filed.
- 2026-11-08 13:16-07:00: filed by the Windows smoke harness run (PR 51)
