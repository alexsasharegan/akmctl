---
id: "261119-5680"
title: radon upright not italic
status: active
created: 2026-11-19
updated: 2026-11-19
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-20T03:20
---

# radon upright not italic

## Goal

The workbench sets Monaspace Radon, the accent face, in italics in two places: the card detail line and the "Raw" word in the capture modal heading. Radon looks good upright, so both become normal style. The top-left wordmark is unchanged (it is already upright and the maintainer likes it). Done when no rule in `app.css` sets italic together with the Radon face.

## Next

Remove italic from the two Radon rules in app.css, add the changelog entry, run the checks, open the pull request.

## Log

- 2026-11-19: filed.
- 2026-11-19 20:20-07:00: claimed (claude-sonnet-5-5).
- 2026-11-19: removed italic from `.card-detail` and `.raw-word` in `app.css`; the wordmark was already upright. Tests with `-race`, records check and format check pass. Not checked by eye in the browser.
