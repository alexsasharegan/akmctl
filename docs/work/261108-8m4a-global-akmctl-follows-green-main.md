---
id: "261108-8m4a"
title: global akmctl follows green main
status: active
kind: feature
priority: 2
created: 2026-11-08
updated: 2026-11-08
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-08T19:45
---

# global akmctl follows green main

## Goal

The `akmctl` on the maintainer's `PATH` (`~/golang/bin/akmctl`) is what every other project's threads run, but it is installed by hand, from whatever working tree a session happened to be in. On 2026-11-08 it was a build of a feature branch, stamped `+dirty` (an untracked `.obsidian/` folder makes every local build dirty), and it lagged `main` by a merge; a project whose skills came from a newer build then lost commands. Make the global binary follow the newest commit of `main` whose CI passed, with nothing built from a working tree.

Design sketch (to confirm): a script `scripts/install-latest.sh` finds the newest `main` commit whose `ci` run succeeded on macOS, Ubuntu and Windows (`gh run list --branch main`), skips if the installed binary already is that commit (the version string carries the 12-character revision), and otherwise runs `GOPROXY=direct go install github.com/alexsasharegan/akmctl/cmd/akmctl@<sha>`, which builds a clean checkout (verified by hand on 2026-11-08: `akmctl 0.3.4-0.20261108190235-06b1081e14da`, no `+dirty`). Run it on a schedule (a launchd agent the maintainer sets up) and after a merge. Sessions in this repository keep using `/tmp/akmctl` and never `go install`.

Done when:

- [ ] `scripts/install-latest.sh` exists, is tested against a fake `gh` and a dummy module, and refuses when no green commit is found.
- [ ] `AGENTS.md` and the shared agent instructions say: nobody runs `go install` from a working tree; the global binary is updated by the script.
- [ ] The maintainer has a schedule for it (launchd or another), documented in `docs/` and set up by the maintainer, not by an agent.
- [ ] The `.obsidian/` folder is excluded locally (`.git/info/exclude`) so local builds are not marked dirty, or the decision not to is recorded.

## Next

Write scripts/install-latest.sh (newest green main commit, skip if installed, go install by sha), a test with a fake gh and go, docs and the AGENTS rule.

## Log

- 2026-11-08: filed.
- 2026-11-08 12:45-07:00: claimed (claude-sonnet-5-5).
