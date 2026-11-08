// Checks of the pure functions in web/js (run by TestPureModulesUnderNode).
import assert from "node:assert/strict";
import { ageDays, daysBefore, ageText, isoDate } from "./util.js";
import * as L from "./logic.js";
import { parseToken, withQuery } from "./store.js";
import * as M from "./moves.js";
import { diffRows } from "./review.js";
import * as C from "./capture.js";

import * as V from "./live.js";
import * as R from "./replay.js";
import * as K from "./controls.js";
import * as P from "./picker.js";
import * as F from "./fields.js";

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
assert.equal(C.buildCapture({ title: "x", cycle: " 2026-s3 " }).body.fields.cycle, "2026-s3");
assert.equal(C.buildCapture({ title: "x", cycle: "" }).body.fields, undefined);
for (const bad of ["-s3", "has space", "a/b", "x".repeat(65)]) assert.equal(C.buildCapture({ title: "x", cycle: bad }).field, "cycle");
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
assert.equal(L.obsidianUrl({ root: "/Users/a b/repo" }, "docs/work/x.md"), "obsidian://open?path=%2FUsers%2Fa%20b%2Frepo%2Fdocs%2Fwork%2Fx.md");
assert.equal(L.obsidianUrl({ root: "/r/" }, "docs/é ü.md"), "obsidian://open?path=%2Fr%2Fdocs%2F%C3%A9%20%C3%BC.md");
assert.equal(L.obsidianUrl({ root: "" }, "a.md"), "");
{
  const base = { from: "blocked", to: "ready", claims: true, sections: { goal: "g", next: "n" }, blockedBy: ["x"], decisions: [] };
  assert.deepEqual(M.planMove(base).steps, ["unblock"]);
  assert.deepEqual(M.buildOps(base, { unblock: true })[0].changes, { status: "ready", blocked_by: null });
  assert.deepEqual(M.buildOps(base, { unblock: false })[0].changes, { status: "ready" });
  const rel = M.buildOps({ ...base, from: "active", blockedBy: [], claimed: true }, {});
  assert.deepEqual(rel.map((o) => o.kind), ["release", "set"]);
}
{
  const k = (o) => M.wizardKey({ key: "Enter", mod: false, tag: "DIV", type: "", phase: "edit", step: "unblock", ...o });
  assert.equal(k({ tag: "INPUT", type: "checkbox" }), "next");
  assert.equal(k({ mod: true, tag: "INPUT", type: "checkbox" }), "next");
  assert.equal(k({ tag: "INPUT", type: "text", step: "identity" }), "next");
  assert.equal(k({ tag: "TEXTAREA", step: "goal" }), "");
  assert.equal(k({ tag: "TEXTAREA", step: "goal", mod: true }), "next");
  assert.equal(k({ tag: "SELECT" }), "");
  assert.equal(k({ tag: "BUTTON" }), "");
  assert.equal(k({ step: "review", phase: "review" }), "write");
  assert.equal(k({ step: "review", phase: "review", mod: true }), "");
  assert.equal(k({ step: "review", phase: "review", tag: "BUTTON" }), "");
  assert.equal(k({ step: "review", phase: "conflict" }), "reload");
  assert.equal(k({ step: "review", phase: "loading" }), "");
  assert.equal(k({ key: "Escape", tag: "TEXTAREA" }), "back");
  assert.equal(k({ key: "a" }), "");
}
console.log("ok");

// the repository's folders come from version.repository, with the defaults when it names none
assert.deepEqual(L.foldersFrom({ repository: { work_dir: "Notes/tasks", adr_dir: "Notes/adr" } }), { work_dir: "Notes/tasks", adr_dir: "Notes/adr" });
assert.deepEqual(L.foldersFrom({ repository: null }), L.DEFAULT_FOLDERS);
assert.deepEqual(L.foldersFrom(undefined), { work_dir: "docs/work", adr_dir: "docs/adr" });
assert.deepEqual(L.foldersFrom({ repository: { work_dir: "", adr_dir: 5 } }), L.DEFAULT_FOLDERS);
assert.equal(C.fileNote("Notes/tasks", "call-plumber"), "File: Notes/tasks/<id>-call-plumber.md");
assert.equal(C.fileNote("docs/work", ""), "");
assert.match(M.opLabel({ kind: "close", status: "done", outcome: "x" }, "i", { work_dir: "Notes/tasks" }), /move the file to Notes\/tasks\/done\/$/);
assert.match(M.opLabel({ kind: "close", status: "dropped" }, "i"), /move the file to docs\/work\/done\/$/);

