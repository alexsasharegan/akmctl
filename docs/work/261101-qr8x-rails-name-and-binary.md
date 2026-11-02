---
id: "261101-qr8x"
title: rails name and binary
status: ready
created: 2026-11-01
updated: 2026-11-01
parent: 261101-vs5c
priority: 1
---

# rails name and binary

## Goal

The rerun showed two wording gaps. Claude claimed as plain `claude` (CL2 risk when sessions share a name); Codex chose a unique name unprompted. Codex also ran the installed `akmctl`, which was stale, got `unknown command "tick"` and filed a raw item for a tool bug that was only an old binary. Done when:

- [ ] `AGENTS.md` says to claim with a name unique to the session (`--as <agent>-<model>`).
- [ ] `AGENTS.md` says which `akmctl` to use in this repository for work commands (the build of the current branch) and that an unknown command may mean a stale binary: check `akmctl version` before filing.

## Next

Edit AGENTS.md (the Rules section): one sentence for a unique --as name, one for building and using /tmp/akmctl and checking akmctl version before filing an unknown-command failure. Mirror the wording in the init templates if they carry the same rules.

## Log

- 2026-11-01: filed.
