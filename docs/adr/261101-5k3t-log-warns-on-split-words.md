---
id: "261101-5k3t"
title: log warns when the shell split the message, and still writes
summary: The log command warns on stderr when it joins several words, showing what it wrote; it does not refuse or change the API
status: proposed
date: 2026-11-01
---

# log warns on split words

## Context

In the first Codex trial an unquoted `akmctl log <id> a b; c d` ran as two commands. `log` wrote the words before the `;`, the shell failed on the rest, and the retry duplicated the line. It happened twice. The tool cannot tell a truncated message from a short one, but it can see that the shell passed several words, which a quoted message never does. The `log` reference itself showed an unquoted message.

## Considered Options

- Warn on stderr when several words are joined, show the text written, and still write (CLI binding only).
- Refuse several words: a quoted message is the only accepted form.
- Add a `warnings` member to the API response, and warn when the new line extends the previous Log line (a probable retry).
- Do nothing and rely on the instruction to quote.

## Decision

Warn and still write, in the CLI binding only. Refusing breaks scripts and the reference's own earlier example, for a failure the warning already makes visible: it shows the exact line, so a cut-short message is plain at the moment it is written. The API gains nothing, so a client is unaffected. The retry-detection warning is left for later; it needs a response member and a rule for what counts as a retry.

## Consequences

- A caller that passes unquoted words sees one stderr line with the text written and the advice to quote. Output on stdout and the exit code are unchanged.
- The `log` reference and the `akm` skill now show a quoted message. `section` already refuses extra words, so it needs no change; a `--text` cut by the shell is still silent.

## Confirmation

`TestLogWarnsWhenTheShellSplitTheMessage` in `internal/cli`.

## Review

Proposed 2026-11-01 under the agent-rails project 261101-vs5c. Accepted by merging the branch that implements it, if Alexander agrees.

## Revert

Remove the warning in `internal/cli/cli.go`, its test, and the quoting sentence in the two `log.md` copies.
