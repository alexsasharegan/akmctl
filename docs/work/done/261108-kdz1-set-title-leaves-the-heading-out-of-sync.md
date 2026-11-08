---
id: "261108-kdz1"
title: set title leaves the heading out of sync
status: done
kind: bug
priority: 3
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
---

# set title leaves the heading out of sync

## Goal

Found by the cold-agent retest of 261101-t8h8. `akmctl set <id> title="New"` changes the frontmatter title and leaves the `# heading` as it was, so `check` later warns WK12 (title differs from the first heading) and nothing at `set` time says so. Expected: `set title=` also rewrites the first heading, or says in its output that the heading is now out of sync (and where to change it). Done when:

- [x] `set title=` keeps WK12 quiet (it rewrites the first heading), or prints the heading it left behind and how to change it.
- [x] A test covers it; the `set` reference says which.

## Next

## Outcome

set title= on a work item renames the first heading when it matched the old title; a heading that already differed is left alone.

## Log

- 2026-11-08: filed.
- 2026-11-08 08:48-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08 08:48-07:00: set now renames the first heading of a work item when it equals the old title (record.RetitleHeading); TestSetTitleRenamesTheHeading; spec, skill reference (both copies) and changelog updated. No fixture: its etags are hand-kept (see the audit, F9, 261108-3kse).
