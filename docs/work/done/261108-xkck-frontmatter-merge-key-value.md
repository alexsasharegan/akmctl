---
id: "261108-xkck"
title: frontmatter merge key value
status: done
kind: bug
created: 2026-11-08
updated: 2026-11-16
closed: 2026-11-16
priority: 3
---

# frontmatter merge key value

## Goal

Found by a 2-second run of `go test ./internal/frontmatter -run "^$" -fuzz ^FuzzAgainstOracle$ -fuzztime 2s` (the new fuzz job of 261108-z9n8). Input: the single line `A: <<`. The tool's frontmatter reader accepts it (value the text `<<`); the YAML oracle in the test reads `<<` as the merge tag (!!merge) and fails, so `FuzzAgainstOracle` reports "tool accepts, oracle fails". Decide which is right for the AKM subset (spec: frontmatter subset, rule S2/H-rules) and fix the reader, the oracle filter or the spec; add the input as a seed.

Done when:

- [x] The input `A: <<` no longer fails `FuzzAgainstOracle` (reader, oracle filter or spec changed, decision recorded).
- [x] The input is a seed of the fuzz target.

## Next

## Outcome

- `A: <<` no longer fails `FuzzAgainstOracle`: the reader refuses a plain scalar that is exactly `<<` (`typePlain` in `internal/frontmatter/parse.go`), with a message that says to quote it. A 20-second fuzz run found nothing.
- The input is a seed: three new cases in `yaml-lite.json` (both copies), which the fuzz target loads.

The standard lists it among the forms that are stricter than YAML. Decision: none (the pick is in the Log).

## Log

- 2026-11-08: filed.
- 2026-11-16 07:09-07:00: claimed (claude-fable-5-2).
- 2026-11-16 07:10-07:00: ticked "The input `A: <<` no longer fails `FuzzAgainstOracle` (reader, oracle filter or spec changed, decision recorded).": typePlain refuses a plain <<; go test -fuzz FuzzAgainstOracle -fuzztime 20s: 115950 runs, no failure
- 2026-11-16 07:10-07:00: ticked "The input is a seed of the fuzz target.": yaml-lite.json case 'merge key as a value' (input a: <<); the fuzz target seeds from the fixtures
- 2026-11-16 07:10-07:00: pick: the reader changes, not the oracle filter or the meaning in the spec. YAML 1.2 core reads << as text, so the tool was right by the letter; but the subset already keeps clear of forms that libraries disagree on (number spellings, capitalized booleans), and a merge tag is such a form. A record that holds a bare << value becomes unreadable until it is quoted; none is known
