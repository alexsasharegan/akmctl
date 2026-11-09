---
type: design
status: draft
created_at: 2026-11-08
updated_at: 2026-11-08
authored_by: agent
---

# Test strategy and coverage audit

Draft for [261030-456q](../work/done/261030-456q-test-strategy-and-coverage-audit.md). Four work items write it: the strategy section below ([261108-dzh4](../work/done/261108-dzh4-audit-test-strategy.md)), a maintainability section ([261108-x2c4](../work/done/261108-x2c4-audit-test-maintainability.md)), a coverage section ([261108-arn2](../work/done/261108-arn2-measure-test-coverage.md)), and the closing item that numbers the findings, groups them, files the tickets and orders the roadmap ([261108-kmwm](../work/done/261108-kmwm-write-audit-report-and-file-tickets.md)). The findings are numbered F1 to F19 and grouped in the closing sections; the strategy findings were candidates C1 to C8, the maintainability ones M1 to M5 and the coverage ones V1 to V6.

Audit date: 2026-11-08, on `main` at the merge of pull request 35. Counts come from the commands named in each place.

## Strategy

### What each layer is for

| Layer                     | Where                                                                                                                                                                                                                            | What it is for                                                                                                                                             | Right layer?                                                                       |
| ------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| Unit tests                | 14 packages under `internal/`, 71 test files (`func Test` counts per package: api 145, body 31, cli 30, server 30, check 27, migrate 24, initrepo 16, frontmatter 13, record 13, depgraph 5, fsx 3, mdlinks 3, index 2, width 2) | Behavior of one package, with temporary copies of a sample repository                                                                                      | Yes. About 11,100 lines of tests against 14,300 of code.                           |
| Conformance fixtures      | `docs/spec/akm-api/.fixtures/cases` (109 files), copy in `testdata/api`                                                                                                                                                          | One contract for the API: a request, the repository it runs on, the full expected response                                                                 | Yes for the API. Runs only in process through `api.Run` (see F1).                  |
| YAML-lite oracle and fuzz | `internal/frontmatter` (`oracle_test.go`, `fuzz_test.go`, `testdata/yaml-lite.json`), `internal/body` (`tree_test.go`, `battery_test.go`)                                                                                        | The hand-written frontmatter parser is compared with `go.yaml.in/yaml/v3` (a test-only dependency); the body tree has fuzz targets                         | Yes. Fuzz targets run only their seeds under `go test`; nothing runs `-fuzz` (F8). |
| CLI binding               | `internal/cli` (6 files)                                                                                                                                                                                                         | Flag parsing, text output, exit codes                                                                                                                      | Yes, but see F1: it does not run the fixtures.                                     |
| HTTP binding              | `internal/server` (`server_test.go` and four more)                                                                                                                                                                               | Auth, host and origin checks, routes, `If-Match`                                                                                                           | Yes, same note.                                                                    |
| Workbench JavaScript      | `internal/server/ui_assets_test.go` runs `testdata/logic.test.mjs` under Node                                                                                                                                                    | The pure modules of the page                                                                                                                               | Partly: 14 of 25 modules (F4).                                                     |
| Skill and template tests  | `internal/initrepo/initrepo_test.go`, `optional_test.go`, `internal/cli/skills_test.go`                                                                                                                                          | The two copies of each skill agree; names follow the grammar; every operation has a reference; `init` writes a repository that passes `check`              | Yes for drift; nothing tests what the text tells an agent (F6).                    |
| Records check             | CI step `akmctl check`, `npm run format:check`, `npm run schemas`                                                                                                                                                                | This repository is its own test case; the schemas validate the fixtures                                                                                    | Yes. Linux only, which is enough.                                                  |
| CI                        | `.github/workflows/ci.yml`                                                                                                                                                                                                       | `go vet`, `go test ./...` on Linux, macOS and Windows; `gofmt`, records, format, schemas and `go test -race` on Linux only; Windows is `continue-on-error` | See F7.                                                                            |

