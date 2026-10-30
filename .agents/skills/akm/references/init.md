# init

Use it once, to make a new AKM repository. `--docs-dir <name>` sets the one folder that holds the records and the manifest (default `docs`); `--wikilinks` makes records link with `[[wikilinks]]`. For an existing repository use `migrate`.

```
$ akmctl init --docs-dir knowledge demo
initialised demo (records in knowledge/)
  wrote    .agents/skills/akm/SKILL.md
  wrote    AGENTS.md
  wrote    knowledge/.layout
  ...
```

It never overwrites: an existing file is listed as `skipped`, so a second run writes nothing. The result passes `check`.

Files changed: creates the layout, `AGENTS.md`, `CLAUDE.md`, skills, manifest and `.gitignore` (only inside a git work tree).

Errors: `invalid_request` for a `docs_dir` that is not one plain folder name or a `dir` with `..`; `refused` with `reason: "not_a_directory"` or `"manifest_exists"` (another manifest is already there).
