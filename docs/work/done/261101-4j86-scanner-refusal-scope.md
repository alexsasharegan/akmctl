---
id: "261101-4j86"
title: "scanner refusals: structural or local, and fewer of them"
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
kind: feature
parent: 261031-gtw6
summary: Locate reports each refusal as structural or local, and claims blockquote and indented code as opaque blocks, so writes refuse less; per ADR 261101-tr8e
---

# scanner refusals: structural or local, and fewer of them

## Goal

From proposed ADR [261101-tr8e](../../adr/261101-tr8e-body-writes-through-the-section-tree.md). A write to a body refuses today for any construct the scanner (`internal/body`) does not read, anywhere in the file. The ADR narrows that: a refusal is _structural_ (it could change where a heading, fence or checklist is) or _local_ (confined to its own lines), and a write refuses only for a structural one anywhere or a local one in the section being edited. Then shrink the list: claim the local constructs the scanner can skip as opaque blocks, so they are not refusals at all. Measured 2026-11-01 on the maintainer's vault: work items 137 files, 18 refused (16 blockquotes, 2 ordered-list tasks), 18 of 568 sections; ADRs 1 of 26 (an HTML line); all notes 31% of files refused (setext or thematic break after text 272, blockquote 227, tab-indented 132, empty list item 46). Done when:

- [x] `UnsupportedError` carries a kind (structural or local), and `Tree` can tell a body's structural refusals from its local ones with the section each local one is in. The ADR's lists are the starting classification: structural are HTML lines, indented fences, headings or code, tab-indented lines, setext underlines or thematic breaks after text, blocks inside list items, control characters and bare carriage returns; local are blockquotes, table-delimiter lookalikes, ordered-list tasks, unusual task boxes and empty list items. Each refusal in `TestRefusals` is tagged, and a test fails if one is not.
- [x] A body with a blockquote outside the edited section splits and edits (a tick in a Goal beside a Notes section with a callout), and one inside the edited section is refused with the construct and line.
- [x] The scanner claims a blockquote region (consecutive `>` lines) as an opaque block of its own kind, so it is no longer refused; the oracle battery, the generator and `FuzzAgainstOracle` agree, and the saved fuzz findings still pass. Indented code blocks are claimed the same way if the battery allows it; each remaining refusal has a reason in the code.
- [x] `check` reports a body with a refused construct (a warning naming the construct and line), so a maintainer finds it before an agent does; the rule is registered in the spec and the glossary.

## Next

Done and verified; close with the Outcome.

## Outcome

- Kinds: `UnsupportedError.Structural` and `IsStructural` (`internal/body/body.go`); `Tree.Refusals` lists the local constructs a split stepped over. `TestEveryRefusalIsClassified` reads `scanner.go` and fails if a refused construct is in neither list, so a new refusal cannot be added untagged. One change to the ADR's draft list: a table-delimiter lookalike is structural, because the oracle showed (generator, tolerant comparison) that it can turn the line before it from a task into plain text.
- Scope: `TestLocalConstructScope` ticks a Goal box and appends a Log line beside a Notes section with a local construct, and refuses a replace of the section that holds it (local refusal, line 9). `TestTreeAgainstOracle`: 11,412 documents split, 6,627 edits confined to their section, 330 refused for a local construct in their own section.
- Blockquote: claimed as an opaque `Quote` block (one block per run of `>` lines), so it is no longer a refusal anywhere; nothing in it is a heading, fence or task (`TestQuoteIsOpaque`). The oracle adapter ignores checkboxes inside a quote to match the claim. The strict battery, `TestGenerated` (30,000 documents) and a 120-second `FuzzAgainstOracle` (11 million executions) agree; `FuzzTolerantAgainstOracle` (new, 90 seconds, 3 million) and `TestTolerantAgainstOracle` (30,000, 2,300 with a stepped-over construct) find no difference in headings and no task the oracle does not know.
- Indented code: tab-indented lines were tried as plain text and the battery refused it at once (saved finding `672f74408e54df82`: the oracle reads `* a` then ` \t* [X]` as a task), so they stay a structural refusal with that reason in the code. Space-indented code is already ignored.
- `check` rule B1 (`internal/check/check.go`, spec table, `TestB1Findings`, table rows, fixture `check-8` in both copies) warns with the construct and file line, once for a structural construct and per section for a local one; blockquotes and closed items are not reported. This repository: 0 findings.
- Effect on the vault (counts only, 2026-11-01): work items, 138 files all split, 2 of 435 sections hold a local construct (both ordered-list tasks), where 18 of 138 files were refused whole before; ADRs, 25 of 26 split, the one structural refusal an HTML line. Follow-ups already filed: 261101-xher (ordered-list tasks) and 261101-v65x (HTML lines).

`gofmt -l .` (empty), `go vet ./...`, `go test -race ./...`, the records check (0 errors, 0 warnings), `npm run format:check` and `npm run schemas` (80 cases, 0 failures) pass; `go list -deps ./cmd/akmctl | grep -c goldmark` prints 0.

Not done: indented-code claiming beyond the spaces case; per-line refusals inside one section are reported but a tolerant scan still stops at the first structural one; the `tick`, `log` and `section` changes themselves (261101-ga4s, 261101-q72k), which are now unblocked.

Decision: none beyond ADR 261101-tr8e; the table-delimiter reclassification is recorded here.

## Log

- 2026-11-01: filed.
- 2026-11-01: filed after Alexander asked why a blockquote anywhere should block a tick.
- 2026-11-01: census on the maintainer's vault added to the Goal; ADR 261101-tr8e gained its principles and the 5% threshold (to be set by Alexander).
- 2026-11-01: claimed (claude).
- 2026-11-01: built the kinds, the tolerant scan, opaque quotes and rule B1; tab-indented lines tried and kept as a refusal; 2 of 435 vault work-item sections still hold a local construct.