### Findings

Each has an area, a severity (gap, risk or nuisance), the evidence and a recommended fix.

- **F1. The bindings do not run the contract.** Area: conformance fixtures. Severity: gap. Evidence: `grep -rn "testdata/api/cases" internal` finds one reader, `internal/api/conformance_test.go`, which calls `api.Run`. The CLI and the HTTP server have their own hand-written tests, and `docs/architecture.md` states that "the CLI and the HTTP binding return the same JSON envelope for the same request"; I found no test that runs one request through both and compares. Fix: run every case file through the HTTP handler and, for the reads, through the CLI with `--format json`, and compare with `expected`.
- **F2. Operations with no fixture.** Area: conformance fixtures. Severity: gap. Evidence: `ls docs/spec/akm-api/.fixtures/cases` has no `commits-*`, `history-*` or `skills-*` file, and the only `new` case (`new-1`) is a refusal. The reason is real: ids are random (API spec section 10) and a fixture directory cannot hold a git history. Fix: an ADR on how to fixture these (a fixed clock and id source for `new`; a fixture that builds a git history from a script for `commits` and `history`); see the draft ADR list in the closing item.
- **F3. Tests that need an outside tool skip without a word.** Area: tooling. Severity: risk. Evidence: `t.Skip` in `internal/server/ui_assets_test.go` (no `node`), four test files in `internal/api` (no `git`; two more use it), `internal/check/safety_test.go` (`AKM_CORPUS` not set, a corpus that no job sets). A runner without Node or git would pass with the workbench and every git feature untested. Fix: in CI, fail when the tool is missing (an environment variable such as `AKM_REQUIRE_TOOLS=1` turns each skip into a failure), and decide whether the corpus test belongs in CI.
- **F4. The workbench has no test for its page code.** Area: workbench. Severity: gap. Evidence: `ls internal/server/web/js` lists 25 modules; `TestPureModulesUnderNode` loads 14 of them (`util logic store live moves project review capture controls effects fields picker toasts replay`). The 11 modules that touch the DOM (`adrs anim board boot check detail flow rows tables ui wizard`) have no automated test; the pull requests that changed them say what was checked by eye in the browser pane. Fix: a DOM test with a small fake `document` for the render functions, or a recorded browser run; at least a smoke test that `boot.js` renders the board against a fixture repository.
- **F5. The scripts that users run have no test.** Area: tooling. Severity: risk. Evidence: `scripts/install.sh` (73 lines, piped from `curl` to `sh` by the README) and `scripts/release.sh` (56 lines) appear in no test; `release.yml` runs only on a tag, so a break in the release script shows up when a release is cut. `cmd/akmctl` also has no test file. Fix: a CI job that runs `scripts/release.sh` for a dummy version and runs `install.sh` against the produced archive in a temporary directory.
- **F6. What the skills tell an agent is not tested.** Area: skills. Severity: gap. Evidence: the tests above check copies, names and operation references. The text is checked only by hand: the cold-agent trial of 261101-baze and its retest in 261101-t8h8, which found two tool snags (261108-h7ds, 261108-kdz1) that no test had pinned. Fix: keep a small scripted cold-agent check as a recorded procedure, and add a test for each snag a retest finds.
- **F7. CI proves less than it looks like on two platforms.** Area: platform. Severity: risk. Evidence: `ci.yml` runs `go test -race` only when `runner.os == 'Linux'`; `AGENTS.md` says to add `-race` on Linux or macOS. The Windows job has `continue-on-error`, and it fails on every pull request (known; Windows debugging is on hold), so a Windows-only regression is invisible among the existing failures. Fix: add `-race` on macOS; when Windows work resumes, make the known failures a named skip list so a new failure shows.
- **F8. Fuzz targets only run their seeds.** Area: tooling. Severity: nuisance. Evidence: `grep -rln "func Fuzz" internal` finds `internal/frontmatter/fuzz_test.go` and the body tests; CI runs plain `go test`. Fix: a scheduled or manual job with `-fuzz` for a minute per target.

