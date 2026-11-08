# Windows first run

A checklist for the first time you use `akmctl` on Windows. It takes about five minutes and uses a throwaway folder. If a step fails, stop and send back what [the last section](#if-a-step-fails) asks for: that is the most useful thing you can do. Windows support is new, so a failure is a finding, not a mistake of yours.

You need PowerShell 7 or Windows PowerShell 5.1. Git is optional for this checklist.

## 1. Install

Follow the [Windows steps in the README](../README.md#windows). The last line of that block prints the version. In a new terminal, check it again:

```powershell
akmctl version
```

It should print `akmctl <version> (API 1.0)`. If `akmctl` is not found, the `PATH` entry did not take: open a new terminal, and if that does not help, see "If a step fails".

## 2. The checklist

Run this block as it is. It makes the folder `akm-first-run` in your home folder, starts a repository there, files one work item, takes it through its life and checks the result. Each command prints one line or a short list; none of it should print an error.

<!-- first-run:start -->

```powershell
$dir = Join-Path $HOME "akm-first-run"
New-Item -ItemType Directory -Force $dir | Out-Null
Set-Location $dir
akmctl init --yes --no-input .
akmctl new work first-task
$id = akmctl ls --format ids | Select-Object -First 1
akmctl section $id Goal --text "Try akmctl on Windows."
akmctl section $id Next --text "Claim it, then close it."
akmctl set $id status=ready
akmctl claim $id --as first-run
akmctl section $id Outcome --text "It worked."
akmctl close $id done
akmctl check
akmctl ls --status done
```

<!-- first-run:end -->

What you should see at the end:

- `akmctl check` ends with `0 errors` (a warning is fine).
- `akmctl ls --status done` lists `first-task`.
- The file is in `docs\work\done\`, and the folder tree uses your normal Windows separators when you look at it in Explorer.

Then try the two things that differ on Windows:

1. Open the work item in your editor, change a word and save it. Run `akmctl set $id priority=2` while the file is still open in the editor. It should work; if it says `Access is denied`, that is the known fault 261108-mwtr: close the file in the editor and run it again, and tell us which editor you used.
2. Run `akmctl check` again after saving. If your editor changed the line endings to CRLF it should still say `0 errors`.

## If a step fails

Copy this into the message, with the failing command and its full output:

```powershell
"akmctl:     $(akmctl version)"
"PowerShell: $($PSVersionTable.PSVersion)"
"Windows:    $([Environment]::OSVersion.VersionString)"
"autocrlf:   $(git config core.autocrlf)"
"folder:     $PWD"
```

Also say: which step, what you typed, what you expected, and whether the folder is on OneDrive, a network drive or a path with spaces (all are supported; they are worth knowing about).

Known faults, so you do not have to report them again: 261108-mwtr (a write is refused while another program holds the file), 261108-ma45 (`check` reports A4 on a CRLF `docs\adr\README.md`), 261108-2z9m (new files are LF in a CRLF repository) and 261108-sfj3 (a link written with backslashes is not checked). The items are under 261108-2ncf.

## For the maintainer

CI runs the block between the `first-run` markers on the Windows runner on every push, so the commands here are known to work there. What CI cannot do is stand in for a person: a real machine, a real editor, OneDrive, antivirus and SmartScreen. Work item 261108-7avk stays open until someone has done the checklist by hand.
