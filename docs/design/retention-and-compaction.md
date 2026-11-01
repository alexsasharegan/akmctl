---
type: design
status: draft
created_at: 2026-11-01
updated_at: 2026-11-01
authored_by: human
---

# AKM: Archiving, Compaction, and Knowledge Lifecycle

**Document type:** Architecture discovery spike  
**Status:** Exploratory — architectural direction emerging  
**Primary objective:** Define a coherent lifecycle for work records and accumulated knowledge, including how completed work is evaluated, archived, compacted, and recovered.

---

## 1. Problem Statement

AKM (Agent Knowledge Management) manages work and knowledge for humans and agents through durable, structured Markdown records.

As projects accumulate work items, decisions, investigations, and historical context, the active working set grows. This creates several problems:

- Agents must navigate increasingly large collections of historical records.
- Relevant knowledge becomes harder to distinguish from obsolete or temporary work history.
- Important context can be lost when completed work is removed.
- Reconstructing the significance of old work becomes increasingly expensive as the original context fades.
- There is no sufficiently defined lifecycle for transitioning completed projects and their accumulated work into a durable archive.

AKM needs a deliberate knowledge lifecycle that keeps the active working set useful while preserving important knowledge, provenance, and the ability to recover historical records.

The goal is not simply to delete old records or summarize them. It is to establish a repeatable process for evaluating completed work, promoting durable knowledge, archiving projects, and safely compacting historical records.

## 2. Foundational Concepts and Design Principles

### 2.1 Markdown as durable persistence

Markdown files remain the canonical persistence format.

- Avoid introducing a separate XML, JSON, or database-oriented persistence model for ordinary AKM records.
- Keep YAML frontmatter clean and focused on established record metadata.
- Use controlled Markdown sections for specialized machine-maintained data where necessary.
- AKM owns the generation, parsing, validation, and repair of machine-maintained sections.
- Human- and agent-authored content remains ordinary Markdown.

### 2.2 Human and agent interoperability

AKM is an interface to shared durable state, not a workflow belonging to one particular agent or model.

- Humans and agents should operate on the same domain concepts.
- The CLI, UI, and Markdown files are complementary interfaces.
- Agents should be interchangeable clients of the AKM domain model.
- Important lifecycle operations should be governed by AKM rather than depending on agents to manually preserve invariants.

### 2.3 Progressive structure

AKM must remain inexpensive to use for small projects while supporting substantial projects with more structure.

The lifecycle and archival model should not require every project to have elaborate documentation or ceremony.

### 2.4 Knowledge should outlive the work that produced it

A completed work item is evidence of work performed. It is not necessarily the appropriate long-term home for the knowledge discovered during that work.

Important knowledge should be promoted into appropriate durable artifacts, such as:

- Architecture documents
- Specifications
- Architecture Decision Records (ADRs)
- Durable knowledge or resource notes

The original work record can eventually be archived without requiring it to remain individually present in the active working set.

### 2.5 Capture knowledge while context is fresh

The cost of recognizing important knowledge is lowest while the person or agent still has the context in which it was discovered.

AKM should capture useful knowledge and relationship signals at natural completion points rather than expecting a future agent to reconstruct their significance.

### 2.6 Preserve recoverability

Archiving and compaction must not make historical evidence irretrievable.

Every archived record must remain identifiable and recoverable through a valid chain of provenance and recovery information.

---

## 3. The Four Questions of Compaction

These are distinct architectural questions and should remain distinct throughout the investigation.

### 3.1 What should be compacted?

This is the **retention policy** question.

Determine which records and information should be:

- Retained individually
- Promoted into durable knowledge artifacts
- Summarized or consolidated
- Archived and removed from the active working set
- Preserved only as recoverable historical evidence

Completed work items should generally be presumed eligible for eventual compaction. Long-term retention should be justified by the information they contain or their ongoing relevance, not simply by the fact that they once represented meaningful work.

### 3.2 When should compaction be considered?

This is the **lifecycle boundary** question.

Potential triggers include:

- Project completion
- Major software version transitions
- Milestone completion
- Temporal boundaries, such as weeks, months, quarters, or years
- Volume thresholds
- Project or initiative closure
- Supersession or replacement of prior work

