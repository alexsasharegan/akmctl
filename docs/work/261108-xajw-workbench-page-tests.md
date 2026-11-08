---
id: "261108-xajw"
title: workbench page tests
status: active
kind: chore
priority: 3
created: 2026-11-08
updated: 2026-11-08
claimed_by: claude-sonnet-5-5-rest
claimed_at: 2026-11-08T22:35
---

# workbench page tests

## Goal

Finding F4 of [the audit](../design/test-strategy-audit.md), from [261030-456q](done/261030-456q-test-strategy-and-coverage-audit.md): 11 of the 25 workbench modules touch the DOM and have no automated test. Add a DOM test with a small fake `document` for the render functions, or a recorded browser run, and at least a smoke test that `boot.js` renders the board against a fixture repository.

Done when:

- [ ] A smoke test renders the board from `boot.js` against a fixture repository.
- [ ] The render functions of the DOM modules are tested, or the Log records why a recorded browser run was chosen instead.

## Next

Try `boot.js` under Node with a fake `document` and see how far it gets.

## Log

- 2026-11-08: filed.
- 2026-11-08 15:35-07:00: claimed (claude-sonnet-5-5-rest).
