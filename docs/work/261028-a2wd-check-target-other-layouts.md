---
id: "261028-a2wd"
title: Let check --target accept layouts other than the tool's own
summary: check --target classifies findings only for the tool's own layout version; it should accept any target a migration chain reaches.
status: active
kind: feature
priority: 1
created: 2026-10-28
updated: 2026-11-09
claimed_by: claude-sonnet-5-5-a2wd
claimed_at: 2026-11-09T16:55
---

# Let check --target accept layouts other than the tool's own

## Goal

`akmctl check --target <version>` classifies each finding as `auto`, `manual` or `info` for a target layout, but it accepts only the tool's own layout version. With migrations embedded, a repository on an older layout should be checkable against any target the chain reaches.

Done when:

- [ ] `check --target` accepts any layout version that the embedded migration chain can reach from the repository's manifest, and rejects others with a clear error.
- [ ] Findings are classified using the steps of that chain.
- [ ] A fixture with a two-step chain proves the classification.

## Next

Accept as a check target the tool's own layout, the repository's own and any version the embedded chain reaches; classify by that chain's steps; refuse others naming the accepted versions. Test with a two-step chain; update the API spec.

## Log

- 2026-10-28: filed.
- 2026-11-09 09:55-07:00: claimed (claude-sonnet-5-5-a2wd).