A boundary should trigger evaluation, not indiscriminate automatic deletion.

### 3.3 How should compaction work?

This is the **operation** question.

Investigate:

- How an archive is represented.
- How the project retains its identity and recovery information.
- How child records are indexed before removal.
- How descriptive summaries and operational metadata are produced.
- How records are removed from the active working set.
- How the operation can be safely interrupted, resumed, verified, or rolled back.

### 3.4 How is recoverability guaranteed?

This is the **integrity** question.

Investigate:

- Stable record identity.
- Git recovery points and exact original paths.
- Reachability through nested or successive compactions.
- Provenance of each compaction operation.
- Integrity of the logical recovery graph.
- Behavior when repositories, branches, or historical records change.

These four questions are related but should not be conflated. Retention policy should not be dictated by Git mechanics, and lifecycle boundaries should not automatically determine what information deserves preservation.

---

## 4. Leading Architectural Direction: Projects as Archive Indexes

### 4.1 Working hypothesis

The leading direction is to use the existing project record as the durable archive and recovery index for its child records.

A completed project already provides:

- A stable identity.
- A natural scope for child work.
- A description of the work and its intended outcomes.
- A lifecycle and completion boundary.
- Recursive containment.
- Existing closure invariants.
- A natural location for synthesized project knowledge and recovery metadata.

Rather than introducing a parallel compaction-record hierarchy, investigate extending the project lifecycle so that an archived project retains the information necessary to understand and recover its historical work.

The project becomes more than a container: it becomes the durable record of the completed initiative and the index into its historical evidence.

### 4.2 Reuse existing project invariants

AKM already has important structural constraints:

- A project cannot close while a child record remains open.
- If a child is itself a project, its children must satisfy the same closure requirements.
- Dependency rules prevent invalid relationships, including certain parent-child dependency relationships that could violate acyclicity.

These existing rules may provide much of the foundation needed for safe archival.

The spike should investigate how much of the existing project machinery can be reused for archival and recovery rather than implementing a second hierarchy with duplicate invariants.

### 4.3 Nested archives

Nested projects introduce an important recursive case.

A parent project may contain a child project that has already been archived. The child project may itself retain an index for its own compacted descendants.

The parent must preserve the child's identity and recovery relationship without needing to duplicate the child's entire recovery index.

Investigate whether the existing hierarchy naturally provides this structure and whether nested archives can be recovered independently.

### 4.4 Evaluate the alternative, but do not make it the starting assumption

A separate compaction record is a possible implementation mechanism for cases that do not fit naturally within a project boundary, such as cross-project or temporal compaction.

However, the investigation should first establish whether the project-as-archive-index model is sufficient.

The question is not a generic competition between two equally preferred designs. It is whether the leading project-based model satisfies the requirements, and what exceptional cases might require another mechanism.

---

## 5. Lifecycle Model

### 5.1 Proposed lifecycle

The working lifecycle is:

`raw → ready → active → done → archived`

This aligns with the natural lifecycle of a project and the PARA model:

- Projects represent active, outcome-oriented work.
- Completed projects transition into Archives.
- Archived projects remain part of the knowledge system but no longer belong to the active working set.

The exact applicability of every state to every record type remains subject to investigation.

### 5.2 Done versus archived

These states represent different points in the lifecycle:

- **Done:** The work itself is complete.
- **Archived:** The record has passed through the archival process and is no longer part of the active working set.

A completed record may remain unarchived temporarily because it is part of a larger active project or because its project has not yet undergone exit refinement.

Archiving is the natural terminal lifecycle state, rather than a separate status layered alongside the ordinary work lifecycle.

### 5.3 Project closure and archival

The existing recursive project-closure invariants are a promising foundation.

Investigate whether the transition into archival should:

1. Require all descendants to be closed.
2. Trigger or require exit refinement.
3. Establish the recovery point and archive index.
4. Verify that important knowledge has been promoted or retained.
5. Compact eligible child records.
6. Mark the project archived only after the operation has been verified.

The exact ordering and whether some steps should be independently resumable must be determined by the spike.

---

## 6. Exit Refinement

### 6.1 Definition

