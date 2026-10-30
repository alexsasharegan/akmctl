# ls

Use it to list records. First step of every task: `akmctl ls --status ready,active,blocked`. Read-only.

```
$ akmctl ls
261029-atxk	2026-10-29	raw	add json output
```

Flags: `--type work|adr`, `--status a,b`, `--parent <id>`, `--blocked`, `--scope <path>`, `--since <date>`, `--cycle <label>`, `--sort`, `--limit`. `--format ids` prints ids only; `--format json` adds the `etag` and the fields.

Files changed: none.

Errors: `invalid_request` for an unknown flag or a `type` or `status` outside its enum.
