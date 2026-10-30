# release

Use it when you stop and will not continue: it drops your claim so another agent can take the item (claims extension).

```
$ akmctl release 261029-atxk
261029-atxk	2026-10-29	active	add json output
```

Rewrite `## Next` first so a cold agent can resume. The status stays as it was; use `set` to move the item back to `ready`.

Files changed: the one record (`claimed_by` and `claimed_at` removed, a Log line added).

Errors: `not_found`; `stale`; `refused` with `reason: "extension_off"`.
