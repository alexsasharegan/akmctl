// Checks of the pure functions in web/js (run by TestPureModulesUnderNode).
import assert from "node:assert/strict";
import { ageDays, daysBefore, ageText, isoDate } from "./util.js";
import * as L from "./logic.js";
import { parseToken, withQuery } from "./store.js";
import * as M from "./moves.js";
import { diffRows } from "./review.js";
import * as C from "./capture.js";

import * as V from "./live.js";

const work = (id, f) => ({ id, type: "work", path: "docs/work/" + id + ".md", etag: "x", fields: { id, title: "T " + id, status: "ready", parent: null, priority: null, blocked_by: null, updated: "2026-10-20", summary: null, kind: null, ...f } });
const adr = (id, f) => ({ id, type: "adr", path: "docs/adr/" + id + ".md", etag: "x", fields: { id, title: "A " + id, status: "accepted", ...f } });

// dates
assert.equal(ageDays("2026-10-20", "2026-10-28"), 8);
assert.equal(ageDays("2026-10-30", "2026-10-28"), 0);
assert.equal(ageDays("nope", "2026-10-28"), null);
assert.equal(daysBefore(30, "2026-10-28"), "2026-09-28");
assert.equal(ageText(0), "today");
assert.equal(ageText(1), "1 day");
assert.equal(ageText(12), "12 days");
assert.equal(ageText(90), "3 months");
assert.equal(isoDate(new Date("2026-10-28T23:59:00Z")), "2026-10-28");

// token and query
assert.equal(parseToken("#token=abc123"), "abc123");
assert.equal(parseToken("#/board"), "");
assert.equal(parseToken(""), "");
assert.equal(parseToken("#a=1&token=zz"), "zz");
assert.equal(withQuery("records", { fields: ["id", "title"], sort: ["-id"], x: null, y: "", z: undefined, n: 5, b: true }), "records?fields=id%2Ctitle&sort=-id&n=5&b=true");
assert.equal(withQuery("records", {}), "records");

// routes
assert.deepEqual(L.parseRoute(""), { name: "board" });
assert.deepEqual(L.parseRoute("#/board"), { name: "board" });
assert.deepEqual(L.parseRoute("#/tables"), { name: "tables", table: "ready-active" });
assert.deepEqual(L.parseRoute("#/tables/blocked"), { name: "tables", table: "blocked" });
assert.deepEqual(L.parseRoute("#/tables/nope"), { name: "tables", table: "ready-active" });
assert.deepEqual(L.parseRoute("#/adrs"), { name: "adrs" });
assert.deepEqual(L.parseRoute("#/check"), { name: "check" });
assert.deepEqual(L.parseRoute("#/record/261020-a2b3"), { name: "record", id: "261020-a2b3" });
assert.deepEqual(L.parseRoute("#/record/../../x"), { name: "board" });
assert.deepEqual(L.parseRoute("#/record/%E0%A4%A"), { name: "board" });
assert.deepEqual(L.parseRoute("#token=abc"), { name: "board" });
for (const r of [{ name: "board" }, { name: "adrs" }, { name: "check" }, { name: "tables", table: "stale" }, { name: "record", id: "261020-a2b3" }]) {
  assert.deepEqual(L.parseRoute(L.routeHref(r)), r);
}

// lanes and focus
const recs = [
  work("261001-aaaa", { status: "raw" }),
  work("261002-bbbb", { status: "active", priority: 3, parent: "261001-aaaa" }),
  work("261003-cccc", { status: "active", priority: 1, blocked_by: ["261002-bbbb", "legal"] }),
  work("261004-dddd", { status: "blocked", blocked_by: ["x"] }),
  work("261005-eeee", { status: "done", updated: "2026-10-27" }),
  adr("261006-ffff", {}),
];
const ls = L.lanes(recs);
assert.deepEqual(ls.map((l) => l.status + l.records.length), ["raw1", "ready0", "active2", "blocked1"]);
assert.deepEqual(L.locate(ls, "261003-cccc"), { lane: 2, index: 1 });
assert.equal(L.locate(ls, "nope"), null);
assert.deepEqual(L.firstLoc(ls), { lane: 0, index: 0 });
assert.equal(L.firstLoc(L.lanes([])), null);
// right skips the empty ready lane; the row is kept where it can be
assert.deepEqual(L.step(ls, { lane: 0, index: 0 }, "right"), { lane: 2, index: 0 });
assert.deepEqual(L.step(ls, { lane: 2, index: 1 }, "right"), { lane: 3, index: 0 });
assert.deepEqual(L.step(ls, { lane: 3, index: 0 }, "right"), { lane: 3, index: 0 });
assert.deepEqual(L.step(ls, { lane: 2, index: 1 }, "left"), { lane: 0, index: 0 });
assert.deepEqual(L.step(ls, { lane: 0, index: 0 }, "left"), { lane: 0, index: 0 });
assert.deepEqual(L.step(ls, { lane: 2, index: 0 }, "down"), { lane: 2, index: 1 });
assert.deepEqual(L.step(ls, { lane: 2, index: 1 }, "down"), { lane: 2, index: 1 });
assert.deepEqual(L.step(ls, { lane: 2, index: 0 }, "up"), { lane: 2, index: 0 });
assert.deepEqual(L.step(ls, { lane: 2, index: 0 }, "last"), { lane: 2, index: 1 });
assert.deepEqual(L.step(ls, { lane: 2, index: 1 }, "first"), { lane: 2, index: 0 });
assert.deepEqual(L.step(ls, null, "down"), { lane: 0, index: 0 });

