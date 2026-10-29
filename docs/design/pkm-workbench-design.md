---
type: design
status: active
created_at: 2026-10-25
updated_at: 2026-10-25
authored_by: agent
---
# pkm-workbench: design for a local PKM server

Result of the design spike 261024-k2rc-local-views-server-design, under Vault Views. Questions 2 to 9 were worked through with Alexander on 2026-10-25 (Plan mode); question 1 (the Bases baseline) was answered the same day.

**Recommendation: build it.** A small Go program, `pkm-workbench` (short forms `pkm-wb`, `pkmwb`), serves the task views from the vault's files and writes changes back to the YAML frontmatter. The files stay the source of truth. Whether to build is Alexander's call; this note is the case for it.

## Purpose

- The views in Tasks (`Tasks.base`) live only inside Obsidian. The wish: the same views in other forms with different styling, and the ability to move things (drag a task across a board, close it) with the change written back to the file.
- The name is chosen for growth. Version 1 is the task board; later it can serve other record types and views of the PKM. "Workbench" names where Alexander works, not what it shows.

Direction, Alexander 2026-10-25, after the Kanban test: the Tasks hub dashboard now leads with the native Kanban board, and he still wants the server. Reasons: it should be prettier and more pleasant to use, and a page left open lets him watch an agent work.
- **Watch mode (wanted, after the MVP):** the server loads the board state; when an agent changes a task, the page animates the move (this task went from this column to that one) when he returns to the page, so changes do not pass unseen. Needs change tracking between visits, which the read-per-request design does not have yet, and probably the file watcher that was left out of version 1.
- **Lifecycle enforcement (the main value):** when he moves tasks himself, for example with no tokens or credits left, the server enforces the bookkeeping an agent would do: an outcome on close, log lines, `updated_at`, the ADR check. He will not remember every step of the task-lifecycle skill; the application should refuse the move until the steps are done. This makes the close prompt and the pinned close cascade and ADR dialog core features, not extras.

## Baseline: what Obsidian Bases gives (question 1)

- On 2026-10-25 three native views were added to `Tasks.base` (Obsidian 1.14.4): Cards by status, Cards by project, List next and active. Obsidian parses them. The checkpoint before the change is commit `3008564`.
- Alexander's verdict: useful, possibly better as a mobile view. The look is a no-go and does not scale; he wants a UI with more control.
- Not verified: dragging a card to another group to change its status. Check it before saying Bases has no board.
- Correction, same day: a subagent read the changelog and Bases docs. Bases has had a **Kanban layout since Obsidian 1.14.0**, so this spike's assumption of no board was wrong. The docs say dragging a card between columns writes the grouped property (for example `status`) to the note's frontmatter. Not documented: order within a column, WIP limits, drag on mobile, per-column colours. The 1.14.5 release (22 Oct) is only an Electron update per the changelog. All of this is from documentation, not tested. Tested the same evening on Alexander's phone: he dragged `2026-10-20 appendix-b-skills` from inbox to next in the Kanban view; the note's `status` changed in the frontmatter (a one-line diff, verified with git), and `updated_at` did not change. So drag-to-change-status works natively on mobile, but it does not set `updated_at`. It is Alexander's call whether it changes the build recommendation; it may, since drag-to-change-status was one of the two gaps (the other is styling and control).
- The three views stay in `Tasks.base` for now.

## Decisions ledger

