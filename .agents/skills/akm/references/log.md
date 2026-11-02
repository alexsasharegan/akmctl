# log

Use it to add one time-stamped line to the `## Log` of a record.

```
$ akmctl log 261029-atxk "tests written"
261029-atxk	2026-10-29	active	add json output
changed: true
```

The line is `- 2026-11-01 16:02-07:00: <text>`: the tool stamps your local time and offset (`--as-of` sets the clock for tests), added after the last Log line. `--at "2026-11-01 11:05"` says when the thing happened; the line ends `(at 11:05-07:00)` and the stamp stays the write time. Do not type a time into the text yourself. An `at` in the future is refused. Quote the whole message as one argument: an unquoted `;`, `&`, `|` or `>` ends it early in the shell, and what follows runs as another command (several words are still joined into the line, with a warning on stderr that shows what was written). The text is one line, not empty. Use `section` for anything longer. A record with no Log gets one.

Files changed: the one record; `updated` is set to today.

Errors: `invalid_request` for empty or multi-line text; `refused` with `reason: "closed"` or `"unsupported_markdown"`; `stale` (`--section-etag`).
