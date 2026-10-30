---
id: 261028-0snt
type: task
status: done
origin: 261028-1z36
created: 2026-10-28
updated: 2026-10-28
closed: 2026-10-29
authored_by: agent
title: Choose a display-width library for the ADR index table
---

# Choose a display-width library for the ADR index table

## Goal

The generated ADR index pads its table cells to display width, as Prettier does (East Asian wide characters and emoji count 2 columns). Go's standard library has no width function; the maintainer allowed one small dependency for it (Y16 in [agent-ready-repo-decision-log](../../design/agent-ready-repo-decision-log.md)). Choose it with the same criteria as the YAML evaluation ([go-yaml-library-evaluation](../../design/go-yaml-library-evaluation.md)), in this order:

1. **Baseline:** output matches Prettier 3.x (the tested 3.9.10 or the then-current pin) for a fixture of cells: ASCII, accented Latin, CJK, fullwidth forms, emoji (single, with variation selectors, ZWJ sequences, flags), combining marks. The test is the index table of [akm-spec-review](../../design/akm-spec-review.md) (appendix) run through `prettier --check`.
2. **Performance:** time per cell on a few thousand rows.
3. **Security and supply chain:** maintenance record, advisories, number of dependencies, behaviour on hostile input (very long strings, invalid UTF-8).
4. **Tiebreaker:** popularity and support.

Candidates to start from, to be verified today: `golang.org/x/text/width`, `github.com/mattn/go-runewidth`, `github.com/rivo/uniseg`; also the no-dependency alternative of an embedded width table, and padding by rune count with the file in `.prettierignore`, as a baseline.

Done when: a short note in the vault compares the candidates with evidence, names a choice, and says how many cells in the fixture each candidate gets different from Prettier.

## Outcome

A recommendation, awaiting Alexander's confirmation: use `github.com/clipperhouse/displaywidth` v0.11.0 behind a one-function wrapper. Against Prettier 3.9.10 it differs on 32 of 206 scored fixture cells (3 of 129 common-text cells); no library reproduces Prettier, whose width rule is not Unicode's. The no-dependency table plus emoji rules matches best (5 of 206) but was fitted to the fixtures; rune count is not good enough (84 of 206). Full comparison: [go-width-library-evaluation](../../design/go-width-library-evaluation.md). Committed with this close.

## Log

- 2026-10-28: filed by Claude Code on the Mac at Alexander's direction, from the spec review interview.
- 2026-10-28: claimed (Claude Code on the Mac, side session; Build mode named in the brief).
- 2026-10-28: ground truth built: Prettier 3.9.10 (18 days old), one-column tables, 208 fixture cells, 206 scored (Prettier rewrites a run of spaces and a tab). Prettier's rule read from its source: code points, controls and U+0300-036F and U+FE00-FE0F are 0, wide 2, emoji regex 2.
- 2026-10-28: results (wrong of 206): displaywidth 32, uniseg 38, x/text/width 40, go-runewidth 42 (71 under ja_JP locale), table+emoji 5, table wide-only 29, rune count 84. Decision: recommend displaywidth because the dependency was the maintainer's call (Y16) and the mismatch is cheap (file is in `.prettierignore`); the fitted table baseline is the fidelity alternative. Added displaywidth as a fifth library (not in the brief) because go-runewidth v0.0.31 now builds on its author's `uax29` module.
- 2026-10-28: note written: [go-width-library-evaluation](../../design/go-width-library-evaluation.md) (design, draft). Performance and hostile-input results are in its appendix; no panics, nothing super-linear.
- 2026-10-28: ADR check: none. The task recommends a dependency but adds none; Y16 in [agent-ready-repo-decision-log](../../design/agent-ready-repo-decision-log.md) records the choice once Alexander confirms.
- 2026-10-28: closed (done). Open tasks that link here ([261027-r3hh-repo-records-go-tool](261027-r3hh-repo-records-go-tool.md), [261028-1z36-akm-spec-open-questions](261028-1z36-akm-spec-open-questions.md)) were not edited, per the brief.
- 2026-10-28: Alexander decided against a dependency: the tool implements Prettier's width rule itself (Y28 in [agent-ready-repo-decision-log](../../design/agent-ready-repo-decision-log.md)). The evaluation stays as the source of the algorithm, the fixtures and the generator.