// staleness
assert.deepEqual(L.staleness(work("261001-aaaa", { updated: "2026-09-01" }), "2026-10-28"), { days: 57, stale: true });
assert.deepEqual(L.staleness(work("261001-aaaa", { updated: "2026-10-27" }), "2026-10-28"), { days: 1, stale: false });
assert.deepEqual(L.staleness(work("261001-aaaa", { updated: null }), "2026-10-28"), { days: null, stale: false });

// search
assert.equal(L.score("foo bar", "a foo and a bar") > 0, true);
assert.equal(L.score("foo baz", "a foo and a bar"), 0);
assert.equal(L.score("", "x"), 1);
const found = L.searchRecords(recs, "T 261002");
assert.equal(found[0].record.id, "261002-bbbb");
assert.deepEqual(L.searchRecords(recs, "   "), []);
assert.equal(L.searchRecords(recs, "zzzz").length, 0);
assert.equal(L.searchRecords(recs, "T", 2).length, 2);

// tables
assert.deepEqual(L.readyActive(recs).map((r) => r.id), ["261002-bbbb", "261003-cccc"]);
const groups = L.openByParent(recs);
assert.deepEqual(groups.map((g) => g.parent), ["261001-aaaa", null]);
assert.equal(groups[0].known, true);
assert.equal(groups[1].title, "No parent");
const orphan = L.openByParent([work("261001-aaaa", { parent: "261099-zzzz" })]);
assert.equal(orphan[0].known, false);
assert.deepEqual(L.blockedRows(recs).map((r) => r.id), ["261003-cccc", "261004-dddd"]);
assert.deepEqual(L.staleRows([work("261001-aaaa", { updated: "2026-10-27" }), work("261002-bbbb", { updated: "2026-09-01" }), work("261003-cccc", { status: "done" })], "2026-10-28").map((x) => x.record.id), ["261002-bbbb", "261001-aaaa"]);
assert.deepEqual(L.closedRows([work("261001-aaaa", { updated: "2026-10-01" }), work("261002-bbbb", { updated: "2026-10-20" })]).map((r) => r.id), ["261002-bbbb", "261001-aaaa"]);
assert.deepEqual(L.adrRows(recs).map((r) => r.id), ["261006-ffff"]);
assert.deepEqual(L.childrenOf(recs, "261001-aaaa").map((r) => r.id), ["261002-bbbb"]);

// check
const findings = [{ severity: "W", rule: "a" }, { severity: "E", rule: "b" }, { severity: "W", rule: "c" }, { severity: "E", rule: "d" }];
assert.deepEqual(L.orderFindings(findings).map((x) => x.rule), ["b", "d", "a", "c"]);
assert.equal(L.isId("261020-a2b3"), true);
assert.equal(L.isId("261020-A2B3"), false);

// ---- moves ----
assert.equal(M.neighbour("raw", "right"), "ready");
assert.equal(M.neighbour("blocked", "right"), null);
assert.equal(M.neighbour("raw", "left"), null);
assert.equal(M.neighbour("active", "left"), "ready");
assert.equal(M.neighbour("done", "left"), null);

