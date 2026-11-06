---
id: "261101-ft79"
title: step trailers not git trailers
status: done
created: 2026-11-01
updated: 2026-11-06
closed: 2026-11-06
parent: 261101-vs5c
priority: 1
---

# step trailers not git trailers

## Goal

In the Codex rerun of 261101-0vh7 the agent passed each line of the tool's commit message as its own `git commit -m`, so AKM-Step and AKM-Item sat in separate paragraphs and git did not read them as trailers (`git interpret-trailers --parse` showed only Co-Authored-By). `akmctl commits --range` still passed them. Done when:

- [x] `commits` reports a step commit whose AKM- trailers are not in the final trailer block (a new CM rule or a CM1 extension), with a test.
- [x] Every command that returns a commit message says to pass it whole (`git commit -F -`), in the skill reference and the printed block.

## Next


## Outcome

commits reports CM5 when a step commit's AKM- trailers are not in the last paragraph, with tests for one -m per line, a blank line before attribution and a missing blank line after the subject. The printed block, akm-commit, the command references and AGENTS.md say to pass the whole message with git commit -F -. Not covered: rewriting old commits; a range over old history reports them.

## Log

- 2026-11-01: filed.
- 2026-11-06 11:12-07:00: claimed (claude-sonnet-5-5).
- 2026-11-06 11:12-07:00: Added CM5 to commits (AKM- trailers must sit in the last paragraph), with tests; the spec, glossary, akm-commit skill, command references, printed block and AGENTS.md say to pass the whole message with git commit -F -.
