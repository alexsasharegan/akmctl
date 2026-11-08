# skills

Added in this release. Use it to see whether the AKM skills in this repository match the ones your `akmctl` ships, and to install or refresh them.

```
akmctl skills
akmctl skills install [--yes] [--dry-run]
akmctl skills update [--yes] [--dry-run]
```

With no argument (or `status`) it reports each skill as `missing`, `current`, `outdated` (the tool ships a newer copy), `edited` (someone changed it), `unmarked` (no `.akm-skill.json` marker yet) or `renamed` (a folder under a retired name).

`install` adds the skills that are missing and puts a marker in the ones that already match. `update` also replaces outdated skills and removes retired ones. An edited skill needs a second yes, so a hand edit is never lost silently. Read the plan before you confirm it.

If the output has a line such as `the skills here were written by akmctl 0.4.0, newer than this akmctl 0.3.4: install akmctl 0.4.0 or newer`, the repository carries skills from a newer tool than the one on your `PATH`, and some commands they name may be missing here. Do not run `skills update` (it leaves those skills alone, and would otherwise downgrade them). Tell the maintainer which version to install, and until then work from the skill text and the plain-file fallbacks in `AGENTS.md`, not from commands your `akmctl` does not list. `akmctl version` shows what you run; `--format json` shows `requires_tool` on each skill concerned.

Files changed: only inside `.agents/skills/<name>/` of a skill the tool manages. The tool never commits; commit the result yourself.

Errors: `refused` with `details.reason` `no_manifest` when the repository has no manifest. If your `akmctl` has no `skills` command, change nothing and tell the maintainer.
