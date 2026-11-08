---
id: "261101-jgdy"
title: log and section guard against shell split text
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
kind: bug
parent: 261101-vs5c
summary: An unquoted log or section message splits at a shell metacharacter; the tool cannot tell and the retry duplicates the line
decisions: [261101-5k3t]
---

# log and section guard against shell split text

## Goal

Seen twice in the first Codex trial: `akmctl log <id> words; more words` ran as two commands. `log` appended the text before the `;`, zsh printed `command not found`, and the retry appended the full sentence, leaving a partial line and a duplicate in the Log. Done when:

- [x] Decide what the tool can safely do: warn when `log` gets several bare words (join and note it), or print the joined text back so the caller sees a truncated message; record the choice (an ADR if it changes the API contract).
- [x] The chosen behavior is implemented, tested and documented in the `log` and `section` references of both skill copies, with the rule "quote the whole message as one argument" in the `akm` skill.

## Next

Done.

## Outcome

`log` warns on stderr when several words were joined and shows the written line; it still writes, and the API is unchanged (ADR 261101-5k3t, proposed). `section` already refuses extra words. Not covered: a `--text` cut by the shell, and a retry that extends the previous Log line. Decision: 261101-5k3t.

## Log

- 2026-11-01: filed.
- 2026-11-01 17:33-07:00: claimed (claude).
- 2026-11-01 17:34-07:00: ticked "Decide what the tool can safely do: warn when `log` gets several bare words (join and note it), or print the joined text back so the caller sees a truncated message; record the choice (an ADR if it changes the API contract).": ADR 261101-5k3t: warn and still write, CLI only; refusing and an API warnings member are the rejected alternatives.
- 2026-11-01 17:34-07:00: ticked "The chosen behavior is implemented, tested and documented in the `log` and `section` references of both skill copies, with the rule "quote the whole message as one argument" in the `akm` skill.": internal/cli/cli.go, TestLogWarnsWhenTheShellSplitTheMessage, log.md in both skill copies, CHANGELOG.md.
