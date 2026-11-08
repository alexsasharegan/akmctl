# Cold-agent check

How a maintainer tests what the skills tell an agent. The unit tests check that the two copies of a skill agree and that every operation has a reference; they cannot tell whether an agent that has never seen the project can follow the text. Only a cold run does (audit finding F6, work item 261108-79t5).

Run it after a change to a skill, to `AGENTS.md` or its template, or to a command a skill teaches, and before a release that touches any of them. A run takes about ten minutes.

## What "cold" means

The agent knows only what a new user's agent knows: the repository's `AGENTS.md`, the skills under `.agents/skills/` and the `akmctl` help. It does not see `docs/design/`, closed work items, the Log of the item it is working, or this procedure.

## Procedure

1. **Build the tool of the branch.** `go build -o /tmp/akmctl-cold ./cmd/akmctl`. Never test with an installed copy.
2. **Make a scratch copy.** Copy the repository without `.git` into a new empty folder (`cp -R` then delete `.git`, or `git archive HEAD | tar -x -C <dir>`), and delete `docs/design/` and `docs/work/done/`. The copy has no git, so the agent cannot commit; that is expected.
3. **Pick a task that uses the skill under test.** One sentence, in the words a user would say: "split this oversized item" (`akm-track-work`), "record that we chose X over Y" (`akm-record-decision`), "start work on the top ready item and close it". Use a real item of the copy.
4. **Start a fresh subagent** (a smaller model finds more snags than a larger one) with only this prompt:

   > You are working in `<scratch dir>`. Read `AGENTS.md` and follow it. Use `/tmp/akmctl-cold` for every akmctl command. The task: `<the sentence>`. Do not read anything under `docs/design/`. When you finish, list every point where you had to guess, every command that failed or surprised you, and every place the skill text and the tool disagreed.

5. **Read the result against the copy, not against the report.** Run `/tmp/akmctl-cold check` in the scratch copy and read the diff the agent left. An agent that says nothing went wrong can still have stored a wrong value.
6. **Sort what you found.** Each is one of:
   - a **tool snag** (a wrong result, a silent success that stored something odd, an error that does not say what to do): file a `raw` work item with the command, the actual and the expected behavior;
   - a **text snag** (the skill is missing a step or names a flag that does not exist): fix the skill in both copies;
   - a **judgement call the rule leaves to the agent**: note it in the Log and leave it.
7. **Pin every snag with a test** before closing the item that fixes it: a Go test for a tool snag, a line in the skill checks for a text snag. A snag without a test will come back.
8. **Record the run** in the Log of the item that asked for it: the date, the model, the task sentence, and the list from step 6. If nothing needed a guess, say so; that is the result a release wants.

## What has been found

| Run                                   | Snags                                                                                  | Pinned by                                                                                   |
| ------------------------------------- | -------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- |
| Split an oversized item (261101-t8h8) | `set blocked_by=[a,b]` stored a mangled list and `check` was silent (261108-h7ds)      | `TestListFieldValueForms`, `TestSetBlockedByStoresACleanList` in `internal/cli/cli_test.go` |
| Split an oversized item (261101-t8h8) | `set title=` left the first heading behind, so `check` warned WK12 later (261108-kdz1) | `TestSetTitleRenamesTheHeading` in `internal/api/retitle_test.go`                           |
| Split an oversized item (261101-t8h8) | `--changes` refuses an object (261107-1gha); `split` cannot give a child a title       | Open                                                                                        |

Add a row for each run that finds something.

## Limits

The check needs a model and a person to read the result, so CI does not run it. A scripted harness would have to fake the agent and would test the fake. What CI can do, and does, is pin each snag once it is found.
