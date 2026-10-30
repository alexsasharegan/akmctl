# version

Use it to learn the tool version, the API version and which extensions the repository keeps. Read-only.

```
$ akmctl version
akmctl 0.0.0-dev+4ca0f0b (API 1.0)
```

Use `--format json` for the full object: `layout`, `extensions`, folders, `reference_style`.

Files changed: none.

Errors: none for the version itself; the repository part is left out when the manifest is missing.
