---
id: "261029-6cp5"
title: The roadmap follows Now / Next / Later, with milestones inside Later
summary: The roadmap follows Now / Next / Later, with milestones inside Later
status: accepted
date: 2026-10-29
---

# The roadmap follows Now / Next / Later, with milestones inside Later

## Context

The roadmap has used three horizons (Now, Next, Later) and a Dropped section since the first release, but the choice was never written down or named. Rules RM1 to RM4 check it, so the convention is already load-bearing. Two questions came up after the 0.2.0 release emptied Now: what happens to finished work, and what "1.0.0" means when it is only a list of items in Later.

## Considered Options

- Dated or quarterly roadmaps. Rejected: dates are promises the tool cannot keep, and a repository's roadmap is edited by whoever and whatever works in it.
- Outcome-based roadmaps, where each slot is a goal and work items sit underneath. Rejected for now: the work item already carries a goal, and a second layer of goals needs a grouping record we do not have.
- Milestone-only roadmaps (a list per release). Rejected as the whole structure: it says nothing about order inside a release, and a single-maintainer project has no use for the Now tier's commitment.
- Now / Next / Later (Janna Bastow's format), with milestones as groups inside Later.

## Decision

The roadmap follows **Now / Next / Later**, plus **Dropped**.

- **Now** holds work committed and under way. **Next** holds work that is probable and ordered. **Later** holds bets. Order within a section is priority order.
- **Milestones** group items inside Later (and Next) as `###` headings, each with exit criteria stated in plain words under it. A milestone says what must be true to cut that version, and an item belongs to it only if it serves one of those criteria. Items that serve no milestone are listed under a heading such as "Workbench and polish" and are not promises.
- Only items with a goal and a done-when belong on the roadmap. A `raw` idea stays a work item until it is refined.
- **Finished work leaves the roadmap.** Its record is the closed work item in `done/` (Outcome and date) and, for released work, the changelog entry (Keep a Changelog). The roadmap never has a Done section. A dropped item moves to Dropped with one line saying why.
- The tool does not move lines on the roadmap; a human or agent edits it, and rules RM1 to RM4 check the result.

`AGENTS.md` says that the roadmap follows this convention and points here.

## Consequences

- A newcomer can look up the format by name, and the milestone definitions give "1.0.0" a testable meaning instead of a list.
- The 1.0.0 milestone and its exit criteria are written into `docs/ROADMAP.md` alongside this ADR.
- No tool or schema change. Compaction of closed work items is a separate decision and is not covered here.
