---
id: "261029-47m0"
title: "migrate plan: simulate move_file steps"
status: done
created: 2026-10-29
updated: 2026-11-09
closed: 2026-11-09
kind: bug
priority: 2
summary: plan counted ADR Candidates as an ADR (22 ids) because it did not apply move_file steps; apply made 21.
---

# migrate plan: simulate move_file steps

## Goal

plan applies the steps on a scratch copy but plans move_file and add_file against the real tree, so a file moved out of a directory is still counted.

Done when:

- [x] plan and apply agree on the number of ids and renames on the vault fixture.
- [x] A test covers a move_file before a mint_ids on the same directory.

## Next

## Outcome

plan and apply now agree: PlanEnv applies move_file and add_file on the scratch copy, mirroring non-Markdown paths on demand. Covered by a synthetic test (move_file before mint_ids on one directory) and one on the vault fixture.

## Log

- 2026-10-29: filed.
- 2026-11-09 09:34-07:00: claimed (claude-sonnet-5-5-47m0).
- 2026-11-09 09:35-07:00: ticked "plan and apply agree on the number of ids and renames on the vault fixture.": TestPlanAgreesWithApplyOnTheVaultFixture
- 2026-11-09 09:35-07:00: ticked "A test covers a move_file before a mint_ids on the same directory.": TestPlanSimulatesMoveFileBeforeMintIDs
- 2026-11-09 09:35-07:00: PlanEnv now plans move_file and add_file on the scratch copy, mirroring the path a step names from the real tree when it is not Markdown (mirrorFile); tests added; the first fails without the fix.
