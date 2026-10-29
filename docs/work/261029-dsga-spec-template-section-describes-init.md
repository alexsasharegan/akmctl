---
id: "261029-dsga"
title: Rewrite the standard's template section to describe init
summary: Section 7 of the standard and the paragraph on the AKM template (a clone template) still describe the dropped template repository; describe akmctl init.
status: ready
kind: chore
priority: 2
created: 2026-10-29
updated: 2026-10-29
---

# Rewrite the standard's template section to describe init

## Goal

Section 7 of the standard (`docs/spec/agent-ready-repo-spec.md`, "Distribution and versions", the paragraph "The AKM template (a clone template)" and the heading "Two repositories, one idea each") and the guide's "A new repo" steps ("Clone the starter") describe a template repository that is dropped ([261029-j0d0](../adr/261029-j0d0-one-tool-spec-and-fixtures-in-this-repository.md)). Describe `akmctl init` with its profiles `code` and `notes` instead, and say that the tool is the only distribution. Do it with [261029-0gb4](261029-0gb4-akmctl-init-with-profiles.md), or before it with the text marked as planned.

Done when:

- [ ] Section 7 of the standard and the guide's new-repo steps describe `init` and no clone template.
- [ ] No other place in `docs/spec/` mentions the template repository or `VERSION`/`CHANGELOG` files of a template, except as history.
- [ ] `akmctl check` and Prettier pass.

## Log

- 2026-10-29: filed.
