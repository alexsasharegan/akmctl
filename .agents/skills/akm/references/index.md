# index

Use it after adding or changing an ADR: it regenerates `docs/adr/README.md`, the routing index. Never edit that file by hand.

```
$ akmctl index
docs/adr/README.md: changed=true written=true
```

`--dry-run` shows the diff and writes nothing.

Files changed: the ADR index only.

Errors: `unreadable` when an ADR cannot be parsed; `refused` when the manifest is not at its usual path.
