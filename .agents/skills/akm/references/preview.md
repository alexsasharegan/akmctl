# preview

Use it to see what a sequence of writes would do before you make any: each write's diff is as it would be after the earlier ones, where `--dry-run` on one write shows it against the file as it is now. It writes nothing. Give it a JSON file holding a list of writes, each `{"op": "section", "id": ..., ...}` with the parameters of that operation.

```
$ akmctl preview --writes writes.json
write 1: section
--- docs/work/261029-atxk-add-json-output.md
+++ docs/work/261029-atxk-add-json-output.md
```

`op` is `new`, `set`, `section`, `tick`, `remove`, `log`, `close`, `claim` or `release`. In a write, `"$new"` stands for the id of the record made by the `new` before it; the id in the preview is not the one the real write will mint.

Files changed: none.

Errors: `invalid_request` for an empty list or an unknown `op`; any error of a write (`stale`, `refused`, `not_found`) with `details.write`, the position of the write that failed.
