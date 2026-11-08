---
id: "261108-8m4a"
title: global akmctl follows green main
status: done
kind: feature
priority: 2
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
---

# global akmctl follows green main

## Goal

The `akmctl` on the maintainer's `PATH` (`~/golang/bin/akmctl`) is what every other project's threads run, but it is installed by hand, from whatever working tree a session happened to be in. On 2026-11-08 it was a build of a feature branch, stamped `+dirty` (an untracked `.obsidian/` folder makes every local build dirty), and it lagged `main` by a merge; a project whose skills came from a newer build then lost commands. Make the global binary follow the newest commit of `main` whose CI passed, with nothing built from a working tree.

Design sketch (to confirm): a script `scripts/install-latest.sh` finds the newest `main` commit whose `ci` run succeeded on macOS, Ubuntu and Windows (`gh run list --branch main`), skips if the installed binary already is that commit (the version string carries the 12-character revision), and otherwise runs `GOPROXY=direct go install github.com/alexsasharegan/akmctl/cmd/akmctl@<sha>`, which builds a clean checkout (verified by hand on 2026-11-08: `akmctl 0.3.4-0.20261108190235-06b1081e14da`, no `+dirty`). Run it on a schedule (a launchd agent the maintainer sets up) and after a merge. Sessions in this repository keep using `/tmp/akmctl` and never `go install`.

Done when:

- [x] `scripts/install-latest.sh` exists, is tested against a fake `gh` and a dummy module, and refuses when no green commit is found.
- [x] `AGENTS.md` says: nobody runs `go install` from a working tree; the global binary is updated by the script. (The shared instructions in the maintainer's dotfiles are the maintainer's to edit.)
- [x] A schedule (a launchd agent) is documented in `docs/release.md`; the maintainer sets it up, not an agent.
- [x] The `.obsidian/` folder is excluded locally (`.git/info/exclude`) so local builds are not marked dirty, or the decision not to is recorded.

## Next

## Outcome

scripts/install-latest.sh installs the newest main commit whose ci run passed in every job, by commit from a clean checkout, and skips when it is installed; tested with fake tools and run for real. AGENTS.md, the glossary, docs/release.md (with a launchd example) and the changelog are updated. Still the maintainer's: setting up the schedule and adding the rule to the shared dotfiles instructions.

## Log

- 2026-11-08: filed.
- 2026-11-08 12:45-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08 12:46-07:00: scripts/install-latest.sh with TestInstallLatest (fake gh, go and akmctl: skips a run with a failed job, installs by full hash, up to date, dry run, no green run, bad flag). Ran it for real: akmctl 0.3.4-0.20261108191956-aea1b146097b, no +dirty. .obsidian/ added to .git/info/exclude (local, not tracked), so a local build is no longer stamped modified. launchd example in docs/release.md, not installed; shared dotfiles instructions not edited.
