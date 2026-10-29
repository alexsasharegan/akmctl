---
id: "261028-nkzy"
title: Click through the workbench actions by hand
summary: Claim, Release, Close and drag-to-blocked in the workbench were never exercised by hand in a browser; do it and fix what breaks.
status: done
kind: chore
priority: 1
created: 2026-10-28
updated: 2026-10-29
---

# Click through the workbench actions by hand

## Goal

The workbench's Claim, Release, Close and drag-to-blocked actions were built and tested through the API but were never clicked by hand in a browser.

Done when:

- [x] Each action is run by hand in a browser against a scratch copy of the sample repository, and the result on disk is checked.
- [x] Anything that fails is fixed or filed as its own item.
- [x] The Outcome says which browser and which actions were exercised.

## Outcome
Claim, Release, Close and drag onto a lane were clicked by hand in the Claude desktop built-in browser against a scratch copy of this repository; each write was checked on disk. One finding (261029-7zcp, a duplicate claim log line) is filed.

## Log

- 2026-10-28: filed.
- 2026-10-29: driven on a scratch copy of this repository with the built-in browser of the Claude desktop app: Claim (wizard and drawer button), Release (drawer button), Close (key d, wizard with an ADR link) and drag of a card onto the Blocked and the Active lane. File contents checked after each write. One finding, filed as 261029-7zcp.