const body = "\n# T\n\n## Goal\nDo it.\n\n```\n## Not a heading\n```\n\n### Sub\nstill goal\n\n## Next\n\n## Log\n- 2026-10-28: x\n";
const secs = M.sectionsOf(body);
assert.deepEqual(Object.keys(secs), ["goal", "next", "log"]);
assert.match(secs.goal, /^Do it\./);
assert.match(secs.goal, /still goal$/);
assert.equal(secs.next, "");
assert.deepEqual(M.sectionsOf(""), {});
assert.deepEqual(M.sectionsOf("## A ##\nx\r\n## b\r\ny"), { a: "x", b: "y" });

assert.deepEqual(M.parseEntries(" a \n\n b\r\na\n"), ["a", "b"]);

// the wizard rules: WK5 to WK7
const base = { from: "raw", to: "ready", claims: false, sections: {}, blockedBy: [] };
assert.deepEqual(M.planMove(base), { kind: "gated", steps: ["goal"] });
assert.deepEqual(M.planMove({ ...base, sections: { goal: "x" } }), { kind: "safe", steps: [] });
assert.deepEqual(M.planMove({ ...base, to: "active", sections: { goal: "x" } }), { kind: "gated", steps: ["next"] });
assert.deepEqual(M.planMove({ ...base, to: "active", sections: { goal: "x", next: "y" } }), { kind: "safe", steps: [] });
assert.deepEqual(M.planMove({ ...base, to: "active", sections: { goal: "x", next: "y" }, claims: true }), { kind: "safe", steps: [] });
assert.deepEqual(M.planMove({ ...base, to: "active", claims: true }), { kind: "gated", steps: ["goal", "next", "identity"] });
assert.deepEqual(M.planMove({ ...base, to: "blocked" }), { kind: "gated", steps: ["goal", "next", "blocked_by"] });
assert.deepEqual(M.planMove({ ...base, to: "blocked", sections: { goal: "x", next: "y" }, blockedBy: ["legal"] }), { kind: "gated", steps: [] });
assert.deepEqual(M.planMove({ ...base, to: "done" }), { kind: "gated", steps: ["goal", "outcome", "decision"] });
assert.deepEqual(M.planMove({ ...base, to: "done", sections: { goal: "g", outcome: "o" } }), { kind: "gated", steps: ["decision"] });
assert.deepEqual(M.planMove({ ...base, to: "dropped" }), { kind: "gated", steps: ["outcome"] });
assert.deepEqual(M.planMove({ ...base, from: "ready", to: "raw" }), { kind: "safe", steps: [] });
assert.equal(M.planMove({ ...base, to: "raw" }).kind, "noop");
assert.equal(M.planMove({ ...base, from: "done", to: "ready" }).kind, "noop");

assert.equal(M.validateStep("goal", { goal: "  " }), "Say what done looks like. A ready, active, blocked or done item needs a Goal.");
assert.equal(M.validateStep("goal", { goal: "x" }), "");
assert.notEqual(M.validateStep("blocked_by", { blocked_by: "\n \n" }), "");
assert.equal(M.validateStep("blocked_by", { blocked_by: "261001-aaaa\nlegal review" }), "");
assert.match(M.validateStep("outcome", { outcome: "" }, "dropped"), /reason/);
assert.match(M.validateStep("outcome", { outcome: "" }, "done"), /outcome/);
assert.notEqual(M.validateStep("decision", { decision: { mode: "link", id: "" } }), "");
assert.equal(M.validateStep("decision", { decision: { mode: "link", id: "261012-p3q4" } }), "");
assert.equal(M.validateStep("decision", { decision: { mode: "none" } }), "");
assert.notEqual(M.validateStep("decision", { decision: { mode: "new", title: "x", slug: "" } }), "");