**Exit refinement:** The process of evaluating a project, version, milestone, or other lifecycle boundary before it is considered complete, to determine what information should be retained, promoted, consolidated, or compacted.

Exit refinement is the knowledge-oriented counterpart to ordinary work refinement.

- Ordinary refinement prepares work for execution.
- Exit refinement prepares completed work and accumulated knowledge for transition into the next lifecycle phase.

It is a judgment-oriented process. It identifies what should happen; it does not itself perform the mechanical compaction operation.

### 6.2 Project completion as a natural trigger

Project completion is a particularly useful boundary because the project already provides a natural scope for its child work items.

When a project is ready to close, AKM could surface an exit-refinement opportunity.

For example:

```text
Project completed: Authentication Rewrite
42 child work items found.

Potential archival:
- 31 completed records eligible for compaction
- 7 records potentially worth retaining
- 4 records referenced by other project knowledge

Would you like to run exit refinement?
```

The numbers are illustrative.

The system should not automatically discard records simply because a project has been completed.

### 6.3 Major version transitions

Major software versions are another useful boundary.

For example:

- At V0 → V1, selectively compact obsolete or low-value V0 work while retaining useful decisions and architectural history.
- At V1 → V2, compact remaining V0 material and selectively retain relevant V1 material.
- At V2 → V3, evaluate older V1 material while preserving useful V2 context.

A major version transition should trigger evaluation of the outgoing version, not indiscriminate deletion.

### 6.4 Temporal boundaries

The same concept applies to personal knowledge management.

Potential boundaries include:

- Weekly → monthly consolidation
- Monthly → quarterly consolidation
- Quarterly → yearly consolidation

Hierarchical consolidation must preserve the ability to recover the original records through successive archival layers.

### 6.5 Exit refinement questions

Exit refinement should evaluate:

- Which decisions still explain the current architecture or operating model?
- Which assumptions were invalidated?
- Which constraints remain relevant?
- Which issues were temporary?
- Which work items are superseded?
- Which knowledge should be promoted into durable artifacts?
- Which unresolved questions should become future work?
- Which records are still referenced by active work?
- Which historical records no longer need to remain individually present?

---

## 7. Context-Local Knowledge Capture

### 7.1 Working concept

**Context-local knowledge capture** is the practice of evaluating and capturing potentially durable knowledge at the moment work is completed, while the relevant context is still available.

This is a working term, not yet a finalized domain term.

The underlying principle is that a small amount of thinking at completion can prevent a much larger reconstruction cost later.

For example, it may take very little effort immediately after a task to recognize that a discovery changes an architectural assumption. Reconstructing that significance months later from a collection of closed tickets may be substantially more expensive.

### 7.2 Completion ceremony

AKM should investigate a lightweight completion or exit-interview interaction for both humans and autonomous agents.

Potential prompts:

- What did we learn that was not obvious beforehand?
- Did we establish or invalidate an assumption?
- Did we make or change an architectural decision?
- Is anything here worth preserving in a durable knowledge artifact?
- Does this work supersede or relate to another record?
- Did we discover unresolved work that should be carried forward?

The ceremony should be adaptive:

- Trivial tasks should not incur unnecessary overhead.
- Significant investigations may warrant richer capture.
- Agents should be able to perform the process autonomously.
- AKM should distinguish between a proposed knowledge artifact and one that has actually been created.

### 7.3 Relationship to exit refinement

Context-local capture and exit refinement solve different problems.

- Context-local capture collects signals while context is fresh.
- Exit refinement evaluates accumulated information at a larger lifecycle boundary.

The first distributes judgment over the lifecycle. The second synthesizes that information before archival.

This reduces the burden on exit refinement and makes archival decisions more reliable.

---

## 8. Record Relationships and Graph Integrity

### 8.1 Existing structural graph

AKM already has a project hierarchy that forms a directed acyclic graph (DAG).

It represents containment and ownership:

- What belongs to a project?
- What must be completed before a project can close?
- What is the recursive scope of a project?

The existing project and dependency invariants provide a foundation for enforcing structural integrity.

### 8.2 Cross-dimensional record links

The hierarchy does not represent every meaningful relationship.

AKM needs to investigate a general record-link primitive that can connect records across project boundaries and record types.

