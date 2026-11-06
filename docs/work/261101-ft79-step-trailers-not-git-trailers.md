---
id: "261101-ft79"
title: step trailers not git trailers
status: active
created: 2026-11-01
updated: 2026-11-06
parent: 261101-vs5c
priority: 1
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-06T18:12
---

# step trailers not git trailers

## Goal

In the Codex rerun of 261101-0vh7 the agent passed each line of the tool's commit message as its own `git commit -m`, so AKM-Step and AKM-Item sat in separate paragraphs and git did not read them as trailers (`git interpret-trailers --parse` showed only Co-Authored-By). `akmctl commits --range` still passed them. Done when:

- [x] `commits` reports a step commit whose AKM- trailers are not in the final trailer block (a new CM rule or a CM1 extension), with a test.
- [x] Every command that returns a commit message says to pass it whole (`git commit -F -`), in the skill reference and the printed block.

## Next

Open the pull request; the maintainer merges.

## Log

- 2026-11-01: filed.
- 2026-11-06 11:12-07:00: claimed (claude-sonnet-5-5).
- 2026-11-06 11:12-07:00: Added CM5 to commits (AKM- trailers must sit in the last paragraph), with tests; the spec, glossary, akm-commit skill, command references, printed block and AGENTS.md say to pass the whole message with git commit -F -.
