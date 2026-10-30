# migrate

Use it to see how far a repository is from the current layout, or to bring it there. `--mode report` (default) only reads; `plan` lists the steps; `apply-safe` runs the safe ones; `handoff` writes the rest for a person.

```
$ akmctl migrate --mode report
layout <version> -> <version>: current
```

`apply-safe` needs `--acknowledge-backup`, a git work tree and a clean tree; it creates a backup tag and a branch. Ask the maintainer before using it.

Files changed: none for report and plan; files and git refs for `apply-safe`.

Errors: `refused` with `reason` `no_backup_ack`, `not_git`, `dirty_tree`, `tag_exists`, `branch_exists` or `no_path`; `invalid_request` for a bad mode or target.