Potential relationship semantics include:

| Relationship   | Meaning                                            |
| -------------- | -------------------------------------------------- |
| `relates_to`   | Records have a meaningful connection               |
| `supersedes`   | One record replaces another                        |
| `derived_from` | One artifact originated from another               |
| `documents`    | A knowledge artifact documents work or its outcome |
| `blocks`       | One work item depends on another                   |

These are examples for exploration, not a finalized vocabulary.

### 8.3 Different relationships have different invariants

Not every relationship should be forced into the containment DAG.

For example:

- Containment must preserve the project's structural acyclicity.
- Dependencies need their own validation rules.
- A generic `relates_to` link may be bidirectional and may participate in cycles without violating structural integrity.
- Supersession and provenance relationships may require stronger directional semantics.

The system should distinguish these relationship classes rather than treating every link as an equivalent edge.

### 8.4 Links to archived records

Archiving raises a reference-integrity question.

Investigate how retained records should resolve links to records that have been compacted.

Potential behavior:

- Links between retained records remain intact.
- Links to archived records resolve through the relevant archive index.
- Links that are no longer operationally relevant may be summarized or removed during exit refinement.
- Recovery metadata is maintained through dedicated, validated archival structures rather than ordinary user-editable semantic links.

Recoverability must never depend on a casual relationship link that can be removed or changed without preserving the recovery chain.

### 8.5 Backlink weighting

Backlinks are not equally valuable.

References from open work should carry more weight when evaluating whether a record should remain individually available.

In particular:

- Raw, ready, and active records indicate ongoing operational relevance.
- Done records may indicate historical relevance.
- Archived records may contain only historical references.

The spike should investigate lifecycle-aware backlink weighting as an input to retention decisions.

A backlink from an active task should generally be stronger evidence for retention than a backlink from a closed task.

---

## 9. Compaction and Recovery Model

### 9.1 Compaction as a mechanism

Compaction is the operation that reduces the active working set while preserving important knowledge and historical recoverability.

It may involve:

- Synthesizing a project summary.
- Promoting durable knowledge into appropriate artifacts.
- Creating or updating a machine-maintained archive index.
- Recording a Git recovery point.
- Recording the exact paths and identities of compacted records.
- Removing eligible records from the active working tree.
- Verifying the resulting archive.

### 9.2 Git-based recovery

Git is the likely recovery substrate.

An archive index should identify:

- The recovery commit.
- The exact original paths of compacted records.
- Stable record identifiers.
- The compaction operation or provenance.
- Any relevant nested archive references.

A recovery commit must represent the exact pre-deletion state of the records being compacted.

Distinguish:

- `created_at_commit`: when a record was originally created or committed.
- `recovery_commit`: the commit that preserves the exact state required to restore it.

Before removal, the operation must verify that all target files are committed and present in the recovery commit.

### 9.3 Illustrative archive index

A controlled Markdown section might look like:

```markdown
## Compaction Index

Schema Version: 1
Recovery Commit: `a81c92f...`

Compacted Records:
- ID: AKM-101
  Path: `work/done/AKM-101.md`
- ID: AKM-104
  Path: `work/done/AKM-104.md`

Record Count: 2

## End Compaction Index
```

This is illustrative, not a final schema.

The section is generated and maintained by AKM. It should be parsed and validated by the domain layer, not treated as an informal block of prose.

YAML frontmatter should remain clean; specialized archive metadata belongs in a controlled Markdown section.

The UI may render this section as structured, read-only information.

### 9.4 Descriptive synthesis versus deterministic metadata

Separate two kinds of archival content:

**Deterministic operational metadata**

- Recovery commit
- Record IDs
- Original paths
- Record counts
- Nested archive references
- Schema version
- Compaction provenance

This should be generated and validated by AKM.

**Descriptive knowledge synthesis**

- Project outcomes
- Significant decisions
- Architectural changes
- Lessons learned
- Current constraints
- Unresolved questions
- Useful historical context

This may require human or agent judgment and synthesis.

The summary does not replace the historical evidence or the recovery index.

### 9.5 Safe and resumable operations

Compaction must be safe to interrupt and resume.

Investigate:

