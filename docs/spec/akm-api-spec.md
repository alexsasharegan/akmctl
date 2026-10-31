---
type: design
status: active
created_at: 2026-10-28
updated_at: 2026-10-28
authored_by: agent
---

# AKM API specification

Version 1.0.

A transport-neutral contract for reading, checking and changing the records of an AKM repository (ADRs, work items, the roadmap), so that more than one implementation can be written against it. The records and rules are defined in the specification ([agent-ready-repo-spec](agent-ready-repo-spec.md)); this document defines the operations over them, the JSON that comes back, and two bindings (CLI and HTTP). The reference implementation is `akmctl`. The schemas are in `akm-api/schemas/` and the conformance fixtures in `akm-api/.fixtures/`.

The words MUST, SHOULD and MAY are used as in RFC 2119.

## 1. Conventions

- **One operation set.** Every operation has a request (named parameters) and a response (one JSON object). A binding only maps its transport to those parameters and back.
- **Names.** Parameters and response fields are `snake_case`. Frontmatter field names are used unchanged for record fields (`blocked_by`, `created`).
- **Dates.** Dates are strings `YYYY-MM-DD` (the local day of the machine that writes them). Timestamps are strings `YYYY-MM-DDTHH:MM` (UTC). Neither is ever a number or an object.
- **Null.** A field defined for a record type but absent from the file is `null` in the response, never omitted. A field defined for neither type is present only if the file has it and the request asked for it (unknown fields, see section 3).
- **Paths.** Relative to the repository root, `/` as separator, UTF-8, no leading `./`.
- **Text.** UTF-8, LF line endings. A tool reading CRLF files reports them as they are.
- **Order.** Arrays have a defined order. Where this spec does not say, the order is by `id` ascending.
- **Unknown fields** in a response MUST be ignored by a client. Unknown request parameters MUST be rejected with `invalid_request`.
- **No side effects without a write operation.** The operations `new`, `set`, `section`, `close`, `claim`, `release`, `index` (unless `dry_run`), `config` (with a `value`, unless `dry_run`), `skills` (with `mode` `install` or `update` and `confirm`, unless `dry_run`), `init` and `migrate` (mode `apply_safe`) write files. All others read only (H6 of the specification). A write operation with `dry_run` true writes nothing (section 7). No operation runs `git commit`, `git push` or any network call. `migrate` is the only operation that writes to git (a backup tag and a branch, section 4.15); `commits` only reads it (section 4.19). While the manifest holds an invalid `work_dir`, `adr_dir` or `roadmap` (rule S3), every write operation is `refused` with `details.reason: "rule"` and `details.rule: "S3"`, because the default folders would be written instead of the intended ones; reads and `check` still run.
- **The `commit` member.** `new`, `claim`, `release`, `close` and a `set` that changes a work item's `status` return `commit`: `{step, item, detail, paths, message}` (ADR 261030-tbkh). The tool describes the commit and never makes it. `step` is one of `file` (`new`), `claim`, `release`, `close` (`done`), `drop` (`dropped`), and, for `set`, `ready` (to `ready`), `start` (to `active`), `block` (to `blocked`) and `unblock` (from `blocked`); a `set` that changes no status, `section`, `index` and every read return none. `item` is the id. `detail` is `claimed by <as>`, `released`, `done`, `dropped`, `ready`, `started`, `unblocked`, `blocked on <the first blocked_by entry>`, and for `file` the title. `paths` are the repository-relative files the step wrote (a `close` or `drop` lists the old path and the new one). `message` is the subject `<item>(<step>): <detail>`, a blank line and the two trailers `AKM-Step: <step>` and `AKM-Item: <item>`, with LF line ends and no final newline. The caller commits with it as given, may put a body between the subject and the trailers, and may add trailers after them. A `dry_run` response carries it too. The member is additive: a client that ignores it is unaffected.

## 2. Envelope and errors

Every response is one JSON object:

| Field              | Type    | Meaning                                                                                     |
| ------------------ | ------- | ------------------------------------------------------------------------------------------- |
| `api_version`      | string  | `"1.0"`: the version of this spec the response conforms to.                                 |
| `operation`        | string  | The operation name.                                                                         |
| `request`          | object  | The request exactly as given, with no defaults added. The defaults are in this spec.        |
| `ok`               | boolean | `true` if the operation ran. A `check` that finds errors is still `ok: true` (see `clean`). |
| _operation fields_ |         | Defined per operation in section 4. Present only when `ok` is `true`.                       |
| `error`            | object  | Present only when `ok` is `false`: `code`, `message`, optional `details`.                   |

Error codes (the closed set for 1.0):

| Code              | Meaning                                                                                                                                                                                                                                                               | Exit | HTTP                             |
| ----------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---- | -------------------------------- |
| `invalid_request` | An unknown or malformed parameter, a value outside its enum, a missing required parameter.                                                                                                                                                                            | 2    | 400                              |
| `unauthorized`    | HTTP only: the token is missing or wrong.                                                                                                                                                                                                                             | n/a  | 401                              |
| `not_found`       | An id, section or path does not resolve.                                                                                                                                                                                                                              | 2    | 404                              |
| `ambiguous`       | A reference resolves to more than one record (a collision, rule I1). `details.paths` lists them.                                                                                                                                                                      | 2    | 409                              |
| `stale`           | A write carried `if_match` and the file changed since (section 7). `details.etag` is the file's current etag.                                                                                                                                                         | 2    | 412                              |
| `refused`         | A precondition failed: closing with an empty `## Outcome`, an unclean git tree on `migrate`, a claim held by someone else, a status move that is not allowed, a section write to a closed record, an HTTP write with no `If-Match`. `details.reason` is a short code. | 2    | 409; 428 for `if_match_required` |
| `unreadable`      | The repository, a file or the manifest could not be read or parsed in a way that stops the operation.                                                                                                                                                                 | 2    | 500                              |
| `internal`        | A bug.                                                                                                                                                                                                                                                                | 2    | 500                              |

Exit codes of the CLI binding: `0` for `ok: true` with `clean: true` (or an operation with no `clean` field); `1` for `ok: true` and `clean: false` (only `check` and `migrate`); `2` for any `ok: false`. This keeps the specification's `0 / 1 / 2` meaning (clean or warnings only / errors found / could not read).

## 3. The record object

`ls`, `search`, `show`, `board` and the write operations return records:

| Field     | Type   | Meaning                                                                                                    |
| --------- | ------ | ---------------------------------------------------------------------------------------------------------- |
| `id`      | string | `YYMMDD-xxxx`.                                                                                             |
| `type`    | string | `work` or `adr`.                                                                                           |
| `path`    | string | The file, for example `docs/work/261027-m3qz-add-json-output.md`.                                          |
| `id_date` | string | The id's date, `YYYY-MM-DD`.                                                                               |
| `etag`    | string | The first 16 lowercase hex characters of the SHA-256 of the file's bytes. Used for `if_match` (section 7). |
| `fields`  | object | The frontmatter.                                                                                           |

`fields` holds frontmatter. `show` and the write operations return every field of the record type, `null` when absent, plus unknown fields the file has. `ls`, `search` and `board` return a short default set (below) unless the request lists fields. The full list, for a work item: `id`, `title`, `summary`, `status`, `kind`, `parent`, `blocked_by`, `priority`, `created`, `updated`, `authored_by`, `decisions`, and with the claims extension `claimed_by` and `claimed_at`; for an ADR `id`, `title`, `summary`, `status`, `date`, `scope`, `supersedes`, `superseded_by`, `authored_by`. Lists are arrays of strings; `priority` is an integer; everything else is a string.

References (`parent`, `blocked_by` entries that are ids, `decisions`, `supersedes`, `superseded_by`) are always returned as the bare id, whichever form the file uses. Free-text `blocked_by` entries are returned as they are. A client that wants to know the file's form reads the file.