### The seed list, checked

| Entry in 261030-456q                                                                                                                                       | Verdict                                | Evidence                                                                                                                                                                                                                                                                                 |
| ---------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `new` has no exact fixture because ids are random; `commits` has none because a fixture directory cannot hold a git history                                | Confirmed, and extended                | `new-1` is a refusal. `commits`, `history` and `skills` have no case file (F2).                                                                                                                                                                                                          |
| `similar`, `dependents` and `commit` on `new` are covered by Go tests only                                                                                 | Corrected                              | `similar` is Go-only (`similar_test.go`, `cli_test.go`). `dependents` is in 6 fixture cases and the `commit` member in 9 (`grep -l dependents docs/spec/akm-api/.fixtures/cases/*.json`).                                                                                                |
| Two copies of the API fixtures; `make_fixtures.py` knows two repositories, so a wrong generator and a wrong implementation can agree                       | Corrected                              | One source since 261029-x4ws, with a test that fails when the copy differs. The generator no longer reproduces the cases (18 files differ in a scratch run) and is marked a frozen reference, so it is not the contract.                                                                 |
| The workbench JavaScript is tested through Node from a Go test; whether it skips silently and how much it covers is unknown                                | Confirmed                              | It skips silently when `node` is missing (F3) and covers 14 of 25 modules (F4).                                                                                                                                                                                                          |
| The skills are tested for drift, not for what they tell an agent                                                                                           | Confirmed                              | F6.                                                                                                                                                                                                                                                                                      |
| `TestCompareAndSwapRace` is flaky on Windows CI (261029-q8z3); `-race` runs only on Linux and macOS; the Windows path beyond CI is uncovered (261028-jzsx) | Partly corrected, partly not rechecked | In CI `-race` runs on Linux only, not macOS (F7). The flaky test lives in `internal/api/write_test.go`; I did not rerun it on Windows, which is on hold. The Windows path is still uncovered, and the whole Windows job fails.                                                           |
| Whether every `check` rule id has a case that trips it and a case that does not is not recorded anywhere                                                   | Partly corrected                       | The rules table of the spec lists 44 rule ids, and every one is named in at least one `internal/check` test or `check-*` fixture (a script over the spec table and those files). Whether each has both a tripping and a non-tripping case is not recorded; that matrix is 261108-arn2's. |

## Maintainability

Measured on 2026-11-08, `main` at 75c8b7d. How hard each suite is to change, and how long a run takes. Findings are numbered F9 to F13 (they were candidates M1 to M5).

### Run times

| Run                           | Command                                    | Result                                                                                                                       |
| ----------------------------- | ------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| Local, all packages           | `go clean -testcache; go test ./...`       | 22.0 s real. `internal/api` takes 19.4 s of it; the next is `internal/server` at 3.4 s.                                      |
| Local, with the race detector | `go clean -testcache; go test -race ./...` | 33.8 s real. `internal/body` takes 31.8 s (1.6 s without `-race`); `internal/api` 22.6 s.                                    |
| CI, last three runs on `main` | `gh run view <id> --json jobs`             | Linux 73, 108 and 73 s; macOS 49, 54 and 27 s; Windows 101, 115 and 104 s (the Windows job fails, see the strategy section). |

The slowest single test is `TestWritesAndDryRunsRace` at 1.4 s. `internal/api` is slow as a sum: 145 tests and 121 subtests, each of which copies a sample repository into a temporary directory (`user` 9.7 s but `sys` 13.0 s for the whole run: the time is file copying). No test calls `t.Parallel`; packages run in parallel, tests inside a package do not.

### Helpers

