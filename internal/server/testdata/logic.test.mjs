// Checks of the pure functions in web/js (run by TestPureModulesUnderNode).
import assert from "node:assert/strict";
import { ageDays, daysBefore, ageText, isoDate } from "./util.js";
import * as L from "./logic.js";
import { parseToken, withQuery } from "./store.js";

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
console.log("ok");
