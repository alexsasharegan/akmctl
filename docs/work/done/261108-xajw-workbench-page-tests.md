---
id: "261108-xajw"
title: workbench page tests
status: done
kind: chore
priority: 3
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
---

# workbench page tests

## Goal

Finding F4 of [the audit](../../design/test-strategy-audit.md), from [261030-456q](261030-456q-test-strategy-and-coverage-audit.md): 11 of the 25 workbench modules touch the DOM and have no automated test. Add a DOM test with a small fake `document` for the render functions, or a recorded browser run, and at least a smoke test that `boot.js` renders the board against a fixture repository.

Done when:

- [x] A smoke test renders the board from `boot.js` against a fixture repository.
- [x] The render functions of the DOM modules are tested, or the Log records why a recorded browser run was chosen instead.

## Next

## Outcome

The workbench page has DOM tests that run in the Go suite under Node with no dependency. A fake document (testdata/fakedom.mjs) and fake API (apphelper.mjs) drive app.js, the views, the drawer and editors, the palette, the move wizard, capture, drag and drop, the project filter and the animator. 524 of 696 functions run. Untested: live polling in the page, the away list, confetti and the animation themes.

## Log

- 2026-11-08: filed.
- 2026-11-08 15:35-07:00: claimed (claude-sonnet-5-5-rest).
- 2026-11-08 15:38-07:00: ticked "A smoke test renders the board from `boot.js` against a fixture repository.": TestDOMModulesUnderNode: testdata/app.test.mjs runs app.js against a fake document and a fake API: token stripped, records fetched, four lanes drawn, title and nav, route change to ADRs, record drawer
- 2026-11-08 15:38-07:00: Added testdata/fakedom.mjs (a small fake document that throws on a member it lacks), dom.test.mjs (boot theme, board, ADRs, check, tables) and app.test.mjs (whole page); checked by breaking board.js and app.js. Not yet covered: ui.js, flow.js, wizard.js, anim.js, capture.js, picker.js, controls.js, and the drawer's sections (detail.js draws fields only in the fake; the sections need the section API). Box 2 stays open.
- 2026-11-08 16:23-07:00: ticked "The render functions of the DOM modules are tested, or the Log records why a recorded browser run was chosen instead.": TestDOMModulesUnderNode: dom, app, keys, drawer, wizard and project tests; 524 of 696 functions in web/ run (v8 coverage); not run: live poll, the away list, confetti and the animation themes
- 2026-11-08 16:23-07:00: Box 2: added keys, drawer, wizard and project tests on the fake DOM (palette, help, theme, safe move and undo, wizard steps and failures, capture, drag and drop, field editors, claim bar, project filter, animator). Function coverage by v8: 524 of 696; the rest is mostly pure code that logic.test.mjs checks, plus live polling, the away list and the animation themes. Mutating anim.js, flow.js and ui.js each failed a test; mutating wizard.js's written counter did not (display only).
