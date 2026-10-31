---
type: design
status: active
created_at: 2026-10-27
updated_at: 2026-10-28
---

# AKM specification

Version 0.3.0 (beta). The layout and this specification follow the beta's version number while it is in beta. Layout 0.3.0 adds the optional work item fields `closed`, `cycle` and `keep` (section 4, "Frontmatter") and the rules WK13 to WK17; an item that has none of them stays valid.

AKM (agent knowledge management) is a standard for a repository that keeps its own memory: its decisions, its direction and its work in flight, as Markdown records with frontmatter, so that any person or agent can pick up where work stopped. This document specifies the files, the work-item schema, the ADR and roadmap formats, and the rules of the checks. It is written for any maintainer or team, solo or not. The tool that implements it is `akmctl`; its operations and JSON are in the API spec ([akm-api-spec](akm-api-spec.md)), and a guide for adopting the standard is in [akm-guide](akm-guide.md).

Every rule here traces to a requirement in section 0 or to a clause of the contract in section 9, so a maintainer who wants to change one can see what it was for.

The words **MUST**, **SHOULD** and **MAY** are used as in RFC 2119. A "maintainer" is a person who may accept decisions and reorder the roadmap; a solo repository has one.

## 0. Requirements and identifiers

### Why these requirements

Repository memory kept in files fails in a few predictable ways: two people claim the same number; a hand-kept list drifts or conflicts on merge; a cold reader pays to read whole documents when a few fields would route them; a name edit breaks links. These requirements come from those failures. Rules elsewhere in this specification cite them as R1 to R10 (identifiers) and H1 to H7 (principles).

### Identifier requirements

| #   | Requirement                                                                                                                                                                                                                                    | Why                                                                                                                                                                                                                                           |
| --- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| R1  | One id space for ADRs and work items. An id names exactly one record in the repo; the type comes from the directory.                                                                                                                           | A reference like `parent` or `decisions` should not need a type prefix to resolve, and an ADR and a work item must never share an id.                                                                                                         |
| R2  | Low collision risk when people or agents create records in parallel branches. No shared counter.                                                                                                                                               | A monotonic counter makes two branches claim the same next number, and both append to the same index lines. The merge conflicts cost time and add noise to merges that were otherwise clean. Agents raise how many records start in parallel. |
| R3  | Fixed width and short (about 12 characters or fewer).                                                                                                                                                                                          | Stable, aligned references in tables, logs and prose.                                                                                                                                                                                         |
| R4  | The creation date is visible in the id, and ids sort chronologically.                                                                                                                                                                          | "When" at a glance, and plain `ls` order is time order.                                                                                                                                                                                       |
| R5  | "What" comes from the title and a slug, which are not part of the id.                                                                                                                                                                          | A readable slug makes an id long; keeping it out of the id keeps the id short.                                                                                                                                                                |
| R6  | The id alone is the reference. The word "ADR" is prose, not part of the id.                                                                                                                                                                    | One form to type, parse and search.                                                                                                                                                                                                           |
| R7  | Ids are minted by the tool or by the shell recipe in section 2, never typed by an agent or person.                                                                                                                                             | Agents are poor random-number generators and two will pick the same letters.                                                                                                                                                                  |
| R8  | An id can be read aloud and typed without confusion: lowercase, and no character that looks like another (`i` and `l` with `1`, `o` with `0`). `u` is also left out, following Crockford's Base32, to avoid ids that accidentally spell words. | Humans and speech transcription both use ids.                                                                                                                                                                                                 |
| R9  | An id is recognisable by pattern.                                                                                                                                                                                                              | `blocked_by` mixes ids with free text, so the check must tell them apart.                                                                                                                                                                     |
| R10 | Renaming a slug never breaks a reference.                                                                                                                                                                                                      | Titles and slugs get reworded; links must not depend on them.                                                                                                                                                                                 |

### Design principles

- **H1 Header layer.** Frontmatter carries every field a reader needs to route (find, filter, decide whether to open the body). A reader can read headers and skip bodies, so the cost of scanning many records stays small, as with HTTP headers and bodies.
- **H2 Generated, not hand-kept.** Anything derivable from headers is generated: the ADR index and the open-work listing. Hand-kept lists drift and conflict on merge.
- **H3 Tool-minted ids.** (R7)
- **H4 One source of truth per fact.** For example, children are found by their `parent` field, not listed in the parent.
- **H5 Frontmatter is YAML-lite, a simple subset of YAML 1.2 core.** Flat `key: value` lines, where a value is a scalar or a list of scalars (inline `[a, b]` or `- item` lines). No nesting, no anchors. An inline list stays on one line. It means the same as YAML 1.2 core, the schema Obsidian reads and writes, wherever it is accepted, and it is an error wherever YAML is richer (grammar below). A tool with no YAML library parses it with a small hand-written parser that edits by splicing bytes.
- **H6 The tool is an accelerator, not a dependency.** Every read and every record creation is possible without it, through plain files, the fallbacks and the id recipe in section 2, so a cold agent without the tool still works (C12). Regenerating the index and running the full checks need the tool.
- **H7 Obsidian-compatible.** Records are Markdown notes with frontmatter that Obsidian reads and edits without making them invalid. Where a design choice and Obsidian's behaviour conflict, Obsidian wins. The standard does not require Obsidian: a repository need not be a vault. Obsidian assigns property types by name for a whole vault, and the standard ships no type settings.

### Frontmatter format: YAML-lite (H5)

A record, and the manifest `docs/.layout`, start at byte 0 with a frontmatter block (the manifest has no delimiter lines; see section 7). This section is the specification of the hand-written parser.

**Relation to YAML.** YAML-lite is a strict subset of YAML 1.2 **core schema**, the schema Obsidian reads and writes. Every block the tool accepts is valid YAML and means the same thing to a YAML 1.2 core parser. Where YAML allows something YAML-lite does not support, the tool reports an error with a line number. It never reinterprets silently. So any YAML 1.2 library can read what the tool writes, and Obsidian can edit the same files without making them invalid.

_Supported, and true to YAML 1.2 core:_ flat `key: value` entries; double-quoted, single-quoted and plain strings; canonical decimal integers; `true` and `false`; `null`, `~` and the empty value (all null, counted as absent); lists of scalars, inline or block; comments.

_Plain strings that YAML 1.1 reads differently but 1.2 core reads as strings, and so are accepted as strings:_ dates and timestamps (`2026-10-28`, `2026-10-28T14:30`), `yes`, `no`, `on`, `off`, `y`, `n`, and base-60 forms such as `12:30`. This is the form Obsidian writes. Dates are therefore written unquoted, and a quoted date reads as the same string.

_Not supported, an error:_

- nested mappings, and block lists of mappings or of lists; flow mappings and nested flow lists; complex keys (`? `)
- multi-line scalars of any style (plain, quoted, literal `|`, folded `>`)
- anchors, aliases, tags and explicit types (`!!str`), directives, more than one document
- every number except a canonical decimal integer (signs other than `-`, leading zeros, underscores, floats, hex, octal, exponents, `.inf`, `.nan`): write them quoted if you mean text
- the capitalised spellings `True`, `TRUE`, `False`, `FALSE`, `Null`, `NULL` (YAML 1.2 core reads them as booleans and null, and the subset keeps one spelling): write `true`, `false`, `null`, or quote them
- empty list items, duplicate keys, and keys spelled `true`, `false`, `null` or `~` in any case (these are not string keys in YAML 1.2 core)

