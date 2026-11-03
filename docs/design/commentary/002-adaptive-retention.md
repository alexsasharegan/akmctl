---
id: "261102-8r4m"
title: "Commentary 002: Adaptive Retention and the Relevance Graph"
summary: "Develops the adaptive retention model through knowledge relationships, progressive semantic formalization, refinement, contextual relevance, and exposure."
type: design/commentary
status: draft
created: 2026-11-02
updated: 2026-11-02
parent: "261101-tpj2"
---

## Session Guidance

This commentary records the output of Session 002 and establishes working context for subsequent architectural exploration. It is a chronological continuation of Commentary 001, not a finalized architecture or immutable specification.

Guidance for subsequent sessions:

1. Treat the accepted working concepts and terminology established in Commentaries 001 and 002 as the current conceptual foundation.
2. Preserve the distinction between established working concepts, proposed architectural directions, and unresolved questions.
3. Continue to challenge assumptions where new exploration reveals contradictions or weaknesses.
4. Do not prematurely translate conceptual models into schemas, algorithms, scoring systems, or implementation requirements.
5. Preserve the separation between knowledge relationships, retention policy, physical placement, and contextual exposure.
6. Favor progressive formalization over requiring a comprehensive ontology or elaborate structure upfront.
7. Maintain attention efficiency: explore deliberately, consolidate useful reasoning, and avoid unnecessary repetition.

Next session: Session 003 — Further Architectural Exploration

The adaptive retention exploration has established useful foundations, but several larger subjects remain unresolved. Subsequent sessions should select and investigate those subjects deliberately rather than extending the current exploration indefinitely.

# Retention and Compaction — Commentary 002

The output of architectural brainstorming Session 002 against `261101-tpj2`.

## Adaptive Retention and the Relevance Graph

### 1. Context

This commentary records the second architectural exploration of the ideas introduced in `docs/design/retention-and-compaction.md` and developed further in Commentary 001.

Commentary 001 established attention management as a central motivation for AKM and distinguished retention, placement, and exposure as separate responsibilities.

Session 002 focused on the knowledge relationships that can inform those responsibilities.

The exploration moved toward a broader conceptual model:

- AKM maintains a logical relevance graph that is independent of physical persistence.
- Relationships provide a foundation for discovering and navigating knowledge.
- Relationship semantics can be progressively formalized as recurring patterns emerge.
- Refinement enriches records and their relationships with judgments about provenance, significance, and future utility.
- Relevance is contextual and cannot be reduced to graph distance, lifecycle state, or observed activity.
- Exposure should distinguish primary information from supporting knowledge and deeper references.

These ideas remain exploratory. This commentary records working concepts and architectural directions without prescribing a complete graph model or adaptive retention algorithm.

### 2. Scope and Management Boundaries

An important distinction emerged between records that participate in AKM's knowledge graph and records whose lifecycle AKM governs.

These are not necessarily the same population.

Work items are a primary candidate population for adaptive retention. Their relevance and placement can change substantially as projects progress, complete, and accumulate historical context.

Other records have different characteristics:

- ADRs may have enduring historical value while their relevance to current implementation changes across software generations.
- Architecture documents, specifications, and glossaries may be maintained as authoritative, continuously useful documents rather than accumulating historical versions as ordinary work records.
- External documents may participate in the relevance graph without being managed by AKM.
- AKM's own infrastructure and metadata may be tracked in the repository without necessarily participating in the knowledge graph.

Adaptive retention should therefore not become a universal document-management policy.

A useful distinction is emerging between:

| Concern             | Governing question                                           |
| ------------------- | ------------------------------------------------------------ |
| Tracking            | Is this record within AKM's tracked filesystem scope?        |
| Graph participation | Does this record participate in the logical knowledge graph? |
| Lifecycle authority | Who governs the record's lifecycle?                          |
| Persistence         | Where and how is the record maintained?                      |

These concerns may overlap, but none should be assumed to imply the others.

The precise terminology and representation remain unresolved. In particular, the distinction between ownership, management responsibility, and lifecycle authority requires further exploration.

