---
id: "261028-1k0e"
title: Say which is which in the ADR index supersedes column
summary: The index column 'supersedes / superseded_by' does not show which id is which; split or label it without breaking the Prettier-stable form.
status: ready
kind: chore
priority: 1
created: 2026-10-28
updated: 2026-10-28
---

# Say which is which in the ADR index supersedes column

## Goal

The generated ADR index has a column `supersedes / superseded_by` that lists the ids with a slash. A row with one id does not say whether the record supersedes it or was superseded by it. A reader has to open the file or read the group (Not in force).

Done when:

- [ ] The index shows which relation each id is: for example separate columns, or ids prefixed by the relation.
- [ ] The specification's index section and the fixtures in `testdata/api/` are updated in the same change (the spec copy is the source).
- [ ] The output stays byte-stable under the pinned Prettier, and rule A4 passes on a repository with a superseded record.

## Log

- 2026-10-28: filed.
