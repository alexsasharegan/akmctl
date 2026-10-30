# section

Use it to write a `##` section of a record: `Goal`, `Next`, `Outcome`. `## Log` only accepts appended lines.

```
$ akmctl section 261029-atxk Goal --text "Print results as JSON.

- [ ] flag works"
261029-atxk	2026-10-29	ready	add json output
changed: true
```

`--text -` reads the text from standard input. For `Log`, the mode is `append`: the text is added as one more line, and a first line that does not start with a list marker (`- `, `* `, `+ ` or `1. `) gets a `- ` bullet, so a bare line does not join the previous bullet. In another section, append adds the text as written.

Files changed: the one record; `updated` is set to today.

Errors: `invalid_request` for an empty name or `replace` on `Log`; `refused` with `reason: "closed"` for a record in `done/`; `stale`; `not_found`.
