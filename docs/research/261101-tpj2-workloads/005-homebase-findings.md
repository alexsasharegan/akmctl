---
id: 261102-m7q2
title: Homebase — Findings
summary: Records the findings of the Homebase workload investigation, examining persistent responsibilities, heterogeneous record lifecycles, progressive structuring, asset provenance, and agent-assisted knowledge refinement.
status: draft
kind: research
created: 2026-11-02
updated: 2026-11-02
parent: 261101-tpj2
---

# Homebase — Findings

## 1. Investigation Summary

Homebase was investigated as a long-lived personal knowledge and life-management environment containing recurring responsibilities, finite projects, enduring resources, external evidence, and historical information.

The investigation began with the distinction between persistent Areas and finite Projects, then expanded into vehicle maintenance, external assets, provenance, and the different lifecycles required by heterogeneous records.

The central finding is that **AKM must distinguish between organizing knowledge, managing work, and governing record lifecycles.** These responsibilities overlap, but they are not interchangeable.

A document can participate in the knowledge network without being governed by a work lifecycle. An enduring responsibility can organize many projects without itself becoming a project. An external object can have a useful digital representation without becoming an object that AKM owns or stores.

The investigation also established that agent-assisted structure should be inexpensive enough to change the practical starting point for knowledge capture. Users should not have to choose between manually maintaining elaborate structures and accumulating unstructured notes.

The accompanying `004-homebase-fixtures.md` preserves the scenarios explored and provides concrete behavioral examples for these findings.

## 2. Observed Requirements

### 2.1 Persistent Knowledge Anchors

An ongoing responsibility needs a durable, discoverable place in the knowledge network.

The Taxes and Vehicles examples demonstrate that these anchors connect:

- Successive projects and executions.
- Enduring entities and resources.
- Historical records.
- Reusable procedures.
- Related responsibilities and supporting information.

An Area hub is not simply a project that never closes. It is a persistent knowledge node whose purpose is to organize and connect information across time.

Its content can change. It may be edited, trimmed, reorganized, and kept current. Historical records and their relationships preserve access to information that is no longer prominent in the hub.

The distinction is between maintaining a useful current representation and preserving historical discoverability. The hub does not need to contain every historical detail to remain a durable organizational anchor.

PARA is a useful organizational lens, but its categories should not automatically become AKM record types or lifecycle requirements.

### 2.2 Different Records Have Different Lifecycle Semantics

Homebase contains several kinds of information with fundamentally different temporal characteristics.

A Project represents finite work toward an outcome. An Area represents an ongoing responsibility. A Resource preserves enduring reference information. An asset proxy represents an external or physical object.

These records should not be forced into one universal work lifecycle.

A completed project and an inactive resource are not equivalent. An Area may remain active indefinitely without becoming a project. An asset proxy may remain useful while its location changes.

AKM should support a subset of strongly typed records with defined lifecycle semantics and corresponding guarantees.

This does not imply that every document needs a managed lifecycle. Many documents need only durable persistence, relationships, and discoverability.

### 2.3 Historical Evidence and Reusable Knowledge Are Distinct

Recurring responsibilities produce successive executions.

Annual tax preparation illustrates the distinction between:

- The persistent Taxes responsibility.
- The current year's preparation project.
- Historical preparation projects.
- The reusable tax procedure.
- Supporting evidence and records.

The previous year's execution is historical evidence. The reusable procedure is accumulated knowledge intended to inform future executions.

These artifacts may be related, but they should not be conflated.

A completed execution does not automatically become the reusable procedure, and updating a procedure should not rewrite historical evidence.

Conditional procedures should also remain available when temporarily inapplicable. A workflow for business income or Schedule K-1 reporting may be useful again even if the current year's return does not require it.

### 2.4 External Evidence Must Participate Without Becoming an AKM Record

Important information frequently originates outside the vault.

Examples include:

- Emails and contractor estimates.
- Google Drive documents.
- Photographs and scanned receipts.
- Dealership service reports.
- Encrypted tax documents.
- Physical documents and objects.

AKM should be able to preserve references to these sources and organize useful information derived from them.

A reference may be a URL, filesystem path, connector-specific identifier, or human-readable breadcrumb. It does not need to be a fully resolvable URL to provide value.

Resolution depends on the reference, the available resolver, the user's environment, and applicable permissions.

The system should not assume that every external source can be opened or read automatically.

