---
id: 261028-1z36
type: task
status: done
origin: 261026-m74k
created: 2026-10-28
updated: 2026-10-28
authored_by: agent
adr: candidate
title: Settle the open questions in the AKM spec
---

# Settle the open questions in the AKM spec

## Goal

Small decisions left open when [261026-m74k-agent-ready-repo-layout-spike](261026-m74k-agent-ready-repo-layout-spike.md) closed, listed in "Still open" in [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md):

- **Formatter:** require Prettier for style (a pinned version and shipped config) with the tool validating structure and no `fmt` command? Settle: the generated ADR index must survive Prettier; Prettier's handling of long frontmatter values must be tested against the simple-subset parser; Node is needed only for formatting.
- **`renumber`:** in the first version of the tool, or a manual repair?
- **`summary` on work items:** required, optional (as now) or dropped in favour of the title?
- **UI name:** `akmctl ui` or a separate AKM Workbench (one binary or two).
- **Name:** AKM and `akmctl` are a working title; the overlap with `itlackey/akm` and the AgenticAKM paper is accepted knowingly ([261027-ensf-akm-an-agent-knowledge-management](../../adr/261027-ensf-akm-an-agent-knowledge-management.md)).
- **Concurrency and claiming:** the spec does not say which agent or person holds an `active` item (a `claimed_by` field?).
- **Content safety:** no mechanical check for secrets or private text in records.

Done when: each question has a decision recorded in the spec (and in a new ADR if it changes how repos are set up), or is dropped with a reason.

## Outcome

Every question is decided or dropped with a reason, recorded in [agent-ready-repo-decision-log](../../design/agent-ready-repo-decision-log.md) (Y1 to Y20) and applied to [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md): Prettier is required (pinned, embedded formatting off); optional `claims` and `safety` extensions; `renumber` is not in v1; work-item `summary` stays optional; the UI is `akmctl ui`; the name stays a working title; the frontmatter grammar's edge rules are confirmed; a small dependency is allowed for display width (follow-up: [261028-0snt-go-width-library-evaluation](261028-0snt-go-width-library-evaluation.md)). The review's mechanical defects are fixed. Commits: `c2de16e` (mechanical fixes), `1b29292` (interview decisions), and the commit of this close.

Why: ADR check below.

## Log

- 2026-10-28: filed by Claude Code on the Mac at Alexander's direction, from the close of [261026-m74k-agent-ready-repo-layout-spike](261026-m74k-agent-ready-repo-layout-spike.md).
- 2026-10-28: claimed (Claude Code on the Mac, side session from the AKM spec session). Scope: review the spec, test Prettier against H5, draft two proposals; no decisions made.
- 2026-10-28: Prettier 3.9.10 test done (fixtures and parser in the appendix of [akm-spec-review](../../design/akm-spec-review.md)). Finding: frontmatter survives only with `embeddedLanguageFormatting: off`; defaults wrap flow lists over 80 columns, `proseWrap: always` folds long summaries; the generated ADR index table is re-aligned (options in the note). Formatter question: evidence supports "require Prettier, no `fmt` command", pending Alexander.
- 2026-10-28: review found 29 defects (9 high); no `summary`, `renumber`, UI-name or name decision is forced by them. Proposals drafted for claiming (`claimed_by`, advisory, stale after 24 h) and secrets checks (built-in patterns plus deny list); recommendations only. Set `blocked`: needs Alexander's decisions.
- 2026-10-28: Alexander asked for the mechanical defects from [akm-spec-review](../../design/akm-spec-review.md) to be fixed and for an interview on the rest. Fixed in [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md): stale names (`wb`, `docs-check`, `--open`, "the script"), id clashes (roadmap rules renamed RM1 to RM4, work rules WK1 to WK12 with the gap closed, H6 moved after H5, suffixed decision ids D2b to D6c removed, P1 removed), the example ADR id, the tree (`.layout`), the `AGENTS.md` headings (rule G1), the fallback greps (quote-tolerant, `ready` row), WK7 (Goal and Outcome), new rules WK2, WK11 and WK12, the collision text, `created` as UTC, the person's name in the spec, metadata; the comparison note banner; the guide's "ADR 260912-3q6j"; the API spec task's operation names. Left for decisions (the interview): review defects 1, 2, 3, 5, 7, 8, 9, 15 (skill `VERSION`), 16, 17 (a close/set command), 20, 25 to 27, and the Prettier and claiming proposals.
- 2026-10-28: interviewed Alexander on the rest of the review. All 14 answers took the recommended option (see [agent-ready-repo-decision-log](../../design/agent-ready-repo-decision-log.md), section C). Applied to the spec: H6 reworded and an id recipe; optional `claims` and `safety` extensions; Prettier required with `embeddedLanguageFormatting: off`; rule A4 dropped (index check is now A4); migrations embedded in the binary; settings in `.layout`; Large tier is nested `AGENTS.md` only; a normative frontmatter grammar with fixtures; `set` and `close` commands; spec history moved to the decision log; ADR 261027-ensf errata (`dropped`, acceptance hash, star count removed). Not in v1: `renumber`.
- 2026-10-28: remaining answers: grammar edge rules confirmed (Y15), display width by a small dependency (Y16), Prettier pin latest 3.x at release (Y17), UI is `akmctl ui` (Y18), name stays a working title (Y19), `authored_by` markers removed from six notes (Y20). Filed [261028-0snt-go-width-library-evaluation](261028-0snt-go-width-library-evaluation.md).
- 2026-10-28: ADR check: candidate (questions 1, 2, 3 and 4: Prettier required, extensions, a new dependency, an exception to "standard library only"). Why: these refine [261027-ensf-akm-an-agent-knowledge-management](../../adr/261027-ensf-akm-an-agent-knowledge-management.md) rather than reverse it; the reasons and alternatives are in [agent-ready-repo-decision-log](../../design/agent-ready-repo-decision-log.md) sections C. A sweep may fold them into an ADR on the AKM v1 spec decisions.
- 2026-10-28: closed.
