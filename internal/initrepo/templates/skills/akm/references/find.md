# find

Use it to turn an id into a file path, so you can open it or stage it by path. Read-only. `--name` does the same for a note name, the target of a wikilink such as `[[Vault Views]]`.

```
$ akmctl find 261029-atxk
docs/work/261029-atxk-add-json-output.md
```

```
$ akmctl find --name "Vault Views"
Projects/Vault Views.md
```

Files changed: none.

Errors: `invalid_request` for a malformed id, `not_found` when no file has it, `ambiguous` when several do (`details.paths`). Give an id or a name, not both. A name matches a Markdown file by its name without `.md` (case ignored, folders optional) anywhere under the repository root except hidden folders.
