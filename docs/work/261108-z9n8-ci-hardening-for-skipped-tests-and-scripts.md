---
id: "261108-z9n8"
title: ci hardening for skipped tests and scripts
status: active
kind: chore
priority: 3
created: 2026-11-08
updated: 2026-11-08
claimed_by: claude-sonnet-5-5-rest
claimed_at: 2026-11-08T20:17
---

# ci hardening for skipped tests and scripts

## Goal

Findings F3, F5, F7 and F8 of [the audit](../design/test-strategy-audit.md), from [261030-456q](done/261030-456q-test-strategy-and-coverage-audit.md): a runner can pass with tests skipped for want of `node` or `git`, the release and install scripts run only on a tag or in a user's shell, `-race` does not run on macOS, and fuzz targets only run their seeds. Make CI fail when a required tool is missing, test `scripts/release.sh` and `scripts/install.sh` against a dummy version, add `-race` on macOS, and add a scheduled or manual `-fuzz` job.

Done when:

- [ ] CI fails when `node` or `git` is missing (each skip becomes a failure when `AKM_REQUIRE_TOOLS=1`), and a decision on the `AKM_CORPUS` test is recorded.
- [ ] A CI job runs `scripts/release.sh` for a dummy version and `scripts/install.sh` against its archive in a temporary directory.
- [ ] `go test -race` runs on macOS as well as Linux.
- [ ] A scheduled or manual job runs each fuzz target with `-fuzz` for a minute.

## Next

Start with the `AKM_REQUIRE_TOOLS` switch (F3): find each `t.Skip` and route it through one helper.

## Log

- 2026-11-08: filed.
- 2026-11-08 13:17-07:00: claimed (claude-sonnet-5-5-rest).
