---
id: "261108-xkck"
title: frontmatter merge key value
status: raw
kind: bug
created: 2026-11-08
updated: 2026-11-08
priority: 3
---

# frontmatter merge key value

## Goal

Found by a 2-second run of `go test ./internal/frontmatter -run "^$" -fuzz ^FuzzAgainstOracle$ -fuzztime 2s` (the new fuzz job of 261108-z9n8). Input: the single line `A: <<`. The tool's frontmatter reader accepts it (value the text `<<`); the YAML oracle in the test reads `<<` as the merge tag (!!merge) and fails, so `FuzzAgainstOracle` reports "tool accepts, oracle fails". Decide which is right for the AKM subset (spec: frontmatter subset, rule S2/H-rules) and fix the reader, the oracle filter or the spec; add the input as a seed.

Done when:

- [ ] The input `A: <<` no longer fails `FuzzAgainstOracle` (reader, oracle filter or spec changed, decision recorded).
- [ ] The input is a seed of the fuzz target.

## Log

- 2026-11-08: filed.
