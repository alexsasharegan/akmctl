---
id: 261028-eqwq
type: task
status: done
origin: 261027-r3hh
created: 2026-10-28
updated: 2026-10-28
authored_by: agent
title: Make the akmctl repository a dogfood candidate of the AKM repo kit
---

# Make the akmctl repository a dogfood candidate of the AKM repo kit

## Goal

`akmctl` (the CLI, which is also the control plane of the AKM standard; the name is still a working title, Y22) should keep its own memory the AKM way: `AGENTS.md`, `docs/ROADMAP.md`, `docs/adr/` with the generated index, `docs/work/`, `docs/.layout`, the pinned Prettier config, and the shared skills. It is the first repository outside the vault to be stamped out from the template, so it tests the template and the standard on a real codebase.

Not yet. It waits for the template repository (the clone template of the layout spec, section 7), which waits for the v1.0 spec cleanup ([261028-j3m9-akm-spec-v1-cleanup](261028-j3m9-akm-spec-v1-cleanup.md)). When the template exists, this task stamps it onto `akmctl` the same way any new repository would get it, and records what had to be fixed in the template.

Done when: `akmctl check` passes on the `akmctl` repository itself with the layout in place; its decisions so far (Y1 to Y37 as they apply to the tool, the licence, the release process) are ADRs there or are linked from them; its remaining work (the gaps listed in [261027-r3hh-repo-records-go-tool](261027-r3hh-repo-records-go-tool.md)) is work items under `docs/work/`; and a note says what the exercise showed about the template.

Open questions for when it starts:

- Whether the tool's repository keeps its own CHANGELOG and README beside `docs/`, or the layout absorbs them.
- Whether the template is applied by copying files or by `akmctl migrate` from a repository with no manifest (layout `0.0.0`).
- Which of the tool's design records (the decision log's Y entries) become ADRs in the repository and which stay as history in the vault.

## Log

- 2026-10-28: filed by Claude Code on the Mac at Alexander's direction. Raw: not for now; it follows the template repository, which follows the v1.0 spec cleanup.
- 2026-10-28: Alexander's order: the spec cleanup first, then this task (the CLI is the first test case of the repo kit, ahead of the vault), then [261027-nrwf-vault-agent-ready-migration](261027-nrwf-vault-agent-ready-migration.md). This task now also covers building the template repository from what the exercise teaches.
- 2026-10-28: promoted to `ready`: the done-when is written above and the spec cleanup is done. Next in the queue with [261028-a149-akmctl-ui-port-pkm-workbench](261028-a149-akmctl-ui-port-pkm-workbench.md) (Alexander: the UI port is the next important task); the two touch different files and can run on separate branches.
- 2026-10-28: done by a sub-agent (Claude Code) on the branch `dogfood-layout` of `~/github.com/alexsasharegan/akmctl` (worktree `~/github.com/alexsasharegan/akmctl-wt/dogfood`), pushed to the private origin, not merged. Open questions settled: README.md and CHANGELOG.md stay at the repository root (the layout lists the README there; a changelog is release notes); the layout was copied, not migrated (layout 0.0.0 has no path to 1.0.0); 14 ADRs from the decision log (Y2, Y3, Y5, Y6, Y8 to Y10, Y13, Y14, Y17, Y18, Y26, Y28 to Y31, Y34 to Y37), the rest stay history. Lessons: [akm-template-lessons](../../design/akm-template-lessons.md).

## Next

Merge the branch `dogfood-layout` after review (Alexander; the ADRs are `accepted` with `authored_by: agent`, the merge is the acceptance). Then build the template repository from [akm-template-lessons](../../design/akm-template-lessons.md) and start the `mint_ids` item in the tool repository (`261028-8dd4`).

## Outcome

- `akmctl check` (built from the branch) passes on the repository: 0 errors, 0 warnings; it is also a CI step on Linux.
- Decisions are 14 ADRs in `docs/adr/` of the tool repository, each citing its decision log entries; the generated index passes rule A4.
- The gaps are 10 work items in `docs/work/` on the roadmap (mint_ids, width ZWJ, `check --target`, G3, index column, Windows, and four UI items); the dogfood itself is closed in `docs/work/done/`.
- `go vet ./...` and `go test ./...` pass; CI run 40783801257 passed on Linux, macOS and Windows. Prettier 3.9.10 was not run (3.8.4, installed globally, reports all documents clean).
- The template lessons note is [akm-template-lessons](../../design/akm-template-lessons.md). One tool bug found and fixed (`new work` was not Prettier-clean).
