# skills

Added in this release. Use it to see whether the AKM skills in this repository match the ones your `akmctl` ships, and to install or refresh them.

```
akmctl skills
akmctl skills install [--yes] [--dry-run]
akmctl skills update [--yes] [--dry-run]
```

With no argument (or `status`) it reports each skill as `missing`, `current`, `outdated` (the tool ships a newer copy), `edited` (someone changed it), `unmarked` (no `.akm-skill.json` marker yet) or `renamed` (a folder under a retired name).

`install` adds the skills that are missing and puts a marker in the ones that already match. `update` also replaces outdated skills and removes retired ones. An edited skill needs a second yes, so a hand edit is never lost silently. Read the plan before you confirm it.

Files changed: only inside `.agents/skills/<name>/` of a skill the tool manages. The tool never commits; commit the result yourself.

Errors: `refused` with `details.reason` `no_manifest` when the repository has no manifest. If your `akmctl` has no `skills` command, change nothing and tell the maintainer.
