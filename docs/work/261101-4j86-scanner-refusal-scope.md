---
id: "261101-4j86"
title: "scanner refusals: structural or local, and fewer of them"
status: active
created: 2026-11-01
updated: 2026-11-01
kind: feature
parent: 261031-gtw6
summary: Locate reports each refusal as structural or local, and claims blockquote and indented code as opaque blocks, so writes refuse less; per ADR 261101-tr8e
claimed_by: claude
claimed_at: 2026-11-01T07:41
---

# scanner refusals: structural or local, and fewer of them

## Goal

From proposed ADR [261101-tr8e](../adr/261101-tr8e-body-writes-through-the-section-tree.md). A write to a body refuses today for any construct the scanner (`internal/body`) does not read, anywhere in the file. The ADR narrows that: a refusal is _structural_ (it could change where a heading, fence or checklist is) or _local_ (confined to its own lines), and a write refuses only for a structural one anywhere or a local one in the section being edited. Then shrink the list: claim the local constructs the scanner can skip as opaque blocks, so they are not refusals at all. Measured 2026-11-01 on the maintainer's vault: work items 137 files, 18 refused (16 blockquotes, 2 ordered-list tasks), 18 of 568 sections; ADRs 1 of 26 (an HTML line); all notes 31% of files refused (setext or thematic break after text 272, blockquote 227, tab-indented 132, empty list item 46). Done when:

- [ ] `UnsupportedError` carries a kind (structural or local), and `Tree` can tell a body's structural refusals from its local ones with the section each local one is in. The ADR's lists are the starting classification: structural are HTML lines, indented fences, headings or code, tab-indented lines, setext underlines or thematic breaks after text, blocks inside list items, control characters and bare carriage returns; local are blockquotes, table-delimiter lookalikes, ordered-list tasks, unusual task boxes and empty list items. Each refusal in `TestRefusals` is tagged, and a test fails if one is not.
- [ ] A body with a blockquote outside the edited section splits and edits (a tick in a Goal beside a Notes section with a callout), and one inside the edited section is refused with the construct and line.
- [ ] The scanner claims a blockquote region (consecutive `>` lines) as an opaque block of its own kind, so it is no longer refused; the oracle battery, the generator and `FuzzAgainstOracle` agree, and the saved fuzz findings still pass. Indented code blocks are claimed the same way if the battery allows it; each remaining refusal has a reason in the code.
- [ ] `check` reports a body with a refused construct (a warning naming the construct and line), so a maintainer finds it before an agent does; the rule is registered in the spec and the glossary.

## Next

Not started. Read ADR 261101-tr8e (Decision, 'Fail closed, narrowly'), internal/body/scanner.go and TestRefusals. Start with the classification and the Tree check; claim blockquote after.

## Log

- 2026-11-01: filed.
- 2026-11-01: filed after Alexander asked why a blockquote anywhere should block a tick.
- 2026-11-01: census on the maintainer's vault added to the Goal; ADR 261101-tr8e gained its principles and the 5% threshold (to be set by Alexander).
- 2026-11-01: claimed (claude).
