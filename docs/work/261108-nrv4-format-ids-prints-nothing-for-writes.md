---
id: "261108-nrv4"
title: format ids prints nothing for writes
status: raw
kind: bug
priority: 3
created: 2026-11-08
updated: 2026-11-08
---

# format ids prints nothing for writes

## Goal

`akmctl new work <slug> --format ids` creates the item but prints nothing; `akmctl set <id> k=v --format ids` also prints nothing. The `ids` format reads like a way to capture the new id in a script, so an agent that gets an empty result may think the call failed and retry, which files a duplicate (this happened on 2026-11-08: one retry left a second item to delete by hand).

Actual: empty output, exit 0, a file written. Expected: the ids of the records the command created or changed, one per line, as `ls --format ids` prints them (or a refusal that says `ids` does not apply, if that is the decision).

Done when:

- [ ] `new` and `set` with `--format ids` print the record id, one per line, or refuse with a message that says what to use instead.
- [ ] A test pins the choice for `new` and `set`, and the help text for `--format` says what `ids` means for a write.
- [ ] The other write commands (`claim`, `close`, `log`, `tick`, `release`, `section`) are checked and the same rule holds.

## Log

- 2026-11-08: filed.
