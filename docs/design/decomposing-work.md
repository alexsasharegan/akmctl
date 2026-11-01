---
type: design
status: draft
created_at: 2026-11-01
updated_at: 2026-11-01
authored_by: agent
---

# Decomposing a work item into a parent and children

Draft for spike 261101-baze. The procedure below is written to drop into a skill (or a section of `akm-track-work`) unchanged; the walk-through and the open questions follow it.

## Procedure (skill-ready text)

### When to split

Split an item when any of these holds:

- Its "done when" list holds independent deliverables, so one could be finished and shipped without the others.
- It is a spike with more than one question, or a spike whose answer needs another spike first.
- It will outlast one sitting and different people or agents could take its parts.

Do not split to tidy a long checklist whose boxes are steps of one change. One change is one item.

### Shape

- The **parent** keeps the why and the umbrella Goal. It holds the ordered list of children as links and a "done when" of its own: every child is `done` or `dropped`, plus anything only the whole can show (the decisions are recorded, the follow-up items are filed). It is never worked directly.
- Each **child** owns one question or one deliverable with its own "done when" and a `parent` field. A child is a spike or a build item; a spike's last box files the build item.
- Use `blocked_by` only between siblings that truly depend on each other. A parent already waits for its children, so never list a descendant in the parent's `blocked_by` (WK19).
- A child may itself be split. Keep the tree as shallow as the work allows.

### Steps

1. Read `ls --parent <id>` first: the item may already have children.
2. Settle the shape on paper: the parent's umbrella Goal and one line per child (title, kind, what it answers or delivers, who it waits for).
3. File each child with `akmctl new work <slug>`, then `akmctl set <child> kind=<kind> parent=<parent> summary="..."` and `akmctl section <child> Goal --text "..."`. Heed the `similar` list `new` prints: a child may already exist.
4. Rewrite the parent's `## Goal` with `akmctl section`: umbrella text, the children in dependency order as links, and the parent-level "done when". Keep the parent's old `## Log`; add a line naming the children (`akmctl log`). Move the old Goal's boxes into the child that owns them; do not delete a box without a trace (`remove` needs a reason).
5. Release the claim on the parent (`akmctl release`) unless you are about to work a child; claim the child you start. Set the parent's `status` to `ready`; `active` and `blocked` need a `## Next`, so write one that points at the first child.
6. Only the parent goes on the roadmap (RM4 asks for top-level items only). Children follow their parent.
7. `akmctl check`, then commit one step at a time: the children, then the parent's rewrite.
8. Close in order: children first, the parent last, with an Outcome that links the children's outcomes.

## Walk-through on 261030-456q (test strategy and coverage audit)

Dry run on a real item; nothing was changed. The item is a `ready` spike with three parts (strategy, maintainability, coverage), six "done when" boxes and a closing step that files tickets by group. It is too big for one spike: the report, the measurements and the filed tickets are separate deliverables.

Result of applying the procedure:

- Parent `456q` keeps the seed list and the framing. Its boxes become: every child is `done`, the seed list is checked entry by entry, the roadmap lists the filed items.
- Children: a spike for strategy, one for maintainability, one for coverage (each adds its section to `docs/design/test-strategy-audit.md` and ends in numbered findings), then a last child that groups the findings and files the tickets, blocked by the other three.
- The old boxes map one to one: the report and its findings go to the three section spikes, the recorded commands and the grouping go to the coverage spike and the ticket child, the seed check stays with the parent.

What the walk-through showed:

1. The split is mostly `new`, `set` and `section` repeated: four children cost twelve commands and twelve names to remember. `split` (261031-ph53) removes most of it; it should take a title and a `--kind` per child.
2. Moving a "done when" box from a parent to a child has no command. `section` rewrites the whole Goal, so the text is retyped. This loses nothing but invites typos.
3. A `ready` item already carries `priority: 3`. The parent keeps it; the children inherit nothing, so ordering by priority puts them below unrelated work unless set by hand.
4. A claimed item becomes a parent: the claim has to be released, and nothing says so today.

## Open questions

- Does the cost of a new skill fit? A new skill file is mirrored in `internal/initrepo/templates/skills` and listed in six `init` fixtures (two copies of each), `make_fixtures.py`, `initrepo_test.go` and `skills_test.go`. A section in `akm-track-work` touches only that skill and its mirror.
- Should `kind` stay `feature | bug | chore | spike` for a parent? 261031-gtw6, the existing hand-built tree, uses `chore`. The glossary says there is no epic record, which this procedure keeps.
