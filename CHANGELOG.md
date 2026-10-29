# Changelog

## Unreleased

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
