---
id: "261028-1k0e"
title: Say which is which in the ADR index supersedes column
summary: The index column 'supersedes / superseded_by' does not show which id is which; split or label it without breaking the Prettier-stable form.
status: active
kind: chore
priority: 1
created: 2026-10-28
updated: 2026-11-16
claimed_by: claude-fable-5-2
claimed_at: 2026-11-16T13:58
---

# Say which is which in the ADR index supersedes column

## Goal

The generated ADR index has a column `supersedes / superseded_by` that lists the ids with a slash. A row with one id does not say whether the record supersedes it or was superseded by it. A reader has to open the file or read the group (Not in force).

Done when:

- [x] The index shows which relation each id is: for example separate columns, or ids prefixed by the relation.
- [x] The specification's index section and the fixtures in `testdata/api/` are updated in the same change (the spec copy is the source).
- [x] The output stays byte-stable under the pinned Prettier, and rule A4 passes on a repository with a superseded record.

## Next

Prefix each id in the last index column with its relation in internal/index, test it with a superseded record, and update the standard's index section.

## Outcome

- The last column prints `supersedes <id>` and `superseded by <id>`: `row` in `internal/index/index.go`. This repository's index shows it on three rows.
- The standard's index section says so. No fixture repository has a superseded record, so the fixtures did not change; the new test builds one from the `basic` fixture.
- `TestSupersedesCellNamesTheRelation` checks the cells and that A4 passes; the pinned Prettier 3.9.10 left the generated index of that test and this repository's index unchanged (run by hand, since the Go tests do not call Prettier).

Decision: none (the pick between prefixes and columns is in the Log).

## Log

- 2026-10-28: filed.
- 2026-11-16 06:58-07:00: claimed (claude-fable-5-2).
- 2026-11-16 06:59-07:00: ticked "The index shows which relation each id is: for example separate columns, or ids prefixed by the relation.": internal/index row(): supersedes <id>, superseded by <id>; this repository's docs/adr/README.md regenerated
- 2026-11-16 06:59-07:00: ticked "The specification's index section and the fixtures in `testdata/api/` are updated in the same change (the spec copy is the source).": standard, The generated index; no fixture repository in testdata/api has a superseded record, so no fixture changed
- 2026-11-16 06:59-07:00: ticked "The output stays byte-stable under the pinned Prettier, and rule A4 passes on a repository with a superseded record.": TestSupersedesCellNamesTheRelation (A4 passes); prettier 3.9.10 --check passed on that index and on docs/adr/README.md
- 2026-11-16 06:59-07:00: pick: prefixed ids, not separate columns. The table keeps its six columns and header, so only a repository that has a superseded record gets a stale index on upgrade; a new column would have made every index stale and the widest table wider
