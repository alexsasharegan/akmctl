// The workbench page end to end under a fake document: app.js reads the token
// from the address, loads the records from a fake API and draws the board
// (run by TestDOMModulesUnderNode).
import assert from "node:assert/strict";
import { installPage } from "./fakedom.mjs";

const doc = installPage();

const work = (id, status, f) => ({ id, status, title: "T " + id, summary: null, kind: null, parent: null, priority: null, created: "2026-10-01", updated: "2026-10-27", blocked_by: null, claimed_by: null, claimed_at: null, ...f });
const records = [
  { id: "261101-aaaa", type: "work", path: "docs/work/261101-aaaa.md", etag: "e1", fields: work("261101-aaaa", "raw") },
  { id: "261101-bbbb", type: "work", path: "docs/work/261101-bbbb.md", etag: "e2", fields: work("261101-bbbb", "ready", { priority: 2 }) },
  { id: "261101-cccc", type: "work", path: "docs/work/261101-cccc.md", etag: "e3", fields: work("261101-cccc", "active") },
  { id: "261029-0001", type: "adr", path: "docs/adr/261029-0001.md", etag: "e4", fields: { id: "261029-0001", title: "An ADR", status: "accepted", date: "2026-10-01" } },
];

const seen = [];
let auth = "";
globalThis.fetch = async (url, init = {}) => {
  seen.push(String(url));
  auth = (init.headers && init.headers.Authorization) || "";
  const u = new URL(url, "http://127.0.0.1");
  const ok = (body) => ({ ok: true, status: 200, statusText: "OK", json: async () => ({ ok: true, ...body }) });
  if (auth !== "Bearer tok123") return { ok: false, status: 401, statusText: "Unauthorized", json: async () => ({ ok: false, error: { code: "unauthorized", message: "no valid token" } }) };
  if (u.pathname === "/ui/config") return ok({ root: "/home/me/demo", identity: "" });
  if (u.pathname === "/api/v1/version") return ok({ work_dir: "docs/work", adr_dir: "docs/adr" });
  if (u.pathname === "/api/v1/records") {
    const type = u.searchParams.get("type");
    return ok({ records: records.filter((r) => !type || r.type === type) });
  }
  const one = u.pathname.match(/^\/api\/v1\/records\/([^/]+)$/);
  if (one) {
    const r = records.find((x) => x.id === decodeURIComponent(one[1]));
    if (!r) return { ok: false, status: 404, statusText: "Not Found", json: async () => ({ ok: false, error: { code: "not_found", message: "no such record" } }) };
    const body = "# " + r.fields.title + "\n\n## Goal\n\nDo <b>it</b>.\n\n## Next\n\nStart.\n\n## Log\n\n- 2026-10-27: filed.\n";
    return ok({ record: { ...r, fields: { ...r.fields }, body } });
  }
  return ok({});
};
const hist = [];
globalThis.history = { replaceState: (_s, _t, url) => hist.push(url) };
globalThis.location = { hash: "#token=tok123", pathname: "/", search: "" };
const listeners = {};
globalThis.window.addEventListener = (t, fn) => (listeners[t] = fn);

await import("./app.js");
for (let i = 0; i < 20; i++) await new Promise((r) => setTimeout(r, 5));

const host = doc.getElementById("view");
assert.deepEqual(hist[0], "/", "the token is stripped from the address bar");
assert.ok(seen.some((u) => u.startsWith("/api/v1/records?")), "the records were asked for");
assert.equal(host.querySelectorAll(".col").length, 4, "the board has four lanes");
assert.deepEqual(host.querySelectorAll(".col").map((c) => c.querySelectorAll(".card").length), [1, 1, 1, 0]);
assert.equal(doc.getElementById("repo").textContent, "demo");
assert.equal(doc.body.dataset.view, "board");
assert.ok(doc.title.includes("Board"), "the tab title names the view: " + doc.title);
assert.equal(doc.querySelector('#nav [data-nav="board"]').getAttribute("aria-current"), "page");

// a route change draws the ADR list
location.hash = "#/adrs";
listeners.hashchange();
assert.equal(doc.body.dataset.view, "adrs");
assert.equal(host.querySelectorAll("tr.row").length, 1);
assert.equal(doc.querySelector('#nav [data-nav="adrs"]').getAttribute("aria-current"), "page");
assert.equal(doc.querySelector('#nav [data-nav="board"]').getAttribute("aria-current"), null);

// a record route opens the drawer; the record text is shown as text
location.hash = "#/record/261101-bbbb";
listeners.hashchange();
for (let i = 0; i < 20; i++) await new Promise((r) => setTimeout(r, 5));
const drawer = doc.getElementById("drawer");
assert.ok(seen.some((u) => u.startsWith("/api/v1/records/261101-bbbb")), "the record was fetched");
assert.match(drawer.text, /261101-bbbb/);
assert.match(drawer.text, /locked: the id never changes/, "the fields table is drawn");
assert.equal(drawer.querySelectorAll("b").length, 0, "markup in a record is text");

console.log("app ok");
process.exit(0);
