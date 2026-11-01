---
id: "261031-mr2q"
title: claimed_at is UTC but nothing an agent reads says so
status: raw
created: 2026-10-31
updated: 2026-11-01
kind: chore
summary: claimed_at is UTC, stated in the spec but not in the guide or skills agents read; fix by documenting it, with a Z suffix recorded as the alternative
decisions: [261029-5pe1, 261101-gw7b]
---

# claimed_at is UTC but nothing an agent reads says so

## Goal

`claimed_at` is a UTC timestamp (`YYYY-MM-DDTHH:MM`, no `Z` or offset), while `created`, `updated` and Log dates are the writer's local calendar day. The UTC rule is stated in the spec (`agent-ready-repo-spec.md`, `akm-api-spec.md`), the glossary, the CHANGELOG and a comment in `internal/api/ops_write.go`, but not in `docs/spec/akm-guide.md`, `.agents/skills/akm/references/claim.md`, the `akm` or `akm-track-work` skills, or ADR 261028-mnc0. An agent that reads only the skills sees a bare time beside local dates; on 2026-10-31 one took `claimed_at: 2026-11-01T05:02` (22:02 local, MST) for a clock error. The tool itself is consistent: it writes and compares in UTC, so the 24-hour staleness rule is unaffected. The risk is to readers.
Read first: ADR 261029-5pe1 (the date and time conventions this documents) and the proposed ADR 261101-gw7b, which adds a third convention (Log lines carry local time with an offset). The wording written here should fit the table in that ADR; the `claimed_at` rule itself does not change.

## Options

1. **Say it where agents read.** Keep the format. State "UTC" in the guide's field description, in `references/claim.md` (and the `internal/initrepo/templates/skills` copy, which a test keeps equal to `.agents/skills`), and in the `akm-track-work` skill's claim step. No format change, no ADR, no fixture or `CL1` change.
2. **Write a `Z` suffix** (`2026-11-01T05:02Z`). The value then says what it is wherever it is read. Costs: a layout and spec change; `reTS` in `internal/check/check.go` (CL1), the fixtures in `docs/spec/akm-api/.fixtures` and `testdata/api/` (both, by hand, see 261029-x4ws), the UI's `.slice(0, 10)` (still fine), the yaml-lite fixture, existing records in the wild (a migration step or a tolerant reader), and a new ADR that supersedes the claim wording in 261028-mnc0, since that ADR is accepted.

## Decision

Option 1 now. It removes the ambiguity for readers at no compatibility cost. Option 2 is recorded as the stronger fix and is not rejected: take it if a real consumer mixes up the two clocks despite the documentation, or at a layout version bump. It needs an ADR first.

Done when:

- [ ] `docs/spec/akm-guide.md`, `.agents/skills/akm/references/claim.md` and `.agents/skills/akm-track-work/SKILL.md` say `claimed_at` is UTC and that dates beside it are local days.
- [ ] The embedded copy under `internal/initrepo/templates/skills` matches (`go test ./internal/initrepo`).
- [ ] `akmctl check` and `npm run format:check` pass.

## Next

Not started. Edit the three files named in Goal (option 1), mirror the skill changes into internal/initrepo/templates/skills, run go test ./... and the records check, and close with Decision: none (option 2 stays recorded here, not rejected).

## Log

- 2026-10-31: filed.
- 2026-10-31: filed after an agent mistook a UTC claimed_at for a clock error (found while resuming 261031-yw5v).
