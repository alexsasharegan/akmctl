---
id: 261102-7KQ9
title: Homebase — Behavioral Fixtures
summary: Concrete scenarios and evolving information structures derived from the Homebase workload investigation, preserved as behavioral fixtures for evaluating AKM's design and implementation.
status: draft
kind: research
created: 2026-11-02
updated: 2026-11-02
parent: 261101-tpj2
---

# Homebase — Behavioral Fixtures

## 1. Purpose

This document preserves the concrete scenarios explored during the Homebase workload investigation.

The original workload describes a personal knowledge-management environment with persistent responsibilities, finite projects, enduring resources, and recurring activities. The investigation expanded those scenarios into more detailed examples of how information is captured, structured, connected, maintained, and rediscovered over time.

These examples are preserved as behavioral fixtures.

They provide concrete situations against which AKM's architectural concepts and eventual implementation can be evaluated. They are not a proposed schema, a collection of mandatory record types, or an implementation specification.

The fixtures deliberately include both ordinary information-management situations and more demanding cases involving external evidence, evolving relationships, agent-assisted refinement, and different record lifecycles.

### 1.1 Relationship to the Workload

This document extends the Homebase workload definition rather than replacing it.

- The workload establishes the environment and its general characteristics.
- These fixtures preserve the concrete scenarios developed through investigation.
- The corresponding findings document will record the conclusions drawn from those scenarios.

The fixtures should remain useful independently of any particular implementation.

### 1.2 Fixture Conventions

Each fixture describes:

- **Initial conditions:** The relevant records, information, and relationships.
- **Events:** What happens in the user's environment.
- **Expected behavior:** What the user and agent should be able to accomplish.
- **Evaluation points:** Properties that a design or implementation should preserve.

Examples of YAML, relationships, and record organization are illustrative. They should not be interpreted as commitments to particular field names or schemas.

---

## 2. Annual Responsibility: Tax Preparation

### 2.1 Initial Conditions

The user has an enduring Taxes responsibility.

The Taxes hub is a persistent knowledge node that connects annual tax projects, relevant entities, historical records, and reusable procedures.

The user has several related entities and resources:

- An Employer record.
- Investment account records.
- Historical tax returns.
- Supporting financial records.
- A collection of receipts and other evidence.
- A tax preparation procedure or runbook.

The hub is not itself an annual execution record. It remains useful across years and can be updated as circumstances change.

### 2.2 Annual Execution

The user begins a new annual tax preparation project.

The project contains work items such as:

- Collect income documents.
- Gather expense records.
- Review deductions.
- Prepare the return.
- File the return.

The project accumulates observations and evidence as the work proceeds.

The previous year's project remains available as historical context. It is not overwritten by the new execution.

### 2.3 Discovering a Changed Payroll Provider

During tax preparation, the user discovers that an expected W-2 has not arrived.

Investigation reveals that the employer changed payroll providers.

The old provider's portal may still be useful for historical documents. The new provider is responsible for current payroll documents.

The observation originates during ordinary work. The user should not need to interrupt the task to perform a separate knowledge-management ceremony.

At a suitable refinement boundary, an agent examines the work history and proposes relevant updates:

- Update the Employer hub with the new provider.
- Preserve the historical relationship with the previous provider.
- Record how to access the new provider's portal.
- Preserve relevant dates and context.
- Update the tax runbook if the change affects future preparation.

The user reviews and approves the proposed changes.

### 2.4 Annual Memory Consolidation

Once the return is filed, the project reaches a meaningful completion boundary.

A separate refinement activity examines the execution history for durable discoveries.

The amount of refinement should depend on what happened:

- An uneventful year may require little or no knowledge maintenance.
- A small procedural discovery may justify a lightweight update.
- A significant change to an employer, financial account, or filing procedure may warrant deliberate review.

The agent should identify potential changes rather than depend entirely on the user remembering them.

### 2.5 Conditional Procedures

The tax runbook contains procedures that are not applicable every year.

For example, the user previously dealt with business income and a Schedule K-1 associated with an LLC.

That procedure may be irrelevant to the current year's return.

It should not automatically be deleted or treated as obsolete. If circumstances change and an LLC is opened again, the procedure may become useful.

The historical execution remains evidence of what happened. The reusable procedure represents accumulated knowledge for future executions.

