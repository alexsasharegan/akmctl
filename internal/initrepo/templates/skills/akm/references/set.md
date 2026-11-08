# set

Use it to change frontmatter fields: status, priority, summary, `blocked_by`, `cycle` (an opaque label such as `2026-PI3-S4`: letters, digits, `.`, `_`, `-`, at most 64 characters; `cycle=null` removes it) and `keep` (`keep=true` marks an item never to be folded or archived; `keep=null` removes it). Not the body (use `section`), not `done` or `dropped` (use `close`), not `claimed_by` (use `claim`), not `closed` (`close` sets it).

```
$ akmctl set 261029-atxk status=ready priority=1 summary="Add a --json flag"
261029-atxk	2026-10-29	ready	add json output
changed: [priority status summary]
```

A list field (`blocked_by`, `decisions`, `scope`) takes entries separated by commas: `akmctl set <id> blocked_by=261029-aaaa,261029-bbbb`. A JSON array (`'blocked_by=["a","b"]'`) or a bracketed list (`blocked_by=[a,b]`) is read as the same list. `blocked_by=null` empties it.

A closed item is read-only: `set` refuses it (`reason: "closed"`) unless you pass `--correction "<why>"` (one line), and then only `title` and `summary` change, with `corrected <fields>: <why>` added to the Log in the same write. `keep` needs no reason. Use it for a typo or a wrong summary, not to change what the item was (kind, priority, parent and the rest are refused with `closed_field`).

Files changed: the one record; `updated` is set to today.

Errors:

- `invalid_request`: `status: must be one of raw, ready, active, blocked, done, dropped`; a managed field such as `closed`; a `cycle` that is not a label.
- `refused`, `reason: "rule"`: `WK5` blocked without `blocked_by`, `WK4` own parent, `WK9` closed item reopened, `A3` superseded without `superseded_by`.
- `stale`: the file changed since you read it; read it again.
