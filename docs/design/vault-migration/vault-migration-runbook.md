---
type: design
status: draft
created_at: 2026-10-28
updated_at: 2026-10-28
authored_by: agent
tags:
    - pkm/rails
---
# Vault migration: dry-run findings and runbook

Dry run of [[261027-nrwf-vault-agent-ready-migration]] by Claude Code on a copy of the vault. The real vault was not modified. Drafts for review: [[vault-migration-drafts]]. Scripts: `Deep Thought/9 scripts/akm-migration/`.

**Resume point.** The copy is `/private/tmp/claude-501/-Users-deepthought-obsidian-vaults/dfafa885-dc8d-435b-9a5a-ca08fd51f7df/scratchpad/copy` (branch `migrate-0.0.0-to-1.0.0`, tag `pre-migrate-0.0.0-to-1.0.0`); `scratchpad/rehearsal` is a second copy from `rehearse.sh` with only the scripted steps. All steps below were run on the copy except the Obsidian rename loop (emulated by `link_repair.py`).

## Result

- Final `akmctl check` on the copy: 0 errors, 47 warnings: RM4 x45 (every open top-level task is not on the roadmap) and L2 x2 (existing: a PNG link and an ambiguous `Common Place Patristics`; not caused by the migration).
- Conformance script on the copy: PASS (after the `id` change below).
- Verification of the link repair (stage after the ADR prose rewrite, `verify.py`): (a) link multiset equal: yes (4825 links before; the only differences are in 8 notes edited by hand afterwards, 2 of them the sweep links below); (b) masked diff: 2355 notes strictly equal, 123 task and ADR notes differ only in frontmatter and the `## Goal` heading: 0 bad; (c) unresolved links 1600 before and after, none added; (d) 69 files with Shavian text changed, 307 links, link spans only (36 task notes, 4 ADRs, 12 other `9 pkm` notes, 8 calendar notes, 1 in `2 PROJECTS`, 3 in `3 AREAS`, 1 in `5 ARCHIVES`, 4 skills); (e) calendar: 8 files (14 links): `2026-05-31`, `2026-09-12`, `2026-10-24`, `2026-10-25`, `2026-09-Monthly-Retro`, `2026-W28`, `2026-W37`, `2026-W43`; `5 ARCHIVES`: `Bullet Journal Blog - Candidate Articles` (1 link); `4 RESOURCES`: none. 1 file under `3 AREAS/AMEX` had one link rewritten (path withheld).
- A mutation test (one link target and one heading changed) made (a) and (b) fail, so the checks bite.
- Overall: 145 renames (124 tasks, 21 ADRs), 821 links in 188 files rewritten, ids deterministic (a second copy gave identical names).

## Step timings (copy on this Mac)

| Step | Time |
| --- | --- |
| copy of the repo (no `.bak`, `.tmp`, `node_modules`) | 4 s |
| migrate report / plan / apply-safe | 0.1 s / 0.8 s / 0.5 s |
| renames list and ADR slugs | 0.1 s |
| link repair and renames (script, stands in for Obsidian) | 1.5 s |
| ADR prose rewrite | 0.7 s |
| 21 summaries with `akmctl set` | 0.7 s |
| `akmctl index`, `akmctl check` | 0.07 s, 0.09 s |
| `verify.py` | 10.7 s |
| conformance script | 0.4 s |
| whole scripted pipeline | 21 s |

Not measured: the 145 Obsidian renames (the CLI cannot run on a copy; the pilot renamed one note). Hand edits (dependents, drafts) took most of the session.

## What changed in the dependents (on the copy; `dependents.patch`)

- `docs/.layout` (new): `work_dir`, `adr_dir`, `roadmap: Deep Thought/9 pkm/design/akm-roadmap.md` (new empty roadmap with the four headings and frontmatter); the tool adds `layout: 1.0.0`.
- `AGENTS.md`: 135 lines to 42, five parts; old text moved verbatim to the new `9 pkm/vault-conventions.md`.
- `Tasks.base`: `updated_at` to `updated` (7). `Tasks.md`: naming, frontmatter block, `done/`, `akmctl close`, examples. `frontmatter-fields.md`: `id` is text; `created`/`updated`, `summary`, `adr_number` registered. `conformance.py`: `KNOWN` for `id` (open, number and text) and `created`; the `NNN` placeholder entry removed. `TMPL - ADR.md`: AKM frontmatter.
- `task-lifecycle` skill (1.3.0) and `adr-check.md`, `monthly-review` SKILL (path of the conventions), `skills-map.md`, `runbooks/Agent Hooks.md` (one bullet). Eight SKILL.md files: block scalars folded to one line (rule S2).
- `design/adr-conventions.md` (moved from `adr/README.md`, rewritten for the AKM fields and statuses; Index table replaced by a pointer), `ADR Candidates.md` moved beside it, two sweep links repointed to `[[adr-conventions]]`.
- `.obsidian/types.json` (untracked, gitignored): `id` text; added `created`, `updated` (date), `summary` (text), `adr_number` (number). Keep a copy before editing.
- Hooks: `session-start.sh` and `stop.sh` need no change (they read `tasks/*.md` and names); run by hand on the copy: the open-task list printed the new names, and a faked transcript made `stop.sh` block for the active task. `user-prompt-submit.sh` and `pre-compact.sh` do not mention task paths.
- Not touched: `pkm-workbench` (does it read `tasks/done/` and `created`/`updated`? not checked).

