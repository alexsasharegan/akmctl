---
id: "261102-k7m2"
title: Adaptive Retention Workloads
summary: Establishes the methodology, isolation boundaries, and orchestration for three independent case studies evaluating AKM's adaptive retention architecture.
type: design
status: draft
created_at: 2026-11-02
updated_at: 2026-11-02
parent: "261101-tpj2"
---

# Adaptive Retention Workloads

## 1. Purpose

The preceding architectural exploration established a conceptual foundation for adaptive retention and knowledge management in AKM.

The exploration consists of three chronological documents:

1. **[Retention and Compaction](https://github.com/alexsasharegan/akmctl/blob/main/docs/design/retention-and-compaction.md)** — The original architectural exploration of work lifecycles, retention, archival, compaction, and recovery.
2. **[Commentary 001: Adaptive Retention and Attention Management](https://github.com/alexsasharegan/akmctl/blob/main/docs/design/commentary/001-adaptive-retention.md)** — Expands the objective to attention management, exposure efficiency, and knowledge placement.
3. **[Commentary 002: Adaptive Retention and the Relevance Graph](https://github.com/alexsasharegan/akmctl/blob/main/docs/design/commentary/002-adaptive-retention.md)** — Explores knowledge relationships, progressive formalization, contextual relevance, and refinement.

These documents established useful working concepts but deliberately explored a broad solution space.

This document introduces a change in methodology.

Rather than continuing to elaborate the conceptual architecture, we will evaluate it against realistic usage patterns.

Three representative workloads will serve as independent case studies for observing how information is created, used, maintained, retained, rediscovered, and eventually disposed of or archived.

The objective is to discover the smallest coherent architectural model that adequately serves these workloads.

## 2. Governing Principles

### 2.1 The Simplest Solution Wins

Architectural complexity must earn its existence through demonstrated practical value.

Previously explored concepts, mechanisms, and abstractions are hypotheses, not commitments. A concept should survive because it helps solve an observed problem, not because considerable effort went into developing it.

We should be willing to discard nearly everything if a simpler solution achieves the desired outcomes.

This applies equally to domain modeling, retention policies, relationship semantics, indexing, automation, and implementation.

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

For example, Ripgrep may provide sufficiently efficient filesystem-wide discovery of YAML front matter without introducing a dedicated indexing subsystem.

This is worth investigating, but does not establish an architectural requirement or implementation commitment.

## 3. Case Study Organization

Each workload has its own definition document and dedicated investigation session.

Workload definitions are independent experimental inputs. They must not be bundled together as shared context for individual investigations.

### 3.1 Workload Definitions

| Workload | Definition              | Findings                         |
| -------- | ----------------------- | -------------------------------- |
| Waypoint | `workloads/waypoint.md` | `workloads/waypoint-findings.md` |
| Atlas    | `workloads/atlas.md`    | `workloads/atlas-findings.md`    |
| Homebase | `workloads/homebase.md` | `workloads/homebase-findings.md` |

Each definition describes its own fictional environment, history, constraints, scenarios, and research objectives.

### 3.2 Investigation Isolation

An individual case-study session should have access to:

- The original design document.
- Commentary 001.
- Commentary 002.
- This research charter.
- Its assigned workload definition.

It should not read the other workload definitions or findings documents.

This is deliberate. The purpose is to evaluate each workload independently, without importing discoveries, assumptions, or proposed solutions from another investigation.

Each session should be able to develop its own conclusions from the same foundational architectural context.

### 3.3 Investigation Launch Prompts

The following prompts establish the launch context for each independent session.

#### Waypoint

> Conduct the Waypoint case study for AKM's adaptive retention architecture.
>
> Read the foundational design documents and research charter referenced in `docs/design/workloads.md`, together with `docs/design/workloads/waypoint.md`.
>
> Treat Waypoint as the sole experimental workload. Do not read the Atlas or Homebase workload definitions or any other case-study findings.
>
> Explore the workload through realistic events and information flows. Evaluate what the system needs to accomplish, how information should evolve over time, and what should happen when the project is resurrected after prolonged dormancy.
>
> Treat previously established architectural concepts as hypotheses to evaluate, not requirements to preserve.
>
> Favor practical behavioral reasoning over premature schema, lifecycle, or algorithm design.
>
> The governing principle is: **the simplest solution wins.**
>
> Produce `docs/design/workloads/waypoint-findings.md` as an independently understandable record of the investigation, its discoveries, validated requirements, challenged assumptions, unnecessary machinery, and recommendations.

#### Atlas

> Conduct the Atlas case study for AKM's adaptive retention architecture.
>
> Read the foundational design documents and research charter referenced in `docs/design/workloads.md`, together with `docs/design/workloads/atlas.md`.
>
> Treat Atlas as the sole experimental workload. Do not read the Waypoint or Homebase workload definitions or any other case-study findings.
>
> Explore the workload through realistic events and information flows. Evaluate how knowledge accumulates during continuous development, how relevance changes, and how the system should handle completed work, recurring investigations, architectural evolution, and historical context.
>
> Treat previously established architectural concepts as hypotheses to evaluate, not requirements to preserve.
>
> Favor practical behavioral reasoning over premature schema, lifecycle, or algorithm design.
>
> The governing principle is: **the simplest solution wins.**
>
> Produce `docs/design/workloads/atlas-findings.md` as an independently understandable record of the investigation, its discoveries, validated requirements, challenged assumptions, unnecessary machinery, and recommendations.

#### Homebase

> Conduct the Homebase case study for AKM's adaptive retention architecture.
>
> Read the foundational design documents and research charter referenced in `docs/design/workloads.md`, together with `docs/design/workloads/homebase.md`.
>
> Treat Homebase as the sole experimental workload. Do not read the Waypoint or Atlas workload definitions or any other case-study findings.
>
> Explore the workload through realistic events and information flows. Evaluate persistent responsibilities, finite projects, recurring activities, enduring resources, archival, and contextual rediscovery across seasonal and annual cycles.
>
> Treat PARA as an organizational lens for exploration, not a predetermined AKM architectural requirement.
>
> Treat previously established architectural concepts as hypotheses to evaluate, not requirements to preserve.
>
> Favor practical behavioral reasoning over premature schema, lifecycle, or algorithm design.
>
> The governing principle is: **the simplest solution wins.**
>
> Produce `docs/design/workloads/homebase-findings.md` as an independently understandable record of the investigation, its discoveries, validated requirements, challenged assumptions, unnecessary machinery, and recommendations.

## 4. Investigation Methodology

### 4.1 Scenario-Driven Exploration

Follow each workload through realistic events rather than beginning with an inventory of record types or predetermined lifecycle rules.

For each significant event, consider:

1. **Event:** What happened in the environment?
2. **Information:** What records or knowledge naturally emerge or change?
3. **System behavior:** What should AKM do, if anything?
4. **Temporal evolution:** What happens as information ages, becomes inactive, or is superseded?
5. **Rediscovery:** What should happen when information becomes relevant again?
6. **Evaluation:** What is the simplest mechanism that achieves the desired outcome?

Do not force every event to exercise every architectural concept.

### 4.2 Evaluating Architectural Concepts

Classify findings into three categories:

- **Must happen:** Necessary to preserve integrity or provide an acceptable experience.
- **Would be useful:** Provides meaningful benefits but is not essential.
- **Unnecessary machinery:** Complexity that does not provide sufficient practical value.

Distinguish observed requirements from proposed implementation approaches.

A workload might demonstrate that old records must be searchable without establishing that AKM needs a relevance graph, dedicated index, or automated retention engine.

### 4.3 Scope Discipline

Avoid prematurely defining:

- Comprehensive record schemas.
- Retention scoring formulas.
- Universal lifecycle state machines.
- Mandatory relationship ontologies.
- Automated compaction policies.
- Dedicated indexing infrastructure.

These may become justified, but should not be assumed.

### 4.4 Findings Documents

Each investigation should produce a findings document containing:

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

Findings should distinguish demonstrated conclusions from hypotheses and unresolved questions.

Do not attempt to produce a comprehensive architecture specification within an individual case study.

## 5. Cross-Workload Synthesis

Synthesis occurs only after the three independent investigations have produced their findings documents.

It should take place in a separate conversation, with all three findings documents deliberately introduced as inputs.

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

## 6. Expected Outcome

This research phase should produce:

- Three independent workload definitions.
- Three independently authored findings documents.
- A synthesis of demonstrated requirements across all workloads.
- A critically reduced set of architectural concepts.
- A clearer understanding of the minimum mechanisms AKM actually needs.

The case studies are evidence-gathering exercises, not a commitment to preserve the existing conceptual architecture.

**Guiding principle: The simplest solution wins.**

The purpose of this phase is to discover what AKM actually needs to do, not to justify everything we previously imagined it might do.
