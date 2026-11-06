package server

import "testing"

// The route exists and passes its query through: outside a git work tree the
// answer is "not available", not an error (ADR 261106-60xv).
func TestHistoryRoute(t *testing.T) {
	ts, _ := setup(t)
	r := do(t, ts, "GET", "/api/v1/history?limit=5", "", nil)
	if r.status != 200 || r.body["operation"] != "history" || r.body["ok"] != true {
		t.Fatalf("history: %d %v", r.status, r.body)
	}
	if r.body["available"] != false || r.body["reason"] == nil {
		t.Errorf("want available false with a reason, got %v", r.body)
	}
	r = do(t, ts, "GET", "/api/v1/history?item=nope", "", nil)
	if r.status != 400 {
		t.Errorf("a bad item: %d %v", r.status, r.body)
	}
}
