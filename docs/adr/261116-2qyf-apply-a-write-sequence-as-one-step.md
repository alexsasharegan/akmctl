---
id: "261116-2qyf"
title: apply makes a write sequence as one step and undoes it on failure
summary: An apply operation makes an ordered list of writes under one lock and, when one fails, restores every file from bytes read before the first; no journal.
status: proposed
date: 2026-11-16
authored_by: agent
---

# apply makes a write sequence as one step and undoes it on failure

## Context

A move in the workbench can need several writes in order: a Goal section, a Next section, a new ADR, the `decisions` change that links it, a release, a status change, a close. The page sent them one request at a time. When a later one failed (a stale record, a refusal, a file that could not be written), the earlier ones stayed: an item with a new Goal and its old status, or a new ADR that nothing linked to (work items 261029-xtqj and 261029-nkqv). The wizard could only report the half and offer to retry the rest. Between two requests another writer could also take the lock and change the record.

Exit criterion 2 for 1.0.0 is that a write that changes several files cannot stop half-way. ADR 261106-60xv decided that git is the move log and that no journal file is added now. ADR 261109-vgg7 added `preview`, which already takes the whole list of a move's writes and runs them on a scratch copy.

Only the server has the code that makes a write, and it holds the lock. The client cannot put a file back: the API has no "write these bytes" call, and should not get one.

## Considered Options

1. **An `apply` operation: the list of writes, made under one lock, with the earlier files restored from memory when one fails** (chosen).
2. Do nothing more than report: keep one request per write, with the retry of the rest. A half-made move stays possible, and the criterion is not met.
3. A journal file: write the intent before the sequence and the completion after it, and replay or reverse on the next start. It also covers a crash in the middle. It adds a file to the layout before the 1.0.0 freeze, and ADR 261106-60xv rejected it for now; the spike 261106-wx2a keeps the question.
4. Make the writes on a scratch copy and then copy the changed files over the repository: the copy step is itself several writes and can stop half-way, so it moves the problem.
5. Undo in the client with compensating writes (set the status back, drop the new ADR): a `close` moves a file and rewrites links, a claim adds a Log line, and a new record cannot be deleted through the API, so the result is never the file as it was.

## Decision

Add the operation `apply` (`POST /apply`; `akmctl apply --writes <file>`). It takes the same request as `preview`: `writes`, an array of `{op, ...parameters}` where `op` is one of `new`, `set`, `section`, `tick`, `remove`, `log`, `close`, `claim` or `release`, and `as_of`. The placeholder `"$new"` and the chaining of `if_match` work as in `preview`: the first write to a record is checked against the etag it carries, and a later write to the same record uses the etag the earlier one produced. A write may not carry `dry_run`.

The operation takes the repository lock once and holds it for the whole sequence. Under the lock, before the first write, it reads into memory every file a write can change: the records and the Markdown files whose links a `close` rewrites. It then makes the writes in order. If one fails, or a `close` reports `link_errors`, it puts the files back: a file that changed or went away is written again from the bytes held, and a file made since is removed. The response is then that write's error with `details.write` (its zero-based position) and `details.rolled_back: true`. If a file cannot be put back, `rolled_back` is `false` and `details.rollback_errors` lists each `{path, error}`.

On success the response is `writes`, one entry per write: the response that operation gives.

Over HTTP the first write to each existing record must carry `if_match`, as the write routes require `If-Match`.

There is no journal. The undo covers a failure the running tool sees. A crash or a power loss in the middle of a sequence leaves the writes made so far, each one whole (every single write is still atomic); git is how that is seen and reversed.

The workbench sends every move of more than one write to `apply`: the wizard, a safe move that releases or claims, and undo.

## Consequences

- A move made through the tool is whole or absent. The wizard's "retry the rest" and the note about an unlinked ADR go away, because there is no half to retry.
- No other tool writer can change a record between two writes of a sequence, since the lock is held throughout. An editor that ignores the lock is still not covered, as for any write.
- The lock is held longer than one file write: for the length of the sequence plus one read of the Markdown files. Other writers wait; on Windows the lock wait is 30 seconds, far above a sequence of a few writes.
- The read before the first write costs one pass over the Markdown files, as `preview` costs one copy. For a repository with thousands of files that is a cost a single write does not have. If it proves slow, the snapshot can be narrowed to the files the writes name without changing the operation.
- A restored file has its content back, not its modification time.
- The API gains an operation (a MINOR addition): a spec section, a schema, fixtures and a CLI and HTTP binding. `details.rolled_back` and `details.rollback_errors` are new `details` keys.
- A plain `close` outside `apply` still reports `link_errors` without undoing the move, as before.
- An agent can make a multi-step change (file a child, set its fields, write its Goal) as one step from the command line.
- Crash safety between writes is still open; the journal spike keeps it.

## Confirmation

`TestApplyRollsBackWhenAWriteFails` (`internal/api`): a sequence that makes an ADR, links it, closes an item and then fails leaves every file as it was. The fixtures `apply-1` to `apply-3`. `rollback.test.mjs` (`internal/server/testdata`): the wizard makes one `apply` call and offers no retry.

## Review

Drafted by an agent under work item 261029-xtqj, whose roadmap line asked for the earlier writes to be restored from memory with no journal. Accepted by the maintainer on merge.

## Revert

Remove the operation, its route and its command, and send the workbench's writes one request at a time again (the commit that added `apply` shows both). Records are not affected: the operation leaves no file of its own.
