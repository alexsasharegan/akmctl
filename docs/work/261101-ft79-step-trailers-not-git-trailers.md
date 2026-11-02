---
id: "261101-ft79"
title: step trailers not git trailers
status: raw
created: 2026-11-01
updated: 2026-11-01
parent: 261101-vs5c
---

# step trailers not git trailers

## Goal

In the Codex rerun of 261101-0vh7 the agent passed each line of the tool's commit message as its own `git commit -m`, so AKM-Step and AKM-Item sat in separate paragraphs and git did not read them as trailers (`git interpret-trailers --parse` showed only Co-Authored-By). `akmctl commits --range` still passed them. Done when:

- [ ] `commits` reports a step commit whose AKM- trailers are not in the final trailer block (a new CM rule or a CM1 extension), with a test.
- [ ] Every command that returns a commit message says to pass it whole (`git commit -F -`), in the skill reference and the printed block.

## Log

- 2026-11-01: filed.