Seven places copy a sample repository to a temporary directory: `copyRepo` and `copyTree` in `internal/api`, `copyDir` in `internal/cli`, `copyTree` and `vaultCopy` in `internal/migrate`, `repoWith` in `internal/check`, and an inline walk in `internal/server`. Most discard the errors of `WalkDir`, `ReadFile` and `WriteFile`, so a failed copy gives an empty repository and a confusing failure later. 15 test files call `filepath.WalkDir`. Three snapshot helpers (`snapshot` in `internal/api` and in `internal/initrepo`, and `snap`) read a tree into a map. There is no shared test-helper package.

### Real time, git and Node

One test uses the real clock (`internal/api/dry_test.go` builds a 48-hour-old time); `check` tests pass a fixed `AsOf`, and the write tests pass `as_of`. Tests that need `git` build a real repository with `git init` in a temporary directory (four files skip without it, see F3). The single Node test (`TestPureModulesUnderNode`) copies 14 modules and a test file to a temporary directory and runs them; it does not appear among the slow tests.

### What breaks when a template or a fixture changes

Measured in a throwaway worktree, one edit at a time, running `go test -count=1 ./...`:

| Edit                                                     | Tests that fail                                                                                                                                    |
| -------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| Reword one sentence in `AGENTS.md.tmpl`                  | None.                                                                                                                                              |
| Rename the `Commands` heading in `AGENTS.md.tmpl`        | `TestCustomDocsDirEndToEnd` (`internal/api`), `TestOptionsCheckClean` and `TestWorkLifecycle` (`internal/initrepo`): the rules check the headings. |
| Reword the "Item too big" heading in one copy of a skill | `TestSkillsMatchRepository`: by design, the two copies must match.                                                                                 |
| Add one reference file to both skill copies              | Six `init` conformance cases (`init-1`, `-2`, `-3`, `-6`, `-7`, `-8`), whose written-file lists are fixed text.                                    |

Of the 109 case files, 38 embed an etag. The etag of one sample work item (`261020-a2b3`) appears in 9 case files, so a change to that file changes 9 cases; I found no tool that rewrites them. During work item 261031-72bd the new `close` members meant editing six case files, which I did with a one-off script.

### Findings

- **F9. No update mode for the expected responses.** Area: conformance fixtures. Severity: risk. Evidence: 38 of 109 cases embed an etag, 9 embed one item's etag, and a response change in 261031-72bd needed six case files edited by a throwaway script. The generator is a frozen reference (261029-x4ws), so the cases are hand-kept. Fix: a `-update` flag on `TestConformance` that rewrites `expected` from the tool's own output for a human to review in the diff, or a placeholder for the etag that the runner fills in.
- **F10. Adding a file to the templates breaks six `init` cases by hand.** Area: conformance fixtures, templates. Severity: nuisance. Evidence: the experiment above. Fix: the same update mode (F9) covers it; or make those cases check the set of written files by pattern.
- **F11. Seven copies of a repository-copy helper and three of a snapshot helper.** Area: tooling and maintainability. Severity: nuisance. Evidence: the helper list above, and the copies ignore their own errors. Fix: one `internal/testutil` package (`CopyTree` that fails the test on error, `Snapshot`, `MustOK`), used by all packages.
- **F12. `internal/api` is serial and file-copy bound.** Area: tooling. Severity: nuisance. Evidence: 19.4 s of a 22.0 s run; `sys` 13.0 s; no `t.Parallel` anywhere. Fix: `t.Parallel` on tests that copy their own repository, and a copy that links instead of copying where the test only reads.
- **F13. The race run is dominated by one package.** Area: tooling. Severity: nuisance. Evidence: `internal/body` takes 31.8 s with `-race` and 1.6 s without; the full race run is 33.8 s. Fix: lower the iteration count of the battery under the race detector (`testing.Short` or a build tag), as CI runs it on every pull request.

## Coverage

Measured on 2026-11-08, `main` at 955895e (pull request 37). Every figure comes from a command in the Log of [261108-arn2](../work/done/261108-arn2-measure-test-coverage.md). Statement coverage says a line ran, not that its result was checked; it finds code no test reaches.

### Per package

