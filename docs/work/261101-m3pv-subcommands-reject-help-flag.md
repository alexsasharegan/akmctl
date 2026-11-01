---
id: "261101-m3pv"
title: Subcommands reject --help instead of printing usage
status: active
created: 2026-11-01
updated: 2026-11-01
kind: bug
summary: akmctl <command> --help fails with unknown parameter "help"; only `akmctl help` lists commands, and it has no per-command page
claimed_by: claude
claimed_at: 2026-11-01T21:41
---

# Subcommands reject --help instead of printing usage

## Goal

Hit twice in the session of spike 261101-baze (`akmctl new work x --help`, `akmctl set --help`): `akmctl: invalid_request: unknown parameter "help"`. An agent or newcomer reaches for `--help` first. Done when:

- [ ] `<command> --help` and `akmctl help <command>` print that command's flags and one example, and exit 0.
- [ ] A test covers every command, so a new command without usage text fails it.

## Next

Find where commands parse flags in `cmd/akmctl`, add a per-command usage table and the `--help` / `help <command>` path, then the every-command test.

## Log

- 2026-11-01: filed.
- 2026-11-01 10:31-07:00: Filed from spike 261101-baze (found by the maintainer's new rule: every akmctl failure is filed).
- 2026-11-01 14:41-07:00: claimed (claude).
