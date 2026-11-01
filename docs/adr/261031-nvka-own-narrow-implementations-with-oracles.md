---
id: "261031-nvka"
title: Write narrow implementations of large specifications and check them against oracles
summary: For large specifications such as Markdown and YAML, write a narrow implementation of our subset, check it against a test-only oracle, and fail closed
status: accepted
date: 2026-10-31
scope: [go.mod, internal]
authored_by: agent
---

# Write narrow implementations of large specifications and check them against oracles

## Context

ADR 261027-ebht chose standard library only for the binary. That was a strong preference, never a strict rule, and the discipline of following it has had a consistent shape: the frontmatter parser is our own with a YAML library as a test oracle (ADR 261028-h99k), and display width is our own rule with Prettier as the reference (ADR 261028-mw1r). The same question came up again for the work item body (spike 261031-g9zk): a spec-compliant Markdown parser (goldmark) against extending our own scanner. Goldmark matched the scanner on all 125 records and differed only on constructs our records do not contain.

## Considered Options

- **Bring in a library for each large specification.** Less code, correct on the long tail, but each one is a runtime dependency to vet and update, and its model of the format becomes ours.
- **Write everything ourselves with no reference.** Nothing to vet, but nothing to say we are right.
- **Write a narrow implementation of the subset we use, and check it against a library kept out of the binary.**

## Decision

For a large specification (Markdown, YAML, and the like), we write our own narrow implementation of the subset our records use and verify it against an established library as an oracle. The oracle is a test-only dependency (as `go.yaml.in/yaml/v3` is today). The practice has four parts:

1. **Narrow scope.** Support what we can be sure about, and name what we do not support.
2. **Fail closed.** Outside the supported subset, the tool refuses with an error that names the construct and writes nothing. Being cautious is correct; being confidently wrong is the bug.
3. **A seam.** The code that reads the format sits behind a small interface that the oracle can also implement, so one test battery runs against both, and swapping the oracle in for production is a change of constructor.
4. **Differential tests.** For any input, either our implementation refuses or it agrees with the oracle on what we claim to understand. The seed corpus is the repository's own records, run on every test; fuzzing, which mutates it, runs on demand.

This is still not a hard rule. A runtime dependency stays possible when it is clearly worth it, but the threshold is high and it needs an ADR that says why a narrow implementation with an oracle does not do. If our implementation hits bugs we cannot fix within the scope, the seam lets us adopt the oracle.

## Consequences

- More code to write and own, and tests that are the main protection against being wrong.
- The binary stays standard library only; `go.mod` carries the oracles, as it does for YAML.
- Each new use of the pattern needs its seam, its refusal behaviour and a differential test, and registers its terms in the glossary.
- ADR 261027-ebht stays accepted; this one states how its preference is applied.
- The Markdown parser decision (261031-cpbg) is the next use of it.
