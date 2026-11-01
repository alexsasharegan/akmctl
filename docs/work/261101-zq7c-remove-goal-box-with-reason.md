---
id: "261101-zq7c"
title: remove goal box with reason
summary: "remove deletes one unchecked Goal box when its requirement is gone, logs the removal with the box text and a reason, and needs confirm because it loses history"
status: ready
kind: feature
parent: "261031-gtw6"
created: 2026-11-01
updated: 2026-11-01
---

# remove goal box with reason

## Goal

A box in the Goal is a requirement. When the work changed and a requirement no longer applies, `close done` refuses while its box is empty, and the only way out is to rewrite the whole Goal with `section`, which leaves no trace of what was removed. Add `remove`, the third way out beside `tick --evidence` (met, perhaps differently) and `close dropped` (the item is abandoned).

Design, from the discussion with Alexander (2026-11-01):

- `remove <id> <n> --reason "..." [--confirm]`. `n` is the box's number in the shared count; `reason` is required and one line.
- It deletes that one box and appends `removed "<box text>": <reason>` to the Log in the same write. The line quotes the whole text, so the breadcrumb is useful without git.
- It removes history. The tool does not track git and the record's previous state is not guaranteed recoverable, which is not a specification invariant. Without `confirm` the call is `refused` with `reason: "confirm_required"`; the message says the operation is lossy and destructive to the record's history, recommends committing first, and says to pass `confirm` to accept that. A human decides whether to commit; an agent surfaces the choice to the user and passes `confirm` only after they decide.
- Unchecked boxes only. A ticked box is the evidence that work was done: `remove` refuses it, naming the Log as the place to say more.
- Open items only. A closed item is read-only (ADR 261101-dxcn).
- Later boxes renumber. The response returns the new numbered list. It takes `section_etag` and `if_match` as `tick` does, so a stale `n` cannot remove the wrong box. Boxes the tool cannot address (a box after prose, a refused construct) are refused as `tick` refuses them.
- `dry_run` shows the diff, so an agent can show the user what will go before asking.
- The `close` refusal for an unchecked Goal names `remove` first.

Done when:

- [ ] `remove` does the above, with tests for confirm, a ticked box, a closed item, renumbering, a stale etag, CRLF and a Goal with a refused construct.
- [ ] API spec (new operation), schema, fixtures in both copies, the `close` refusal text and its fixture.
- [ ] A `references/remove.md` and the `akm` and `akm-track-work` skills (and embedded copies) say that it loses history, to commit first, and who decides.
- [ ] `docs/spec/akm-guide.md` says the same; glossary (Remove) and changelog.

## Next

Start with the `Tree` operation next to `SetChecked` in `internal/body`, then the API op in `internal/api/ops_write.go` beside `opTick`.

## Log

- 2026-11-01: filed from the discussion after 261031-z0n8 (the waive path its Outcome did not cover).
