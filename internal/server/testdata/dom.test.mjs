// Render checks of the workbench's DOM modules against a fake document (run by
// TestDOMModulesUnderNode). The records are a small repository in the shape the
// API returns.
import assert from "node:assert/strict";
import { installDOM, mountPoint, mql } from "./fakedom.mjs";

const doc = installDOM();

const work = (id, f) => ({ id, type: "work", path: "docs/work/" + id + ".md", etag: "x", fields: { id, title: "T " + id, status: "ready", parent: null, priority: null, blocked_by: null, updated: "2026-10-20", summary: null, kind: null, ...f } });
const adr = (id, f) => ({ id, type: "adr", path: "docs/adr/" + id + ".md", etag: "x", fields: { id, title: "A " + id, status: "accepted", date: "2026-10-01", ...f } });
const records = [
  work("261101-aaaa", { status: "raw" }),
  work("261101-bbbb", { status: "ready", priority: 1, summary: "the summary" }),
  work("261101-cccc", { status: "active", claimed_by: "claude-x", claimed_at: "2026-10-27T10:00" }),
  work("261101-dddd", { status: "blocked", blocked_by: ["261101-bbbb"], parent: "261101-aaaa" }),
  adr("261029-0001", {}),
  adr("261029-0002", { status: "proposed", title: "<b>not html</b>" }),
];
const byId = new Map(records.map((r) => [r.id, r]));
const pages = [];
const opened = [];
const ctx = { records: () => records, byId: () => byId, today: () => "2026-10-28", setPage: (p) => pages.push(p), open: (id) => opened.push(id) };
const key = (k, target) => ({ key: k, shiftKey: k.length === 1 && k !== k.toLowerCase(), target: target || doc.body });

// boot.js: the theme before first paint
{
  const boot = (await import("./js/boot.js"), doc.documentElement.dataset.theme);
  assert.equal(boot, "light", "no saved choice and no dark preference gives light");
  localStorage.setItem("akm-theme", "dark");
  await import("./js/boot.js?again");
  assert.equal(doc.documentElement.dataset.theme, "dark", "a saved choice wins");
  localStorage.removeItem("akm-theme");
  globalThis.matchMedia = () => mql(true);
  await import("./js/boot.js?dark");
  assert.equal(doc.documentElement.dataset.theme, "dark", "the system preference decides without a saved choice");
  globalThis.matchMedia = () => mql(false);
}

// the board
{
  const { mountBoard, cardNode } = await import("./js/board.js");
  const host = mountPoint();
  const b = mountBoard(host, ctx);
  const lanes = host.querySelectorAll(".col");
  assert.equal(lanes.length, 4);
  assert.deepEqual(lanes.map((c) => c.querySelectorAll(".card").length), [1, 1, 1, 1]);
  assert.match(lanes[0].querySelector(".col-head").text, /Raw/);
  const card = host.querySelector('.card[data-id="261101-bbbb"]');
  assert.ok(card, "a card for each work item");
  assert.equal(card.getAttribute("href"), "#/record/261101-bbbb");
  assert.match(card.text, /P1/);
  assert.match(card.text, /the summary/);
  assert.match(host.querySelector('.card[data-id="261101-dddd"]').text, /waits on 261101-bbbb/);
  assert.match(host.querySelector('.card[data-id="261101-dddd"]').text, /parent 261101-aaaa/);
  assert.match(host.querySelector('.card[data-id="261101-cccc"]').text, /held by claude-x/);
  assert.equal(host.querySelectorAll(".card").filter((c) => c.tabIndex === 0).length, 1, "exactly one card is reachable by Tab");
  assert.equal(pages.at(-1).mode, "BOARD");

  // keys move across lanes and open the focused card
  const first = b.current();
  assert.equal(first, "261101-aaaa");
  assert.equal(b.keys(key("l")), true);
  assert.equal(b.current(), "261101-bbbb");
  assert.equal(b.keys(key("j")), true);
  assert.equal(b.keys(key("o")), true);
  assert.equal(opened.at(-1), b.current());
  assert.equal(b.keys(key("z")), false, "an unknown key is not used");
  b.destroy();
  assert.equal(host.kids.length, 0);

  // an empty lane says so
  const empty = mountPoint("v2");
  mountBoard(empty, { ...ctx, records: () => [] });
  assert.match(empty.text, /Nothing in raw\./);

  // text from the repository is text, never markup
  const evil = cardNode(work("261101-eeee", { title: "<img src=x onerror=boom>" }), byId, "2026-10-28");
  assert.equal(evil.querySelector(".card-title").text, "<img src=x onerror=boom>");
  assert.equal(evil.querySelectorAll("img").length, 0);
}

