# claim

Use it to take a `ready` item before working it (claims extension). Sets `status: active`, `claimed_by` and `claimed_at`, and appends a Log line.

```
$ akmctl claim 261029-atxk --as demo-agent
261029-atxk	2026-10-29	active	add json output
$ akmctl claim 261029-atxk --as other
akmctl: refused: 261029-atxk is claimed by demo-agent
```

If someone else holds a fresh claim, take another item. `--takeover` is for a stale claim, with the maintainer's say-so.

Returns a `commit` (message and paths to stage) for the commit; the optional `akm-commit` skill says how to use it.

Files changed: the one record.

Errors: `refused` with `reason: "claimed"` or `"extension_off"`; `invalid_request` without `--as`; `not_found`; `stale`.
