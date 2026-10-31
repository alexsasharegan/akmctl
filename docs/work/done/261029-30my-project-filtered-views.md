---
id: "261029-30my"
title: Filter the board and tables by project
status: done
created: 2026-10-29
updated: 2026-10-31
closed: 2026-10-31
kind: feature
priority: 2
summary: A vault with many projects (for example AEI Fabrication) needs the board and tables narrowed to one project.
---

# Filter the board and tables by project

## Goal

Alexander's request: in the vault the board and tables show too much to give a good signal. Add a project filter so he can view one project, for example AEI Fabrication, on the board and in the tables. In the vault, the project is the extra frontmatter field project, a quoted wikilink to the project note (101 of 124 migrated tasks have it, 13 distinct values). AKM's own grouping field is parent, so design the filter around a field choice, with project as the first one, and remember the last choice per repository.

Open questions to settle first: (1) does the API need a way to filter by an extra field (today ls filters by status, kind, parent, text and similar; check the spec), or is a client-side filter over the loaded records enough at this size; (2) show a project picker with counts, plus a no-project option; (3) the pkm-workbench filters work item in the vault (261026-zg0g) holds earlier thinking.

Done when:

- [x] A control on the board and tables picks a project (or all), with counts per project and a no-project choice.
- [x] The filter applies to the board lanes, the tables and the palette search, and shows in the URL hash so a view can be bookmarked.
- [x] The choice is remembered per repository (in try/catch storage) and is easy to clear.
- [x] Works on the vault (project is a wikilink; match on the link target, aliases ignored) and on a repository with no project field (the control is hidden).
- [x] Pure-function tests for the matching; a browser check on a copy of the vault with AEI Fabrication.

## Outcome

The project filter works on the board, the tables and the palette (client-side, so the API is unchanged): logic.js has projectTarget, projectCounts, filterByProject, projectFromHash and withProject (pure, tested under node, including wikilink alias and heading stripping); project.js is the header select with counts per project and a no-project choice; app.js applies it to the views and palette, writes ?project= into the hash (kept when switching views) and the hash wins over the remembered choice, which localStorage keeps per repository inside try/catch. Cleared by choosing All projects or the palette action. Checked in the workbench on a fixture with three projects: counts, board and table filtering, the URL, and the choice surviving a reload. Not checked: a copy of the real vault (13 projects, 124 tasks), and the control being hidden in a repository with no project field (covered by a unit test only); the empty-table message still reads "Nothing is ready or active" when a filter hides everything. Decision: none.

## Log

- 2026-10-29: filed.
- 2026-10-31: claimed (claude).
