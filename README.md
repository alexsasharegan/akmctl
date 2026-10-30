# Agent Knowledge Management Kit (`akmctl`)

Give your project a memory that you and your AI agents share: decisions, plans and tasks kept as plain Markdown files in the repository.

**AKM** stands for **agent knowledge management**. The kit is a standard for how those files are laid out, the `akmctl` command-line tool that creates and checks them, skills that teach an agent the routine, and an optional board in your browser. There is no database and no service: it is a folder of text files.

> **Status: unstable, pre-1.0 (0.3.x), for new (greenfield) projects.** You start a fresh folder with `akmctl init`. Bringing an existing repository into AKM (`akmctl migrate`) is documented but not supported yet. Expect changes before 1.0.0: any release below it may change the layout or the tool. macOS and Linux are supported; Windows is best effort.

## What you get

Run one command in an empty folder and you get this:

```text
my-project/
├── AGENTS.md              instructions an agent reads first
├── CLAUDE.md              points Claude Code at AGENTS.md
├── .agents/skills/        skills: akm, akm-work-items, akm-record-decision (and akm-commit, optional)
└── docs/
    ├── .layout            the manifest: layout version and settings
    ├── ROADMAP.md         the order of work
    ├── adr/               decisions, with a generated index
    └── work/              tasks; finished ones move to work/done/
```

`init` also writes a stub `README.md` (only if you have none; fill in its TODO), and, for tools that format Markdown, `.prettierrc.json` and `.prettierignore`: Prettier must leave record frontmatter alone and skip the generated ADR index. Ignore them if you do not use Prettier.

Why bother:

- **Agents pick up where the last session stopped.** Each task records its goal, its next step and its outcome, so a new session does not start from zero.
- **Decisions are written down.** A choice that is hard to reverse gets a short record, so nobody (human or agent) has to rediscover why.
- **It is just files.** You can read, edit, search and version them with the tools you already use.

## Install

Pick one method.

### Go install

```sh
go install github.com/alexsasharegan/akmctl/cmd/akmctl@latest
```

Needs Go 1.27 or later. The binary lands in `$(go env GOBIN)`, or in `$(go env GOPATH)/bin` when `GOBIN` is not set; make sure that folder is on your `PATH`.

### Download a release