// the ADR list
{
  const { mountAdrs } = await import("./js/adrs.js");
  const host = mountPoint("adrs");
  const v = mountAdrs(host, ctx);
  const rows = host.querySelectorAll("tr.row");
  assert.equal(rows.length, 2);
  assert.match(host.querySelector("h1").text, /ADRs \(2\)/);
  assert.equal(host.querySelector("td.c-name") && host.querySelectorAll("td.c-name").some((c) => c.text === "<b>not html</b>"), true, "a title with markup shows as text");
  assert.equal(host.querySelectorAll("b").length, 0);
  assert.equal(pages.at(-1).mode, "ADRS");
  v.focus();
  assert.equal(v.keys(key("j")), true);
  assert.equal(v.keys(key("Enter")), true);
  assert.ok(opened.at(-1).startsWith("261029-"));
  rows[0].dispatch("click");
  v.destroy();
  const none = mountPoint("adrs2");
  mountAdrs(none, { ...ctx, records: () => [] });
  assert.match(none.text, /no ADRs/);
}

// the check view
{
  const { mountCheck } = await import("./js/check.js");
  const finding = { severity: "E", rule: "WK1", file: "docs/work/x.md", line: 3, message: "bad", id: "261101-aaaa" };
  const host = mountPoint("chk");
  const v = mountCheck(host, { ...ctx, check: async () => ({ clean: false, errors: 1, warnings: 1, findings: [{ severity: "W", rule: "WK12", message: "w" }, finding], reconcile: { state: "ok", steps: [] } }) });
  assert.match(host.text, /Running check/);
  await new Promise((r) => setTimeout(r, 0));
  assert.match(host.querySelector(".check-sum").text, /Not clean\. 1 errors, 1 warnings\. Manifest: ok\./);
  const rows = host.querySelectorAll("tr.row");
  assert.equal(rows.length, 2);
  assert.match(rows[0].text, /WK1/, "errors come first");
  assert.equal(rows[0].dataset.id, "261101-aaaa");
  assert.equal(rows[0].querySelector("a.rowlink").getAttribute("href"), "#/record/261101-aaaa");
  v.destroy();

  const clean = mountPoint("chk2");
  mountCheck(clean, { ...ctx, check: async () => ({ clean: true, errors: 0, warnings: 0, findings: [] }) });
  await new Promise((r) => setTimeout(r, 0));
  assert.match(clean.text, /No findings\./);

  const bad = mountPoint("chk3");
  mountCheck(bad, { ...ctx, check: async () => { throw new Error("boom"); } });
  await new Promise((r) => setTimeout(r, 0));
  assert.equal(bad.querySelector(".err").getAttribute("role"), "alert");
  assert.match(bad.querySelector(".err").text, /boom/);
}

// the tables
{
  const { mountTables } = await import("./js/tables.js");
  const wentTo = [];
  const tctx = { ...ctx, go: (slug) => wentTo.push(slug), loadClosed: async () => [work("261101-zzzz", { status: "done", updated: "2026-10-27" })], closedNow: () => null };

  const host = mountPoint("t1");
  const v = mountTables(host, tctx, "ready-active");
  assert.equal(host.querySelectorAll("nav.tabs a.tab").length, 5);
  assert.equal(host.querySelectorAll("a.tab").filter((a) => a.getAttribute("aria-current") === "page").length, 1);
  const ids = host.querySelectorAll("tr.row").map((r) => r.dataset.id);
  assert.deepEqual(ids.sort(), ["261101-bbbb", "261101-cccc"], "ready and active work only");
  assert.match(host.querySelector("tr.row").text, /T 261101-/);
  assert.equal(pages.at(-1).mode, "TABLES");
  assert.equal(v.keys(key("l")), true);
  assert.deepEqual(wentTo, ["by-parent"], "l goes to the next table");
  assert.equal(v.keys(key("h")), true);
  assert.deepEqual(wentTo.at(-1), "closed", "h wraps to the last table");
  v.destroy();

  const blocked = mountPoint("t2");
  mountTables(blocked, tctx, "blocked");
  assert.deepEqual(blocked.querySelectorAll("tr.row").map((r) => r.dataset.id), ["261101-dddd"]);
  assert.match(blocked.querySelector("tr.row").text, /261101-bbbb/, "the Waits on column names the blocker");

  const parents = mountPoint("t3");
  mountTables(parents, tctx, "by-parent");
  assert.ok(parents.querySelectorAll("tr.group").length >= 1, "by-parent has group rows");

  const none = mountPoint("t4");
  mountTables(none, { ...tctx, records: () => [] }, "blocked");
  assert.match(none.text, /Nothing is blocked\./);

  const closed = mountPoint("t5");
  mountTables(closed, tctx, "closed");
  assert.match(closed.text, /Loading closed work/);
  await new Promise((r) => setTimeout(r, 0));
  assert.deepEqual(closed.querySelectorAll("tr.row").map((r) => r.dataset.id), ["261101-zzzz"]);

  const failed = mountPoint("t6");
  mountTables(failed, { ...tctx, loadClosed: async () => { throw new Error("offline"); } }, "closed");
  await new Promise((r) => setTimeout(r, 0));
  assert.match(failed.querySelector(".err").text, /offline/);
}

console.log("dom ok");
