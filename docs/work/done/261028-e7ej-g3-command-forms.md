---
id: "261028-e7ej"
title: Make rule G3 check more command forms than scripts/ paths
summary: G3 resolves only scripts/ paths named in AGENTS.md Commands; make and npm run targets and other forms go unchecked.
status: done
kind: feature
priority: 1
created: 2026-10-28
updated: 2026-11-16
closed: 2026-11-16
---

# Make rule G3 check more command forms than scripts/ paths

## Goal

Rule G3 warns when a command in the `AGENTS.md` Commands section names something that does not exist, but it checks only backticked `scripts/...` paths. The specification also names `make` and `npm run` targets and similar forms.

Done when:

- [x] G3 resolves `make <target>` against a `Makefile`, `npm run <name>` against `package.json` scripts, and `just <name>` against a `justfile`, when the file is present.
- [x] A target that does not exist is a warning that names the line.
- [x] A command it cannot resolve is skipped, not reported.
- [x] Fixtures cover each form.

## Next

## Outcome

- G3 resolves `make <target>`, `npm run <script>` and `just <recipe>` in code spans and fenced blocks of the Commands section: `commandPaths` in `internal/check/check.go`.
- A missing target is a warning with the line: four new cases in `TestRules`.
- A command it cannot resolve is skipped: ten near-miss cases in `internal/check/near_miss_test.go`.
- The cases edit the `basic` fixture, as the other rule cases do; no conformance case was added, because `check` cases need a whole fixture repository each.

The standard's G3 row names the forms. Decision: none.

## Log

- 2026-10-28: filed.
- 2026-11-16 06:56-07:00: claimed (claude-fable-5-2).
- 2026-11-16 06:57-07:00: ticked "G3 resolves `make <target>` against a `Makefile`, `npm run <name>` against `package.json` scripts, and `just <name>` against a `justfile`, when the file is present.": makeTargets, npmScripts, justTargets in internal/check/check.go
- 2026-11-16 06:57-07:00: ticked "A target that does not exist is a warning that names the line.": TestRules: G3 missing make target, npm script, just recipe, make target in a fence
- 2026-11-16 06:57-07:00: ticked "A command it cannot resolve is skipped, not reported.": near misses: no Makefile, an include, a pattern rule, make -C, package.json that is not JSON, a just import, the word make in prose
- 2026-11-16 06:57-07:00: ticked "Fixtures cover each form.": the G3 cases of TestRules and nearMisses edit the basic fixture, one or more for each form
- 2026-11-16 06:57-07:00: picks: only files at the repository root are read; a Makefile with an include or a pattern rule and a justfile with an import are skipped, since the scan cannot see every target; make's built-in implicit rules are not modeled
