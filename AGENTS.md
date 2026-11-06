# akmctl

The command-line tool and API of AKM (agent knowledge management), written in Go for maintainers and agents who keep a repository's decisions, direction and work as Markdown records.

## Commands

- Build: `go build -o /tmp/akmctl ./cmd/akmctl` (Go 1.27; standard library only, plus `go.yaml.in/yaml/v3` in tests).
- Test: `go vet ./... && go test ./...`; add `-race` on Linux or macOS. `gofmt -l .` prints nothing.
- Width tables: `go generate ./internal/width` (Unicode 17.0.0 data in `internal/width/data`).
- Release archives: `scripts/release.sh <version>`; CI runs it on a `v*` tag. The steps are in `docs/release.md`.
- Records check: build the tool as above, then `/tmp/akmctl check` (this repo is its own test case; use the build of the current branch, not an installed copy).
- Format records and schemas: `npm ci` once, then `npm run format:check` and `npm run schemas` (Prettier 3.9.10 and Ajv 8 are dev dependencies in `package.json`); `/tmp/akmctl index` regenerates `docs/adr/README.md`.
- Full check before a commit: the test (with `-race`), records check and format commands above. If the sandbox blocks the Go build cache, Git or the network, ask for the escalation; never skip a check or report one as run that was not.

## Where things live

- `docs/architecture.md`: the map of packages and the invariants.
- `docs/adr/`: decisions; `docs/adr/README.md` is the generated routing index.
- `docs/glossary.md`: the words the project uses and avoids (record, work item, workbench parts, Obsidian Bases).
- `docs/ROADMAP.md`: the order of work, in the Now / Next / Later convention with milestones in Later (ADR 261029-6cp5); finished work leaves it for `done/` and the changelog. Open work: `/tmp/akmctl ls --status ready,active,blocked`.
- `docs/release.md`: how a version is cut (prepare, tag, verify); the maintainer tags.
- `docs/work/`: work items; `docs/work/done/` holds closed ones.
- `.agents/skills/`: the procedures (`akm`, `akm-track-work`, `akm-record-decision`).
- `docs/spec/`: the AKM standard, the guide, the API spec and `akm-api/` (schemas, validator, fixtures). `docs/design/`: the design history. `testdata/api/` is a copy of `docs/spec/akm-api/.fixtures`; the spec copy is the source, so change both by hand (work item 261029-x4ws files the fix).

## Rules

- Before starting, list the open work. Resume the `active` item you hold; if you hold none, take the top `ready` item; if several `active` items are unclaimed, ask. The `claims` extension is on: use `claim` and `release`. Claim with a name unique to your session, `--as <agent>-<model>` (for example `claude-sonnet-5-5`), never a bare `claude` or `codex`, so two sessions never share a claim.
- Work only under a work item. If none fits, create a `raw` one (`akmctl new work <slug>`) and ask which item to work. Before committing a new item, fill its `## Goal` (for a tool failure: the command, the actual and the expected behavior, a checkable done-when); `new` leaves it empty and `check` allows a `raw` item to be.
- A request only to find or recommend work gets an answer, not a claim. A branch and its pull request carry one work item (items you file in passing ride along); once a branch has a claim commit it keeps that item, and a different item means releasing this one and starting a new branch.
- A hard-to-reverse choice, or one that rejects a real alternative, needs an ADR draft with status `proposed`. Never edit the body of an `accepted` ADR; supersede it.
- Update the item's `## Next` and `## Log` in the same commit as the code.
- A change a user or a client can see (a command, the API, the spec, a rule, a fixture) gets an entry under `## Unreleased` in `CHANGELOG.md` in the same commit as the code, as the entries there do: what changed, who is affected, the work item and ADR ids. A fix with no visible change does not.
- Mint ids with `akmctl new`. Without the tool: `id="$(date +%y%m%d)-$(LC_ALL=C tr -dc '0-9a-hjkmnp-tv-z' </dev/urandom | head -c 4)"`, then `ls docs/work docs/work/done docs/adr | grep "^$id-"` must print nothing.
- Use the build of the current branch for work commands: `/tmp/akmctl`, built as in Commands, not an installed `akmctl`. An `unknown command` or an unknown flag usually means a stale binary: run `/tmp/akmctl version` (rebuild if it is older than the branch) before filing it.
- Surface every `akmctl` failure: a crash, a wrong result, an error that does not say what to do next, or a refusal you did not expect. Tell the maintainer in your next message and file it at once as at least a `raw` work item (`akmctl new work <slug>`; read the `similar` list and add to an existing item if one fits). Do this before you work around it. A refusal that did what the spec says and told you why (a cycle, a closed item) is the tool working: no item.
- Code and tests win over prose. When a document is wrong, fix it or say so.
- This file wins over a shared skill where they disagree. The tool never stages or commits; stage by path, one logical commit per step (the `akm-commit` skill is on here). Land work on a branch through a pull request, never by pushing to `main`, unless the maintainer asks you to merge it locally; push the verified branch (the skill's default of never pushing does not apply in this repository) and do not merge the pull request yourself. Keep the configured Git author and pass each commit message whole to one `git commit -F -` (never one `-m` per line, so git reads the trailers; work item 261101-ft79) and end it with `Assisted-by: <agent and model>` directly after the `AKM-` trailers, with no blank line between, and each pull request description with the line `Assisted by <agent>` (ADR 261106-zvpm; never "generated by", no emoji). Tags, releases and visibility changes still need the maintainer.
- Without the tool: decisions are in `docs/adr/README.md`; `grep -lE '^status: ?"?active"?$' docs/work/*.md` lists active work (same for `ready`, `blocked`); close an item by setting `status`, writing `## Outcome`, setting `updated` and `git mv` to `docs/work/done/`. Never hand-edit the ADR index.
- The `akmctl` binary must work on macOS and Linux. Windows is best effort and not required below 1.0.0 (ADR 261029-eegx): CI runs it without blocking.

- Register every new term in `docs/glossary.md` in the same change: a feature, command, field, rule or workbench part. Write "Obsidian Bases" in full, never "base" alone.

## Do not, without asking

- Push a tag, create a release, change repository visibility or settings, or publish: releases and going public are the maintainer's call.
- Add a dependency to the tool (a test-only dependency is fine with a reason in an ADR).
- Put credentials, tokens or private material in any file: the `safety` extension checks known patterns only.
- Edit the body of an accepted ADR, or hand-edit `docs/adr/README.md`.
