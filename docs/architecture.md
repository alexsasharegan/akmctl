# Architecture

`akmctl` reads a repository's AKM records (work items, ADRs, the roadmap, the manifest `.layout`, `docs/.layout` by default) and offers one set of operations over them through two bindings: a command line and a loopback HTTP server. The contract is the AKM API spec in `docs/spec/akm-api-spec.md`; the fixtures in `testdata/api/` are its executable form, a copy of `docs/spec/akm-api/.fixtures` kept in step by hand (work item 261029-x4ws).

The words used here are defined in [glossary.md](glossary.md).

## Packages

| Package                | Role                                                                                                                                                                          |
| ---------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `cmd/akmctl`           | The binary; calls `internal/cli`.                                                                                                                                             |
| `internal/frontmatter` | The YAML-lite parser and the byte-splice editor. Fixtures: `testdata/yaml-lite.json`.                                                                                         |
| `internal/record`      | The loader: records, roadmap, manifest, sections, etags, wikilink-to-id resolution.                                                                                           |
| `internal/api`         | Every operation (`version ls search show board find when check commits index new set close claim release migrate init config`), request validation and the response envelope. |
| `internal/check`       | The rules WK, I, A, RM, G, L, S, CL, SF.                                                                                                                                      |
| `internal/index`       | The ADR index generator, stable under Prettier.                                                                                                                               |
| `internal/width`       | Prettier's string-width rule; tables generated from Unicode 17.0.0 (`data/`, `gen/main.go`).                                                                                  |
| `internal/initrepo`    | `init`: the files for the options (docs dir, wikilinks, git), the embedded templates and skills, the no-overwrite writer.                                                     |
| `internal/migrate`     | Declarative layout migration steps (including `mint_ids`), safe application, handoff. No migrations are embedded yet.                                                         |
| `internal/fsx`         | The advisory lock (`.akm.lock` in the manifest folder) and atomic writes.                                                                                                     |
| `internal/server`      | The HTTP binding and the workbench page (`web/`).                                                                                                                             |
| `internal/cli`         | Flags to request parameters; `--format json\|table\|ids`; the `ui` command.                                                                                                   |

## Invariants

- Standard library only in the tool; `go.yaml.in/yaml/v3` is a test oracle and never linked into the binary.
- `init` never overwrites a file; versions in what it writes come from the tool's constants (`check.LayoutVersion`), not from templates. `internal/initrepo/templates/skills` is a copy of `.agents/skills`, kept equal by a test.
- The manifest, `.akm-deny`, `.akm.lock` and `.migrate-log` share one folder (`docs/` by default), found by `record.Locate`: `docs/.layout`, else exactly one `<name>/.layout` among the root's immediate non-hidden subdirectories; two is an error. No code hard-codes `docs/` for them.
- The tool never stages, commits or pushes; it describes a step's commit (`commit` member) and `commits` reads git history, read-only. `close` renames; `new adr` does not regenerate the index (`index` does).
- A write is a compare-and-swap on an etag under the repository lock, written by temp file and rename.
- Unknown frontmatter keys, comments, quoting and line endings survive every edit byte for byte.
- Records the tool writes use bare ids unless the manifest sets `reference_style: wikilink`.
- The CLI and the HTTP binding return the same JSON envelope for the same request.
- `testdata/` is byte-exact: `.gitattributes` turns line-ending conversion off for it.

## Where the standard lives

This repository holds the whole of AKM: the standard (`docs/spec/agent-ready-repo-spec.md`), the guide (`docs/spec/akm-guide.md`), the API spec (`docs/spec/akm-api-spec.md`) with its schemas, validator and fixtures (`docs/spec/akm-api/`), the decisions that shaped them (`docs/adr/`, `docs/design/agent-ready-repo-decision-log.md`) and the closed work (`docs/work/done/`). There is no template repository; `akmctl init` and `akmctl migrate` are the way in (ADR 261029-j0d0). `scripts/vault-migration/` is the history of the first migration, kept as the seed of the `migrate` discovery phase.

## Versions

The tool, the standard (layout) and the API are versioned independently: the tool is 0.x, the layout in the manifest is 0.3.0, the API is 1.0. The tool's version is stamped at build time from the git tag; there is no `VERSION` file.
