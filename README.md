# akmctl

The command-line tool for AKM (agent knowledge management), a standard that keeps a repository's decisions, direction and work in flight as Markdown records. The name is a working title.

Status: under construction. The contract is the AKM API spec (`akm-api-spec`), with schemas and fixtures, kept in the maintainer's notes until a public home exists. Nothing here is released.

## Layout

- `cmd/akmctl`: the binary.
- `internal/frontmatter`: the YAML-lite parser and byte-splice editor.
- `internal/record`: loads work items, ADRs, the roadmap and the manifest.
- `internal/api`: the operations (`version`, `ls`, `search`, `show`, `board`, `find`, `when` so far) and their request checking.
- `internal/cli`: the command-line binding.
- `testdata`: conformance fixtures shared with the spec (`yaml-lite.json`, `api/`).

## Use

```sh
akmctl ls --status ready,active --sort roadmap,-priority,id
akmctl show 261020-a2b3 --section next
akmctl board --format json
```

`--format json` prints one JSON document; the default is a plain table. Errors exit 2.

## Development

```sh
go test ./...
```

Standard library only, apart from `go.yaml.in/yaml/v3` in tests, used as an oracle for the parser.
