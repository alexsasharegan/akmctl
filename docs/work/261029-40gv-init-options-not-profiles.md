---
id: "261029-40gv"
title: Replace init profiles with options; move the manifest into the docs dir
status: ready
created: 2026-10-29
updated: 2026-10-29
kind: feature
priority: 1
summary: init has --docs-dir, --wikilinks and git options instead of profiles; the manifest, deny list, lock and migrate log live in the docs dir.
---

# Replace init profiles with options; move the manifest into the docs dir

## Goal

Per ADR 261029-hsfd: replace `init --profile code|notes` with ordinary options and move the manifest into the docs dir. Details in the ADR.

Done when:

- [ ] `init` has `--docs-dir <single name>` (default `docs`), `--wikilinks`/`--no-wikilinks`, `--git-commits`/`--no-git-commits`, `--git-init`/`--no-git-init`, `--yes`; `profile` is gone from the API, CLI, HTTP, schema, fixtures, spec, guide, README and skill references; a bare `init` with no terminal succeeds with the defaults.
- [ ] The terminal questions are: where do the records go (default `docs`), use wikilinks (y/N), install the commit skill (y/N), run git init (Y/n, only when not in a repository).
- [ ] The manifest, `.akm-deny`, `.akm.lock` and `.migrate-log` live in the docs dir; discovery finds `docs/.layout` first, then exactly one `<name>/.layout` among the root's immediate non-hidden subdirectories, and two candidates is an error naming both; check, index, new, set, close, claim, migrate, config and the workbench follow it.
- [ ] `AGENTS.md` from init has no build, test, lint or run lines; the format-check line uses the chosen docs dir.
- [ ] `.gitignore` is written only in a git repository: the lock file under the docs dir, plus Obsidian's per-machine files only when `.obsidian/` exists.
- [ ] A repository whose manifest is at `docs/.layout` behaves exactly as before (this repository's own checks pass unchanged).
- [ ] Tests and fixtures (both copies) cover a custom docs dir, ambiguity, a repository found from a subdirectory, wikilinks on/off, and the git cases; `npm run schemas` passes.

## Next

Not started.

## Log

- 2026-10-29: filed.
