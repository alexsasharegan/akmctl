---
id: "261109-4wtb"
title: Reconcile the roadmap after the theming and Windows work
status: done
kind: chore
priority: 1
created: 2026-11-09
updated: 2026-11-09
closed: 2026-11-09
---

# Reconcile the roadmap after the theming and Windows work

## Goal

Bring `docs/ROADMAP.md` in line with the board after the theming stack (PRs 67, 70 to 74) and the Windows and test-hardening work landed, and sort the 1.0.0 items by how much of the maintainer's input each needs, so an unattended agent can take the ones that need none.

Done when:

- [x] `Now` names what is in flight or waiting on the maintainer; `Next` is the queue of items an agent can take without a decision, in order.
- [x] The 1.0.0 items are grouped by exit criterion, each marked with what decision it needs (none, a recorded pick, or the maintainer), and criterion 5 has an item.
- [x] Nothing `done` or `dropped` is listed; every `ready` or `active` top-level item is listed once (RM1 to RM4 pass).
- [x] `akmctl check` reports no new finding.

## Next

## Outcome

docs/ROADMAP.md now matches the board: Now holds the Windows thread, Next is the unattended queue (eight items with no decision, four 1.0.0 items with a pick), Later holds the 1.0.0 criteria and what needs the maintainer. Criterion 5 got its item (261109-vxdr). Not covered: the raw items are untriaged, and 261101-12c0 still warns WK12 (title and heading differ).

## Log

- 2026-11-09: filed.
- 2026-11-09 01:22-07:00: claimed (claude-sonnet-5-5-roadmap).
- 2026-11-09 01:23-07:00: Theming stack merged (PRs 67, 70 to 74); the stack's tip needed main merged in and the DOM tests of 261108-xajw moved to the theme model (data-mode, akm-mode). Spelling pass (PR 69) is green and left open for the maintainer. Rewrote the roadmap: Now is the Windows thread, Next is the unattended queue split into no-decision and pick items, Later holds what needs the maintainer. Filed 261109-vxdr (criterion 5) and 261109-e540 (CI 504s).
