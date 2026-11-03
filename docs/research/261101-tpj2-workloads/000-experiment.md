---
id: "261102-k7m2"
title: Adaptive Retention Workloads
summary: Establishes the methodology, isolation boundaries, and orchestration for three independent case studies evaluating AKM's adaptive retention architecture.
type: design
status: draft
created_at: 2026-11-02
updated_at: 2026-11-02
parent: 261101-tpj2
---

# Adaptive Retention Workloads

## 1. Purpose

The preceding architectural exploration established a conceptual foundation for adaptive retention and knowledge management in AKM.

The exploration consists of three documents:

1. Retention and Compaction — The original architectural exploration of work lifecycles, retention, archival, compaction, and recovery.
2. Commentary 001 — Expands the objective to attention management, exposure efficiency, and knowledge placement.
3. Commentary 002 — Explores knowledge relationships, progressive formalization, contextual relevance, and refinement.

These documents established useful working concepts but deliberately explored a broad solution space.

This research phase changes the methodology.

Rather than continuing to elaborate the conceptual architecture, we will evaluate it against realistic usage patterns.

Three representative workloads will serve as independent case studies for observing how information is created, used, maintained, retained, rediscovered, and eventually disposed of or archived.

The objective is to discover the smallest coherent architectural model that adequately serves these workloads.

## 2. Governing Principles

### 2.1 The Simplest Solution Wins

Architectural complexity must earn its existence through demonstrated practical value.

Previously explored concepts, mechanisms, and abstractions are hypotheses, not commitments. A concept should survive because it solves an observed problem, not because considerable effort went into developing it.

We should be willing to discard nearly everything if a simpler solution achieves the desired outcomes.

### 2.2 Workloads, Not Record Types

A workload describes a representative pattern of information creation, modification, access, inactivity, rediscovery, and disposition.

The three workloads intentionally exercise different characteristics:

| Workload | Activity profile                                    | Primary challenge                                                     |
| -------- | --------------------------------------------------- | --------------------------------------------------------------------- |
| Waypoint | Episodic, prolonged dormancy, eventual resurrection | Recovering historical rationale and adapting to changed circumstances |
| Atlas    | Continuous, high-frequency development              | Managing accumulated knowledge and changing relevance                 |
| Homebase | Cyclical, recurring, heterogeneous                  | Preserving useful knowledge across responsibilities and contexts      |

These are independent test environments for one coherent knowledge-management system, not proposals for three different architectures.

### 2.3 Separate Architecture from Implementation

The case studies are primarily architectural and behavioral explorations.

Implementation ideas may arise naturally, but should remain hypotheses until supported by demonstrated requirements.

For example, Ripgrep may provide sufficiently efficient filesystem-wide discovery of YAML front matter without introducing a dedicated indexing subsystem. This is worth investigating, but does not establish an architectural requirement.

## 3. Research Organization

### 3.1 Document Location

All workload research documents belong in:

`docs/research/261101-tpj2-workloads/`

The experiment charter and workload documents are maintained together.

Each workload has an input document describing its experimental environment. Its investigation produces a separate findings document.

The filename of an input document does not determine its role. In particular, `001-waypoint-findings.md` is the Waypoint workload definition, not an existing investigation report.

Investigation findings must be created as new sibling documents, using the next available numeric prefix and a descriptive workload-specific filename.

Do not assume paths from previous directory conventions or invent filenames for documents that have not been identified.

### 3.2 Foundational Documents

Every investigation uses the same architectural foundation:

- `docs/design/retention-and-compaction.md`
- `docs/design/commentary/001-adaptive-retention.md`
- `docs/design/commentary/002-adaptive-retention.md`

It also reads:

- This experiment charter.
- Its assigned workload definition.

The exact workload input must be supplied by the session's launch prompt.

### 3.3 Investigation Isolation

Each workload receives a dedicated investigation session.

An individual session must not read:

- Other workload definitions.
- Findings from other workload investigations.
- Cross-workload synthesis documents.

This is deliberate. Each investigation must develop its conclusions independently, without importing discoveries, assumptions, or proposed solutions from another investigation.

The shared architectural foundation is intentional; the isolation of experimental inputs is equally intentional.

## 4. Investigation Protocol

### 4.1 Discussion Before Documentation

An investigation is a collaborative exploration, not a report-generation exercise.

The agent must not immediately produce a findings document, comprehensive analysis, or predetermined architectural recommendation at launch.

Instead, it should:

1. Establish a shared understanding of the workload.
2. Introduce realistic scenarios and events.
3. Examine how information emerges, changes, and accumulates.
4. Challenge assumptions and explore alternative interpretations.
5. Ask questions and pursue promising lines of inquiry.
6. Identify implications for system behavior and architectural responsibilities.
7. Revisit earlier conclusions when new scenarios challenge them.

The human and agent should develop the investigation together through discussion.

The agent should lead the exploration where useful, but must leave room for the human to introduce scenarios, challenge reasoning, and redirect attention.

Do not turn the investigation into a sequence of questionnaires. Ask focused questions when an answer would materially influence the direction of exploration.

### 4.2 Scenario-Driven Exploration

Follow the workload through realistic events rather than beginning with an inventory of record types or predetermined lifecycle rules.

For each significant event, consider:

