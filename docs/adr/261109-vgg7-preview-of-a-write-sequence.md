---
id: "261109-vgg7"
title: a preview operation applies a write sequence on a scratch copy
summary: A read-only preview operation applies a list of writes on a scratch copy and returns each one's dry-run response as it is after the earlier ones.
status: proposed
date: 2026-11-09
---

# a preview operation applies a write sequence on a scratch copy

## Context

A move in the workbench can need several writes in order: a Goal section, a Next section, a new ADR, a `decisions` change, a release, a status change, a close. The wizard shows the diff of every write before any is made (work item 261029-g8x2). Each diff comes from the write's own `dry_run`, which reads the file as it is on disk. Write 2's dry run therefore ignores write 1, and the review shows a diff the real write will not make. The dry run also returns the current etag, not the etag the write would produce, so a client cannot chain dry runs. The wizard works around this by not previewing the writes that depend on an earlier one (a claim after a status change, a `set` that names the new ADR).

The review needs the answer to "what does write n do after writes 1 to n-1". Only the server has the code that makes a write, and ADR 261101-tr8e keeps that code in one place. `migrate plan` already solves the same problem for migration steps: it applies each step to a scratch copy of the Markdown files after planning it, so a step that depends on an earlier one is planned as it will run.

## Considered Options

1. **A `preview` operation: an ordered list of writes, applied on a scratch copy** (chosen).
2. Apply the writes to the file text in the browser and diff there: no API change, but it copies the section and field logic of the server into JavaScript, which can drift from what the server writes, and ADR 261101-tr8e exists to prevent that.
3. Let a dry run take a base text and return the would-be text, so a client can chain the calls: a smaller change, but every client must carry the text between calls, the response grows by a whole file per write, and a write that moves a file (`close`) or makes one (`new`) has no single base text.
4. Show one combined diff of the first file against the last: it loses the per-write review and still needs the writes applied somewhere.

## Decision

Add the operation `preview` (`POST /preview`; `akmctl preview --writes <file>`). It takes `writes`, an array of objects, each `{op, ...parameters}` where `op` is one of `new`, `set`, `section`, `tick`, `remove`, `log`, `close`, `claim` or `release` and the other members are the parameters of that operation, and `as_of`, the clock for any write that has none. It never writes to the repository.

The server copies the manifest and the Markdown files to a scratch directory, as `migrate plan` does, then handles the writes in order. For each it runs the operation as a `dry_run` on the scratch copy (for the response and its `diff`), then makes the write there, so the next one sees it. The response is `writes`, an array with one entry per write: the response object the operation gives for a dry run, as it would be after the earlier writes. A write's `if_match` is checked as given for the first write to a record; for a later write to the same record it is replaced by the etag the previous write on the scratch copy produced, so a client sends the etag it read once.

The placeholder `"$new"` in a string of a write stands for the id of the record made by the most recent `new` before it. The id made in the preview is not the id the real write will mint; the diff shows the provisional id, and the entry of the `new` write says so in a `note`.

A write that fails stops the preview. The response is that write's error with `details.write` set to its zero-based position, so the client shows the error where the write would have failed. A `stale` first write keeps the code `stale`.

The scratch copy holds the manifest and Markdown files only, so a write that needs another file fails there as it would for real. The operation makes no claim about files other than the records it writes.

## Consequences

- The review shows each diff as it will apply after the earlier writes, and every write can be previewed, including a claim after a status change and a `set` that names a new ADR. The workbench's `previewable` goes away, apart from the creation of the ADR file, which still has no earlier file to diff against.
- One request replaces one request per write, and the scratch copy is made once.
- The API gains an operation: a spec section, a schema, fixtures and a CLI and HTTP binding. It is read-only, so it needs no `if_match` rule and no lock on the real repository.
- The scratch copy is made per request. For a repository with thousands of Markdown files that is a cost the per-write dry run did not have; the workbench's move is one request on a click. If it proves slow, a cache keyed on the repository state can come later without changing the operation.
- Any other client that needs a dry run of a sequence (an agent planning a multi-step change) gets it from the same operation.
- A time-of-check gap remains as before: the real writes happen later, and a change made in between is `stale` at write time.
