# when

Use it to read the date inside an id: its calendar date, weekday and age. The id need not exist as a record. Read-only.

```
$ akmctl when 261029-atxk
2026-10-29 Thursday, 0 days ago
```

Files changed: none.

Errors: `invalid_request` when the argument does not match the id pattern.
