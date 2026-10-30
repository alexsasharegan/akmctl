# akmctl

The command-line tool for AKM (agent knowledge management), a standard that keeps a repository's decisions, direction and work in flight as Markdown records. The name is a working title.

**Status: 0.2.0 beta, for new (greenfield) repositories.** `akmctl init` creates a repository that is ready for people and AI agents to keep their decisions, direction and work as records. Adopting AKM in an existing repository (`akmctl migrate`) is documented but not supported in this beta. Expect changes before 1.0.0.

The contract is the AKM API spec ([`docs/spec/akm-api-spec.md`](docs/spec/akm-api-spec.md)), with its schemas and fixtures in [`docs/spec/akm-api/`](docs/spec/akm-api/); the standard is [`docs/spec/agent-ready-repo-spec.md`](docs/spec/agent-ready-repo-spec.md) and a short introduction is [`docs/spec/akm-guide.md`](docs/spec/akm-guide.md).

## Quick start

Needs Go 1.27 to build from a clone (a release archive needs nothing).

```sh
go build -o ~/bin/akmctl ./cmd/akmctl   # from a clone of this repository; put the binary on your PATH
mkdir my-project && cd my-project && git init
akmctl init --profile code    # a software repository; or --profile notes for a notes repository / Obsidian vault
akmctl check                  # 0 errors on a new repository
akmctl new work first-task --title "First task"
akmctl ls --status raw,ready,active
```

`init` never overwrites a file, so running it again changes nothing. It writes `AGENTS.md` (the instructions an agent reads first), `CLAUDE.md`, the `work-items` and `record-decision` skills under `.agents/skills/`, the `docs/` layout and a `.gitignore`. Fill in the commands in `AGENTS.md`, then point your agent at the repository. The work cycle is `akmctl claim <id> --as <name>`, work, then `akmctl close <id> done` after writing the `## Outcome`. Nothing stages or commits: you do.

## Layout

- `cmd/akmctl`: the binary.
- `internal/frontmatter`: the YAML-lite parser and byte-splice editor.
- `internal/record`: loads work items, ADRs, the roadmap and the manifest.
- `internal/api`: the operations (`version`, `ls`, `search`, `show`, `board`, `find`, `when`, `check`, `index`, `new`, `set`, `section`, `close`, `claim`, `release`, `init`, `migrate`) and their request checking.
- `internal/check`: the rules (WK, I, A, RM, G, L, S, the claims rules CL and the safety rules SF).
- `internal/index`: the ADR index generator, stable under Prettier.
- `internal/width`: Prettier's string-width rule; tables are generated from the Unicode data in `data/` (`go generate ./internal/width`).
- `internal/migrate`: layout migrations (embedded steps, safe application, handoff).
- `internal/fsx`: the repository lock and atomic writes.
- `internal/server`: the HTTP binding (loopback only, bearer token) and the workbench page.
- `internal/cli`: the command-line binding.
- `testdata`: conformance fixtures shared with the spec (`yaml-lite.json`, `api/`); a copy of `docs/spec/akm-api/.fixtures`, kept in step by hand.
- `docs/spec`: the AKM standard, the guide, the API spec, its schemas and fixtures.
- `docs/design`: the design history (decision log, comparison, reviews, library evaluations), including `vault-migration/`.
- `scripts/vault-migration`: the scripts of the run that migrated the maintainer's vault; history, vault-specific.

## Use

```sh
akmctl ls --status ready,active --sort roadmap,-priority,id
akmctl show 261020-a2b3 --section next
akmctl board --format json
akmctl check            # exit 1 when a rule finds an error
akmctl index            # regenerate docs/adr/README.md
akmctl ui --open        # the workbench in a browser, on 127.0.0.1 only
```

`--format json` prints one JSON document; the default is a plain table. Errors exit 2.

## Project records

This repository keeps its own memory the AKM way, as the first test case of the standard: `AGENTS.md` (start here, for agents), `docs/adr/` (the decisions, with a generated index), `docs/work/` (open and closed work), `docs/ROADMAP.md` and `docs/architecture.md`. `akmctl check` must pass on it.

## Development

```sh
go vet ./... && go test ./...
go build -o /tmp/akmctl ./cmd/akmctl && /tmp/akmctl check
```

Standard library only, apart from `go.yaml.in/yaml/v3` in tests, used as an oracle for the parser.

## Licence

MIT; see `LICENSE`. The Unicode data files in `internal/width/data` are under the Unicode licence.
