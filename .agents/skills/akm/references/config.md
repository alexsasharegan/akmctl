# config

Added in this release. Use it to read or change the repository settings kept in `docs/.layout`.

```
akmctl config
akmctl config git-commits on|off
```

With no argument it prints the settings. `git-commits on` makes the optional `akm-commit` skill commit after each logical step; `off` stops that.

Files changed: `docs/.layout` only, when you set a value.

Errors: `invalid_request` for an unknown setting or value. If your `akmctl` has no `config` command, change nothing and tell the maintainer.