// ops and requests
const mv = { id: "261022-e5f6", from: "ready", to: "active", claims: false, decisions: [] };
assert.deepEqual(M.buildOps(mv, {}), [{ kind: "set", changes: { status: "active" } }]);
assert.deepEqual(M.buildOps({ ...mv, from: "raw", to: "ready" }, { goal: " Done when x " }), [{ kind: "section", name: "Goal", text: "Done when x", mode: "replace" }, { kind: "set", changes: { status: "ready" } }]);
assert.deepEqual(M.buildOps({ ...mv, claims: true }, { identity: "alex" }), [{ kind: "claim", as: "alex" }]);
assert.deepEqual(M.buildOps({ ...mv, claims: true, from: "blocked" }, { next: "n", identity: "alex" }), [{ kind: "section", name: "Next", text: "n", mode: "replace" }, { kind: "set", changes: { status: "active" } }, { kind: "claim", as: "alex" }]);
assert.deepEqual(M.buildOps({ ...mv, to: "blocked" }, { next: "n", blocked_by: "a\nb" }), [{ kind: "section", name: "Next", text: "n", mode: "replace" }, { kind: "set", changes: { status: "blocked", blocked_by: ["a", "b"] } }]);
assert.deepEqual(M.buildOps({ ...mv, to: "blocked" }, {}), [{ kind: "set", changes: { status: "blocked" } }]);
assert.deepEqual(M.buildOps({ ...mv, to: "dropped" }, { outcome: " no longer needed " }), [{ kind: "close", status: "dropped", outcome: "no longer needed" }]);
assert.deepEqual(M.buildOps({ ...mv, to: "done" }, {}), [{ kind: "close", status: "done" }]);
assert.deepEqual(M.buildOps({ ...mv, to: "done", decisions: ["261012-p3q4"] }, { outcome: "o", decision: { mode: "link", id: "261012-p3q4" } }), [{ kind: "set", changes: { decisions: ["261012-p3q4"] } }, { kind: "close", status: "done", outcome: "o" }]);
const withNew = M.buildOps({ ...mv, to: "done" }, { outcome: "o", decision: { mode: "new", title: " Pick X ", slug: "pick-x" } });
assert.deepEqual(withNew.map((o) => o.kind), ["new", "set", "close"]);
assert.deepEqual(withNew[0], { kind: "new", what: "adr", slug: "pick-x", title: "Pick X" });
assert.deepEqual(M.resolveOp(withNew[1], "261028-abcd").changes.decisions, ["261028-abcd"]);
assert.equal(M.previewable(withNew[0], mv), false);
assert.equal(M.previewable(withNew[1], mv), false);
assert.equal(M.previewable(withNew[2], mv), true);
assert.equal(M.previewable({ kind: "claim", as: "a" }, { from: "ready" }), true);
assert.equal(M.previewable({ kind: "claim", as: "a" }, { from: "raw" }), false);

assert.deepEqual(M.requestFor({ kind: "set", changes: { status: "ready" } }, "261022-e5f6"), { method: "PATCH", path: "records/261022-e5f6", body: { changes: { status: "ready" } } });
assert.deepEqual(M.requestFor({ kind: "set", changes: { status: "ready" } }, "261022-e5f6", { dry_run: true }).body, { changes: { status: "ready" }, dry_run: true });
assert.deepEqual(M.requestFor({ kind: "section", name: "Next", text: "t", mode: "replace" }, "i"), { method: "PUT", path: "records/i/sections/Next", body: { text: "t", mode: "replace" } });
assert.equal(M.requestFor({ kind: "section", name: "A b/c", text: "", mode: "replace" }, "i").path, "records/i/sections/A%20b%2Fc");
assert.deepEqual(M.requestFor({ kind: "close", status: "done", outcome: "o" }, "i", { dry_run: true }), { method: "POST", path: "records/i/close", body: { status: "done", outcome: "o", dry_run: true } });
assert.deepEqual(M.requestFor({ kind: "close", status: "dropped" }, "i").body, { status: "dropped" });
assert.deepEqual(M.requestFor({ kind: "claim", as: "a" }, "i", { takeover: true }), { method: "POST", path: "records/i/claim", body: { as: "a", takeover: true } });
assert.deepEqual(M.requestFor({ kind: "release" }, "i"), { method: "DELETE", path: "records/i/claim", body: undefined });
assert.deepEqual(M.requestFor({ kind: "new", what: "adr", slug: "s", title: "T" }, "i"), { method: "POST", path: "records", body: { kind: "adr", slug: "s", title: "T" } });
assert.throws(() => M.requestFor({ kind: "nope" }, "i"));
assert.match(M.opLabel({ kind: "close", status: "done", outcome: "x" }, "i"), /Close as done/);

// ---- review rows ----
const d = diffRows({ path: "docs/work/a.md", new_path: "docs/work/done/a.md", hunks: [{ header: "@@ -1,3 +1,3 @@", lines: [{ op: "ctx", text: "a", old: 1, new: 1 }, { op: "rem", text: "status: ready", old: 2, new: null }, { op: "add", text: "status: done", old: null, new: 2 }, { op: "add", text: "", old: null, new: 3 }] }] });
assert.deepEqual(d.rows.map((r) => r.kind), ["hunk", "ctx", "rem", "add", "add"]);
assert.deepEqual(d.rows.map((r) => r.sign), ["", " ", "−", "+", "+"]);
assert.deepEqual(d.rows.map((r) => r.word), ["hunk header", "unchanged", "removed", "added", "added"]);
assert.equal(d.added, 2);
assert.equal(d.removed, 1);
assert.equal(d.renamed, true);
assert.equal(d.rows[2].old, 2);
assert.equal(d.rows[2].new, null);
assert.equal(diffRows({ path: "p", new_path: "p", hunks: [] }).rows.length, 0);
assert.equal(diffRows({ path: "p", new_path: "p", hunks: [] }).renamed, false);
assert.equal(diffRows(null).rows.length, 0);