// page title: the repository's folder name, with the parent folder when another root shares it
assert.equal(L.repoLabel("/Users/me/github.com/alexsasharegan/akmctl"), "akmctl");
assert.equal(L.repoLabel("C:\\Users\\me\\notes.vaults\\"), "notes.vaults");
assert.equal(L.repoLabel("/home/me/work/docs", ["/home/me/play/docs"]), "work/docs");
assert.equal(L.repoLabel("/home/me/work/docs", ["/home/me/work/docs"]), "docs");
assert.equal(L.repoLabel("/home/me/work/docs", ["/home/me/play/akmctl"]), "docs");
assert.equal(L.repoLabel("docs", ["/x/docs"]), "docs");
assert.equal(L.repoLabel(""), "");
const longName = "a-repository-with-an-extremely-long-name-that-goes-on-and-on-for-a-while-indeed";
assert.equal(L.pageTitle({ view: "Board", repo: longName }), "Board - " + longName + " - akmctl");
assert.equal(L.pageTitle({ view: "Record 261029-9ep9: Name the page", repo: "akmctl" }), "Record 261029-9ep9: Name the page - akmctl - akmctl");
assert.equal(L.pageTitle({ view: "ADRs" }), "ADRs - akmctl");

// toast queue
{
  const T = await import("./toasts.js");
  const sp = (o) => T.toastSpec("m", o);
  assert.deepEqual(sp(), { msg: "m", level: "info", persistent: false, ms: 15000 });
  assert.equal(sp({ level: "error" }).persistent, true);
  assert.equal(sp({ level: "error", persistent: false }).persistent, false);
  assert.equal(sp({ level: "success", persistent: true }).persistent, true);
  assert.equal(sp({ ms: 12000 }).ms, 12000);
  assert.equal(sp({ ms: -1 }).ms, 15000);
  assert.throws(() => T.toastSpec("m", "bad"), TypeError);
  assert.throws(() => T.toastSpec("m", { level: "bad" }), RangeError);
  assert.equal(T.roleFor("error"), "alert");
  assert.equal(T.roleFor("warning"), "alert");
  assert.equal(T.roleFor("success"), "status");
  assert.equal(T.roleFor("info"), "status");

  // the limit evicts the oldest ephemeral toast
  let q = T.createToastQueue();
  const a = q.add(sp(), 0).id;
  const b = q.add(sp(), 1).id;
  const c = q.add(sp(), 2).id;
  const d = q.add(sp(), 3);
  assert.deepEqual(d.evicted, [a]);
  assert.deepEqual(q.ids(), [b, c, d.id]);

  // a persistent toast is not evicted by an ephemeral one
  q = T.createToastQueue();
  const p1 = q.add(sp({ level: "error" }), 0).id;
  const e1 = q.add(sp(), 0).id;
  const p2 = q.add(sp({ level: "error" }), 0).id;
  let r = q.add(sp(), 0);
  assert.deepEqual(r.evicted, [e1]);
  assert.deepEqual(q.ids(), [p1, p2, r.id]);
  q.add(sp({ level: "error" }), 0); // evicts the ephemeral one again
  r = q.add(sp(), 0); // full of persistent: the newcomer is dropped
  assert.equal(r.id, null);
  assert.equal(q.ids().length, 3);
  r = q.add(sp({ level: "error" }), 0); // a persistent newcomer evicts the oldest persistent
  assert.deepEqual(r.evicted, [p1]);
  assert.equal(q.newestPersistent(), r.id);
  assert.equal(q.dismiss(r.id), true);
  assert.equal(q.dismiss(r.id), false);
  assert.equal(q.newestPersistent(), q.ids().at(-1));
  assert.equal(T.createToastQueue().newestPersistent(), null);

  // timers: expiry, pause, resume, two holds
  q = T.createToastQueue();
  const t = q.add(sp({ ms: 1000 }), 100).id;
  const pe = q.add(sp({ level: "error" }), 100).id;
  assert.equal(q.remaining(t, 100), 1000);
  assert.equal(q.remaining(t, 400), 700);
  assert.equal(q.remaining(pe, 400), null);
  assert.deepEqual(q.expired(1099), []);
  assert.deepEqual(q.expired(1100), [t]);
  q.pause(t, "hover", 400);
  assert.equal(q.paused(t), true);
  assert.equal(q.remaining(t, 5000), 700);
  assert.deepEqual(q.expired(5000), []);
  q.pause(t, "focus", 5000);
  q.resume(t, "hover", 6000);
  assert.equal(q.paused(t), true); // focus still holds it
  assert.equal(q.remaining(t, 7000), 700);
  q.resume(t, "focus", 7000);
  assert.equal(q.paused(t), false);
  assert.equal(q.remaining(t, 7200), 500);
  assert.deepEqual(q.expired(7700), [t]);
  q.resume(t, "focus", 8000); // a release without a hold changes nothing
  assert.equal(q.remaining(t, 8000), 0);
  assert.deepEqual(q.clear().sort(), [t, pe].sort());
  assert.deepEqual(q.ids(), []);
}

