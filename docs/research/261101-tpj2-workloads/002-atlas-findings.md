---
id: "261102-8q4c"
title: Atlas — Continuous Development and Knowledge Accumulation
summary: Defines an independent workload for evaluating knowledge accumulation, changing relevance, and attention management in a continuously developed software product.
status: draft
kind: research
created: 2026-11-02
updated: 2026-11-02
parent: "261101-tpj2"
---

# Atlas — Continuous Development and Knowledge Accumulation

## 1. Scenario

Atlas is a fictional software product developed and maintained by a small software team.

It is a continuously developed commercial product with an established user base, ongoing feature development, and regular maintenance requirements.

Development follows an agile software development lifecycle. Work is organized around projects, sprints, releases, and individual work items.

The team works on Atlas daily. New requirements arrive, bugs are reported, investigations produce findings, and architectural decisions accumulate as the product evolves.

The project has a sustained development history spanning multiple releases and approximately one year of representative activity.

Unlike a newly created project, Atlas already has a substantial body of historical knowledge. New work regularly intersects with previous work, existing architecture, and decisions made under earlier requirements.

The system must support the immediate demands of active development without allowing accumulated historical information to overwhelm the team's working context.

## 2. Development Environment

### 2.1 Team and Workflow

Assume:

- A small development team.
- Daily development activity.
- Two-week sprint cycles.
- Multiple concurrent development efforts.
- Regular product releases.
- Ongoing bug fixes and maintenance.
- Shared ownership of architectural knowledge.

The team uses work items to track features, bugs, investigations, and other development activities.

Projects and initiatives organize related work toward larger outcomes.

Architectural decisions and durable technical documentation provide context that may outlive individual development efforts.

The investigation should not assume that every meaningful activity produces a formal record or that every record needs a specialized type.

### 2.2 Development History

The product progresses through several broad phases:

| Phase                   | Activities                                                          |
| ----------------------- | ------------------------------------------------------------------- |
| Initial development     | Authentication, accounts, billing, and core functionality           |
| Feature development     | New capabilities driven by product requirements                     |
| Maintenance             | Bug reports, recurring defects, and root-cause investigations       |
| Performance             | Profiling, bottleneck analysis, and optimization                    |
| Architectural evolution | Subsystem redesigns, superseded decisions, and changing constraints |
| Major release           | V1 ships and development transitions toward V2                      |

These phases overlap. They should be understood as a representative history rather than rigid sequential stages.

## 3. Information Produced

Atlas naturally generates several populations of information:

- Feature requests and product requirements.
- Projects, sprint plans, and work items.
- Bug reports and troubleshooting records.
- Technical investigations and experimental findings.
- Architecture and design decisions.
- API and subsystem documentation.
- Performance measurements and optimization findings.
- Testing and release records.
- Deferred ideas and unresolved questions.
- Historical implementation context.

Some information is short-lived and closely tied to an active sprint. Other information remains valuable throughout the life of the product.

The investigation should establish how these populations behave over time without assuming that they require separate retention mechanisms.

## 4. Development Scenarios

### 4.1 A Recurring Bug

A new bug report resembles a problem investigated and resolved eight months earlier.

The original issue was fixed, the associated work was completed, and the relevant development activity is no longer prominent.

The new report raises several questions:

- Can the previous investigation be discovered?
- Can the original symptoms and root cause be recovered?
- Is the previous solution still applicable?
- Did the original fix establish a durable technical insight?
- Should the system recognize the relationship between the two incidents?

The goal is to avoid repeating expensive investigative work without unnecessarily keeping every completed bug report prominent.

### 4.2 Performance Investigation Changes Architectural Assumptions

A performance investigation reveals a bottleneck in a major subsystem.

The team discovers that an architectural assumption made during initial development is no longer appropriate at the current scale.

The investigation produces:

- Profiling results.
- Experimental findings.
- An explanation of the bottleneck.
- A proposed architectural change.
- A revised design decision.

Explore how the investigation should contribute to durable architectural knowledge and how the previous assumptions should remain understandable.

### 4.3 A Deferred Feature Becomes Important

A feature request is considered but deferred across several sprints because of competing priorities.

