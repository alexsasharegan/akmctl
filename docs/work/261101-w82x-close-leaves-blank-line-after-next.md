---
id: "261101-w82x"
title: close leaves blank line after next
status: active
created: 2026-11-01
updated: 2026-11-06
parent: 261101-vs5c
priority: 2
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-06T18:44
---

# close leaves blank line after next

## Goal

Since close empties a non-empty ## Next (261101-avg8), the closed file keeps an extra blank line after the empty heading, and Prettier rejects it. Codex hit this on 261101-3erc and the Claude session on 261101-0vh7, each fixing it by hand with a formatting-only commit. Done when:

- [x] close leaves the file Prettier-clean when it empties Next, with a test.
- [x] A CHANGELOG entry notes the fix if the output changes.

## Next

Open the pull request; the maintainer merges.

## Log

- 2026-11-01: filed.
- 2026-11-06 11:44-07:00: claimed (claude-sonnet-5-5).
- 2026-11-06 11:46-07:00: Fixed ReplaceSection when clearing: drops the blank lines after the text if a blank line follows the heading; tests for the three shapes (before a section, last section, CRLF).
