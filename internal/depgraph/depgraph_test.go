package depgraph

import (
	"reflect"
	"testing"
)

func graph(edges ...[3]any) *Graph {
	g := &Graph{Out: map[string][]string{}, kinds: map[[2]string]Kind{}}
	for _, e := range edges {
		g.add(e[0].(string), e[1].(string), e[2].(Kind))
	}
	return g
}

func TestCyclesOnePerComponent(t *testing.T) {
	// a<->b mutual block; c -> d -> e -> c a longer loop; f hangs off both
	g := graph(
		[3]any{"b", "a", Blocked}, [3]any{"a", "b", Blocked},
		[3]any{"e", "c", Blocked}, [3]any{"c", "d", Blocked}, [3]any{"d", "e", Blocked},
		[3]any{"f", "a", Blocked}, [3]any{"f", "c", Blocked},
	)
	got := g.Cycles(true)
	want := [][]string{{"a", "b", "a"}, {"c", "d", "e", "c"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%v, want %v", got, want)
	}
}

func TestCyclesSelfAndAcyclic(t *testing.T) {
	if got := graph([3]any{"a", "a", Blocked}).Cycles(true); !reflect.DeepEqual(got, [][]string{{"a", "a"}}) {
		t.Fatalf("%v", got)
	}
	if got := graph([3]any{"a", "b", Blocked}, [3]any{"a", "c", Blocked}, [3]any{"b", "c", Blocked}).Cycles(true); len(got) != 0 {
		t.Fatalf("a diamond is not a cycle: %v", got)
	}
}

func TestParentOnlyLoopsAreWK4s(t *testing.T) {
	g := graph([3]any{"a", "b", Child}, [3]any{"b", "a", Child})
	if got := g.Cycles(true); len(got) != 0 {
		t.Fatalf("%v", got)
	}
	if got := g.Cycles(false); !reflect.DeepEqual(got, [][]string{{"a", "b", "a"}}) {
		t.Fatalf("%v", got)
	}
}

func TestParentWaitingOnItsChildIsACycle(t *testing.T) {
	// parent p waits for child c (Child edge); p is also blocked_by c, or c by p
	for _, g := range []*Graph{
		graph([3]any{"p", "c", Child}, [3]any{"c", "p", Blocked}),
		graph([3]any{"p", "c", Child}, [3]any{"p", "c", Blocked}, [3]any{"c", "g", Child}, [3]any{"g", "p", Blocked}),
	} {
		if got := g.Cycles(true); len(got) != 1 {
			t.Fatalf("%v", got)
		}
	}
}

func TestPathAndWith(t *testing.T) {
	g := graph([3]any{"a", "b", Blocked}, [3]any{"b", "c", Blocked})
	if got := g.Path("a", "c"); !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Fatalf("%v", got)
	}
	if g.Path("c", "a") != nil {
		t.Fatal("no path expected")
	}
	h := g.With("c", "a", Blocked)
	if got := h.Cycles(true); !reflect.DeepEqual(got, [][]string{{"a", "b", "c", "a"}}) {
		t.Fatalf("%v", got)
	}
	if len(g.Cycles(true)) != 0 {
		t.Fatal("With must not change the original")
	}
}
