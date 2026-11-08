---
id: "261108-h7ds"
title: set blocked by brackets store a mangled list
status: raw
kind: bug
priority: 3
created: 2026-11-08
updated: 2026-11-08
---

# set blocked by brackets store a mangled list

## Goal

Found by the cold-agent retest of 261101-t8h8. `akmctl set <id> 'blocked_by=[a,b]'` and `blocked_by=["a","b"]` store a list of mangled strings (`["[a", "b]"]`) and report `changed`; `check` says nothing. Only `blocked_by=a,b` stores a clean list, and neither `set --help` nor the skill shows that syntax. Expected: a bracketed or JSON list is read as a list (or refused with the form to use), and the help or the `set` reference shows `blocked_by=a,b`. Done when:

- [ ] `set blocked_by=[a,b]` and `blocked_by=["a","b"]` store `[a, b]`, or are refused with `invalid_request` naming `blocked_by=a,b`.
- [ ] A test covers the bracketed, JSON and comma forms; `set --help` or the `set` reference shows the list form.

## Log

- 2026-11-08: filed.
