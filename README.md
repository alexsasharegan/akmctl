# akmctl

The command-line tool for AKM (agent knowledge management), a standard that keeps a repository's decisions, direction and work in flight as Markdown records. The name is a working title.

Status: under construction. The contract is the AKM API spec (`akm-api-spec`), with schemas and fixtures, kept in the maintainer's notes until a public home exists. Nothing here is released.

## Layout

- `cmd/akmctl`: the binary.
- `internal/frontmatter`: the YAML-lite parser and byte-splice editor.
- `testdata`: conformance fixtures shared with the spec.

## Development

```sh
go test ./...
```

Standard library only, apart from `go.yaml.in/yaml/v3` in tests, used as an oracle for the parser.
