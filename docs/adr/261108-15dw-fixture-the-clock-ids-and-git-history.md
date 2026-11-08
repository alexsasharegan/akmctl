---
id: "261108-15dw"
title: Fixtures fix the clock, ids and git history with replayed inputs
summary: Conformance fixtures bind random ids and git hashes with matchers and replay git history from the case, so any binding in any language can run them
status: proposed
date: 2026-11-08
---

# Fixtures fix the clock, ids and git history with replayed inputs

## Context

The conformance fixtures (`docs/spec/akm-api/.fixtures`) run only through `api.Run` (finding F1 of the test-strategy audit), and four operations have none: `new` (a random id), `commits` and `history` (they read git), and `skills` (its markers carry the tool version and file hashes). To run every case through the HTTP handler and the reads through the CLI (work item 261108-f1tv), and to add those four operations, the fixtures need a way to fix what is not a function of the request: the clock, the random part of an id, the git history and the text of a shipped skill.

The fixtures are the contract for any implementation of the API, in any language (the validator is Node). A Go-only seam (a package variable for the random source, as `migrate` has) would fix the values in one implementation and leave the fixture unusable by the others, so the expected response has to say what must hold, not which random draw the Go code makes.

## Considered Options

- **Go test seams** for the random source, the clock and git: simple, but only the Go tests can use them; the fixtures would stay unreadable for another implementation and the HTTP and CLI bindings, run as processes, could not reach the seam.
- **A `seed` or `now` request parameter**: fixes everything in any language, but adds a public parameter to the API that exists only for tests, and `as_of` already is the clock.
- **Matchers in `expected` plus inputs the case replays**: the case carries the git history as data and the expected response marks the values it does not fix. Larger fixture format, no change to the API, any runner can implement it.
- **A frozen `.git` directory in a fixture repository**: hashes are then fixed, but a repository inside a repository is awkward (git ignores nested `.git` contents on add), and the hashes depend on the git version's object format.

## Decision

Matchers in `expected` plus replayed inputs.

1. **Clock.** Unchanged: a write case passes `as_of`. The HTTP and CLI runners pass it through (body or query, `--as-of`).
2. **Random ids.** A string in `expected` (or in `after` paths and texts) may be `$id(YYMMDD)`: it matches an id of that date with four characters from the id alphabet. The first occurrence binds the actual value and every other occurrence of the same marker in the case means that value, so the file name, `id`, `title` links and the commit message agree.
3. **Anything not fixed.** `$any` matches any value of any type but must be present. It stands for shipped text and tool-dependent members: the `diff` of a skills plan, a marker's `tool` and hashes, a version string.
4. **Git history.** A case may carry `git`: a list of commits, each `{message, date, files}` (`files` maps a path to its new text, `null` deletes it). The runner replays them with `git init -b main`, fixed author and committer name and email, `GIT_AUTHOR_DATE` and `GIT_COMMITTER_DATE` from the case, `commit.gpgsign` off and `core.autocrlf` false. Hashes are not fixed: `$hash` matches 7 to 40 hex digits and binds like `$id`. A runner without git skips these cases and says so (CI sets `AKM_REQUIRE_TOOLS=1`).
5. **Skills.** A `skills` case pins `state`, `clean`, the plan's `action` and `path` and `files` states, and uses `after_paths` for what a write creates; text and hashes are `$any`. Changing a shipped skill then does not break a fixture.
6. **Bindings.** A case runs through `api.Run`, the HTTP handler and, for read operations, the CLI with `--format json`, and the three responses must equal `expected` after matcher binding. The runner holds the route table of the specification (section 6) and the exit-code table. A write case without `if_match` gets the record's current etag from a `GET` first, since the HTTP binding refuses a write without one; a case may set `bindings` to restrict itself (for example `["http"]` for the `if_match_required` refusal, which has no API form), with the reason in `name`.

## Consequences

- The fixture format grows three matchers and two case members (`git`, `bindings`); `docs/spec/akm-api/README.md` and the spec section on conformance say so, and the validator accepts them.
- `new`, `commits`, `history` and `skills` get fixtures, and so does the HTTP-only refusal `if_match_required`.
- A runner in another language implements the matchers and, to run `commits` and `history`, git. Cases that need git can be skipped with a notice.
- `-update` (work item 261108-3kse) must keep matchers it finds in `expected` instead of replacing them with the actual values.
- Results that differ between bindings are real defects of one binding; they will surface when the runners are first turned on, and each is filed as its own work item.

## Confirmation

The conformance test runs every case through each binding; a case whose `$id` or `$hash` does not bind, or whose binding returns a different envelope, fails.

## Review

Proposed by the agent on 2026-11-08; the maintainer decides.

## Revert

Remove the matchers and the `git` and `bindings` members, drop the cases that use them and supersede this record with the option chosen instead.
