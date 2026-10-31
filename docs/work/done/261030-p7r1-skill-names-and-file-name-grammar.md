---
id: "261030-p7r1"
title: Skill names and file-name grammar
summary: Rename akm-work-items to a verb-noun name, write the naming grammar for skills and files into the glossary, and let the installer retire the old name.
status: done
kind: chore
priority: 2
decisions: [261030-1wj5]
created: 2026-10-30
updated: 2026-10-30
closed: 2026-10-30
---

# Skill names and file-name grammar

## Goal

Names are the interface agents and people read first. `akm-record-decision` is verb and noun (and reads as "AKM decision record"); `akm-work-items` is a noun only. The maintainer wants the set to follow one grammar and the grammar to be written down. The rename belongs before the skills installer (261030-8d8t) is built: nothing outside this repository has the skills yet (the maintainer's vault has none), so the installer can ship with final names instead of migrating twice.

Proposed grammar, for the glossary:

- **Skills.** The router is `akm`. Every other skill is `akm-<verb>[-<object>]`, an imperative the agent can read as an instruction: `akm-record-decision`, `akm-commit`, and the renamed work skill. Prefix `akm-` marks tool ownership (ADR 261029-tb3d); the extension named `work-items` in the manifest keeps its name.
- **Files.** User-facing files (`.layout`, records, roadmap, ADR index) have plain names; system files carry `.akm-` (ADR 261030-1wj5, proposed).
- **Words.** The glossary already says "record" is any file with an id and "work item" is a record of type `work`; skill and command names follow it where they can.

The new name. Recommended: **`akm-track-work`**. The skill covers the whole life of an item: find, claim, checkpoint, close. "Record" is the wrong verb for that (it would read as creating, and "record" is already the glossary's umbrella noun for any file). Alternatives: `akm-record-work` (the maintainer's idea; matches `akm-record-decision` but implies only filing), `akm-work-item` (still a noun), `akm-manage-work` (vague). The maintainer decides in the ADR.

Risk: low. About 25 files name the skill: both copies of the skills, the router table, the `init` templates and `AGENTS.md.tmpl`, this repository's `AGENTS.md`, `README.md`, the spec and guide, the `init` fixtures and their generator, and tests. Accepted ADRs name it too and stay unedited; a new ADR records the change. Repositories made by 0.3.x `init` keep the old folder until `skills update` runs; the installer needs a `renamed` table in the binary (old name, new name, first version) and one new rule, below.

Done when:

- [x] An ADR (`proposed`) records the grammar, the new name, and the installer rule for a renamed skill: `status` shows `renamed` for an old folder; `update` writes the new skill and removes the old folder only when its files equal what its marker says was written (or with `confirm_edited` after a diff), and reports it otherwise. This narrows ADR 261030-h6t7 ("neither deletes"). No version-range resolution is built: the table and the marker's tool version are enough while the tool is unstable.
- [x] The skill is renamed in `.agents/skills/` and `internal/initrepo/templates/skills/` (the drift test passes), its `name:` and heading follow, and every mention in the files above is updated; accepted ADRs and closed items are not edited.
- [x] The `init` fixtures and `make_fixtures.py` are regenerated and synced to `testdata/api`; `go test ./...` passes.
- [x] The glossary has a "Skill and file vocabulary" section with the grammar and the two file classes.
- [x] A test lists the dotfiles written by `init` and by the write operations and fails on one that is neither `.layout` nor `.akm*` (except the grandfathered `.migrate-log`), as ADR 261030-1wj5's Confirmation says.
- [x] The changelog says the skill was renamed and how an existing repository gets the new name once the installer exists.

## Outcome

- ADR 261030-wsqb (`proposed`): the skill grammar, the rename to `akm-track-work` (the maintainer chose the name on 2026-10-30) and the installer rule for a renamed skill; ADR 261030-1wj5 (system files carry the `akm` namespace) is accepted. Both are in `decisions`.
- `akm-work-items` is `akm-track-work` in `.agents/skills/` and `internal/initrepo/templates/skills/` (renamed with `git mv`, heading and `name:` updated), and in the router, `init`'s `AGENTS.md` template, this repository's `AGENTS.md` and the README. The `init` fixtures and `make_fixtures.py` are regenerated and synced to `testdata/api`. Accepted ADRs, closed work items, the changelog history and `docs/design/` are unchanged.
- `.migrate-log` is `.akm-migrate-log` (`record.MigrateLogPath`, spec, architecture); a migration that finds the old file renames it first (`TestMigrateLogKeepsHistoryUnderTheOldName`). The layout version is not bumped.
- Tests: `TestSkillNamesFollowTheGrammar` (every skill folder is `akm` or `akm-<verb>[-<object>]`) and `TestHiddenFilesFollowTheNamespaceRule` (every hidden file `init` and the writes leave is `.layout`, `.akm*` or a known file of another tool).
- The glossary has "Skill and file vocabulary"; the changelog has both changes.
- Not built here: the installer's `renamed` status and the removal of the old folder; that is in 261030-8d8t, which the wsqb ADR specifies and which now builds on the final names.
- Checks: `go vet`, `go test ./...`, `gofmt -l`, `akmctl check`, `npm run format:check`, `npm run schemas`.

Decision: ADR 261030-wsqb, proposed; ADR 261030-1wj5, accepted.

## Log

- 2026-10-30: filed, from the maintainer's naming feedback on `akm-work-items` and the question of whether a record states the file-namespace rule (it did not: ADR 261030-1wj5 is the draft).
- 2026-10-30: claimed (claude).
- 2026-10-30: done; closed.
- 2026-10-30: released (claude).
