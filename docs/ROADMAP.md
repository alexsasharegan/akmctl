# Roadmap

## Now

- [261029-pzb1] Write the glossary and require every new feature to register its terms

## Next

Workbench polish and the Obsidian Bases file, in order (the styling freeze is lifted for these):

- [261029-kb33] Workbench form controls: a custom select and a summary text area with a counter
- [261029-aezq] Capture modal: header, width, a cycle field and the submit button
- [261029-bdt9] Workbench parent picker: a searchable list with match highlighting
- [261028-esyt] Show and edit a record's fields in the drawer
- [261029-kb88] Generate an Obsidian Bases file that mirrors the workbench

## Later

### 1.0.0

Exit criteria (draft, for Alexander to edit). The layout, the API and the spec are declared stable at 1.0.0, and:

1. **Adoption works.** An existing repository can be brought into AKM with `migrate`: discovery finds its records, the plan simulates every step, and apply reports what it leaves broken.
2. **Writes are safe.** A write that changes several files cannot stop half-way, and the last operation can be undone.
3. **Other layouts are first-class.** `check` can validate a repository against a layout other than the tool's own.
4. **The standard has no known gaps.** One copy of the API fixtures, and the known rule defects (L2, S2, orphaned ADRs from the done wizard) are closed.
5. **Compatibility is stated.** The spec says what a MINOR and a MAJOR change to the layout, the API and the extensions mean.

Items:

- [261029-tgfe] Give migrate a discovery phase for repositories with no manifest
- [261029-47m0] migrate plan: simulate move_file steps
- [261029-m0vm] move_file: report the inbound links it leaves broken
- [261029-xtqj] Multi-write moves can stop half-way with no rollback
- [261029-shz5] Undo covers only the last safe move or capture
- [261028-a2wd] Let check --target accept layouts other than the tool's own
- [261029-x4ws] Keep one copy of the API fixtures
- [261029-z4gk] Rule L2: understand escaped pipes in wikilinks inside tables
- [261029-wpak] Rule S2: accept block scalars in SKILL.md frontmatter
- [261029-nkqv] A new ADR made by the done wizard can be orphaned

### Not promised for 1.0.0

Workbench and polish:

- [261029-95f2] Replace the init questionnaire with a Bubble Tea wizard
- [261029-xeph] Measure the workbench poll, then add a conditional poll with a change token
- [261029-12sx] Push updates to the workbench with a watcher, keeping polling as the fallback

- [261029-30my] Filter the board and tables by project
- [261029-ctzc] Make wikilinks in drawer fields clickable and open them in Obsidian
- [261029-rmrk] Fix the card border that shows on only three sides
- [261029-57jq] The closed table flickers on reload
- [261029-8948] Show the claim holder on workbench cards
- [261030-3vgb] status lifecycle graphics
- [261029-g8x2] The review diff is against the current file, not after earlier writes
- [261029-jm46] Restart of akmctl ui forces a new URL
- [261026-mrdr] pkm-workbench: sort options
- [261026-np3c] pkm-workbench: delete a task, with its backlinks listed
- [261026-zg0g] pkm-workbench: filters
- [261029-mj49] After a capture, offer to open the new record

Hardening and small rules:

- [261030-456q] Test strategy and coverage audit
- [261029-q8z3] TestCompareAndSwapRace is flaky on Windows CI
- [261028-jzsx] Cover the Windows path beyond CI
- [261028-e7ej] Make rule G3 check more command forms than scripts/ paths
- [261028-1k0e] Say which is which in the ADR index supersedes column
- [261028-x1b4] Match Prettier's width on non-standard emoji ZWJ pairs
- [261029-2k7z] spike for shell completions
- [261030-55cs] akmctl self update
- [261030-43jw] Evaluate and potentially migrate to Oxc toolchain (oxlint / oxfmt)

## Dropped
