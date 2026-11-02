---
id: "261101-4gk3"
title: acceptance criteria section
summary: the acceptance criteria (done-when) should have its own section and more rigid schema
status: raw
parent: 261101-vs5c
created: 2026-11-01
updated: 2026-11-01
---

# acceptance criteria section

## Goal

The `Done when:` words open a set of checkboxes which form pass/fail criteria.
Work items have mixed placement: 
- `\n\nDone when:\n\n` surrounded by newlines
- `... Done when:\n\n` occurs at the end of a goal. the "Done when" should be a managed AKM section. (261031-02ff, 261031-72bd, 261031-fqbg, etc.)
The tool need only parse the checkboxes in the section.
Consider rejecting all non-checkbox content. 
This could be too restrictive and lead to frustrating UX.
Consider implementing an ownership marker for managed sections.
Ownership marker might also be called "section ownership declaration".
Settle the name and register in the glossary.
The goal is to alert the user and convey the consequences of hand-edits in a managed section. 
`Managed Section` could be a new glossary entry.

The rough sketch with options:

```md
## Done When

<!-- managed-by: <owner> -->
_Managed by: <owner>_
_[<owner>-managed]_

- [ ] incomplete item
- [x] complete item
```

The raw comment `<!-- managed-by: <owner> -->` is the front-runner since it naturally gets hidden when markdown is rendered. 
Obsidian has been tested and shows the comment line in both edit modes (live preview/source).

## Log

- 2026-11-01: filed.
