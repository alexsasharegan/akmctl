---
id: "261031-nep4"
title: release v0 3 3
status: done
created: 2026-10-31
updated: 2026-10-31
closed: 2026-10-31
kind: chore
priority: 10
summary: Cut v0.3.3
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

Close this item.

## Outcome

- Changelog: `## 0.3.3 - 2026-10-31` heading in `CHANGELOG.md`.
- Tag `v0.3.3` (annotated) pushed; the `release` workflow run 41233978127 passed. CI on the prepared commit passed on Linux and macOS (the non-blocking Windows job failed, as before).
- `gh release view v0.3.3` lists the five archives (darwin amd64 and arm64, linux amd64 and arm64, windows amd64) and `SHA256SUMS`.
- `go install github.com/alexsasharegan/akmctl/cmd/akmctl@v0.3.3` then `akmctl version` prints `akmctl 0.3.3 (API 1.0)`.

Decision: none.

## Log

- 2026-10-31: filed.
- 2026-10-31: claimed (claude).
- 2026-10-31: tagged and released; verified.
