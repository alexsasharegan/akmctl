---
id: "261031-g9zk"
title: spike markdown parser dependency
status: done
created: 2026-10-31
updated: 2026-10-31
closed: 2026-10-31
kind: spike
parent: 261031-gtw6
summary: "Spike: which Markdown parser, if any, gives exact block spans for splice edits, and whether a dependency is worth it"
decisions: [261031-cpbg]
---

# spike markdown parser dependency

## Goal

From the body-model discussion of 2026-10-31 (spike [261031-qbat](261031-qbat-spike-work-item-body-model.md), ADR [261031-gam0](../../adr/261031-gam0-work-item-body-model.md)). The body should be a Markdown tree with a typed layer over it, and edits should change only the bytes they mean to. Writing a spec-compliant parser ourselves is a high cost; the tool is standard-library only, and a dependency needs an ADR. Decide whether a parser dependency is worth it, and which.

What the tool needs from a parser is not a tree it can print back. It needs the source offset where each block starts and ends, so a write replaces just that range of the file and leaves every other byte as it was (a "span", used for a splice). A parser that gets block boundaries wrong (a `## ` inside a code fence, a setext heading, an HTML block, a nested list) would make the tool edit the wrong bytes.

Candidates (from memory; confirm each): goldmark with its GFM extensions, gomarkdown/markdown, blackfriday v2 (a baseline, not CommonMark), and our own scanner extending `SetSection`. Ruled out up front: tree-sitter-markdown, because it needs cgo and breaks the one static binary on macOS, Linux and Windows. Oxc is Rust and JavaScript and cannot be linked into a Go binary; its place is dev tooling, in [261030-43jw](../261030-43jw-evaluate-oxc-toolchain.md).

Done when:

- [x] Each candidate is run over every item in `docs/work/` and `docs/work/done/`, and the span of every `##` section, checklist item and log line is compared with the file; mismatches are counted and shown.
- [x] Each is run over a set of hostile cases: a `## ` inside a code fence, setext headings, HTML blocks, nested lists with task items, tables, CRLF files, a file with no trailing newline.
- [x] A splice edit (tick one item, append one log line) is done through each candidate on a scratch copy, and the git diff shows only the intended lines.
- [x] Each candidate is costed: lines of code to adopt, added binary size, licence, maintenance activity, release cadence, and how often the spec it follows changes.
- [x] Whether a parser plus our own printer could replace Prettier for records is noted as a finding only; replacing Prettier is a separate decision.
- [x] A recommendation (adopt one, extend our scanner, or defer) is in `## Outcome`; the dependency or the decision not to take one has a `proposed` ADR (the dependency rule in AGENTS.md); `261031-gam0` and `261031-j4bc` are updated to match; follow-up items are filed.

## Next

Build a scratch Go module outside the repository (the scratchpad) with goldmark, gomarkdown/markdown and blackfriday; write a harness that prints each candidate's block spans, then run the corpus and the hostile cases.

## Outcome

Recommendation: adopt goldmark (with GFM) as a reader behind an internal body package. The dependency exception needs your decision; it is drafted as the proposed ADR [261031-cpbg](../../adr/261031-cpbg-goldmark-markdown-parser.md). I added nothing to the repository's `go.mod`.

- **Method.** A scratch Go module outside the repository (goldmark v1.8.7, gomarkdown/markdown at its 2026-10-31 commit, blackfriday v2.1.0) and a harness; throwaway, not kept.
- **No offsets.** gomarkdown and blackfriday nodes carry no source position, so they cannot locate a span; ruled out for this need. blackfriday's last release is 2020 and it keeps `\r` in CRLF headings. gomarkdown counted a heading after an unclosed fence.
- **Corpus.** Over all 125 records, goldmark's `##` headings and task-list lines equal the current scanner's.
- **Hostile cases (21).** goldmark is right and the scanner wrong on: a fence indented two spaces, a four-backtick fence holding three, a `## ` inside an HTML block, a task in a blockquote, and a task in an ordered list. Setext `---` is a heading to goldmark and not to the scanner (we write ATX). Both agree on unclosed fences, indented code, headings inside quotes and lists, CRLF and no final newline.
- **Splice.** Ticking one item and appending one log line through goldmark changed only the intended lines, in LF, CRLF, no-final-newline and a file with a task-like line inside a fence. A file with no final newline gains one when a line is appended.
- **Cost.** No dependencies of its own, MIT, about 13,700 non-test lines (akmctl has about 11,200), +1.36 MB on a 10.1 MB binary, 100+ releases. A parser's spec moves slowly, so upgrades are mostly about edge cases, not new syntax.
- **Prettier.** goldmark's core has no Markdown printer. A parser plus our own printer could replace Prettier for records, but matching its output is a long tail of its own; kept as a finding, not a decision.
- **Not done.** I did not test goldmark's nested-heading offsets for setext headings, tables or HTML blocks beyond the cases above; I did not check CommonMark or GFM spec versions (from memory: CommonMark 0.31.2, GFM 0.29); I did not build the Go package.
- Follow-ups: [261031-j4bc](../261031-j4bc-body-section-tree-package.md) now builds on this once the ADR is accepted; the open spikes 261031-6bzq (operations) and 261031-x50y (schema) can start from the section model.

## Log

- 2026-10-31: filed from the body-model discussion; the section tree package 261031-j4bc waits for it.
- 2026-10-31: claimed (claude).
- 2026-10-31: made ready and claimed to start the spike.
- 2026-10-31: closed done; drafted proposed ADR 261031-cpbg.
