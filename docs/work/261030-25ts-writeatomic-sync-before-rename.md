---
id: "261030-25ts"
title: writeatomic sync before rename
summary: WriteAtomic renames an unsynced temp file over the record, so a power loss can leave a short file; sync the file before the rename and the folder after.
status: active
kind: bug
priority: 1
created: 2026-10-30
updated: 2026-10-30
claimed_by: claude
claimed_at: 2026-10-31T04:44
---

# writeatomic sync before rename

## Goal

`fsx.WriteAtomic` writes a temp file and renames it over the target (API spec, section 7: "a crash leaves the old file or the new file, never a mix"). It does not `Sync` the temp file first. That holds for a process crash. After a power loss or kernel crash, a filesystem that orders the rename before the data (ext4 with delayed allocation, and others) can leave the new name pointing at a short or empty file, and the old content is gone. Records are the repository's memory, so the promise should hold there too.

Found while explaining the write path to the maintainer (2026-10-30).

Done when:

- [ ] `WriteAtomic` calls `Sync` on the temp file before closing it and renaming; a failed sync removes the temp file, leaves the target untouched and returns the error.
- [ ] After a successful rename it syncs the containing folder on platforms that allow it (not Windows), best effort: an error from the folder sync is not returned, because the rename has happened, and a filesystem that refuses it (`EINVAL`) is not a failure.
- [ ] A test with the sync injected (a package variable) shows the order: sync, rename, folder sync; and that a sync error leaves the old file and no `.akm-*.tmp` behind.
- [ ] The other direct writers are looked at and the Log says what was decided: `migrate`'s scratch copy (`internal/migrate/migrate.go`) and the `--out` handoff file (`internal/cli/cli.go`) are not records and stay as they are, unless the look finds a reason.
- [ ] The API spec's Atomicity bullet says the file is flushed before the rename; the changelog has a line.

## Next

Add `syncFile` and `syncDir` package variables in `internal/fsx/write.go`, call them in `WriteAtomic`, then the order and failure tests.

## Log

- 2026-10-30: filed, from the maintainer's question about the write path.
- 2026-10-30: claimed (claude).
