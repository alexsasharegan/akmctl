# remove

Use it when a requirement in the `Goal` no longer applies: it deletes that one unchecked box and logs it. It loses history.

```
$ akmctl remove 261029-atxk 2 --reason "out of scope after the redesign"
akmctl: refused: remove deletes item 2 ("...") ... This loses history ... pass confirm to accept the loss
$ akmctl remove 261029-atxk 2 --reason "out of scope after the redesign" --confirm
261029-atxk	2026-10-29	active	add json output
changed: true
```

The tool does not track git, so the removed line is recoverable only from a commit. Without `--confirm` the call is `refused` (`reason: "confirm_required"`) and writes nothing. Commit first if the record has uncommitted work you may want back. If you are a person, decide. If you are an agent, tell the user what will go, say it cannot be undone from the tool, and pass `--confirm` only after they decide. `--dry-run` shows the diff and needs no confirm.

Which way out of an unchecked box: `tick --evidence` if the requirement was met, even differently; `remove` if it is gone; `close <id> dropped` if the whole item is abandoned.

`<n>` is the box's position, counted as for `tick` (bullet and ordered boxes share one count); later boxes move up by one, and the response lists the checklist after the write. `--reason` is required, one line. The Log line is `removed "<box text>": <reason>`, quoting the whole text. `--section <name>` picks another section; `--section-etag` guards against a changed list.

Files changed: the one record (the box line, a Log line, and `updated`).

Errors: `invalid_request` for a missing or multi-line reason; `not_found` for a missing section or item; `refused` with `reason: "confirm_required"`, `"ticked_box"` (a ticked box records that work was done), `"closed"` or `"unsupported_markdown"`; `stale`.
