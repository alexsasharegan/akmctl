---
id: "261029-hsfd"
title: init takes options (docs dir, wikilinks, git) instead of profiles; the manifest lives in the docs dir
summary: init takes options (docs dir, wikilinks, git) instead of profiles; the manifest lives in the docs dir
status: accepted
date: 2026-10-29
---

# init takes options (docs dir, wikilinks, git) instead of profiles; the manifest lives in the docs dir

## Context

ADR 261029-j0d0 said `akmctl init` takes a profile, `code` or `notes`. Looking at what the profiles differ in showed three independent choices bundled together: where the records go, whether links are wikilinks, and project-specific lines in `AGENTS.md`. `notes` also hard-coded a folder named `notes/`, a contested name, and left a `docs/` folder holding only the manifest and other hidden files. Nothing has been released with profiles.

## Considered Options

- Keep two profiles. Rejected: the bundles do not match real needs (a code project may want wikilinks; a vault may want a folder named anything).
- Options on `init`, with the manifest staying in `docs/`. Rejected: a custom docs dir would leave a `docs/` folder with only hidden files.
- Options on `init`, with the manifest, deny list, lock file and migrate log in the chosen dir, found by scanning one level.

## Decision

`akmctl init` has no profile. Its options are `--docs-dir <name>` (a single folder name, default `docs`; work items, ADRs and the roadmap go under it), `--wikilinks` (sets `reference_style: wikilink`, default off), `--git-commits` and `--git-init`. A bare `init` works with defaults and needs no terminal. The `AGENTS.md` it writes has no placeholder build, test, lint or run lines.

The manifest `.layout`, the deny list `.akm-deny`, the lock `.akm.lock` and the migrate log `.migrate-log` live in the docs dir. To find them the tool looks for `docs/.layout` first, then for exactly one `<name>/.layout` among the repository root's immediate, non-hidden subdirectories. Two candidates are an error that names both. The repository root is still the nearest ancestor with a manifest or `AGENTS.md`. The `work_dir`, `adr_dir` and `roadmap` keys stay, now relative to the repository root as before; `init` writes them from `--docs-dir`.

`init` writes a `.gitignore` only when the directory is or becomes a git repository. It lists the lock file, and Obsidian's per-machine files only when an `.obsidian/` folder already exists.

## Consequences

- The `code` and `notes` profiles, the `profile` parameter, its schema, fixtures, spec text, README and skill references change. The `notes/` folder name disappears.
- Repository discovery gains the one-level scan; every place that hard-codes `docs/.layout`, `docs/.akm.lock`, `docs/.akm-deny` or `docs/.migrate-log` follows the manifest's directory instead.
- A docs dir name is a single segment; nested names are refused so discovery stays one level deep.
- 261029-j0d0's sentence about profiles is superseded by this record; the rest of it stands.
- Existing repositories whose manifest is at `docs/.layout` are unchanged.
