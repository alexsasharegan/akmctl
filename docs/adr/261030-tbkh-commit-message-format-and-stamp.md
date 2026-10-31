---
id: "261030-tbkh"
title: AKM step commits use one subject and trailers, returned by the tool as a commit object
summary: AKM step commits use one subject and AKM-Step and AKM-Item trailers; the tool returns them as a commit object and still never commits.
status: proposed
date: 2026-10-30
---

# AKM step commits use one subject and trailers, returned by the tool as a commit object

## Context

Commits for AKM steps (claim, release, close, filing an item) carry no sign of whether an AKM step made them, and their subjects drift ("close", "mark ready", "file the ..."). A reader of `git log` cannot tell a step commit from any other, and nothing can check one. ADR 261029-nyv6 decided that the tool never stages or commits and that only the `akm-commit` skill does; that holds here, because the tool cannot know which code belongs to a step, hooks, signing and dirty trees would become its error paths, and the HTTP server would commit on a click.

## Considered Options

- Leave it to the skill's prose (`<id>: <what happened>`). Rejected: the format drifts, and an agent builds each message on its own.
- Have the tool commit. Rejected, for the reasons above; kept open as an opt-in `--commit` flag in a later item that stages only the reported paths.
- A plain message string only. Rejected: the touched paths have no text form, and a check would have to parse a subject.
- Conventional Commits order `<step>(<id>): ...`. Considered; the id first reads better in `git log --oneline`, where every line starts with the id as today, and no Conventional Commits tooling is in use. A change of order is one regular expression.
- The tool returns a structured `commit` object, and `message` is derived from it.

## Decision

The tool returns a `commit` object from the operations that make an AKM step, and still never stages or commits. It narrows ADR 261029-nyv6 only by adding output; that ADR stands.

The object is `{step, item, detail, paths, message}`.

- `step` is one of a closed set: `file` (`new`), `ready`, `start`, `block`, `unblock` (a `set` that changes a work item's `status`), `claim`, `release`, `close`, `drop` (`close` with status `dropped`). A `set` that changes no status returns no `commit`; neither do `section` (a `## Next` or `## Log` checkpoint is not a status change) or `index`.
- `message` is the subject `<id>(<step>): <detail>`, a blank line and the trailers `AKM-Step: <step>` and `AKM-Item: <id>`. The caller commits with it as it is, may put a body between the subject and the trailers, and may add other trailers (`Co-Authored-By`) after them.
- `paths` are the repository-relative files the step wrote (a `close` lists the old and the new path). The caller stages those and the code it changed.
- `detail` is derived: the title for `file`, `claimed by <as>` for `claim`, `released`, `ready`, `started`, `blocked on <first blocked_by entry>`, `unblocked`, `done`, `dropped`.

The structured object is the contract, in the API schema and fixtures; `message` is derived from it by the tool, so the format has one source. A trailer is self-reported, so a later work step (`commits`, a read-only operation that runs `git log` over a range, because `check` is a function of the working tree and takes no history) verifies each commit that carries `AKM-Step` against its diff. Its rules are written with that operation, not here.

## Consequences

- One additive member in the responses of `new`, `set` (status changes), `claim`, `release` and `close`; clients that ignore it are unaffected. The API stays 1.0.
- `akm-commit` says to use `message` and `paths` verbatim.
- The `commits` operation needs `git` on the path, as `migrate` already does; it is the second operation that touches git, and read-only.
- A repository with its own commit conventions can ignore the object; nothing in `check` requires it.

## Confirmation

The conformance fixtures for `claim`, `release` and `close` carry the `commit` object, and `new` is tested for shape. A test fails if a subject or trailer in `message` disagrees with `step` and `item`.

## Review

Proposed 2026-10-30 after a design discussion with the maintainer; accepted when the maintainer merges or sets `accepted`.

## Revert

Remove the member from the schemas and fixtures and from the five operations; the `akm-commit` skill returns to its `<id>: <what happened>` message.
