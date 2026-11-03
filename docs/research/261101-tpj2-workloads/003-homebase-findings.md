---
id: "261102-9v6d"
title: Homebase — Persistent Responsibilities and Cyclical Knowledge
summary: Defines an independent workload for evaluating persistent areas, finite projects, enduring resources, and contextual rediscovery across recurring personal responsibilities.
status: draft
kind: research
created: 2026-11-02
updated: 2026-11-02
parent: "261101-tpj2"
---

# Homebase — Persistent Responsibilities and Cyclical Knowledge

## 1. Scenario

Homebase is a personal knowledge and life-management system maintained over several years.

It serves as a central repository for information about the user's responsibilities, interests, projects, learning, and personal administration.

Unlike a conventional task manager, Homebase contains both actionable work and enduring knowledge.

Some information receives frequent attention. Other information remains dormant for months or years before becoming useful again.

Activity also follows recognizable annual and seasonal cycles.

The user organizes information using Tiago Forte's PARA method:

- Projects
- Areas
- Resources
- Archives

This organizational model provides a useful starting point for exploring how persistent responsibilities, finite projects, and enduring reference material relate to one another.

The investigation should determine which aspects of this organizational model are useful to AKM and which require no special architectural treatment.

## 2. Organizational Environment

### 2.1 Projects

Projects represent finite efforts with identifiable outcomes.

Examples include:

- Preparing an annual tax return.
- Replacing vehicle tires.
- Completing a home improvement.
- Planning a vacation.
- Organizing a major purchase.
- Completing a learning objective.

Projects naturally contain related work items.

For example:

```text
2026 Tax Preparation (Project)
├── Collect W-2s
├── Gather expense records
├── Review deductions
└── File return
```

Individual work items progress through their own lifecycles. Once the work is complete, the parent project can also be completed and archived.

Projects are therefore natural candidates for ordinary work-item management.

### 2.2 Areas

Areas represent ongoing responsibilities that have no predetermined completion date.

Examples include:

- Taxes.
- Vehicle maintenance.
- Home maintenance.
- Personal finances.
- Health and fitness administration.
- Professional development.
- Household management.

An area persists because the responsibility persists.

For example, taxes are an ongoing responsibility, but preparing the 2026 tax return is a finite project.

Similarly, vehicle maintenance is an ongoing responsibility, while replacing the tires is a finite project.

An area can therefore be understood conceptually as a project that never closes.

One possible organizational structure is:

```text
Taxes (Area)
└── 2026 Tax Preparation (Project)
    ├── Collect W-2s
    ├── Gather expense records
    ├── Review deductions
    └── File return
```

The project completes and moves into the archives. The Taxes area remains.

This raises an important question: does an area require a distinct entity or lifecycle, or is a persistent organizational record that anchors related projects sufficient?

The investigation should explore this without presupposing the answer.

### 2.3 Resources

Resources are enduring reference materials that do not necessarily represent work or responsibility.

Examples include:

- Learning notes.
- Technical references.
- Book notes.
- Research and discoveries.
- Equipment documentation.
- Household information.
- Personal reference material.
- Architecture and design documents.

A resource may remain useful for years without participating in a conventional work-item lifecycle.

For example, an appliance's model number, purchase information, warranty, and repair history may be valuable whenever the appliance requires maintenance.

Similarly, a technical learning note may become relevant to a project long after it was originally created.

Resources are not necessarily organized around completion or archival.

The investigation should explore how they relate to projects and areas, and whether they need any lifecycle management beyond ordinary discoverability and maintenance.

### 2.4 Archives

Archives contain completed projects and inactive historical material.

Examples include:

- Previous years' tax preparation.
- Completed home improvement projects.
- Past travel planning.
- Abandoned personal projects.
- Historical learning notes.
- Superseded reference documents.

Archived information is not necessarily obsolete.

A completed project may provide useful context for future work. An old checklist may be reused in the following year. A previously abandoned project may become interesting again.

Archival should therefore be distinguished from deletion and permanent irrelevance.

## 3. Annual and Seasonal Activity

Homebase exhibits recurring patterns of attention throughout the year.

| Period            | Representative activity                         |
| ----------------- | ----------------------------------------------- |
| January           | Annual review, resolutions, and planning        |
| February          | Regular projects and learning                   |
| March–April       | Tax preparation                                 |
| May–July          | Home projects and maintenance                   |
| August            | Back-to-school preparation                      |
| September–October | Regular activities and ongoing responsibilities |
| November          | Holiday planning and gift research              |
| December          | Year-end review and preparation                 |

These are illustrative patterns rather than rigid scheduling requirements.

Individual responsibilities and projects may follow different cycles.

The important characteristic is that information can remain dormant for most of the year and still be highly valuable during a particular period.

## 4. Information Produced

Homebase naturally contains several populations of information:

- Personal projects and their work items.
- Ongoing areas of responsibility.
- Learning notes and discoveries.
- Research and external references.
- Household and equipment records.
- Recurring maintenance information.
- Annual planning and review records.
- Administrative records.
- Historical project documentation.
- Reusable checklists and procedures.

Some information changes frequently. Other information is relatively stable.

Some records are important because they describe current responsibilities. Others are important because they preserve the history of previous work.

The investigation should explore these distinctions without assuming that every population requires a specialized record type.

## 5. Scenarios to Explore

### 5.1 Completing One Year's Tax Preparation