- Precondition validation.
- Recovery-point creation.
- Verification that all target records are recoverable.
- Index generation.
- File removal.
- Post-operation verification.
- Rollback behavior.
- Idempotency and resumability.

The project should not transition to archived until the archival operation has reached a verified state.

---

## 10. Recovery Graph and Integrity Invariants

The logical compaction structure may form a DAG:

Original records → first-level archive indexes → higher-level archive indexes.

A higher-level archive may summarize or index lower-level archives while preserving the ability to recover the original records.

The logical recovery graph and the Git commit graph are separate structures. A valid Git history does not automatically guarantee valid logical recovery.

### 10.1 Required invariants to investigate

1. **Acyclicity**  
   Archive references must not create direct or transitive cycles.

2. **Recoverability and reachability**  
   Every compacted original record must remain reachable through a valid recovery chain.

3. **Stable identity**  
   Original record IDs must remain identifiable after physical removal, relocation, consolidation, or restoration.

4. **Reference integrity**  
   Recovery indexes must not point to missing records, invalid commits, or unavailable recovery locations.

5. **Provenance**  
   The system must be able to identify which archival operation compacted a record and which recovery point preserves it.

6. **Safe nesting**  
   Higher-level archival operations must preserve the recovery information of lower-level archives.

7. **Independent restoration**  
   A record should be restorable without requiring unrelated records or entire projects to be restored unnecessarily.

### 10.2 Relationship to existing AKM invariants

Investigate whether the project's existing recursive hierarchy and closure rules can enforce most of these requirements.

In particular, determine whether the existing DAG structure can serve as the logical archive graph, or whether a separate recovery relationship is necessary for some cases.

The goal is to reuse existing invariants wherever possible without assuming that structural containment automatically guarantees recovery integrity.

---

## 11. Retention Policy and Candidate Signals

The spike should investigate a practical, explainable retention policy rather than prematurely designing a sophisticated scoring algorithm.

Potential signals include:

- Backlinks from raw, ready, and active work.
- Backlinks from done or archived work.
- Whether a record's knowledge has been promoted into an ADR, specification, or architecture document.
- Whether a decision remains relevant to the current system.
- Whether a record has been superseded.
- Whether a record belongs to a recently completed project or milestone.
- Whether it contains unresolved questions or important constraints.
- Whether it is already represented by a higher-level summary or archive.

A useful initial policy may be rule-based and transparent.

The goal is to identify which records deserve continued individual availability and which can be safely represented through durable knowledge artifacts and recoverable archival history.

---

## 12. Operational Considerations and Failure Modes

The investigation must consider the operational lifecycle of archival data.

### Git and repository behavior

- Git garbage collection and object retention.
- Branch deletion and history rewriting.
- Repository migration or cloning.
- Recovery when the original branch no longer exists.
- Dirty working trees.
- Missing or uncommitted target files.
- Recovery commits that do not actually contain every target record.

### Nested and successive archives

- Parent projects containing archived child projects.
- Archiving a parent whose descendants have already been compacted.
- Higher-level archives referencing lower-level recovery indexes.
- Restoration of a child project independently.
- Modification or restoration of a lower-level archive after a higher-level archive references it.

### Interruption and failure

- Process termination after creating the recovery commit.
- Process termination after writing the archive index.
- Partial file removal.
- Failed verification.
- Repeated execution of the same archival operation.
- Rollback after partial completion.

### Long-term integrity

- Archive index corruption or manual edits.
- Schema evolution.
- Record ID stability.
- Broken references.
- Recovery after repository relocation.
- Whether archive indexes themselves eventually become eligible for higher-level archival.

---

## 13. Proposed Investigation Decomposition

The spike should be decomposed into natural investigation units. These are conceptual workstreams, not yet implementation tickets.

Each workstream should produce findings, explicit recommendations, unresolved questions, and any necessary follow-up spikes.

### A. Lifecycle and archival semantics

**Question:** What does it mean for an AKM record or project to move from active work into an archive?

Investigate:

- The proposed `raw → ready → active → done → archived` lifecycle.
- Which record types participate in each state.
- The distinction between done and archived.
- Project closure versus archival.
- Whether archival is a terminal lifecycle state.
- Preconditions and postconditions for archival.
- How existing project invariants can be reused.

