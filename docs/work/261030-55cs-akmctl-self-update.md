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
- The daily check never runs inside an API operation, the HTTP server or `ui`. A command may start it in a goroutine with a short timeout after its own work is done, and exits without waiting. It records the time and the latest version it saw in the user cache directory; a notice is printed to stderr (never stdout, never inside JSON) from that cache on a later run, at most once a day.
- Off by `AKMCTL_NO_UPDATE_CHECK=1`, off when `CI` is set or stderr is not a terminal, and a config switch in the ADR's decision if one is wanted.
- Standard library only (`net/http`, `crypto/sha256`, `archive/tar`, `archive/zip`); no new dependency.

Constraints:

- Releases are made by the maintainer's tag and the `release` workflow (`docs/release.md`); the archive names and `SHA256SUMS` format come from `scripts/release.sh`, and the tool must accept exactly those.
- Fetching a release without credentials needs the repository to be public. Until the maintainer makes it so (AGENTS.md: that is their call), build and test against a local test server and a fake release; the check against GitHub itself waits.
- Windows cannot overwrite a running executable; it is best effort and not required below 1.0.0 (ADR 261029-eegx), so the rename-aside method is allowed to be untested there.

Done when:

- [ ] An ADR (`proposed`) records: the command and its flags, that the check never runs in an API operation, the server or `ui`, the cache file and opt-outs, what is verified (checksum from the same release; whether a signature is wanted), and the narrowed guarantee. The specification's "needs no network" sentence and the API spec's "no network call" sentence are reworded to say which commands may use the network and that none of the API operations does.
- [ ] `self-update` replaces the binary from a test server's fake release: wrong checksum, wrong archive name, no newer version and a `go install` binary each behave as the ADR says (a test for each).
- [ ] The daily check is non-blocking (a test with a server that never answers shows the command exits on time), prints one stderr line from the cache at most once a day, and prints nothing under `CI`, a non-terminal or the opt-out.
- [ ] The glossary, the README install section, `docs/release.md` (the archive names are now an interface) and the changelog are updated.
- [ ] After the maintainer makes the repository public: `self-update --check` against the real latest release is run once and the Log records the result.

## Next

Write the ADR (`proposed`) first: `akmctl new adr tool-self-update`. Then the checksum verify and replace logic against a test server.

## Log

- 2026-10-30: filed.
- 2026-10-30: refined by Claude from the one-line request: shape, constraints (the no-network guarantee, the public-repository dependency, Windows) and done-when written; marked ready.
