# Roadmap

## Now

Nothing in flight. Agent rails hardening (261101-vs5c) is done; the rerun on other agent harnesses waits for a first release (261106-xrpr).

## Next

Nothing queued. Work decomposition (261101-t8h8) is done; raw items wait for triage (`akmctl ls --status raw`).

## Later

### 1.0.0

Exit criteria (draft, for Alexander to edit). The layout, the API and the spec are declared stable at 1.0.0, and:

1. **Adoption works.** An existing repository can be brought into AKM with `migrate`: discovery finds its records, the plan simulates every step, and apply reports what it leaves broken.
2. **Writes are safe.** A write that changes several files cannot stop half-way, and the last operation can be undone.
3. **Other layouts are first-class.** `check` can validate a repository against a layout other than the tool's own.
4. **The standard has no known gaps.** One copy of the API fixtures (done, 261029-x4ws), and the known rule defects (L2 done, 261029-z4gk; S2, orphaned ADRs from the done wizard) are closed.
5. **Compatibility is stated.** The spec says what a MINOR and a MAJOR change to the layout, the API and the extensions mean.

Items:

- [261029-tgfe] Give migrate a discovery phase for repositories with no manifest
- [261029-47m0] migrate plan: simulate move_file steps
- [261029-xtqj] Multi-write moves can stop half-way with no rollback
- [261029-shz5] Undo covers only the last safe move or capture
- [261028-a2wd] Let check --target accept layouts other than the tool's own
- [261029-nkqv] A new ADR made by the done wizard can be orphaned

### Not promised for 1.0.0

Workbench and polish:

- [261029-95f2] Replace the init questionnaire with a Bubble Tea wizard
- [261029-xeph] Measure the workbench poll, then add a conditional poll with a change token
- [261029-12sx] Push updates to the workbench with a watcher, keeping polling as the fallback

- [261030-3vgb] status lifecycle graphics
- [261029-g8x2] The review diff is against the current file, not after earlier writes
- [261029-jm46] Restart of akmctl ui forces a new URL
- [261026-mrdr] pkm-workbench: sort options
- [261026-np3c] pkm-workbench: delete a task, with its backlinks listed
- [261026-zg0g] pkm-workbench: filters
- [261029-mj49] After a capture, offer to open the new record
- [261107-49ta] Animation settings in the workbench
- [261107-qsd4] More lane-move themes
- [261108-gvzn] Workbench theming
- [261107-7v9y] Dropped fizzle and close streaks

Hardening and small rules:

- [261108-2ncf] Basic Windows support
- [261108-qv1w] CLI and check test gaps
- [261108-rwqj] Shared test helpers and run time
- [261108-f1tv] Run fixtures through both bindings
- [261108-79t5] Skill text checks
- [261108-xajw] Workbench page tests
- [261028-e7ej] Make rule G3 check more command forms than scripts/ paths
- [261028-1k0e] Say which is which in the ADR index supersedes column
- [261028-x1b4] Match Prettier's width on non-standard emoji ZWJ pairs
- [261029-2k7z] spike for shell completions
- [261030-55cs] akmctl self update
- [261030-43jw] Evaluate and potentially migrate to Oxc toolchain (oxlint / oxfmt)

## Dropped
