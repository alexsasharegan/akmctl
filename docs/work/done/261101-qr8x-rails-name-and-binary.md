---
id: "261101-qr8x"
title: rails name and binary
status: done
created: 2026-11-01
updated: 2026-11-06
closed: 2026-11-06
parent: 261101-vs5c
priority: 1
---

# rails name and binary

## Goal

The rerun showed two wording gaps. Claude claimed as plain `claude` (CL2 risk when sessions share a name); Codex chose a unique name unprompted. Codex also ran the installed `akmctl`, which was stale, got `unknown command "tick"` and filed a raw item for a tool bug that was only an old binary. Done when:

- [x] `AGENTS.md` says to claim with a name unique to the session (`--as <agent>-<model>`).
- [x] `AGENTS.md` says which `akmctl` to use in this repository for work commands (the build of the current branch) and that an unknown command may mean a stale binary: check `akmctl version` before filing.

## Next


## Outcome

AGENTS.md, the init template and the claim skill text now ask for a unique --as name (<agent>-<model>) and say an unknown command may be a stale binary: use /tmp/akmctl here, run version before filing. Not verified with a model run; the rerun protocol (261101-0vh7) would test it.

## Log

- 2026-11-01: filed.
- 2026-11-06 11:00-07:00: claimed (claude-sonnet-5-5).
- 2026-11-06 11:00-07:00: AGENTS.md, the init template, both copies of the claim skill text and the help example now say a unique --as name and to check the binary version before filing an unknown-command failure.
