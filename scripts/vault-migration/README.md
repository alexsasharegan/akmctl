# AKM migration scripts (vault, dry run of 2026-10-28)

Used by the dry run on a copy of the vault; the real run uses the same files. Python 3 standard library and bash only. Runbook: `Deep Thought/9 pkm/design/vault-migration-runbook.md`; drafts: `vault-migration-drafts.md`.

| File | Does |
| --- | --- |
| `vault-0.0.0-to-1.0.0.json` | The one-off migration for `akmctl migrate --migration`. A sample: regenerate it with `make_migration.py` on the day (it reads task statuses to emit the `done/` moves). |
| `make_migration.py` | Writes that JSON. |
| `adr_slugs.py` | Replaces the ADR slugs in the renames list with title-based ones. |
| `obsidian_renames.py` | The Obsidian CLI rename loop (real run only; untested against Obsidian). |
| `link_repair.py` | Emulates Obsidian's link repair, then renames (dry run stand-in for the loop). |
| `adr_prose.py` | Rewrites bare `ADR-019` mentions to `ADR <id>`; dry run by default. |
| `apply_summaries.py`, `adr_summaries.json` | Writes the drafted ADR summaries with `akmctl set`. |
| `verify.py` | Checks (a) link multiset, (b) masked diff, (c) no unresolved link added, (d) Shavian files, (e) calendar, resources, archives. Paths and counts only. |
| `fold_skill_frontmatter.py` | Folds block scalars in SKILL.md frontmatter (rule S2). |
| `goal_review.py` | Lists open tasks with their scaffolded Goal and flags. |
| `dependents.patch` | The dependent edits (AGENTS.md, skills, Tasks.md, registry, conformance, templates, runbook), made on the copy. `git apply --check` it first. |
| `rehearse.sh` | Runs the scripted steps on a fresh copy and prints step timings. Never point it at the real repository as destination. |
