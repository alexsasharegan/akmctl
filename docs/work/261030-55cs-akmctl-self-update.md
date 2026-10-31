---
id: "261030-55cs"
title: akmctl self update
summary: akmctl updates itself from a GitHub release, and tells the user at most once a day, without blocking, that a newer version exists.
status: ready
kind: feature
priority: 3
created: 2026-10-30
updated: 2026-10-30
---

# akmctl self update

## Goal

The binary can replace itself with the latest release, and a short, non-blocking daily check tells the user when one exists.

Today a user updates by downloading an archive by hand or by `go install`. The tool also promises "needs no network" (the specification, H6) and "no operation makes a network call" (API spec, section 1), so this item changes a stated guarantee and needs an ADR before any code.

Shape (open points go to the ADR):

- `akmctl self-update` is a CLI command like `ui`, not an API operation. It finds the latest release of this repository, downloads the archive for the running OS and architecture and `SHA256SUMS`, verifies the checksum, and replaces the running binary atomically (write beside it, then rename). `--check` only reports. A binary installed by `go install` or a package manager reports how to update that way and changes nothing.
- The daily check never runs inside an API operation, the HTTP server or `ui`. The API operations still make no network call; the carve-out is that the CLI binding may, after it has printed its result, start a check (see "The check" below). The tool does not need the network to do its job: a failed or skipped check changes nothing about any command.
- Off by `AKMCTL_NO_UPDATE_CHECK=1`, off when `CI` is set or stderr is not a terminal, and a manifest key a repository can set at `init` (for example a team that does not want its agents spawning a network process); the ADR names the key.
- Standard library only (`net/http`, `crypto/sha256`, `archive/tar`, `archive/zip`); no new dependency.

Constraints:

- Releases are made by the maintainer's tag and the `release` workflow (`docs/release.md`); the archive names and `SHA256SUMS` format come from `scripts/release.sh`, and the tool must accept exactly those.
- Fetching a release without credentials needs the repository to be public. Until the maintainer makes it so (AGENTS.md: that is their call), build and test against a local test server and a fake release; the check against GitHub itself waits.
- Windows cannot overwrite a running executable; it is best effort and not required below 1.0.0 (ADR 261029-eegx), so the rename-aside method is allowed to be untested there.

Done when:

- [ ] An ADR (`proposed`) records: the command and its flags, that the check never runs in an API operation, the server or `ui`, the state file in the user cache directory (not the manifest), the detached child, the opt-outs, what is verified (checksum from the same release; whether a signature is wanted), and the narrowed guarantee. The specification's "needs no network" sentence and the API spec's "no network call" sentence are reworded to say which commands may use the network and that none of the API operations does.
- [ ] `self-update` replaces the binary from a test server's fake release: wrong checksum, wrong archive name, no newer version and a `go install` binary each behave as the ADR says (a test for each).
- [ ] The daily check is non-blocking (a test with a server that never answers shows the command exits on time), prints one stderr line from the cache at most once a day, and prints nothing under `CI`, a non-terminal or the opt-out.
- [ ] The glossary, the README install section, `docs/release.md` (the archive names are now an interface) and the changelog are updated.
- [ ] After the maintainer makes the repository public: `self-update --check` against the real latest release is run once and the Log records the result.

## The check

Decided in discussion with the maintainer (2026-10-30); the ADR records it.

- **Never on the hot path.** A command reads one small state file; if it says the last check is older than a day, the CLI starts the check after the command has printed its result and does not wait for it. A stale or missing file, an unreadable file, no network and a failed fetch all do nothing and say nothing.
- **State is per user, not per repository.** The file lives in the user cache directory (`os.UserCacheDir()/akmctl/update.json`): `checked_at`, `latest` (the newest version seen), `notified_at` and `failed_at`. It is not the repository's `.layout`: that file is committed and shared, so stamping it on every run would dirty the working tree, make merge noise and break "a read-only command writes nothing", and the tool's version is a fact about a machine, not a repository. The repository's manifest holds only a policy (the opt-out), never a timestamp.
- **How it runs.** A short-lived CLI cannot rely on a goroutine finishing before exit, so the CLI starts a detached child of itself (a hidden internal subcommand, no stdio, a few seconds' timeout) that fetches the latest release, writes the state file atomically and exits. The parent returns at once. A failure writes `failed_at` so an offline machine backs off for hours instead of retrying on every command.
- **The notice.** On a later run, if the cache says a newer version exists and none was announced in the last day, one line goes to stderr, only when stderr is a terminal. Agents run without a terminal, so they never see it and never update unprompted. The check only ever announces; only `self-update` installs.

### Prior art (from memory; verify while writing the ADR)

- `npm` (`update-notifier`): detached child checks the registry on an interval, result is cached in a user config file, notice on the next run; off in CI and by an environment variable. This is the closest model to the one above.
- `gh`: checks for a new release in the background, stores `checked_for_update_at` in its per-user state file, once a day, notice on stderr only on a terminal, off by an environment variable. Same shape in Go.
- Oh My Zsh: a timestamp file records the last check; at shell start it prompts or reminds every N days, or updates automatically. The timestamp pattern is the same, but its prompt is interactive, which a CLI used by agents must not copy.
- Homebrew: updates itself before a command, throttled. That is the blocking design to avoid.
- `kubectl` and the AWS CLI: no self-update; the package manager or installer is the update path, and `kubectl` only warns about version skew. Terraform uses an opt-out version check (`CHECKPOINT_DISABLE`). `rustup`, `deno` and `bun` update only on an explicit command.
- Conclusion: an explicit `self-update`, plus a cached, detached, stderr-only notice, is the common shape for a single static binary; a package-manager install should defer to the package manager.

## Next

Write the ADR (`proposed`) first: `akmctl new adr tool-self-update`. Then the checksum verify and replace logic against a test server.

## Log

- 2026-10-30: filed.
- 2026-10-30: refined by Claude from the one-line request: shape, constraints (the no-network guarantee, the public-repository dependency, Windows) and done-when written; marked ready.
- 2026-10-30: recorded the maintainer's design: never on the hot path, a staleness timestamp read from a state file, the check queued after the command is served, `init` can opt a repository out, failures are silent. Moved the timestamp from the manifest to a per-user cache file and added prior art (npm, gh, Oh My Zsh, Homebrew, kubectl, AWS CLI).
