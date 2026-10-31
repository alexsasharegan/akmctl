---
id: "261030-1wj5"
title: Files the tool keeps for itself carry the akm namespace; files a person recognises do not
summary: System files the CLI maintains are hidden and named .akm-*; the files a person recognises (.layout, records, roadmap) are plain; .migrate-log is grandfathered.
status: proposed
date: 2026-10-30
---

# Files the tool keeps for itself carry the akm namespace; files a person recognises do not

## Context

The tool writes more kinds of file than it used to: the manifest, records, the lock, the deny list, temporary files, the migration log, and now a marker in each skill folder (ADR 261030-h6t7) and a planned update-check state file (work item 261030-55cs). The names grew by habit. ADR 261029-tb3d fixes the `akm-` prefix for skills only; no record says how the other files are named, so each new file asks the question again. Today `.akm-deny`, `.akm.lock` and the temp files `.akm-*.tmp` carry the namespace, `.layout` does not, and `.migrate-log` does not either although it is metadata the CLI keeps.

## Considered Options

- Name each file as it comes. Rejected: it is what produced `.akm.lock` next to `.akm-deny`, and `.migrate-log` beside them.
- Prefix every file the tool touches, `.layout` and the records included. Rejected: the manifest is the one configuration file a person is meant to recognise, and records are the person's own content.
- Prefix nothing. Rejected: a bare `.migrate-log` or `.lock` in a repository gives no hint of its owner, and the tool must not mistake another tool's file for its own.
- Two classes: files a person recognises or edits keep plain names; files the CLI keeps for itself carry the namespace.

## Decision

There are two classes of file.

**User-facing files** are the ones a person opens, edits or is meant to recognise: the manifest `.layout` (the one special configuration file), the docs folder and what it holds (records, `ROADMAP.md`, the ADR index), `AGENTS.md` and `CLAUDE.md`. They keep plain names, chosen with the person at `init`.

**System files** are custodial metadata the CLI handles and a person need not care about. They are hidden (a leading dot) and carry the namespace as `.akm-<what>`, with an extension only when the content needs one: `.akm-skill.json` in each managed skill folder, `.akm-*.tmp` for a temporary file next to its target, and any new file of this kind. Outside the repository, in the user's cache or config folder, the folder is `akmctl` and a file is named `akm-<what>.json` (the update-check state is `akm-update.json`).

Skills keep the rule of ADR 261029-tb3d (`akm` and `akm-<name>`); a skill is read by an agent and is not a system file in this sense.

Existing names stay: `.akm-deny` (a list a person may edit, but only meaningful to this tool's safety extension), `.akm.lock` (documented and listed in `.gitignore` templates) and `.migrate-log`. They are grandfathered, not examples. `.migrate-log` is the one that breaks the rule; it moves to `.akm-migrate-log` at the next layout version, with a migration step, and until then the tool reads and writes only the old name.

The test for a new file: would the maintainer open it or reason about it? Then it is user-facing. Is it written and read only by the CLI? Then it is a system file and takes the prefix.

## Consequences

- New system files have a name decided in advance; reviewers can point at this record.
- Nothing is renamed now. One rename is scheduled, with the layout bump that carries it.
- `.layout` stays unprefixed, so a person still recognises it, and the standard keeps a plain name for the file other tools may read.
- The glossary lists the two classes.
- A person looking at a repository can tell from the name what is safe to ignore and what to read.

## Confirmation

A test lists the dotfiles the tool writes in a repository made by `init` and a run of the write operations, and fails if one is neither `.layout` nor starts with `.akm`, except the grandfathered `.migrate-log` named in the test.

## Review

Proposed 2026-10-30 after the maintainer asked whether a record states the namespace rule and described the two classes. The grandfathering of `.migrate-log` is my recommendation; the maintainer may prefer to rename it now.

## Revert

Delete this record. Nothing depends on it; new files would again be named case by case.
