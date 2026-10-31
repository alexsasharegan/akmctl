// Package depgraph builds the dependency graph of work items (spec WK18): an edge
// from an item to each work item in its blocked_by (it waits for that item) and
// from a parent to each of its children (a parent waits for its children).
package depgraph

import (
	"sort"

	"github.com/alexsasharegan/akmctl/internal/frontmatter"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// Kind says why an item waits for another.
type Kind int

const (
	Blocked Kind = iota // the waiting item names the other in blocked_by
	Child               // the waiting item is the other's parent
)

// Graph maps an item id to the ids it waits for, each list sorted.
type Graph struct {
	Out   map[string][]string
	kinds map[[2]string]Kind
}

// Build reads the edges from the repository. Entries that do not resolve to a
// work item (free text, unknown ids, ADRs) are ignored.
func Build(repo *record.Repo) *Graph {
	g := &Graph{Out: map[string][]string{}, kinds: map[[2]string]Kind{}}
	work := func(ref string) *record.Record {
		id := record.RefID(ref)
		if !record.ValidID(id) {
			return nil
		}
		if recs := repo.ByID(id); len(recs) > 0 && recs[0].Type == record.Work {
			return recs[0]
		}
		return nil
	}
	for _, r := range repo.Records {
		if r.Type != record.Work {
			continue
		}
		for _, it := range r.Fields["blocked_by"].Items {
			if it.Kind == frontmatter.String {
				if t := work(it.Str); t != nil {
					g.add(r.ID, t.ID, Blocked)
				}
			}
		}
		if p := work(r.Str("parent")); p != nil {
			g.add(p.ID, r.ID, Child)
		}
	}
	for k := range g.Out {
		sort.Strings(g.Out[k])
	}
	return g
}

func (g *Graph) add(from, to string, k Kind) {
	e := [2]string{from, to}
	if _, ok := g.kinds[e]; ok {
		if k == Blocked { // an explicit block outranks the implicit edge
			g.kinds[e] = Blocked
		}
		return
	}
	g.kinds[e] = k
	g.Out[from] = append(g.Out[from], to)
}

// ParentLinks returns the graph reduced to the parent-to-child edges, whose loops
// are WK4.
func (g *Graph) ParentLinks() *Graph {
	h := &Graph{Out: map[string][]string{}, kinds: map[[2]string]Kind{}}
	for e, k := range g.kinds {
		if k == Child {
			h.add(e[0], e[1], Child)
		}
	}
	for k := range h.Out {
		sort.Strings(h.Out[k])
	}
	return h
}

// Without returns a copy of the graph minus one edge.
func (g *Graph) Without(from, to string) *Graph {
	h := &Graph{Out: map[string][]string{}, kinds: map[[2]string]Kind{}}
	for e, k := range g.kinds {
		if e != [2]string{from, to} {
			h.add(e[0], e[1], k)
		}
	}
	for k := range h.Out {
		sort.Strings(h.Out[k])
	}
	return h
}

// With returns a copy of the graph with one more edge (the proposed change).
func (g *Graph) With(from, to string, k Kind) *Graph {
	h := &Graph{Out: map[string][]string{}, kinds: map[[2]string]Kind{}}
	for f, ts := range g.Out {
		h.Out[f] = append([]string(nil), ts...)
	}
	for e, k := range g.kinds {
		h.kinds[e] = k
	}
	h.add(from, to, k)
	sort.Strings(h.Out[from])
	return h
}

// Path returns a path of ids from `from` to `to` (both included), or nil when
// `to` cannot be reached. A depth-first search in id order, so it is stable.
func (g *Graph) Path(from, to string) []string {
	seen := map[string]bool{}
	var walk func(cur string, path []string) []string
	walk = func(cur string, path []string) []string {
		path = append(path, cur)
		if cur == to && len(path) > 1 {
			return path
		}
		seen[cur] = true
		for _, n := range g.Out[cur] {
			if n == to {
				return append(path, n)
			}
			if !seen[n] {
				if p := walk(n, path); p != nil {
					return p
				}
			}
		}
		return nil
	}
	return walk(from, nil)
}

// Cycles returns one cycle for each strongly connected component that has a loop,
// as the ordered ids with the first repeated at the end (a -> b -> a), found by
// Tarjan's algorithm in O(items + edges). Cycles are sorted by their first id, and
// each starts at its lowest id. With onlyBlocked, a component whose loops are made
// only of parent links is left out (WK4 reports those, on ParentLinks).
func (g *Graph) Cycles(onlyBlocked bool) [][]string {
	ids := make([]string, 0, len(g.Out))
	for k := range g.Out {
		ids = append(ids, k)
	}
	sort.Strings(ids)
	index, low := map[string]int{}, map[string]int{}
	on := map[string]bool{}
	var stack []string
	var sccs [][]string
	n := 0
	var strong func(v string)
	strong = func(v string) {
		index[v], low[v] = n, n
		n++
		stack = append(stack, v)
		on[v] = true
		for _, w := range g.Out[v] {
			if _, ok := index[w]; !ok {
				strong(w)
				low[v] = min(low[v], low[w])
			} else if on[w] {
				low[v] = min(low[v], index[w])
			}
		}
		if low[v] == index[v] {
			var comp []string
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				on[w] = false
				comp = append(comp, w)
				if w == v {
					break
				}
			}
			sccs = append(sccs, comp)
		}
	}
	for _, v := range ids {
		if _, ok := index[v]; !ok {
			strong(v)
		}
	}
	var out [][]string
	for _, comp := range sccs {
		sort.Strings(comp)
		in := map[string]bool{}
		for _, c := range comp {
			in[c] = true
		}
		start := comp[0]
		cyc := g.cycleIn(start, in, onlyBlocked)
		if cyc != nil {
			out = append(out, cyc)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i][0] < out[j][0] })
	return out
}

