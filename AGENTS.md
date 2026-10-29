# akmctl

The command-line tool and API of AKM (agent knowledge management), written in Go for maintainers and agents who keep a repository's decisions, direction and work as Markdown records.

## Commands

- Build: `go build -o /tmp/akmctl ./cmd/akmctl` (Go 1.27; standard library only, plus `go.yaml.in/yaml/v3` in tests).
- Test: `go vet ./... && go test ./...`; add `-race` on Linux or macOS. `gofmt -l .` prints nothing.
- Width tables: `go generate ./internal/width` (Unicode 17.0.0 data in `internal/width/data`).
- Release archives: `scripts/release.sh <version>`; CI runs it on a `v*` tag.
- Records check: build the tool as above, then `/tmp/akmctl check` (this repo is its own test case; use the build of the current branch, not an installed copy).
- Format records: `npx prettier@3.9.10 --check AGENTS.md "docs/**/*.md"`; `/tmp/akmctl index` regenerates `docs/adr/README.md`.
- Full check before a commit: the test, records check and format commands above.

## Where things live

- `docs/architecture.md`: the map of packages and the invariants.
- `docs/adr/`: decisions; `docs/adr/README.md` is the generated routing index.
- `docs/ROADMAP.md`: the order of work. Open work: `/tmp/akmctl ls --status ready,active,blocked`.
- `docs/work/`: work items; `docs/work/done/` holds closed ones.
- `.agents/skills/`: the procedures (`work-items`, `record-decision`).
- `docs/spec/`: the AKM standard, the guide, the API spec and `akm-api/` (schemas, validator, fixtures). `docs/design/`: the design history. `testdata/api/` is a copy of `docs/spec/akm-api/.fixtures`; the spec copy is the source, so change both by hand (work item 261029-x4ws files the fix).

## Rules

- Before starting, list the open work. Resume the `active` item you hold; if you hold none, take the top `ready` item; if several `active` items are unclaimed, ask. The `claims` extension is on: use `claim` and `release`.
- Work only under a work item. If none fits, create a `raw` one (`akmctl new work <slug>`) and ask which item to work.
- A hard-to-reverse choice, or one that rejects a real alternative, needs an ADR draft with status `proposed`. Never edit the body of an `accepted` ADR; supersede it.
- Update the item's `## Next` and `## Log` in the same commit as the code.
- Mint ids with `akmctl new`. Without the tool: `id="$(date -u +%y%m%d)-$(LC_ALL=C tr -dc '0-9a-hjkmnp-tv-z' </dev/urandom | head -c 4)"`, then `ls docs/work docs/work/done docs/adr | grep "^$id-"` must print nothing.
- Code and tests win over prose. When a document is wrong, fix it or say so.
- This file wins over a shared skill where they disagree. The tool never stages or commits; stage by path, one logical commit per step.
- Without the tool: decisions are in `docs/adr/README.md`; `grep -lE '^status: ?"?active"?$' docs/work/*.md` lists active work (same for `ready`, `blocked`); close an item by setting `status`, writing `## Outcome`, setting `updated` and `git mv` to `docs/work/done/`. Never hand-edit the ADR index.
- The `akmctl` binary must work on macOS, Linux and Windows; CI covers all three.

## Do not, without asking

- Push a tag, create a release, change repository visibility or settings, or publish: releases and going public are the maintainer's call.
- Add a dependency to the tool (a test-only dependency is fine with a reason in an ADR).
- Change `internal/server` while the UI port is being planned against it.
- Put credentials, tokens or private material in any file: the `safety` extension checks known patterns only.
- Edit the body of an accepted ADR, or hand-edit `docs/adr/README.md`.
