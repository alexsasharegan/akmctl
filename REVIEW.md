# Review policy

`akmctl` edits a repository's Markdown records for people and agents who also edit them by hand. A defect here corrupts someone's notes or lets a record and the tool disagree. Review for those failures. The formatter, `go vet`, Prettier, the records check, the schemas and the tests already run in CI; do not repeat them.

Report only what the diff introduces, or what the diff touches or builds on. Mention older debt only in a file or function the diff changes, and then as a Note.

## Severity

- **Blocking**: data loss or corruption, a broken invariant below, a trust boundary opened, or a change a user or client can see that lacks its record (see Records).
- **Note**: a test gap, or older debt the diff touches. Nothing else.

A Blocking finding names the file and line, the input or state, and the wrong result. If you cannot trace a failure, stay silent. Do not report style, naming, comments, Windows-only behavior, or generated and copied files that a test or generator owns.

## Invariants

1. **Bytes outside an edit never change.** A write is a splice into the span it edits. Unknown frontmatter keys, comments, quoting and line endings survive. Flag a write that re-serialises a record, or a new construct with no round-trip test.
2. **One scanner, one fence rule.** Heading, section and fence detection live in `internal/body` (`FenceTracker`, `Locate`). Flag a new loop outside it that tests for ` ``` `, `~~~` or `## ` prefixes.
3. **Fail closed.** When the tool cannot read a construct it refuses, names it and writes nothing. A change that makes the scanner or the YAML-lite parser accept a new construct needs the oracle battery and a fuzz run to agree, with the evidence in the work item's Outcome. A unit test alone is Blocking.
4. **Writes are lock, reload, compare, atomic write.** A write path takes `fsx.LockRepo`, reloads, checks the etag, and writes through `fsx.WriteAtomic`. Flag a repo-path write by `os.WriteFile`, `os.Rename` or `os.Create` outside `internal/fsx`, except in `internal/migrate` and `internal/initrepo`.
5. **A multi-step write has a defined failure state.** A new write that can stop half-way needs a rollback or an error that says which files changed. Do not flag the existing restore in `edit.commit`: work item 261029-xtqj owns it.
6. **Spec, `check` and the write path agree.** A new or changed rule must be enforced by both `check` and the operation that could break it. Missing one side is Blocking.
7. **Standard library only in the binary.** `go.yaml.in/yaml/v3` stays test-only; a new non-test import outside the standard library is Blocking.
8. **The server stays on loopback.** Blocking: a non-loopback bind, a weakened Host, Origin or token check, a request-supplied id or path that can resolve outside the repo root, or a secret or private path in a response, log or fixture.

## Records

A change a user or a client can see (a command, the API, the spec, a rule, a fixture, the bytes a write produces) needs, in the same change:

- an entry under `## Unreleased` in `CHANGELOG.md`;
- a `docs/glossary.md` entry for each new feature, command, field, rule or workbench part;
- an ADR with status `proposed` for a hard-to-reverse choice or one that rejects a real alternative.

A missing one is Blocking. A fix with no visible change needs none. Either fixture tree edited alone (`testdata/api` or `docs/spec/akm-api/.fixtures`) is Blocking until work item 261029-x4ws makes a test enforce it.

Never edit the body of an `accepted` ADR; supersede it. The one exception is a scrub of private material (names, paths, identifiers), which must say so in the commit message.

## Tests

A fix has a test that fails without it. A regression test for a parser or scanner change follows invariant 3. A gap here is a Note.

## Out of scope

Commit trailers, claims and the one-item-per-branch rule (`akmctl check` and `akm-commit` cover them), and workbench styling with no data, security or accessibility effect.

## Output

Start with `N actionable findings` (or `No actionable findings.`). Group under **Blocking** then **Note**. Each finding: `file:line`, the failure scenario, why the diff causes it, and a fix. Do not invent findings.