### 3. The Relevance Graph

AKM is establishing a logical relevance graph: a representation of meaningful connections between records that is independent of their physical persistence.

The graph is not merely an implementation detail of adaptive retention.

It can support several broader capabilities:

- Knowledge discovery and navigation.
- Contextual retrieval and exposure.
- Provenance and historical reconstruction.
- Lifecycle evaluation.
- Knowledge consolidation.
- Relationship validation and governance.
- Adaptive retention and placement.

A record's participation in the graph does not imply that AKM owns, manages, or physically stores it.

For example, an authoritative architecture document may be an important graph node even when its lifecycle is governed outside AKM.

Likewise, moving a record between physical persistence tiers should not change its logical identity or relationships.

This reinforces the distinction established in Commentary 001:

**The physical representation determines where knowledge resides. The logical representation determines how knowledge is identified, related, and discovered.**

The relevance graph should remain independent of the storage mechanisms used to preserve records.

### 4. Relationship Model

#### 4.1 A Generic Relationship Primitive

The initial relationship model should favor a simple, generic connection between records.

At its most basic:

`A ↔ B`

This connection establishes association without necessarily asserting:

- Hierarchy.
- Directionality.
- Priority.
- Temporal applicability.
- Behavioral consequences.
- A particular retrieval policy.

The generic relationship is symmetric. Directionality and additional semantics belong to characterized relationships.

Cycles should be permissible in the generic graph. A global acyclicity constraint would incorrectly impose the rules of particular relationship types on the entire knowledge graph.

#### 4.2 Characterized Relationships

A relationship can acquire more precise semantics when its meaning becomes understood and useful.

Relationship characterization may introduce several dimensions:

| Dimension     | Concern                                                         |
| ------------- | --------------------------------------------------------------- |
| Structural    | How records are organized or contained.                         |
| Directional   | Whether the relationship has a source and target.               |
| Temporal      | Whether the relationship changes over time.                     |
| Behavioral    | Whether the relationship imposes rules or invariants.           |
| Informational | What the relationship communicates about the connected records. |
| Retrieval     | How the relationship can inform discovery or exposure.          |

These are conceptual dimensions, not necessarily mutually exclusive relationship categories.

Examples:

- `parent` has structural and behavioral semantics.
- `supersedes` can be directional, temporal, informational, and relevant to retrieval.
- `related` may remain a minimally characterized association.

A characterized relationship is more than a descriptive label. Its semantics may establish meaningful constraints on the system.

For example, a parent-child relationship can impose lifecycle behavior:

- A parent cannot become `done` until all its children are `done`, recursively.
- A record cannot be blocked by a parent, regardless of how far up the parent hierarchy the relationship extends.

These rules belong to the semantics of the parent relationship, not to generic graph connectivity.

#### 4.3 Progressive Formalization

The relationship vocabulary should emerge from actual use rather than requiring a comprehensive ontology upfront.

A possible progression:

1. A human or agent establishes an ad hoc connection.
2. Repeated use reveals a recurring semantic pattern.
3. Someone recognizes that the connection has meaning beyond generic association.
4. The relationship is characterized and its semantics formalized.
5. AKM can begin validating the relationship and enforcing any established invariants.

This supports an extensible knowledge model without making arbitrary connections appear more authoritative or meaningful than they are.

Existing metadata conventions can provide migration opportunities. For example, recurring values in `document.type` might eventually be recognized as candidates for formalized types understood by AKM.

The checker could support this evolution by:

- Identifying recurring ad hoc values.
- Warning about recognized but unformalized conventions.
- Suggesting migration to characterized types.
- Validating newly established invariants.
- Helping detect inconsistencies in existing records.

The intended progression is:

**Discovery remains open. Characterization is extensible. Recognized semantics are constrained. Enforcement applies only where semantics have been established.**

This is a proposed direction, not a finalized schema or migration mechanism.

### 5. Relationship Provenance and Authority

Relationships can originate from different sources, and their provenance matters.

Three useful distinctions emerged:

