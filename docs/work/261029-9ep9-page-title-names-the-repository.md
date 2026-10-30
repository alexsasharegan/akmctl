---
id: "261029-9ep9"
title: Put the repository name in the page title
status: active
created: 2026-10-29
updated: 2026-10-30
kind: feature
priority: 2
summary: Two workbench tabs on different repositories look the same in the browser tab strip.
claimed_by: ui-agent
claimed_at: 2026-10-30T00:42
---

# Put the repository name in the page title

## Goal

Alexander runs one akmctl ui per repository (for example the vault and the akmctl repo) and cannot tell the tabs apart: the document title is the same in each. Put the repository name in the title (the same name the header shows, from version.repository.name or GET /ui/config root), and the view in it, for example 'Board - obsidian.vaults - akmctl'. Keep it current on every route change and on a record open. If two repositories have the same folder name, add the parent folder (decide the rule and write it down).

Done when:

- [ ] document.title names the repository and the current view, and updates on hash route changes.
- [ ] Two servers on different repositories show different tab titles.
- [ ] The title is set with textContent-safe APIs (document.title only), with no record text beyond the id and title.
- [ ] A pure function builds the title; tests cover long names and duplicate folder names.

## Next

- Coordinator: open two workbenches on different repositories and check the tab titles in the browser (not done here).
- Duplicate folder names: `repoLabel(root, others)` in web/js/logic.js adds the parent folder when another root shares the name, but the page passes `[]` because it only knows its own root. Wiring it needs the server to report the other open roots; decide that before closing.

## Log

- 2026-10-29: filed.
- 2026-10-30: claimed (ui-agent).
- 2026-10-30: `repoLabel` and `pageTitle` are pure functions in web/js/logic.js, tested in testdata/logic.test.mjs (run by TestPureModulesUnderNode, with long names and duplicate folder names). app.js sets document.title from `setTitle` on every route and on a record open, and once the repository name arrives from GET /ui/config. Titles use `view - repo - akmctl`; a record shows its id and title only.
