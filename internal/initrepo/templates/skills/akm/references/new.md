# new

Use it to create a work item (`work`) or an ADR draft (`adr`). It mints the id; never type one by hand.

```
$ akmctl new work add-json-output
261029-atxk	2026-10-29	raw	add json output
docs/work/261029-atxk-add-json-output.md
```

A work item starts `raw`, an ADR `proposed`. The slug is lowercase kebab-case, at most 60 characters; `--title` overrides the title. Fill the file with `set` and `section`.

Files changed: one new file in the work or ADR folder.

Errors: `invalid_request` for a bad slug, kind or initial field; `not_found` when a `parent` does not resolve.
