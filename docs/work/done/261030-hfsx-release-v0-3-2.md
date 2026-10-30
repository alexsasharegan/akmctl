---
id: "261030-hfsx"
title: Release v0.3.2 with the Log append fix
status: done
created: 2026-10-30
updated: 2026-10-30
closed: 2026-10-30
priority: 1
summary: Cut v0.3.2 so the Log append bullet fix ships
---

# Release v0.3.2 with the Log append fix

## Goal

Ship the Log append fix (261030-4d9b) as v0.3.2 (Alexander asked for a dot release).

Done when:

- [ ] CHANGELOG.md has a 0.3.2 heading holding the entry
- [ ] v0.3.2 is tagged and pushed by Alexander
- [ ] the release has its archives and `go install` reports 0.3.2

## Next

Alexander tags and pushes v0.3.2 (git tag v0.3.2 && git push origin v0.3.2), then confirms the release archives and go install; close this item then.

## Outcome

- CHANGELOG.md has a 0.3.2 heading holding the entry: commit 1a00673.
- v0.3.2 is tagged and pushed by Alexander: tag v0.3.2.
- The release has its archives and go install reports 0.3.2: gh release view v0.3.2 lists the five archives and SHA256SUMS, and go install ...@v0.3.2 then akmctl version prints 0.3.2.

## Log

- 2026-10-30: filed.
- 2026-10-30: claimed (claude).
- 2026-10-30: changelog heading 0.3.2 added; waiting for the tag, which is the maintainer's.
- 2026-10-30: v0.3.2 tag pushed by Alexander; wrote docs/release.md and linked it from AGENTS.md.
- 2026-10-30: first release run failed at publish with a GitHub HTTP 500 (no release created); re-ran the failed job and it succeeded.
