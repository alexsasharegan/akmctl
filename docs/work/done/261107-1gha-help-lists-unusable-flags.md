---
id: "261107-1gha"
title: help lists --fields and --changes flags that refuse a JSON object
summary: new --help and set --help list --fields and --changes, but a JSON object gives 'must be an object'; the CLI takes --field key=value
status: done
kind: bug
created: 2026-11-07
updated: 2026-11-08
closed: 2026-11-08
---

# help lists --fields and --changes flags that refuse a JSON object

## Goal

Command: `akmctl new work x --fields '{"kind":"bug"}'` and `akmctl set <id> --changes '{"status":"ready"}'`. Actual: `invalid_request: fields must be an object` (and `changes must be an object`), although `new --help` and `set --help` list `--fields <text>` and `--changes <text>`. Expected: help names the flag that works (`--field key=value`, repeatable; `set <id> key=value`), or the flags accept a JSON object, and the error says which form to use.

Done when:

- [x] `new --help` and `set --help` list only forms that work, with an example of each
- [x] the error for a JSON object names the working form

## Next

## Outcome

--fields and --changes accept a JSON object, merged with --field and key=value; an invalid value is refused with a message naming the working forms; help shows <json>.

## Log

- 2026-11-07: filed.
- 2026-11-08 00:03-07:00: The cold-agent retest of 261101-t8h8 also hit this: 'set <id> --changes {"priority":2}' returns 'changes must be an object' for an object. Related: 261108-h7ds (bracketed blocked_by lists).
- 2026-11-08 08:50-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08 08:50-07:00: the CLI now reads a JSON object for --fields (new) and --changes (set), merged with --field and key=value; bad values name the working forms; help shows <json>; TestJSONObjectFlags; changelog.
