---
id: 261026-np3c
type: task
status: ready
project: "[[Vault Views]]"
origin: "[[261025-g68j-build-pkm-workbench]]"
created: 2026-10-26
updated: 2026-10-26
authored_by: agent
title: "pkm-workbench: delete a task, with its backlinks listed"
---
# pkm-workbench: delete a task, with its backlinks listed

## Goal

Alexander captured a test raw task in the workbench, then deleted it in Obsidian because no delete exists and `dropped` would leave a note in the folder. He wants a hard delete in the workbench.

Done when:

The method is chosen and recorded; deleting a task from the board lists its backlinks and removes the file only after confirmation; tested on copies; the real vault checked by Alexander.

## Spec (Alexander, 2026-10-26)

- A hard delete, for **any** task, not only raw stubs.
- Always an "are you sure?" dialog.
- If other notes link to the task, the dialog **lists those notes by name**, not a count.
- Backlinks do not block the delete: he can still confirm after seeing the list.
- Dead links are accepted, including in past daily notes (which stay read-only; the delete does not edit them).
- Recovery is git: autosave commits several times a day, so a delete is recoverable unless it falls between commits.

## Open question: how to find backlinks

Alexander wants a deep dive on how Obsidian itself does it, because he worries about getting it wrong. The Obsidian CLI is the most accurate route but needs the app running; a vault scan by the server has no dependency but may miss what Obsidian counts.

1. Document how Obsidian resolves links: wikilinks and Markdown links, aliases, links with headings, blocks or display text, embeds, links in frontmatter properties, case, path versus shortest name, links inside code blocks, and what Bases do. Include what the `obsidian backlinks` command returns and what the metadata cache holds. Check against the installed version and the real vault, not memory.
2. Compare the CLI route with a server-side scan on the real vault: list the differences note by note.
3. Pick a method with Alexander. A hybrid is possible (use the CLI when it answers, fall back to the scan, and say which one was used).

## Then build

Delete route, confirmation dialog listing the backlinks, a key (palette action and a board key), tests on copies of the notes, README line, roadmap update. The delete is a hash-checked removal, like undo. It does not touch other notes.

## Log

- 2026-10-26: filed by Claude Code on the Mac at Alexander's direction. `inbox`, not claimed.
- 2026-10-26: status inbox to ready (vocabulary change, [[261026-p5cz-task-status-vocabulary]]).
