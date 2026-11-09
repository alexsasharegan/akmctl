# Roadmap

## Now

One thread is in flight, Windows support. What is left in it waits on the maintainer or on the first Windows user:

- [261108-2ncf] Basic Windows support (the parent; its installers, harness and audit are done)
  - [261108-wg98] The support-tier ADR is drafted (`proposed`); on acceptance it supersedes 261029-eegx and the README, `AGENTS.md` and `docs/release.md` follow.
  - [261108-7avk] The first Windows user runs `docs/windows-first-run.md`.
  - [261108-ct4s] The same user checks `init` and `akmctl ui` in Windows Terminal and PowerShell.

Bugs the Windows work found are filed as `raw` items under the same parent and are triaged by whoever takes the Windows thread: 261108-2z9m (new files are LF in a CRLF repository), 261108-ma45 (`check` A4 fails on a CRLF index), 261108-mwtr (a write fails while an editor holds the record) and 261108-sfj3 (backslash link targets are ignored).

The rerun of the agent rails trial on other harnesses (261106-xrpr) waits for a first release.

## Next

The queue for an agent working unattended. Take the top item that is `ready`; none needs the maintainer's decision first. Items marked _pick_ ask the agent to choose between small options and record the choice in the item's Log and the pull request; the maintainer reviews it at merge.

Without a decision:

- [261029-47m0] migrate plan: simulate move_file steps (1.0.0, criterion 1)
- [261028-a2wd] Let check --target accept layouts other than the tool's own (1.0.0, criterion 3)
- [261029-g8x2] The review diff is against the current file, not after earlier writes
- [261101-3che] show reports edit states
- [261028-e7ej] Make rule G3 check more command forms than scripts/ paths
- [261028-1k0e] Say which is which in the ADR index supersedes column (_pick_: separate columns or prefixed ids)
- [261028-x1b4] Match Prettier's width on non-standard emoji ZWJ pairs
- [261109-e540] CI steps that download from GitHub fail on a 504 and block the PR run (_pick_: warn or gate after the retries)

With a pick, for 1.0.0:

- [261029-nkqv] A new ADR made by the done wizard can be orphaned (criterion 4; _pick_: say so and offer a retry, or write the link first)
- [261029-xtqj] Multi-write moves can stop half-way with no rollback (criterion 2; _pick_: restore the earlier writes from the bytes held in memory, no journal; a `proposed` ADR if it changes a guarantee)
- [261029-shz5] Undo covers only the last safe move or capture (criterion 2; _pick_: document the limit in the key help, or read git history as ADR 261106-60xv does)
- [261029-tgfe] Give migrate a discovery phase for repositories with no manifest (criterion 1; _pick_: the report-only phase first, which changes nothing, then the proposed steps)

## Later

### 1.0.0

Exit criteria (draft, for Alexander to edit). The layout, the API and the spec are declared stable at 1.0.0, and:

1. **Adoption works.** An existing repository can be brought into AKM with `migrate`: discovery finds its records, the plan simulates every step, and apply reports what it leaves broken. Items: 261029-tgfe, 261029-47m0.
2. **Writes are safe.** A write that changes several files cannot stop half-way, and the last operation can be undone. Items: 261029-xtqj, 261029-shz5. ADR 261106-60xv chose git as the move log and kept the journal as the spike 261106-wx2a; the design note on xtqj and shz5 about a journal is an alternative to that, not the plan.
3. **Other layouts are first-class.** `check` can validate a repository against a layout other than the tool's own. Item: 261028-a2wd.
4. **The standard has no known gaps.** One copy of the API fixtures (done, 261029-x4ws), and the known rule defects (L2 done, 261029-z4gk; S2, orphaned ADRs from the done wizard, 261029-nkqv) are closed.
5. **Compatibility is stated.** The spec says what a MINOR and a MAJOR change to the layout, the API and the extensions mean. Item: 261109-vxdr, which needs the maintainer.

Needs the maintainer, in this order:

- [261109-vxdr] State what a MINOR and a MAJOR change mean: a `proposed` ADR to accept, then the spec text.
- [261101-5zf4] v1.0.0 architecture review: a collaborative review, run last so it sees the finished shape of criteria 1 to 4.

### Not promised for 1.0.0

Decisions the maintainer owns before an agent can take them:

- [261029-95f2] Replace the init questionnaire with a Bubble Tea wizard (a dependency, so an ADR first)
- [261030-55cs] akmctl self update (changes the "no network" guarantee, so an ADR first)
- [261030-43jw] Evaluate and potentially migrate to Oxc toolchain (oxlint / oxfmt): the evaluation can run unattended, the migration is the maintainer's call
- [261030-3vgb] status lifecycle graphics: a spike that needs the maintainer to look at two prototypes
- [261026-mrdr] pkm-workbench: sort options
- [261026-zg0g] pkm-workbench: filters
- [261026-np3c] pkm-workbench: delete a task, with its backlinks listed

With a pick, when the 1.0.0 queue is empty:

- [261029-xeph] Measure the workbench poll, then add a conditional poll with a change token (measure first; build only if the number is bad)
- [261029-jm46] Restart of akmctl ui forces a new URL (fix it, or record the trade-off)
- [261029-mj49] After a capture, offer to open the new record (the key-binding variant needs no toast actions and can ship first)
- [261029-2k7z] spike for shell completions
- [261029-12sx] Push updates to the workbench with a watcher, keeping polling as the fallback
- [261107-49ta] Animation settings in the workbench
- [261107-qsd4] More lane-move themes
- [261107-7v9y] Dropped fizzle and close streaks

## Dropped