// project filter
assert.equal(L.projectTarget('[[Project Notes]]'), "Project Notes");
assert.equal(L.projectTarget("[[Acme Fabrication|Acme]]"), "Acme Fabrication");
assert.equal(L.projectTarget("[[Notes#Plan]]"), "Notes");
assert.equal(L.projectTarget("  plain  "), "plain");
assert.equal(L.projectTarget(null), "");
assert.equal(L.projectTarget("[[]]"), "");
{
  const rs = [work("P", { title: "Alpha" }), work("a", { parent: "P" }), work("b", { parent: "a" }), work("Q", { title: "Beta" }), work("c", { parent: "Q" }), work("d", {}), adr("z", {})];
  const c = L.projectCounts(rs);
  assert.deepEqual(c.list, [{ target: "P", label: "Alpha", count: 3 }, { target: "Q", label: "Beta", count: 2 }]);
  assert.equal(c.none, 1);
  assert.equal(L.hasProjects(rs), true);
  assert.equal(L.hasProjects([work("d", {}), adr("z", {})]), false);
  assert.deepEqual(L.filterByProject(rs, null), rs);
  assert.deepEqual(L.filterByProject(rs, "P").map((r) => r.id), ["P", "a", "b", "z"]);
  assert.deepEqual(L.filterByProject(rs, L.NO_PROJECT).map((r) => r.id), ["d", "z"]);
  assert.deepEqual(L.filterByProject(rs, "Nope").map((r) => r.id), ["z"]);
}
assert.equal(L.projectFromHash("#/board?project=Project%20Notes"), "Project Notes");
assert.equal(L.projectFromHash("#/board"), null);
assert.equal(L.projectFromHash("#/board?project="), null);
assert.equal(L.projectFromHash("#/tables/blocked?project=~"), "~");
assert.equal(L.withProject("#/board", "Project Notes"), "#/board?project=Project%20Notes");
assert.equal(L.withProject("#/board?project=x", null), "#/board");
assert.deepEqual(L.parseRoute("#/tables/blocked?project=x"), { name: "tables", table: "blocked" });
assert.deepEqual(L.parseRoute("#/record/261020-a2b3?project=x"), { name: "record", id: "261020-a2b3" });

// wikilinks in field values
assert.deepEqual(L.wikilinkParts("[[Project Notes]]"), { target: "Project Notes", alias: "", heading: "" });
assert.deepEqual(L.wikilinkParts("[[Acme Fabrication|Acme]]"), { target: "Acme Fabrication", alias: "Acme", heading: "" });
assert.deepEqual(L.wikilinkParts("[[Notes#Plan|the plan]]"), { target: "Notes", alias: "the plan", heading: "Plan" });
assert.deepEqual(L.wikilinkParts(" [[261025-g68j-build-pkm-workbench]] "), { target: "261025-g68j-build-pkm-workbench", alias: "", heading: "" });
assert.equal(L.wikilinkParts("see [[A]]"), null);
assert.equal(L.wikilinkParts("[[]]"), null);
assert.equal(L.wikilinkParts("[[a]] [[b]]"), null);
assert.equal(L.wikilinkParts(5), null);
assert.equal(L.wikilinkParts(null), null);
assert.equal(L.wikilinkLabel(L.wikilinkParts("[[A|b]]")), "b");
assert.equal(L.wikilinkLabel(L.wikilinkParts("[[A]]")), "A");

