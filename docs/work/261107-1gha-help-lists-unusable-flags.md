---
id: "261107-1gha"
title: help lists --fields and --changes flags that refuse a JSON object
summary: new --help and set --help list --fields and --changes, but a JSON object gives 'must be an object'; the CLI takes --field key=value
status: raw
kind: bug
created: 2026-11-07
updated: 2026-11-08
---

# help lists --fields and --changes flags that refuse a JSON object

## Goal

Command: `akmctl new work x --fields '{"kind":"bug"}'` and `akmctl set <id> --changes '{"status":"ready"}'`. Actual: `invalid_request: fields must be an object` (and `changes must be an object`), although `new --help` and `set --help` list `--fields <text>` and `--changes <text>`. Expected: help names the flag that works (`--field key=value`, repeatable; `set <id> key=value`), or the flags accept a JSON object, and the error says which form to use.

Done when:

- [ ] `new --help` and `set --help` list only forms that work, with an example of each
- [ ] the error for a JSON object names the working form

## Log

- 2026-11-07: filed.
- 2026-11-08 00:03-07:00: The cold-agent retest of 261101-t8h8 also hit this: 'set <id> --changes {"priority":2}' returns 'changes must be an object' for an object. Related: 261108-h7ds (bracketed blocked_by lists).
