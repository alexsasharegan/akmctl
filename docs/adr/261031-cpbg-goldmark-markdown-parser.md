---
id: "261031-cpbg"
title: Use goldmark as the Markdown parser behind a body package
summary: Take goldmark (GFM) as the tool's one runtime dependency, used only to find block structure and line offsets, behind an internal package
status: proposed
date: 2026-10-31
---

# Use goldmark as the Markdown parser behind a body package

## Context

ADR 261031-gam0 models a work item body as a lossless list of sections with a typed core, and edits it by changing only the lines an edit means to change. That needs a parser that says where each block starts, not one that can print Markdown back. The tool is standard library only; a runtime dependency needs an ADR (AGENTS.md). Spike 261031-g9zk compared three Go parsers and our own scanner (the one `SetSection` uses).

Findings (scratch module outside the repository; goldmark v1.8.7, gomarkdown/markdown at its 2026-10-31 commit, blackfriday v2.1.0):

- gomarkdown/markdown and blackfriday v2 expose no source offsets on their nodes, so they cannot locate a span; blackfriday's last release is 2020 and it leaves `\r` in CRLF headings; gomarkdown counted a `## ` after an unclosed fence as a heading.
- goldmark gives each block's source lines. Its `##` headings and task-list lines match the scanner on all 125 records in `docs/work/`.
- On 21 hostile cases goldmark is right where the scanner is wrong in four: a fence indented by two spaces, a four-backtick fence holding three, a `## ` inside an HTML block, and a task item inside a blockquote or an ordered list. Both agree on unclosed fences, indented code, headings inside quotes and list items, CRLF and a missing final newline.
- A splice edit (tick one item, append one log line) through goldmark changed only the intended lines in LF, CRLF, no-final-newline and fenced-task files.
- Cost: no dependencies of its own, MIT, about 13,700 lines, +1.36 MB on a 10.1 MB binary (+13%), 100+ releases and the latest on 2026-10-30.
- goldmark's core has no Markdown printer, so it is a reader, not a formatter.

## Considered Options

- **Our own scanner, extended.** No dependency. Needs rules for indented and longer fences, HTML blocks and nested lists; it would keep growing toward a parser.
- **gomarkdown/markdown or blackfriday v2.** Ruled out: no offsets.
- **goldmark.** One dependency, exact enough, actively kept.

## Decision

Use goldmark with its GFM extension for reading. All use is inside one internal package (`internal/body`, work item 261031-j4bc), which exposes our section and node types and offsets; nothing else imports goldmark, so it can be replaced. Frontmatter is split off by the tool before the body is parsed. Edits splice the lines the parser located and do not re-print the document. Whether the tool later prints Markdown itself, or replaces Prettier, is a separate decision.

## Consequences

- The tool gains its first runtime dependency; `go.mod` and `go.sum` grow, and the binary grows by about 1.4 MB.
- Block structure beyond what our records use (nested lists, HTML blocks, odd fences) is handled by a maintained parser.
- Upgrades to goldmark can change offsets or what counts as a task item; the tests over `docs/work/` and the hostile cases in 261031-g9zk guard that.
- If the dependency is later refused, the fallback is the extended scanner behind the same package.