1. Event: What happened in the environment?
2. Information: What records or knowledge naturally emerge or change?
3. System behavior: What should AKM do, if anything?
4. Temporal evolution: What happens as information ages, becomes inactive, or is superseded?
5. Rediscovery: What should happen when information becomes relevant again?
6. Evaluation: What is the simplest mechanism that achieves the desired outcome?

Do not force every event to exercise every architectural concept.

### 4.3 Evaluating Architectural Concepts

Classify findings into three categories:

- **Must happen:** Necessary to preserve integrity or provide an acceptable experience.
- **Would be useful:** Provides meaningful benefits but is not essential.
- **Unnecessary machinery:** Complexity that does not provide sufficient practical value.

Distinguish observed requirements from proposed implementation approaches.

A workload might demonstrate that old records must be searchable without establishing that AKM needs a relevance graph, dedicated index, or automated retention engine.

### 4.4 Scope Discipline

Avoid prematurely defining:

- Comprehensive record schemas.
- Retention scoring formulas.
- Universal lifecycle state machines.
- Mandatory relationship ontologies.
- Automated compaction policies.
- Dedicated indexing infrastructure.

These may become justified, but should not be assumed.

### 4.5 Reaching a Conclusion

An investigation should continue as long as meaningful questions remain unexplored.

Do not declare completion merely because the initial scenarios have been discussed or a plausible architecture has emerged.

A meaningful stopping point occurs when:

- The important workload scenarios have been examined.
- Significant assumptions have been challenged.
- The major behavioral requirements are understood.
- Architectural implications have been explored.
- Remaining uncertainties can be recorded explicitly rather than requiring further immediate discussion.

At that point, the agent should consolidate the investigation into a standalone findings document.

## 5. Findings Documents

Each investigation produces an independently understandable document containing:

1. **Scenario summary:** The workload and events explored.
2. **Observed requirements:** What the system actually needs to accomplish.
3. **Information lifecycle:** How records and knowledge evolve.
4. **Important discoveries:** Behaviors or requirements that were not obvious beforehand.
5. **Architectural concepts validated:** Existing concepts that proved useful.
6. **Architectural concepts challenged:** Concepts requiring modification.
7. **Unnecessary machinery:** Concepts that did not justify their complexity.
8. **Simpler alternatives:** Existing tools, conventions, or mechanisms that could accomplish the same outcomes.
9. **Open questions:** Issues requiring further exploration.
10. **Recommendations:** The smallest set of capabilities that adequately supports the workload.

Findings should distinguish:

- Demonstrated conclusions.
- Reasonable inferences.
- Unvalidated hypotheses.
- Unresolved questions.

Preserve important reasoning, not merely conclusions. A future reader should understand why a requirement was established and which scenarios support it.

Do not attempt to produce a comprehensive architecture specification within an individual case study.

## 6. Investigation Launch Prompt

Use the following template to start an independent workload investigation. Supply the exact workload name and document path.

> Investigate the **[WORKLOAD]** workload.
>
> Read these documents:
>
> - `docs/design/retention-and-compaction.md`
> - `docs/design/commentary/001-adaptive-retention.md`
> - `docs/design/commentary/002-adaptive-retention.md`
> - `docs/research/261101-tpj2-workloads/000-experiment.md`
> - `[EXACT WORKLOAD DOCUMENT PATH]`
>
> Do not read the other workload definitions or their investigation results.
>
> Use this as the starting point for a collaborative investigation. Explore the workload through discussion, examining scenarios, challenging assumptions, and identifying implications.
>
> Do not jump directly to a report or findings document. Ask questions and pursue promising lines of inquiry. Let the investigation develop through conversation.
>
> Treat previously established architectural concepts as hypotheses to evaluate, not requirements to preserve.
>
> Favor practical behavioral reasoning over premature schema, lifecycle, or algorithm design.
>
> When the investigation reaches a meaningful conclusion, record the findings as a standalone document in the workload research directory, following the experiment charter. Use the next available numeric filename prefix.

## 7. Cross-Workload Synthesis

Synthesis occurs only after all three independent investigations have produced their findings documents.

It takes place in a separate conversation, with all three findings documents deliberately introduced as inputs.

The synthesis should:

1. Compare common requirements across the workloads.
2. Identify requirements unique to particular workloads.
3. Reevaluate the concepts established in Commentaries 001 and 002.
4. Identify concepts that can be discarded.
5. Identify mechanisms that can be simplified or combined.
6. Determine which responsibilities belong to AKM and which can remain ordinary filesystem or Git operations.
7. Propose the smallest coherent architectural model that satisfies the demonstrated requirements.

A useful comparison should distinguish genuine shared requirements from superficially similar behaviors that have different underlying causes.

The synthesis may require substantial orchestration and should not be assumed to fit into any individual case-study session.

Its output should inform subsequent architectural commentary or revisions to the original design document.

## 8. Expected Outcome

This research phase should produce:

- Three independent workload definitions.
- Three independently authored findings documents.
- A synthesis of demonstrated requirements across all workloads.
- A critically reduced set of architectural concepts.
- A clearer understanding of the minimum mechanisms AKM actually needs.

The case studies are evidence-gathering exercises, not a commitment to preserve the existing conceptual architecture.

**Guiding principle: The simplest solution wins.**

The purpose of this phase is to discover what AKM actually needs to do, not to justify everything we previously imagined it might do.
