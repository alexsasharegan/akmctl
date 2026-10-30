# close

Use it when the work is finished (`done`) or abandoned (`dropped`), after writing `## Outcome`.

```
$ akmctl close 261029-atxk done
akmctl: refused: 261029-atxk has no ## Outcome; write it before closing
```

Write the Outcome with `section` (or pass it as `outcome`), then close again. The file moves to `docs/work/done/` and the claim is removed. Update the roadmap if the item is on it.

Files changed: the record, renamed into `done/`.

Errors: `refused` with `reason: "empty_outcome"` or `"open_children"` (a parent closes after its children); `not_found`; `stale`.