The user completes the 2026 tax preparation project.

The project contains multiple work items, supporting information, and references.

Once the return is filed, the project is completed and archived.

The Taxes area remains active, and the following year's tax preparation will eventually become a new project.

Explore:

- How the persistent area relates to each annual project.
- Whether previous projects remain useful as historical references.
- Which supporting records should remain associated with the completed project.
- Whether reusable checklists or discoveries should be carried forward.
- How the following year's project can benefit from previous experience.

### 5.2 Returning to Vehicle Maintenance

The user has a vehicle maintenance area containing information about the vehicle, its service history, and previous maintenance projects.

Several months after the last maintenance activity, the user needs to replace the tires.

The previous work may contain useful information about the vehicle, tire specifications, preferred service providers, and maintenance history.

Explore how the system should surface relevant information from the persistent area and previous projects.

### 5.3 Rediscovering a Gift Idea

The user captures a gift idea in February.

The note is not relevant to any active project at the time.

In November, the user begins planning holiday gifts and needs to recover the idea.

Explore how information that was captured casually and remained dormant for months can become discoverable in the appropriate context.

### 5.4 Reusing an Annual Checklist

The user completed a recurring annual activity the previous year using a checklist.

The following year, the user needs to perform the same activity again.

The previous execution record is useful, but the user may also want to reuse the checklist as a starting point.

Explore the distinction between:

- The reusable procedure.
- The previous year's execution.
- The new year's project.
- The historical evidence of what happened previously.

Determine whether these distinctions require specialized entities or can be handled through simpler organization and relationships.

### 5.5 Returning to a Dormant Personal Project

The user begins a personal project but loses interest before completing it.

A year later, circumstances change and the project becomes interesting again.

The user needs to recover the original goals, previous work, research, decisions, and unresolved questions.

Explore how an archived or inactive project can be reactivated without losing its historical context.

### 5.6 Reusing an Old Learning Note

The user captures a technical concept while studying an unrelated subject.

Two years later, the concept becomes relevant to a new project.

The original note was never associated with the new project and may not have been accessed since its creation.

Explore how durable knowledge can become relevant again through changing interests and project requirements.

### 5.7 Maintaining a Durable Resource

The user maintains a reference document containing information about a household appliance.

The document is updated when the appliance is repaired, serviced, or otherwise changes.

The resource persists across multiple maintenance projects and years.

Explore whether this kind of document needs any managed lifecycle at all, and how it can remain discoverable without participating in ordinary task completion and archival workflows.

### 5.8 Revisiting an Archived Project

A previously completed project becomes relevant because circumstances change.

The user needs information from the original project, but does not necessarily intend to resume the original work.

Explore how historical information can be reused without requiring the project to be reactivated.

## 6. Temporal Characteristics

Homebase contains several overlapping timescales:

| Timescale  | Representative activity                                           |
| ---------- | ----------------------------------------------------------------- |
| Daily      | Capturing notes, managing tasks, and routine administration       |
| Weekly     | Ongoing projects and personal responsibilities                    |
| Monthly    | Maintenance, reviews, and recurring activities                    |
| Seasonal   | Holidays, school preparation, and other periodic activities       |
| Annual     | Taxes, resolutions, annual planning, and recurring reviews        |
| Multi-year | Long-lived resources, historical knowledge, and personal projects |

The significance of an individual record can vary substantially across these timescales.

A note may be unimportant for months and then become immediately useful. A resource may remain continuously valuable without being accessed frequently. A completed project may be inactive but provide important historical context for a future project.

The investigation should avoid equating inactivity, low access frequency, and obsolescence.

## 7. Research Objectives

The investigation should establish what AKM needs to accomplish in this environment.

Pay particular attention to:

1. **Persistent responsibilities:** How should ongoing areas relate to finite projects?
2. **Project completion:** What happens when a project and its child work items are completed?
3. **Recurring activity:** How should successive executions of the same responsibility relate to one another?
4. **Enduring resources:** How should reference material exist independently of ordinary work-item lifecycles?
5. **Archival:** What does archival mean when historical projects and information remain useful?
6. **Seasonal rediscovery:** How can dormant information become discoverable at the appropriate time?
7. **Contextual relevance:** How should information become useful again when circumstances or interests change?
8. **Knowledge reuse:** How can previous projects, checklists, and discoveries contribute to future work?
9. **Organizational simplicity:** Can persistent anchors and ordinary relationships adequately represent areas, projects, resources, and archives?
10. **Lifecycle simplicity:** Which information genuinely needs lifecycle management, and which can remain as enduring reference material?

## 8. Boundaries and Assumptions

- Homebase is a personal knowledge and life-management system.
- It is maintained over several years.
- It contains projects, recurring responsibilities, reference material, learning notes, and historical information.
- PARA is the user's established organizational model.
- Projects represent finite outcomes.
- Areas represent persistent responsibilities.
- Resources represent enduring reference material.
- Archives preserve completed and inactive information.
- Activity follows heterogeneous weekly, monthly, seasonal, and annual patterns.
- Historical information remains available and may become useful again.

These assumptions establish the experimental environment. They do not prescribe how AKM should organize or manage its information.

## 9. Central Question

**How can a knowledge-management system support persistent responsibilities, finite projects, and enduring reference material while making historically accumulated information useful again across recurring activities and changing contexts?**

The investigation should identify the simplest mechanisms that adequately support this outcome.
