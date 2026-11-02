---
type: design
status: draft
created_at: 2026-11-01
updated_at: 2026-11-01
authored_by: agent
---

# The first Codex trial: what the rails held and what they did not

Findings for 261101-vs5c. The maintainer gave Codex a loose prompt on 2026-11-01 ("find a nice work item that would be a good test… follow a PR flow instead of committing directly to main"). It produced PR 1. The code is secondary; this document is about how well the repository's instructions steered a model that was not the one they were written with.

## Method

1. Read the PR and its six commits against `AGENTS.md` and the `.agents/skills/` procedures.
2. Interview the same Codex session, read-only, with questions that ask for quoted rule text and what it saw, not for motives.
3. Check the interview against the session's rollout log (`~/.codex/sessions/`): which files it opened, in what order, and which commands it ran.

Step 3 matters because a model's account of its own reasoning is a reconstruction. The interview was accurate on every fact the log could confirm (reads, the sandbox failure that stopped the first claim, the quoted npm output); it was incomplete, not wrong.

## What worked

- Codex discovered `AGENTS.md` and `.agents/skills/` on its own, read the skills and the references it needed, and used `akmctl` for every record change. Nothing had to be pointed at it.
- `akmctl check` caught the WK6 violation (an active item with no `## Next`) at once. `akmctl commits --range` checked the step commits before the push.
- It staged by path, ran the full check list, never committed on `main`, and surfaced the tool failure as the rules ask.

## What went wrong, and why

| Deviation                                                                 | Cause in the evidence                                                                                                                                               | Layer that should hold it           |
| ------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------- |
| Implemented an item when asked to _find_ one                              | The prompt was loose and `AGENTS.md` has no rule for a request to find or recommend work                                                                            | Wording                             |
| Claimed w0mz, then switched to 3erc on one branch                         | No rule says a branch carries one item, or what to do after a claim commit exists                                                                                   | Wording                             |
| `claim` left 3erc `active` with no `## Next`, which `check` then rejected | The tool allows a state its own checker forbids                                                                                                                     | Tool (261101-p2m8)                  |
| Duplicate `## Log` lines, twice                                           | An unquoted `akmctl log <id> a; b` ran as two commands, zsh printed `command not found`, and the retry appended again. It repeated after seeing the same error once | Tool and wording (261101-jgdy)      |
| 2n55 filed with an empty `## Goal`                                        | `new` says to fill the file; nothing checks it                                                                                                                      | Check (261101-avg8)                 |
| 3erc closed with its pre-work `## Next`                                   | The rule says to update `## Next` in the code commit; nothing checks a closing commit                                                                               | Check (261101-avg8)                 |
| `go test` without `-race`                                                 | `AGENTS.md` says to add `-race`; the full-check line does not                                                                                                       | Wording (261101-7ghb)               |
| w0mz release folded into an unlabelled checkpoint commit                  | The checkpoint format has no trailers, so `commits` skips it                                                                                                        | Check (261101-avg8)                 |
| No `CHANGELOG.md` entry for a visible change (the block-detail output)    | `AGENTS.md` never says to write one; the convention lived only in the existing entries                                                                              | Wording (this commit adds the rule) |
| No agent attribution on commits                                           | `AGENTS.md` names no convention                                                                                                                                     | Wording (261101-7ghb)               |

Not a deviation but a portability fact: the sandbox made the Git metadata read-only and blocked the Go build cache and the network, so branch, test and PR each needed an escalation. `AGENTS.md` assumes all three are writable.

## Principle for the fixes

A rule that needs judgment (which item, how big a PR) gets one concrete sentence in `AGENTS.md`. A rule that can be computed (a state the checker already forbids, an empty Goal, a stale `## Next`) goes into `akmctl`, so no model has to remember it. Prose is the weakest layer; the trial's only reliable catches were the tool's.

## Next

Rerun the same task on a clean checkout with a reworded prompt, on at least Codex and Claude Code, and score both with the table above (261101-0vh7).