// cycleIn finds a cycle through start inside one component. With onlyBlocked it
// prefers a cycle that uses a blocked_by edge, and returns nil when the component
// has a loop only of parent links.
func (g *Graph) cycleIn(start string, in map[string]bool, onlyBlocked bool) []string {
	// a component with one item is a loop only if it waits for itself
	if len(in) == 1 {
		for _, w := range g.Out[start] {
			if w == start {
				return []string{start, start}
			}
		}
		return nil
	}
	if !onlyBlocked {
		return g.loopThrough(start, in)
	}
	// try each blocked_by edge u -> v inside the component: a cycle exists through it
	var froms []string
	for u := range in {
		froms = append(froms, u)
	}
	sort.Strings(froms)
	for _, u := range froms {
		for _, v := range g.Out[u] {
			if in[v] && g.kinds[[2]string{u, v}] == Blocked {
				back := g.pathWithin(v, u, in)
				if back != nil {
					return rotate(append([]string{u}, back...))
				}
			}
		}
	}
	return nil
}

func (g *Graph) loopThrough(start string, in map[string]bool) []string {
	p := g.pathWithin(start, start, in)
	if p == nil {
		return nil
	}
	return rotate(append([]string{start}, p...))
}

// pathWithin is Path restricted to the ids in `in`; from == to finds a loop.
func (g *Graph) pathWithin(from, to string, in map[string]bool) []string {
	seen := map[string]bool{}
	var walk func(cur string) []string
	walk = func(cur string) []string {
		seen[cur] = true
		for _, n := range g.Out[cur] {
			if !in[n] {
				continue
			}
			if n == to {
				return []string{n}
			}
			if !seen[n] {
				if p := walk(n); p != nil {
					return append([]string{n}, p...)
				}
			}
		}
		return nil
	}
	if from == to {
		return walk(from)
	}
	if p := walk(from); p != nil {
		return append([]string{from}, p...)
	}
	return nil
}

// rotate takes a closed path (first == last) and starts it at its lowest id.
func rotate(p []string) []string {
	body := p[:len(p)-1]
	lo := 0
	for i, id := range body {
		if id < body[lo] {
			lo = i
		}
	}
	out := append(append([]string{}, body[lo:]...), body[:lo]...)
	return append(out, out[0])
}