| Origin   | Meaning                                                                 |
| -------- | ----------------------------------------------------------------------- |
| Asserted | Explicitly established by a human or authorized process.                |
| Inferred | Proposed by an agent or automation based on observed context.           |
| Derived  | Follows from established relationship semantics or other defined rules. |

These describe the origin of a relationship, not its type.

An inferred relationship should not silently acquire the same authority as an explicitly asserted one.

In particular:

- Inference can suggest connections worth investigating.
- Confirmation can establish a stronger relationship.
- Derived relationships can be trusted to the extent that their underlying rules are dependable.
- Observed activity or repeated exposure should not automatically turn an inferred relationship into an authoritative one.

The authority of a relationship should reflect how it was established and whether its semantics have been recognized.

### 6. Relationship Lifecycle and Contextual Applicability

Relationships can remain historically valid while becoming less applicable to current work.

A relationship need not be deleted simply because its operational relevance changes.

For example, an ADR may be superseded by a later architectural decision. Its former relationship to implementation work remains historically meaningful, even if it should no longer guide current implementation.

Likewise, a completed work item may lose relevance to ongoing execution while retaining value as historical evidence.

This introduces two distinct concerns:

- **Relationship history:** What connection existed, and how did it evolve?
- **Contextual applicability:** Is that connection useful for the purpose currently being served?

Applicability is not necessarily binary.

A superseded ADR may be:

- Inapplicable to current implementation.
- Relevant to understanding historical decisions.
- Useful when planning a migration.
- Important when establishing non-goals.
- Relevant again when revisiting an abandoned approach.

Similarly, a completed record can lose relevance to execution while gaining relevance to historical analysis.

Relevance is therefore purpose-dependent, and relationship history should not be confused with current applicability.

### 7. Refinement as Knowledge Enrichment

An important insight from Session 002 is that relationships should not be treated as passive observations of graph connectivity.

AKM has natural lifecycle opportunities to enrich records and establish meaningful connections.

Two particularly important junctures are:

**Input refinement**

When work enters a more organized state, refinement can establish:

- Initial context and purpose.
- Parent or originating work.
- Relevant relationships and dependencies.
- Scope and intended outcomes.
- Known references and supporting material.

**Exit refinement**

When work concludes, refinement can organize what was learned:

- Decisions and their rationale.
- Architectural discoveries.
- Resulting or newly identified ADRs.
- Knowledge that should be promoted into durable artifacts.
- Supersession and dependency relationships.
- Provenance and supporting evidence.
- Follow-up work and newly discovered connections.

Exit refinement is especially important for connecting transient work to durable knowledge.

For example, a work item may produce an architectural discovery that should become an ADR or inform an authoritative architecture document.

The relationship graph is enriched through these activities rather than being left to emerge solely from backlinks or access patterns.

Refinement should be understood as a recurring knowledge-quality process, not merely metadata cleanup or a one-time lifecycle ceremony.

### 8. Consolidation, Evidence, and Provenance

Session 001 distinguished memory consolidation from physical compaction. Session 002 further explored the consequences of synthesis and refinement.

Synthesis is inherently selective and lossy.

If a synthesis preserves every detail without selection or abstraction, it is effectively concatenation rather than synthesis.

Consolidation requires judgment about:

- What matters.
- What is redundant.
- What can be generalized.
- What should be omitted.
- What must remain traceable to its source.
- Which original evidence remains independently valuable.

An ADR may consolidate knowledge surfaced across one or more work items. It can preserve the resulting decision and rationale without reproducing every detail of the work that led to it.

However, consolidation does not automatically invalidate or replace its source records.

A useful distinction is:

| Operation     | Concern                                                                                               |
| ------------- | ----------------------------------------------------------------------------------------------------- |
| Refinement    | Improving the organization, relationships, and quality of a record.                                   |
| Consolidation | Transforming accumulated experience into a more useful knowledge representation.                      |
| Compaction    | Reducing physical representational overhead or removing eligible records from the active working set. |
| Archival      | Preserving records and their recoverability outside the active working set.                           |

