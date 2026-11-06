---
id: "261106-xrpr"
title: Rerun the rails trial on other agent harnesses
status: raw
created: 2026-11-06
updated: 2026-11-06
---

# Rerun the rails trial on other agent harnesses

## Goal

Rerun the rails trial (`docs/design/rails-trial-protocol.md`, the method of 261101-0vh7) on Codex and any other agent harness the maintainer cares about, and score it against the protocol's rows. The first Codex run and the Claude rerun found the gaps that vs5c closed; the PR row (R14: open a pull request with the agent line and do not merge it) was reached by neither, and the rules changed since: a unique `--as` name, CM5, `Assisted-by`. Deferred by the maintainer on 2026-11-06: the focus is Claude, and other harnesses wait for a first release. Done when:

- [ ] The protocol is run on Codex (and any other harness named then) from a clean checkout of a current release, with the maintainer's own prompt.
- [ ] Every row is scored, the PR row included, and each deviation is moved to the layer that can prevent it, as in vs5c.

## Log

- 2026-11-06: filed.
