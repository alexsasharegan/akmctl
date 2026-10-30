# show

Use it to read one record by id. Read-only. The JSON form carries the `etag`.

```
$ akmctl show 261029-atxk
id: 261029-atxk
title: add json output
status: raw
...
## Goal

## Log

- 2026-10-29: filed.
```

Files changed: none.

Errors:

- `invalid_request`: the id has the wrong form.
- `not_found`: no record has this id.
- `ambiguous`: two records share the id (rule I1); `details.paths` lists them.