The goal is to preserve important knowledge, not necessarily every original record indefinitely.

However, some original evidence cannot be safely replaced by a summary. A synthesized interpretation may preserve conclusions while omitting details that later become important.

Provenance and source traceability should therefore survive consolidation, even when the original details are not reproduced in the resulting artifact.

The retention decision must distinguish between knowledge captured elsewhere and evidence that remains independently valuable.

### 9. Importance, Relevance, and Demand

Session 002 identified a need to distinguish several concepts that can otherwise become conflated.

| Concept    | Meaning                                             |
| ---------- | --------------------------------------------------- |
| Importance | The enduring significance of knowledge or evidence. |
| Relevance  | Its usefulness to a particular context or purpose.  |
| Demand     | Observed activity, access, or retrieval behavior.   |

These are related signals, not interchangeable measurements.

Examples:

- A routine active task may be highly relevant to current work but have little enduring importance.
- An old ADR may have substantial historical importance but little relevance to current implementation.
- A detailed bug report may remain important as primary evidence even when its RCA is the more relevant starting point for subsequent work.
- A frequently accessed record may be popular without being intrinsically important.

This last point is especially important for adaptive retention.

Access frequency, relationship density, and visibility can become self-reinforcing. A record that receives more exposure may accumulate more activity, which could then cause it to receive even more exposure.

Popularity is evidence of demand, not proof of intrinsic importance.

Adaptive retention should avoid treating observed activity as an unquestioned measure of knowledge value.

### 10. Contextual Relevance and Exposure

#### 10.1 Graph Distance Is Not Relevance

Graph distance is a property of connectivity. Relevance is a property of meaning in context.

A record's relevance cannot be determined solely by how many relationships separate it from the current task.

Graph proximity may be a useful signal, but it is not a reliable substitute for semantic judgment.

Knowledge lineage is not relevance inheritance.

A record several relationships away from active work may remain highly relevant because it contains primary evidence, historical context, or a detail that cannot be recovered from a derived artifact.

Conversely, a directly connected record may have little practical value for the current task.

#### 10.2 Progressive Disclosure

Exposure should help consumers begin with the most useful information and retrieve additional detail when necessary.

Consider a bug investigation:

1. A bug report contains reproduction steps and original observations.
2. An RCA refines the evidence into an explanation of the underlying cause.
3. An ADR records a resulting architectural or corrective decision.
4. Implementation work follows from that decision.

When implementing the fix, the RCA may be the most useful starting point. The ADR provides the decision context, and the original bug report provides detailed evidence when the implementation requires closer investigation.

The original bug report does not necessarily deserve primary exposure merely because it is the source of the RCA.

However, it may deserve continued supporting relevance: it should remain readily discoverable when the RCA or its resulting work is exposed.

This suggests three useful exposure roles:

| Role       | Purpose                                                                                         |
| ---------- | ----------------------------------------------------------------------------------------------- |
| Primary    | The record most directly useful to the current context.                                         |
| Supporting | A record that provides meaningful context or substantiates a primary record.                    |
| Reference  | A record that is not ordinarily needed but should be readily retrievable for deeper inspection. |

These are conceptual exposure roles, not proposed record types or finalized metadata.

#### 10.3 Supporting Relevance and Refinement

The relationship between two records does not inherently determine their exposure roles.

For example, an RCA may be derived from a bug report, but that relationship alone should not guarantee that the bug report always appears alongside the RCA.

Instead, refinement can establish that the original report contains details worth keeping readily accessible.

This separates three concerns:

1. **Relationship:** How are the records connected?
2. **Refinement judgment:** What significance does that connection have for future use?
3. **Exposure behavior:** How should the connected record be made available in relevant contexts?

This suggests a possible concept of a relevance floor: certain records retain a minimum level of supporting discoverability even when ordinary graph proximity or ranking would otherwise place them much lower.

The precise mechanism is unresolved. It may be more appropriate to represent this as a qualitative exposure relationship or refinement judgment than as a numeric graph weight.

The important distinction is that a supporting record need not compete with the primary record for attention in order to remain accessible.

