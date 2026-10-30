---
id: "261029-2k7z"
title: spike for shell completions
summary: "add shell autocompletion command for bash, zsh, fish, and powershell"
status: raw
kind: spike
created: 2026-10-29
updated: 2026-10-29
---

# spike for shell completions

## Goal

### Problem / Context

Users currently have to manually type or copy-paste subcommands, flags, and dynamic arguments. Adding shell tab-completion reduces friction, aids command discovery, and prevents syntax errors.

### Proposed Solution

Expose a top-level `completion` command that emits autocompletion scripts for supported shells. If dynamic arguments exist (e.g., completing resource IDs, cluster names, or config profiles), wire them up using dynamic flag completion hooks.

- **Supported Shells:**
- Bash (`v4.0+`)
- Zsh
- Fish
- PowerShell

### User Experience

```bash
# Example usage in ~/.zshrc
source <(akmctl completion zsh)

# Or persisted to system completion directories
akmctl completion bash | sudo tee /etc/bash_completion.d/akmctl > /dev/null
```

### Scope & Technical Details

- **Completion Pattern (Self-Delegating Runtime Model):**
  - Because the Go standard library (`flag`) does not provide built-in shell completion generators, the CLI will adopt the delegating pattern: the CLI serves static shell wrapper scripts that route active tab-completion requests back into the binary for dynamic evaluation.
- **Subcommand & Hidden Flag Surface:**
  - **Public Command:** Add `akmctl completion [bash|zsh|fish]` (implemented using a sub-`flag.FlagSet`) to print the shell wrapper scripts.
  - **Hidden Hook:** Add an internal flag or hidden subcommand (e.g., `akmctl __complete <args...>` or inspection of `COMP_LINE` and `COMP_POINT` environment variables) to evaluate the current token stream and output completion candidates to `stdout`.
- **Flag & Command Introspection:**
  - Implement an internal completion registry or an inspector for registered `flag.FlagSet` instances.
  - For each sub-`FlagSet`, iterate through registered flags using `fs.VisitAll(func(f *flag.Flag) { ... })` to extract names and usage strings.
  - Match against the current word prefix (`COMP_WORDS[COMP_CWORD]` or equivalent) and output matching flag names (`--flag`) or subcommand names.
- **Script Generation via `//go:embed`:**
  - Store shell glue templates in an embedded filesystem (`embed.FS`):
    - `completions/bash.sh`
    - `completions/zsh.sh`
    - `completions/fish.fish`
  - Scripts will hook into the respective shell completion engines (`complete -C` for Bash, `compdef` for Zsh, `complete -c` for Fish) and forward the cursor context to `akmctl __complete`.
- **Performance & Early Exit Guardrails:**
  - The completion hook must execute before heavy runtime operations (e.g., reading remote configs, authenticating network clients, or loading heavy dependencies).
  - Check for the completion trigger at the very top of `main()` to guarantee terminal responses within $< 30\text{ ms}$.

Does the CLI use a custom subcommand dispatcher (e.g., switching on `os.Args[1]` with separate `flag.FlagSet` instances), or a single flat set of flags?

#### Acceptance Criteria

- [ ] Running `akmctl completion <shell>` generates a valid, working completion script for Bash, Zsh, Fish, and PowerShell.
- [ ] Passing an invalid or missing shell argument returns an actionable error and exits with code 1.
- [ ] Command descriptions/descriptions flags are displayed correctly in shells that support them (Zsh, Fish).
- [ ] Shell scripts do not emit stdout/stderr warnings when sourced or triggered.
- [ ] Documentation added to `README.md` or docs site covering installation instructions for all four shells across macOS and Linux.
- [ ] Unit tests verify that `completion [shell]` outputs valid text and exits with status `0`.

## Log

- 2026-10-29: filed.