### 2.5 Digital Proxies Preserve Useful Representations

Some external or physical objects deserve a durable digital representation.

A vehicle title, for example, can have an AKM record describing the document, identifying the associated vehicle, preserving useful extracted information, and indicating where the physical title is stored.

The proxy is not the original object. It is a stable, discoverable representation that makes the object useful to the knowledge network.

Its identity persists when its location changes.

The current location can be updated while the history of previous locations remains available. The same principle applies to other mutable information about the represented object.

A proxy can remain useful even when the original source is temporarily inaccessible.

### 2.6 Structured Capture Should Be Inexpensive

The vehicle maintenance investigation demonstrated that information can begin as a chronological note and later evolve into a hub with individual service records and a queryable collection.

However, this progression should not be interpreted as a requirement to begin with unstructured notes.

When an agent can cheaply extract a dealership report into a structured service record, there is little reason to require the user to maintain a manually curated junk drawer first.

The user should be able to provide source material and delegate extraction, organization, linking, and routine bookkeeping to an agent.

The important distinction is between the cost of creating structure and the decision to create additional artifacts.

Structure should be inexpensive. Additional artifacts should still have a meaningful purpose.

### 2.7 Knowledge Must Be Discoverable Across Changing Contexts

Information can remain useful while dormant or become relevant again under different circumstances.

A historical service recommendation may become important when planning a repair. A previous payroll provider may matter when retrieving an old W-2. An old procedure may become applicable when circumstances change.

Inactivity does not establish obsolescence.

Likewise, historical relevance and current operational applicability are different concerns. A record can cease to guide current work while remaining valuable for verification, historical reconstruction, or future decisions.

## 3. Information Lifecycle and Evolution

### 3.1 Execution Produces Knowledge

Work items and projects naturally accumulate observations, decisions, evidence, and references.

Not every discovery needs to be fully classified at the moment it occurs. The user is trying to accomplish the work, not continuously maintain the knowledge architecture.

Work history provides a natural place for observations to accumulate while context is fresh.

At meaningful boundaries, those observations can be evaluated for durable value.

### 3.2 Exit Refinement Is a Knowledge-Maintenance Opportunity

Project completion provides a natural opportunity to examine what was learned.

This refinement is distinct from simply marking work complete or archiving records.

It may identify:

- Changes to persistent hubs.
- New or modified relationships.
- Reusable procedural knowledge.
- New resources or asset proxies.
- Information worth extracting from historical records.
- Historical relationships that should remain discoverable.
- No changes worth making.

Refinement should not be a mandatory, elaborate ceremony after every task or work item.

Its depth should scale with the significance of the boundary, the maturity of the process, and the discoveries made.

An uneventful execution may require no meaningful refinement. A substantial annual project may justify a deliberate review.

### 3.3 Refinement Should Adapt to Process Maturity

A first-time responsibility and an established recurring process require different kinds of assistance.

A novel process may require conversational discovery:

- What information is needed?
- Which entities and resources matter?
- What procedures should be established?
- What organizational structure is useful?

An established process may require targeted updates:

- Replace a payroll provider.
- Update an access procedure.
- Preserve an old relationship.
- Create a new record for an emerging resource.

These are different modes of the same refinement responsibility, not necessarily different lifecycle mechanisms.

### 3.4 Agents Propose; Users Authorize

The agent should examine relevant work history and relationships, identify candidate knowledge changes, and present focused proposals or questions.

It should not depend entirely on the user remembering what deserves preservation.

The user remains responsible for approving meaningful changes.

The user's preference for an independent refinement agent is a useful interaction pattern: a fresh context can review the work without simply inheriting the assumptions of the agent that performed it.

AKM's responsibility is to provide the lifecycle opportunities, relevant context, relationships, and safe capabilities. The agent supplies the contextual synthesis and user-specific conventions.

This is an agent-assisted responsibility, not a deterministic retention algorithm.

## 4. Important Discoveries

### 4.1 Knowledge Organization and Lifecycle Authority Are Separate

The investigation reinforced a distinction already emerging in the adaptive-retention exploration.

A record can participate in the logical knowledge network without AKM governing its lifecycle.

Similarly, a persistent hub can be an important organizational anchor without being a specialized lifecycle-managed record.

Tracking, graph participation, lifecycle authority, and physical persistence should remain distinct concerns.

### 4.2 Structure and Formalization Have Different Costs

Progressive formalization remains valuable, but its interpretation needs refinement.

