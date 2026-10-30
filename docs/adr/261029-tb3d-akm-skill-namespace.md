---
id: "261029-tb3d"
title: Tool-owned skills live under the akm- prefix, with a router skill
summary: Tool-owned skills live under the akm- prefix, with a router skill
status: accepted
date: 2026-10-29
---

# Tool-owned skills live under the akm- prefix, with a router skill

## Context

`akmctl init` writes agent skills into `.agents/skills/`. Until now they were called `work-items` and `record-decision`, which can collide with a repository's own skills, and `AGENTS.md` explained how to use the tool in prose. Agents load a skill by its description and read its references only when needed, so a small router costs few tokens.

## Considered Options

- Keep the plain names. Rejected: they collide easily and cannot be recognised as the tool's.
- Explain every command in `AGENTS.md`. Rejected: it is read in full on every session and drifts from the tool.
- Prefix every tool-owned skill with `akm-` and add a router skill.

## Decision

Every skill the tool installs or may remove is named `akm-<name>`: `akm-work-items`, `akm-record-decision` and the optional `akm-commit`. A router skill named `akm` has a short `SKILL.md` (a table of commands) and one reference file per command under `references/`, holding when to use it, an example, the files it changes and common errors. The extension named `work-items` in the manifest keeps its name; only skill directories are renamed. The prefix marks what the tool owns: the tool never writes or removes a skill outside it.

## Consequences

- Skill directories are renamed in this repository and in what `init` writes; the spec, the guide and `AGENTS.md` follow.
- A repository's own skills cannot clash with the tool's unless they also start with `akm-`; if they do, `init` reports a conflict and does not overwrite.
- Testers who already copied the old names must rename them; the beta is the first release that writes skills, so this is cheap now and costly later.
- The references need a test that every operation has one.