_Stricter than YAML (valid YAML that is rejected anyway):_ a tab anywhere in the block, a byte-order mark, any backslash escape other than `\"` and `\\`, and any value that does not fit on one line; any control character or non-character code point (YAML's printable set only); a key longer than 255 characters; in an inline list, a plain item that contains `?` or starts with `:`, or that is only `-`, `?` or `:`; and every spelling of a number that some YAML library reads as one: `0x`, `0o` and `0b` prefixes in any case, and a string that becomes such a number once underscores are removed. Libraries disagree on these corners, so the subset keeps clear of them.

_Other YAML parsers._ YAML 1.1 parsers, and some libraries when they decode into a generic value, read an unquoted date as a timestamp and `yes`, `no`, `on`, `off` as booleans. The tool does not. A test oracle must use YAML 1.2 core meanings and read date-shaped values as text (for example `go.yaml.in/yaml/v3`, comparing the node text of a timestamp-tagged scalar).

_Conformance test._ For every file an implementation accepts, a YAML 1.2 core parser (the oracle) returns the same data, with date-shaped values compared as text; for every reject fixture, the implementation errors. Every implementation of the standard runs the same fixtures.

_What Obsidian does to a note it edits:_ it rewrites the whole block from a parsed object. It keeps dates unquoted, drops comments, turns inline lists into block lists and removes unneeded quotes. All of that stays inside YAML-lite. The tool's own edits preserve everything else byte for byte; Obsidian's do not, which is cosmetic.

**Structure**

- The first line is exactly `---`; the block ends at the next line that is exactly `---` (a Markdown convention: inside the delimiters is one YAML document). No byte-order mark. Text is UTF-8. Lines end with LF or CRLF; a tool keeps whichever the file uses.
- Each line in the block is blank, a comment (`#` at the start of the line, after optional spaces), a `key: value` entry, or a list item belonging to the entry above.
- A key matches `[A-Za-z_][A-Za-z0-9_.-]*`, then `:`, then a space and a value, or the end of the line. A duplicate key is an error. A key spelled `true`, `false`, `null` or `~` in any case is an error.
- Indentation uses spaces only; a tab is an error. Only list items may be indented.

**Values**

- A scalar is a double-quoted string, a single-quoted string or a plain scalar.
- Double-quoted: the only escapes are `\"` and `\\`; any other backslash is an error. Single-quoted: `''` is one quote.
- Plain: non-empty text that does not start with `,` `[` `]` `{` `}` `#` `&` `*` `!` `|` `>` `%` `@` or a backtick; does not start with `-`, `?` or `:` followed by a space or the end of the line; does not contain `: ` and does not end with `:`; and does not contain ` #` (a space then `#` starts a comment). Inside an inline list a plain scalar also cannot contain `,` `[` `]` `{` `}`.
- Typing: `-?(0|[1-9][0-9]*)` is an integer; `true` and `false` are booleans; `null`, `~` and the empty value are null; any other unquoted text is a string, unless it is one of the number forms or capitalised spellings listed as errors above. A quoted value is always a string. Version strings such as `1.0.0` and ids such as `261027-m3qz` are ordinary strings; `1.2` is a float and needs quotes.
- Fields of type date hold `YYYY-MM-DD` as a plain string, the local calendar day of the machine that writes it, and timestamps hold `YYYY-MM-DDTHH:MM` in UTC. The tool writes them unquoted, as Obsidian does.
- A list is an inline list on one line, `[` scalars separated by commas (no empty items, an optional trailing comma) `]`; or block form: the key with no value, followed by lines `- scalar`, all indented by the same number of spaces (zero or more).
- An entry with no value and no list items is empty (null) and counts as absent.

**Editing.** To change a value, rename a key, add a key (at the end of the block) or remove one, the tool changes only the bytes of that entry. Everything else in the block, including order, comments, quoting and line endings, is byte-for-byte unchanged. Unknown keys are kept. Errors carry the file and the line number.

**Fixtures.** The machine-readable set is `yaml-lite.json` (52 cases: input, accept or reject, and the expected value), shipped with the API fixtures (API spec, section 9). It is the conformance suite for every implementation. The table below is a human-readable subset:

| Input                                              | Result                                                                                       |
| -------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `status: active`                                   | accept: string                                                                               |
| `priority: 2`                                      | accept: integer                                                                              |
| `priority: 02`                                     | reject: leading zero                                                                         |
| `priority: +2`                                     | reject: sign other than `-`                                                                  |
| `created: 2026-10-28`                              | accept: string (a date field)                                                                |
| `created: "2026-10-28"`                            | accept: the same string                                                                      |
| `claimed_at: 2026-10-28T14:30`                     | accept: string (a timestamp field)                                                           |
| `layout: 0.3.0`                                    | accept: string                                                                               |
| `layout: 1.2`                                      | reject: float (quote it)                                                                     |
| `time: 12:30`                                      | accept: string (YAML 1.1 reads a base-60 integer)                                            |
| `done: yes`                                        | accept: string (YAML 1.1 reads a boolean)                                                    |
| `flag: true`                                       | accept: boolean                                                                              |
| `flag: True`                                       | reject: capitalised boolean                                                                  |
| `gone: ~`                                          | accept: null                                                                                 |
| `gone: Null`                                       | reject: capitalised null                                                                     |
| `id: 261027-m3qz`                                  | accept: string                                                                               |
| `summary: "Fix the \"x\" case"`                    | accept                                                                                       |
| `summary: 'it''s fine'`                            | accept                                                                                       |
| `title: Add: a thing`                              | reject: `: ` in a plain scalar (quote it)                                                    |
| `parent: "[[261027-m3qz-fix-login\|261027-m3qz]]"` | accept: string (a wikilink reference)                                                        |
| `parent: [[261027-m3qz-fix-login]]`                | reject: nested list (quote it)                                                               |
| `blocked_by: [261015-k7f2, "a person"]`            | accept: inline list                                                                          |
| `blocked_by: [a, [b]]`                             | reject: nested list                                                                          |
| `blocked_by: [a,,b]`                               | reject: empty item                                                                           |
| `tags:` then `- a` and `- b`                       | accept: block list                                                                           |
| `tags:` then `- a` and `   - b`                    | reject: items not at the same indentation (YAML itself reads this as the one string `a - b`) |
| `owner: x` twice                                   | reject: duplicate key                                                                        |
| `  nested: 1` outside a list                       | reject: indented line                                                                        |
| `note: &x value`                                   | reject: anchor                                                                               |
| `text: \|`                                         | reject: block scalar                                                                         |
| `key: value # why`                                 | accept: the comment is not part of the value                                                 |
| `url: http://example.com`                          | accept: string (`:` not followed by a space)                                                 |
| `x: ]a`                                            | reject: starts with a flow indicator                                                         |
| `on: x`                                            | accept: the key `on` is a string key in YAML 1.2 core                                        |
| `true: x`                                          | reject: the key is a boolean                                                                 |

### Identifier format

`YYMMDD-xxxx`, for example `261027-k7f2`. Width 11.

- `YYMMDD` is the creation date, two-digit year (20YY), as the **local calendar day** of the machine that mints the id (or the `as_of` day when the request gives one), so an evening in a western time zone does not carry tomorrow's date.
- `xxxx` is four random characters from the alphabet `0123456789abcdefghjkmnpqrstvwxyz` (Crockford's Base32 alphabet, lowercase; 32 characters; `i`, `l`, `o` omitted as lookalikes for `1` and `0`, `u` omitted to avoid accidental words; no check digit, and no decoding of lookalike or uppercase characters).
- Pattern: `^\d{6}-[0-9a-hjkmnp-tv-z]{4}$`.
- Collision odds: 1,048,576 suffixes per day. With 20 records created across parallel branches in one day, the chance of any collision is about 1 in 5,000. Two files with the same id and different slugs merge cleanly in git, so a collision is caught when `akmctl check` runs (rule I1). Repair it by hand: rename the later file, edit its `id`, and search for references to the old id.

### Filenames and references

- File: `<id>-<slug>.md`, in `docs/work/`, `docs/work/done/` or `docs/adr/`. The `slug` is lowercase kebab-case and SHOULD stay under 40 characters. Only the filename carries the slug.
- A reference is the id alone: in frontmatter (`parent`, `blocked_by`, `decisions`, `supersedes`, `superseded_by`), in roadmap lines, and in prose or commit messages. In prose an ADR is introduced as "ADR 261015-k7f2".
- In frontmatter a reference MAY instead be a wikilink in a quoted string, so Obsidian shows it as a link: `parent: "[[261027-m3qz-fix-login|261027-m3qz]]"`. The target is the record's filename without `.md`, because Obsidian resolves links by whole filename, not by id prefix; the alias after `|` keeps the id visible. The id is the leading id of the target, and that is all the tool reads: the slug is not checked, so a renamed slug is not an error (Obsidian rewrites the link itself). The tool reads both forms everywhere a frontmatter reference is allowed (WK4, WK5, WK11), writes the bare id unless the manifest sets `reference_style: wikilink`, and never rewrites a reference it did not create. Roadmap lines, prose and commit messages use the bare id only.
- The `id` frontmatter field equals the filename's id part. To open a record from its id, find the file in the three directories whose name starts with `<id>-`. The tool offers `akmctl find <id>`.
- Renaming the slug (renaming the file) does not change the id.

### Minting and utilities

| Command                  | Does                                                                    |
| ------------------------ | ----------------------------------------------------------------------- |
| `akmctl new work <slug>` | Mint an id, write the file from the work-item template, print the path. |
| `akmctl new adr <slug>`  | Same for an ADR, status `proposed`.                                     |
| `akmctl find <id>`       | Print the path of the record.                                           |
| `akmctl when <id>`       | Print the id's date as `YYYY-MM-DD`, the weekday, and the age in days.  |

There is no command to renumber a record. Repair a collision by hand, as above.

Listings (`akmctl ls`, the ADR index) show the expanded date next to the id, so a reader never has to decode `261027` by eye. The id's date is the local calendar day of the machine that minted it; it is not a timestamp.

## 1. Layout

```text
AGENTS.md                  60 to 80 lines; identity, commands, boundaries, pointers
CLAUDE.md                  one line: @AGENTS.md   (a tool-specific import; optional)
README.md                  for humans
.prettierrc.json           {"embeddedLanguageFormatting": "off"} (section 6)
.prettierignore            lists docs/adr/README.md (section 4)
docs/
  ROADMAP.md               hand-ordered lists (section 5)
  architecture.md          the map: layout, data flow, invariants (optional)
  adr/                     decision records (section 4)
    README.md              GENERATED routing index; do not edit
    261015-k7f2-use-sqlite.md
  work/                    work items (section 3)
    261027-m3qz-add-json-output.md
    done/                  closed items
  .layout                  manifest: layout and extension versions (section 7)
  .akm-deny                the deny list of the safety extension (section 6)
  .akm.lock                advisory lock the tool takes while it writes; not committed (list it in .gitignore)
  .akm-migrate-log         one JSON line per migration step applied (section 7)
.agents/skills/
  <skill>/                 shared skill copy: SKILL.md, scripts/, assets/
```

- The minimal tier (section 8) is the smallest set of these files.
- **Moved locations.** The work folder, the ADR folder and the roadmap file sit at `docs/work/`, `docs/adr/` and `docs/ROADMAP.md` unless the manifest says otherwise: the keys `work_dir`, `adr_dir` and `roadmap` (section 6, "Settings") move them, for a repository that keeps its records inside another structure, such as a notes vault. `done/` is always `<work_dir>/done/`; the rules below name the default paths and apply to the configured ones.
- **The manifest folder.** The manifest `.layout`, the deny list `.akm-deny`, the lock `.akm.lock` and the migration log `.akm-migrate-log` (named `.migrate-log` before 2026-10-30; a tool that finds the old name renames it before it appends) live together in one folder, `docs/` unless the repository chose another name (`akmctl init --docs-dir`). A tool finds it from the repository root, which is the nearest ancestor of the working directory that has a manifest or an `AGENTS.md`: first `docs/.layout`; if that is absent, the one immediate subdirectory of the root, not hidden and not `node_modules`, that holds a `.layout` file. Two or more such subdirectories are an error that names them; none means the repository has no manifest. The folder name is a single path segment, so discovery looks one level deep. Where this specification says `docs/.layout`, `docs/.akm-deny`, `docs/.akm.lock` or `docs/.migrate-log`, it means that file in the manifest folder.
- Long policies (privacy, publishing, deployment) live in named docs under `docs/` and are linked from the Do-not list of `AGENTS.md`.

## 2. `AGENTS.md` contract

`AGENTS.md` MUST have these parts, in this order, and nothing else of length. Part 1 is a one-sentence paragraph under the `#` title. Parts 2 to 5 are `##` headings with exactly the titles in bold below (rule G1):

1. **Identity** (no heading): one sentence: what this is, who uses it.
2. **Commands**: build, test, lint, run, and the full check before a commit (including `akmctl check`, section 6).
3. **Where things live**: pointers to `docs/architecture.md`, `docs/adr/`, `docs/ROADMAP.md`, and the command that lists open work.
4. **Rules**: the working rules for agents (below).
5. **Do not, without asking**: boundaries. Each links a policy doc when one exists.

Required rules (in the Rules part):

- Before starting, list the open work (`akmctl ls --status ready,active,blocked`, or the fallbacks below); resume the `active` item you hold (Claiming, section 3); if you hold none, take the top `ready` item; if several `active` items are unclaimed, ask.
- Work only under a work item. If none fits, create a `raw` one and ask which item to work.
- A hard-to-reverse choice, or one that rejects a real alternative, needs an ADR draft with status `proposed`. Never edit the body of an `accepted` ADR; supersede it.
- Update the item's `## Next` and `## Log` in the same commit as the code.
- Run `akmctl check` before committing docs, and format with the repository's pinned Prettier.
- Code and tests win over prose. When a document is wrong, fix it or say so; do not work around it.
- A repo's `AGENTS.md` wins over the shared skill where they disagree.
- If the tool is not installed, use the plain-file fallbacks below and tell the maintainer how to install it.

Plain-file fallbacks (H6), which the `AGENTS.md` template carries:

| Need                                 | Without the tool                                                                                                                    |
| ------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------- |
| Which decisions exist                | Read `docs/adr/README.md` (committed, generated).                                                                                   |
| What work is active                  | `grep -lE '^status: ?"?active"?$' docs/work/*.md`                                                                                   |
| What work is ready                   | `grep -lE '^status: ?"?ready"?$' docs/work/*.md`                                                                                    |
| What is blocked, and on whom         | `grep -lE '^status: ?"?blocked"?$' docs/work/*.md`, then read `blocked_by` in each.                                                 |
| Where to resume an item              | Read its `## Next` section.                                                                                                         |
| Mint an id                           | The recipe below.                                                                                                                   |
| Create a record                      | Copy the template, name the file `<id>-<slug>.md`, set `id` to the same id.                                                         |
| Close an item                        | Set `status`, write `## Outcome`, set `updated` and `closed` to today's date, and move the file to `docs/work/done/` with `git mv`. |
| Claim an item (claims extension)     | Make the first `## Log` line `claimed (where)`.                                                                                     |
| Find a record by id                  | `ls docs/work docs/work/done docs/adr \| grep '^<id>-'`                                                                             |
| The ADR index is stale or conflicted | Do not edit it by hand; tell the maintainer to regenerate it with the tool.                                                         |

The id recipe, which `AGENTS.md` carries. It mints the same form as the tool; run it again if the second command prints a match:

```sh
id="$(date +%y%m%d)-$(LC_ALL=C tr -dc '0-9a-hjkmnp-tv-z' </dev/urandom | head -c 4)"
ls docs/work docs/work/done docs/adr | grep "^$id-"
```

A fallback is slower and uses more tokens; it must stay correct.

Must-hold boundaries SHOULD also be enforced by a tool-level control (permissions or a hook) where the agent tool supports one; text alone is advisory. This is outside the tool's checks.

## 3. Work items

A work item is the one unit of in-flight work: an issue, a task, a feature, a project, a spike. A feature or project is a work item that other items name as `parent`.

### File

`docs/work/<id>-<slug>.md`, with `<id>` as defined in section 0 (for example `261027-m3qz-add-json-output.md`). Ids are unique across work items and ADRs together and never reused. A closed item moves to `docs/work/done/` with the same name; its `id` and links do not change.

### Frontmatter

| Field         | Type        | Req          | Values and default                                                                                                                                                                                                                                                                                                           |
| ------------- | ----------- | ------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id`          | string      | yes          | The id part of the filename: `YYMMDD-xxxx`.                                                                                                                                                                                                                                                                                  |
| `title`       | string      | yes          | Imperative or noun phrase; matches the first heading.                                                                                                                                                                                                                                                                        |
| `summary`     | string      | no           | One sentence, up to 160 characters: what this is and why it matters. Shown in listings so the body need not be opened.                                                                                                                                                                                                       |
| `status`      | enum        | yes          | `raw`, `ready`, `active`, `blocked`, `done`, `dropped`                                                                                                                                                                                                                                                                       |
| `kind`        | enum        | no           | `feature`, `bug`, `chore`, `spike`. Default `chore`.                                                                                                                                                                                                                                                                         |
| `parent`      | string      | no           | Reference to another item (an `id` or a wikilink to it). Absent for a top-level item.                                                                                                                                                                                                                                        |
| `blocked_by`  | list        | if `blocked` | Entries are item ids (matching the id pattern) or free text naming a person, role or outside dependency.                                                                                                                                                                                                                     |
| `priority`    | integer     | no           | Higher is sooner. Absent counts as 0. Sorts the listing only; the roadmap supersedes it.                                                                                                                                                                                                                                     |
| `created`     | date        | yes          | The local date, `YYYY-MM-DD`, a plain string (YAML-lite).                                                                                                                                                                                                                                                                    |
| `updated`     | date        | yes          | The local date, `YYYY-MM-DD`; set whenever the file changes.                                                                                                                                                                                                                                                                 |
| `closed`      | date        | no           | The local date, `YYYY-MM-DD`, on which the item was closed. Set only by `close`, for `done` and `dropped` alike; `set` cannot change it. Present only on a closed item (WK13 to WK15).                                                                                                                                       |
| `cycle`       | string      | no           | The cycle (sprint, program increment, quarter) the item is planned into, for example `2026-PI3-S4`: 1 to 64 characters of `A-Z`, `a-z`, `0-9`, `.`, `_` and `-`, starting with a letter or digit (WK16). The tool treats it as an opaque label: it sets, shows, lists and filters by it, and does not order or interpret it. |
| `keep`        | boolean     | no           | `true` marks the item as never to be folded or archived by a later compaction; `false` and absence mean the same (WK17).                                                                                                                                                                                                     |
| `authored_by` | string      | no           | `agent` while a maintainer has not reviewed it; removed on review (C11).                                                                                                                                                                                                                                                     |
| `decisions`   | list of ids | no           | ADR ids that constrain or result from this item.                                                                                                                                                                                                                                                                             |

Unknown fields are allowed and ignored (a team MAY add its own). An item MUST NOT use a status outside the enum.

`closed`, `cycle` and `keep` are a MINOR addition, made in layout 0.3.0 (ADR 261029-sg82): an item without them is valid, so no item of a 0.2.0 repository needs a change (only its manifest versions, which the maintainer raises). `closed` is the day the item closed, in the local day like `updated` (`as_of` wins); an item closed before 0.3.0 has none until someone backfills it (rule WK13 warns). A calendar of cycles, a check of `cycle` against one, and compaction of closed work are not part of this version.

### Statuses

| Status    | Meaning                    | Needs                                                       |
| --------- | -------------------------- | ----------------------------------------------------------- |
| `raw`     | Captured, not specified.   | A title.                                                    |
| `ready`   | Specified; can be started. | `## Goal` with a "done when" list.                          |
| `active`  | Being worked.              | A non-empty `## Next`.                                      |
| `blocked` | Cannot proceed.            | `blocked_by` non-empty; `## Next` says what unblocks it.    |
| `done`    | Finished and verified.     | `## Outcome` with evidence for each "done when" line (C10). |
| `dropped` | Decided against.           | `## Outcome` with the reason.                               |

Legal transitions: `raw` to `ready` to `active`; `active` to `blocked` and back; `active` to `done`; any open status to `dropped`. A check can see only the current state, so transitions are guidance in the skill, not checked.

### Sections

In this order. Headings are exactly as written.

```markdown
# <Title>

## Goal
One paragraph, then a checklist of "done when" lines, each independently checkable.

## Next
The single next concrete step. Overwritten at checkpoints. Required when active or blocked.

## Log
- YYYY-MM-DD: dated lines, appended, never rewritten.

## Outcome
Filled at close. For `done`: for each "done when" line, the test, command or file that proves it.
For `dropped`: the reason.
```

- `## Goal` is required for `ready`, `active`, `blocked` and `done`. `## Next` is required for `active` and `blocked`. `## Outcome` is required for `done` and `dropped`. A `dropped` item needs only an Outcome, so a `raw` item can be dropped without a Goal.
- A parent item keeps its own `## Next`. Its `## Goal` MAY list its children by id.
- `## Next` MAY end with a `Saved: YYYY-MM-DD HH:MM` line.
- The Log is append-only by convention; the check does not diff it.

### Parent and child

- A child names its parent with `parent`. Children are found by search, not listed in the parent's frontmatter (one source of truth).
- A parent MUST NOT be `done` while any child is not `done` or `dropped`. A parent MAY be `dropped` only if every open child is dropped or re-parented first.
- Depth is unlimited by the schema; the skill recommends two levels.

### Claiming (extension `claims`, optional)

- Adds two frontmatter fields to a work item: `claimed_by` (a free string such as `claude-code@mac`, a person or a session name) and `claimed_at` (a UTC timestamp, `YYYY-MM-DDTHH:MM`).
- `akmctl claim <id> --as <name>` sets both when it moves the item to `active`. The tool clears them on any move out of `active` and with `akmctl release <id>`. Each claim is also logged.
- A claim is advisory and never a lock: it is visible to others only after it is committed (and pushed, if there is a remote). A claim older than `claim_stale_hours` (default 24) shows as stale, and any agent may take the item over by logging the takeover.
- If another holder has a claim that is not stale, do not work the item: tell the maintainer, and take another `ready` item.
- Rules (only when the extension is enabled): CL1 and CL2 in section 6.
- Without the tool, the first `## Log` line says `claimed (where)`.
- Gaps: agents have no stable identity (use the tool plus the branch or session); two branches can claim the same item and meet only at merge, where the field conflicts, which at least makes it visible.

## 4. Decision records (ADRs)

`docs/adr/<id>-<slug>.md`, with `<id>` as defined in section 0 (for example `261015-k7f2-use-sqlite.md`). ADRs and work items share one id space (R1).

### Frontmatter (the header layer)

Every field the index and a routing read need is here (H1):

| Field           | Req             | Values                                                                                                                                                             |
| --------------- | --------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `id`            | yes             | The id part of the filename.                                                                                                                                       |
| `title`         | yes             | Short statement of the decision.                                                                                                                                   |
| `summary`       | yes             | One sentence, up to 160 characters: what was decided and when to read this record. The index shows it.                                                             |
| `status`        | yes             | `proposed`, `accepted`, `rejected`, `deprecated`, `superseded`                                                                                                     |
| `date`          | yes             | The date of the last status change, `YYYY-MM-DD`.                                                                                                                  |
| `scope`         | no              | List of paths, globs or topic words this decision applies to (for example `[src/store/, cli]`). Lets a reader working in an area find the records that bear on it. |
| `supersedes`    | no              | ADR id.                                                                                                                                                            |
| `superseded_by` | if `superseded` | ADR id.                                                                                                                                                            |
| `authored_by`   | no              | `agent` until a maintainer reviews.                                                                                                                                |

Sections: `## Context`, `## Considered Options` (at least two, including doing nothing), `## Decision`, `## Consequences`, `## Confirmation` (optional: a test or grep that fails if the decision is violated), `## Review` and `## Revert` (optional).

- Acceptance is a maintainer's commit that sets `status: accepted`, or the merge of the PR that carries it. An agent drafts as `proposed` and does not accept.
- An `accepted` body MUST NOT be edited except status, `superseded_by`, an errata section and typos (convention only; no check). A change of mind adds a new ADR and sets the old one to `superseded`.

### The generated index

`docs/adr/README.md` is a routing table built from the frontmatter of every ADR, reading headers only and never bodies (H1, H2). It is committed, so a reader with no tooling can open it. It starts with `<!-- generated by akmctl index; do not edit -->`.

- Groups, in order: **In force** (`accepted`), **Proposed**, **Not in force** (`rejected`, `deprecated`, `superseded`).
- Columns: `id`, `title`, `date`, `scope`, `summary`, `supersedes / superseded_by`.
- Rows sort by `id` (so by date) within a group.
- The id column links to the file with a relative Markdown link, not a wikilink, so the index reads the same on a Git host and in Obsidian.
- Layout: each group is a `##` heading, a blank line and a pipe table; an empty group is omitted. The header row, a separator row `| --- |` and the data rows are padded so that each column is as wide as its widest cell, measured in display width (East Asian wide characters and emoji count 2). A `|` in a cell is written `\|`. The file is the header comment, a blank line, then the groups, and ends with one newline. This is the form Prettier leaves unchanged. Display width follows Prettier's string-width rule: ASCII counts its length; East Asian wide and fullwidth characters count 2; an emoji sequence counts 2 (1 for a bare text-style emoji); control characters, combining marks U+0300 to U+036F and variation selectors count 0.
- The file is also listed in `.prettierignore`, and rule A4 compares it byte for byte with the generated output.
- A merge conflict in the index is resolved by running `akmctl index`, never by hand. A repo MAY add a `.gitattributes` merge driver that does this.

The open-work listing is the work equivalent. It is not committed: work records change too often, and a committed copy would conflict constantly. It is derived on demand with `akmctl ls --status ready,active,blocked` (header fields only).

## 5. Roadmap

`docs/ROADMAP.md`: hand-ordered lists only. No prose blocks, no dated "where things stand".

```markdown
# Roadmap

## Now
- [261027-m3qz] Add JSON output
## Next
- [261028-b8wd] Config file support
## Later
- [261101-q4ts] Plugin API
## Dropped
- [261020-x9ha] Windows installer: out of scope; see ADR 261015-k7f2
```

- Sections: `Now`, `Next`, `Later`, `Dropped`. Order within a section is the priority.
- A line is `- [id] text`, where `id` is a work-item id. It MUST resolve to a work item. Text after the id is free; a `Dropped` line SHOULD give a reason.
- The roadmap lists top-level work, usually parents and standalone items. A child item need not appear.
- The roadmap supersedes `priority`; the derived listing shows roadmap order first, then `priority`, then `id`.

## 6. The tool: checks and queries

The tool is `akmctl`: one Go binary with subcommands, with a stable JSON output and an API (the API spec, [akm-api-spec](akm-api-spec.md)). It reads the repo's files and git history, writes nothing unless a subcommand says so, and needs no network. It uses the Go standard library only: its own display-width code and its own frontmatter parser, with no YAML library. It is run by the agent, and it is an accelerator, not a dependency (H6). The rules in "Rules" below are its checks.

### Subcommands

| Command                                         | Does                                                                                                                                                                                                                                   |
| ----------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `akmctl check`                                  | Run all checks; print findings; exit non-zero on any error. `--target <version>` classifies each finding for a layout version (section 7).                                                                                             |
| `akmctl ls`                                     | List records. Filter by `--type work\|adr`, `--status`, `--parent`, `--blocked`, `--scope`, `--since`; sort with `--sort`; choose columns with `--fields`; `--limit`; `--count-only`; `--format table\|ids\|json`. Header fields only. |
| `akmctl search <text>`                          | Match title, summary and scope; `--body` also searches bodies.                                                                                                                                                                         |
| `akmctl show <id>`                              | Print one record; `--section next` prints only that section.                                                                                                                                                                           |
| `akmctl board`                                  | A text Kanban by status; the JSON of `ls` is what a UI consumes.                                                                                                                                                                       |
| `akmctl index`                                  | Regenerate `docs/adr/README.md`.                                                                                                                                                                                                       |
| `akmctl migrate`                                | Reconcile a repo with a target layout version: `--mode report\|plan\|apply-safe\|handoff` (section 7).                                                                                                                                 |
| `akmctl set <id> field=value`                   | Change a field and refresh `updated` (not `closed`, which only `close` sets).                                                                                                                                                          |
| `akmctl close <id> done\|dropped`               | Set the status, move the file to `done/`, refresh `updated` and set `closed`; refuses an empty `## Outcome`.                                                                                                                           |
| `akmctl claim <id> --as <name>`, `release <id>` | Claims extension (section 3).                                                                                                                                                                                                          |
| `akmctl new work\|adr <slug>`, `find`, `when`   | The identifier commands of section 0.                                                                                                                                                                                                  |
| `akmctl ui`                                     | Start the workbench UI.                                                                                                                                                                                                                |

Exit codes: `0` clean or warnings only; `1` at least one error; `2` the command could not run or was refused (an unreadable repo, an unknown id, a refused write).

### Output contract

- `--format json` is the stable output: a versioned envelope with `api_version`, the request echoed, `records` or `findings`, and `count`. Field names are the frontmatter names; fields absent on a record type are `null`.
- The CLI is one binding of a transport-neutral operation set. An HTTP binding returns the same JSON. The operations, request and response schemas, the CLI and HTTP bindings and the conformance fixtures are in [akm-api-spec](akm-api-spec.md).

### Open-work listing (`akmctl ls --status ready,active,blocked`)

Items with status `ready`, `active` or `blocked`, grouped as: **Active**, **Blocked** (grouped by `blocked_by`, so what waits on whom is visible), **Ready**, in this sort order: roadmap order, then `priority` descending, then `id`. Each row: id, expanded date, title, status, `parent` if any, `updated`, and `summary` if present. Header fields only (H1); `ls --next` adds the first line of each `## Next`, which is in the body. `raw` items are counted, not listed. This listing replaces a hand-written status board. It is never committed.

### Rules

Severity: **E** error (exit 1), **W** warning (printed, exit 0). Rule ids: `WK` work items, `I` ids, `A` ADRs, `RM` roadmap, `G` `AGENTS.md`, `L` links, `S` skills and versions, `CL` claims, `SF` safety.

| ID   | Sev | Rule                                                                                                                                                                                                                                                                                                                             |
| ---- | --- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| WK1  | E   | Every work file parses; frontmatter present and within the simple subset (H5); `id` equals the filename's id part and matches the id pattern.                                                                                                                                                                                    |
| WK2  | W   | The date part of the id is a valid date and equals `created`.                                                                                                                                                                                                                                                                    |
| I1   | E   | Ids are unique across `docs/work/`, `docs/work/done/` and `docs/adr/` together (R1).                                                                                                                                                                                                                                             |
| WK3  | E   | `status` is in the enum; required fields (`title`, `created`, `updated`) present; dates are valid.                                                                                                                                                                                                                               |
| WK4  | E   | `parent`, if present, resolves to an existing item and is not the item itself; no parent cycles.                                                                                                                                                                                                                                 |
| WK5  | E   | `status: blocked` has a non-empty `blocked_by`; every `blocked_by` entry that matches the id pattern resolves to an item; other entries are free text and are not checked.                                                                                                                                                       |
| WK6  | E   | `status: active` or `blocked` has a non-empty `## Next`.                                                                                                                                                                                                                                                                         |
| WK7  | E   | `ready`, `active`, `blocked` or `done` has a `## Goal`; `done` or `dropped` has a non-empty `## Outcome`.                                                                                                                                                                                                                        |
| WK8  | E   | A parent with `status: done` has no child that is not `done` or `dropped`.                                                                                                                                                                                                                                                       |
| WK9  | E   | An item in `docs/work/done/` has status `done` or `dropped`; a `done` or `dropped` item is in `done/`.                                                                                                                                                                                                                           |
| WK10 | W   | `active` item with `updated` older than 30 days (configurable).                                                                                                                                                                                                                                                                  |
| WK11 | E   | Field types: `kind` is in its enum, `priority` is an integer, `blocked_by` and `decisions` are lists, every `decisions` entry resolves to an ADR, `summary` is at most 160 characters.                                                                                                                                           |
| WK12 | W   | `title` matches the first heading.                                                                                                                                                                                                                                                                                               |
| WK13 | W   | An item in `docs/work/done/` has no `closed`.                                                                                                                                                                                                                                                                                    |
| WK14 | E   | An item whose `status` is not `done` or `dropped` has `closed`.                                                                                                                                                                                                                                                                  |
| WK15 | E   | `closed` is a valid date (`YYYY-MM-DD`) and is not earlier than `created`.                                                                                                                                                                                                                                                       |
| WK16 | E   | `cycle` matches `[A-Za-z0-9][A-Za-z0-9._-]{0,63}`.                                                                                                                                                                                                                                                                               |
| WK17 | E   | `keep` is `true` or `false`.                                                                                                                                                                                                                                                                                                     |
| A1   | E   | ADR files parse; frontmatter within the simple subset (H5); `id` equals the filename's id part and matches the id pattern; `title`, `summary`, `status`, `date` present; `summary` is at most 160 characters.                                                                                                                    |
| A2   | E   | ADR `status` is in the enum; `date` valid.                                                                                                                                                                                                                                                                                       |
| A3   | E   | `superseded` has `superseded_by`; the target exists, is an ADR, and is not itself; if the target has `supersedes`, it points back. `supersedes` and `superseded_by` resolve to ADRs.                                                                                                                                             |
| A4   | E   | `docs/adr/README.md` equals, byte for byte, what `akmctl index` would generate from the current frontmatter.                                                                                                                                                                                                                     |
| RM1  | E   | Every `[id]` in `ROADMAP.md` resolves to a work item (not an ADR).                                                                                                                                                                                                                                                               |
| RM2  | E   | Roadmap sections are among `Now`, `Next`, `Later`, `Dropped`; an id appears once.                                                                                                                                                                                                                                                |
| RM3  | E   | An item listed under `Now` or `Next` or `Later` is not `done` or `dropped`; an item under `Dropped` is `dropped`.                                                                                                                                                                                                                |
| RM4  | W   | A `ready` or `active` top-level item (no `parent`) is not on the roadmap.                                                                                                                                                                                                                                                        |
| G1   | E   | `AGENTS.md` exists and has the headings of section 2.                                                                                                                                                                                                                                                                            |
| G2   | W/E | `AGENTS.md` length: warn over 80 lines, error over 150.                                                                                                                                                                                                                                                                          |
| G3   | W   | A command named in `AGENTS.md` Commands is not a script or target that exists (checked for `scripts/…` paths and `make`/`npm run`/similar targets it can resolve).                                                                                                                                                               |
| L1   | E   | Relative Markdown links (`[text](path)`) in `docs/` and `AGENTS.md` resolve to existing files. External URLs and same-file anchors are skipped, as are code spans and code blocks. A path with a space is written `%20`-encoded or in angle brackets.                                                                            |
| L2   | W   | A wikilink (`[[name]]`, `[[name\|alias]]`, `[[name#heading]]`) in `docs/` or `AGENTS.md` does not resolve: `name` must match exactly one file in the repository by filename without `.md` (or by path), as Obsidian resolves it. Zero matches and more than one match are both findings. Code spans and code blocks are skipped. |
| CL1  | W   | (claims enabled) An `active` item has no `claimed_by`.                                                                                                                                                                                                                                                                           |
| CL2  | W   | (claims enabled) One `claimed_by` holds more than one `active` item.                                                                                                                                                                                                                                                             |
| SF1  | E   | (safety enabled) A record, `AGENTS.md` or `ROADMAP.md` matches a built-in secret pattern: private-key header, cloud access key id, token prefix, JWT shape, `password:` assignment. The finding prints the file, line and rule, never the matched text.                                                                          |
| SF2  | E   | (safety enabled) A file matches a line of `docs/.akm-deny`. Same output rule.                                                                                                                                                                                                                                                    |
| SF3  | W   | (safety enabled) A high-entropy string that matches no pattern.                                                                                                                                                                                                                                                                  |
| SF4  | W/E | (safety enabled) An inline allow marker `akm-allow: SF1 <reason>` (also SF2, SF3) on a line suppresses that rule there and is itself listed as a warning with its reason. A marker with no reason is an error and suppresses nothing.                                                                                            |
| S1   | W   | The manifest (`docs/.layout`) reconciles as anything other than `current` (`behind`, `ahead`, `unsupported`, `conflict`); the finding names the steps or the fix. Absent manifest: warning.                                                                                                                                      |
| S2   | E   | Each skill under `.agents/skills/` has a `SKILL.md` with `name` equal to its directory and a `description`.                                                                                                                                                                                                                      |
| S3   | E   | `work_dir`, `adr_dir`, `roadmap` and `git_commits` in the manifest are valid (section 6, "Settings"). The finding names the key; until it is fixed the defaults are read and the tool refuses writes.                                                                                                                            |
| S4   | W   | `git_commits` is `true` and `.agents/skills/akm-commit` is missing, or `git_commits` is `false` or absent and the directory exists. `akmctl config git-commits on\|off` fixes it.                                                                                                                                                |

Not checked, on purpose: whether a record is true; whether a Log line is honest; whether an agent actually ran the check (the agent is trusted to run the tool); `accepted` ADR bodies; whether a `## Goal` has a done-when list or an `## Outcome` gives evidence for each line (only that the sections exist and are non-empty); who wrote a record.

### Formatter

The standard requires Prettier 3.x, pinned by the template, with `.prettierrc.json` containing `{"embeddedLanguageFormatting": "off"}` so Prettier leaves frontmatter (long inline lists and long `summary` values included) untouched. The generated index is listed in `.prettierignore`, and the generator also emits the form Prettier leaves unchanged. There is no `fmt` command: the tool validates structure and Prettier owns style. Node is needed only to run Prettier. The template pins the latest 3.x at the time it ships, and each bump of the pin is re-tested against the tool's Prettier fixtures (the index and width tests).

### Settings

Optional flat keys in the manifest (`docs/.layout`), in the manifest's format. Three keys place the records: `work_dir` (default `docs/work`), `adr_dir` (`docs/adr`) and `roadmap` (`docs/ROADMAP.md`, a `.md` file). A value is a path relative to the repository root, with `/` as the separator; it is not absolute, has no drive letter or backslash, no empty, `.` or `..` segment, and no control character, and it may contain spaces. `work_dir` and `adr_dir` must not be the same folder or inside one another, and `roadmap` must not be inside either. A value that breaks these is rule S3. The generated ADR index is `<adr_dir>/README.md`. The other keys: `stale_days` (default 30; rule WK10), `claim_stale_hours` (24), `agents_md_warn_lines` (80) and `agents_md_max_lines` (150; rule G2) and `reference_style` (`id` by default, or `wikilink`; section 0, "Filenames and references"), and `git_commits` (`true` or `false`, default `false`; whether the optional skill `akm-commit` is installed, which tells an agent to commit after each logical step; any other value is rule S3, and a mismatch with `.agents/skills/akm-commit` is rule S4). `git_commits` is a MINOR addition to layout 0.2.0: an absent key means `false`, so no repository needs a migration. The status enum and the roadmap sections are fixed in v1. Absent keys mean the defaults.

### Extensions in this version

A repo opts in by listing an extension in the manifest (for example `extension.claims: 0.3.0`).

- `claims`: section 3, rules CL1 and CL2, commands `claim` and `release`.
- `safety`: rules SF1 to SF4 and the deny list `docs/.akm-deny`: one entry per line, a literal substring, or a regular expression when the line starts with `re:`; `#` lines and blank lines are ignored. Files scanned: records, `AGENTS.md` and `ROADMAP.md`. SF1 knows private-key headers, cloud access key ids, GitHub, Slack, payment-provider, Google and npm token formats, JSON Web Tokens, and `password`, `secret`, `api_key` and token assignments whose value is at least 8 characters with two kinds of character and is not an obvious placeholder (`changeme`, `example`, `${...}`, `<...>`, `xxxx`). SF3 flags a run of 32 or more base64 characters with upper case, lower case, a digit and at least 4.5 bits of entropy per character, skipping all-hex strings (hashes), values after `h1:` or `sha256-` style prefixes, and anything inside a URL. Findings print the file, line and rule and the kind of pattern, never the matched text. A pass means "no known pattern", not "safe". Measured on about 15,000 lines of ordinary notes, the rules raised no finding once hashes and URLs were skipped.

## 7. Distribution and versions

One tool, one repository, no template repository. The standard, the guide, the API spec, the fixtures and the tool live together.

- **The tool.** The Go binary and UI, `akmctl`. Released as **prebuilt binaries**, so installing needs no Go toolchain; building from source is possible. It has its own version.
- **New repositories.** `akmctl init` writes the layout, the `AGENTS.md` and record templates, the skills with their instructions and the manifest into a new directory. It never overwrites an existing file. Its options are `--docs-dir <name>` (a single folder name, default `docs`: the folder that holds the records, the roadmap and the manifest; a name other than `docs` is written to the manifest as `work_dir`, `adr_dir` and `roadmap`), `--wikilinks` (writes `reference_style: wikilink`), `--git-commits` and `--git-init`. It writes a `.gitignore` only when the directory is, or becomes, a git work tree: it lists the lock file, and Obsidian's per-machine files (`.obsidian/workspace.json`, `.obsidian/workspaces.json`, `.trash/`) only when an `.obsidian/` folder already exists. The `AGENTS.md` it writes has no build, test, lint or run lines for the maintainer to fill in. It refuses to create a second manifest next to an existing one in another folder.
- **Existing repositories.** `akmctl migrate` brings a repository to the layout (below). Adoption of an existing repository is the migration from layout `0.0.0`. In the 0.3.0 beta, `init` is the supported way in; `migrate` is documented but not yet supported for adoption.

### Versions, manifest and reconciliation

A repo records which layout and which extensions it conforms to, so it can be reconciled with a newer one.

**Manifest.** `docs/.layout`, in the same simple key-value subset as frontmatter (H5), no nesting and no delimiter lines:

```text
layout: 0.3.0
extension.adr-index: 0.3.0
extension.work-items: 0.3.0
stale_days: 30
```

**Extensions.** Each extension carries an `EXTENSION` file with `name`, `version` and `requires_layout` (a range such as `>=1.2.0 <2.0.0`). A repo opts in by listing the extension in its manifest.

**Version rules** (semantic versioning, for the layout and for each extension):

- PATCH: a fix; no change to a repo.
- MINOR: a new optional field, rule (off or a warning by default) or extension; backward compatible, so no migration.
- MAJOR: may change the schema or meaning of a rule. Only a MAJOR step needs a migration.
- An extension's MAJOR may require a higher layout version (`requires_layout`).

**Reconcile.** The tool compares the manifest with the layout and extension versions it supports and reports one of: `current`; `behind` (a migration path exists, with its steps); `ahead` (the repo is newer than the tool: upgrade the tool); `unsupported` (no path); `conflict` (an extension needs a layout the repo is not on). Rule S1 reports this.

**Migration (`akmctl migrate`).** A MAJOR step is data: an ordered list of steps in a file `migrations/<from>-to-<to>.json` (format in the API spec, section 4.15). Each tool release embeds the migrations it supports, so migrating needs no network and a newer migration needs a newer tool. Three parts, in this order:

1. **Conformance report.** `akmctl check --target <version>` classifies each finding: `auto` (a step can fix it), `manual` (needs a judgement), `info`. Output is structured JSON. This is the static analysis, and it is useful by itself.
2. **Safe steps (codemod-style).** Steps from a small set of declarative primitives: rename, copy, remove or default a frontmatter field, set a title from the first heading, scaffold a `## Goal`, rename a section heading, move or rename a file (keeping history), add a file from a template, mint ids (the list is in the API spec, section 4.15). Properties: dry run by default, `--mode apply-safe` to write; idempotent; never deletes or overwrites; preserves unknown fields and prose bodies; writes by temp file and rename; each applied step recorded in an append-only log so the change is traceable. Steps that need real code are written in the tool, but the primitives cover the common cases and are auditable and testable on fixtures.
   A repository with no layout history (no manifest, or one without `layout`) can run a one-off migration file the maintainer supplies (`akmctl migrate --migration <file>`; API spec, section 4.15): the same steps and the same safeguards, read from the working tree instead of the binary.

3. **Handoff.** `akmctl migrate --mode handoff` returns an agent prompt for the `manual` remainder. It contains the target spec pointers, the report, the steps left with file paths, a stop condition (`akmctl check --target` clean), and the safeguards below.

**Safeguards** (in the tool before it writes, and again in the prompt):

- Refuse unless the git tree is clean; create a backup tag (`pre-migrate-<from>-to-<to>`) and work on a new branch.
- Require an explicit acknowledgement that a backup exists and the maintainer accepts the risk of data loss.
- Stop and ask on any ambiguity; never delete a file; show the diff before any commit.
- The acknowledgement is a safety step, not a legal waiver. Liability is handled by the tool's licence, and the wording of any disclaimer is the maintainer's call.

The manifest's `layout` is the version of this specification's layout. A repo records any deliberate deviation from the layout in its `AGENTS.md`; its `AGENTS.md` wins.

The tool itself is installed separately (above); `AGENTS.md` names the minimum tool version; until the tool is installed, the plain-file fallbacks apply (H6).

## 8. Tiers

| Tier     | When                                    | Has                                                                                                                                                                             |
| -------- | --------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Minimal  | A small tool or library; one maintainer | `AGENTS.md`, `docs/adr/` (index), `docs/work/` with a few items, `docs/ROADMAP.md`, `docs/.layout`. No `architecture.md`. The tool still runs.                                  |
| Standard | Most repos                              | Everything in section 1.                                                                                                                                                        |
| Large    | Multi-component or team                 | Records stay flat at the repo root. Components get their own nested `AGENTS.md` (nearest wins) and are tied to records with `scope` and `parent`. Everything else is unchanged. |

A repo moves up a tier when its records outgrow a handful of items or when a new maintainer cannot find their way in a first read. Moving up is adding files, never reshaping the schema.

## 9. The contract and how the standard meets it

A repository's memory exists so that a cold reader (an agent with no memory, or a maintainer back after months away) can work well and not repeat past mistakes. It lets that reader:

| #   | Clause                        | Meaning                                                                                                                                                                                                                               |
| --- | ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| C1  | Orientation                   | Learn what the project is, how to build, test and run it, where things live.                                                                                                                                                          |
| C2  | Rationale                     | Learn why it is the way it is, including rejected alternatives, so settled decisions are not re-litigated.                                                                                                                            |
| C3  | Direction                     | Learn what is planned, in what order, and what was dropped.                                                                                                                                                                           |
| C4  | Current state                 | Learn what is in flight, where it stands, the single next step; work is resumable.                                                                                                                                                    |
| C5  | Boundaries                    | Learn what must not be touched or done without asking.                                                                                                                                                                                |
| C6  | Procedures                    | Repeat recurring work the same way each time.                                                                                                                                                                                         |
| C7  | Trust                         | Records stay accurate; drift is prevented or detected; it is clear what may be edited and by whom.                                                                                                                                    |
| C8  | Cost                          | Cheap to read cold, cheap to keep up.                                                                                                                                                                                                 |
| C9  | Authority and ownership       | A declared source of truth when documents conflict (code and spec win over prose; the repo's `AGENTS.md` wins over a skill); who decides (agents propose, a maintainer accepts); a human review surface (what waits on a maintainer). |
| C10 | Verification                  | Records say what was and was not verified; nothing is "done" because code exists; done has evidence.                                                                                                                                  |
| C11 | Content safety and provenance | What a record may never contain (credentials, private or confidential material, copyrighted text); which records an agent wrote and which a maintainer has reviewed.                                                                  |
| C12 | Portability and handoff       | The setup works across tools and machines without relying on one tool's private configuration; modes, claiming and concurrent threads are handled.                                                                                    |

The mechanisms of this standard that serve each clause:

| Clause                            | Mechanisms                                                                                                                                  |
| --------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| C1 Orientation                    | Section 2 `AGENTS.md` parts; `docs/architecture.md`; G3 command check                                                                       |
| C2 Rationale                      | ADRs with considered options and Confirmation; A1 to A4                                                                                     |
| C3 Direction                      | Roadmap with Dropped and reasons; RM1 to RM4                                                                                                |
| C4 Current state                  | `## Next` per item; `ls` listing; WK5 to WK7, WK10                                                                                          |
| C5 Boundaries                     | `AGENTS.md` Do-not list linking policies; tool-level controls advised                                                                       |
| C6 Procedures                     | Shared skill with templates; the tool; S1, S2                                                                                               |
| C7 Trust                          | One tool with checks; conventions on immutability; WK1 to WK12, L1, L2                                                                      |
| C8 Cost                           | 60 to 80 line `AGENTS.md`; history on demand; per-item files; G2                                                                            |
| C9 Authority and ownership        | `AGENTS.md` wins over skill; code wins over prose; maintainer accepts ADRs and orders the roadmap; blocked listing shows what waits on whom |
| C10 Verification                  | `## Outcome` evidence for each "done when" line (WK7)                                                                                       |
| C11 Content safety and provenance | `authored_by`; policy docs linked from Do-not; optional `safety` extension (known patterns and a deny list)                                 |
| C12 Portability and handoff       | Plain files, a single binary with plain-file fallbacks, `CLAUDE.md` as a one-line import; optional `claims` extension for active items      |

### Gaps this standard leaves

- **Concurrency (C12):** claims are optional and advisory; two branches can claim the same item and meet only at merge.
- **Content safety (C11):** the optional `safety` extension catches known patterns and declared strings only. A pass is not proof of safety.
- **Silent skipped checks (C7):** the agent runs the tool; nothing forces it.
- **Truth:** a well-formed, wrong record passes every rule.