`go clean -testcache; go test -count=1 -cover ./...` counts only the tests of the package itself. `go test -count=1 -coverpkg=./... -coverprofile=cover.out ./...` counts every test that reaches a package; its blocks were merged by file and range, since each test binary lists every block (my first merge double-counted, and the total was cross-checked with `go tool cover -func`).

| Package                | Own tests | All tests                   | Statements |
| ---------------------- | --------- | --------------------------- | ---------- |
| `cmd/akmctl`           | 0.0 %     | 0.0 %                       | 5          |
| `internal/api`         | 88.0 %    | 88.4 %                      | 2,615      |
| `internal/body`        | 87.2 %    | 95.4 %                      | 438        |
| `internal/check`       | 89.2 %    | 93.0 %                      | 629        |
| `internal/cli`         | 68.2 %    | 68.2 %                      | 1,023      |
| `internal/depgraph`    | 75.3 %    | 97.5 %                      | 158        |
| `internal/frontmatter` | 84.6 %    | 88.6 %                      | 449        |
| `internal/fsx`         | 38.0 %    | 68.0 %                      | 50         |
| `internal/index`       | 96.1 %    | 96.1 %                      | 51         |
| `internal/initrepo`    | 78.7 %    | 89.1 %                      | 202        |
| `internal/mdlinks`     | 74.2 %    | 98.5 %                      | 132        |
| `internal/migrate`     | 71.5 %    | 86.3 %                      | 779        |
| `internal/record`      | 70.9 %    | 95.7 %                      | 302        |
| `internal/server`      | 92.0 %    | 92.0 %                      | 226        |
| `internal/width`       | 96.6 %    | 96.6 %                      | 59         |
| **All**                |           | **86.9 %** (6,184 of 7,118) |            |

