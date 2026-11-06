# claim

Use it to take a `ready` item before working it (claims extension). Sets `status: active`, `claimed_by` and `claimed_at`, and appends a Log line.

```
$ akmctl claim 261029-atxk --as demo-agent
261029-atxk	2026-10-29	active	add json output
$ akmctl claim 261029-atxk --as other
akmctl: refused: 261029-atxk is claimed by demo-agent
```

Claim with a name unique to the session, such as `claude-sonnet-5-5`: two sessions that share a name look like one holder.

`claimed_at` is UTC (`YYYY-MM-DDTHH:MM`, no `Z`). The dates beside it (`updated`) are local days and the Log line carries local time with its offset, so a `claimed_at` that looks hours ahead or behind is not a clock error. Do not edit it.

An item with an empty `## Next` is refused, because it would become `active` and fail rule WK6: write the `Next` first with `akmctl section <id> Next --text "..."`, then claim.

If someone else holds a fresh claim, take another item. `--takeover` is for a stale claim, with the maintainer's say-so.

Returns a `commit` (message and paths to stage). Use `commit.message` as the commit message exactly as given, with its `AKM-Step` and `AKM-Item` trailers, and stage `commit.paths` with the code that belongs to the step. Pass the whole message to one `git commit -F -`, not one `-m` per line, or git will not read the trailers. Do not write your own subject. One commit per step.

Files changed: the one record.

Errors: `refused` with `reason: "claimed"`, `"extension_off"` or `"rule"` (WK6); `invalid_request` without `--as`; `not_found`; `stale`.
