---
id: "261101-w82x"
title: close leaves blank line after next
status: raw
created: 2026-11-01
updated: 2026-11-01
parent: 261101-vs5c
---

# close leaves blank line after next

## Goal

Since close empties a non-empty ## Next (261101-avg8), the closed file keeps an extra blank line after the empty heading, and Prettier rejects it. Codex hit this on 261101-3erc and the Claude session on 261101-0vh7, each fixing it by hand with a formatting-only commit. Done when:

- [ ] close leaves the file Prettier-clean when it empties Next, with a test.
- [ ] A CHANGELOG entry notes the fix if the output changes.

## Log

- 2026-11-01: filed.
