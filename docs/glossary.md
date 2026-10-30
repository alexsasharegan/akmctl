# Glossary

The words this project uses, and the ones it avoids. Every new feature, command and workbench part registers its terms here in the same change that adds it. When two words mean the same thing, this file picks one.

## Project and names

| Term                     | Meaning                                                                                                                                     |
| ------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------- |
| **AKM**                  | Agent knowledge management: the standard for keeping a repository's decisions, direction and work as Markdown records.                      |
| **`akmctl`**             | The command-line tool, and the name of this repository and its binary. It implements AKM.                                                   |
| **The kit**              | AKM, `akmctl`, the skills and the workbench together. Say "the standard" for AKM alone and "the tool" for `akmctl` alone.                   |
| **Specification** (spec) | `docs/spec/agent-ready-repo-spec.md`. The **guide** is the plain introduction; the **API spec** describes the operations.                   |
| **Repository**           | The folder (usually a git repository) that holds a project and its records. Say "repository", not "repo" or "project", in user-facing text. |
| **Docs dir**             | The folder inside the repository that holds the records, the roadmap and the manifest. `docs` by default (`init --docs-dir`).               |
| **Unstable**             | Any release below 1.0.0. The layout or the tool may change.                                                                                 |

## Files

| Term            | Meaning                                                                                                                                                               |
| --------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Record**      | Any Markdown file with an id and frontmatter in the AKM schema. The umbrella term. It has a **type**: `work` (a work item) or `adr`.                                  |
| **Document**    | A Markdown file that is not a record: `README.md`, `AGENTS.md`, `ROADMAP.md`, this glossary.                                                                          |
| **Note**        | Obsidian's word for any Markdown file. Every record is a note; most notes are not records. Avoid "note" for a record.                                                 |
| **Work item**   | A record of type `work`: one unit of in-flight work, in `work/` or, once closed, `work/done/`. **Task** is the informal word for it and is fine in user-facing prose. |
| **Kind**        | A work item's category: `feature`, `bug`, `chore` or `spike`. A **feature** is a work item of kind `feature`, not a product capability.                               |
| **ADR**         | Architecture decision record: a record of type `adr` in `adr/`, with a status. The **ADR index** (`adr/README.md`) is generated; never edit it by hand.               |
| **Id**          | `YYMMDD-xxxx`: the local day plus four random characters. A reference to a record is its id.                                                                          |
| **Frontmatter** | The YAML block at the top of a record. Its **fields** (status, kind, parent, cycle, closed, keep, and so on) are defined in the spec.                                 |
| **Section**     | A `##` part of a work item: `Goal`, `Next`, `Log`, `Outcome`. The `section` command reads and writes one.                                                             |
| **Status**      | A work item's state: `raw`, `ready`, `active`, `blocked`, `done`, `dropped`. `done` and `dropped` are the **closed** statuses.                                        |
| **Closed**      | A date field set by `close`; also the adjective for a `done` or `dropped` item.                                                                                       |
| **Cycle**       | An opaque label for the period a work item is planned into (a sprint, a program increment). The tool does not interpret it.                                           |
| **Keep**        | A flag that pins a work item so compaction (not yet built) never folds it.                                                                                            |
| **Claim**       | The `claimed_by` and `claimed_at` fields, set by `claim`, that say who holds an active item. `claimed_at` is the one UTC timestamp.                                   |
| **Parent**      | A work item's `parent`: the id of the item it belongs to. There is no separate epic record.                                                                           |
| **Roadmap**     | `ROADMAP.md`: the order of work in **Now / Next / Later** and **Dropped**, with **milestones** under Later (ADR 261029-6cp5).                                         |
| **Milestone**   | A `###` heading in the roadmap with exit criteria, such as 1.0.0.                                                                                                     |
| **Compaction**  | Folding old closed work items into a digest. A design in progress, not a feature.                                                                                     |

## Format and versions

| Term                | Meaning                                                                                                                                           |
| ------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Manifest**        | The file `docs/.layout` (in the docs dir): the layout version, the extension versions and the settings. Not "the layout file".                    |
| **Layout**          | The folder structure and rules a repository follows, identified by a version (`layout: 0.3.0`). Say "manifest" for the file.                      |
| **Extension**       | An optional rule set the manifest switches on: `adr-index`, `work-items`, `claims`, `safety`. Not a file extension.                               |
| **Rule**            | One check with an id such as `WK9`, a severity (error or warning) and a message. `akmctl check` reports **findings**.                             |
| **Skill**           | A procedure folder an agent reads, under `.agents/skills/`. Ours start with `akm-`; `akm` is the **router** skill with a reference per operation. |
| **Reference style** | How records link to each other: by `id` (plain) or `wikilink`.                                                                                    |

## Operations

| Term          | Meaning                                                                                                                         |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| **Operation** | One thing the API does (`new`, `set`, `close`, `check`, and so on). The CLI command with the same name calls it.                |
| **Envelope**  | The JSON response every operation returns, from the CLI with `--format json` or from HTTP.                                      |
| **Etag**      | A hash of a file's contents. A write sends it as `If-Match`; if the file changed, the answer is `stale` and nothing is written. |
| **Capture**   | Filing a raw work item quickly: `akmctl new`, or the workbench's capture modal.                                                 |
| **Migrate**   | Moving a repository from one layout to another by **steps**. Not supported for adopting an existing repository before 1.0.0.    |

## The workbench

Register every new view, control and message here when it is added. Names are as in the code (`internal/server/web/js/`).

| Term                    | Meaning                                                                                                                                               |
| ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Workbench**           | The browser app that `akmctl ui` serves: local, with a one-time token in the address. Older notes say "pkm-workbench"; it is the same thing.          |
| **Board**               | The workbench view that shows open work items in status lanes (Raw, Ready, Active, Blocked). `board` is also an API operation that returns the lanes. |
| **Lane**                | One status column on the board.                                                                                                                       |
| **Tables, ADRs, Check** | The other three workbench views.                                                                                                                      |
| **Drawer**              | The side panel that shows one record: its fields, what points at it, and its sections.                                                                |
| **Capture modal**       | The form (key `c`) that files a raw work item.                                                                                                        |
| **Move**                | Changing a work item's status from the board. A **safe move** writes at once and offers undo (`u`); any other goes through the wizard.                |
| **Move wizard**         | The one-question-per-step flow that collects what a move still needs, then shows the review.                                                          |
| **Review**              | The diff of every write a move will make, shown before any is made.                                                                                   |
| **Live board**          | The board kept current by the **poll**: a request every two seconds for each work item's id, path, etag and status.                                   |
| **Toast**               | A short message shown over the page. (Placement and kinds are being redesigned in 261029-kmre.)                                                       |
| **Palette**             | The search box (`⌘K`).                                                                                                                                |

## Obsidian

| Term               | Meaning                                                                                                                                                                   |
| ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Obsidian vault** | A folder Obsidian opens (it has an `.obsidian` folder). A repository can be, or live inside, one. Do not use "vault" for a repository.                                    |
| **Wikilink**       | `[[name]]` or `[[name\|alias]]`: Obsidian's link form, used when `reference_style` is `wikilink`.                                                                         |
| **Obsidian Bases** | Always written in full: **Obsidian Bases**, capitalised. Never "base" or "a base" alone, because the workbench code also has a variable called `base` (the current view). |

## Words to avoid

- "task" in a spec or an API doc (use "work item"); "repo" or "project" for the repository; "vault" for the repository; "layout file" for the manifest; "epic" (use `parent`); "wizard" alone (say "move wizard"); "beta" (say "unstable").
