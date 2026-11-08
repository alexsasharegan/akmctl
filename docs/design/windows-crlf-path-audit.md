---
type: design
status: final
created_at: 2026-11-08
updated_at: 2026-11-08
authored_by: agent
---

# Windows CRLF and path audit

Result of work item 261108-kbve (part of 261108-2ncf). The audit changed no behavior. Line endings were run with the built binary on a repository whose Markdown files were all converted to CRLF (what `core.autocrlf=true` produces on Windows), and with the migrate steps in a scratch test. Path forms ran on the Windows runner through `scripts/windows-smoke.ps1`.

## Line endings

| Operation                                                                                                                 | Existing CRLF record    | Result                                           |
| ------------------------------------------------------------------------------------------------------------------------- | ----------------------- | ------------------------------------------------ |
| `set` (fields, title)                                                                                                     | keeps CRLF              | pass                                             |
| `claim`, `release`                                                                                                        | keeps CRLF              | pass                                             |
| `log`                                                                                                                     | keeps CRLF              | pass                                             |
| `tick`, `tick --undo`                                                                                                     | keeps CRLF              | pass                                             |
| `section` (multi-line text, a new section)                                                                                | keeps CRLF              | pass                                             |
| `remove`                                                                                                                  | keeps CRLF              | pass                                             |
| `close` (the moved item, an inbound link in a record, an inbound link in the roadmap)                                     | keeps CRLF              | pass                                             |
| `split` (the parent)                                                                                                      | keeps CRLF              | pass                                             |
| `new work`, `new adr`, `split` (the new files)                                                                            | new file is LF          | fault 261108-2z9m (rule to decide)               |
| `index` (`docs/adr/README.md`)                                                                                            | rewrites the file as LF | fault 261108-ma45 (same cause as the `check` A4) |
| `check` on a CRLF `docs/adr/README.md`                                                                                    | A4 error                | fault 261108-ma45                                |
| `migrate` steps `rename_field`, `copy_field`, `set_default`, `title_from_h1`, `remove_field`, `goal_scaffold`, `mint_ids` | keeps CRLF              | pass                                             |

No operation left a file with mixed endings.

## Paths

| Case                                                                                  | Where run    | Result                                                         |
| ------------------------------------------------------------------------------------- | ------------ | -------------------------------------------------------------- |
| `--repo` with backslashes, forward slashes, mixed, a trailing separator, dot segments | Windows (CI) | pass                                                           |
| `--repo` doubled separators, and relative with a backslash (`.\repo`)                 | Windows (CI) | pass                                                           |
| A repository path with spaces                                                         | Windows (CI) | pass                                                           |
| `--repo` pointing into a subfolder (`docs`)                                           | macOS        | refused with a message that names the real root                |
| `find <id>` prints a path                                                             | Windows (CI) | forward slashes                                                |
| `find --name` with backslashes                                                        | Windows (CI) | refused: `name is not a note name` (a clear message, no crash) |
| Link repair by `close`                                                                | Windows (CI) | rewritten links use `/`                                        |
| A link target written with backslashes (`[x](..\adr\README.md)`)                      | macOS        | not checked and not rewritten: fault 261108-sfj3               |
| Wikilinks inside records with backslashes                                             | not run      | `find --name` (the resolver) was run, see above                |

## Not covered

- Drive-letter paths (`C:\...`) other than the runner's temp directory, and UNC paths (`\\server\share`).
- A console: prompts, colors and `ui` (work item 261108-ct4s).
- Reader and rename behavior: 261108-mwtr.