// ---- slug (vectors from pkm-workbench internal/move/block_test.go) ----
for (const [input, want] of Object.entries({
  "Fix the sink, again!": "fix-sink-again",
  "Café run ☕ for Erin's list of things to buy": "cafe-run-erin-s-list-things",
  "Crème brûlée à la Zoë": "creme-brulee-la-zoe",
  "the to of": "the-to-of",
  "???": "",
  "Straße and Łódź": "strasse-lodz",
  "": "",
  "  Call   the plumber  ": "call-plumber",
  "Ünïcode ÇÅ": "unicode-ca",
  "v2.0 release 2026": "v2-0-release-2026",
  "日本語": "",
})) {
  assert.equal(C.slugify(input), want, JSON.stringify(input));
}
const long = C.slugify("abcdefgh ".repeat(8));
assert.ok(long.length <= 48 && !long.endsWith("-"), long);
assert.equal(long, "abcdefgh-abcdefgh-abcdefgh-abcdefgh-abcdefgh-abc");
assert.equal(C.slugify("one two three four five six seven eight"), "one-two-three-four-five-six");
// a typed slug keeps its filler and may run to eight words and 60 characters
assert.equal(C.cleanSlug("The Plumber Thing!"), "the-plumber-thing");
assert.equal(C.cleanSlug("a b c d e f g h i j"), "a-b-c-d-e-f-g-h");
assert.ok(C.cleanSlug("abcdefghij ".repeat(8)).length <= 60);
assert.equal(C.validSlug("fix-it-2"), true);
assert.equal(C.validSlug("Fix-it"), false);
assert.equal(C.validSlug("fix--it"), false);
assert.equal(C.validSlug("x".repeat(61)), false);

// the capture request
assert.deepEqual(C.buildCapture({ title: "Call the plumber" }), { body: { kind: "work", slug: "call-plumber", title: "Call the plumber" } });
assert.deepEqual(C.buildCapture({ title: " Call  the plumber ", slug: "My Own Slug", summary: " s ", kind: "bug", priority: "3", parent: "261015-c4d5" }, ["261015-c4d5"]), { body: { kind: "work", slug: "my-own-slug", title: "Call the plumber", fields: { summary: "s", kind: "bug", priority: 3 , parent: "261015-c4d5" } } });
assert.equal(C.buildCapture({ title: "  " }).field, "title");
assert.equal(C.buildCapture({ title: "???" }).field, "slug");
assert.equal(C.buildCapture({ title: "x", priority: "1.5" }).field, "priority");
assert.equal(C.buildCapture({ title: "x", priority: "-2" }).body.fields.priority, -2);
assert.equal(C.buildCapture({ title: "x", priority: "007" }).field, "priority");
assert.equal(C.buildCapture({ title: "x", kind: "epic" }).field, "kind");
assert.equal(C.buildCapture({ title: "x", parent: "261099-zzzz" }, []).field, "parent");
assert.equal(C.buildCapture({ title: "x", summary: "y".repeat(161) }).field, "summary");


