---
id: "261101-xyfy"
title: spike version flag instead of version command
status: raw
kind: spike
summary: Decide whether -v and --version should print the version as well as the version command
created: 2026-11-01
updated: 2026-11-01
---

# spike version flag instead of version command

## Goal

A spike, probably ending in "do not": Alexander keeps reaching for `-v` or `--version` instead of `version`. Today `--version` already works as an alias in `internal/cli/cli.go`; `-v` does not (it fails as an unknown command). Decide whether to add `-v`, and whether the bare `version` command stays the documented form, weighing: `-v` is often `--verbose` elsewhere, the spec and API name the operation `version`, and a second spelling is one more thing to document and test.

Done when:

- [ ] the Outcome records the decision and the reason, and any follow-up item is filed

## Log

- 2026-11-01: filed.
