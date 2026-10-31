# commits

Use it to check that step commits say what they did. Read-only; needs a git repository. Exit code 0 is clean, 1 means a problem.

```
$ akmctl commits --range origin/main..HEAD
3f2a91c4	ok	261030-q03z(claim): claimed by claude
9c04d1ee	FAIL	261030-q03z(close): done
	[CM3] close: the commit does not move the file of 261030-q03z into docs/work/done/
2 AKM step commits checked, 5 other commits skipped
```

It examines commits whose subject is `<id>(<step>): <detail>` or that carry an `AKM-` trailer, and skips the rest. Rules CM1 (subject), CM2 (the two trailers agree with the subject) and CM3 (the diff does what the step says) are in the API specification. `--limit` caps the count (default 50).

Files changed: none.

Errors: `refused` with `reason: "not_git"`; `invalid_request` for a bad `range`.