// live: the snapshot diff
const rec = (id, status, etag, path) => ({ id, type: "work", path: path || "docs/work/" + id + ".md", etag, fields: { status } });
const snap = (...r) => V.snapshotOf(r);
const none = { has: () => false };
{
  const a = snap(rec("a", "ready", "e1"), rec("b", "active", "e2"), rec("c", "raw", "e3"), rec("d", "ready", "e4"));
  assert.deepEqual(V.diffSnapshots(a, a, none), []);
  assert.equal(V.snapshotOf([{ id: "x", type: "adr", etag: "z", fields: { status: "accepted" } }]).size, 0);
  // moved, added, removed, edited, closed
  const b = snap(
    rec("a", "active", "e1b"),
    rec("b", "active", "e2b"),
    rec("d", "done", "e4b", "docs/work/done/d.md"),
    rec("e", "raw", "e5"),
  );
  assert.deepEqual(V.diffSnapshots(a, b, none), [
    { type: "removed", id: "c", from: "raw", source: "outside" },
    { type: "moved", id: "a", from: "ready", to: "active", closed: false, source: "outside" },
    { type: "edited", id: "b", source: "outside" },
    { type: "moved", id: "d", from: "ready", to: "done", closed: true, source: "outside" },
    { type: "added", id: "e", to: "raw", source: "outside" },
  ]);
  // a file moved into done/ with the same status is a move; a move within done/ is not
  const closedOnly = V.diffSnapshots(snap(rec("a", "done", "e1")), snap(rec("a", "done", "e1", "docs/work/done/a.md")), none);
  assert.deepEqual(closedOnly, [{ type: "moved", id: "a", from: "done", to: "done", closed: true, source: "outside" }]);
  assert.deepEqual(V.diffSnapshots(snap(rec("a", "done", "e1", "docs/work/done/a.md")), snap(rec("a", "done", "e2", "docs/work/done/a.md")), none).map((e) => e.type), ["edited"]);
  // own writes are "you"; an outside change after one is not
  const own = V.createOwn();
  own.note("a", "e1b");
  assert.equal(V.diffSnapshots(a, b, own).find((e) => e.id === "a").source, "you");
  assert.equal(V.diffSnapshots(a, b, own).find((e) => e.id === "b").source, "outside");
  const c = snap(rec("a", "blocked", "e1c"));
  assert.equal(V.diffSnapshots(snap(rec("a", "active", "e1b")), c, own)[0].source, "outside");
  own.note("n", "e9");
  assert.equal(V.diffSnapshots(new Map(), snap(rec("n", "raw", "e9")), own)[0].source, "you");
  // the words
  assert.equal(V.describeEvent({ type: "moved", id: "a", from: "ready", to: "active", source: "outside" }, "Title"), "Outside: a Title moved, ready to active.");
  assert.equal(V.describeEvent({ type: "moved", id: "a", from: "ready", to: "active", source: "you" }, ""), "You moved a: ready to active.");
  assert.equal(V.describeEvent({ type: "moved", id: "a", from: "done", to: "done", closed: true, source: "outside" }, ""), "Outside: a moved, closed.");
  assert.equal(V.summarize([{ type: "moved", source: "outside" }, { type: "moved", source: "outside" }, { type: "added", source: "outside" }]), "3 changes outside (2 moved, 1 added).");
}

// live: the loop (injected fetch, no DOM)
{
  const wait = (ms) => new Promise((r) => setTimeout(r, ms));
  const seen = [];
  const states = [];
  let mode = "same";
  let calls = 0;
  const base = [rec("a", "ready", "e1")];
  const lv = V.createLive({
    interval: 5,
    onStatus: (s) => states.push(s),
    onEvents: (ev) => {
      if (mode === "reject") {
        mode = "same2";
        return Promise.reject(new Error("refresh failed"));
      }
      seen.push(ev);
    },
    fetchWork: async () => {
      calls++;
      if (mode === "down") throw Object.assign(new Error("x"), { code: "unreachable" });
      return { records: mode === "moved" || mode === "reject" || mode === "same2" ? [rec("a", "active", "e2")] : base };
    },
  });
  lv.start(base);
  await wait(30);
  assert.equal(seen.length, 0);
  mode = "down";
  await wait(250);
  assert.ok(states.includes("offline"), "offline after repeated failures: " + states);
  mode = "same";
  await wait(250);
  assert.equal(states[states.length - 1], "ok");
  assert.equal(seen.length, 0, "no events when nothing changed while offline");
  mode = "reject";
  await wait(150);
  assert.equal(seen.length, 1, "the events are offered again after a failed refresh");
  assert.equal(seen[0][0].to, "active");
  await wait(60);
  assert.equal(seen.length, 1, "and only once after they were taken");
  lv.stop();
  const n = calls;
  await wait(40);
  assert.equal(calls, n, "stopped");
}
assert.equal(L.obsidianUrl({ root: "/Users/a b/repo", obsidian_vault: "" }, "docs/work/x.md"), "obsidian://open?path=%2FUsers%2Fa%20b%2Frepo%2Fdocs%2Fwork%2Fx.md");
assert.equal(L.obsidianUrl({ root: "/r/", obsidian_vault: "My Vault" }, "docs/é ü.md"), "obsidian://open?vault=My%20Vault&file=docs%2F%C3%A9%20%C3%BC.md");
assert.equal(L.obsidianUrl({ root: "", obsidian_vault: "v" }, "a.md"), "");
console.log("ok");
