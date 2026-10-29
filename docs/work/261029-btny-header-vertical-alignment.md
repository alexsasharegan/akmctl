---
id: "261029-btny"
title: Center the WORKBENCH label and repository name with the title
status: ready
created: 2026-10-29
updated: 2026-10-29
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

## Log

- 2026-10-29: filed.