**Evaluation points**

- A persistent responsibility can connect multiple finite annual projects.
- Completing a project does not complete the responsibility.
- Historical executions remain available without being confused with current procedures.
- Knowledge can be promoted from work history into reusable guidance.
- Conditional procedures can remain useful while currently inapplicable.
- Refinement can scale with the significance of discoveries.

---

## 3. Vehicle Knowledge: From Accumulated Notes to Structured Records

### 3.1 Initial Conditions

The user owns multiple vehicles.

Initially, vehicle information is captured in a single note per vehicle.

For example, a note titled `2026 Honda Pilot` contains:

- Purchase date.
- VIN.
- Initial mileage.
- Chronological service entries.
- Service location and personnel.
- Work performed.
- Mileage at service.
- Recommendations.
- Costs.

This is a practical starting point. It allows the user to capture information without first designing a complete record system.

The note is effectively a chronological collection of vehicle history.

### 3.2 Progressive Structuring

Over time, the user finds that the accumulated information benefits from a more structured organization.

The vehicle becomes a hub note.

Individual service visits become separate dated records, generated from a reusable template.

An Obsidian Base exposes the service records as a queryable collection.

Conceptually:

```text
Vehicles/
├── Pilot/
│   ├── Pilot.md
│   └── Service/
│       ├── 2025-03-12-service.md
│       ├── 2025-09-18-service.md
│       └── 2026-04-02-service.md
├── F-150/
│   ├── F-150.md
│   └── Service/
│       └── ...
└── Vehicle Events.base
```

The hub provides a discoverable organizational anchor.

The service records preserve individual events.

The Base exposes the records in a structured view.

This structure is not required at the beginning. It becomes worthwhile because the agent can make the bookkeeping inexpensive.

### 3.3 Structured Ingestion

The dealership provides a printed service report containing findings, recommendations, and work performed.

The user photographs or screenshots the report and gives it to an agent.

The agent extracts the relevant information and creates a structured service record.

The record should preserve the original evidence or a reference to it.

The user should not have to manually transcribe every finding into a database-like format.

### 3.4 Retrieval

Later, the user asks:

- Who helped me at the dealership last time?
- When was the cabin air filter last changed?
- What recommendations did I decline?
- What did the dealership recommend during the last two visits?

The agent should be able to retrieve and synthesize the relevant service records.

**Evaluation points**

- An initially simple note can evolve into a hub and structured collection.
- The transition should not require the user to manually reconstruct existing information.
- Service records remain independent historical events.
- The hub provides a stable organizational anchor.
- Structured views make records easier to query without replacing their underlying documents.
- Original evidence remains distinguishable from extracted information.

---

## 4. F-150: Declined Findings and the Emergence of a Planning Record

This fixture explores how information can remain in its original context until a distinct purpose justifies extracting it.

### 4.1 Initial Service Visit

The F-150 receives a dealership inspection.

The service report identifies several issues, including substantial suspension work.

The user declines the suspension recommendations for the time being.

The service record contains:

- The inspection findings.
- Recommended work.
- Work performed.
- Work declined.
- Estimated costs.
- Mileage.
- Dealership and service personnel.
- The original service report.

The user drives relatively few miles and expects to address the suspension work around the time of an upcoming spring bonus.

At this point, no separate planning record is necessary.

The declined findings remain in the service record.

### 4.2 Simple Deferral

The user takes the F-150 in for another routine service before the bonus arrives.

The dealership repeats the suspension recommendations.

The user retrieves the previous service record and compares the findings.

The user discusses whether the vehicle has sufficient remaining mileage before the work needs to be performed.

The recommendations remain in the original service record.

Nothing requires extraction merely because the recommendations are still relevant.

### 4.3 Active Planning Begins

The user decides to investigate the suspension work more seriously.

They ask an agent:

> I want to make a planning note. Pull the last service record for my F-150. Focus on all the vehicle findings that were not performed on the vehicle.

The agent retrieves the service record and presents the declined findings.

The user filters the results:

- Suspension findings are relevant.
- The cabin visor recommendation is not relevant.

The user and agent then explore the repair:

- Can some parts of the work be performed without others?
- Would splitting the work result in duplicated labor?
- Would it be cheaper to perform everything together?
- What parts and labor are involved?
- How should the work fit within the available budget?

