---
id: "261108-h7ds"
title: set blocked by brackets store a mangled list
status: done
kind: bug
priority: 3
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
---

# set blocked by brackets store a mangled list

## Goal

Found by the cold-agent retest of 261101-t8h8. `akmctl set <id> 'blocked_by=[a,b]'` and `blocked_by=["a","b"]` store a list of mangled strings (`["[a", "b]"]`) and report `changed`; `check` says nothing. Only `blocked_by=a,b` stores a clean list, and neither `set --help` nor the skill shows that syntax. Expected: a bracketed or JSON list is read as a list (or refused with the form to use), and the help or the `set` reference shows `blocked_by=a,b`. Done when:

- [x] `set blocked_by=[a,b]` and `blocked_by=["a","b"]` store `[a, b]`, or are refused with `invalid_request` naming `blocked_by=a,b`.
- [x] A test covers the bracketed, JSON and comma forms; `set --help` or the `set` reference shows the list form.

## Next

## Outcome

The command line reads a JSON array or a bracketed list as a list, so blocked_by=[a,b] and blocked_by=["a","b"] store [a, b]. A wikilink entry stays one entry. set --help and the set reference show the comma form. The API took arrays all along; the fault was only in the CLI.

## Log

- 2026-11-08: filed.
- 2026-11-08 00:11-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08 00:11-07:00: ticked "`set blocked_by=[a,b]` and `blocked_by=["a","b"]` store `[a, b]`, or are refused with `invalid_request` naming `blocked_by=a,b`.": listEntries in internal/cli/cli.go; TestListFieldValueForms, TestSetBlockedByStoresACleanList
- 2026-11-08 00:11-07:00: ticked "A test covers the bracketed, JSON and comma forms; `set --help` or the `set` reference shows the list form.": both tests; set --help and the set reference (both skill copies) show the comma form
- 2026-11-08 00:11-07:00: The API was never at fault (it takes arrays); only the CLI split the raw text on commas. Fixed in the CLI: a JSON array and a bracketed list are read as lists, a wikilink entry is left alone. Chose reading over refusing so a pasted list just works.
