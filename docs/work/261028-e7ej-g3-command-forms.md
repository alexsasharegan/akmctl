---
id: "261028-e7ej"
title: Make rule G3 check more command forms than scripts/ paths
summary: G3 resolves only scripts/ paths named in AGENTS.md Commands; make and npm run targets and other forms go unchecked.
status: active
kind: feature
priority: 1
created: 2026-10-28
updated: 2026-11-16
claimed_by: claude-fable-5-2
claimed_at: 2026-11-16T13:56
---

# Make rule G3 check more command forms than scripts/ paths

## Goal

Rule G3 warns when a command in the `AGENTS.md` Commands section names something that does not exist, but it checks only backticked `scripts/...` paths. The specification also names `make` and `npm run` targets and similar forms.

Done when:

- [ ] G3 resolves `make <target>` against a `Makefile`, `npm run <name>` against `package.json` scripts, and `just <name>` against a `justfile`, when the file is present.
- [ ] A target that does not exist is a warning that names the line.
- [ ] A command it cannot resolve is skipped, not reported.
- [ ] Fixtures cover each form.

## Next

Extend commandPaths in internal/check/check.go to resolve make, npm run and just targets; add rule and near-miss cases; name the forms in the standard's G3 row.

## Log

- 2026-10-28: filed.
- 2026-11-16 06:56-07:00: claimed (claude-fable-5-2).