The 2026 Honda Pilot example demonstrated that a simple chronological note can evolve into a hub and structured service collection.

The F-150 example demonstrated that an individual service record can remain the correct home for declined findings until a separate planning activity emerges.

These examples establish two complementary principles:

- Make useful structure inexpensive to create.
- Create additional artifacts when they serve a distinct purpose.

The system should not preserve unnecessary manual friction in the name of progressive formalization.

Nor should it create a new record for every fact merely because doing so is technically easy.

### 4.3 Extraction Is Not Relocation

When the F-150 suspension findings become the subject of active planning, selected information is extracted into a planning note.

The original service record remains historical evidence.

The planning note represents a new purpose and preserves a reference to its source.

This is not simply moving information from one note to another. It is creating a new representation for a different activity while retaining the original evidence.

### 4.4 Provenance Is Broader Than Document Relationships

Provenance can be a simple reference in a work-item log, a link to an external file, a breadcrumb to an email, or a durable relationship to another AKM record.

A useful reference need not be a formally structured provenance object.

A digital proxy is appropriate when an external or physical object deserves its own stable identity and useful representation.

The investigation considered a typed reference model using a scheme and value, with resolvers for recognized schemes. This is a plausible direction, but the workload does not establish a need for a comprehensive provenance subsystem.

The important requirement is that sources remain identifiable and useful, with resolution capabilities determined by the available environment.

### 4.5 A Proxy Is a Stable Logical Reference

The vehicle-title example exposed an important property of digital proxies.

The identity of the represented object must not depend on its current location.

A proxy can retain its identity while its physical location, external URI, or other means of access changes.

Its metadata can provide a current projection, while its history preserves meaningful changes.

The proxy should not be confused with the source or with a particular locator.

### 4.6 Relationships Can Express Succession Without Proxies

The payroll-provider example demonstrated that changing services do not automatically require asset proxies.

Distinct provider records can link to the Employer. Historical provider relationships remain discoverable, while the Employer's current relationship identifies the active provider.

The graph itself can express the association and its history.

This is a useful reminder not to introduce a proxy abstraction where ordinary records and relationships already provide the necessary behavior.

### 4.7 The Document Environment Must Remain Open

AKM operates in a document environment where users and agents can add arbitrary YAML front matter and content.

It cannot assume complete control over every field or convention.

The system should be opinionated about the semantics it manages while tolerating information outside that managed domain.

This suggests three levels of interpretation:

1. **Managed semantics:** Strongly typed record categories, defined lifecycle behavior, and fields whose invariants AKM governs.
2. **Cooperative conventions:** User-defined metadata and classifications that AKM can use for generic filtering, grouping, and presentation.
3. **Unmanaged document space:** Arbitrary content and fields whose semantics AKM does not claim to understand.

These are levels of interpretation and authority, not necessarily separate document classes.

AKM should be conspicuous when users or agents cross into managed territory incorrectly. It should not silently accept conflicting declarations or allow unmanaged fields to override managed invariants.

At the same time, unknown fields and document content should not make the vault unusable.

### 4.8 Type, State, and Presentation Are Distinct

A record's type determines what kind of thing it represents and which managed lifecycle semantics apply.

Its state describes its current condition within that lifecycle.

Its presentation determines how it is surfaced and organized for a particular use.

A generic Resources view might filter by user-defined dimensions such as `resource_kind`, without AKM needing to implement specialized lifecycle behavior for every subtype.

A Project view can be more opinionated because Projects participate in a defined work lifecycle.

This provides extensibility in classification and presentation without requiring arbitrary user-defined lifecycle state machines.

## 5. Architectural Concepts Validated

The following concepts from the preceding architectural exploration proved useful in the Homebase workload.

### 5.1 Logical Relationships Independent of Physical Storage

Meaningful relationships support discovery and navigation even when the connected objects are stored outside the vault.

The relationship network can include AKM records, external resources, and references to physical objects.

Graph participation does not imply lifecycle ownership or physical persistence.

### 5.2 Progressive Formalization

Relationships and classifications can begin as ordinary links, references, or ad hoc metadata.

Recurring patterns can justify more formal semantics.

The user should not need a comprehensive ontology before capturing useful information.

However, the availability of inexpensive agent-assisted structure means that progressive formalization should not be confused with requiring unstructured capture.

### 5.3 Contextual Relevance and Rediscovery