### 11. Architectural Principles Emerging from This Session

The following principles are accepted as working foundations for subsequent exploration:

1. **Separate graph participation from lifecycle authority.** A record can participate in the relevance graph without being managed by AKM.
2. **Separate logical identity from physical persistence.** Moving a record between storage tiers should not alter its logical identity or relationships.
3. **Use a generic relationship primitive.** Basic connectivity should not imply hierarchy, directionality, priority, or behavior.
4. **Allow progressive relationship formalization.** More precise semantics should emerge from recurring patterns rather than requiring a comprehensive ontology upfront.
5. **Constrain recognized semantics.** AKM should enforce only the invariants associated with established relationship types.
6. **Preserve relationship provenance.** Asserted, inferred, and derived connections have different origins and should not silently acquire equivalent authority.
7. **Distinguish relationship history from contextual applicability.** A historically valid relationship may no longer be applicable to current work.
8. **Treat refinement as knowledge enrichment.** Input and exit refinement are important opportunities to improve relationships, provenance, and durable knowledge.
9. **Distinguish consolidation from compaction.** Knowledge synthesis is selective and lossy; physical compaction is a separate operation.
10. **Preserve provenance through consolidation.** Synthesized knowledge should remain traceable to important source records and evidence.
11. **Distinguish importance, relevance, and demand.** Observed activity does not establish enduring knowledge value.
12. **Do not equate graph distance with relevance.** Semantic significance and purpose can outweigh proximity.
13. **Support progressive disclosure.** Primary information should be exposed without losing access to supporting evidence and deeper references.
14. **Treat exposure implications as refinement judgments.** Relationship types should not automatically encode every contextual exposure decision.

### 12. Open Questions and Deferred Exploration

This session established conceptual foundations but deliberately did not resolve several larger architectural subjects.

#### 12.1 Graph and Management Boundaries

- How should graph participation be represented and governed?
- What is the precise distinction between ownership, management responsibility, and lifecycle authority?
- How should AKM-managed infrastructure records be distinguished from records that participate in the knowledge graph?

#### 12.2 Relationship Semantics

- How should characterized relationships be represented?
- How should temporal changes and supersession be recorded?
- How should inferred relationships be reviewed, confirmed, or rejected?
- What should the checker validate, and how should relationship semantics evolve through migration?

#### 12.3 Adaptive Retention and Placement

- How should importance, relevance, demand, and relationship semantics influence retention?
- How should promotion and demotion work?
- How should different record populations receive different retention policies?
- How should the system avoid self-reinforcing exposure and popularity?
- What is the relationship between storage placement and contextual exposure?

#### 12.4 Refinement and Consolidation

- Which refinement activities should be explicit lifecycle responsibilities?
- How should humans and agents share responsibility for enrichment?
- When should consolidation trigger reevaluation of source records?
- What evidence must remain independently preserved?
- How should provenance and recovery work across successive consolidation operations?

#### 12.5 Contextual Retrieval and Exposure

- How should a consumer's purpose influence relevance?
- How should primary, supporting, and reference records be exposed?
- How should supporting relevance or a relevance floor be represented?
- How should retrieval balance relevance, coverage, and attention cost?
- How should the system evaluate exposure efficiency in practice?

These subjects are intentionally deferred. The next architectural exploration should select a coherent subset rather than attempting to resolve all of them simultaneously.

### 13. Relationship to Previous Design Work

This commentary extends the conceptual foundations established in:

- `docs/design/retention-and-compaction.md`
- Commentary 001: Adaptive Retention and Attention Management

It does not replace either document or establish a finalized architecture.

Commentary 001 remains the record of the initial attention-management exploration. Commentary 002 records the subsequent exploration of the relevance graph, relationship semantics, refinement, and contextual exposure.

The original design remains the source of its existing architectural exploration. These commentaries preserve the chronological evolution of the ideas and identify concepts that may eventually justify revisions to the original design.

Future work should build on these foundations while continuing to distinguish exploratory concepts from established architectural decisions.