| Q | Decision | Options considered | Why |
| --- | --- | --- | --- |
| 2 Watch scope | Tasks only: notes with `type: task` in `9 pkm/tasks/`. Read the folder on every request. No file watcher in version 1. | Watch the whole vault; watch with FSEvents and push live | The task schema is the only one with a stable status vocabulary. About 100 files read in milliseconds. A watcher is added if live push is wanted. |
| 3 Reading views | The server defines its own views from the spec below. `Tasks.base` stays for Obsidian. Both implement the same spec. | Parse `.base` files; reuse the Bases expression language | Re-implementing Bases formulas (`today() - updated_at`, `file.inFolder`) is a trap. The cost is two implementations, kept honest by the spec. |
| 4 Write-back | Edit the one frontmatter line in place. Compare a SHA-256 of the file (SHA-1 is fine: change detection only) against the hash the page was rendered from, and refuse on a mismatch. Write a temp file in the same folder and rename it over the original. | Parse and re-serialise YAML | A YAML library reorders keys, changes quoting such as `"[[Vault Views]]"` and drops comments. Task frontmatter is flat, so a line edit is safe. |
| 5 Write path | Direct file writes. One write path for the whole system. | `obsidian property:set` | The CLI needs Obsidian running. The server should work with Obsidian closed. |
| 6 Security | Bind `127.0.0.1` only. A per-run random token in a header on reads and writes. Check the `Origin` and `Host` headers. Phone access is out of scope. | LAN or Tailscale bind with real auth | A page in a browser can send requests to localhost, so Origin and token checks are required. Notes can hold private detail. |
| 7 Rendering | The server has its own templating layer (`html/template`, embedded CSS and JS) and shares no code or CSS with Vault Views. The UI is designed from scratch by an agent with no context on the kit. | Reuse the kit and the Vault Views build | A board needs live JS for dragging and the close prompt. The kit's pages are static documents. Alexander does not like the kit's design. This supersedes the task's original question 7, "reuse whatever html-views-poc produces". |
| 8 Stack | Go, standard library only (`net/http` with the 1.22+ `ServeMux` path patterns, `html/template`, `embed.FS`). Vanilla JS for drag and drop. A dependency comes in only when it earns its place, with the reason logged. | Node, Python, Rust | A single static binary suits a local daemon under `launchd`, and `git-autosave` is already a Go binary. |
| 9 Inbox lane | A new status `raw` for unrefined captures, shown as its own lane. Refining writes `project`, `origin` and a done-when line, then moves the note to `inbox` or a later status. | Keep raw notes as `inbox` notes missing the refine fields | Explicit, and the filter is trivial. `inbox` keeps its meaning: filed, not claimed. |
| Closing | A drag to `done` or `dropped` opens a prompt that requires an outcome. The server sets `status` and `updated_at`. | Allow any status change with no outcome | A card must not leave the board without an outcome. |
| Repo | One new repo of Alexander's own: `~/github.com/alexsasharegan/pkm-workbench`, Go module `github.com/alexsasharegan/pkm-workbench`. Separate from the vault and from the Vault Views fork. | One repo with the fork; the vault's `9 scripts/` | Real code with tests and versions. The server draws ideas from the kit but copies nothing. |

## Architecture

- One Go binary. It reads `9 pkm/tasks/*.md` per request, parses the frontmatter (a small hand parser, or `yaml.v3` read-only if the reason is logged), applies a view from the spec, and renders a page from embedded templates.
- Every page and API route is under `127.0.0.1` with the token and Origin checks above.
- Routes use the standard library's method and path patterns, for example `GET /views/{name}` and `POST /tasks/{id}/status`.
- Config (vault path, port) comes from flags or environment. Views come from the spec, not from `Tasks.base`.

## View spec (authoritative)

The two implementations, `Tasks.base` and the server, both implement this table. Change the spec first, then both. Version 0 lives here; move it into the Tasks hub when the repo exists, so one place defines scope, statuses and views.

Shared definitions:
- **open**: `status` is not `done` and not `dropped`.
- **age_days**: today minus `updated_at`, in days.
- **stale**: open, and `age_days` greater than 30 (or no `updated_at`).
- Base filter: `type == "task"` and the note is in `9 pkm/tasks/`.

| View | Filter | Group | Sort | Columns |
| --- | --- | --- | --- | --- |
| Open by project | open | none | `project` ascending, `status` ascending | project, name, status, age_days |
| Ready and active | `status` is `ready` or `active` | none | none | name, status, project, `updated_at` |
| Blocked | `status` is `blocked` | none | none | name, project, `updated_at` |
| Stale | stale | none | none | name, status, project, age_days |
| Recently closed | `status` is `done` or `dropped`, and age_days is 14 or less | none | none | name, status, project, `adr`, `updated_at` |
| Cards by status | open | `status` | `updated_at` descending | name, project, age_days |
| Cards by project | open | `project` | `status` ascending | name, status, age_days |
| List ready and active | `status` is `ready` or `active` | `project` | none | name, status |

The server adds one view `Bases` lacks: **Raw**, the notes with `status: raw`.

## Write path

