---
id: "261101-ga4s"
title: tick and log operations on the section tree
status: raw
created: 2026-11-01
updated: 2026-11-01
kind: feature
parent: 261031-gtw6
summary: Add tick and log operations (CLI, API, HTTP) built on body.Tree, with section_etag, per proposed ADR 261101-tr8e
---

# tick and log operations on the section tree

## Goal

From spike [261031-6bzq](done/261031-6bzq-spike-body-writer-api.md) and proposed ADR [261101-tr8e](../adr/261101-tr8e-body-writes-through-the-section-tree.md). Add the two operations the spike recommends, on `body.Tree` ([261031-j4bc](done/261031-j4bc-body-section-tree-package.md)):

- `tick`: `id`, `item` (1-based, in the Goal checklist; `section` defaults to `goal`), `done` (default true), `if_match`, `section_etag`, `dry_run`, `as_of`. CLI `akmctl tick <id> <n> [--undo]`; HTTP `PUT /records/{id}/sections/{name}/items/{n}`.
- `log`: `id`, `text` (one line, not empty), date from `as_of` or today's local day. CLI `akmctl log <id> <text>`; HTTP `POST /records/{id}/log`.

Both read under the lock, split with `Split`, splice, set `updated`, and write once, like `section`. Both refuse a closed record (`closed`), an unsupported construct (`refused`, with the construct and line), a missing section or item (`not_found`) and a stale `section_etag` (`stale`, with the current one). Responses carry the record, the section with its etag, `changed` and, on a dry run, the `diff`. Accept the ADR when this lands. Done when:

- [ ] `tick` and `log` work from the CLI and the HTTP binding with the same JSON envelope, with `dry_run` and `diff`.
- [ ] Tests: a tick changes one byte (LF and CRLF); `log` adds one line at the right place; every refusal above; a stale and a good `section_etag`; `updated` is set in the same write.
- [ ] Spec (operations and HTTP table), schemas, fixtures in both copies, the `akm` and `akm-track-work` skills (and the embedded copy), the glossary and the changelog are updated.
- [ ] ADR 261101-tr8e is `accepted` in the PR that lands this.

## Next

Not started. Read the spike's Outcome and ADR 261101-tr8e; make ready once the ADR is agreed.

## Log

- 2026-11-01: filed.
- 2026-11-01: filed by spike 261031-6bzq.