Fourteen functions are never run by any test: `main`, `usage`, `runUI` and `openBrowser` (a browser window cannot be opened in a test), `progressText` (the text view of a parent's progress, added by 261031-72bd), `scopeMatch`, `unreadable`, `migrate.Plan`, `record.SetMetaDir`, `present`, `choice`, and three `Error` methods. The weakest functions that do run are `cli.text` and `cli.recordsOf` (30 %: the table view of every command; the tests use `--format json`) and `api.matchesFilters` (54 %).

### Operations against the fixtures

The 25 operations take 120 request parameters in all (`api.Params`). The 109 fixture cases name 73 of them (61 %). The cells that no fixture reaches:

| Operation                | Fixture cases | Parameters no case names                                      |
| ------------------------ | ------------- | ------------------------------------------------------------- |
| `ls`                     | 7             | `blocked`, `cycle`, `limit`, `offset`, `scope`, `since`       |
| `search`                 | 1             | `body`, `fields`, `limit`, `offset`, `status`, `type`         |
| `claim`                  | 3             | `dry_run`, `if_match`, `takeover`                             |
| `log`                    | 5             | `dry_run`, `if_match`, `section_etag`                         |
| `release`                | 1             | `dry_run`, `if_match`                                         |
| `tick`                   | 6             | `done`, `if_match`, `section`                                 |
| `remove`                 | 4             | `if_match`, `section`, `section_etag`                         |
| `close`, `section`       | 9 each        | `if_match`                                                    |
| `migrate`                | 4             | `acknowledge_backup`, `as_of`, `target`                       |
| `new`                    | 1             | `as_of`, `title`                                              |
| `config`, `skills`       | 10 and 0      | `confirm`, `confirm_edited` (`skills` also `dry_run`, `mode`) |
| `bases`, `board`, `init` | 4, 1, 9       | `include_content`; `statuses`; `git_init`                     |
| `commits`, `history`     | 0             | all of them (`range`, `limit`, `item`)                        |

`if_match`, the compare-and-swap check of section 7, is named in a fixture only for `set` and `tick`'s stale case; the eight other writes have no stale case. Error codes: the spec lists eight (`invalid_request`, `unauthorized`, `not_found`, `ambiguous`, `stale`, `refused`, `unreadable`, `internal`); the fixtures expect four of them (`invalid_request` 12, `refused` 20, `not_found` 4, `stale` 2). The rest are covered by Go tests (`unauthorized` by the server tests, `ambiguous` and `unreadable` in three and four test files). Of the 18 `refused` reasons the spec names, 12 have a fixture; `extension_off`, `if_match_required`, `needs_confirmation`, `no_manifest`, `not_git` and `ticked_box` are in Go tests only.

Go tests cover more than the fixtures do, so a cell above is a gap in the contract, not necessarily in the code. The code gaps are the ones in the profile: lines 271 to 277, 285 to 287 and 290 to 296 of `internal/api/ops.go` are the bodies of the `ls` filters `blocked`, `scope` and `since`, and no test in the repository runs them. They are in the API spec (section 4.2), so a reader can rely on three filters that nothing checks.

### Rules against the tests

The spec's rules table lists 44 rule ids. The code emits 43 of them from `internal/check`; `S3` comes from `internal/record` (reading the manifest) and is tested in `internal/record` and `internal/api`. `TestRules` has 62 rows covering 38 rules, each a case that raises exactly that rule; `CL1`, `CL2` and `SF1` to `SF4` are raised in their own tests (`TestClaims`, `safety_test.go`). So every rule has a test that raises it. What no test pins is the opposite for a single rule: a near miss that must not raise it. `TestCleanRepository` shows the clean sample raises nothing, which is one negative case for all rules at once; per-rule near misses exist only where a bug put one there (`TestLinksThatResolve`, `TestParentListingItsChildInBlockedByIsNotACycle`, the dot-folder test).

### What the spec says it does not test

- **Section 10, "Fixtures for `commits` and `history`".** A fixture directory cannot hold git history; an implementation builds repositories with `git`. `TestCommits*` and `TestHistory*` do that. No shared contract exists, so a second implementation has nothing to run.
- **Section 10, "Exact fixtures for `new`".** Ids are random, so `new` is tested for shape only; the one `new` case is a refusal.
- **Section 10, "More write fixtures".** A case for each refusal reason and for the wikilink reference style: the 6 reasons above, and no wikilink-style case.
- **Section 7, "Compare and swap".** "An editor that writes the file between the final check and the rename is not covered"; the window is one rename. Not testable without a fault-injection hook.
- **Section 7, "Atomicity".** The link repair after `close` writes other files one at a time and is not rolled back (work item 261029-xtqj). The failure path is covered by one test of `link_errors` that needs a folder that refuses writes, which it skips as root and on Windows.

### Findings

- **F14. Three `ls` filters have no test at all.** Area: API operations. Severity: gap. Evidence: `go tool cover` shows no execution of `matchesFilters` lines 271 to 277 (`blocked`), 285 to 287 (`scope`) and 290 to 296 (`since`); `scopeMatch` is at 0 %. Fix: three fixture cases and a Go test for the path-prefix rule of `scope`.
- **F15. The text view of the CLI is untested.** Area: CLI. Severity: gap. Evidence: `cli.text` 30 %, `cli.recordsOf` 30 %, `progressText` 0 %; `internal/cli` is at 68 %, the lowest of the large packages. People and agents read this view; the tests read `--format json`. Fix: golden-output tests for `ls`, `show`, `close`, `split` and the commit message block.
- **F16. 47 of 120 request parameters are never named in a fixture.** Area: conformance fixtures. Severity: gap. Evidence: the table above; `if_match` has no stale case for eight writes, `dry_run` none for `claim`, `log` and `release`. Fix: one case per cell where the behavior is in the contract, starting with the stale case of every write.
- **F17. Six `refused` reasons exist only in Go tests.** Area: conformance fixtures. Severity: nuisance. Evidence: `extension_off`, `if_match_required`, `needs_confirmation`, `no_manifest`, `not_git`, `ticked_box`. Fix: the "More write fixtures" item of section 10.
- **F18. No near-miss case per rule.** Area: `check`. Severity: nuisance. Evidence: the rules section above. Fix: for each rule, a case that must not raise it, next to the case that does.
- **F19. `cmd/akmctl` and `usage` are untested.** Area: CLI. Severity: nuisance. Evidence: `main` and `usage` at 0 %; the 14 zero-coverage functions include `migrate.Plan` and `record.SetMetaDir`, which may be dead code. Fix: a smoke test that runs the built binary with `version` and `help`; delete dead code or test it.

## Findings by group, and the tickets

The 19 findings above are numbered `F1` to `F19`: F1 to F8 from the strategy section (were C1 to C8), F9 to F13 from maintainability (were M1 to M5), F14 to F19 from coverage (were V1 to V6). Each group has one work item; the item names its findings and this report names the item. The order is the one recommended for the roadmap: cheap risk reduction first, the shared tools next, the design-heavy pieces last.

| Order | Work item                                                                                                                        | Findings               | Why this place                                                                          |
| ----- | -------------------------------------------------------------------------------------------------------------------------------- | ---------------------- | --------------------------------------------------------------------------------------- |
| 1     | [261108-z9n8](../work/done/261108-z9n8-ci-hardening-for-skipped-tests-and-scripts.md) CI hardening for skipped tests and scripts | F3, F5, F7, F8         | Risks (F3, F5, F7); small, and later work is only trusted if a skip cannot pass quietly |
| 2     | [261108-3kse](../work/done/261108-3kse-fixture-update-mode-and-gap-cases.md) Fixture update mode and gap cases                   | F9, F10, F14, F16, F17 | The update mode (F9) makes every later fixture cheap to add                             |
| 3     | [261108-qv1w](../work/done/261108-qv1w-cli-and-check-test-gaps.md) CLI and check test gaps                                       | F15, F18, F19          | Independent of the fixtures; the CLI text is what people and agents read                |
| 4     | [261108-rwqj](../work/done/261108-rwqj-shared-test-helpers-and-run-time.md) Shared test helpers and run time                     | F11, F12, F13          | Nuisances; pay off before the suite grows more                                          |
| 5     | [261108-f1tv](../work/done/261108-f1tv-run-fixtures-through-both-bindings.md) Run fixtures through both bindings                 | F1, F2                 | Larger and needs an ADR first; benefits from 2                                          |
| 6     | [261108-79t5](../work/done/261108-79t5-skill-text-checks.md) Skill text checks                                                   | F6                     | Needs a decision on where the cold-agent procedure lives                                |
| 7     | [261108-xajw](../work/done/261108-xajw-workbench-page-tests.md) Workbench page tests                                             | F4                     | Largest effort; the workbench is not promised for 1.0.0                                 |

### ADR drafts the findings call for

Listed, not decided. Each would be a `proposed` ADR written inside the item named.

1. **How to fixture non-deterministic and git-dependent operations** (F2, in 261108-f1tv): a fixed clock and id source for `new`; a fixture that builds a git history from a script for `commits` and `history`; whether `skills` needs a fixture at all.
2. **What the bindings must prove** (F1, in 261108-f1tv): whether the contract is that the CLI and the HTTP server return the same envelope as `api.Run`, and which operations the CLI is allowed to differ on.
3. **Whether a skipped test may pass** (F3, in 261108-z9n8): the policy that a missing outside tool fails CI, and whether the `AKM_CORPUS` test belongs in CI.
4. **Where the cold-agent check lives** (F6, in 261108-79t5): a recorded procedure, a script, or a scheduled job; and what it costs to run.

### Seed list

The checked seed list of [261030-456q](../work/done/261030-456q-test-strategy-and-coverage-audit.md) is the table "The seed list, checked" above, one row per entry. Its last row, the rule-id matrix, is answered by the Coverage section: every rule is raised by a test; near misses are F18.
