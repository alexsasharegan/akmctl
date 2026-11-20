---
id: "261119-pk23"
title: drawer field table layout
status: done
kind: bug
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
closed: 2026-11-19
summary: "The drawer field table drifts apart: name column, dividers and Edit links misalign"
---

# drawer field table layout

## Goal

The record drawer's field table lays out wrong (found while trying the centered drawer in the spike 261031-18t8). .fm-name is display:flex on a th, so the th stops being a table cell: the name column, the row dividers and the Edit links drift apart. Seen on a work item with a claim: a wide gap between the field names and the values, dividers under the names that stop short of the value rows, and Edit sitting lower than its value.

Fix in app.css: leave the th a table cell, stack its two spans as blocks, align the cells on their first line (vertical-align:baseline), and trim the Edit button's vertical padding.

Done when:

- [x] The name column is narrow, the row dividers run the full width, and Edit lines up with its value, in light and dark, on a record with a claim (before and after screenshots on the pull request).
- [x] Editing a field still works (the existing drawer tests pass).

## Next

## Outcome

1. The name column is narrow, the row dividers run the full width, and Edit lines up with its value, in light and dark, on the record where the defect was seen: before and after screenshots of that record were taken and Alexander inspected the running build. He decided they need not go on the pull request. The locked label and the Edit links also share one right edge (measured: every row ends at x 1187, where some ended at 1170), and the not-set marker is Radon, upright.
2. Editing a field still works: go test -race ./internal/server/... passes, including the drawer tests that edit fields.

Decision: none.

## Log

- 2026-11-19: filed.
- 2026-11-19 22:01-07:00: claimed (claude-sonnet-5-5-drawer-table).
- 2026-11-19 22:02-07:00: ticked "Editing a field still works (the existing drawer tests pass).": go test -race ./internal/server/... passes, including the drawer tests that edit fields
- 2026-11-19 22:02-07:00: Applied the four CSS changes; before and after screenshots of record 261108-wg98 (the record the defect was seen on) are in .evidence/drawer-table-fix; go vet, gofmt and go test -race ./internal/server pass. Left: attach the screenshots to the pull request.
- 2026-11-19 22:18-07:00: Follow-up from Alexander inspecting the tip: the locked label now ends on the same right edge as Edit (negative right margin on the button, which keeps its hover area), and not set is Radon with no italic, as the merged upright-Radon change (261119-5680) does elsewhere. The diff-hunk italic in the check view was left alone.
- 2026-11-19 22:36-07:00: ticked "The name column is narrow, the row dividers run the full width, and Edit lines up with its value, in light and dark, on a record with a claim (before and after screenshots on the pull request).": before and after screenshots of record 261108-wg98 (light and dark) taken and inspected by Alexander in the running build; every right-edge measured at x 1187 after; not attached to the pull request, by his decision
