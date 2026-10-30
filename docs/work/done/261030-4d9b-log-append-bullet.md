---
id: "261030-4d9b"
title: "section --mode append: a Log line without a bullet joins the previous bullet"
status: done
created: 2026-10-30
updated: 2026-10-30
closed: 2026-10-30
summary: A Log line appended without a bullet joins the previous bullet; add the bullet in append mode
---

# section --mode append: a Log line without a bullet joins the previous bullet

## Goal

`akmctl section <id> Log --mode append --text "..."` writes the text verbatim as one line (`AppendToSection`, `internal/record/fields.go`). A line without a leading `- ` lands under the previous bullet; Markdown reads it as a lazy continuation, `akmctl check` does not flag it, and Prettier later indents it into that bullet. The skills never say the caller supplies the bullet.

Fix: in append mode for `Log`, add `- ` when the line does not already start with a list marker.

Done when:

- [ ] a bare Log line is written as a `- ` bullet; a line that already starts with `- ` is unchanged
- [ ] a test next to `internal/record/section_test.go` covers both
- [ ] `.agents/skills/akm/references/section.md` says what append does with the bullet

## Outcome

- A bare Log line is written as a `- ` bullet and a line that already has a list marker is unchanged: `TestSetSection` cases "append adds the bullet to a bare Log line" and "append keeps a numbered item" (`internal/record/section_test.go`).
- The test covers both, plus a multi-line text and another section left verbatim: the same test.
- `.agents/skills/akm/references/section.md` and the `mode` row of `docs/spec/akm-api-spec.md` say what append does with the bullet.

## Log

- 2026-10-30: filed.
- 2026-10-30: found while filing 261030-43jw; marked ready.
- 2026-10-30: claimed (claude).
  claimed; implemented in SetSection (internal/record/fields.go), with tests, section.md and the API spec line updated.
