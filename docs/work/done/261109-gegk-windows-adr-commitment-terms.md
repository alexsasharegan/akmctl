---
id: "261109-gegk"
title: Rewrite the Windows support tier ADR to the maintainer's commitment terms
status: done
parent: 261108-2ncf
priority: 2
created: 2026-11-09
updated: 2026-11-09
closed: 2026-11-09
---

# Rewrite the Windows support tier ADR to the maintainer's commitment terms

## Goal

The maintainer read the draft of ADR [261108-31p6](../../adr/261108-31p6-windows-support-tier.md) (work item [261108-wg98](../261108-wg98-windows-support-tier-adr.md), held by another session) and gave different terms on 2026-11-09: what is landed and passing in CI on Windows is committed to and stays green; new features are tested on Windows by default; a new feature that fails on Windows is held, not merged, until it works or the maintainer decides it is not supported there; and such a feature is switched off on Windows, as a named, visible exception, while the rest of the tool still works.

This item rewrites the draft's Decision, Consequences and title to those terms. It stays `proposed`; the maintainer accepts it by merging the pull request. It does not touch 261108-wg98, which another session holds and which owns the follow-up edits to the README, `AGENTS.md` and `docs/release.md` on acceptance.

Done when:

- [x] ADR 261108-31p6 states the four terms above, keeps the list of what is not promised yet as the first entries of the unsupported list, and keeps its Review and Revert sections.
- [x] The ADR says how an unsupported feature behaves (a plain refusal with a distinct exit code, a named skip in the Windows tests checked against the list, one list in one place) and what it means for the API spec's operations.
- [x] `docs/adr/README.md` is regenerated with `akmctl index`, and `akmctl check` and `npm run format:check` pass.

## Next

## Outcome

ADR 261108-31p6 now carries the maintainer's commitment terms and stays proposed; the maintainer accepts it by merging the PR. Open points left in the ADR on purpose: the exit code and mechanism for a switched-off feature, and how the API spec declares an operation unavailable on a platform, both decided by the first feature that needs them.

## Log

- 2026-11-09: filed.
- 2026-11-09 17:32-07:00: claimed (claude-sonnet-5-5-31p6).
- 2026-11-09 17:32-07:00: Rewrote ADR 261108-31p6 to the maintainer's four terms (2026-11-09): committed once green, tested by default, held when failing, switched off by name when decided. Option numbering after the rewrite: option 4 is recommended. wg98 is untouched; its follow-up edits (README, AGENTS.md, release.md, superseding 261029-eegx) still wait for acceptance.
