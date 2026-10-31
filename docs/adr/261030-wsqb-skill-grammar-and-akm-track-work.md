---
id: "261030-wsqb"
title: Skills are named akm-verb-object, and akm-work-items becomes akm-track-work
summary: Skill names are an imperative, akm-<verb>[-<object>]; akm-work-items is renamed akm-track-work, and skills update retires the old folder when it is unedited.
status: accepted
date: 2026-10-30
---

# Skills are named akm-verb-object, and akm-work-items becomes akm-track-work

## Context

ADR 261029-tb3d fixed the `akm-` prefix and the router `akm`, but not what follows the prefix. The skills read `akm-record-decision` (a verb and its object, which also reads as "AKM decision record"), `akm-commit` (a verb) and `akm-work-items` (a noun alone). A skill name is what an agent matches against a situation and what a person types or reads in a listing, so the set should follow one grammar. No repository outside this one has the skills yet (the maintainer's vault has none; ADR 261030-h6t7 adds the installer), so a rename is cheap now and costly later.

## Considered Options

- Keep `akm-work-items`. Rejected: it is the only name without an action, which the maintainer found jarring, and it names a thing where the other names say what to do.
- `akm-record-work`. Considered, to match `akm-record-decision`. Rejected: the skill covers the whole life of an item (find, claim, checkpoint, close), and "record" suggests only filing it; the glossary also uses "record" as the umbrella noun for any file with an id.
- `akm-manage-work`. Rejected: vague.
- `akm-track-work`.
- Resolve names through version ranges so old and new names both work. Rejected for now: the tool is unstable, and a small table of renames plus the version in each skill's marker is enough.

## Decision

**Grammar.** The router is `akm`. Every other skill is `akm-<verb>[-<object>]`, an imperative an agent can read as an instruction: `akm-record-decision`, `akm-commit`, `akm-track-work`. The `akm-` prefix marks tool ownership (ADR 261029-tb3d). The extension named `work-items` in the manifest keeps its name; only the skill is renamed. The grammar and the file-name classes of ADR 261030-1wj5 are written in the glossary.

**The rename.** `akm-work-items` becomes `akm-track-work`, in this repository's `.agents/skills/`, in what `init` writes, and in every document and test that names it, except accepted ADRs and closed work items, which stay as written.

**Retiring the old name.** The binary carries a table of renames (old name, new name, first version that has the new name). `akmctl skills status` shows a folder with an old name as `renamed`. `skills update` installs the new skill and then removes the old folder only when every file in it equals what its `.akm-skill.json` marker says was written (or, for a folder with no marker, equals the shipped file of the old name); otherwise it shows the difference and removes it only with `confirm_edited`. A folder it cannot remove is reported, not an error. This is the single deletion that ADR 261030-h6t7 does not allow; it narrows that record for renamed skills only, and `config git-commits off` keeps its own removal rule. No version-range resolution is built.

## Consequences

- About 25 files change in one commit; the `init` fixtures change because they list the files `init` writes.
- A repository made by the 0.3.x `init` keeps `akm-work-items` until `skills update` runs, which exists only once work item 261030-8d8t is built; nothing breaks in between, and the old skill keeps working.
- Accepted ADRs and closed items still say `akm-work-items`. The glossary notes the rename.
- The next skill that is added or renamed has a grammar to follow, and a rename has a mechanism.

## Confirmation

The skill drift test and `TestEveryOperationHasReference` keep both copies equal and complete; a test lists the skill directories and fails on a name that is not `akm` or does not match `akm-<verb>[-<object>]` in lowercase letters and hyphens. When the installer is built, a test shows that a folder named `akm-work-items` with an unedited marker is removed by `update` and an edited one is not.

## Review

Proposed 2026-10-30. The maintainer chose the name `akm-track-work` on 2026-10-30 after the lifecycle argument, and accepted this record on 2026-10-30.

## Revert

Rename the folder back in both copies and the documents and regenerate the `init` fixtures; delete the renames table.