// ---- replay of the moves made while away (261101-w0mz) ----
{
  const mv = (item, step, from, to, extra) => ({ sha: "s" + item + step, time: "2026-11-06T12:22:15-07:00", author: "Alex Regan", item, step, from, to, trailers: ["Assisted-by: Claude Code"], ...extra });
  // who: the Assisted-by value, else the author
  assert.equal(R.whoOf(mv("a", "claim", "ready", "active")), "Claude Code");
  assert.equal(R.whoOf(mv("a", "claim", "ready", "active", { trailers: [] })), "Alex Regan");
  assert.equal(R.whoOf(mv("a", "claim", "ready", "active", { trailers: ["Reviewed-by: x"] })), "Alex Regan");
  // a lane change, a close, a new file and a claim become different events
  assert.deepEqual(R.moveToEvent(mv("a", "start", "ready", "active")).type, "moved");
  const closed = R.moveToEvent(mv("a", "close", "active", "done"));
  assert.equal(closed.type, "moved");
  assert.equal(closed.closed, true);
  assert.equal(R.moveToEvent(mv("a", "file", "", "raw")).type, "added");
  assert.equal(R.moveToEvent(mv("a", "claim", "ready", "ready")).type, "step");
  assert.equal(R.moveToEvent(mv("a", "claim", "ready", "ready")).source, "outside");

  // the words name who made the move
  assert.equal(V.describeEvent(R.moveToEvent(mv("a", "start", "ready", "active")), "Title"), "Claude Code moved a Title: ready to active.");
  assert.equal(V.describeEvent(R.moveToEvent(mv("a", "claim", "ready", "ready")), ""), "Claude Code claimed a.");
  assert.equal(V.describeEvent(R.moveToEvent(mv("a", "release", "active", "active")), ""), "Claude Code released a.");
  assert.equal(V.describeEvent({ type: "moved", id: "a", from: "ready", to: "active", source: "outside" }, "Title"), "Outside: a Title moved, ready to active.");

  // ordering: the moves replace the poll's net event, in commit order (not id order), path included
  const net = [
    { type: "moved", id: "a", from: "ready", to: "done", closed: true, source: "outside" },
    { type: "moved", id: "b", from: "ready", to: "active", closed: false, source: "outside" },
  ];
  const moves = [mv("b", "start", "ready", "active"), mv("a", "start", "ready", "active"), mv("a", "close", "active", "done")];
  const merged = R.mergeWithHistory(net, moves);
  assert.deepEqual(merged.map((e) => e.id + ":" + e.from + ">" + e.to), ["b:ready>active", "a:ready>active", "a:active>done"]);
  assert.ok(merged.every((e) => e.who === "Claude Code"));

  // what the history does not explain stays as the poll saw it, after the replayed moves
  const stay = R.mergeWithHistory(
    [
      { type: "moved", id: "a", from: "ready", to: "active", source: "outside" },
      { type: "moved", id: "c", from: "ready", to: "blocked", source: "outside" }, // never committed
      { type: "edited", id: "d", source: "outside" },
      { type: "moved", id: "e", from: "ready", to: "active", source: "you" }, // this page's own move
    ],
    [mv("a", "start", "ready", "active"), mv("e", "start", "ready", "active")],
  );
  assert.deepEqual(stay.map((e) => e.id + ":" + e.type + ":" + (e.who ? "history" : e.source)), ["a:moved:history", "c:moved:outside", "d:edited:outside", "e:moved:you"]);

  // a history that ends somewhere else than the poll does not replace it
  const off = R.mergeWithHistory([{ type: "moved", id: "a", from: "ready", to: "done", closed: true, source: "outside" }], [mv("a", "start", "ready", "active")]);
  assert.deepEqual(off.map((e) => e.source), ["outside"]);
  assert.equal(off[0].who, undefined);
  assert.deepEqual(R.mergeWithHistory(net, []), net);
  assert.deepEqual(R.mergeWithHistory(net, null), net);

  // the batching rule: the animator plays at most AWAY_MAX, so the panel lists what it will not
  assert.equal(R.AWAY_MAX, 10);
  assert.equal(R.AWAY_MAX, 10);
  assert.equal(R.awayHeading(1, false), "1 move while you were away");
  assert.equal(R.awayHeading(12, true), "12 moves while you were away (older ones are not listed)");
  const lines = R.awayLines([R.moveToEvent(mv("a", "start", "ready", "active"))], () => "T", new Date("2026-11-06T20:00:00Z"));
  assert.deepEqual(lines, [{ time: "12:22", text: "Claude Code moved a T: ready to active." }]);
  assert.equal(R.timeLabel("2026-11-05T09:03:00-07:00", new Date("2026-11-06T20:00:00Z")), "2026-11-05 09:03");
  assert.equal(R.timeLabel("", new Date()), "");
  // a burst summary counts claims and releases
  assert.equal(V.summarize([{ type: "step", source: "outside" }, { type: "moved", source: "outside" }]), "2 changes outside (1 moved, 1 claimed or released).");
}

