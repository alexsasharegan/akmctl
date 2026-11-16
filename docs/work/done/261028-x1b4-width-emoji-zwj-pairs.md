---
id: "261028-x1b4"
title: Match Prettier's width on non-standard emoji ZWJ pairs
summary: The width code differs from Prettier for emoji joined by a zero-width joiner when the pair is not a standard sequence; the ADR index could misalign.
status: done
kind: bug
priority: 1
created: 2026-10-28
updated: 2026-11-16
closed: 2026-11-16
---

# Match Prettier's width on non-standard emoji ZWJ pairs

## Goal

`internal/width` implements Prettier's string-width rule from Unicode 17.0.0 data. It differs from Prettier on emoji pairs joined by a zero-width joiner that are not standard emoji sequences: the index table could be padded differently from what Prettier produces, and rule A4 would then disagree with Prettier. No record has hit this.

Done when:

- [x] A test case set of non-standard ZWJ pairs with Prettier's measured widths (run with the pinned Prettier 3.9.10).
- [x] `internal/width` returns Prettier's width for every case, or the difference is documented in the code with the reason.
- [x] A cell with such a pair in an ADR title passes `prettier --check` on the generated index.

## Next

## Outcome

- Test cases: `internal/width/testdata/prettier-3.9.10-zwj.json`, 4465 cases with the widths Prettier 3.9.10 gives (`prettier.util.getStringWidth`), written by `internal/width/gen/zwj-cases.mjs`.
- `internal/width` returns Prettier's width for every one (`TestPrettierZWJ`): a chain is one emoji only when it is a standard sequence of `data/emoji-zwj-sequences.txt`; otherwise each emoji counts 2 and each joiner 1. The one known miss of `TestPrettier` is gone.
- An ADR title with such pairs: the generated index passes `prettier --check` with this build and fails with the build before it (run by hand; the Go tests do not call Prettier).

Decision: none (the choice of source is in the Log).

## Log

- 2026-10-28: filed.
- 2026-11-16 07:30-07:00: claimed (claude-fable-5-2).
- 2026-11-16 07:31-07:00: ticked "A test case set of non-standard ZWJ pairs with Prettier's measured widths (run with the pinned Prettier 3.9.10).": internal/width/testdata/prettier-3.9.10-zwj.json: 4465 cases (2725 non-standard pairs, prefixes, extended chains, bare joiners, and every standard sequence with and without U+FE0F), written by gen/zwj-cases.mjs with Prettier 3.9.10
- 2026-11-16 07:31-07:00: ticked "`internal/width` returns Prettier's width for every case, or the difference is documented in the code with the reason.": TestPrettierZWJ: every case matches; the known-miss list of TestPrettier is empty now
- 2026-11-16 07:31-07:00: ticked "A cell with such a pair in an ADR title passes `prettier --check` on the generated index.": by hand: an ADR title with two non-standard pairs and one standard sequence; the index of the old build fails prettier --check, the index of this build passes, and check reports no A4
- 2026-11-16 07:31-07:00: the maintainer approved the download on 2026-11-16: emoji-zwj-sequences.txt, Emoji 17.0, from unicode.org/Public/17.0.0/emoji/. Source pick: Unicode's file, not a list derived from Prettier's emoji pattern; measured against Prettier 3.9.10 the two agree on all 1614 sequences, and Prettier also accepts each sequence without its variation selectors, so the table stores them without U+FE0F
- 2026-11-16 07:31-07:00: noticed, not changed: gen/main.go writes to standard output, so go generate ./internal/width as AGENTS.md gives it does not rewrite tables.go; the tables were regenerated with a redirect
