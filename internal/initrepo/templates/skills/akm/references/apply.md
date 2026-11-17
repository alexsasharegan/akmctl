# apply

Use it to make several writes as one step: a change that needs more than one write (a section, then a claim; a new ADR, the field that links it, then a close) is made whole or not at all. The writes run in order under one lock. If one fails, the ones before it are undone and every file is as it was. Give it the same JSON file `preview` takes: a list of writes, each `{"op": "section", "id": ..., ...}` with the parameters of that operation. Run `preview` first to see the diffs.

```
$ akmctl apply --writes writes.json
write 1: section
261029-atxk	2026-10-29	ready	add json output
changed: true
write 2: claim
261029-atxk	2026-10-29	active	add json output
$ akmctl apply --writes bad.json
akmctl: not_found: no record with id 261001-zzzz
write 2 failed; the writes before it were undone and the repository is as it was
```

`op` is `new`, `set`, `section`, `tick`, `remove`, `log`, `close`, `claim` or `release`. In a write, `"$new"` stands for the id of the record made by the `new` before it. An `if_match` is checked on the first write to a record; later writes to it follow on from that one. A write may not carry `dry_run`.

Each write that changes a status or a claim prints its own commit message; commit them as the `akm-commit` skill says, one step each, or stage all the paths for the last step when the writes are one logical step.

The undo covers a failure the tool sees (a refusal, a stale record, a file it cannot write). It does not cover a crash or a power loss in the middle of a sequence: after one, read `git status` and restore with git.

Files changed: the files of each write, or none when a write fails.

Errors: `invalid_request` for an empty list, an unknown `op` or a `dry_run` in a write; any error of a write (`stale`, `refused`, `not_found`, `unreadable`) with `details.write`, the position of the write that failed, and `details.rolled_back`. If `rolled_back` is `false`, `details.rollback_errors` lists the files that could not be put back: tell the maintainer and restore them with git.
