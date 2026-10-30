# Roadmap

## Now

- [261030-2abz] Beta release for greenfield use (v0.2.0)

## Next

Beta (greenfield only; adoption is not claimed). Order of work:

- [261029-0gb4] Add akmctl init with the profiles code and notes
- [261029-dsga] Rewrite the standard's template section to describe init
- [261029-sjyg] Prettier is not run by CI

## Later

1.0.0 requirements (adoption and safety):

- [261029-tgfe] Give migrate a discovery phase for repositories with no manifest
- [261028-a2wd] Let check --target accept layouts other than the tool's own
- [261029-47m0] migrate plan: simulate move_file steps
- [261029-m0vm] move_file: report the inbound links it leaves broken
- [261029-xtqj] Multi-write moves can stop half-way with no rollback
- [261029-shz5] Undo covers only the last safe move or capture
- [261029-x4ws] Keep one copy of the API fixtures
- [261029-z4gk] Rule L2: understand escaped pipes in wikilinks inside tables
- [261029-wpak] Rule S2: accept block scalars in SKILL.md frontmatter
- [261029-nkqv] A new ADR made by the done wizard can be orphaned

Workbench and polish:

- [261029-30my] Filter the board and tables by project
- [261029-ctzc] Make wikilinks in drawer fields clickable and open them in Obsidian
- [261029-paw7] Show a record's frontmatter as a typed table in the drawer
- [261029-rmrk] Fix the card border that shows on only three sides
- [261029-57jq] The closed table flickers on reload
- [261029-8948] Show the claim holder on workbench cards
- [261029-g8x2] The review diff is against the current file, not after earlier writes
- [261029-jm46] Restart of akmctl ui forces a new URL
- [261028-k0ek] Give the workbench live refresh
- [261028-esyt] Edit more than status in the workbench
- [261028-fcc1] Add a capture form to the workbench
- [261026-mrdr] pkm-workbench: sort options
- [261026-np3c] pkm-workbench: delete a task, with its backlinks listed
- [261026-zg0g] pkm-workbench: filters

Hardening and small rules:

- [261029-q8z3] TestCompareAndSwapRace is flaky on Windows CI
- [261028-jzsx] Cover the Windows path beyond CI
- [261028-e7ej] Make rule G3 check more command forms than scripts/ paths
- [261028-1k0e] Say which is which in the ADR index supersedes column
- [261028-x1b4] Match Prettier's width on non-standard emoji ZWJ pairs
- [261029-ap3y] Keep docs/.akm.lock out of commits

## Dropped
