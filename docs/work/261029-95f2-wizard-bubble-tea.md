---
id: "261029-95f2"
title: Replace the init questionnaire with a Bubble Tea wizard
status: ready
created: 2026-10-29
updated: 2026-10-29
kind: feature
priority: 3
summary: The beta init questionnaire uses plain y/N prompts; a Bubble Tea wizard is a possible later enhancement.
---

# Replace the init questionnaire with a Bubble Tea wizard

## Goal

The beta asks its `init` questions with plain line prompts (`Install the git skill? [y/N]`, the capitalised letter is the default), using only the standard library. A terminal wizard built on Charm's Bubble Tea would look better and handle more questions, but it adds a dependency tree to a standard-library-only tool and more terminal handling to maintain, so it was deferred (decided with Alexander, 2026-10-29).

Revisit when the questionnaire has more than a few questions or users ask for it. Adding the dependency needs its own ADR.

Done when:

- [ ] An ADR records the dependency decision (versions, size, transitive modules, macOS and Linux behaviour).
- [ ] The wizard runs only on a terminal and never blocks an agent; flags and `--yes` still skip it.
- [ ] The questions and answers stay in the CLI layer; the API `init` operation is unchanged.

## Next

Not started; deferred until after the 0.2.0 beta.

## Log

- 2026-10-29: filed.