- **One-line edit.** Find the key inside the frontmatter block, replace its value, leave every other byte alone. If the key is missing, insert it before the closing `---`. The server writes only known scalar fields: `status`, `updated_at`, later `priority`.
- **Conflict check.** The page carries the hash of the file it was rendered from. A write re-reads the file, hashes it, and refuses on a mismatch, showing the change that got in the way.
- **Atomic.** Temp file in the same folder, then rename.
- **Close is the one structured edit.** It puts `## Outcome` above `## Log`, deletes the `## Next` block, adds a closing log line, and sets `status` and `updated_at`, as one write under the same hash check.
- **Obsidian open with unsaved edits.** Obsidian reloads a changed file if its editor is clean. If the note has unsaved edits, the write can race with Obsidian's own save, and the server cannot see that. Accepted for a single-user tool: the hash check covers the common case, and git is the undo.
- **Autosave.** `git-autosave` commits whatever is on disk, so a server write lands in the next autosave commit. No conflict.

## Inbox lane

- Capture: a text box on the Raw lane creates a stub note in `9 pkm/tasks/` with `status: raw`.
- Refine (an agent or Alexander): write a done-when line (required) and `project` and `origin` (asked, may be skipped), then set `ready`.
- Raw notes stay out of every other view, so they do not pollute the queue.
- Updated 2026-10-26: the statuses are now `raw`, `ready`, `active`, `blocked`, `done`, `dropped` (definitions in Tasks, which is the single source for the wording). `inbox` and `next` are gone; the Raw lane is always shown. See 261026-p5cz-task-status-vocabulary.

## UI approach

- Designed from scratch, not from the kit. An agent with no context on the kit gets a brief that covers what the screens must do and the constraints. It does not get the kit's look.
- Screens: the board by status, the project view, the Raw lane, the close prompt.
- Constraints: light and dark; dense but readable; keyboard use; works at laptop and tablet widths.
- The close prompt is where the outcome is entered. The page never closes a card without it.

## Moves are transactions (decided 2026-10-25)

A move between columns is a transaction, not an edit. Dragging a card opens it; nothing is written until the user confirms.

1. **Check.** The rules table (below) says what the target status requires. If the task already meets it, the move commits with no modal.
2. **Wizard for what is missing.** Only unmet requirements become steps: one question per step, "2 of 3" progress, Enter to advance, Back keeps answers. The server fills in everything mechanical (`updated_at`, the "claimed" and "closed" log lines, deleting `## Next`); the user supplies judgment only.
3. **Review.** The last step shows the exact change to the note as a standard diff view of the raw Markdown (unified diff, added and removed lines coloured). The user accepts or cancels. Cancel snaps the card back.
4. **Commit, then animate.** The write goes through the usual hash check and atomic rename; then the move animates. A move the user just committed takes priority over any queued background animation.

Rules table (gated moves approved by Alexander; this is part of the view spec):

| Move to | Required (steps shown only if unmet) |
| --- | --- |
| `active` (claim) | Nothing typed. One confirm, with a prefilled "claimed" log line and a `## Next` stub. Not gated. |
| `blocked` | Reason; who or what it waits on. Gated. |
| `done` | Outcome (1 to 3 lines); the ADR check (five toggles under one "None of these" button; any toggle adds a `Why:` step); only if open tasks link to this one, the cascade step. Gated. |
| `dropped` | Reason. Gated. |
| `raw` to `ready` | Done-when line (required); project and origin (asked, may be skipped); each skipped if present. Gated for `raw` only. |
| any other move | No gate. |

Diff view: implement it in the Go standard library (a small line diff and an HTML renderer), not a library, in keeping with the stdlib-only rule. Revisit if it proves poor.

## Watch mode and the animation queue (decided 2026-10-25, after the MVP)

- The server keeps an ordered log of changes since the user's last visit (a change is a task whose status or fields differ from the last state it served).
- On return to the page the changes **play one at a time, in order**, never together: if an agent touched three tasks, the first animates, finishes, then the second, then the third. The queue is FIFO; a user-committed move jumps to the front.
- Needs change tracking plus the file watcher that version 1 leaves out.

## Keyboard first (decided 2026-10-25)

The server should be driven from the keyboard like a terminal UI, and every action discoverable on screen. This, with the rule gating, is the reason to use it over the native Obsidian Kanban board.

