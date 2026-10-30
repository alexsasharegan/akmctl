# init

Use it once, to make a new AKM repository. `--profile code` keeps records under `docs/`, `--profile notes` under `notes/` (wikilinks). For an existing repository use `migrate`.

```
$ akmctl init --profile code demo
initialised demo (profile code)
  wrote    .agents/skills/akm/SKILL.md
  wrote    AGENTS.md
  wrote    docs/.layout
  ...
```

It never overwrites: an existing file is listed as `skipped`, so a second run writes nothing. The result passes `check`.

Files changed: creates the layout, `AGENTS.md`, `CLAUDE.md`, skills, manifest and `.gitignore`.

Errors: `invalid_request` for an unknown profile or a `dir` with `..`; `refused` with `reason: "not_a_directory"`.
