# section

Use it to write a `##` section of a record: `Goal`, `Next`, `Outcome`. `## Log` only accepts appended lines.

```
$ akmctl section 261029-atxk Goal --text "Print results as JSON.

- [ ] flag works"
261029-atxk	2026-10-29	ready	add json output
changed: true
```

`--text -` reads the text from standard input. For `Log`, the mode is `append`.

Files changed: the one record; `updated` is set to today.

Errors: `invalid_request` for an empty name or `replace` on `Log`; `refused` with `reason: "closed"` for a record in `done/`; `stale`; `not_found`.
