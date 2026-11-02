# Review policy

`akmctl` edits a repository's Markdown records for people and agents who also edit them by hand. A defect here corrupts someone's notes or lets a record and the tool disagree. Review for those failures. The formatter, `go vet`, Prettier, the records check, the schemas and the tests already run in CI; do not repeat them.

Report only what the diff introduces or touches. Report older debt only inside a function the diff changes, and then as a Note.

## Severity

- **Blocking**: data loss or corruption, a broken invariant below, a trust boundary opened, or a change a user or client can see that lacks its record (see Records).
- **Note**: older debt inside a function the diff changes. Nothing else.

Each Blocking rule is marked _(behavior)_ or _(process)_. A _(behavior)_ finding names the file and line, the input or state, and the wrong result. A _(process)_ finding names the missing artifact or the rule broken, and needs no failure scenario. If you cannot do either, stay silent. Do not report style, naming, comments, Windows-only behavior, or generated and copied files that a test or generator owns.

## Invariants

1. **Bytes outside an edit never change.** _Blocking (behavior)._ A write is a splice into the span it edits. Unknown frontmatter keys, comments, quoting and line endings survive. Flag a write that re-serialises a record.
2. **One scanner, one fence rule.** _Blocking (behavior)._ In work-item and ADR bodies, heading, section and fence detection live in `internal/body` (`FenceTracker`, `Locate`). Flag a new loop outside it that finds them by testing line prefixes (` ``` `, `~~~`, `## `); name a body, such as a four-backtick fence, where it disagrees with `FenceTracker`.
3. **Fail closed.** _Blocking (process)._ When the tool cannot read a construct it refuses, names it and writes nothing. A change that makes the scanner or the YAML-lite parser accept a new construct needs `FuzzAgainstOracle` (`internal/body`, `internal/frontmatter`) and the oracle battery tests (`internal/body/battery_test.go`, `internal/frontmatter/oracle_test.go`) to pass, and the work item's Outcome to name the verification run. A unit test alone is not enough.
4. **Writes: lock, reload, check any supplied `if_match`, then write atomically.** _Blocking (behavior)._ A write path takes `fsx.LockRepo`, reloads, checks the etag when one was supplied, and writes through `fsx.WriteAtomic`. Flag a write to a file inside the repository root by `os.WriteFile`, `os.Rename` or `os.Create` outside `internal/fsx`, except in `internal/migrate`, `internal/api/ops_migrate.go` and `internal/initrepo`. A file outside the root, such as the `--out` handoff file, is not covered.
5. **A multi-step write has a defined failure state.** _Blocking (behavior)._ A new write that can stop half-way needs a rollback or an error that says which files changed.
6. **Spec, `check` and the write path agree.** _Blocking (behavior)._ A rule the spec says an operation refuses, or that an operation could itself make false, must be refused by that operation as well as reported by `check`. A `check` warning for hand edits needs no write-side refusal.
7. **Standard library only in the binary.** _Blocking (process)._ `go.yaml.in/yaml/v3` stays test-only; a new non-test import outside the standard library breaks the rule.
8. **The server stays on loopback.** _Blocking (behavior)._ A non-loopback bind, a weakened Host, Origin or token check, a request-supplied id or path that can resolve outside the repo root, or a secret or private path in a response, log or fixture.
9. **Closed work items stay closed.** _Blocking (behavior)._ An operation that edits a closed work item's fields or body refuses with `reason: "closed"` unless the call carries `correction`. Only `set` (title and summary) and `section` (not the Log) accept one. `keep`, `close`, `claim`, `release` and `migrate` are exempt (ADR 261101-dxcn).

## Records

_Blocking (process)._ A change a user or client can see needs, in the same change:

- an entry under `## Unreleased` in `CHANGELOG.md`;
- a `docs/glossary.md` entry for each new feature, command, field, rule or workbench part;
- an ADR for a hard-to-reverse choice or one that rejects a real alternative, drafted as `proposed`.

Visible means a command, the API, the spec, a rule, a conformance fixture (`docs/spec/akm-api/.fixtures`, `testdata/api`), the workbench, the install script, or the bytes a write produces. A refactor, or a test-only fixture (`testdata/migrate`, `testdata/yaml-lite.json`), needs none unless it changes what a user or client sees.

Edit either conformance fixture tree alone and the change is Blocking until work item 261029-x4ws makes a test enforce it.

An agent may draft an ADR as `proposed`. Setting it to `accepted` needs the maintainer's explicit authorization, recorded in the implementing work item's Log; the status change itself may be in the implementing commit. Never edit the body of an `accepted` ADR's decision; supersede it. These edits are allowed: `status`, `superseded_by`, an errata section, typos, and a scrub of private material (names, paths, identifiers) that says so in the commit message.

## Out of scope

Commit trailers, claims and the one-item-per-branch rule (`akmctl check` and `akm-commit` cover them), and workbench styling with no data, security or accessibility effect.

## Output

Start with `N actionable findings` (or `No actionable findings.`). Group under **Blocking** then **Note**. Each finding: `file:line`, the failure scenario or the missing artifact, why the diff causes it, and a fix. Do not invent findings.
