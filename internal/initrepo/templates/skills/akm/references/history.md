# history

Use it to see what happened to the work, in order and by whom: the AKM step commits of a range as moves. Read-only; needs a git repository and `git_commits` on.

```
$ akmctl history --range origin/main..HEAD
2026-11-06T12:08:22-07:00	261101-w0mz	claim	ready to active	Alex Regan + Claude Code
2026-11-06T12:22:15-07:00	261101-w0mz	block	active to blocked	Alex Regan + Claude Code
```

Moves come oldest first. Each has the commit, its time and author, the item, the step, the item's status before and after (`from`, `to`) and the other trailers of the message, such as `Assisted-by`. Only a commit whose subject is `<id>(<step>): <detail>` is a move; `commits` is the one that judges it. `--item <id>` shows one item's path, `--limit` keeps the newest moves (default 100) and says when it left older ones out.

To read only what is new, pass the `head` of the last answer: `--range <head>..HEAD` (use `--format json`).

Files changed: none.

Without a git work tree, or with `git_commits` off, it answers `available: false` with a `reason` (`not_git`, `git_commits_off`) and exit code 0: it is not an error. `invalid_request` for a bad `range`, `item` or `limit`.