Download the archive for your system from the [Releases](https://github.com/alexsasharegan/akmctl/releases) page (`akmctl_<version>_<os>_<arch>.tar.gz`, or `.zip` for Windows), unpack it and put `akmctl` on your `PATH`. The binaries are not code-signed, so macOS may block one that a browser downloaded; fetching it with `curl` avoids that, or use Go install.

### Build from source

```sh
git clone https://github.com/alexsasharegan/akmctl && cd akmctl
go build -o akmctl ./cmd/akmctl
```

Then move `akmctl` onto your `PATH`.

### Check it worked

```sh
akmctl version
```

## Start a project

Make a folder (or use one you already have), go into it and run:

```sh
mkdir my-project && cd my-project
akmctl init
```

That is the whole setup with the defaults: the records go in a `docs/` folder, links between records are plain ids, and nothing is committed for you. Add a first task and look at it:

```sh
akmctl new work first-task --title "My first task"
akmctl ls
akmctl ui        # optional: a board of your tasks in the browser (prints a local address)
```

Then point your agent at the folder; it reads `AGENTS.md` first and uses the skills in `.agents/skills/`. Add your project's own build and test commands to `AGENTS.md` yourself if you want agents to run them.

### Options

| Option              | What it does                                                                                              | Default   |
| ------------------- | --------------------------------------------------------------------------------------------------------- | --------- |
| `--docs-dir <name>` | The folder for the records: tasks, decisions, the roadmap, and the manifest. One folder name, no slashes. | `docs`    |
| `--wikilinks`       | Write links between records as `[[wikilinks]]`, which Obsidian shows in its links and graph.              | plain ids |
| `--git-init`        | Run `git init` if the folder is not yet a git repository.                                                 | off       |
| `--git-commits`     | Install the `akm-commit` skill so agents commit after each step. Needs a git repository.                  | off       |
| `--yes`             | Ask nothing; use the options you gave and the defaults for the rest.                                      |           |

Run bare `akmctl init` in a terminal and it asks about these in turn: where the records go, wikilinks, whether to run `git init` (only when needed), and the commit skill. Answer with `y` or `n`; the capital letter is the default if you just press Enter. When there is no terminal, for example when an agent runs it, it never asks. A `.gitignore` is written only in a git repository, and it adds Obsidian's per-machine files only if the folder already has an `.obsidian` folder.

`init` never overwrites a file, so running it again changes nothing.

### Examples

An Obsidian vault is any folder; Obsidian's own `.obsidian` folder is what makes it one. To keep the records in a folder called `Meta` with wikilinks:

```sh
cd my-vault
akmctl init --docs-dir Meta --wikilinks
```

The tool finds `Meta` on its own the next time you run it from anywhere inside the vault. In Obsidian, **Open folder as vault** on `my-vault` if you have not already.

A code project that uses git and wants agents to commit as they go:

```sh
cd my-project
akmctl init --git-init --git-commits
```

### Settings you can change later

```sh
akmctl config                      # list the settings
akmctl config git-commits on       # install the akm-commit skill (needs a git repository)
akmctl config git-commits off      # remove it again; it shows what it will remove and asks first
```

With `git-commits` on, agents commit after each logical step (a status change, or a "next step" checkpoint). The tool itself never stages or commits anything; only the skill tells the agent to.

## A task from start to finish

This is a real run in a fresh `code` project. The tool prints one line per change; nothing is staged or committed for you.

```sh
akmctl new work add-login --title "Add a login page"
# docs/work/261029-6p3c-add-login.md        a new file, status "raw"

akmctl set 261029-6p3c summary="A login page with a form" kind=feature priority=2 status=ready
akmctl section 261029-6p3c Goal --text - <<'EOF'
Users can sign in.

Done when:

- [ ] the form posts and shows errors
EOF
# the task is now "ready": it has a goal and a way to know when it is done

akmctl claim 261029-6p3c --as alex
# status: active, claimed_by: alex. Other people and agents see that it is taken

akmctl section 261029-6p3c Next --text "Build the form."
akmctl section 261029-6p3c Log --mode append --text "- 2026-10-29: form built"
# keep Next current and add dated Log lines as you work

akmctl section 261029-6p3c Outcome --text "Form works; see the tests."
akmctl close 261029-6p3c done
# moved to docs/work/done/261029-6p3c-add-login.md

akmctl check
# 0 errors, 0 warnings
```

The file's status goes `raw`, `ready`, `active`, `done`. `close` also moves it into `docs/work/done/` and refuses if the Outcome is empty. In practice your agent runs these steps for you using the `akm-work-items` skill; you read the result.

## What the unstable releases do not do

- It does not support adding AKM to an existing repository (`migrate` is there but is not a supported path before 1.0.0).
- It never stages, commits or pushes; the optional skill does the committing, only when you turn it on.
- Windows is best effort: the tool builds and the tests run there, but a Windows-only problem does not block a release before 1.0.0.

---

## For agents

Read the repository's `AGENTS.md` first; it wins over everything below. The skills in `.agents/skills/` are the procedures:

- `akm`: a routing table of every `akmctl` command with a reference file for each; read the reference before the first use of a command.
- `akm-work-items`: find, claim, work, checkpoint and close a work item.
- `akm-record-decision`: write an ADR for a hard-to-reverse choice.
- `akm-commit` (only when `git_commits` is `true` in `docs/.layout` and the folder is a git repository): commit after each logical step.

If `akmctl` is not installed, tell the maintainer and use the plain-file fallbacks in `AGENTS.md`; do not invent ids or hand-edit the ADR index.

Everyday commands: `akmctl ls --status ready,active,blocked`, `claim`, `section`, `set`, `close`, `check`. Add `--format json` for the full response. The operations, their JSON and the error codes are in the [API spec](docs/spec/akm-api-spec.md), with schemas and fixtures in [`docs/spec/akm-api/`](docs/spec/akm-api/). The standard (layout, record formats, rules) is the [specification](docs/spec/agent-ready-repo-spec.md); a plainer introduction is the [guide](docs/spec/akm-guide.md).

## Contributing to this repository

This repository keeps its own memory the AKM way and is the first test case of the standard: `AGENTS.md` (start here), `docs/adr/` (decisions), `docs/work/` (tasks), `docs/ROADMAP.md` and `docs/architecture.md`. `akmctl check` must pass on it.

```sh
go vet ./... && go test -race ./...
go build -o /tmp/akmctl ./cmd/akmctl && /tmp/akmctl check
npm ci && npm run format:check && npm run schemas   # Prettier and the API schema check
```

The tool uses only the standard library (the YAML library in tests is an oracle for the parser). The packages and invariants are in [`docs/architecture.md`](docs/architecture.md).

## Licence

MIT; see `LICENSE`. The Unicode data files in `internal/width/data` are under the Unicode licence.