The usefulness of information depends on the current purpose.

A historical record may be irrelevant to routine work but important for a later repair, tax question, or verification.

Relationships and durable references help preserve discoverability without requiring every historically useful record to remain prominent in the active working context.

### 5.4 Exit Refinement

Project completion and other meaningful boundaries provide opportunities to enrich durable knowledge, preserve provenance, and identify reusable information.

The Homebase workload validates refinement as a responsibility distinct from ordinary work execution.

It does not establish that refinement must be a large ceremony or a deterministic automated operation.

### 5.5 Attention Management

The goal is not to minimize the number of records at all costs.

A structured collection may contain more records than a single accumulated note while making relevant information much easier to retrieve.

The useful measure is the effort required to discover, interpret, and use the appropriate information.

Making hubs, records, and structured views inexpensive to maintain can improve exposure efficiency without requiring aggressive compaction.

## 6. Architectural Concepts Challenged

### 6.1 Areas as Projects That Never Close

This is an inadequate model.

An Area hub is a persistent organizational and knowledge node, not simply a Project with an indefinitely extended lifecycle.

Areas may have current or inactive conditions, but their semantics should not be inherited automatically from Projects.

### 6.2 Universal Work Lifecycles

A single lifecycle cannot adequately describe Projects, Areas, Resources, asset proxies, and historical event records.

AKM should provide a managed subset of strongly typed records with appropriate lifecycle semantics rather than treating every document as a work item.

### 6.3 Delaying Structure Until Retrieval Becomes Painful

This principle needs qualification.

It remains useful when deciding whether to extract information into an additional artifact.

It is not a reason to avoid creating structured records when an agent can do so cheaply during initial capture.

The system should minimize the cost of useful organization, not merely tolerate disorganization until it becomes painful.

### 6.4 Universal Provenance Structures

Not every source reference needs a dedicated structured field or provenance record.

A log entry, ordinary link, file path, or human-readable breadcrumb can be sufficient.

Formalization is justified when it provides meaningful capabilities such as reliable resolution, validation, or queryability.

### 6.5 Treating Every Durable Object as a Proxy

A proxy is useful when an external or physical object needs a stable digital representation.

It is not a general-purpose replacement for ordinary entities and relationships.

Changing payroll providers, for example, is better represented by distinct provider records and their relationships to the employer.

## 7. Unnecessary Machinery

The investigation did not establish a need for:

- A specialized lifecycle for every kind of document.
- A universal status vocabulary shared by all record types.
- A comprehensive ontology of Areas, Resources, Assets, and their subtypes.
- A dedicated asset-proxy engine separate from ordinary AKM record capabilities.
- A mandatory structured provenance object for every external reference.
- A universal resolver capable of opening every possible external source.
- A separate proxy for every changing service or provider.
- A mandatory hub and Base for every collection of related information.
- A new planning record whenever a recommendation remains unresolved.
- A mandatory retrospective after every task or work item.
- A deterministic algorithm that decides which information must be promoted into durable knowledge.

These mechanisms may be useful in particular circumstances, but the Homebase workload does not justify requiring them as universal infrastructure.

## 8. Simpler Alternatives

Several scenarios can be supported through combinations of existing document and relationship capabilities:

| Need                                        | Sufficient starting approach                               |
| ------------------------------------------- | ---------------------------------------------------------- |
| Persistent responsibility                   | Durable hub note with relationships                        |
| Historical service events                   | Individual event records linked to a vehicle               |
| Queryable event collection                  | Generic Base or equivalent view                            |
| External evidence                           | Ordinary URI, path, or breadcrumb                          |
| Physical object representation              | Ordinary AKM record with descriptive metadata and location |
| Changing provider                           | Distinct entity records and relationships                  |
| Selected information extracted for planning | New note linked to its source                              |
| User-defined classifications                | Front matter fields and generic views                      |
| Durable knowledge refinement                | Agent-assisted review at meaningful boundaries             |

These approaches should remain the starting point unless demonstrated requirements justify additional mechanisms.

## 9. Must Happen, Would Be Useful, and Unnecessary Machinery

### 9.1 Must Happen

The workload establishes the following behavioral requirements:

