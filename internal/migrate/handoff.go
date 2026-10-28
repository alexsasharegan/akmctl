package migrate

import (
	"fmt"
	"strings"
)

// Handoff writes the prompt for an agent that finishes a migration by hand:
// what to read, what is left, where, when to stop, and the safeguards.
func Handoff(from, to string, results []Result, findings []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Finish the migration from layout %s to %s\n\n", from, to)
	b.WriteString("You are finishing a migration of this repository's AKM records. The tool has applied the steps it can do safely. The rest needs judgement.\n\n")
	b.WriteString("## Read first\n\n- `AGENTS.md` and the layout spec the repository points to (docs/ROADMAP.md, docs/adr/, docs/work/).\n- `docs/.layout`, the manifest, and `docs/.migrate-log` if it exists: the steps already applied.\n\n")
	b.WriteString("## What is left\n\n")
	n := 0
	for _, r := range results {
		if r.Step.Class == "auto" && r.Status != "pending" {
			continue
		}
		n++
		fmt.Fprintf(&b, "%d. **%s** (%s): %s\n", n, r.Step.ID, r.Step.Class, r.Step.Description)
		for _, f := range r.Files {
			fmt.Fprintf(&b, "   - %s\n", f)
		}
	}
	if n == 0 {
		b.WriteString("No steps are left. Run the check below.\n")
	}
	if len(findings) > 0 {
		b.WriteString("\n## Findings from the target version's check\n\n")
		for _, f := range findings {
			b.WriteString("- " + f + "\n")
		}
	}
	fmt.Fprintf(&b, "\n## Stop when\n\n`akmctl check --target %s` reports no errors and every step above is done.\n\n", to)
	b.WriteString("## Safeguards\n\n")
	b.WriteString("- Work on the migration branch; the backup tag is `pre-migrate-" + from + "-to-" + to + "`. Do not delete the tag.\n")
	b.WriteString("- Never delete a file. Stop and ask on any ambiguity.\n")
	b.WriteString("- Change only what a step names. Keep unknown fields and prose.\n")
	b.WriteString("- Show the diff to the maintainer before any commit. Do not push.\n")
	return b.String()
}
