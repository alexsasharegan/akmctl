---
id: "261101-ga4s"
title: tick and log operations on the section tree
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
kind: feature
parent: 261031-gtw6
summary: Add tick and log operations (CLI, API, HTTP) built on body.Tree, with section_etag, per proposed ADR 261101-tr8e
---

# tick and log operations on the section tree

## Goal

From spike [261031-6bzq](261031-6bzq-spike-body-writer-api.md) and proposed ADR [261101-tr8e](../../adr/261101-tr8e-body-writes-through-the-section-tree.md). Add the two operations the spike recommends, on `body.Tree` ([261031-j4bc](261031-j4bc-body-section-tree-package.md)):

- `tick`: `id`, `item` (1-based, in the Goal checklist; `section` defaults to `goal`), `done` (default true), `if_match`, `section_etag`, `dry_run`, `as_of`. CLI `akmctl tick <id> <n> [--undo]`; HTTP `PUT /records/{id}/sections/{name}/items/{n}`.
- `log`: `id`, `text` (one line, not empty), date from `as_of` or today's local day. CLI `akmctl log <id> <text>`; HTTP `POST /records/{id}/log`.

Both read under the lock, split with `Split`, splice, set `updated`, and write once, like `section`. Both refuse a closed record (`closed`), an unsupported construct that matters (`refused`, with the construct and line: a structural one anywhere, or a local one in the edited section, per the ADR), a missing section or item (`not_found`) and a stale `section_etag` (`stale`, with the current one). Responses carry the record, the section with its etag, `changed` and, on a dry run, the `diff`. Done when:

- [x] `tick` and `log` work from the CLI and the HTTP binding with the same JSON envelope, with `dry_run` and `diff`.
- [x] Tests: a tick changes one byte (LF and CRLF); `log` adds one line at the right place; every refusal above; a stale and a good `section_etag`; `updated` is set in the same write.
- [x] Spec (operations and HTTP table), schemas, fixtures in both copies, the `akm` and `akm-track-work` skills (and the embedded copy), the glossary and the changelog are updated.
- [x] The ADR's Confirmation holds: the tests it names pass for the new operations.

## Next

Done; see Outcome.

## Outcome

`tick` and `log` are API operations with CLI and HTTP bindings, built on `Tree.SetChecked` and `Tree.AppendLog`.

- **Behaviour.** `tick` takes `id`, `item`, `section` (default `goal`), `done` (default true); `log` takes `id`, `text`. Both set `updated` in the one write, take `if_match`, `section_etag`, `dry_run` and `as_of`, and answer with the record, the section (`name`, `text`, `etag`), `changed`, the operation's `item` or `entry` and the `diff` on a dry run. A tick that changes nothing writes nothing and leaves `updated` alone.
- **Refusals.** `closed`, `unsupported_markdown` (structural anywhere, local in the edited section), `not_found` (section or item), `stale` (`details.etag` for `if_match`, `details.section_etag` for a stale section) and `invalid_request`. `log` creates a missing Log as `claim` does; with a `section_etag` and no Log it is `not_found`.
- **Bindings.** CLI `akmctl tick <id> <n> [--section s] [--undo]` and `akmctl log <id> <words...>`. HTTP `PUT /records/{id}/sections/{name}/items/{n}` and `POST /records/{id}/log`; a `section_etag` in the body stands in for `If-Match`.
- **Tests.** `internal/api/tick_test.go` (one byte differs, LF and CRLF; one log line; every refusal; a stale and a good `section_etag`, including a change elsewhere in the file not staling it), `TestTickAndLogRoutes`, fixtures `tick-1`..`4` and `log-1`..`2` in both copies (the section etag in them is computed by the Python generator, independently), schemas `tick.json` and `log.json`, `TestEveryOperationHasReference`.
- **Docs.** Spec 4.21 and 4.22, the CLI and HTTP bindings, skill references `tick` and `log` and the embedded copies, glossary, changelog, architecture.
- **Noted, not fixed.** A section added at the end of a file with no Log adds a blank line to the section above it, so that section's etag changes (spec 4.21 says so). Adding before an existing Log does not. The ADR's Confirmation holds: the tests it names pass.

## Log

- 2026-11-01: filed.
- 2026-11-01: filed by spike 261031-6bzq.
- 2026-11-01: refusal rule narrowed in ADR 261101-tr8e; blocked_by 261101-4j86 (the classification).
- 2026-11-01: 261101-4j86 closed; blocked_by removed.
- 2026-11-01: claimed (claude).
- 2026-11-01: tick and log built; closing.