## Findings, ranked

1. **`id` collision.** 121 Field Notes use a numeric `id`; AKM ids are text. Retyping `id` to text in `types.json` and marking the collision known in the registry and the conformance script makes the script pass. Alternative: rename the Field Notes field. Alexander to choose.
2. **The `ADR-NNN` mentions are not all vault ADRs.** Of 94 mentions: 40 rewritten; 40 sit in 14 files that contain Shavian text (skipped; includes `LifeNotes_ADR` and the Tagging Manifest, which may use their own numbering); 4 in 2 calendar notes (not touched); 10 are in `3 AREAS/AMEX`, a different ADR series (excluded, must not be rewritten). 9 more mentions are numbers with no vault ADR (ADR-042 and others). 10 occur in code spans or link text and were left.
3. **ADR slugs.** `mint_ids` builds `adr-019` from the old file name; `adr_slugs.py` rewrites them from the titles (40 characters at most, trailing stop words dropped, e.g. `261027-ensf-akm-an-agent-knowledge-management`). Some are cut mid-phrase. Alexander to approve the 21 names (the renames list in the dry-run copy).
4. **Two date pairs.** ADRs and tasks now use `created`/`updated`; every other note keeps `created_at`/`updated_at`. ADR `date` is a copy of the old `created_at`. Decide whether ADRs should keep `created_at`.
5. **RM4 x45.** The roadmap is empty. Fill it (a hand-ordered list of ids) or accept 45 warnings.
6. **Skills.** Rule S2 rejects the block scalars in 8 SKILL.md files; folded in the dry run. Alternative: change the tool (work item 3).
7. **Conventions note.** The tool's `move_file` does not repair links; two `[[9 pkm/adr/README#Sweep watermark|...]]` links broke and were edited by hand. The alias text still says "the README's sweep section". Obsidian's `move` would repair them (not piloted).
8. **`docs/.akm.lock`** is created by the tool and was committed by `git add -A`; add `/docs/.akm.lock` to `.gitignore` first. `docs/.migrate-log` is 106 KB (every file of every step).
9. Statuses of old ADR vocabulary (`active`, `error`) are retired in `adr-conventions.md`; no ADR used them.
10. Nothing in the prose of the ADR files is edited except mentions; ADR H1s read `# ADR 261027-ensf: Title`.
11. The Next draft was needed for one task only, and 124 (not 115) tasks got a Goal scaffold.

## Proposed `akmctl` work items

1. `mint_ids`: option to build the slug from the `title` (and a maximum length), for notes whose file name is a number.
2. `migrate --mode plan` does not apply `move_file` steps, so it counted `ADR Candidates` as an ADR (22 ids, 22 renames); apply made 21.
3. Rule S2: accept YAML block scalars (`>-`) in SKILL.md frontmatter, or parse skills with a separate rule.
4. A `move_file` variant by status (or a `done/` step) instead of 76 generated steps; and `docs/.migrate-log` without full file lists.
5. The lock file: keep it out of the working tree (or write the ignore line).
6. `mint_ids rename:false` could write the renames to a file (`--renames-out`); today they are read from the JSON response.
7. Optional: a `move_file` that also reports inbound links it leaves broken.

## Runbook for the real run

Do not start before Alexander has answered the decisions above (findings 1 to 5, the 14 Shavian files, the 21 slugs, the roadmap). Obsidian must be running. Do not type in the vault during steps 6 to 9. `H=/Users/deepthought/obsidian.vaults`; `A="Deep Thought/9 scripts/akm-migration"`; work files go in `~/akm-run/` (outside the repo).

