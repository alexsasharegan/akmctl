---
id: "261029-0r4h"
title: go install reports 0.0.0-dev instead of the module version
status: done
created: 2026-10-29
updated: 2026-10-29
closed: 2026-10-29
kind: bug
priority: 1
summary: akmctl version after go install of a tag printed 0.0.0-dev; read the module version from build info
---

# go install reports 0.0.0-dev instead of the module version

## Goal

`go install ...@v0.2.0` printed `0.0.0-dev`, because only the commit was read from the build info.

Done when:

- [x] a module version from build info is the version, a link-time version wins, a development build keeps `0.0.0-dev+<commit>`

## Outcome

resolveVersion in internal/cli/version.go reads Main.Version; unit test covers four cases. Shipped as v0.2.1.

## Log

- 2026-10-29: filed.
- 2026-10-29: claimed (claude).
