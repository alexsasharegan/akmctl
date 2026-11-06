# close

Use it when the work is finished (`done`) or abandoned (`dropped`), after writing `## Outcome`.

```
$ akmctl close 261029-atxk done
akmctl: refused: 261029-atxk has no ## Outcome; write it before closing
```

Write the Outcome with `section` (or pass it as `outcome`), then close again. The file moves to `docs/work/done/`, the claim is removed, and `closed` is set to today (for `done` and `dropped` alike). `close` is the only operation that writes `closed`; `set` refuses it. It also empties a non-empty `## Next`: a closed item has no next step, and the old plan would mislead. Update the roadmap if the item is on it.

When open items name the closed one in `blocked_by`, `close` lists them (`dependents` in JSON) after the move:

```
$ akmctl close 261029-atxk done
261029-atxk	2026-10-30	done	add json output
moved to docs/work/done/261029-atxk-add-json-output.md
open items that were waiting on this one (unblock or re-point each):
  261029-m3pq	blocked	Publish the first release
```

`close` never edits them: unblock or re-point each with `set`, in the same commit.

Returns a `commit` (message and paths to stage). Use `commit.message` as the commit message exactly as given, with its `AKM-Step` and `AKM-Item` trailers, and stage `commit.paths` with the code that belongs to the step. Pass the whole message to one `git commit -F -`, not one `-m` per line, or git will not read the trailers. Do not write your own subject. One commit per step.

The move breaks relative Markdown links in and to the item, so `close` rewrites them (the item's own links, and inline links in other records, the roadmap, AGENTS.md and the docs that point at it) and prints them under `rewrote Markdown links` (`links_rewritten` in JSON). Wikilinks, code and anchors are left alone. A file it cannot write is listed as `link_errors`; fix those by hand.

Files changed: the record, renamed into `done/`, with `updated` and `closed` set to today; every file whose links to it were rewritten. `commit.paths` lists them all.

`close done` is refused while the `## Goal` has an empty box (`reason: "unchecked_goal"`, each box named with its line and, when `tick` can address it, its number). Tick each one with its proof (`akmctl tick <id> <n> --evidence "..."`), `remove` it (`references/remove.md`) if the requirement is gone, change the Goal with `section`, or close the item `dropped`. There is no flag to close past it. `dropped` is not checked.

Errors: `refused` with `reason: "empty_outcome"`, `"unchecked_goal"` or `"open_children"` (a parent closes after its children); `not_found`; `stale`.
