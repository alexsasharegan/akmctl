---
id: "261102-cg6g"
title: "Commentary 001: Adaptive Retention and Attention Management (Retention and Compaction)"
summary: "Explores adaptive retention, attention management, exposure efficiency, and knowledge placement in AKM."
type: design/commentary
status: draft
created: 2026-11-02
updated: 2026-11-02
parent: "261101-tpj2"
---

# Retention and Compaction — Commentary 001

The output of architectural brainstorming Session 001 against [Design: Retention and Compaction](../retention-and-compaction.md).

## Session Guidance

This commentary establishes working context for subsequent architectural exploration, not immutable architectural requirements.

**Guidance for Session 002 and subsequent sessions:**

1. Treat the concepts, terminology, and principles established here as accepted working context.
2. Preserve the distinction between accepted working concepts, proposed architectural directions, and unresolved questions.
3. Challenge existing assumptions when new ideas expose weaknesses or contradictions. Do not treat this commentary as an authority that overrides sound architectural reasoning.
4. Continue from the conceptual foundations established here rather than reopening settled terminology without cause.
5. Prioritize architectural exploration and conceptual clarity over implementation details. Do not prematurely translate exploratory ideas into implementation requirements.
6. Maintain attention efficiency throughout the design process: explore deliberately, preserve important reasoning, and avoid unnecessary repetition.
7. Record subsequent sessions as chronological, curated design commentaries rather than silently revising earlier explorations.

**Next session: Session 002 — Adaptive Retention Model**

Focus on how AKM should determine what knowledge to retain, where to place it, and which signals should influence those decisions.

Build upon the foundations of attention management, exposure efficiency, and memory consolidation established in this commentary. The objective is to develop the adaptive retention model, not yet finalize its implementation.

## Adaptive Retention and Attention Management

### 1. Context

This commentary records a subsequent architectural exploration of the ideas introduced in `docs/design/retention-and-compaction.md`.

The original design exploration established the foundations for retention, compaction, archival storage, and project lifecycle integration.

This session explored a broader motivation for those mechanisms: managing the attention demanded of knowledge consumers.

Rather than revising the original design document in place, this commentary preserves the chronological evolution of the architectural thinking. It introduces concepts and working principles for subsequent exploration without claiming to supersede the original document.

### 2. Foundational Motivation: Attention Management

AKM serves two fundamentally different but complementary consumers: humans and agents.

Both have limited resources available for interacting with knowledge.

**Human attention costs:**

- Cognitive load.
- Navigation complexity.
- Reading and interpretation effort.
- Difficulty maintaining a mental model of the knowledge system.

**Agent attention costs:**

- Context-window capacity.
- Token consumption.
- Tool calls and retrieval iterations.
- Reasoning effort.
- The cost of reconciling irrelevant or historical context.

The knowledge system should minimize unnecessary exposure to information while making relevant knowledge readily discoverable.

This suggests a broader architectural framing:

**AKM is a knowledge system that manages the attention of its consumers.**

The filesystem is one interface to knowledge. Git is one persistence mechanism. Compaction is one operation. The overarching objective is to expose the right knowledge at the right level of detail for the current context.

An important consequence follows:

> The number of records is not the measure of complexity experienced by a consumer.

A large repository can be highly efficient if it consistently surfaces the small subset of knowledge needed. A much smaller repository can be inefficient if consumers must inspect substantial amounts of irrelevant information to find what they need.

Optimize for attention, not cardinality.

### 3. Exposure Efficiency

**Working definition:**

Exposure efficiency is the ability of a knowledge system to deliver the information relevant to a consumer's current context while minimizing the attention required to discover, interpret, and use it.

Exposure efficiency has three independent evaluation dimensions.

#### 3.1 Relevance

Are the surfaced records actually useful to the consumer's current context?

Low relevance means the system is exposing information that does not meaningfully contribute to the task.

#### 3.2 Coverage

Did the system surface all the sufficiently relevant records necessary to accomplish the task?

A retrieval operation can return highly relevant results while still omitting important knowledge.

Coverage is therefore distinct from relevance and deserves independent consideration.

#### 3.3 Attention Cost

How much effort did the consumer expend discovering, interpreting, and using the surfaced knowledge?

For agents, this includes tokens, tool calls, retrieval iterations, and reasoning effort.

For humans, this includes navigation, reading, cognitive load, and mental-model maintenance.

These dimensions should not prematurely be collapsed into a single numerical score.

#### 3.4 Cardinality Paradox

Repository cardinality and exposure efficiency are orthogonal.

Consider two hypothetical repositories:

|                 | Repository A | Repository B |
| --------------- | -----------: | -----------: |
| Total records   |      100,000 |        1,000 |
| Records exposed |            3 |          200 |
| Outcome         |    Efficient |  Inefficient |

Repository A can provide excellent exposure efficiency despite its size.

Repository B can impose substantial attention costs despite being much smaller.

The number of records alone tells us very little about the experience of consuming knowledge.

The effectiveness of retrieval must also account for relevance and coverage. Returning three relevant records when five are necessary is not successful exposure.

### 4. Retention, Placement, and Exposure

Three distinct concepts emerged as a useful foundation for the architecture.

| Concept   | Governing question                                          |
| --------- | ----------------------------------------------------------- |
| Retention | What knowledge must we preserve?                            |
| Placement | Where should that knowledge physically reside?              |
| Exposure  | What knowledge should be presented to a consumer right now? |

These responsibilities are related but independent.

A historical architectural decision might have:

- High historical importance and retention value.
- Low current relevance and placement value.
- High exposure value when an agent works on the affected subsystem.