At this point, the information has acquired a distinct planning purpose.

The user creates a separate planning note.

### 4.4 Extraction and Provenance

The planning note contains the selected suspension findings and the information needed to evaluate the work.

It links back to the original service record.

The original service record remains unchanged as historical evidence of what the dealership found and what the user declined.

The planning note does not replace that evidence. It represents the new work of evaluating and arranging the repair.

### 4.5 Execution

Eventually, the user schedules and completes the repair.

A new service record documents what actually happened.

The planning note remains available as the history of how the work was considered, evaluated, and arranged.

**Evaluation points**

- A recommendation can remain in its original service record while merely deferred.
- Continued relevance does not automatically require a new record.
- A distinct planning activity justifies extracting selected information.
- Extraction should preserve provenance to the original service evidence.
- The original event record should not be rewritten to reflect later planning.
- Planning and execution are distinct activities with distinct records.

---

## 5. External Evidence and Asset References

Homebase contains important information whose sources are not necessarily AKM records or even files inside the vault.

External evidence may include:

- Google Drive documents.
- Photographs.
- PDFs.
- Email messages.
- Contractor estimates.
- Physical documents.
- Physical objects and their storage locations.

The following fixtures explore how these objects can participate in the knowledge network.

### 5.1 Contractor Estimate

A contractor provides an estimate for installing a soft-start mechanism on the home's air-conditioning system.

The estimate is approximately $2,000.

The user captures a small amount of information:

- What the proposed work is.
- The estimated cost.
- A reference to the contractor's email containing the quote.

The email may not have a convenient permanent URL.

A useful reference might instead be a breadcrumb:

```text
Email:
  Gmail
  Date: [message date]
  Subject: [estimate subject]
  Sender: [contractor]
```

The agent may be able to resolve this through an email connector. Otherwise, the breadcrumb remains useful to the user.

The important information—the estimated cost and nature of the work—is available without reopening the email.

The source remains available for verification when necessary.

### 5.2 Medical Expense Evidence

The user has a note recording approximately $1,200 in medical expenses.

The supporting evidence consists of photographs and PDFs of receipts stored in Google Drive.

The note records the useful aggregate information and references the supporting documents.

The user does not need every receipt reproduced inside the vault.

The evidence should remain retrievable when verification is necessary.

### 5.3 Encrypted Tax Filing Document

The user has a tax filing PDF that is password-encrypted.

Including the original document in an agent's working context may be inconvenient or impossible.

The user instead extracts selected information:

- Important filing figures.
- Relevant details needed for future tax preparation.
- Screenshots of pages that contain useful information but do not expose unnecessary personally identifiable information.

The resulting representation is a useful projection of the original document.

It does not pretend to be the complete original filing.

The original remains the authoritative source.

### 5.4 Physical Vehicle Title

The user possesses a physical vehicle title.

The title itself cannot be opened or searched by an agent.

The user creates a digital proxy record representing the title.

The proxy contains:

- What the document is.
- Which vehicle it concerns.
- Relevant extracted information.
- Where the physical title is stored.

For example:

```text
Vehicle Title — F-150

Physical location:
  User's safe
```

The proxy is a durable digital representation of the physical document.

It provides an agent with useful information without requiring access to the physical object.

### 5.5 Moving the Physical Document

The user moves the vehicle title from the home safe to a garage safe.

Later, the user leaves it with their parents for safekeeping.

The proxy retains its identity.

Its current location changes, and the movement can be recorded in its history.

The proxy can therefore answer both:

- Where is the title now?
- How did it get there?

The physical title is not recreated or replaced by the proxy.

### 5.6 Asset Identity and Resolution

An asset proxy provides a stable logical identity for an external or physical object.

Its current location or means of access can change without changing the identity of the proxy.

Examples include:

- A physical document whose storage location changes.
- A digital document whose external location changes.
- A contractor estimate whose original email is difficult to resolve.
- An encrypted source document represented by selected extracted information.

The source may be directly accessible, accessible through a connector, accessible only to the user, or temporarily unavailable.

A useful proxy should remain meaningful even when the original source cannot currently be resolved.

**Evaluation points**

