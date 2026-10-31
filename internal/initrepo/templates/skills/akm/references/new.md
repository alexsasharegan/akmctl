# new

Use it to create a work item (`work`) or an ADR draft (`adr`). It mints the id; never type one by hand.

```
$ akmctl new work add-json-output
261029-atxk	2026-10-29	raw	add json output
docs/work/261029-atxk-add-json-output.md
```

When open work items resemble the new one, the command lists them under the path (`similar` in JSON), most shared words first; no match prints nothing extra:

```
$ akmctl new work json-output-for-board
261030-yct9	2026-10-30	raw	json output for board
docs/work/261030-yct9-json-output-for-board.md
similar open items (add to one instead of filing, if it fits):
  261029-ab12	active	Add JSON output to ls
```

Read the list: if an item covers the work, add to it. `new` never refuses because of a match.

A work item starts `raw`, an ADR `proposed`. The slug is lowercase kebab-case, at most 60 characters; `--title` overrides the title. Fill the file with `set` and `section`.

Returns a `commit` (message and paths to stage) for the commit; the optional `akm-commit` skill says how to use it.

Files changed: one new file in the work or ADR folder.

Errors: `invalid_request` for a bad slug, kind or initial field; `not_found` when a `parent` does not resolve.