// shared controls: the summary counter
{
  const lv = (n) => K.counterState(n).level;
  assert.equal(lv(0), "ok");
  assert.equal(lv(143), "ok"); // 90% of 160 is 144
  assert.equal(lv(144), "near");
  assert.equal(lv(159), "near");
  assert.equal(lv(160), "limit");
  assert.equal(lv(161), "over");
  assert.equal(K.counterState(150).text, "150 / 160 - 10 left");
  assert.equal(K.counterState(160).text, "160 / 160 - limit reached");
  assert.equal(K.counterState(165).text, "165 / 160 - 5 over the limit");
  assert.equal(K.counterState(9, 10).level, "near"); // another limit: ceil(9) = 9
  assert.equal(K.oneLine("a\nb\r\n  c"), "a b c");

  // the select's key handling
  const opts = [{ value: "", label: "(default: chore)" }, { value: "feature", label: "feature" }, { value: "bug", label: "bug" }, { value: "chore", label: "chore" }, { value: "spike", label: "spike" }];
  const closed = { open: false, active: 0, buf: "", at: 0 };
  const open = (active) => ({ open: true, active, buf: "", at: 0 });
  const press = (st, key, now = 0) => K.selectKey(st, key, opts, now);
  assert.deepEqual([press(closed, "ArrowDown").st.open, press(closed, "ArrowDown").st.active], [true, 0]); // opens first, on the chosen option
  assert.equal(press(closed, "Enter").st.open, true);
  assert.equal(press(closed, " ").handled, true);
  assert.equal(press(open(0), "ArrowDown").st.active, 1);
  assert.equal(press(open(4), "ArrowDown").st.active, 4); // stops at the end
  assert.equal(press(open(0), "ArrowUp").st.active, 0);
  assert.equal(press(open(2), "Home").st.active, 0);
  assert.equal(press(open(2), "End").st.active, 4);
  const pick = press(open(2), "Enter");
  assert.deepEqual([pick.choose, pick.st.open], [2, false]);
  assert.equal(press(open(2), " ").choose, 2);
  // Escape closes an open list and is used; on a closed select it is left to the dialog
  assert.deepEqual([press(open(2), "Escape").st.open, press(open(2), "Escape").handled, press(open(2), "Escape").choose], [false, true, null]);
  assert.equal(press(closed, "Escape").handled, false);
  // Tab closes without choosing and is not swallowed
  const tab = press(open(2), "Tab");
  assert.deepEqual([tab.st.open, tab.handled, tab.choose], [false, false, null]);
  // type-ahead: a prefix typed fast, a restart after a pause, one letter repeated cycles
  assert.equal(press(open(0), "b", 1000).st.active, 2);
  const s = press(open(0), "s", 1000);
  assert.equal(s.st.active, 4);
  assert.equal(press(s.st, "p", 1100).st.active, 4); // "sp" -> spike
  assert.equal(press(s.st, "f", 1100).st.active, 4); // nothing starts with "sf": the active option stays
  assert.equal(press(open(1), "c", 1000).st.active, 3);
  assert.equal(press({ ...open(3), buf: "c", at: 1000 }, "c", 1100).st.active, 3); // only one c option: stays
  // a modified key and an unrelated one are not handled
  assert.equal(press(open(0), "F5").handled, false);
  assert.equal(press(open(0), "x", 1000).st.active, 0); // no match leaves the active option
}

