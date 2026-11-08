# Releasing

How a version of `akmctl` is cut. The maintainer tags and pushes (AGENTS.md); an agent prepares everything up to the tag.

## When to release

Before 1.0.0 there are no hard rules; these are hints, and the maintainer decides.

- **Version.** A change that breaks callers (a new refusal on an existing operation, a changed output, a layout change) is a minor bump (`0.x.0`). A patch (`0.x.y`) carries fixes and docs. A new operation or parameter that breaks nothing is the maintainer's call.
- **1.0.0** means the API is frozen. It is not what a breaking release becomes; breaking changes are expected until then.
- **A good moment.** A logical chunk is finished: the code, fixtures, spec and docs agree and `check` is clean (the body-schema work was one). Many unreleased changelog entries, or a fix someone is waiting on, are also reasons.

## Prepare

1. File a `chore` work item for the release (`akmctl new work release-v0-3-2`), claim it, and put it under Now in `docs/ROADMAP.md`. Its done-when: the changelog heading, the tag pushed by the maintainer, the archives attached and `go install` reporting the version. Leave `## Next` as the tag step.
2. Pick the version. Below 1.0.0 a dot release (`0.3.1` to `0.3.2`) carries fixes and small additions; a MINOR bump (`0.3.x` to `0.4.0`) carries a layout or API change.
3. In `CHANGELOG.md`, put a `## <version> - <date>` heading under `## Unreleased` and move the entries beneath it. `Unreleased` stays, empty.
4. Check the template. Run `git diff --stat <last-tag> -- internal/initrepo/templates/AGENTS.md.tmpl`. If it prints a change, an existing repository's `AGENTS.md` (which `init` never rewrites) is behind, so the release's `CHANGELOG.md` must carry an entry with the fixed phrase **The `AGENTS.md` template changed: run `akmctl agents`** and say what changed. A change to a skill template needs no phrase: `akmctl skills status` reports it. No test enforces this; the step is the check (261101-hsh4).
5. Run the full check from AGENTS.md: `go vet ./... && go test ./...`, `gofmt -l .`, `/tmp/akmctl check`, `npm run format:check`, `npm run schemas`.
6. Optional dry run of the build: `SKIP_TESTS=1 scripts/release.sh <version>` writes `dist/<version>/` (git-ignored) with the five archives and `SHA256SUMS`.
7. Commit the changelog, the item and the roadmap together, and push `main`. Wait for CI on that commit to pass.

## Tag (the maintainer)

Tags are annotated, with the message `v<version> (unstable pre-1.0)`; a lightweight tag does not match the earlier ones.

```sh
git switch main && git pull --ff-only
git tag -a v<version> -m "v<version> (unstable pre-1.0)"
git push origin v<version>
```

Pushing a `v*` tag is the trigger. The `release` workflow (`.github/workflows/release.yml`) runs `scripts/release.sh`, which runs `go vet` and `go test`, builds the archives for darwin and linux (amd64, arm64) and windows (amd64), and links the version in with `-X .../internal/cli.Version`. It then creates the GitHub release with the archives, `SHA256SUMS` and generated notes. A run takes about two minutes.

`scripts/install.sh` is what the README tells users to pipe to `sh`. It is read from `main` and finds the newest release itself, so a release needs no change to it; a change to the archive names in `release.sh` does.

## Verify and close

1. `gh release view v<version>` lists the five archives and `SHA256SUMS`.
2. `curl -fsSL https://raw.githubusercontent.com/alexsasharegan/akmctl/main/scripts/install.sh | AKMCTL_VERSION=<version> AKMCTL_INSTALL_DIR=$(mktemp -d) sh` installs the archive for this machine.
3. `go install github.com/alexsasharegan/akmctl/cmd/akmctl@v<version> && akmctl version` prints `<version>`.
4. Close the release item with an Outcome that names the checks, take it off the roadmap, and commit.

## Packages (not published yet)

Package manifests are generated from a release, and CI checks the generators; nothing is submitted until the maintainer decides the build has earned it (work item 261108-2ncf).

**WinGet.** After a release, `gh release download v<version> --pattern SHA256SUMS --output SHA256SUMS`, then `scripts/winget-manifests.sh <version> SHA256SUMS out/` writes the three manifests. CI does this for the published v0.3.3 on the Windows runner, runs `winget validate --manifest` and installs from the local files (`winget install --manifest`). To submit, once decided: fork `microsoft/winget-pkgs`, copy the three files to `manifests/a/alexsasharegan/akmctl/<version>/`, and open a pull request; Microsoft's bots validate it and a moderator merges it (the first review takes days). Later versions can be sent with `wingetcreate update alexsasharegan.akmctl --version <version> --urls <zip URL> --submit`, which needs a GitHub token. The maintainer submits; nothing in this repository does.

**Homebrew.** `scripts/homebrew-formula.sh <version> SHA256SUMS Formula/akmctl.rb` writes the formula for the four macOS and Linux archives. CI packages the published v0.3.3 and installs it from a local tap (`brew tap-new --no-git local/akm`, copy the file into its `Formula/` folder, `brew install local/akm/akmctl`) on the macOS and Ubuntu runners, then runs `akmctl version`. To publish, once decided: the maintainer creates the repository `alexsasharegan/homebrew-tap`, commits the generated `Formula/akmctl.rb` for each release (users then run `brew install alexsasharegan/tap/akmctl`), and, to automate it, adds a token with write access to that repository as an Actions secret for a `release.yml` step. No credential belongs in this repository. **Never run `brew test` or `brew audit` on a maintainer's Mac.** They write into `/opt/homebrew/Library/Homebrew` (they install development gems into its tracked `vendor/bundle`) and left `brew` unable to run some commands (work item 261108-43gk). CI does not run them either. To check the formula beyond the install, use a CI job, a container or a throwaway user.

## If it goes wrong

- The `publish` step fails with an HTTP 5xx from GitHub and no release exists (`gh release view v<version>` says not found): re-run the failed job (`gh run rerun <id> --failed`). It reuses the same tag; v0.3.2 needed this once.
- The workflow fails for another reason after the tag is pushed: fix forward and release the next dot version. Do not move or delete a pushed tag; `go install` and the module proxy cache it.
- A push of the tag fails with `Permission denied (publickey)`: the SSH agent needs unlocking or approving; nothing has been published.
