---
id: "261030-q03z"
title: commit message format and stamp
status: active
created: 2026-10-30
updated: 2026-10-30
kind: feature
priority: 2
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

- [ ] An ADR (`proposed`) records the format, the closed step set, the trailer names, that the structured `commit` object is the contract with `message` derived from it, and that the tool still does not commit (it narrows ADR 261029-nyv6, it does not supersede it); whether Conventional Commits order (`<step>(<id>)`) was considered is stated.
- [ ] The listed commands return `commit` in the CLI (plain `message` under the record path) and in the API; schema and fixtures cover each step.
- [ ] A `check` rule (or a `log` verb; the ADR picks) verifies a commit's trailers and subject against its diff: `close` moves a file into `done/`, `claim` sets `active`, and so on. Its rule id and a spec row exist.
- [ ] `akm-commit` says to use the returned message and paths verbatim and add code paths; both skill copies are updated if the skill is embedded by then, and the drift test passes.
- [ ] Specification, guide, glossary (commit object, step, trailer) and changelog are updated.
- [ ] Open for later, not in this item: an opt-in `--commit` flag that stages only the reported `paths`.

## Next

Nothing started. Take it after 261030-dw71 (it shares the API schema and fixtures with akc4 and dw71).

## Log

- 2026-10-30: filed.
- 2026-10-30: refined from a design discussion with the maintainer; format, trailers, structured commit object and no tool commits approved.
- 2026-10-30: claimed (claude).
