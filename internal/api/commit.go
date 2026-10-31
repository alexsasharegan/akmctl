package api

import "strings"

// The commit object is how a step's commit message is made one way (ADR
// 261030-tbkh). The tool returns it; it never stages or commits.

// commitSteps is the closed set of steps.
var commitSteps = map[string]bool{
	"file": true, "ready": true, "start": true, "block": true, "unblock": true,
	"claim": true, "release": true, "close": true, "drop": true,
}

// commitObj builds the `commit` member: step, item, detail, paths and the
// message (subject, blank line, trailers).
func commitObj(step, id, detail string, paths ...string) Obj {
	detail = strings.TrimSpace(strings.Join(strings.Fields(detail), " "))
	subject := id + "(" + step + "): " + detail
	msg := subject + "\n\nAKM-Step: " + step + "\nAKM-Item: " + id
	ps := make([]any, len(paths))
	for i, p := range paths {
		ps[i] = p
	}
	return O("step", step, "item", id, "detail", detail, "paths", ps, "message", msg)
}

// statusStep names the step for a work item's status change made by `set`, or
// "" when the change is not a step.
func statusStep(from, to string) string {
	switch {
	case from == to:
		return ""
	case to == "blocked":
		return "block"
	case from == "blocked":
		return "unblock"
	case to == "ready":
		return "ready"
	case to == "active":
		return "start"
	}
	return ""
}
