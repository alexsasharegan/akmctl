---
id: "261101-m3pv"
title: Subcommands reject --help instead of printing usage
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
kind: bug
summary: akmctl <command> --help fails with unknown parameter "help"; only `akmctl help` lists commands, and it has no per-command page
---

# Subcommands reject --help instead of printing usage

## Goal

Hit twice in the session of spike 261101-baze (`akmctl new work x --help`, `akmctl set --help`): `akmctl: invalid_request: unknown parameter "help"`. An agent or newcomer reaches for `--help` first. Done when:

- [x] `<command> --help` and `akmctl help <command>` print that command's flags and one example, and exit 0.
- [x] A test covers every command, so a new command without usage text fails it.

## Next

Done; close after Alexander's review.

## Outcome

akmctl <command> --help, -h and akmctl help <command> print the command's flags (from the API parameter table) and one example, exit 0. TestEveryCommandHasHelp fails for any command without an example. ADR check: none.

## Log

- 2026-11-01: filed.
- 2026-11-01 10:31-07:00: Filed from spike 261101-baze (found by the maintainer's new rule: every akmctl failure is filed).
- 2026-11-01 14:41-07:00: claimed (claude).
- 2026-11-01 15:09-07:00: ticked "`<command> --help` and `akmctl help <command>` print that command's flags and one example, and exit 0.": akmctl close --help, set -h, ui --help and help <command> print flags and an example, exit 0 (TestEveryCommandHasHelp)
- 2026-11-01 15:09-07:00: ticked "A test covers every command, so a new command without usage text fails it.": TestEveryCommandHasHelp fails for a command with no example
- 2026-11-01 15:09-07:00: Added internal/cli/help.go: --help, -h and help <command> print flags (from the API parameter table) and an example; TestEveryCommandHasHelp covers every command.
