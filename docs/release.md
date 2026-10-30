# Releasing

How a version of `akmctl` is cut. The maintainer tags and pushes (AGENTS.md); an agent prepares everything up to the tag.

## Prepare

1. File a `chore` work item for the release (`akmctl new work release-v0-3-2`), claim it, and put it under Now in `docs/ROADMAP.md`. Its done-when: the changelog heading, the tag pushed by the maintainer, the archives attached and `go install` reporting the version. Leave `## Next` as the tag step.
2. Pick the version. Below 1.0.0 a dot release (`0.3.1` to `0.3.2`) carries fixes and small additions; a MINOR bump (`0.3.x` to `0.4.0`) carries a layout or API change.
3. In `CHANGELOG.md`, put a `## <version> - <date>` heading under `## Unreleased` and move the entries beneath it. `Unreleased` stays, empty.
4. Run the full check from AGENTS.md: `go vet ./... && go test ./...`, `gofmt -l .`, `/tmp/akmctl check`, `npm run format:check`, `npm run schemas`.
5. Optional dry run of the build: `SKIP_TESTS=1 scripts/release.sh <version>` writes `dist/<version>/` (git-ignored) with the five archives and `SHA256SUMS`.
6. Commit the changelog, the item and the roadmap together, and push `main`. Wait for CI on that commit to pass.

## Tag (the maintainer)

Tags are annotated, with the message `v<version> (unstable pre-1.0)`; a lightweight tag does not match the earlier ones.

```sh
git switch main && git pull --ff-only
git tag -a v<version> -m "v<version> (unstable pre-1.0)"
git push origin v<version>
```

Pushing a `v*` tag is the trigger. The `release` workflow (`.github/workflows/release.yml`) runs `scripts/release.sh`, which runs `go vet` and `go test`, builds the archives for darwin and linux (amd64, arm64) and windows (amd64), and links the version in with `-X .../internal/cli.Version`. It then creates the GitHub release with the archives, `SHA256SUMS` and generated notes. A run takes about two minutes.

## Verify and close

1. `gh release view v<version>` lists the five archives and `SHA256SUMS`.
2. `go install github.com/alexsasharegan/akmctl/cmd/akmctl@v<version> && akmctl version` prints `<version>`.
3. Close the release item with an Outcome that names both checks, take it off the roadmap, and commit.

## If it goes wrong

- The workflow fails after the tag is pushed: fix forward and release the next dot version. Do not move or delete a pushed tag; `go install` and the module proxy cache it.
- A push of the tag fails with `Permission denied (publickey)`: the SSH agent needs unlocking or approving; nothing has been published.
