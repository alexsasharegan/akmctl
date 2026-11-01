# log

Use it to add one dated line to the `## Log` of a record.

```
$ akmctl log 261029-atxk tests written
261029-atxk	2026-10-29	active	add json output
changed: true
```

The line is `- <today>: <text>` (`--as-of` sets the day), added after the last Log line. The words after the id are the text: one line, not empty. Use `section` for anything longer. A record with no Log gets one.

Files changed: the one record; `updated` is set to today.

Errors: `invalid_request` for empty or multi-line text; `refused` with `reason: "closed"` or `"unsupported_markdown"`; `stale` (`--section-etag`).