Eventually, product requirements change and the feature becomes important.

The team needs to recover:

- The original request.
- Previous discussions and investigations.
- Reasons for deferral.
- Known constraints and dependencies.
- Any work already completed toward the feature.

Determine how deferred work should remain discoverable without occupying the same attention as currently active work.

### 4.4 Returning to an Untouched Subsystem

A developer is assigned to a subsystem that has not received meaningful attention in six months.

The developer needs to understand:

- The subsystem's purpose.
- Its current architecture.
- Important historical decisions.
- Known limitations.
- Previous investigations.
- Unresolved questions.
- Relevant completed work.

The challenge is to provide sufficient context without requiring the developer to inspect the entire historical record.

### 4.5 An Architectural Decision Is Superseded

An existing architecture decision record is superseded following a redesign.

The original decision is no longer authoritative, but its rationale remains useful.

For example, understanding the old decision may explain:

- Why the previous architecture existed.
- Which constraints influenced it.
- Why an earlier alternative was rejected.
- Which limitations the new design addresses.

Explore how historical decisions should remain available without confusing them with current guidance.

### 4.6 A Project Completes

A development project finishes after several sprints.

It has accumulated:

- Numerous completed work items.
- Several investigations.
- A collection of implementation decisions.
- A few unresolved questions.
- New architectural documentation.
- Some ideas deferred to future work.

Explore what should happen when the project closes.

Which information should remain associated with the historical project? Which discoveries should become part of the product's durable knowledge? Which records can safely become less prominent?

### 4.7 An Abandoned Feature Returns

A feature was previously abandoned because it was too expensive or technically impractical.

A new product requirement or technological improvement makes the feature desirable again.

The team needs to recover the previous work and understand why it was abandoned.

Determine how the system should make historical work useful again without treating old conclusions as current facts.

## 5. Temporal Characteristics

Atlas exhibits several overlapping timescales:

| Timescale | Representative activity                                         |
| --------- | --------------------------------------------------------------- |
| Daily     | Implementation, reviews, bug fixes, and investigation           |
| Sprint    | Feature delivery, project progress, and work-item completion    |
| Release   | Consolidation of changes and product evolution                  |
| Quarterly | Larger initiatives and architectural improvements               |
| Annual    | Accumulated technical history and longer-term product evolution |

The investigation should consider that information can change in practical importance at each timescale.

A work item can be highly relevant during a sprint and relatively unimportant after completion. An architectural decision may remain important across numerous releases. A previous investigation may become relevant again when a similar problem emerges.

## 6. Research Objectives

The investigation should establish what AKM needs to accomplish in a continuously active development environment.

Pay particular attention to:

1. **Knowledge accumulation:** How should information accumulate without overwhelming active working contexts?
2. **Changing relevance:** How should the system respond as work progresses and priorities change?
3. **Completion:** What should happen to records and knowledge when projects and work items finish?
4. **Historical discovery:** How can old investigations, decisions, and implementation history remain useful?
5. **Recurring problems:** How should related incidents and investigations be connected?
6. **Architectural evolution:** How should superseded decisions and previous assumptions remain understandable?
7. **Deferred work:** How should postponed and abandoned work remain discoverable?
8. **Context restoration:** What information should be surfaced when a developer returns to an unfamiliar or inactive subsystem?
9. **Attention efficiency:** How can the system provide sufficient historical context without exposing unnecessary information?
10. **Knowledge consolidation:** Which discoveries should become durable product knowledge, and what original evidence should remain available?

## 7. Boundaries and Assumptions

- Atlas is a continuously developed software product.
- It has a small development team.
- Development occurs daily.
- Work is organized around projects, sprints, releases, and individual work items.
- The product has approximately one year of representative development history.
- Multiple workstreams can be active concurrently.
- The product's architecture and requirements evolve over time.
- Historical work remains available even after completion.

These assumptions establish the experimental environment. They do not prescribe how AKM should organize or manage its information.

## 8. Central Question

**How can a knowledge-management system support continuous development while allowing a growing body of historical work, decisions, and discoveries to remain useful without overwhelming the team's current attention?**

The investigation should identify the simplest mechanisms that adequately support this outcome.
