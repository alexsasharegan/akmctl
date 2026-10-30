---
id: "261029-kmre"
title: "Workbench toasts: placement, levels, persistence and a countdown border"
summary: "Rework the workbench toasts: bottom right, four levels, ephemeral and persistent flavours, a countdown border, hover pause and a close button"
status: ready
kind: feature
priority: 1
created: 2026-10-29
updated: 2026-10-29
---

# Workbench toasts: placement, levels, persistence and a countdown border

## Goal

1. toasts should arrive in the bottom right of the window
2. toasts should come in two flavors: ephemeral & persistent. persistent toasts must be cleared by the user (should have a keyboard binding). ephemeral should clear after a longer time frame. suggest 8 seconds. persistent toasts may be extended to support actions (buttons which take an action like undoing something)
3. toasts should come in level flavors. a suggested set would be: success, info, warning, error.
4. toasts should have a thick top border line which doubles as the countdown animation letting the user know when the toast will be automatically dismissed
5. if a user hovers over the toast, the countdown is paused. animation should reflect.
6. toasts should have a close icon for individual dismissal

Refinement (2026-10-29): the styling freeze is lifted for this item; use the existing design tokens. The first cut has no action buttons; undo actions are a follow-up once the toast shell exists.

Proposals to confirm when picked up:

- level and flavor are independent; errors are persistent by default and the other levels ephemeral
- the existing `toast(msg, "bad")` callers map to the error level and the rest to info; the API becomes `toast(msg, { level, persistent, ms })`, and the callers are migrated
- placement sits above the key hint bar and never covers a dialog's buttons (the existing code clears toasts when a dialog opens)
- Escape dismisses the newest persistent toast when no overlay is open
- hover and keyboard focus both pause the countdown; with reduced motion the border is static and the timeout still runs
- warning and error toasts use role alert, the others role status, so screen readers announce them; the countdown is never the only signal
- a stacking limit stays (three today); persistent toasts are not evicted by ephemeral ones

Done when:

- [ ] toasts appear bottom right, above the hint bar
- [ ] four levels (success, info, warning, error) look distinct without relying on colour alone
- [ ] ephemeral toasts clear after 8 seconds; persistent ones stay until dismissed by the close icon or the keyboard binding
- [ ] the top border counts down, pauses on hover and focus, and respects reduced motion
- [ ] every existing toast caller still shows its message at a sensible level
- [ ] unit tests cover the queue (limit, persistence, timers); checked in the browser in light and dark
- [ ] the Toast row in docs/glossary.md lists the levels and flavors

## Next

Ready; no dependencies. Start here: it has the highest priority.

## Log

- 2026-10-29: filed.
- 2026-10-29: refined with Alexander (first cut without actions).
