---
id: "261029-kmre"
title: better toast msgs
summary: address toast placement, persistent alerts, and styling
status: raw
kind: feature
priority: 1
created: 2026-10-29
updated: 2026-10-29
---

# better toast msgs

## Goal

1. toasts should arrive in the bottom right of the window
2. toasts should come in two flavors: ephemeral & persistent. persistent toasts must be cleared by the user (should have a keyboard binding). ephemeral should clear after a longer time frame. suggest 8 seconds. persistent toasts may be extended to support actions (buttons which take an action like undoing something)
3. toasts should come in level flavors. a suggested set would be: success, info, warning, error.
4. toasts should have a thick top border line which doubles as the countdown animation letting the user know when the toast will be automatically dismissed
5. if a user hovers over the toast, the countdown is paused. animation should reflect.
6. toasts should have a close icon for individual dismissal

## Log

- 2026-10-29: filed.
