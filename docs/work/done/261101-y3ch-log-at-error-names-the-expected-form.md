---
id: "261101-y3ch"
title: A pattern-checked parameter's error does not say the expected form
status: done
created: 2026-11-01
updated: 2026-11-06
closed: 2026-11-06
kind: bug
summary: log --at 09:46 fails with 'has the wrong form' and never says what form is valid; the same message serves every pattern parameter
parent: 261101-vs5c
---

# A pattern-checked parameter's error does not say the expected form

## Goal

`akmctl log <id> text --at 09:46` fails with `invalid_request: at: "09:46" has the wrong form`. The message comes from `checkString` (`internal/api/params.go:216`) for every parameter with a pattern, and the form is in none of them, so the caller has to find a reference page. Done when:

- [x] Each pattern parameter carries a short description of its form (for `at`: `YYYY-MM-DD HH:MM`), and the error appends it.
- [x] A test fails for a pattern parameter that has no description.

## Next

## Outcome

Each pattern parameter has a form description in params.go and the invalid_request message appends it ('expected YYYY-MM-DD HH:MM ...'). A test fails for a pattern with no description, and one checks three messages. Not changed: the slug and cycle checks in ops_write.go keep their own messages.

## Log

- 2026-11-01: filed.
- 2026-11-01 10:31-07:00: Filed from spike 261101-baze.
- 2026-11-06 11:48-07:00: claimed (claude-sonnet-5-5).
- 2026-11-06 11:49-07:00: Added a form description beside each pattern in params.go and appended it to the error; two tests (every pattern has one; three messages name the form).
