---
id: "261108-kdz1"
title: set title leaves the heading out of sync
status: raw
kind: bug
priority: 3
created: 2026-11-08
updated: 2026-11-08
---

# set title leaves the heading out of sync

## Goal

Found by the cold-agent retest of 261101-t8h8. `akmctl set <id> title="New"` changes the frontmatter title and leaves the `# heading` as it was, so `check` later warns WK12 (title differs from the first heading) and nothing at `set` time says so. Expected: `set title=` also rewrites the first heading, or says in its output that the heading is now out of sync (and where to change it). Done when:

- [ ] `set title=` keeps WK12 quiet (it rewrites the first heading), or prints the heading it left behind and how to change it.
- [ ] A test covers it; the `set` reference says which.

## Log

- 2026-11-08: filed.
