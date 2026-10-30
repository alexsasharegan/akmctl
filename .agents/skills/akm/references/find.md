# find

Use it to turn an id into a file path, so you can open it or stage it by path. Read-only.

```
$ akmctl find 261029-atxk
docs/work/261029-atxk-add-json-output.md
```

Files changed: none.

Errors: `invalid_request` for a malformed id, `not_found` when no file has it, `ambiguous` when several do (`details.paths`).