**Expected output:** A lifecycle model and transition/invariant specification.

### B. Context-local knowledge capture

**Question:** What information should humans and agents capture when work is completed, while context is fresh?

Investigate:

- Completion ceremonies and exit interviews.
- Adaptive prompts based on task significance.
- Agent-autonomous knowledge capture.
- Durable knowledge promotion.
- Recording unresolved questions.
- Relationship and supersession signals.
- How to avoid excessive ceremony for trivial work.

**Expected output:** A lightweight completion-capture model and proposed interaction contract.

### C. Exit refinement and retention policy

**Question:** How should AKM evaluate accumulated work at project, version, milestone, and temporal boundaries?

Investigate:

- Exit-refinement triggers.
- Project completion as a natural archival boundary.
- Major-version transitions.
- Temporal consolidation.
- Backlink weighting by lifecycle state.
- Knowledge promotion and supersession.
- Candidate retention rules.
- Human approval versus agent autonomy.

**Expected output:** An exit-refinement process and an initial explainable retention policy.

### D. Project-based archive and compaction architecture

**Question:** Can the existing project record become the durable archive and recovery index for its child work?

Investigate:

- Reuse of project identity and hierarchy.
- Reuse of recursive closure invariants.
- Archive index representation.
- Nested projects and nested archives.
- Cross-project and temporal compaction cases.
- Whether a separate compaction record is needed for exceptional cases.
- How much existing project machinery can be reused.

**Expected output:** An architectural recommendation centered on the project-as-archive-index hypothesis, including its limitations and any justified exceptions.

### E. Record relationships and graph integrity

**Question:** How should AKM represent and validate relationships across records without confusing semantic links with structural containment?

Investigate:

- A general record-link primitive.
- Relationship types and directionality.
- Structural DAG versus semantic relationship graph.
- Dependency and supersession constraints.
- Links to archived records.
- Backlink discovery and weighting.
- Recovery graph reachability and acyclicity.

**Expected output:** A relationship model and a set of graph-integrity invariants.

### F. Compaction mechanics and Git recovery

**Question:** How can archival safely remove records from the active working set while preserving exact historical recovery?

Investigate:

- Recovery commit creation.
- Exact path and stable ID indexing.
- Archive index generation and validation.
- Removal and verification.
- Resumability and idempotency.
- Rollback.
- Nested recovery.
- Git retention and repository lifecycle concerns.

**Expected output:** A mechanical compaction design, recovery contract, and failure-handling specification.

### G. Integration and conformance

**Question:** How should the lifecycle, knowledge, relationship, and recovery models fit together across AKM interfaces and agent clients?

Investigate:

- Domain-layer enforcement of invariants.
- CLI and UI operations.
- Markdown representation and parsing.
- Agent workflows.
- Consistent behavior across different model and agent clients.
- Tests for interrupted operations and cross-agent recovery.

**Expected output:** An integrated domain contract and a set of conformance scenarios.

---

## 14. Suggested Investigation Dependencies

The workstreams should be separable, but their outputs have useful dependencies.

A suggested dependency structure:

```text
A. Lifecycle and archival semantics
        |
        +--------------------+
        |                    |
B. Context-local       E. Record relationships
   knowledge capture      and graph integrity
        |                    |
        +----------+---------+
                   |
          C. Exit refinement
          and retention policy
                   |
          D. Project-based archive
          and compaction architecture
                   |
          F. Compaction mechanics
          and Git recovery
                   |
          G. Integration and conformance
```

This is a suggested investigation sequence, not a requirement that every workstream finish before another can begin.

For example:

- Graph integrity can define invariants independently of the final archive representation.
- Knowledge capture can be investigated independently of Git recovery.
- The project-as-index hypothesis can be explored early using the existing project model.
- Mechanical Git recovery should not be finalized until the archive representation and recovery invariants are understood.

---

## 15. Expected Outcomes of the Discovery Spike

The discovery spike should produce:

1. **Lifecycle specification**  
   Record states, transitions, archival semantics, and invariants.

2. **Knowledge-capture model**  
   A proposal for context-local capture at work completion.

3. **Exit-refinement model**  
   Triggers, process, retention signals, and decision responsibilities.

