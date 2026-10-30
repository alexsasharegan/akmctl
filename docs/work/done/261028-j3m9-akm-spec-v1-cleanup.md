---
id: 261028-j3m9
type: task
status: done
origin: 261028-1z36
created: 2026-10-28
updated: 2026-10-28
closed: 2026-10-29
authored_by: agent
title: Prepare a clean v1.0 spec
---

# Prepare a clean v1.0 spec

## Goal

The spec, the guide and the AKM ADR carry the path to v1.0 (Y27 in [agent-ready-repo-decision-log](../../design/agent-ready-repo-decision-log.md)). For the release, rewrite them as clean normative documents that state the result, not how it was reached.

Strip from [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md) and [akm-guide](../../spec/akm-guide.md): status banners and "revised", "accepted", "reversed" notes; decision-id citations (D, O, X, Y) unless kept as stable links into the log; assumption tables; the pointer section; mentions of the vault, `pkm-workbench` and competitor names; "proposed", "draft" and "working title" flags that no longer apply; the first-person history. Keep: the normative rules, the grammar and fixtures, the contract, the examples.

For the decisions: [261027-ensf-akm-an-agent-knowledge-management](../../adr/261027-ensf-akm-an-agent-knowledge-management.md) (with its errata) and [261028-zh0w-obsidian-compatibility-is-a-stated-goal](../../adr/261028-zh0w-obsidian-compatibility-is-a-stated-goal.md) stay as the record. Decide whether a single consolidated ADR for AKM v1.0 supersedes them, or whether they stay as they are with the clean spec beside them. The log, the evaluation notes and git hold the history.

When: after the open decisions settle and before the template repository is created or anything is published (the same trigger as the name, Y22). Also at that point: start a `CHANGELOG` at 1.0.0, set the spec's `VERSION`, and make the name final or change it.

Done when: the spec and guide read as a first edition with no trace of the decisions behind them; every rule id and fixture still resolves; the decision log links to the clean spec; the ADR question is settled.

## Decisions (Alexander, 2026-10-28)

- The name is final: the standard is AKM (agent knowledge management), the tool is `akmctl`.
- The clean text replaces the working text in place: [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md), [akm-guide](../../spec/akm-guide.md) and [akm-api-spec](../../spec/akm-api-spec.md) keep their note names, so links stay valid; git holds the earlier text and the decision log holds the reasons.
- [261027-ensf-akm-an-agent-knowledge-management](../../adr/261027-ensf-akm-an-agent-knowledge-management.md) (with its errata) and [261028-zh0w-obsidian-compatibility-is-a-stated-goal](../../adr/261028-zh0w-obsidian-compatibility-is-a-stated-goal.md) stay as the record; no consolidated ADR.
- The API spec is in this cleanup. `CHANGELOG` and `VERSION` files belong to the template repository, not the vault: the documents state "version 1.0.0" in their text.

## Outcome

- Rewritten in place as first editions, version 1.0.0: [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md) (572 lines), [akm-guide](../../spec/akm-guide.md) and [akm-api-spec](../../spec/akm-api-spec.md). No decision ids (D, O, X, Y), no banners or "draft" and "working title" flags, no links into the vault's tasks or evaluation notes, no history. Kept: every rule id, requirement, principle, fixture and example (checked by a script: the 34 rules, R1 to R10, H1 to H7 and C1 to C12 all resolve; the 52 YAML-lite fixtures and 25 API cases are untouched).
- Added to the specification: the twelve contract clauses (section 9), which the comparison note used to define. Corrected to match the built tool: the CLI table (`--count-only`, `migrate --mode`, the `ui` command, exit code 2), the migration file format (`migrations/<from>-to-<to>.json`), the API spec's fixture list, and the guide's adoption and migrate sections.
- The decision log links to the clean text and records the cleanup as Y38. [261027-ensf-akm-an-agent-knowledge-management](../../adr/261027-ensf-akm-an-agent-knowledge-management.md) and [261028-zh0w-obsidian-compatibility-is-a-stated-goal](../../adr/261028-zh0w-obsidian-compatibility-is-a-stated-goal.md) are unchanged; no consolidated ADR (Alexander's decision). `CHANGELOG` and `VERSION` are for the template repository: [261028-eqwq-akmctl-dogfood-akm-layout](261028-eqwq-akmctl-dogfood-akm-layout.md) is the next task and builds it.
- ADR check: none (the name is final and the ADR question is settled in the log; nothing else changed a convention or procedure).

## Log

- 2026-10-28: filed by Claude Code on the Mac at Alexander's direction. Not for now: the notes are useful while decisions are still moving.
- 2026-10-28: Alexander's order for the next work: this cleanup first, then [261027-nrwf-vault-agent-ready-migration](261027-nrwf-vault-agent-ready-migration.md). The template repository and [261028-eqwq-akmctl-dogfood-akm-layout](261028-eqwq-akmctl-dogfood-akm-layout.md) come after.
- 2026-10-28: Alexander's order changed: this cleanup, then [261028-eqwq-akmctl-dogfood-akm-layout](261028-eqwq-akmctl-dogfood-akm-layout.md) (the CLI is the first test case of the repo kit, so the kinks get worked out before the vault), then [261027-nrwf-vault-agent-ready-migration](261027-nrwf-vault-agent-ready-migration.md).
- 2026-10-28: claimed (Claude Code on the Mac). Interview: name final, rewrite in place, ADRs stay, API spec included.
- 2026-10-28: done. The spec was rewritten section by section from the working text; the guide and API spec edited in place; leftover ids, links and banners searched for and removed; rule ids and fixtures verified by script.
