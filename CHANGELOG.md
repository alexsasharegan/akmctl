# Changelog

## Unreleased

- `new`, `claim`, `release`, `close` and a `set` that changes a work item's `status` return a `commit` object (`step`, `item`, `detail`, `paths`, `message`): the subject `<id>(<step>): <detail>` and the trailers `AKM-Step` and `AKM-Item`, for the caller to commit with. The tool still never stages or commits (ADR 261030-tbkh, proposed). The CLI prints the message and the paths to stage. The API spec, schemas, fixtures (new case `set-8`) and the `akm-commit` skill are updated; `akmctl commits [--range <revs>] [--limit <n>]` is new: it reads `git log` and checks each commit that claims to be an AKM step. CM1 the subject, CM2 the two trailers, CM3 the diff (a `close` moves the file into `done/`, a `claim` leaves it active and claimed, and so on). Read-only, exit 1 when a commit fails; `GET /commits` over HTTP; schema `commits.json`; no fixture cases, since a fixture directory cannot hold a git history (tested against repositories the tests build with `git`).

- `close` returns `dependents`: the open work items (`raw`, `ready`, `active`, `blocked`) whose `blocked_by` names the closed item, for `done` and `dropped` and in a dry run. It only reports and never refuses; the CLI lists them after the move, and a close with none is unchanged. The API spec (section 4.12), `close.json`, the `close` fixtures and the `akm-work-items` and `akm` skills are updated. A `check` rule for a stale blocker was considered and left out (it would bump the layout).

- `new work` returns `similar`: up to five open work items (`raw`, `ready`, `active`, `blocked`) that share at least two words with the slug and title, most shared first. It never refuses, and a response with no match is unchanged. The CLI prints them under the new path. The API spec (section 4.10) states the matching rule, `new.json` has the member, and the `akm-work-items` and `akm` skills tell the agent to add to a match instead of filing a duplicate.

## 0.3.2 - 2026-10-30

- `section`: appending to `## Log` adds a `- ` bullet to a first line that has no list marker, so a bare line no longer joins the previous bullet. The API spec and the skill reference say so.

## 0.3.1 - 2026-10-29

- `ui`: toasts are reworked. They sit bottom right above the hint bar, come in four levels (success, info, warning, error) with a glyph and the level word above the message, and are ephemeral (8 seconds) or persistent (errors stay until dismissed). The top border counts down and pauses on hover and keyboard focus; each toast has a close icon, and Escape dismisses the newest persistent one.

## 0.3.0 - 2026-10-29

- Layout 0.3.0 (the layout, the four extensions and the specification; the API stays 1.0): work items gain three optional fields, ADR 261029-sg82. `closed` is the day the item closed (local day, `as_of` wins); only `close` writes it, for `done` and `dropped`, on the line after `updated`, and `set` refuses it. `cycle` is an opaque label matching `[A-Za-z0-9][A-Za-z0-9._-]{0,63}` and `keep` a boolean; both are set with `set` (`null` removes them), shown by `show`, `ls --fields` and the board, and `ls --cycle <label>` filters by the label. New `check` rules: WK13 (a warning: an item in `done/` has no `closed`), WK14, WK15, WK16 and WK17. Items closed earlier have no `closed` until someone backfills it. The schemas, fixtures (new cases for `close` dropped, `set` of `cycle` and `keep`, and the refused `closed`), the specification, the guide and the skill references are updated. `migrate`'s `move_closed` does not invent a date.

## 0.2.1 - 2026-10-29

- `akmctl version` after `go install` reports the module version (it printed `0.0.0-dev`).

## 0.2.0 - 2026-10-29

- `check`: SF3 no longer flags record file names or Markdown link targets as secrets; L2 reads `[[name\|alias]]` inside a table cell as the note `name`.
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