- **Navigation:** `h` `j` `k` `l` and the arrows; `g g` and `G` for first and last.
- **Moves:** `H` and `L` (shift) move the focused card one column left or right; `d` closes (done), `x` drops. A move opens the transaction (wizard if needed, then the review diff).
- **Command palette:** `Cmd+K` or `Ctrl+K`: type to filter every action and every task, each row showing its shortcut (the pattern Raycast, Linear and GitHub's palette share; "cmdk"-style components implement it in React, here it is hand-written).
- **Discoverability:** shortcuts are shown beside the action they trigger, in a footer bar that follows the context, and `?` lists them all.
- **Wizard:** `Ctrl+Enter` or `Cmd+Enter` next step, `Esc` cancel or back, `Enter` on the review writes.
- **Open the note:** `o` opens it in Obsidian through the `obsidian://` link.
- **Mouse** drag and click still work; the keyboard is primary, not the only path.

Styling is driven by `docs/style-brief.md` in the repo: three directions are being mocked up from the brief by an agent with no knowledge of the Vault Views kit, to pick from.

## Fast follow: animations, a brainstorm (2026-10-25)

Not in the MVP, but expected to be high value: the page you leave open and watch while agents move tasks. Alexander's wish: playful, with choices, possibly configurable. Captured as ideas, none decided.

- **What a theme receives.** Animation is presentation only. The server emits a move event (task, from column, to column, kind: agent or user, time) and the active theme draws it. The queue plays events one at a time (see above). A theme is a self-contained CSS and JS bundle embedded in the binary, so adding one never touches the server logic.
- **Ideas from Alexander:**
  - Vegas or "power mode": the lights dim, coloured spotlights sweep, confetti or small explosions on a close.
  - A rocket ship that turns into a car and drives to the next column.
  - Super Mario World: a little figure scoots along the lane from one column to the next.
  - Frogger: the card hops across lanes, dodging traffic.
- **More to consider:** Pac-Man eating the card and spitting it out in the new column; a Tetris-style drop; a train carrying cards between stations; a marble run; confetti only for `done`; a quiet "fizzle" for `dropped`; streaks (several moves in a row build a combo).
- **Controls.** A theme picker; a speed control (a slow "popcorn" watch speed); an off switch. Quiet by default, with `prefers-reduced-motion` honoured. A busy page can be hard on a sensory-sensitive user, so power themes are opt-in. Sound, if ever, is off by default.
- **Plain default.** Without a theme the card slides to its new column and briefly highlights, which costs almost nothing and is enough to follow what happened.

## Later (pinned)

- **Close cascade.** On close, list the open tasks that link to the closing one (a search of the tasks folder for its name) and offer, per task, to close it with its own outcome or to log that the linked task closed. The search is programmatic; the decision to close each one stays with Alexander.
- **ADR check dialog.** Show the five questions of the ADR check (`.agents/skills/task-lifecycle/adr-check.md`) as checkboxes in the close prompt. None checked logs `ADR check: none`. Any checked sets `adr: candidate` and requires a one-line `Why:` in the outcome. The server records answers and never invents them. Until then a closed task with `adr` unset shows "ADR check pending".
- Live updates with a file watcher.
- Other record types and views (projects, HSA claims, field notes).
- Phone access: needs a LAN or Tailscale bind and real auth. Static snapshots from Vault Views already work on the phone.
- `priority` as an editable field (drag to raise).
- **Override for a gated move** ("move anyway"). Not in the MVP. Design: opt-in, with a visible cost. The move writes a flag in the note's frontmatter (name to be chosen and added to the field registry) and a log line saying what was skipped; the card shows as incomplete. Agents read the flag and treat the flagged tasks as technical debt to fix (not backlog), clearing the flag once the missing outcome, log lines or ADR check are supplied.

## Recommendation and first milestone

Build. The Bases baseline shows the native views are a useful mobile view but not the workbench Alexander wants, and the design above is small: one binary, no dependencies, one write path.

First milestone, in the new repo:
1. Read the tasks folder and render the Next and active view, read-only.
2. Add the other views from the spec, and the check that compares the server's lists with `Tasks.base` on the same files.
3. Add the status change by drag, with the hash check and atomic write.
4. Add the move wizard (see Moves are transactions), starting with `done`: required outcome, ADR check, review diff, then commit.

Open at build time: whether `yaml.v3` is worth a dependency for reads, where the spec moves to, and the design brief for the UI agent.
