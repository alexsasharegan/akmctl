# Changelog

## Unreleased

- Beta 0.2.0 (greenfield): `akmctl init [--docs-dir <name>] [--wikilinks] [--git-init] [--git-commits]` writes the layout, `AGENTS.md`, `CLAUDE.md`, the `akm` skills, the manifest and, in a git repository, a `.gitignore` for the lock into a new directory, never overwriting (API operation `init`, schema and fixtures). The manifest, deny list, lock and migrate log live in the docs dir, which the tool finds on its own. `akmctl config git-commits on|off` installs or removes the optional `akm-commit` skill. The layout, the extensions and the specification are versioned 0.2.0; adopting an existing repository is not supported in the beta.
- New ids and record dates use the local calendar day (`as_of` still wins); `claimed_at` stays UTC. Re-claiming an item you hold adds no log line. `mint_ids` takes `slug_from` and `max_slug`. New `migrate` step `move_closed`.
- `ui`: the page title names the view and repository; the header is centred; an SVG favicon.
- Manifest keys `work_dir`, `adr_dir` and `roadmap` in `docs/.layout` (repository-relative, forward slashes, no `..`) move the work, ADR and roadmap locations; `done/` is `<work_dir>/done`. Loading, `check`, `index`, `new`, `set`, `section`, `close`, `claim`, `release`, `migrate` and the workbench follow them. New rule S3 (error) reports a bad value; writes are refused while it holds. `version.repository` carries `work_dir` and `adr_dir`, which the workbench reads instead of hard-coding the paths. The manifest and `docs/.migrate-log` stay under `docs/`.
- `migrate` steps: `title_from_h1`, `goal_scaffold`, `remove_field`, `copy_field`; `rename_field`, `set_default` and the new field steps take `dir` and then reach notes that have no id yet. `mint_ids` takes `rename: false` (write the `id`, leave the names, report `renames: [{from, to}]` on the step) and `seed` (ids repeatable between plan and apply). `plan` now applies each step on a scratch copy of the Markdown files, so a step that depends on an earlier one is planned as it will run.
- `migrate --migration <file>` and the API parameter `migrations` run a one-off declarative migration for a repository with no layout history; the safeguards are unchanged.
- `close` and `section` leave a blank line after the heading of a section they create (a Prettier-clean `## Outcome`).
- `migrate`: the declarative step `mint_ids` (`dir`, `date_field`, `recursive`) gives a record with no `id` a minted id (date from a field, unique in the batch and the repository), renames the file to `<id>-<slug>.md`, writes `id` first, and reports references to the old name without rewriting them. A migrate step object may carry `notes`.
- API 1.0 (edited in place): `section` writes one `##` section; `close` takes `outcome`; `dry_run` returns a `diff` for `set`, `section`, `close`, `claim` and `release`; `version` carries a `repository` object; `details.etag` is documented.
- `ui`: the workbench is a client of the API: board with vim keys, tables, command palette, record drawer, moves with undo, a wizard and a review diff for gated moves, capture, live updates by polling, claim controls, open in Obsidian (by absolute path), the Quiet paper theme and bundled Monaspace fonts. `GET /ui/config` gives the root and the identity.
- `new work` leaves a blank line after the `## Log` heading, so the file passes `prettier --check` as created.

## 0.1.0

First release. Everything in the AKM API 1.0 draft is built.

- Records: `ls`, `search`, `show`, `board`, `find`, `when`; `new`, `set`, `close`, `claim`, `release` with compare-and-swap writes under an advisory lock.
- `check`: the work item, ADR, roadmap, `AGENTS.md`, link, manifest and skill rules, the claims rules, and the safety extension (secret patterns, a deny list, allow markers).
- `index`: the ADR routing table, byte-stable under Prettier 3.9.10.
- `migrate`: report, plan, apply_safe (backup tag and branch, declarative steps, a log) and handoff. No migrations are embedded yet: layout 1.0.0 is the first.
- `ui`: an HTTP binding on the loopback interface with a per-run token, and a Kanban workbench page.
- A YAML-lite parser and byte-splice editor; Prettier's width rule implemented from Unicode 17.0.0 data.