4. **Archive architecture recommendation**  
   A reasoned recommendation for using projects as archive indexes, including any exceptions.

5. **Record-link specification**  
   Relationship semantics and constraints across the structural and semantic graphs.

6. **Recovery specification**  
   Recovery index structure, stable identity, reachability, provenance, and nested recovery rules.

7. **Compaction operation design**  
   Preconditions, execution steps, verification, interruption, resumption, and rollback.

8. **Implementation decomposition**  
   Smaller implementation projects and tasks derived from the findings.

The spike should not prematurely implement the entire archival system. Its purpose is to resolve the important domain and architectural questions well enough that implementation can be decomposed confidently.

---

# Appendix A: Working Glossary

| Term                                | Working definition                                                                                                                                  |
| ----------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Archiving**                       | The lifecycle transition that moves completed work out of the active working set and into the archive.                                              |
| **Compaction**                      | The mechanism for consolidating historical records, preserving important knowledge, and maintaining recovery information during archival.           |
| **Exit refinement**                 | Evaluation of accumulated work and knowledge at a lifecycle boundary to determine what should be retained, promoted, consolidated, or compacted.    |
| **Context-local knowledge capture** | Capturing potentially durable knowledge while the context in which it was discovered is still available.                                            |
| **Archive index**                   | Machine-maintained information retained by an archived project to identify and recover compacted child records.                                     |
| **Recovery point**                  | A Git commit preserving the exact state of records before compaction.                                                                               |
| **Structural DAG**                  | The acyclic project containment hierarchy and its associated structural invariants.                                                                 |
| **Semantic relationship graph**     | Cross-dimensional links expressing relationships between records that are not necessarily parent-child relationships.                               |
| **Stable identity**                 | The ability to identify a record consistently even after it has been removed, relocated, consolidated, or restored.                                 |
| **Knowledge promotion**             | Moving important information from transient work history into a durable knowledge artifact such as an ADR, specification, or architecture document. |
| **Nested archive**                  | An archived project whose descendants may themselves include archived projects and recovery indexes.                                                |

# Appendix B: Architectural Distinctions to Preserve

These distinctions should remain explicit during design and implementation:

1. **What vs. when vs. how vs. recoverability:** retention policy, lifecycle trigger, operation mechanics, and integrity are separate questions.
2. **Done vs. archived:** completion of work is not the same as completion of its archival lifecycle.
3. **Work history vs. durable knowledge:** the original task is evidence; important knowledge should be promoted into the appropriate artifact.
4. **Exit refinement vs. compaction:** judgment about what to preserve is separate from the mechanical operation.
5. **Structural DAG vs. semantic links:** containment and dependency invariants are not interchangeable with general record relationships.
6. **Git graph vs. logical recovery graph:** Git history alone does not guarantee that every original record remains logically recoverable.
7. **Operational metadata vs. descriptive synthesis:** recovery indexes should be deterministic and validated; summaries may require human or agent judgment.
8. **Archiving vs. deletion:** removal from the active working set must not mean loss of identity or recoverability.

# Appendix C: Open Architectural Questions

These questions should be resolved by the relevant investigation workstreams rather than assumed in advance.

- Does every record type use the same archival lifecycle?
- Does archiving always occur at project scope, or are there legitimate record-level exceptions?
- What exactly must be present before a project can transition to archived?
- Should exit refinement be required, recommended, or conditionally required?
- How should the system distinguish durable knowledge from useful historical evidence?
- What is the initial retention policy, and how should lifecycle-aware backlinks influence it?
- Which relationship types should be first-class?
- How should links to archived records resolve?
- Can the existing project hierarchy guarantee sufficient recovery reachability?
- What archive metadata must be retained when nested projects are compacted?
- How should archive indexes evolve as their schema changes?
- What Git retention assumptions can AKM safely make?
- Can every archival operation be made resumable and idempotent?
- Under what circumstances should an archive itself be compacted?

---

**Guiding architectural principle:** AKM should make it inexpensive to capture knowledge while it is fresh, straightforward to evaluate that knowledge at natural lifecycle boundaries, and safe to archive completed work without sacrificing durable understanding or historical recoverability.