1. **Window.** `launchctl list | grep autosave` (leave the job alone). It fires at 09:00, 13:00, 17:00, 21:00, 23:50 and does `git add .` then `git commit` on the current branch. Start five minutes after a run and finish well before the next (the scripted part takes under a minute; budget 45 minutes for the renames and checks). A run that fires mid-way only commits the half state on the migration branch, which the tag undoes.
2. **Clean tree.** `cd $H && NOTES=$H git-autosave && git status --porcelain` must print nothing (the tool refuses a dirty tree). `git tag pre-akm-migration`. `mkdir -p ~/akm-run && cp "Deep Thought/.obsidian/types.json" ~/akm-run/types.json.before`.
3. **Prepare on main.** Add `/docs/.akm.lock` to `.gitignore`; write `docs/.layout` (three keys, as above) and `Deep Thought/9 pkm/design/akm-roadmap.md`; `python3 "$A/make_migration.py" . > "$A/vault-0.0.0-to-1.0.0.json"`; `git add -- .gitignore docs/.layout "Deep Thought/9 pkm/design/akm-roadmap.md" "$A"` and commit. Build the tool: `cd ~/github.com/alexsasharegan/akmctl && go build -o ~/akm-run/akmctl ./cmd/akmctl`.
4. **Plan.** `~/akm-run/akmctl migrate --mode report --migration "$A/vault-0.0.0-to-1.0.0.json"`, then `--mode plan`. Expect the three errors classes only: S2, G1, A4 (and S1). `plan` counts 22 ADRs; ignore.
5. **Apply.** `~/akm-run/akmctl migrate --mode apply-safe --acknowledge-backup --migration "$A/vault-0.0.0-to-1.0.0.json" --format json > ~/akm-run/apply.json`. It makes tag `pre-migrate-0.0.0-to-1.0.0` and branch `migrate-0.0.0-to-1.0.0`, moves the two conventions notes and 76 closed tasks, edits the frontmatter, writes ids without renaming. Then extract the renames and set ADR slugs (snippet in `rehearse.sh`) into `~/akm-run/renames.json`; `python3 "$A/adr_slugs.py" --root . --adr-dir "Deep Thought/9 pkm/adr" --renames ~/akm-run/renames.json`. Check `jq` count is 145 (`124` tasks, `21` ADRs). Commit by path: `git add -- "Deep Thought/9 pkm" docs` and commit "AKM migration: tool steps".
6. **Renames through Obsidian.** `python3 "$A/obsidian_renames.py" --renames ~/akm-run/renames.json --root . --dry-run | head`, then without `--dry-run --log ~/akm-run/renames.log`. Each call is `obsidian --vault="Deep Thought" rename file="<old name>" name="<new name>.md"`; it stops at the first failure, a rerun resumes. Then `git status --short | wc -l` (about 205 changed or renamed files) and commit "AKM migration: Obsidian renames".
7. **Verify links.** `python3 "$A/verify.py" --repo . --before pre-migrate-0.0.0-to-1.0.0 --after HEAD --migration "$A/vault-0.0.0-to-1.0.0.json" --renames ~/akm-run/renames.json --tool-dirs "Deep Thought/9 pkm/tasks/" "Deep Thought/9 pkm/adr/" --out ~/akm-run/verify.json`. Expect (a) false by exactly 2 (the sweep links), (b) and (c) true, Shavian files about 67, calendar 8. Anything else: stop and compare with `git diff pre-migrate-0.0.0-to-1.0.0` before going on.
8. **ADR mentions.** `python3 "$A/adr_prose.py" --root . --adr-dir "Deep Thought/9 pkm/adr" --include "Deep Thought/9 pkm" --include "Deep Thought/2 PROJECTS" --include "Deep Thought/3 AREAS" --include "Deep Thought/0 Inbox" --include .agents --exclude "Deep Thought/3 AREAS/AMEX" --report ~/akm-run/prose.json` (dry); with `--apply` once it reads 40 (plus the Shavian files if Alexander allows). Commit.
9. **Dependents.** `git apply --check "$A/dependents.patch" && git apply "$A/dependents.patch"`; edit `types.json` (`id` text; add `created`, `updated`, `summary`, `adr_number`); `python3 "$A/fold_skill_frontmatter.py" .agents/skills/*/SKILL.md` is already in the patch. Apply the summaries: `python3 "$A/apply_summaries.py" ~/akm-run/akmctl . "Deep Thought/9 pkm/adr" "$A/adr_summaries.json"` (after review), the `## Next` with `akmctl section <id> Next --text "..."`, then `~/akm-run/akmctl index`. Commit.
10. **Checks.** `~/akm-run/akmctl check` (0 errors; 45 RM4 and 2 L2 warnings), `python3 "Deep Thought/9 scripts/frontmatter/conformance.py"` (PASS), the two hooks by hand as in the dry run, open `Tasks` in Obsidian and look at the board.
11. **Merge.** `git switch main && git merge --ff-only migrate-0.0.0-to-1.0.0`. If an autosave commit landed on main meanwhile, use a normal merge. Log the run in the task note and the decision log (Y43).
12. **Roll back** (before the merge): `git switch main` and `git branch -D migrate-0.0.0-to-1.0.0`; restore `types.json` from `~/akm-run/types.json.before`. After the merge only with Alexander's yes: `git reset --hard pre-akm-migration`.

## Not verified

The Obsidian CLI rename loop at scale (timing, whether `rename` accepts these names, how it treats table `\|` and frontmatter links in all 188 files beyond the one pilot); `obsidian move` for the conventions notes; `pkm-workbench` and `akmctl ui` against `done/`; Prettier (not run); the hooks inside a real Claude Code session; that G3 sees the conformance command (it does not); the ADR summaries and Next text beyond the Decision sections I read.
