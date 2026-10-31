---
id: "261031-nep4"
title: release v0 3 3
status: active
created: 2026-10-31
updated: 2026-10-31
kind: chore
priority: 10
summary: Cut v0.3.3
claimed_by: claude
claimed_at: 2026-10-31T19:21
---

# release v0 3 3

## Goal

Cut v0.3.3 following `docs/release.md`.

Done when:

- [ ] `CHANGELOG.md` has a `## 0.3.3 - 2026-10-31` heading with the Unreleased entries beneath it.
- [ ] The tag `v0.3.3` is pushed and the release workflow passes.
- [ ] `gh release view v0.3.3` lists the five archives and `SHA256SUMS`.
- [ ] `go install github.com/alexsasharegan/akmctl/cmd/akmctl@v0.3.3 && akmctl version` prints `0.3.3`.

## Next

Prepare the changelog and run the full check, then push main, tag and verify.

## Log

- 2026-10-31: filed.
- 2026-10-31: claimed (claude).
