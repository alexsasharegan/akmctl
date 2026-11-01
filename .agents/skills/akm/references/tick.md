# tick

Use it to tick or clear one checkbox of the `Goal` (or another section's checklist) without sending the section again.

```
$ akmctl tick 261029-atxk 2
261029-atxk	2026-10-29	active	add json output
changed: true
$ akmctl tick 261029-atxk 2 --undo
```

`<n>` is the 1-based position of the box in the checklist: the first run of `- [ ]` or `1. [ ]` lines with no blank line between; bullet and ordered boxes share one count (`- [ ] a`, `1. [ ] b` are items 1 and 2). `--section <name>` picks another section (default `goal`). `--evidence "TestX passes"` adds a Log line (`ticked "<box text>": TestX passes`) in the same write: say what proves the box. A box already in the asked state is `changed: false` and writes nothing. To guard against a changed list, pass the `section_etag` from an earlier response (`--section-etag`): a mismatch is `stale`.

Files changed: the one record; one byte of it (plus the Log line with `--evidence`), and `updated` is set to today.

Errors: `not_found` for a missing section or item; `refused` with `reason: "closed"` or `"unsupported_markdown"` (the construct and line are in `details`); `stale`.
