# search

Use it to find records by words, for example before opening a new ADR. Read-only.

```
$ akmctl search json
261029-atxk	2026-10-29	ready	add json output
```

Files changed: none.

Errors: `invalid_request` when the text is missing. No match is not an error: the list is empty.