- Provenance can be captured as an ordinary reference in a work-item log or record.
- URLs, file paths, breadcrumbs, and other references can all be useful.
- External sources do not need to become AKM records.
- A digital proxy can represent an external or physical object.
- The proxy's identity should remain stable when its location changes.
- Current location and historical movement can both be useful.
- Extracted information can provide a useful projection without reproducing the complete source.
- Verification should remain possible when the source is accessible.

---

## 6. Payroll Provider Succession

This fixture explores how relationships can represent changing associations without introducing an unnecessary proxy.

### 6.1 Initial Conditions

An Employer record represents the user's employer.

The employer uses Payroll Provider A.

A provider record represents Payroll Provider A and links to the Employer.

### 6.2 Provider Changes

The employer changes payroll providers.

Payroll Provider B replaces Payroll Provider A.

Later, Payroll Provider C replaces Payroll Provider B.

The historical provider records remain available.

Conceptually:

```text
Payroll Provider A ──┐
Payroll Provider B ──┼──> Employer
Payroll Provider C ──┘

Employer ───────────────> Payroll Provider C
```

The provider records link to the employer.

The employer's current outgoing relationship identifies Payroll Provider C.

The graph therefore provides both historical association and current selection.

There is no need to create an asset proxy for the changing payroll service.

### 6.3 Rediscovery

During a future tax preparation project, the user asks where to retrieve the current W-2.

The agent follows the employer's current relationship to Payroll Provider C.

If historical payroll information is needed, the previous provider records remain discoverable through their relationships to the employer.

**Evaluation points**

- Changing service providers can be represented as distinct entities.
- Historical associations remain discoverable.
- Current relationships can distinguish the active provider from historical providers.
- A proxy is not necessary when ordinary entity relationships adequately represent the situation.
- The graph should not require historical records to be overwritten when an association changes.

---

## 7. Agent-Assisted Refinement and Memory Consolidation

This fixture explores the agent responsibilities surrounding the lifecycle of work.

### 7.1 During Work

The user performs a project or work item.

The working agent helps accomplish the task.

As events occur, the agent records observations, references, discoveries, and decisions in the work history.

The user should not have to stop every time something noteworthy happens to determine its final knowledge-management destination.

### 7.2 At a Meaningful Boundary

When a substantial project completes, a separate refinement activity examines the accumulated work.

The refinement agent receives the relevant work history and relationships.

It identifies potential durable changes, such as:

- Updating a hub.
- Creating a resource record.
- Updating a runbook.
- Preserving a historical relationship.
- Creating an asset proxy.
- Extracting a planning record from an earlier event.
- Leaving the information exactly where it is.

The agent proposes changes and asks focused questions.

The user retains authority over what is promoted, changed, or created.

### 7.3 Refinement Depth

Not every work item requires a formal retrospective.

The level of attention should correspond to the significance of the boundary and the discoveries made.

An ordinary task may need no refinement.

A substantial annual project may justify a more deliberate review.

A novel responsibility may require conversational discovery and the creation of initial structures.

An established process may need only targeted updates to existing records.

### 7.4 Independent Review Context

The user prefers the refinement activity to be performed by a different agent from the one that conducted the work.

This provides a fresh perspective, similar to an independent review.

The distinction is between responsibilities, not necessarily a requirement for different models or a particular orchestration mechanism.

**Evaluation points**

- Work history can serve as the input to later knowledge refinement.
- Refinement should be an available lifecycle opportunity.
- The agent should proactively identify candidate changes.
- The user should not need to remember every discovery.
- The agent proposes; the user decides.
- Refinement depth should be proportional to the significance of the work and discoveries.
- Different agent responsibilities can provide useful separation of execution and review.

---

## 8. Record Types, Lifecycle Semantics, and Open Metadata

This fixture explores how the same knowledge environment contains records that should not be forced into a common work-item lifecycle.

### 8.1 Different Kinds of Records

The Homebase environment includes:

- Projects with finite outcomes.
- Tasks with actionable completion.
- Areas representing ongoing responsibilities.
- Resources that preserve enduring reference information.
- Asset proxies representing external or physical objects.
- Entities such as employers and service providers.
- Historical service records and other event records.

These objects have different purposes and temporal characteristics.

A completed project is not equivalent to an inactive resource.

