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

- [x] A smoke test renders the board from `boot.js` against a fixture repository.
- [ ] The render functions of the DOM modules are tested, or the Log records why a recorded browser run was chosen instead.

## Next

Cover the rest of the DOM modules with the fake document: ui.js (palette, help overlay, toasts), flow.js (moves, capture), wizard.js, anim.js, picker.js, controls.js, and the drawer's sections (extend app.test.mjs's fake API with the section route). Add any DOM member a module needs to fakedom.mjs. Then tick box 2.

## Log

- 2026-11-08: filed.
- 2026-11-08 15:35-07:00: claimed (claude-sonnet-5-5-rest).
- 2026-11-08 15:38-07:00: ticked "A smoke test renders the board from `boot.js` against a fixture repository.": TestDOMModulesUnderNode: testdata/app.test.mjs runs app.js against a fake document and a fake API: token stripped, records fetched, four lanes drawn, title and nav, route change to ADRs, record drawer
- 2026-11-08 15:38-07:00: Added testdata/fakedom.mjs (a small fake document that throws on a member it lacks), dom.test.mjs (boot theme, board, ADRs, check, tables) and app.test.mjs (whole page); checked by breaking board.js and app.js. Not yet covered: ui.js, flow.js, wizard.js, anim.js, capture.js, picker.js, controls.js, and the drawer's sections (detail.js draws fields only in the fake; the sections need the section API). Box 2 stays open.
