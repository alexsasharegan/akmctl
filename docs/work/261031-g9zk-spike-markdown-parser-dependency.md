---
id: "261031-g9zk"
title: spike markdown parser dependency
status: raw
created: 2026-10-31
updated: 2026-10-31
kind: spike
parent: 261031-gtw6
summary: "Spike: which Markdown parser, if any, gives exact block spans for splice edits, and whether a dependency is worth it"
---

# spike markdown parser dependency

## Goal

From the body-model discussion of 2026-10-31 (spike [261031-qbat](done/261031-qbat-spike-work-item-body-model.md), ADR [261031-gam0](../adr/261031-gam0-work-item-body-model.md)). The body should be a Markdown tree with a typed layer over it, and edits should change only the bytes they mean to. Writing a spec-compliant parser ourselves is a high cost; the tool is standard-library only, and a dependency needs an ADR. Decide whether a parser dependency is worth it, and which.

What the tool needs from a parser is not a tree it can print back. It needs the source offset where each block starts and ends, so a write replaces just that range of the file and leaves every other byte as it was (a "span", used for a splice). A parser that gets block boundaries wrong (a `## ` inside a code fence, a setext heading, an HTML block, a nested list) would make the tool edit the wrong bytes.

Candidates (from memory; confirm each): goldmark with its GFM extensions, gomarkdown/markdown, blackfriday v2 (a baseline, not CommonMark), and our own scanner extending `SetSection`. Ruled out up front: tree-sitter-markdown, because it needs cgo and breaks the one static binary on macOS, Linux and Windows. Oxc is Rust and JavaScript and cannot be linked into a Go binary; its place is dev tooling, in [261030-43jw](261030-43jw-evaluate-oxc-toolchain.md).

Done when:

- [ ] Each candidate is run over every item in `docs/work/` and `docs/work/done/`, and the span of every `##` section, checklist item and log line is compared with the file; mismatches are counted and shown.
- [ ] Each is run over a set of hostile cases: a `## ` inside a code fence, setext headings, HTML blocks, nested lists with task items, tables, CRLF files, a file with no trailing newline.
- [ ] A splice edit (tick one item, append one log line) is done through each candidate on a scratch copy, and the git diff shows only the intended lines.
- [ ] Each candidate is costed: lines of code to adopt, added binary size, licence, maintenance activity, release cadence, and how often the spec it follows changes.
- [ ] Whether a parser plus our own printer could replace Prettier for records is noted as a finding only; replacing Prettier is a separate decision.
- [ ] A recommendation (adopt one, extend our scanner, or defer) is in `## Outcome`; the dependency or the decision not to take one has a `proposed` ADR (the dependency rule in AGENTS.md); `261031-gam0` and `261031-j4bc` are updated to match; follow-up items are filed.

## Log

- 2026-10-31: filed from the body-model discussion; the section tree package 261031-j4bc waits for it.