Unknown frontmatter fields (a team's own) are returned in `fields` only when the request lists them in `fields`, or asks for `"*"`, which returns every field the file has.

**Default set** for `ls`, `search` and `board` (keys always present, `null` when absent): work item `id`, `title`, `summary`, `status`, `kind`, `parent`, `priority`, `updated`; ADR `id`, `title`, `summary`, `status`, `date`, `scope`. `board` adds `blocked_by`, `cycle` and `keep` to work items. The `fields` parameter replaces the set; `["*"]` returns every field of the type. The set leaves out the fields that are mostly empty or rarely routing-relevant (`created`, `closed`, `blocked_by`, `cycle`, `keep`, `decisions`, `authored_by`, the claim fields; ask for them with `fields`) to keep a row small.

Header fields only: `ls`, `search` and `board` never return the body (H1). `show` returns it.

Two derived members MAY appear on a record, only when asked: `next_line` (the first line of `## Next`, from `ls` with `next: true`) and `matches` (from `search`).

## 4. Operations

For each operation: the request parameters (all optional unless marked **required**), the response fields (besides the envelope), the errors, and the side effect. Schemas: `akm-api/schemas/<operation>.json`, with `$defs.request` and `$defs.response`.

### 4.1 `version`

Request: none. Response: `tool` (`name`, `version`), `api_versions` (array of supported versions), `layout_versions` (array of layout versions the tool can check and migrate between), `extensions` (array of extension names it implements) and `repository`. CLI: `akmctl version`.

`repository` describes the repository the call ran against, so that a client can adapt without reading the manifest. It is an object with `name` (the name of the repository's root directory), `layout` (the `layout` version in the manifest, `docs/.layout` by default, or `"0.0.0"` when the manifest is missing), `extensions` (the names of the extensions the manifest enables, in ascending order), `reference_style` (`id` or `wikilink`; `id` when the manifest sets none) and `claim_stale_hours` (integer; 24 when the manifest sets none), `work_dir` and `adr_dir` (the folders of work items and of ADRs, repository-relative, the manifest's `work_dir` and `adr_dir` or `docs/work` and `docs/adr`; closed items are in `<work_dir>/done`) and `git_commits` (boolean; the manifest's `git_commits`, `false` when it sets none). A client builds paths from these instead of hard-coding them. It is `null` when the tool runs outside any repository, which only the CLI can do. `tool`, `api_versions`, `layout_versions` and `extensions` describe the implementation, not the repository.

### 4.2 `ls`

List records, header fields only.

| Parameter    | Type              | Default         | Meaning                                                                                                                                                                                                                                                                                                                                                                           |
| ------------ | ----------------- | --------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `type`       | `work` or `adr`   | both            |                                                                                                                                                                                                                                                                                                                                                                                   |
| `status`     | array of status   | any             | Matches if the record's status is in the list. Work and ADR statuses are both allowed.                                                                                                                                                                                                                                                                                            |
| `parent`     | id                | any             | Work items whose `parent` is this id.                                                                                                                                                                                                                                                                                                                                             |
| `blocked`    | boolean           | any             | `true`: work items with a non-empty `blocked_by`. `false`: those without.                                                                                                                                                                                                                                                                                                         |
| `scope`      | string            | any             | ADRs with a matching `scope` entry. An entry ending in `/` is a path prefix: it matches when the value starts with it (`src/store/open.go` matches `src/store/`). Any other entry matches only when equal to the value (`cli` matches `cli`).                                                                                                                                     |
| `cycle`      | string            | any             | Work items whose `cycle` equals this label (`[A-Za-z0-9][A-Za-z0-9._-]{0,63}`; anything else is `invalid_request`).                                                                                                                                                                                                                                                               |
| `since`      | date              | none            | Records whose `updated` (work) or `date` (ADR) is on or after it.                                                                                                                                                                                                                                                                                                                 |
| `sort`       | array of sort key | `["id"]`        | A key is a field name, `-` before it for descending, or `roadmap`. Sortable: `id`, `title`, `status`, `kind`, `priority`, `created`, `updated`, `date`, `roadmap`. Absent values sort last in both directions. `roadmap` sorts by position in `ROADMAP.md` (Now, Next, Later in file order), records not on it after, so `["roadmap", "-priority", "id"]` is the open-work order. |
| `fields`     | array of string   | the default set | Which keys of `fields` to return. `"*"` for all fields of the type.                                                                                                                                                                                                                                                                                                               |
| `limit`      | integer ≥ 1       | none            |                                                                                                                                                                                                                                                                                                                                                                                   |
| `offset`     | integer ≥ 0       | 0               |                                                                                                                                                                                                                                                                                                                                                                                   |
| `count_only` | boolean           | false           | Return `total` and `count` with no records.                                                                                                                                                                                                                                                                                                                                       |
| `next`       | boolean           | false           | Add `next_line` to work items.                                                                                                                                                                                                                                                                                                                                                    |

Response: `total` (matches before `limit`), `count` (records returned), `records`. An unknown `type` or `status` value is `invalid_request`. A `status` that does not apply to the chosen `type` matches nothing.

### 4.3 `search`

| Parameter                                     | Type    | Default      | Meaning                                                                        |
| --------------------------------------------- | ------- | ------------ | ------------------------------------------------------------------------------ |
| `text`                                        | string  | **required** | A case-insensitive substring. Whitespace-separated words must all match (AND). |
| `body`                                        | boolean | false        | Also search bodies. Otherwise title, summary and scope.                        |
| `type`, `status`, `fields`, `limit`, `offset` |         |              | As for `ls`.                                                                   |

Response: as `ls`, with `matches` on each record: an array of `{field, line, snippet}`; `field` is a frontmatter field name or `body`; `line` is the 1-based line in the file (frontmatter lines count); `snippet` is the matching line, trimmed to 160 characters. A record matches when every word appears in one of title, summary or scope (or in the body with `body`). `matches` has one entry for each of those fields, or for each body line, that contains at least one of the words. Records are ordered by number of matches descending, then `id`.

### 4.4 `show`

| Parameter | Type   | Meaning                                                                                                          |
| --------- | ------ | ---------------------------------------------------------------------------------------------------------------- |
| `id`      | id     | **required**                                                                                                     |
| `section` | string | A section heading, case-insensitive, without the `##` (`next`, `goal`, `log`, `outcome`, `context`, `decision`). |

Response: `record` (with all fields the file has), `sections` (the file's `##` headings in order, lowercase), and either `body` (the whole body, if `section` is absent) or `section` (`{name, text}`, the text without its heading). A missing section is `not_found`.

### 4.5 `board`

The Kanban grouping, as data. A client draws it; `akmctl board` draws it as text.

| Parameter  | Type                 | Default                        | Meaning                 |
| ---------- | -------------------- | ------------------------------ | ----------------------- |
| `statuses` | array of work status | `["active","blocked","ready"]` | Columns, in this order. |

Response: `columns` (array of `{status, count, records}`, in the order requested; `records` sorted `roadmap`, `-priority`, `id`) and `totals` (an object with the count of every work status, including those with no column). `blocked_by` is in the default set here, so a client can group blocked items by what they wait on; `cycle` and `keep` are too, so a client can group by cycle and mark pinned items.

### 4.6 `check`

Run the rules of section 6 of the specification.

| Parameter | Type              | Default        | Meaning                                                                                             |
| --------- | ----------------- | -------------- | --------------------------------------------------------------------------------------------------- |
| `target`  | version string    | the tool's own | Run the rules of another layout version and classify each finding (section 7 of the specification). |
| `rules`   | array of string   | all            | Rule ids (`WK4`) or prefixes (`WK`, `L`).                                                           |
| `as_of`   | timestamp or date | now            | The clock for staleness rules (WK10, claims). Tests set it.                                         |

Response: `clean` (true when no finding has severity `E`), `errors`, `warnings` (counts), `count` (findings), `findings`, and `reconcile` (`state`: `current`, `behind`, `ahead`, `unsupported`, `conflict` or `absent`; and `steps`, an array of strings, when `behind`).

A finding is `{rule, severity, file, line, id, message, class}`: `severity` is `E` or `W`; `file` is a path or `null`; `line` is the 1-based line of the offending field or entry when there is one, and `null` for a finding about a whole file or a missing section; `id` is the record id or `null`; `message` is one sentence; `class` is `auto`, `manual` or `info` when `target` is given and `null` otherwise. Safety findings (SF1 to SF3) never include the matched text in `message`. Findings are ordered by `file`, `line`, `rule`.

### 4.7 `index`

Regenerate `docs/adr/README.md` (section 4 of the specification).

| Parameter         | Type    | Default | Meaning                    |
| ----------------- | ------- | ------- | -------------------------- |
| `dry_run`         | boolean | false   | Compute, do not write.     |
| `include_content` | boolean | false   | Return the generated text. |

Response: `path`, `changed` (the generated text differs from the file), `written` (the file was written), and `content` when asked. Side effect: writes the file unless `dry_run`.

### 4.8 `find`

Gives the path of a record by `id`, or of a note by `name`. Exactly one of the two is needed.

- `id` returns `{id, type, path}`. `not_found` if no file has it; `ambiguous` if several do.
- `name` is the target of a wikilink: a file name without `.md`, optionally with the folders it sits in (`Projects/Vault Views`). The tool looks at every Markdown file under the repository root, skipping hidden folders (`.git`, `.obsidian`, `.trash`, `.agents`) and `node_modules`, and never follows a link. The match ignores case and takes a trailing run of folders. It returns `{name, path}`, and `id` and `type` as well when the file is a record. `not_found` if no note has the name; `ambiguous` (`details.paths`) if several do. A name that is empty or holds `..` or a backslash is `invalid_request`.

`id` and `name` together are `invalid_request`. CLI: `akmctl find <id>` and `akmctl find --name <note>`. HTTP: `GET /find/{id}` and `GET /find?name=<note>`.

### 4.9 `when`

| Parameter | Type | Default       | Meaning                                                             |
| --------- | ---- | ------------- | ------------------------------------------------------------------- |
| `id`      | id   | **required**  | Need not exist as a record; any string that matches the id pattern. |
| `as_of`   | date | today (local) |                                                                     |

Response: `id`, `date`, `weekday` (English name), `age_days` (integer, negative if the id is dated after `as_of`).

### 4.10 `new`

| Parameter | Type            | Meaning                                                                                                                                                      |
| --------- | --------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `kind`    | `work` or `adr` | **required**                                                                                                                                                 |
| `slug`    | string          | **required**. Lowercase kebab-case, at most 60 characters.                                                                                                   |
| `title`   | string          | Default: the slug with hyphens as spaces.                                                                                                                    |
| `fields`  | object          | Initial values: for `work`, `summary`, `kind`, `parent`, `priority`, `blocked_by`, `decisions`; for `adr`, `summary`, `scope`. Others are `invalid_request`. |
| `as_of`   | timestamp       | The clock for the id's date. Tests set it.                                                                                                                   |

Response: `record` and `path`. The new file comes from the template, with a fresh id, `created` and `updated` (or `date`) set to today, status `raw` (work) or `proposed` (ADR). A `parent` that does not resolve is `not_found`. The id is chosen so that no existing file in the three directories has it. For `kind: work`, the response also has `similar` when open items resemble the new one, and omits it otherwise (a response with no match is as it was before `similar`). `new` never refuses because of a match.

The matching rule: a _word_ is a maximal run of letters and digits (Unicode), lower-cased; the stop words `a an and are as at be by for from in is it of on or the to with` are dropped, and a word counts once. The new item's words are those of the slug (hyphens as spaces) and of the title. A candidate is a work item whose `status` is `raw`, `ready`, `active` or `blocked` (a `done` or `dropped` item never matches), and its words are those of its `title` and `summary`. The candidate matches when it shares at least two words with the new item, or every word when the new item has fewer than two. No stemming. The matches are ordered by the number of shared words, most first, then by id, and the first five are returned. Each is `{id, status, title, path, words}`, `words` being the shared words in the order of the new item. The new item itself is never a match. `new adr` returns no `similar`.

Side effect: writes one file.

### 4.11 `set`

| Parameter  | Type                   | Meaning                                                             |
| ---------- | ---------------------- | ------------------------------------------------------------------- |
| `id`       | id                     | **required**                                                        |
| `changes`  | object                 | **required**. Field to new value. A `null` value removes the field. |
| `if_match` | etag                   | The record's `etag` as the client last saw it (section 7).          |
| `dry_run`  | boolean, default false | Compute and return the change; write nothing (section 7).           |
| `as_of`    | date or timestamp      | The clock for `updated`. Tests set it.                              |

Response: `record` (after) and `changed` (the field names that differ), and `diff` when `dry_run` is true. Allowed fields are those of the record type except `id`, `created`, `closed` (work; only `close` sets it), `status` values `done` and `dropped` for work (use `close`) and `claimed_by`/`claimed_at` (use `claim`). `updated` (work) is set to today whenever the file changes; `date` (ADR) only when `status` changes. A value outside the field's type or enum is `invalid_request`: a work item's `cycle` is a string matching `[A-Za-z0-9][A-Za-z0-9._-]{0,63}` and `keep` is a boolean (a `null` removes either). A change that breaks a rule the tool can see from the record alone is `refused` with `details.reason: "rule"` and `details.rule` naming it: `WK5` (status `blocked` with no `blocked_by`), `WK4` (a record as its own parent), `WK9` (a closed item given an open status) and `A3` (an ADR set to `superseded` with no `superseded_by`). Content rules such as WK6 (an empty `## Next`) are left to `check`, because an agent sets the status first and writes the section next. Unknown fields cannot be set; `date` (ADR), `updated` and `closed` (work) are managed by the tool, so a change to one is `invalid_request` with `details.parameter` `changes` (`changes: "closed" cannot be set on a work (id, created, updated, closed, date and claim fields are managed by the tool)`). Unchanged values are not written, so `changed` can be empty. Only frontmatter fields change here; the body is changed by `section` (section 4.16). Side effect: rewrites one file, unless `dry_run`.

### 4.12 `close`

| Parameter                      | Type                | Meaning                                                                                                                                                                           |
| ------------------------------ | ------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id`                           | id                  | **required**; a work item.                                                                                                                                                        |
| `status`                       | `done` or `dropped` | **required**                                                                                                                                                                      |
| `outcome`                      | string              | The text of the `## Outcome` section. When given and not blank, it replaces the section (or creates it, as `section` does) in the same write as the status change and the rename. |
| `if_match`, `dry_run`, `as_of` |                     | As for `set`.                                                                                                                                                                     |

Response: `record`, `moved_from`, `moved_to` (paths), and `diff` when `dry_run` is true (its `new_path` is `moved_to`). Sets the status, removes any `claimed_by` and `claimed_at`, sets `updated` and `closed` (both to the local day, or the `as_of` day; `closed` is placed on the line after `updated`, for `done` and `dropped` alike), and moves the file to `done/` inside the work folder (`docs/work/done/` by default) with a plain rename (Git sees a rename once the change is staged; the tool does not stage, section 7). `refused` with `details.reason: "empty_outcome"` if `## Outcome` is empty or missing after `outcome` is applied (neither the file nor the parameter supplies one); `refused` with `details.reason: "open_children"` when closing `done` a parent with children not `done` or `dropped` (WK8). The response also has `dependents` when open work items name the closed item in `blocked_by`, and omits it otherwise: an array of `{id, status, title, path}`, ordered by id, of the work items whose status is `raw`, `ready`, `active` or `blocked` and whose `blocked_by` holds the closed item's id (bare or as a wikilink). It is the same for `done` and `dropped` and appears in a `dry_run` response. `close` only reports: it edits none of those items and is never refused because of them; the caller unblocks or re-points each one (`set`). Free-text entries in `blocked_by` are not matched. Side effect: rewrites one file and renames it, unless `dry_run`.

### 4.13 `claim`

Claims extension. `id` (**required**), `as` (**required**, a free string), `if_match`, `takeover` (boolean, default false), `dry_run`, `as_of` (timestamp; the clock). Sets `status: active` if the item is `ready`, sets `claimed_by` and `claimed_at`, appends a `claimed` line to `## Log`, and sets `updated`. `refused` with `details.reason: "claimed"` when another holder's claim is not stale (older than `claim_stale_hours`) unless `takeover` is true; a takeover adds a `took over from <name>` line. `refused` with `details.reason: "extension_off"` when the manifest does not enable `claims`. Response: `record`, `previous_holder` (string or `null`), and `diff` when `dry_run` is true.

### 4.14 `release`

`id` (**required**), `if_match`, `dry_run`, `as_of`. Clears `claimed_by` and `claimed_at`, appends a `released` line to `## Log`, sets `updated`. The status is not changed. Response: `record`, and `diff` when `dry_run` is true.

### 4.15 `migrate`

Reconcile the repository with a layout version (section 7 of the specification).

| Parameter            | Type                                      | Default        | Meaning                                                                                                                                                                               |
| -------------------- | ----------------------------------------- | -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `mode`               | `report`, `plan`, `apply_safe`, `handoff` | `report`       |                                                                                                                                                                                       |
| `target`             | version string                            | the tool's own |                                                                                                                                                                                       |
| `acknowledge_backup` | boolean                                   | false          | Must be `true` for `apply_safe` (the safeguard of the specification).                                                                                                                 |
| `migrations`         | array of migration objects                | none           | One-off migrations to search besides the embedded ones: `{from, to, steps}` as below. The CLI flag `--migration <file>` reads a file holding one object or a list and passes it here. |
| `as_of`              | date or timestamp                         | now            | The clock for the migration log. Tests set it.                                                                                                                                        |

`from` is the `layout` version in the manifest, or `"0.0.0"` when the manifest is missing. The path from `from` to `target` is found among the migrations the tool embeds (one step per `<from>-to-<to>`, searched breadth-first); `state` is `current` (equal), `behind` (a path exists), `ahead` (the repository is newer) or `unsupported` (older, no path). This release embeds none, because no layout before 0.2.0 was published; the machinery is tested with a made-up 0.1.0. A repository with no layout history is on `"0.0.0"`; the `migrations` parameter adds a one-off migration (`{from: "0.0.0", to: "0.2.0", steps}`) to the search, so such a repository can be converted from a file kept in the working tree. The parameter is validated like an embedded file: an unknown member, step kind or class, a step without its required parameters, a path outside the repository, a repeated step id or a bad version is `invalid_request` with `details.parameter` `migrations`. Every safeguard of `apply_safe` applies unchanged, and the migration is not recorded in the repository except through `docs/.akm-migrate-log`, the tag and the branch.

Response: `mode`, `from`, `to`, `state`, `steps` (array of `{id, kind, class, files, status}`, plus `notes` when a step has any and `renames` when a `mint_ids` step lists file moves; `class` is `auto`, `manual` or `info`; `status` is `pending`, `applied` or `skipped`), `findings` (as `check`, with `class` set: `auto` when an auto step lists the rule in its `fixes`, `manual` when a manual step does or the finding is an error, `info` for a warning; only when `target` is the tool's own layout version, because this release has one rule set), `clean` (state is `current` and no `auto` or `manual` step is pending after the mode ran), and, for `handoff`, `handoff` (the prompt text).

- `report`: `state`, the steps with empty `files`, and the classified findings. Writes nothing.
- `plan`: the same with `files`, the files each `auto` step would change. Writes nothing. The steps are planned in order on a scratch copy of the repository's Markdown files, each applied there before the next is planned, so a step that depends on an earlier one (a field renamed before ids are minted from it) is planned as it will run; a `move_file` or `add_file` step is planned against the working tree.
- `handoff`: the plan plus `handoff`, the prompt for an agent: what to read, the steps left with their files, the errors of the target's check, the stop condition (`check --target` clean) and the safeguards. Writes nothing; the CLI writes it to a file only when given `--out <file>`.
- `apply_safe`: `refused` unless every safeguard holds, with `details.reason`: `no_backup_ack` (no `acknowledge_backup`), `not_git` (not a git work tree), `dirty_tree` (uncommitted changes), `tag_exists` or `branch_exists` (the names below are taken), `no_path` (state is `ahead` or `unsupported`). When the repository is already `current` it does nothing and answers `clean: true`. Otherwise it takes the lock, creates the tag `pre-migrate-<from>-to-<to>`, switches to the new branch `migrate-<from>-to-<to>`, runs the `auto` steps in order (each is idempotent and ends `applied` or `skipped`), appends one JSON line per applied step to `.akm-migrate-log` in the manifest folder (`docs/.akm-migrate-log` by default) (time, from, to, step, files), and, only when no `manual` step is pending, sets `layout` in the manifest to the target. It does not commit; the caller reviews and commits. If a step fails the response is `unreadable` with `details.tag` and `details.branch` so the state before the migration can be restored.

**Step primitives** (the `kind` of a step; all declarative, all idempotent, none deletes, none overwrites): `rename_field` (`from`, `to`; skipped when `from` is absent or `to` exists), `copy_field` (`from`, `to`; the value of `from` is copied to `to`, skipped when `from` is absent or `to` exists), `remove_field` (`field`; skipped where it is absent), `set_default` (`field`, `value`; only where the field is absent), `title_from_h1` (sets `title` from the first `# ` heading outside code fences, only where `title` is absent; a file with no heading is reported in `notes`), `goal_scaffold` (below), `rename_heading` (`## From` to `## To`, outside code fences), `move_file` (`from`, `to`; skipped when the source is missing or the target exists), `add_file` (`path`, `content`; skipped when the file exists), `mint_ids`, `move_closed` (`dir`; optional `to`, default `<dir>/done`; optional `statuses`, default `done` and `dropped`; `recursive`; moves each closed record to `to`, never overwrites, skips `README.md`, idempotent) and `manual` (a description only).

The field steps (`rename_field`, `copy_field`, `remove_field`, `set_default`, `title_from_h1`, `goal_scaffold`) choose their files in one of two ways. With `dir` (repository-relative; `recursive`, default false), every Markdown file in that folder with a frontmatter block, whatever its name (not `README.md`, not a dot file): this reaches notes that have no id yet. A file with no frontmatter block is left alone; one whose block cannot be read is reported in `notes`; a `dir` that does not exist is a note and a skipped step. Without `dir`, the records the tool has loaded, narrowed by `scope` (`work`, `adr` or `all`).

`goal_scaffold` (`dir` or `scope`; `statuses`, a list: only files whose `status` is one of them, default all) inserts `## Goal` after the first `# ` heading of a file that has no `## Goal` section. The text between that heading and the next `## ` heading moves under the new `## Goal`; if the file has a `## Done when` section, its heading becomes the lead-in line `Done when:` and its text stays, as a list, at the end of the Goal; a file with neither gets an empty `## Goal`. Every file it changes adds one entry to `notes` asking for a review of the text, because the tool cannot know whether the moved paragraph states the goal. A second run changes nothing.

`mint_ids` (`dir`, required, repository-relative; `date_field`, default `created`; `recursive`, default false) (`rename`, default true; `seed`, optional; `slug_from`, `name` (default, the file name) or `title` (the `title` field); `max_slug`, 1 to 100, default 60, a cut ends on a whole word) gives an id to each Markdown file directly in `dir` that has no `id` field (not `README.md`, not a dot file; subdirectories only with `recursive`). The id is `YYMMDD-xxxx`: the date part is the UTC date of the `date_field` value (a date or an ISO timestamp), or today's date (the `as_of` clock) when the field is absent or not a date, which a note reports; the suffix is four random characters of the id alphabet, drawn again until the id is new in the batch and among the ids of the ADR folder, the work folder, its `done/` and `dir` (file name prefixes and `id` fields). The file is renamed to `<id>-<slug>.md`, the slug being the old name without extension and date prefix, lower-case words joined by hyphens, at most 60 characters (the title when the name gives no word, then `untitled`), and `id: "<id>"` is written as the first frontmatter field, quoted as `new` writes it; a file with no frontmatter gets a block, and the note says so. A name that already starts with an id keeps it: only the field is written. A file whose frontmatter cannot be read is left alone and reported. Nothing is deleted and no reference is rewritten: links and wikilinks that name an old file name are reported. A second run changes nothing. With `rename` false no file is renamed: only the `id` field is written, and the step reports `renames`, an array of `{from, to}` (repository-relative paths) holding the moves the caller should make, for a tool that rewrites links when it renames (Obsidian); each `to` is exactly the name `rename: true` would give, and a note says the file was not renamed. With `rename` true `renames` lists the moves made. `seed` makes the draw repeatable: the suffix is the first four bytes of the SHA-256 of the seed, a NUL, the file's path, a NUL and a try counter, each byte reduced to its low five bits as an index into the alphabet; on a collision the counter is raised and the draw repeated. The same seed, files and clock give the same ids in `plan` and in `apply_safe`. Without `seed` the ids in `plan` are examples: `apply_safe` draws again. A file that has any `id` field is skipped, so a legacy `id` of another form (an ADR numbered `1`) must be renamed first: run `rename_field` (`id` to `adr_number`, with `dir`) before `mint_ids`; `plan` shows the right result only because it applies the steps in order. Because the date comes from a field, run `mint_ids` after any step that renames that field, or name the old field in `date_field`.

A step object may carry `notes`, an array of strings, present only when there is something to report: for `mint_ids`, one line per file (`old -> new (id …, date …)`), the files skipped, and the references not rewritten. The migration log carries them too. Paths must stay inside the repository. A migration file is JSON, `<from>-to-<to>.json`, `{from, to, steps: [{id, kind, class, description, fixes, ...}]}`, embedded in the binary, or supplied for one run with `migrations`.

### 4.16 `section`

Replace or extend one `##` section of a record's body. `set` changes frontmatter only; this operation is how a client writes the text of a record (the goal, the next step, the outcome, a log line).

| Parameter                      | Type                  | Default      | Meaning                                                                                                                                                                                                              |
| ------------------------------ | --------------------- | ------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id`                           | id                    | **required** | A work item or an ADR.                                                                                                                                                                                               |
| `name`                         | string                | **required** | The text of a `##` heading, compared without regard to case, without the hashes. One line, not empty.                                                                                                                |
| `text`                         | string                | **required** | The new text. May be empty. Line endings are written as the file has them.                                                                                                                                           |
| `mode`                         | `replace` or `append` | `replace`    | `replace` makes `text` the whole section. `append` adds `text` on new lines after the last non-blank line of the section; on the Log section a first line that does not start with a list marker gets a `- ` bullet. |
| `if_match`, `dry_run`, `as_of` |                       |              | As for `set`.                                                                                                                                                                                                        |

A section that does not exist is created before `## Log`, or at the end when there is no Log, with a blank line after its heading (the form Prettier leaves unchanged); `close` with `outcome` creates `## Outcome` the same way. Response: `record` (after), `section` (`{name, text}`: the section as `show` returns it after the write, `name` in lower case), `changed` (boolean: the file text differs) and `diff` when `dry_run` is true.

- **Where the text goes.** The heading line and the blank lines around the text are kept as the file has them; a section that has no blank line before the next heading gets one. Leading and trailing newlines of `text` are ignored for `replace`; trailing newlines are ignored for `append`, and a leading blank line in `text` starts a new paragraph. An `append` of empty text changes nothing. A heading inside a code fence is not a heading.
- **A missing section** is created, with the heading written as `## Name` (the first letter in upper case when `name` is all lower case, otherwise as given): before the `## Log` section when the body has one, otherwise at the end of the body.
- **`updated`** (work items) is set to today when the file changes. Nothing is written, and `updated` is not touched, when `changed` is false. An ADR has no `updated`; its `date` is not changed.
- **Errors.** `not_found` if the id does not resolve. `invalid_request` for an empty or multi-line `name`, a `name` that starts with `#`, and for `replace` on the Log section, which is append-only (use `mode: "append"`). `stale` as for `set`. `refused` with `details.reason: "closed"` for a record in the done folder (`docs/work/done/` by default): a closed record is a record of what happened, and its text is not edited (`close` with `outcome` is how its Outcome is written).
- Side effect: rewrites one file, unless `dry_run`. CLI: `akmctl section <id> <name> --text <text>`; `--text -` reads the text from standard input.

### 4.17 `init`

Create an AKM repository in a directory: the layout, `AGENTS.md`, `CLAUDE.md`, a `README.md` stub, the skills and the manifest, and, when the directory is inside a git work tree, a `.gitignore` that keeps the lock file out of commits. It is for a new repository; adopting an existing one is `migrate`.

| Parameter     | Type    | Default | Meaning                                                                                                                                                                                                                                                            |
| ------------- | ------- | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `docs_dir`    | string  | `docs`  | The folder, one path segment, that holds the work items (`<docs_dir>/work`), the ADRs (`<docs_dir>/adr`), the roadmap (`<docs_dir>/ROADMAP.md`) and the manifest. Not empty, `.` or `..`, no path separator, no leading `.`, no character unsafe in a folder name. |
| `wikilinks`   | boolean | `false` | Write `reference_style: wikilink` to the manifest.                                                                                                                                                                                                                 |
| `dir`         | string  | `.`     | The directory to initialise, relative to the base directory, without `..`. It is created when absent.                                                                                                                                                              |
| `git_commits` | boolean | `false` | Write `git_commits: true` and install the optional skill `akm-commit` (section 4.18). Only effective inside a git work tree (below).                                                                                                                               |
| `git_init`    | boolean | `false` | Run `git init` in the directory when it is not inside a git work tree. A failure is reported in the response, not as an error.                                                                                                                                     |

Response: `docs_dir` and `wikilinks` (the values in effect), `dir`, `layout` (the version written to the manifest), `written` and `skipped` (repository-relative paths), `git_commits` (boolean; the value in effect), `git_commits_forced_off` (boolean; `git_commits` was requested but the directory is not inside a git work tree after the optional `git init`, so the key and the skill were not written), `git` (`none`, `initialised`, `existing` when the directory was already in a work tree, or `failed`) and, when `git` is `failed`, `git_error` (string). An invalid `docs_dir` is `invalid_request` with `details.parameter` `docs_dir`. A name other than `docs` is written to the manifest as `work_dir`, `adr_dir` and `roadmap`; the manifest, the deny list and (once a write runs) the lock live in that folder. A manifest that already exists is not changed by `init`; `config` changes it. An existing file is never overwritten: it is listed in `skipped`, so a second run writes nothing and succeeds. `.gitignore` is written only when the directory is inside a git work tree after the optional `git init`: it lists `<docs_dir>/.akm.lock`, and `.obsidian/workspace.json`, `.obsidian/workspaces.json` and `.trash/` only when an `.obsidian` folder already exists in the directory. `AGENTS.md` has no build, test, lint or run lines; its format-check line uses `<docs_dir>`. `refused` with `details.reason` `not_a_directory` when `dir` exists and is not a directory, and with `manifest_exists` when the directory already has a manifest in a folder other than `docs_dir` (a second one would make discovery ambiguous). The written repository passes `check` with no errors. Side effect: creates files. CLI: `akmctl init [--docs-dir <name>] [--wikilinks|--no-wikilinks] [--git-commits|--no-git-commits] [--git-init|--no-git-init] [--yes|--no-input] [dir]`; on a terminal the CLI asks for what the flags leave open, in this order, with plain prompts: where the records go (`docs_dir`), whether to use wikilinks, `Run git init? [Y/n]` (only when the directory is not in a git work tree) and whether to install the commit skill (only when a work tree exists or will be created). A non-terminal, `--yes` and `--no-input` never ask: the defaults apply, so a bare `akmctl init` succeeds without a terminal. The prompts are in the CLI only; the operation never asks. HTTP: `POST /init`.

### 4.18 `config`

Read or change a repository setting kept in the manifest. Without `key` it lists the settings; with `key` alone it returns that one; with `key` and `value` it sets it. The one setting is `git_commits`: a boolean that makes the tool install (`true`) or remove (`false`) the optional skill `.agents/skills/akm-commit/`, which tells an agent to commit after each logical step. The tool itself never commits.

| Parameter        | Type          | Default | Meaning                                                                                                 |
| ---------------- | ------------- | ------- | ------------------------------------------------------------------------------------------------------- |
| `key`            | `git_commits` | none    | The setting.                                                                                            |
| `value`          | boolean       | none    | The new value. Needs `key`.                                                                             |
| `dry_run`        | boolean       | false   | Return the plan and write nothing (no lock).                                                            |
| `confirm`        | boolean       | false   | Allow a removal.                                                                                        |
| `confirm_edited` | boolean       | false   | Allow a removal when installed files differ from the shipped copy (asked for in addition to `confirm`). |

Response: `settings`, a list of `{key, value, description, paths}` (`paths`: what the setting controls). A set adds `key`, `value`, `changed`, `applied` (false for a `dry_run`), `needs_confirmation` (a removal would happen), `plan` (`{write, remove, modify}`, repository-relative paths: files to create, files to delete, files to edit; `modify` is the manifest path, `["docs/.layout"]` by default, when the key changes) and `edits` (see below); `settings` then shows the values after the change. Setting the value it already has changes nothing (`changed` false, empty plan), and `off` with no skill directory only changes the key.

Rules, all of which hold before anything is written:

- **Only its own.** The tool writes or removes only inside `.agents/skills/<name>/` of a known optional skill (`akm-commit`; the name starts with `akm-`). `on` never overwrites: a path that exists and is not identical to the shipped file, or is a link or not a regular file or directory, makes the call `refused` with `details.reason: "conflict"` and `details.paths`, and nothing is changed. Links are never followed.
- **Removal is confirmed.** `off` with the directory present is `refused` with `details.reason: "needs_confirmation"` (with `details.plan`, `details.edited` and `details.edits`) unless `confirm` is true; a `dry_run` returns the plan instead. When an installed file differs from the shipped copy, or is not part of it, `edits` lists a diff per file (the shipped text to the installed text, in the form of section 7) and the call needs `confirm_edited: true` as well. Nothing outside the directory is touched.
- **A work tree.** `on` is `refused` with `details.reason: "not_a_git_repo"` when the repository root is not inside a git work tree. `off` always works.
- **A manifest.** `refused` with `details.reason: "no_manifest"` when the manifest does not exist.
- **Order and atomicity.** The call holds the repository lock; every file is written by temporary file and rename. `on` writes the skill, then the key; `off` removes the directory, then sets the key to `false`. A crash between the two leaves a state rule S4 reports and a repeat of the call repairs.

Side effect: writes files unless `dry_run`. CLI: `akmctl config` lists; `akmctl config git-commits on|off [--yes] [--dry-run]` (the CLI shows the plan, asks `y/N` on a terminal and then sends `confirm`, and `confirm_edited` after showing the diff and asking again; `--yes` sends both). HTTP: `GET /config` (query `key`) and `POST /config` (JSON body).

### 4.19 `commits`

Reads the git history and checks it against the `commit` member (section 1; ADR 261030-tbkh). Read-only; it runs `git log`, `git diff-tree` and `git show` and writes nothing. A repository that is not a git work tree is `refused` with `details.reason: "not_git"`.

| Parameter | Type    | Meaning                                                                                                                                                                                                    |
| --------- | ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `range`   | string  | A git revision range such as `origin/main..HEAD`. Default: `HEAD`. A value that starts with `-` or holds anything but letters, digits and `_./~^@{}-` is `invalid_request`; so is a range git cannot read. |
| `limit`   | integer | At most this many commits, newest first, 1 to 500. Default 50.                                                                                                                                             |

Response: `checked` (the AKM step commits examined), `skipped` (other commits, which make no claim and are not judged), `clean` (true when no examined commit has a problem), and `commits`, an array, newest first, of `{sha, subject, step, item, ok, problems}` for the examined commits. `problems` is an array of `{rule, message}`, empty when `ok` is true. Merge commits are not examined.

A commit is _examined_ when its subject has the form `<id>(<word>): <text>` or its message has a trailer that starts with `AKM-`; any other commit is skipped. The rules, all errors:

| Rule | The commit                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| ---- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| CM1  | has the subject `<id>(<step>): <detail>` with `<step>` one of the nine steps of the `commit` member.                                                                                                                                                                                                                                                                                                                                                                                                     |
| CM2  | has exactly one `AKM-Step` and one `AKM-Item` trailer, and they equal the step and id of the subject.                                                                                                                                                                                                                                                                                                                                                                                                    |
| CM3  | changed what the step says, in the files of the item (a file in the work or ADR folder whose name starts with `<id>-`): `file` adds one; `close` and `drop` move it into `done/` with `status` `done` or `dropped`; `claim` leaves `status: active` and a `claimed_by`; `release` removes a `claimed_by`; `ready`, `start` and `block` change `status` to `ready`, `active` and `blocked` from another value; `unblock` changes it from `blocked`. A commit that changes no file of the item breaks CM3. |

The exit code of the CLI is 1 when `clean` is false, as for `check`. The step is the caller's claim: the rules catch a message that does not match what the commit did, not a commit that made the right change under another message.

### 4.20 `skills`

Reports, installs and updates the skills the tool manages in `.agents/skills/` (ADR 261030-h6t7, ADR 261030-wsqb): `akm`, `akm-record-decision` and `akm-track-work`, and `akm-commit` while `git_commits` is on. Each managed skill folder holds a marker, `.akm-skill.json`, with the tool version that wrote it and the SHA-256 of each file as written; `init` and `config` write the same markers.

| Parameter        | Type                          | Default  | Meaning                                                                                                          |
| ---------------- | ----------------------------- | -------- | ---------------------------------------------------------------------------------------------------------------- |
| `mode`           | `status`, `install`, `update` | `status` | `status` reads. `install` adds what is missing. `update` also replaces outdated skills and retires renamed ones. |
| `dry_run`        | boolean                       | false    | Return the plan and write nothing (no lock). Needs a write mode.                                                 |
| `confirm`        | boolean                       | false    | Do the plan. Without it a write mode returns the plan and writes nothing. Needs a write mode.                    |
| `confirm_edited` | boolean                       | false    | Also replace or remove a skill that was edited or has no marker. Needs a write mode.                             |

Each skill is in one state: `missing` (no folder), `current` (every file equals the shipped one and the marker is present), `outdated` (the files differ from the shipped ones but equal what the marker recorded, so only the tool's older version wrote them), `edited` (a file differs from both the shipped text and the marker, or the folder has no marker and differs), `unmarked` (every file present equals the shipped one but there is no marker) or `conflict` (a path in the way is a file or a link). A folder under a retired name (`akm-work-items`) is `renamed`; it is `unedited` only when its marker is valid and every file matches the marker.

Response: `mode`, `clean` (every managed skill is `current` and no retired folder is left) and `skills`, one `{name, path, state}` per skill with `marker`, `files` (`{path, state}`), `note` and `stale` where they apply; a retired folder adds `replaced_by` and `unedited`. A write mode adds `plan` (`{skill, action, path, needs_confirm_edited}` with `action` `write`, `replace`, `marker` or `remove`, and a `diff` of section 7 for a `replace`), `applied`, `needs_confirmation` (a plan exists and was not applied), `written`, `replaced`, `removed` and `skipped` (`{skill, reason}`: a skill left alone, such as an edited one without `confirm_edited`, or an outdated one in `install`).

Rules: the tool writes or removes only inside the folders of the skills it manages, never follows links, and takes the repository lock for a write that applies. An unmarked folder gets its missing files and a marker; the files it has are kept. The retired folder `akm-work-items` that has no marker (the 0.3.x `init` wrote none) counts as edited, so removing it needs `confirm_edited`. `notes` (strings, present only when non-empty) says that the optional skill `akm-commit` is not installed while `git_commits` is off, and how to turn it on. A write mode is `refused` with `details.reason` `no_manifest` when the manifest is missing. The tool never commits.

Side effect: writes files in a write mode with `confirm` and without `dry_run`. CLI: `akmctl skills [status|install|update] [--yes] [--dry-run]` (on a terminal it shows the plan and asks, and asks again for edited skills; `--yes` sends both). HTTP: `GET /skills` (status only) and `POST /skills` (JSON body).

## 5. CLI binding

- The command is `akmctl <operation> [positional] [flags]`. Operation names are the subcommands of the specification (`board`, `check`, `claim`, `close`, `commits`, `config`, `find`, `index`, `init`, `ls`, `migrate`, `new`, `release`, `search`, `section`, `set`, `show`, `skills`, `version`, `when`). `ui` is not an API operation.
- Parameter `foo_bar` is the flag `--foo-bar`. A boolean flag with no value means true; `--no-foo-bar` means false. A list is a comma-separated value (`--status ready,active`) or a repeated flag. The required parameter that is a single id or text is positional (`akmctl show <id>`, `akmctl search <text>`, `akmctl find <id>`, `akmctl when <id>`). `new` takes `kind` and `slug` positionally (`akmctl new work add-json-output`). `set` takes `<id> field=value ...` and `close` takes `<id> done|dropped`; `claim` takes `<id> --as <name>`, and `section` takes `<id> <name>` with the text in `--text`.
- `--format json` prints the response object as one JSON document, then a newline, on stdout. The default format is a human text table (a `dry_run` response prints its `diff` in unified form) and is not part of the contract; it does not change with whether stdout is a terminal, so a command prints the same thing in a pipe and by hand. `AGENTS.md` names `--format json` explicitly. `--format ids` prints one `id` per line (for `ls`, `search`, `board`).
- With `--format json`, a failure also prints the response object (`ok: false`) on stdout, and nothing on stdout otherwise. Diagnostics and progress go to stderr and are never part of the contract.
- `--repo <dir>` sets the repository root; the default is the nearest ancestor of the working directory that has a manifest (`docs/.layout`, or the one `<name>/.layout` among its immediate non-hidden subdirectories) or `AGENTS.md`. None found is `unreadable`, and so are two candidate manifests, which the message names.
- `migrate` also takes `--migration <file>`: a JSON file holding one migration object or a list of them, sent as the `migrations` parameter (the flag can repeat). An unreadable or invalid file is exit 2.
- Exit codes: section 2.

## 6. HTTP binding

`akmctl ui` serves the binding on the loopback interface. It is for one local user and a local client (the workbench UI, an editor extension), not a network service.

**Security (all MUST):**

- Bind to `127.0.0.1` or `::1` only; there is no option to bind elsewhere in 1.0.
- Every request carries `Authorization: Bearer <token>`. The token is random (at least 128 bits), chosen at start and kept in memory only. `akmctl ui` prints it once on stderr inside the URL it opens, as a fragment (`http://127.0.0.1:PORT/#token=...`). A URL fragment is never sent to the server, so the page itself holds no token and no data: it is served without authentication, reads the token from the fragment, removes it from the address bar and sends it as the bearer header. A missing or wrong token is `unauthorized` (401), with no other detail.
- The server checks `Host` against its own loopback address and port, and rejects any request with an `Origin` header that is not the same origin (a defence against DNS rebinding and cross-site requests).
- No CORS headers are sent, and a preflight is not answered.
- The page and its two assets (`/`, `/app.js`, `/app.css`) are the only routes served without the token. They are sent with a `Content-Security-Policy` of `default-src 'none'` plus `'self'` for script, style and fetch, so a record's text cannot run as script; the page shows record text with `textContent` only. Every response carries `Cache-Control: no-store` and `X-Content-Type-Options: nosniff`.
- A request with a wrong `Host` or a foreign `Origin` is answered 403 with the error code `unauthorized`.
- A write route needs `Content-Type: application/json`.
- A write to an existing record (`set`, `section`, `close`, `claim`, `release`) MUST carry `If-Match` (or `if_match` in the body). Without it the server answers `refused` with `details.reason: "if_match_required"` and status `428`. A client that has no etag reads the record first (the response carries it). This is the compare-and-swap of section 7, and it is required here because the UI holds records for a long time. A request with `dry_run: true` is exempt: it writes nothing, so it needs no `If-Match`; when it carries one, the server checks it and answers `stale` if it differs.

**Routes** (base `/api/v1`; the version is in the path; the response `api_version` is `1.0`):

| Method and path                     | Operation          | Parameters                                                                                                                      |
| ----------------------------------- | ------------------ | ------------------------------------------------------------------------------------------------------------------------------- |
| `GET /version`                      | `version`          |                                                                                                                                 |
| `GET /records`                      | `ls`               | query: the `ls` parameters; lists comma-separated (`status=ready,active`); `sort=roadmap,-priority`                             |
| `GET /search`                       | `search`           | query: `text`, `body`, `type`, `status`, `fields`, `limit`, `offset`                                                            |
| `GET /records/{id}`                 | `show`             | query: `section`                                                                                                                |
| `GET /board`                        | `board`            | query: `statuses`                                                                                                               |
| `GET /commits`                      | `commits`          | query: `range`, `limit`                                                                                                         |
| `GET /check`                        | `check`            | query: `target`, `rules`, `as_of`                                                                                               |
| `GET /`                             | the workbench page | no token                                                                                                                        |
| `GET /find/{id}`                    | `find`             |                                                                                                                                 |
| `GET /find?name=`                   | `find`             | query: `name` (a note name; the path form is `GET /find/{id}`)                                                                  |
| `GET /when/{id}`                    | `when`             | query: `as_of`                                                                                                                  |
| `POST /records`                     | `new`              | JSON body: the request object                                                                                                   |
| `PATCH /records/{id}`               | `set`              | JSON body: `changes` (and `if_match`, `dry_run`)                                                                                |
| `PUT /records/{id}/sections/{name}` | `section`          | JSON body: `text`, `mode`, `if_match`, `dry_run`, `as_of`; `name` is the last path segment (a `name` in the body must equal it) |
| `POST /records/{id}/close`          | `close`            | JSON body: `status`, `outcome`, `if_match`, `dry_run`                                                                           |
| `POST /records/{id}/claim`          | `claim`            | JSON body: `as`, `if_match`, `takeover`, `dry_run`                                                                              |
| `DELETE /records/{id}/claim`        | `release`          | `If-Match` header; optional JSON body: `dry_run`                                                                                |
| `POST /index`                       | `index`            | JSON body: `dry_run`, `include_content`                                                                                         |
| `POST /migrate`                     | `migrate`          | JSON body: the request object                                                                                                   |
| `POST /init`                        | `init`             | JSON body: the request object                                                                                                   |
| `GET /config`                       | `config`           | query: `key`                                                                                                                    |
| `POST /config`                      | `config`           | JSON body: `key`, `value`, `dry_run`, `confirm`, `confirm_edited`                                                               |
| `GET /skills`                       | `skills`           | query: none (status only)                                                                                                       |
| `POST /skills`                      | `skills`           | JSON body: `mode`, `dry_run`, `confirm`, `confirm_edited`                                                                       |

- The response body is the response object, with the HTTP status of section 2 (`200` for `ok: true`, `201` for a successful `new`). The body is the same JSON the CLI prints. `check` and `migrate` return `200` even when `clean` is `false`.
- A record response carries `ETag: "<etag>"`. A write route reads `If-Match: "<etag>"` as `if_match` when the body has none.
- `GET` responses are never cached (`Cache-Control: no-store`).
- A change to a record file made outside the server (an editor, `git pull`) is seen by the next request; the server keeps no cache that can go stale, or revalidates it against file modification times and sizes. A client that wants push updates polls `GET /records?since=...` or uses `etag`; there is no event stream in 1.0.

## 7. Writes and concurrency

- **Atomicity.** A write goes to a temporary file in the same directory and is renamed over the target. The temporary file is flushed to storage before the rename, so a power loss as well as a crash leaves the old file or the new file, never a mix or a short file; the folder is flushed after the rename, as far as the filesystem allows. A move (`close`) renames, then rewrites; if the rewrite fails the file is back where it was.
- **Edit by splice.** Writes change only the bytes of the fields and lines they must (the parser edits by byte-splice), so comments, ordering and the body of a record are kept. A write MUST NOT reformat what it did not change.
- **Compare and swap.** `if_match` is the record's `etag`. When present and different from the file's current etag, the write is `stale` and nothing is written. The check and the rename happen while the tool holds an exclusive advisory lock (the file `.akm.lock` in the manifest folder, `docs/.akm.lock` by default, which a template lists in `.gitignore`), so two tool processes cannot both pass the check; an editor that writes the file between the final check and the rename is not covered, and that window is the length of one rename. In the CLI binding `if_match` is optional (`--if-match <etag>`): without it the write applies to the current file. In the HTTP binding it is required for writes to an existing record (section 6).
- **Locking.** The lock above is held only for the duration of one write. Claims are advisory (specification section 3).
- **Dry run.** A write operation (`set`, `section`, `close`, `claim`, `release`) with `dry_run` true does everything up to the write and then stops: it reads the file, applies the `if_match` check and every refusal rule, and builds the new text. It writes no file, renames nothing, and takes no lock (the lock file is not created or touched), so it never waits on a writer and no file's modification time changes. A dry run reads whole files only: a concurrent write is seen entirely or not at all, because writes replace by rename. The response is the response the write would give, with these differences: `record` is the would-be record (its fields, and for a `close` its `path`, as after the write) but its `etag` is the current etag of the file as it is now, not the etag the write would produce, so a client can use it as `if_match` for the real write; and the response has a `diff`.
- **The `diff` member** (`dry_run` responses only): `{path, new_path, hunks}`. `path` is the file now and `new_path` the file after the write; they differ for a `close`. `hunks` is an array, empty when nothing would change, of `{header, lines}`: `header` is the unified-diff header (`@@ -12,7 +12,9 @@`, a count of 1 is written without `,1`, an empty side is `0,0`) and `lines` is an array of `{op, text, old, new}` where `op` is `ctx`, `add` or `rem`; `text` is the line without its LF (a CR before the LF stays in `text`); `old` is the 1-based line number in the file now (`null` for `add`) and `new` the number after the write (`null` for `rem`). A hunk has up to three `ctx` lines before its first change and after its last; two changes with six or fewer unchanged lines between them are one hunk. A file that gains or loses its final newline shows its last line as removed and added. Applying the hunks to the file now gives exactly the file the write would produce. Where a run of lines could stand at more than one position because the lines around it repeat, which position the diff chooses is not defined; a client shows the diff and does not interpret it further.
- **Git.** No write operation stages, commits, pushes or branches, except `migrate` `apply_safe` (tag and branch). The caller decides when to commit.
- **After a write**, the response `record` is read back from the file, so `etag` is the new one.

## 8. Versions

- `api_version` is `MAJOR.MINOR`. A MINOR change adds an operation, an optional request parameter, a response field, an error `details` key or an enum value in a response; it never changes the meaning of what exists. A MAJOR change removes or renames something, adds a required parameter, or changes a type or meaning.
- The HTTP path carries the MAJOR (`/api/v1`). A tool that serves two MAJOR versions serves both paths. The CLI default is the tool's newest version; `--api-version 1` pins a MAJOR.
- A client MUST ignore response fields it does not know and MUST tolerate enum values it does not know in response fields (`class`, `state`).
- The API version is separate from the layout version and the extension versions (specification, section 7). `version` reports all of them.
- The schemas carry `$id` URIs of the form `https://akm.invalid/api/1.0/<name>.json` as stable identifiers only. They are not served; `.invalid` is a reserved top-level domain, so no real host is implied.

## 9. Conformance

An implementation conforms to API 1.0 if, for every case in the fixtures it claims, its response equals the expected response (after the normalisation below), and every response it produces validates against the schemas.

- **Layout of the fixtures.** `akm-api/.fixtures/` holds `repos/<name>/` (`basic`, a clean repository; `broken`, the same with three errors; `claims`, with the claims extension on and one claimed item; `closable`, with an item that has an Outcome; `vault`, the clean repository with its folders moved by the manifest under a folder whose names have spaces; `legacy`, a repository with no layout history and notes without ids; `commits`, the clean repository with `git_commits: true` and no skill; `conflict`, the clean repository with a file where the skill directory would go), `cases/<operation>-<n>.json` (67 cases: reads, `check`, `index`, `migrate`, `version` and the writes, including dry runs) and `make_fixtures.py`, which regenerates the sample repositories and the expected responses; it is a reference for the sample data only, not an implementation. Each case is `{name, repo, request, expected}`; a case with `git_work_tree: true` runs in a copy that is made a git work tree first (an empty `.git` directory). `repo` is the name of a directory in `repos/`; a runner that copies it to run a case names the copy the same, because `version` reports the name. `request` is the request object; the runner passes it through the binding under test. `expected` is the full response.
- **Normalisation.** Compared as JSON values, so key order and whitespace do not matter. Arrays are compared in order. `api_version` is compared; `request` is the request as given.
- **`version`.** Only the envelope and `repository` are compared. `tool`, `api_versions`, `layout_versions` and `extensions` describe the implementation, which the fixtures cannot know.
- **Clock.** Cases that depend on time set `as_of`.
- **Writes.** A write case runs against a copy of its repository. Besides `expected` it may carry `after` (files that must exist afterwards, by path and exact content) and `removed` (paths that must be gone). A case with neither must leave every file as it was (refusals, stale writes and dry runs). A dry-run case carries `would_be` (the path in the response's `diff.new_path`, mapped to the exact text the write would produce) instead of `after`. The runner compares the response without its `diff` member, then checks the `diff`: `path` and `new_path` equal the expected ones, the hunks are well formed (line numbers, headers, context and removed lines match the file as it is), and applying them to the file gives the `would_be` text. The expected `diff` in the fixture is one valid answer; equal diffs may place a run of lines differently (section 7). `new` has no exact case because ids are random: an implementation tests it for shape (the id pattern, the path, the fields).
- **YAML-lite.** The frontmatter fixtures (`yaml-lite.json`, described in section 0 of the specification) are the parser's conformance set and are run on the parser directly, not through this API.
- **Schemas.** `akm-api/schemas/` holds `common.json` (shared definitions) and one file per operation, in JSON Schema 2020-12. `akm-api/validate.js` validates every case's request and response against them (`npm install ajv@8`, then `node akm-api/validate.js akm-api`; the deliberate `invalid_request` request is skipped). This checks the schemas and the fixtures, not an implementation.
- **Dot-folder.** The fixtures are in a folder whose name starts with a dot, which Obsidian hides, so a vault that holds this document does not index the sample records as notes. In an implementation's repository the folder can be ordinary.

## 10. Not in 1.0

- **Fixtures for `commits`.** The operation reads git history, which a fixture directory cannot hold (a nested `.git` is not committed). An implementation tests it against repositories it builds with `git`: the tool's own messages are clean, and each rule is broken once.
- **Exact fixtures for `new`.** Ids are random; `new` is tested for shape only.
- **More write fixtures.** The set has `set`, `section`, `close`, `claim` and `release` cases; a case for each refusal reason and for the wikilink reference style can be added without changing the contract.
- **`search` ranking.** Ordered by match count, then id; a real index may want field weights.
- **Pagination by cursor.** `limit` and `offset` are enough for a repository of hundreds of records.
- **Event stream.** None; the UI polls. A MINOR version can add one if polling is a problem.
- **`renumber`.** Not in 1.0. When added it is a write operation with `dry_run` defaulting to true.
- **Extension operations beyond `claims`.** The `safety` extension adds rules and findings, no new operation.
- **Error `details` keys.** `paths` (`ambiguous`), `reason` (`refused`), `rule` (`refused` with reason `rule`) and `etag` (`stale`: the file's current etag, so a client can retry or show the conflict) are defined; more are a MINOR addition.
- **Lock file.** `.akm.lock` is in the specification's layout, in the manifest folder; a template lists it in `.gitignore`.
- **`mint_ids` migration step.** Implemented after 0.1 (work item 261028-8dd4); see the step primitives in 4.15. No migration embeds it yet. The options `rename`, `seed` and the steps `title_from_h1`, `goal_scaffold`, `remove_field`, `copy_field`, the `dir` of the field steps and the `migrations` parameter came with the vault migration (work item 261029-9aa2).
