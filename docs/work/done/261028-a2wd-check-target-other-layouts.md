---
id: "261028-a2wd"
title: Let check --target accept layouts other than the tool's own
summary: check --target classifies findings only for the tool's own layout version; it should accept any target a migration chain reaches.
status: done
kind: feature
priority: 1
created: 2026-10-28
updated: 2026-11-09
closed: 2026-11-09
---

# Let check --target accept layouts other than the tool's own

## Goal

`akmctl check --target <version>` classifies each finding as `auto`, `manual` or `info` for a target layout, but it accepts only the tool's own layout version. With migrations embedded, a repository on an older layout should be checkable against any target the chain reaches.

Done when:

- [x] `check --target` accepts any layout version that the embedded migration chain can reach from the repository's manifest, and rejects others with a clear error.
- [x] Findings are classified using the steps of that chain.
- [x] A fixture with a two-step chain proves the classification.

## Next

## Outcome

check --target accepts the tool's own layout, the repository's own and any version the embedded chain reaches, and classifies by that chain's steps; others are refused with the accepted versions named. Proven with a two-step chain fixture; API spec updated.

## Log

- 2026-10-28: filed.
- 2026-11-09 09:55-07:00: claimed (claude-sonnet-5-5-a2wd).
- 2026-11-09 09:55-07:00: ticked "`check --target` accepts any layout version that the embedded migration chain can reach from the repository's manifest, and rejects others with a clear error.": TestCheckTargetAcceptsAnyLayoutTheChainReaches, TestCheckTargetRefusesALayoutTheChainDoesNotReach
- 2026-11-09 09:55-07:00: ticked "Findings are classified using the steps of that chain.": TestCheckTargetAcceptsAnyLayoutTheChainReaches
- 2026-11-09 09:55-07:00: ticked "A fixture with a two-step chain proves the classification.": twoStepRegistry in internal/api/check_target_test.go
- 2026-11-09 09:55-07:00: check --target accepts the tool's own layout, the repository's own and any version the embedded chain reaches; others are refused with the accepted list. migrate.Reachable added. Findings stay those of this release's rule set (the migrate op says the same); only class comes from the chain.