// the parent picker: matching, ranking and highlighting
{
  const items = [
    { id: "261029-bdt9", title: "Workbench parent picker: a searchable list", summary: "Replace the datalist" },
    { id: "261028-esyt", title: "Show and edit a record's fields in the drawer", summary: "A typed table" },
    { id: "261029-kb33", title: "Workbench form controls", summary: "A custom select and a parent field helper" },
    { id: "261101-4n0v", title: "Spike animations", summary: "" },
  ];
  const ids = (q) => P.searchItems(items, q).map((r) => r.item.id);
  assert.deepEqual(P.tokens("  Foo  bar "), ["foo", "bar"]);
  assert.deepEqual(P.matchRanges("Parent parent", ["parent"]), [[0, 6], [7, 13]]);
  assert.deepEqual(P.matchRanges("abcd", ["ab", "bc"]), [[0, 3]]); // overlapping ranges merge
  assert.deepEqual(P.matchRanges("abc", []), []);
  assert.deepEqual(P.highlight("Hello world", [[0, 5]]), [{ text: "Hello", hit: true }, { text: " world", hit: false }]);
  assert.deepEqual(P.highlight("a<b>c", [[1, 4]]), [{ text: "a", hit: false }, { text: "<b>", hit: true }, { text: "c", hit: false }]); // text stays text
  assert.deepEqual(P.highlight("", []), []);
  // no query lists everything as given
  assert.deepEqual(ids(""), items.map((i) => i.id));
  // an id match, by a prefix first
  assert.deepEqual(ids("261028"), ["261028-esyt"]);
  assert.deepEqual(ids("esyt"), ["261028-esyt"]);
  assert.equal(P.searchItems(items, "esyt")[0].where.join(), "id");
  // every token must match somewhere: title and summary together
  assert.deepEqual(ids("workbench parent"), ["261029-bdt9", "261029-kb33"]);
  assert.deepEqual(ids("workbench datalist"), ["261029-bdt9"]);
  // a title match from a word start before a match inside a word; a summary-only match last
  assert.deepEqual(ids("pick"), ["261029-bdt9"]);
  assert.deepEqual(ids("parent"), ["261029-bdt9", "261029-kb33"]);
  assert.deepEqual(P.searchItems(items, "parent").map((r) => r.where.join()), ["title", "summary"]);
  assert.deepEqual(ids("nothing here"), []);
  // case does not matter; the limit applies
  assert.deepEqual(ids("TYPED"), ["261028-esyt"]);
  assert.equal(P.searchItems(items, "", 2).length, 2);
  assert.equal(P.countText(0), "No open work item matches");
  assert.equal(P.countText(1), "1 result");
  assert.equal(P.countText(7), "7 results");
  // stage 3: the chosen items of a multi-select
  assert.deepEqual(P.addChip(["a"], "b"), ["a", "b"]);
  assert.deepEqual(P.addChip(["a"], "a"), ["a"]); // once
  assert.deepEqual(P.removeChip(["a", "b"], "a"), ["b"]);
  assert.deepEqual(P.removeChip(["a"], "z"), ["a"]);
  assert.deepEqual(P.offered(items, ["261029-bdt9"], "261028-esyt").map((i) => i.id), ["261029-kb33", "261101-4n0v"]); // not chosen, not itself
  assert.deepEqual(P.offered(items, []).length, 4);
}

