# bases

Use it when the maintainer wants Obsidian views of the work: it writes an Obsidian Bases file that mirrors the workbench (table, cards and Kanban views of the work items, and a table of the ADRs).

```
$ akmctl bases
docs/akm.base: written=true
```

It writes once and never overwrites: a file that exists is `refused` with `reason: "exists"`. `--path <file.base>` picks another repository-relative file; `--dry-run --include-content` prints the text and writes nothing. The folder filters start at the nearest `.obsidian` folder above the repository, or at the repository root when there is none (the command says so).

Files changed: the one `.base` file.

Errors: `invalid_request` for a path that is not repository-relative or does not end in `.base`; `refused` when the file exists or the manifest is not at its usual path. The tool cannot open Obsidian: ask the maintainer to check that the views render.
