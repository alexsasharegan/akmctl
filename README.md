# Agent Knowledge Management Kit (`akmctl`)

Give your project a memory that you and your AI agents share: decisions, plans and tasks kept as plain Markdown files in the repository.

**AKM** stands for **agent knowledge management**. The kit is a standard for how those files are laid out, the `akmctl` command-line tool that creates and checks them, skills that teach an agent the routine, and an optional board in your browser. There is no database and no service: it is a folder of text files.

> **Status: 0.2.0 beta, for new (greenfield) projects.** You start a fresh folder with `akmctl init`. Bringing an existing repository into AKM (`akmctl migrate`) is documented but not supported yet. Expect changes before 1.0.0. macOS and Linux are supported; Windows is best effort.

## What you get

Run one command in an empty folder and you get this:

```text
my-project/
├── AGENTS.md              instructions an agent reads first (fill in your build and test commands)
├── CLAUDE.md              points Claude Code at AGENTS.md
├── .agents/skills/        skills: akm, akm-work-items, akm-record-decision (and akm-commit, optional)
└── docs/
    ├── .layout            the manifest: layout version and settings
    ├── ROADMAP.md         the order of work
    ├── adr/               decisions, with a generated index
    └── work/              tasks; finished ones move to work/done/
```

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

`akmctl init` needs a **profile**. Both profiles create the same kit: `AGENTS.md`, the skills and the `docs/.layout` manifest. They differ in where the records go and what is set up around them:

|                                     | `--profile code`                                       | `--profile notes`                                                         |
| ----------------------------------- | ------------------------------------------------------ | ------------------------------------------------------------------------- |
| For                                 | a software project (a repository with source code)     | a notes folder, a knowledge base or an Obsidian vault                     |
| Records (tasks, decisions, roadmap) | under `docs/`                                          | under `notes/`                                                            |
| Links between records               | plain ids (`261029-6p3c`)                              | Obsidian `[[wikilinks]]`                                                  |
| `AGENTS.md`                         | has Build, Test, Lint and Run lines for you to fill in | no build commands; it is about notes                                      |
| `.gitignore`                        | the tool's lock file                                   | also Obsidian's per-machine files (`.obsidian/workspace.json`, `.trash/`) |

Pick `code` when the folder is a project you build and test. Pick `notes` when the folder is mainly Markdown notes you read and link, whether or not you use Obsidian. You cannot switch later without moving the records by hand, so choose the one that matches the folder.

### If you use Obsidian (or just want a notes folder)

An Obsidian vault is any folder; Obsidian's own `.obsidian` folder is what makes it one. So:

```sh
mkdir my-notes && cd my-notes
akmctl init --profile notes
```

Then in Obsidian choose **Open folder as vault** and pick `my-notes`. The records live under `notes/` and link to each other with `[[wikilinks]]`, so they show up in Obsidian's links and graph. Add a first task:

```sh
akmctl new work first-task --title "My first task"
akmctl ui        # optional: a board of your tasks in the browser (prints a local address)
```

### If you have a code project

```sh
mkdir my-project && cd my-project     # or an existing folder you want to start fresh in
akmctl init --profile code
```

Open `AGENTS.md` and fill in the commands for building, testing and linting. Then point your agent at the folder; it reads `AGENTS.md` first and uses the skills in `.agents/skills/`.

### Questions `init` asks

Run bare `akmctl init` in a terminal and it asks a few plain questions: the profile, whether to install the optional git-commit skill, and whether to run `git init` if the folder is not yet a git repository. Answer with `y` or `n`; the capital letter is the default if you just press Enter. Flags answer them ahead of time (`--profile`, `--git-commits` or `--no-git-commits`, `--git-init` or `--no-git-init`), and `--yes` skips the questions. When there is no terminal, for example when an agent runs it, it never asks.

`init` never overwrites a file, so running it again changes nothing.

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

## What the beta does not do

- It does not support adding AKM to an existing repository (`migrate` is there but is not a supported path in this beta).
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
