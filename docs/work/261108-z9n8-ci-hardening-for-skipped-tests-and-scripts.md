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

- [x] CI fails when `node` or `git` is missing (each skip becomes a failure when `AKM_REQUIRE_TOOLS=1`), and a decision on the `AKM_CORPUS` test is recorded.
- [x] A CI job runs `scripts/release.sh` for a dummy version and `scripts/install.sh` against its archive in a temporary directory.
- [x] `go test -race` runs on macOS as well as Linux.
- [x] A scheduled or manual job runs each fuzz target with `-fuzz` for a minute.

## Next

Close after CI runs the new jobs. Follow-up: fix 261108-xkck so the weekly fuzz job is green; rwqj will grow internal/testutil.

## Log

- 2026-11-08: filed.
- 2026-11-08 13:17-07:00: claimed (claude-sonnet-5-5-rest).
- 2026-11-08 13:17-07:00: AKM_REQUIRE_TOOLS=1 set for the test job; internal/testutil.RequireTool replaces the git, node and sh skips (AKM_CORPUS stays a skip: a manual measurement on notes CI does not have). New packaging job: release.sh 0.0.0-ci, then install.sh from the local archive (new AKMCTL_BASE_URL). -race on macOS too. New fuzz.yml (weekly and manual, each target 1m). A 2s local fuzz run found 261108-xkck (frontmatter reader vs oracle on 'A: <<'); the fuzz job fails until that is fixed. Not run here: the workflow files themselves, until CI runs them.
