---
id: "261030-q03z"
title: commit message format and stamp
status: done
created: 2026-10-30
updated: 2026-10-30
closed: 2026-10-30
kind: feature
priority: 2
decisions: [261030-tbkh]
summary: A standard commit subject and AKM-* trailers for status steps, returned by the tool as a commit object and checked against the diff.
---

# commit message format and stamp

## Goal

Commits for AKM steps follow one format, and the tool produces it. `claim`, `release`, `close` and `new` (and `block` and `unblock`, if they exist as commands) return a `commit` object: `step`, `item`, `detail`, `paths` (the files the step touched) and `message`. The agent commits with that message verbatim and adds the code paths it changed. The tool still never stages or commits.

Format:

- Subject: `<id>(<step>): <detail>`, for example `261030-akc4(claim): start`. Steps: `file`, `ready`, `claim`, `release`, `block`, `unblock`, `close done`, `close dropped` (the closed set is settled in the ADR).
- Trailers: `AKM-Step: <step>` and `AKM-Item: <id>`. The body may carry free text; the subject and trailers stay as returned.

Why: today a commit carries no sign of whether an AKM step made it, subjects drift ("close", "mark ready", "file the ..."), and the tool cannot commit for the reasons in ADR 261029-nyv6 (it cannot know which code belongs to the step; hooks, signing and dirty trees become its error paths; the HTTP server would commit on a click). A trailer is self-reported, so a rule checks it against the diff.

Done when:

- [x] An ADR (`proposed`, 261030-tbkh) records the format, the closed step set, the trailer names, that the structured `commit` object is the contract with `message` derived from it, and that the tool still does not commit (it narrows ADR 261029-nyv6, it does not supersede it); whether Conventional Commits order (`<step>(<id>)`) was considered is stated.
- [x] The listed commands return `commit` in the CLI (plain `message` under the record path) and in the API; schema and fixtures cover each step.
- [x] A read-only `commits` operation (the ADR picks it over a `check` rule: `check` is a function of the working tree and takes no history) runs `git log` over a range verifies a commit's trailers and subject against its diff: `close` moves a file into `done/`, `claim` sets `active`, and so on. Its rule id and a spec row exist.
- [x] `akm-commit` says to use the returned message and paths verbatim and add code paths; both skill copies are updated if the skill is embedded by then, and the drift test passes.
- [x] Specification, guide, glossary (commit object, step, trailer) and changelog are updated.
- [x] Open for later, not in this item: an opt-in `--commit` flag that stages only the reported `paths`.

## Outcome

- ADR `261030-tbkh` (`proposed`) records the format, the closed set of nine steps, the trailers, the structured `commit` object as the contract with `message` derived from it, the choice of an operation over a `check` rule, and that the tool still never commits (it narrows `261029-nyv6`). Listed in `decisions`.
- `new`, `claim`, `release`, `close` and a `set` that changes a work item's status return `commit` (`internal/api/commit.go`); the CLI prints the message and the paths to stage. Fixtures: `claim-1`, `release-1`, `close-2` to `close-6` carry it and `set-8` is new; `new` is tested for shape (`internal/api/commit_test.go`).
- `akmctl commits [--range] [--limit]` (`internal/api/ops_commits.go`, `GET /commits`, schema `commits.json`) reads `git log` and applies CM1 (subject), CM2 (trailers agree with the subject) and CM3 (the diff does what the step says). Tests build real git repositories (`internal/api/commits_test.go`): the tool's own messages for claim, release, a ready change, new and close are clean; drop is clean; each rule is broken at least once; ordinary commits are skipped; a range starting with `-` or holding shell characters, a non-git directory and a limit over 500 are refused. No fixture cases: a fixture directory cannot hold a git history; the API spec (section 10) says so.
- `akm-commit` says to use `message` and `paths` verbatim and to run `commits` before a push; the `akm` skill has `references/commits.md`; both copies are equal and the drift tests pass. Specification (sections 1, 4.19, 5, 6, 10), guide, architecture, glossary (commit object, step, trailer, step commit) and changelog are updated.
- Checks: `go vet`, `go test ./...`, `gofmt -l`, `akmctl check`, `npm run format:check` and `npm run schemas` (20 schemas, 68 cases).
- Not done, by design: the opt-in `--commit` flag that stages the reported paths (a later item if the manual step proves annoying), and a `check` rule that warns about hand-written subjects.

Decision: ADR 261030-tbkh, proposed; the maintainer accepts it.

## Log

- 2026-10-30: filed.
- 2026-10-30: refined from a design discussion with the maintainer; format, trailers, structured commit object and no tool commits approved.
- 2026-10-30: claimed (claude).
- 2026-10-30: claimed. ADR 261030-tbkh drafted (`proposed`). The `commit` object is built for `new`, `claim`, `release`, `close` and `set` status changes: code, schemas, fixtures (`set-8` new, five regenerated), tests, spec, the `akm-commit` skill and the reference pages. Left: the `commits` verifier.
- 2026-10-30: `commits` verifier built and tested; closed.