1. Persistent knowledge anchors must remain useful across multiple projects and executions.
2. Historical evidence must remain distinguishable from current procedures and representations.
3. Different kinds of records must not be forced into an inappropriate universal work lifecycle.
4. External and physical objects must be representable through useful references or digital proxies.
5. A proxy's identity must remain stable when its location changes.
6. Extracted information must retain a path back to its source when that source matters.
7. Historical relationships must remain discoverable when their current operational relevance changes.
8. AKM must enforce the invariants of the semantics it manages without claiming authority over arbitrary document content.
9. Agent-assisted organization must not require the user to perform extensive manual bookkeeping.

### 9.2 Would Be Useful

The following capabilities provide meaningful benefits:

1. Strongly typed record categories with distinct lifecycle semantics.
2. Generic views that filter and group records using user-defined metadata.
3. Convenient extraction of structured records from documents, photographs, and other evidence.
4. Recognized reference schemes with environment-dependent resolution capabilities.
5. Durable asset proxies with current metadata and historical logs.
6. Agent-assisted entry and exit refinement.
7. Independent refinement contexts that can review accumulated work.
8. Contextual exposure of historical records and supporting evidence.
9. Reusable procedures that preserve conditional and currently inapplicable workflows.

### 9.3 Unnecessary Machinery

The workload does not justify making the following universal requirements:

1. A lifecycle state machine for every document.
2. A mandatory record type for every organizational category or subtype.
3. A special proxy mechanism for every external service or changing entity.
4. A structured provenance field for every reference.
5. A complete resolver framework for arbitrary external sources.
6. Automatic extraction of every unresolved recommendation into a separate planning record.
7. A mandatory refinement ceremony at every completion boundary.
8. A universal ontology or fully user-extensible lifecycle engine.

## 10. Recommendations

The smallest coherent capability set that adequately supports Homebase is:

1. **Persistent records and relationships.** Support durable organizational anchors and meaningful relationships independently of work-item containment.
2. **A managed subset of record types.** Provide defined lifecycle semantics for the record categories where AKM can offer meaningful guarantees.
3. **Open metadata.** Tolerate arbitrary user-defined document content and metadata while enforcing managed invariants.
4. **Generic presentation.** Allow records to be filtered and grouped using metadata that AKM can understand without claiming ownership of its semantics.
5. **References to external objects.** Preserve useful locators and breadcrumbs without requiring external sources to become AKM records.
6. **Ordinary records as digital proxies.** Allow an AKM record to represent an external or physical object through descriptive metadata, relationships, and location information.
7. **Inexpensive structured capture.** Enable agents to transform source material into useful records and maintain their relationships without imposing manual bookkeeping.
8. **Lifecycle refinement opportunities.** Make accumulated work and its discoveries available for agent-assisted knowledge consolidation at meaningful boundaries.

These are behavioral capabilities, not a prescribed schema or implementation architecture.

The workload does not determine the exact record taxonomy, lifecycle state vocabulary, resolver interface, or generic view implementation.

Those decisions should be evaluated against the broader set of workloads rather than settled within Homebase alone.

## 11. Open Questions

The investigation leaves several details intentionally unresolved:

1. Which record categories warrant first-class lifecycle semantics?
2. What is the smallest useful set of managed lifecycle behaviors?
3. How should managed metadata and user-defined conventions be distinguished and validated?
4. What minimum field-shape understanding is necessary for useful generic views?
5. How should external references be represented when structured resolution is available?
6. What capabilities should recognized reference resolvers expose, and how should permissions and unavailable sources be handled?
7. What metadata and history are useful for a digital proxy without turning it into a specialized subsystem?
8. How should agent-assisted refinement opportunities be exposed without creating unnecessary ceremony?

These questions do not prevent the findings from establishing the behavioral boundaries of the Homebase workload.

## 12. Conclusion

Homebase demonstrates that a long-lived knowledge environment cannot be modeled entirely as a collection of projects and tasks.

It contains persistent responsibilities, enduring resources, historical events, external evidence, physical objects, and finite work. These things have different purposes, relationships, and lifecycle characteristics.

AKM should be opinionated about the record semantics and lifecycle guarantees it manages while remaining accommodating toward the open document environment in which it operates.

Its value comes from making information inexpensive to structure, connect, retrieve, and refine—not from requiring every piece of information to conform to a comprehensive ontology.

The most important distinction established by this investigation is:

**AKM must distinguish what a record represents, how it participates in the knowledge network, who governs its lifecycle, and how it is exposed to a particular context.**

Those concerns can reinforce one another without being collapsed into a universal record model.

The governing principle remains:

**The simplest solution that adequately supports the workload wins.**