The knowledge does not need to be duplicated or rewritten merely because its relevance to current work changes.

This suggests the following working principle:

> Historical importance informs retention. Current relevance informs placement. Task relevance determines exposure.

### 5. Cognitive Abstractions

Human cognition provides useful functional abstractions for thinking about knowledge management.

These are architectural inspirations, not claims that AKM should reproduce biological mechanisms literally.

| Cognitive abstraction   | Potential AKM equivalent                                |
| ----------------------- | ------------------------------------------------------- |
| Working memory          | Active working set                                      |
| Long-term memory        | Retained knowledge                                      |
| Memory consolidation    | Synthesis and organization of accumulated knowledge     |
| Associative recall      | Graph-based retrieval                                   |
| Attention               | Contextual relevance and exposure                       |
| Rehearsal or activation | Repeated use and promotion                              |
| Forgetting              | Reduced exposure without necessarily deleting knowledge |

**Memory consolidation** is particularly important as an upstream concept.

It describes transforming accumulated experience into organized, durable knowledge.

This is distinct from compaction, which changes the physical representation of records.

For example, closing a project might trigger knowledge consolidation, producing durable knowledge artifacts. Adaptive retention then evaluates the original records, and compaction changes their physical placement.

### 6. Physical and Logical Representations

A useful architectural distinction emerged between the physical and logical representations of knowledge.

**Physical representation**

- Where records actually reside.
- How they are persisted.
- Which storage tier contains them.
- How they are recovered.

**Logical representation**

- How records are identified.
- How relationships are traversed.
- How knowledge is discovered.
- How records are resolved regardless of physical location.

`akmctl` acts as the engine and facade for the logical representation.

Consumers should not need to know whether a record resides in the active working tree, Git history, or a deeper archival substrate.

They should be able to request knowledge by stable identity and let AKM resolve its physical location.

This also introduces the importance of progressive disclosure: consumers should be able to begin with a concise result, inspect provenance, and request additional detail when necessary.

### 7. Adaptive Retention

Adaptive retention is the overarching policy for preserving knowledge and managing its availability.

It extends beyond the original project-closure compaction trigger.

Potential signals include:

- Historical backlinks and architectural significance.
- Current backlinks from active records.
- Active project relationships.
- Recent access patterns.
- Retrieval frequency.
- Knowledge promotion and supersession.

A useful distinction emerged:

**Historical importance informs retention. Current relevance informs storage placement.**

A completed record may remain historically important while having little current relevance. Conversely, an older record may become central to a new project.

Project closure is therefore one opportunity to reevaluate placement, not the only one.

The initial implementation need not depend on sophisticated access telemetry. Deterministic relationship-based signals can provide an initial foundation.

### 8. Working Terminology

The following vocabulary emerged during the exploration.

| Term                 | Working definition                                                                               |
| -------------------- | ------------------------------------------------------------------------------------------------ |
| Attention management | Managing the attention demanded of knowledge consumers                                           |
| Exposure efficiency  | Delivering relevant knowledge with minimal attention required to discover, interpret, and use it |
| Retention            | Determining what knowledge must be preserved                                                     |
| Placement            | Determining where knowledge resides                                                              |
| Exposure             | Determining what knowledge is presented in a particular context                                  |
| Working set          | Knowledge maintained for efficient direct access                                                 |
| Memory consolidation | Transforming accumulated experience into organized, durable knowledge                            |
| Adaptive retention   | Dynamically managing knowledge preservation and storage placement                                |
| Promotion            | Moving knowledge into a more accessible storage tier                                             |
| Compaction           | Removing eligible records from the working set while preserving recoverability                   |

These definitions are working terminology, not necessarily final specification language.

### 9. Architectural Principles Emerging from This Session

The following principles are accepted as working foundations for subsequent exploration:

1. **Optimize for attention, not cardinality.** Repository size is not a reliable measure of consumer complexity or retrieval efficiency.
2. **Separate retention, placement, and exposure.** They answer different architectural questions and should not be conflated.
3. **Treat relevance, coverage, and attention cost as independent dimensions of exposure efficiency.**
4. **Preserve knowledge independently of its current prominence.** Reduced exposure does not imply deletion.
5. **Separate physical and logical representations.** Consumers should not need to understand storage mechanics to retrieve knowledge.
6. **Use cognitive abstractions as architectural inspiration.** Memory consolidation, associative recall, working memory, and attention offer useful models.
7. **Serve humans and agents through the same underlying knowledge architecture.** Optimize for their complementary needs rather than privileging one consumer.

### 10. Open Questions for Subsequent Exploration

This commentary deliberately does not resolve the following questions.

- How should adaptive retention determine placement?
- Which relevance signals should influence promotion and demotion?
- How should historical importance be calculated or classified?
- How should active backlinks be weighted relative to historical backlinks?
- What role should access telemetry play?
- How should memory consolidation integrate with project closure and exit refinement?
- Should compaction operate on individual records, project subtrees, or both?
- What should the logical retrieval interface expose through `akmctl`?
- How should exposure efficiency be measured in practice?
- What constitutes sufficient relevance and coverage for a successful retrieval?
- How should storage tiers evolve beyond the current working-tree/Git-history model?

### 11. Relationship to the Original Design

This commentary extends the conceptual foundations of `docs/design/retention-and-compaction.md`.

It does not replace that document or establish a finalized architecture.

The original design remains the source of its existing architectural exploration. This commentary records subsequent thinking that may eventually justify revisions to the original model.

Future commentaries should preserve chronological continuity and distinguish newly introduced proposals from established architectural decisions.
