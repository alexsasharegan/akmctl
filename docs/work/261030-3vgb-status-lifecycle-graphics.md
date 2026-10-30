---
id: "261030-3vgb"
title: status lifecycle graphics
summary: Prototype an inline transition and a lifecycle strip under the move toast's sentence, then pick one
status: ready
kind: spike
created: 2026-10-30
updated: 2026-10-30
priority: 0
---

# status lifecycle graphics

## Goal

A spike: prototype two ways to show a status move as a graphic in the workbench's toasts, compare them in the browser, and record which to keep. The graphic is a scannable addition **below** the toast's sentence; it never replaces the sentence.

### The two prototypes

1. **Inline transition:** one row under the sentence with the two statuses as the existing `.stat` markers, `[◌ raw] → [◐ ready]`.
2. **Lifecycle strip:** a row of raw, ready, active, blocked and done markers with the from and to markers highlighted and an arc between them. A skipped stage shows as skipped, and a close ends at done or dropped.

Both reuse the status shapes and `--st-*` colours in `app.css`. Optional on either: the marker fills in over about 300 ms (a half-filled marker becomes full for ready to active), off under `prefers-reduced-motion`; a backward move (active to ready) uses a reversed or muted arrow.

### Reach

The graphic appears wherever a toast reports a status move, through one helper, so a move looks the same wherever it came from:

- your own moves (`doneText` and `undoToast` in `js/flow.js`, including the board's keyboard moves)
- the undo toast
- moves made outside the workbench (`describeEvent` in `js/live.js`)

### Constraints

- The toast keeps its full sentence (it has `role=status`); the graphic is `aria-hidden` and carries no meaning that the sentence lacks. The marker shapes already differ without colour.
- A pure function in `js/logic.js` (for example `transition(from, to)`) builds the graphic spec, so it is unit tested; `toastSpec` in `js/toasts.js` accepts the extra content.
- The Undo button keeps working in the same toast.

Done when:

- [ ] both prototypes render under the sentence in the move toast, selectable while comparing (a temporary switch is fine)
- [ ] the same helper serves own moves, the undo toast and outside moves, and the sentence is unchanged in each
- [ ] unit tests cover the transition spec (forward, backward, skipped stages, close as done and dropped) and `toastSpec` with the extra content
- [ ] both checked in the browser, light and dark, with reduced motion on and off
- [ ] the Outcome records which prototype to keep, with a screenshot or the reason, and a follow-up feature item is filed for the chosen one; the losing prototype's code is removed
- [ ] new terms are registered in `docs/glossary.md`

## Log

- 2026-10-30: filed.
- 2026-10-30: refined with Alexander: spike on two prototypes (inline transition, lifecycle strip), graphic below the sentence and never replacing it; marked ready.
