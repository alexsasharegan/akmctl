# set

Use it to change frontmatter fields: status, priority, summary, `blocked_by`. Not the body (use `section`), not `done` or `dropped` (use `close`), not `claimed_by` (use `claim`).

```
$ akmctl set 261029-atxk status=ready priority=1 summary="Add a --json flag"
261029-atxk	2026-10-29	ready	add json output
changed: [priority status summary]
```

Files changed: the one record; `updated` is set to today.

Errors:

- `invalid_request`: `status: must be one of raw, ready, active, blocked, done, dropped`.
- `refused`, `reason: "rule"`: `WK5` blocked without `blocked_by`, `WK4` own parent, `WK9` closed item reopened, `A3` superseded without `superseded_by`.
- `stale`: the file changed since you read it; read it again.