An Area may remain active for years without ever becoming completed.

A physical document proxy may remain useful indefinitely while its current location changes.

A service provider may become inactive because an employer changes providers.

### 8.2 Different Lifecycle Semantics

AKM should support a subset of strongly typed records with defined lifecycle behavior.

Those types can receive lifecycle rules and guarantees appropriate to their purpose.

Not every document needs to participate in a project or task lifecycle.

A resource might have a useful active/inactive distinction, but that does not make it a project.

An asset proxy might also have active/inactive status without being a work item.

The lifecycle should express the meaning of the record, not force every record into a universal Kanban workflow.

### 8.3 User-Defined Dimensions

The user may want to classify resource records by dimensions that AKM does not define.

For example:

```yaml
type: resource
resource_kind: payroll-provider
provider_region: north-america
```

Another resource might use:

```yaml
type: resource
resource_kind: vehicle-document
vehicle: F-150
```

The user-defined dimensions should be available to generic views.

A Resources view might filter by `resource_kind` and display only payroll providers, vehicle documents, or another user-defined category.

The generic view does not need to understand the semantics of every custom classification.

### 8.4 Managed and Unmanaged Document Space

AKM operates on Markdown documents with YAML front matter.

Users and agents can add arbitrary fields and content.

AKM must therefore coexist with information outside its managed schema.

The boundary is one of authority:

- AKM can enforce guarantees for fields and lifecycle semantics it manages.
- User-defined fields can participate in generic filtering, grouping, and presentation when their structure is sufficiently understandable.
- Unknown fields and unrelated document content should remain tolerated.
- Conflicts with managed fields should be conspicuous rather than silently accepted.

This is not unrestricted execution of arbitrary user-defined lifecycle behavior.

It is an open document environment with a deliberately managed semantic subset.

**Evaluation points**

- Strongly typed records can have distinct lifecycle rules.
- Not every document should be forced onto a Kanban board.
- User-defined metadata can support useful generic views.
- Record type, lifecycle state, and presentation are separate concerns.
- AKM must tolerate arbitrary document content outside its managed domain.
- Managed invariants must not be silently overridden by user-defined fields.
- Unknown metadata should not make an otherwise usable document or vault unusable.

---

## 9. Cross-Fixture Evaluation

The fixtures collectively exercise several behaviors that should be considered together.

| Concern                                         | Fixtures |
| ----------------------------------------------- | -------- |
| Persistent responsibilities and finite projects | 2, 7     |
| Memory consolidation and reusable procedures    | 2, 7     |
| Progressive structuring                         | 3        |
| Historical events and derived planning records  | 3, 4     |
| External provenance and source resolution       | 5        |
| Stable identity with changing location          | 5        |
| Relationship-based succession                   | 6        |
| Heterogeneous lifecycle semantics               | 8        |
| Open metadata and generic views                 | 8        |
| Agent-assisted discovery and refinement         | 2, 4, 7  |

### 9.1 Combined Behavioral Expectations

A design that adequately supports Homebase should be able to accommodate the following without requiring every scenario to use the same mechanism:

1. A persistent responsibility can connect multiple finite projects and remain useful between them.
2. A structured collection of event records can be queried without turning the collection's hub into a duplicate of its contents.
3. A source document can remain external while its useful information is represented inside AKM.
4. A digital proxy can preserve the identity of an object while its location changes.
5. An original event record can remain intact while selected information is extracted for a new purpose.
6. Historical relationships can remain discoverable while current relationships identify the present state.
7. Records with different purposes can have different lifecycle semantics.
8. User-defined metadata can support useful generic presentation without requiring AKM to understand every custom classification.
9. Agents can help create and maintain structured information without making the user perform all of the bookkeeping manually.
10. Historical information can become useful again when a new project, event, or question provides the context for retrieval.

### 9.2 Architectural Restraint

These fixtures should not be used to justify a universal ontology, a mandatory provenance subsystem, a special-purpose proxy engine, or a comprehensive user-defined lifecycle framework.

Several scenarios can be handled through ordinary documents, relationships, references, and agent capabilities.

More specialized mechanisms should be justified by the behavior they enable, not by the conceptual elegance of the abstraction.

The governing principle remains:

**The simplest solution that adequately supports the workload wins.**
