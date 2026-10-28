---
id: "261028-x1b4"
title: Match Prettier's width on non-standard emoji ZWJ pairs
summary: The width code differs from Prettier for emoji joined by a zero-width joiner when the pair is not a standard sequence; the ADR index could misalign.
status: ready
kind: bug
priority: 1
created: 2026-10-28
updated: 2026-10-28
---

# Match Prettier's width on non-standard emoji ZWJ pairs

## Goal

`internal/width` implements Prettier's string-width rule from Unicode 17.0.0 data. It differs from Prettier on emoji pairs joined by a zero-width joiner that are not standard emoji sequences: the index table could be padded differently from what Prettier produces, and rule A4 would then disagree with Prettier. No record has hit this.

Done when:

- [ ] A test case set of non-standard ZWJ pairs with Prettier's measured widths (run with the pinned Prettier 3.9.10).
- [ ] `internal/width` returns Prettier's width for every case, or the difference is documented in the code with the reason.
- [ ] A cell with such a pair in an ADR title passes `prettier --check` on the generated index.

## Log

- 2026-10-28: filed.
