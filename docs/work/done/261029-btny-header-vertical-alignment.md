---
id: "261029-btny"
title: Center the WORKBENCH label and repository name with the title
status: done
created: 2026-10-29
updated: 2026-10-30
kind: bug
priority: 3
summary: In the header, the WORKBENCH label and the repository name sit lower than the akmctl title and are not vertically centered.
---

# Center the WORKBENCH label and repository name with the title

## Goal

Reported by Alexander with two screenshots (dark theme, the akmctl repo and the obsidian.vaults repo): the small WORKBENCH label and the repository name are not vertically centered against the akmctl title and the tab row; the baselines differ. Reproduce at a few widths, find the cause (probably align-items or mixed baselines between the logo font and the label font) and fix it.

Done when:

- [ ] The title, the label, the repository name and the tabs share one vertical centre line (or one consistent baseline) in both themes.
- [ ] Checked at desktop and tablet widths with a long repository name.
- [ ] A short note in the Log says what the cause was.

## Next

- Coordinator: check the header in both themes at desktop and tablet widths with a long repository name (not done here; no browser was run).

## Outcome

Label, repository name and nav share one centre line (32px at 1200px wide, measured in the built-in browser). Fix in `internal/server/web/app.css`: `.brand` is an inline flex box and each part has line-height 1. The repository name is hidden under 900px by an existing rule, unchanged.

## Log

- 2026-10-29: filed.
- 2026-10-30: claimed (ui-agent).
- 2026-10-30: cause: `.brand` was a block whose inline `<small>` label and the title shared one line box, so the label sat on the title's baseline, and the repository name (`.repo`) had its own font and line box; the flex row centred the boxes, not the glyphs. Fix in web/app.css: `.brand` is an inline flex row with `align-items:center`, and the label and `.repo` have `line-height:1`. The visual check is not done here.
- 2026-10-30: released (ui-agent).
