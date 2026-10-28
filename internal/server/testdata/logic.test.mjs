// Checks of the pure functions in web/js (run by TestPureModulesUnderNode).
import assert from "node:assert/strict";
import { ageDays, daysBefore, ageText, isoDate } from "./util.js";
import * as L from "./logic.js";
import { parseToken, withQuery } from "./store.js";
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
console.log("ok");
