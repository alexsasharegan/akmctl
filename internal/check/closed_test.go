package check

import "testing"

func TestClosedCycleKeepValid(t *testing.T) {
	res := Run(repoWith(t, map[string]func(string) string{
		a2b3: sub("priority: 2", "priority: 2\ncycle: 2026-PI3-S4\nkeep: true"),
		e5f6: sub("status: ready", "status: ready\ncycle: a.b_c-1\nkeep: false"),
		j9k0: sub("closed: 2026-10-12", "closed: 2026-10-10"),
	}), Options{AsOf: asOf})
	if len(res.Findings) != 0 {
		t.Fatalf("%+v", res.Findings)
	}
}
