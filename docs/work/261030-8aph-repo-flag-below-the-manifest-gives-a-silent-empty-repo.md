---
id: "261030-8aph"
title: --repo pointing inside a repository is read as an empty one, with no warning
summary: --repo takes the folder as the root and does not look upward for the manifest, so a folder inside a repository reads as one with default folders.
status: ready
kind: bug
priority: 3
created: 2026-10-30
updated: 2026-10-31
---

# --repo pointing inside a repository is read as an empty one, with no warning

## Goal

`--repo <dir>` where `<dir>` is inside a repository, below the folder that holds the manifest, gives a wrong answer without saying so.

Found 2026-10-31 while auditing a vault. The vault's manifest is at the git root (`docs/.layout`, `work_dir: Deep Thought/9 pkm/tasks`). `akmctl commits --repo "<root>/Deep Thought"` took that folder as the root, found no manifest, used the default `docs/work`, and reported `CM3 no file of <id> changed` on correct commits. Run with no `--repo` from the same folder, the tool climbed to the manifest and passed. `record.Load` (internal/record/repo.go) takes the folder as given and never looks upward; the CLI climbs with `FindRoot` only when `--repo` is absent.

Decide, then do one of:

- `--repo` climbs to the nearest manifest, as the default does (what `git -C` and `--repo` suggest to a user), or
- `--repo` stays literal, and a root with no manifest says so: a `no_manifest` note on reads, and a refusal from `commits`, which cannot judge anything without the folders.

Related: 261029-tgfe covers repositories that have no manifest at all; this is a manifest that exists above the given folder.

Done when:

- [ ] The chosen behaviour is stated in the API spec (conventions) and the CLI usage.
- [ ] `commits --repo <folder inside a repository>` no longer reports a false CM3.
- [ ] A test covers a manifest above the given folder, for the CLI and the API.

## Log

- 2026-10-30: filed.
