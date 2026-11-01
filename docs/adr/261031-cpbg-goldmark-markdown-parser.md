---
id: "261031-cpbg"
title: Read work item bodies with our own scanner, checked against goldmark as a test oracle
summary: A narrow Markdown scanner behind a Locate interface, fail closed outside its subset, with goldmark as a test-only oracle that can be promoted to production
status: proposed
date: 2026-10-31
---

# Read work item bodies with our own scanner, checked against goldmark as a test oracle

## Context

ADR 261031-gam0 models a work item body as sections with a typed core and edits it by splicing only the lines an edit changes. That needs a reader that says where each block starts, not one that prints Markdown. Spike 261031-g9zk compared goldmark, gomarkdown/markdown, blackfriday v2 and our own scanner (the one `SetSection` uses). ADR 261031-nvka sets the project's approach to large specifications: a narrow implementation, a test-only oracle, and fail closed.

Findings (scratch module outside the repository; goldmark v1.8.7):

- gomarkdown/markdown and blackfriday v2 expose no source offsets, so they cannot be the reader or the oracle.
- goldmark's headings and task lines equal the scanner's on all 125 records. Costs if used at run time: +1.36 MB on a 10.1 MB binary, about 13,700 lines, MIT, no dependencies of its own.
- A generator of random Markdown (60,000 documents in three runs, 30% of lines from a hostile pool) found the scanner wrong in two ways. Mixed fences (a backtick fence closed by `~~~`) and longer fences are mishandled by today's `SetSection` (work item 261031-2j14). Tracking fence character and length fixed both.
- Abstaining on constructs the scanner does not claim to understand (indented fences or headings, HTML lines, blockquotes, tabs, tables, ordered-list tasks, setext or thematic underlines after text) left no case where the scanner was confidently wrong against goldmark: 0 of about 30,000 compared documents. About half of the generated documents abstain, because the generator is hostile on purpose; none of the 125 real records does.

## Considered Options

- **Adopt goldmark at run time.** Exact and maintained; the first runtime dependency, and its model of the format becomes ours.
- **Our own scanner, no oracle.** No dependency; nothing says it is right.
- **Our own scanner behind an interface, goldmark as a test-only oracle.**

## Decision

The third. The body package (261031-j4bc) defines a small interface, `Locate(body []byte) (Layout, error)`, where a `Layout` lists blocks: kind (heading and level, list item, task item, fence, other), byte start and end, and line range. The scanner implements it and returns an error naming the construct when it meets one it does not claim to understand (no edit is written). A goldmark adapter implements the same interface in test files.

Tests: one battery runs against both implementations. For any input, either the scanner refuses or its `Layout` equals the oracle's on the blocks the scanner claims. The 125 records and a set of hostile cases run as ordinary tests; a fuzz target (Go's native fuzzing, seeded with the records, plus a generator over a grammar of Markdown fragments) runs on demand.

If the scanner's scope outgrows our ability to keep it right, promoting the goldmark adapter from a test file to production is a change of constructor, and a new ADR records the dependency.

## Consequences

- `go.mod` gains goldmark as a test-only dependency; the binary does not grow.
- We own a small scanner and its tests; its refusals are a visible measure of whether it is too narrow.
- A refusal is a feature: the writer edits only what it is sure of.
- Frontmatter is split off before the body is read.
- The generator and the abstention list from the prototype are the starting point of the test battery; the prototype itself is not kept.
