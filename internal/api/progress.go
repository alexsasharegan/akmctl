package api

import (
	"sort"

	"github.com/alexsasharegan/akmctl/internal/record"
)

// childIndex maps a work item's id to its children, ordered by id. It is computed
// from the parent links on each call and never stored (H4: one source of truth).
func childIndex(repo *record.Repo) map[string][]*record.Record {
	idx := map[string][]*record.Record{}
	for _, r := range repo.Records {
		if r.Type != record.Work {
			continue
		}
		if p := record.RefID(r.Str("parent")); p != "" {
			idx[p] = append(idx[p], r)
		}
	}
	for _, l := range idx {
		sort.Slice(l, func(i, j int) bool { return l[i].ID < l[j].ID })
	}
	return idx
}

// progressOf counts the children of id: done, dropped and all, direct and then
// through every descendant. It is nil for an item with no children. A parent loop
// (which check reports) is walked once.
func progressOf(idx map[string][]*record.Record, id string) Obj {
	kids := idx[id]
	if len(kids) == 0 {
		return nil
	}
	var cd, cx int
	for _, k := range kids {
		switch k.Status() {
		case "done":
			cd++
		case "dropped":
			cx++
		}
	}
	var dd, dx, dt int
	seen := map[string]bool{id: true}
	var walk func(string)
	walk = func(p string) {
		for _, k := range idx[p] {
			if seen[k.ID] {
				continue
			}
			seen[k.ID] = true
			dt++
			switch k.Status() {
			case "done":
				dd++
			case "dropped":
				dx++
			}
			walk(k.ID)
		}
	}
	walk(id)
	return O("children_done", cd, "children_dropped", cx, "children_total", len(kids),
		"descendants_done", dd, "descendants_dropped", dx, "descendants_total", dt)
}

// treeOrder orders recs depth first: each record is followed by its children that are
// also in recs, in the order recs gives them. A record whose parent is not in recs is a
// root. It returns the records and the depth of each.
func treeOrder(recs []*record.Record) ([]*record.Record, map[string]int) {
	in := map[string]bool{}
	for _, r := range recs {
		in[r.ID] = true
	}
	kids := map[string][]*record.Record{}
	var roots []*record.Record
	for _, r := range recs {
		p := ""
		if r.Type == record.Work {
			p = record.RefID(r.Str("parent"))
		}
		if p != "" && in[p] && p != r.ID {
			kids[p] = append(kids[p], r)
		} else {
			roots = append(roots, r)
		}
	}
	depth := map[string]int{}
	var out []*record.Record
	var walk func(r *record.Record, d int)
	walk = func(r *record.Record, d int) {
		if _, dup := depth[r.ID]; dup {
			return
		}
		depth[r.ID] = d
		out = append(out, r)
		for _, k := range kids[r.ID] {
			walk(k, d+1)
		}
	}
	for _, r := range roots {
		walk(r, 0)
	}
	// records only in a parent loop have no root: list them after the rest
	for _, r := range recs {
		walk(r, 0)
	}
	return out, depth
}
