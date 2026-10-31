---
id: "261030-q03z"
title: commit message format and stamp
status: active
created: 2026-10-30
updated: 2026-10-30
kind: feature
priority: 2
decisions: [261030-tbkh]
summary: A standard commit subject and AKM-* trailers for status steps, returned by the tool as a commit object and checked against the diff.
claimed_by: claude
claimed_at: 2026-10-31T04:08
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
- [ ] A read-only `commits` operation (the ADR picks it over a `check` rule: `check` is a function of the working tree and takes no history) runs `git log` over a range verifies a commit's trailers and subject against its diff: `close` moves a file into `done/`, `claim` sets `active`, and so on. Its rule id and a spec row exist.
- [x] `akm-commit` says to use the returned message and paths verbatim and add code paths; both skill copies are updated if the skill is embedded by then, and the drift test passes.
- [ ] Specification, guide, glossary (commit object, step, trailer) and changelog are updated.
- [ ] Open for later, not in this item: an opt-in `--commit` flag that stages only the reported `paths`.

## Next

Build `commits`: an API operation and CLI command (`akmctl commits [--range <rev-range>]`, default the last 50 commits) that runs `git log` and, for each commit with an `AKM-Step` trailer, reports whether the subject, the two trailers and the diff agree (`close` moves a file into `done/` with status done, `drop` the same with dropped, `claim` sets `active` and `claimed_by`, `file` adds a record, and so on). Write its rules in the spec, add the schema `commits.json` and fixtures (a small git repository built by `make_fixtures.py`), then the glossary, guide and changelog lines. After that, finish the remaining done-when lines and close.

## Log

- 2026-10-30: filed.
- 2026-10-30: refined from a design discussion with the maintainer; format, trailers, structured commit object and no tool commits approved.
- 2026-10-30: claimed (claude).
- 2026-10-30: claimed. ADR 261030-tbkh drafted (`proposed`). The `commit` object is built for `new`, `claim`, `release`, `close` and `set` status changes: code, schemas, fixtures (`set-8` new, five regenerated), tests, spec, the `akm-commit` skill and the reference pages. Left: the `commits` verifier.
