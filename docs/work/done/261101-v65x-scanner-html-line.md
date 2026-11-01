---
id: "261101-v65x"
title: "scanner: an HTML line in a record is refused"
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
kind: feature
parent: 261031-gtw6
summary: Locate refuses any HTML line (structural); 1 of the vault's ADRs has one. Decide whether a comment or inline tag can be claimed, or the refusal stays
---

# scanner: an HTML line in a record is refused

## Goal

Found by the refusal census in ADR [261101-tr8e](../../adr/261101-tr8e-body-writes-through-the-section-tree.md) (principle 4). The scanner refuses any HTML line (`<div>`, `<!-- c -->`, `<br>`) because an HTML block can swallow a `## ` heading until a blank line, which makes it structural: a write refuses whenever the body has one. One of the maintainer's vault ADRs has one. Done when:

- [x] A decision, recorded in this item's Outcome: claim the safe cases (a one-line comment `<!-- c -->` or an inline tag inside a text line, which do not open an HTML block) and keep refusing the rest, or keep refusing all HTML with a reason in the code.
- [x] If claimed: the oracle battery, the generator and the fuzz target agree on the claimed cases, and `TestRefusals` keeps a row for each case still refused.

## Next

Done; ready to close.

## Outcome

Claimed the safe cases; the rest stay structural refusals.

- Claimed: a comment that ends on its own line (`<!-- c -->`) and a line that starts with an inline-only tag (`a b i u em strong code kbd mark span sub sup small del ins`) with text after it. Neither can open an HTML block. `plainHTML` in `internal/body/scanner.go`.
- Still refused: block-level tags, a tag alone on its line (type 7 runs to the next blank line), a comment that spans lines. `TestRefusals` has a row for each.
- Evidence: `TestPlainHTML`, the generator (fragments added, 30000 documents), `FuzzAgainstOracle` 45s and `FuzzTolerantAgainstOracle` 30s with no finding; `go test -race ./...` passes.

Decision: none.

## Log

- 2026-11-01: filed.
- 2026-11-01: filed from the refusal census (ADR 261101-tr8e, principle 4).
- 2026-11-01: 261101-4j86 closed; blocked_by removed.
- 2026-11-01 01:56-07:00: claimed (claude).