// the drawer's field table: value types, edit rules, edit requests
{
  assert.equal(F.valueType("2026-11-07"), "date");
  assert.equal(F.valueType("2026-11-07T20:59"), "date");
  assert.equal(F.valueType(3), "number");
  assert.equal(F.valueType(false), "flag");
  assert.equal(F.valueType(["a"]), "list");
  assert.equal(F.valueType("[[a note]]"), "link");
  assert.equal(F.valueType("261029-bdt9"), "id");
  assert.equal(F.valueType("feature"), "text");
  assert.equal(F.valueType("2026-11-07x"), "text");
  assert.equal(F.valueType("[[x]]", () => false), "text"); // the caller says what a wikilink is

  const st = (n, t, s) => F.editState(n, t, s).state;
  // an open work item
  for (const s of ["raw", "ready", "active", "blocked"]) {
    for (const f of ["summary", "blocked_by", "cycle", "keep"]) assert.equal(st(f, "work", s), "yes", f + " " + s);
  }
  assert.equal(st("kind", "work", "raw"), "yes");
  assert.equal(st("kind", "work", "ready"), "yes");
  assert.equal(st("kind", "work", "active"), "warn");
  assert.equal(st("kind", "work", "blocked"), "warn");
  assert.equal(st("priority", "work", "ready"), "yes");
  assert.equal(st("priority", "work", "active"), "warn");
  assert.equal(st("priority", "work", "blocked"), "yes");
  assert.equal(st("parent", "work", "active"), "warn");
  assert.equal(st("parent", "work", "raw"), "yes");
  assert.match(F.editState("kind", "work", "active").why, /rewrites/);
  // a closed item keeps only `keep`
  for (const f of ["summary", "kind", "priority", "parent", "blocked_by", "cycle"]) assert.equal(st(f, "work", "done"), "no", f);
  assert.equal(st("keep", "work", "dropped"), "yes");
  // the tool's own fields, and ADRs
  for (const f of ["id", "created", "updated", "closed", "status", "claimed_by", "claimed_at"]) assert.equal(st(f, "work", "ready"), "no", f);
  assert.match(F.editState("claimed_by", "work", "active").why, /claim/);
  assert.equal(st("summary", "adr", "accepted"), "no");
  assert.equal(st("origin", "work", "ready"), "no"); // no control

  // edit requests
  const ctx = { ids: ["261029-bdt9", "261028-esyt"], self: "261028-esyt" };
  const ch = (n, v) => F.buildChanges(n, v, ctx);
  assert.deepEqual(ch("summary", " a line "), { changes: { summary: "a line" } });
  assert.deepEqual(ch("summary", ""), { changes: { summary: null } });
  assert.match(ch("summary", "x".repeat(161)).error, /161 characters/);
  assert.deepEqual(ch("kind", "bug"), { changes: { kind: "bug" } });
  assert.match(ch("kind", "epic").error, /one of/);
  assert.deepEqual(ch("priority", "2"), { changes: { priority: 2 } });
  assert.deepEqual(ch("priority", "-1"), { changes: { priority: -1 } });
  assert.deepEqual(ch("priority", ""), { changes: { priority: null } });
  assert.match(ch("priority", "1.5").error, /whole number/);
  assert.match(ch("priority", "01").error, /whole number/);
  assert.deepEqual(ch("parent", "261029-bdt9"), { changes: { parent: "261029-bdt9" } });
  assert.deepEqual(ch("parent", ""), { changes: { parent: null } });
  assert.match(ch("parent", "261028-esyt").error, /own parent/);
  assert.match(ch("parent", "261000-zzzz").error, /No open work item/);
  assert.deepEqual(ch("blocked_by", ["261029-bdt9", "a few words"]), { changes: { blocked_by: ["261029-bdt9", "a few words"] } });
  assert.deepEqual(ch("blocked_by", []), { changes: { blocked_by: null } });
  assert.deepEqual(ch("cycle", "2026-w45"), { changes: { cycle: "2026-w45" } });
  assert.deepEqual(ch("cycle", ""), { changes: { cycle: null } });
  assert.match(ch("cycle", "-bad").error, /label/);
  assert.deepEqual(ch("keep", true), { changes: { keep: true } });
  assert.deepEqual(ch("keep", false), { changes: { keep: false } });
  assert.match(ch("status", "done").error, /cannot be edited/);

  // the rows: present fields in file order, no title, then the editable ones it lacks
  const rec = { type: "work", id: "x", fields: { id: "x", title: "T", status: "ready", kind: "bug", summary: null, blocked_by: [], created: "2026-11-07" } };
  assert.deepEqual(F.rowNames(rec), ["id", "status", "kind", "created", "summary", "priority", "parent", "blocked_by", "cycle", "keep"]);
  const done = { type: "work", id: "x", fields: { id: "x", title: "T", status: "done", kind: "bug" } };
  assert.deepEqual(F.rowNames(done), ["id", "status", "kind", "keep"]);
  const adr = { type: "adr", id: "a", fields: { id: "a", title: "T", status: "accepted" } };
  assert.deepEqual(F.rowNames(adr), ["id", "status"]);
  assert.equal(F.TYPE_LEGEND.length, 7);
}
