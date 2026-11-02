---
id: "261101-jgdy"
title: log and section guard against shell split text
status: ready
created: 2026-11-01
updated: 2026-11-01
kind: bug
parent: 261101-vs5c
summary: An unquoted log or section message splits at a shell metacharacter; the tool cannot tell and the retry duplicates the line
---

# log and section guard against shell split text

## Goal

Seen twice in the first Codex trial: `akmctl log <id> words; more words` ran as two commands. `log` appended the text before the `;`, zsh printed `command not found`, and the retry appended the full sentence, leaving a partial line and a duplicate in the Log. Done when:

- [ ] Decide what the tool can safely do: warn when `log` gets several bare words (join and note it), or print the joined text back so the caller sees a truncated message; record the choice (an ADR if it changes the API contract).
- [ ] The chosen behaviour is implemented, tested and documented in the `log` and `section` references of both skill copies, with the rule "quote the whole message as one argument" in the `akm` skill.

## Log

- 2026-11-01: filed.
