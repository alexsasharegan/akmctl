package api

import (
	"strings"

	"github.com/alexsasharegan/akmctl/internal/record"
)

// The edit states of a work item: what a client may offer to change, by
// status. The matrix is the Outcome of work item 261031-z0n8; show serves it
// (work item 261101-3che) so that no client holds a copy.
//
// "yes" is editable, "warn" is editable and the client shows the reason first
// (the API never refuses on a warning), "no" is read-only. A closed item's
// title, summary and sections are "no" with correction true: a write that
// carries a correction is accepted (ADR 261101-dxcn).
const (
	editYes  = "yes"
	editWarn = "warn"
	editNo   = "no"
)

// editState is one cell of the matrix.
type editState struct {
	State      string
	Reason     string // set for warn and no
	Correction bool   // a closed item: a write with a correction is accepted
	Append     bool   // the Log: lines are added, never rewritten
}

func (s editState) obj() Obj {
	o := O("state", s.State)
	if s.Reason != "" {
		o = append(o, KV{"reason", s.Reason})
	}
	if s.Correction {
		o = append(o, KV{"correction", true})
	}
	if s.Append {
		o = append(o, KV{"append", true})
	}
	return o
}

const (
	reasonManaged      = "managed by the tool"
	reasonNotSettable  = "not a field set can write"
	reasonClosed       = "the item is closed"
	reasonCorrection   = "the item is closed; a correction with a one-line reason changes it"
	reasonUnderWay     = "the item is under way; this describes it as it was before work began"
	reasonRanked       = "the item is active; its priority ranked it before work began"
	reasonGoalposts    = "the item is active; changing the goal moves the goalposts, so add a log line that says why"
	reasonClosedNoNext = "a closed item has no next step"
	reasonLogClosed    = "the item is closed; its Log is written by the tool"
	reasonStatusClosed = "a closed item stays done or dropped"
)

// fieldEditState is the matrix row of one frontmatter field of a work item in
// the given status. closed is true for done and dropped.
func fieldEditState(field, status string, closed bool) editState {
	if !contains(setFields[record.Work], field) {
		if contains(record.TypeFields(record.Work, true), field) {
			return editState{State: editNo, Reason: reasonManaged}
		}
		return editState{State: editNo, Reason: reasonNotSettable}
	}
	if closed {
		switch field {
		case "keep":
			return editState{State: editYes}
		case "title", "summary":
			return editState{State: editNo, Reason: reasonCorrection, Correction: true}
		case "status":
			return editState{State: editNo, Reason: reasonStatusClosed}
		}
		return editState{State: editNo, Reason: reasonClosed}
	}
	switch field {
	case "status":
		return editState{State: editYes}
	case "kind", "parent":
		if status == "active" || status == "blocked" {
			return editState{State: editWarn, Reason: reasonUnderWay}
		}
	case "priority":
		if status == "active" {
			return editState{State: editWarn, Reason: reasonRanked}
		}
	}
	return editState{State: editYes}
}

// sectionEditState is the matrix row of one section (lower-case name).
func sectionEditState(name, status string, closed bool) editState {
	if name == "log" {
		if closed {
			return editState{State: editNo, Reason: reasonLogClosed}
		}
		return editState{State: editYes, Append: true}
	}
	if closed {
		if name == "next" {
			return editState{State: editNo, Reason: reasonClosedNoNext}
		}
		return editState{State: editNo, Reason: reasonCorrection, Correction: true}
	}
	if name == "goal" && status == "active" {
		return editState{State: editWarn, Reason: reasonGoalposts}
	}
	return editState{State: editYes}
}

// editSections are the sections a work item can always be given, listed in
// edit even when the file does not have them yet.
var editSections = []string{"goal", "next", "outcome", "log"}

// editObj is the edit member of show for a work item: one state for each
// field the response lists and for each section the file has or can be given.
func editObj(repo *record.Repo, r *record.Record, fields Obj, sections []string) Obj {
	status := r.Status()
	closed := repo.IsDone(r.Path) || status == "done" || status == "dropped"
	fo := Obj{}
	for _, kv := range fields {
		fo = append(fo, KV{kv.K, fieldEditState(kv.K, status, closed).obj()})
	}
	names := append([]string(nil), sections...)
	for _, n := range editSections {
		if !contains(names, n) {
			names = append(names, n)
		}
	}
	so := Obj{}
	for _, n := range names {
		n = strings.ToLower(n)
		so = append(so, KV{n, sectionEditState(n, status, closed).obj()})
	}
	return O("fields", fo, "sections", so)
}
