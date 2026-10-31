# close

Use it when the work is finished (`done`) or abandoned (`dropped`), after writing `## Outcome`.

```
$ akmctl close 261029-atxk done
akmctl: refused: 261029-atxk has no ## Outcome; write it before closing
```

Write the Outcome with `section` (or pass it as `outcome`), then close again. The file moves to `docs/work/done/`, the claim is removed, and `closed` is set to today (for `done` and `dropped` alike). `close` is the only operation that writes `closed`; `set` refuses it. Update the roadmap if the item is on it.

When open items name the closed one in `blocked_by`, `close` lists them (`dependents` in JSON) after the move:

```
$ akmctl close 261029-atxk done
261029-atxk	2026-10-30	done	add json output
moved to docs/work/done/261029-atxk-add-json-output.md
open items that were waiting on this one (unblock or re-point each):
  261029-m3pq	blocked	Publish the first release
```

`close` never edits them: unblock or re-point each with `set`, in the same commit.

Returns a `commit` (message and paths to stage). Use `commit.message` as the commit message exactly as given, with its `AKM-Step` and `AKM-Item` trailers, and stage `commit.paths` with the code that belongs to the step. Do not write your own subject. One commit per step.

Files changed: the record, renamed into `done/`; `updated` and `closed` are set to today.

Errors: `refused` with `reason: "empty_outcome"` or `"open_children"` (a parent closes after its children); `not_found`; `stale`.
