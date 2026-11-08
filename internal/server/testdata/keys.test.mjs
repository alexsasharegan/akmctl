// The chrome and the keyboard of the workbench: help, palette, theme, view
// keys, capture and the move wizard, driven by key presses on the fake document
// (run by TestDOMModulesUnderNode).
import assert from "node:assert/strict";
import { startApp, settle } from "./apphelper.mjs";

let staleNext = false;
const wrote = (u, init, records) => {
  const m = u.match(/\/records\/([^/]+)/);
  const r = records.find((x) => x.id === (m && decodeURIComponent(m[1])));
  const body = init.body ? JSON.parse(init.body) : {};
  if (init.method === "POST" && u.endsWith("/records")) return { ok: true, status: 201, statusText: "Created", json: async () => ({ ok: true, record: { id: "261109-new1", type: "work", path: "docs/work/261109-new1.md", etag: "n0", fields: { title: body.title, status: "raw" } } }) };
  if (staleNext) { staleNext = false; return { ok: false, status: 412, statusText: "Precondition Failed", json: async () => ({ ok: false, error: { code: "stale", message: "changed" } }) }; }
  const fields = { ...(r ? r.fields : {}), ...(body.changes || {}) };
  if (body.dry_run) return { ok: true, status: 200, statusText: "OK", json: async () => ({ ok: true, dry_run: true, record: r, diff: { path: r.path, new_path: r.path, hunks: [] } }) };
  if (u.endsWith("/close")) fields.status = body.status;
  return { ok: true, status: 200, statusText: "OK", json: async () => ({ ok: true, record: { ...(r || {}), fields, etag: "n" + Math.random().toString(36).slice(2, 6) } }) };
};
const app = await startApp({ onWrite: wrote });
const { doc, press } = app;
const view = doc.getElementById("view");
const overlay = doc.getElementById("overlay-root");

// help
press("?");
assert.match(overlay.text, /keyboard/i, "? opens the shortcut help");
press("Escape");
assert.equal(overlay.kids.length, 0, "Esc closes it");

// theme
const before = doc.documentElement.dataset.theme;
press("t");
assert.notEqual(doc.documentElement.dataset.theme, before, "t switches the theme");
assert.equal(localStorage.getItem("akm-theme"), doc.documentElement.dataset.theme);

// palette
press("/");
const input = overlay.querySelector("#pal-q");
assert.ok(input, "/ opens the palette");
assert.ok(overlay.querySelectorAll("li").length > 3, "it lists the actions");
input.value = "ready";
input.dispatch("input");
await settle(3);
assert.match(overlay.querySelector("#pal-list").text, /Ready and active/);
press("Enter");
await settle(3);
assert.equal(doc.body.dataset.view, "tables", "Enter runs the selected action");
assert.equal(overlay.kids.length, 0);

// view keys: 1 to 4 switch views
for (const [k, v] of [["1", "board"], ["3", "adrs"], ["4", "check"], ["2", "tables"]]) {
  press(k);
  await settle(3);
  assert.equal(doc.body.dataset.view, v, k + " goes to " + v);
}
press("1");
await settle(3);

// a safe move: L sends the focused card (261101-aaaa, raw) one lane right
const toasts = doc.getElementById("toasts");
assert.equal(view.querySelectorAll(".card").length, 3);
assert.equal(app.calls.length, 0);
press("L");
await settle(10);
assert.equal(app.calls.length, 1, "one write");
assert.equal(app.calls[0].method, "PATCH");
assert.match(app.calls[0].url, /records\/261101-aaaa/);
assert.deepEqual(app.calls[0].body.changes, { status: "ready" });
assert.match(app.calls[0].headers["If-Match"], /"e1"/, "the write carries the etag the page loaded");
assert.match(toasts.text, /Moved .* from raw to ready/);

// undo moves it back
press("u");
await settle(10);
assert.equal(app.calls.length, 2);
assert.deepEqual(app.calls[1].body.changes, { status: "raw" });
assert.match(toasts.text, /back to raw/);

// u with nothing to undo says so
press("u");
await settle(3);
assert.match(toasts.text, /Nothing to undo/);

// a gated move: d on the active card asks for an outcome, then reviews, then writes
press("l");
press("l");
const before2 = app.calls.length;
press("d");
await settle(10);
assert.match(overlay.text, /Close as done/, "d opens the close wizard");
assert.match(overlay.text, /What was the outcome/);
const area = overlay.querySelector("textarea");
assert.ok(area, "the outcome step has a text area");
press("Enter", { metaKey: true, target: area });
await settle(3);
assert.match(overlay.text, /Outcome/, "an empty outcome is refused and the step stays");
const area2 = overlay.querySelector("textarea"); // the step was drawn again
area2.value = "It works.";
area2.dispatch("input");
press("Enter", { metaKey: true, target: area2 });
await settle(3);
assert.match(overlay.text, /Decision|ADR/i, "next comes the decision step");
press("Enter", { target: overlay.querySelector("input:checked") || overlay });
await settle(10);
assert.match(overlay.text, /Review/, "then the review");
assert.equal(app.calls.length, before2 + 1, "one dry run so far");
assert.equal(app.calls.at(-1).body.dry_run, true, "the review is a dry run");
assert.equal(app.calls.at(-1).body.status, "done");
press("Enter", { target: overlay.querySelector(".diff") || overlay });
await settle(10);
assert.equal(app.calls.length, before2 + 2, "and one real write");
assert.equal(app.calls.at(-1).body.outcome, "It works.");
assert.equal(app.calls.at(-1).body.dry_run, undefined);
assert.equal(overlay.kids.length, 0, "the wizard closes");
assert.match(toasts.text, /Closed .* as done/);

// capture: c opens the form; a title makes a raw item
const callsBefore = app.calls.length;
press("c");
await settle(3);
assert.match(overlay.text, /Capture/, "c opens the capture form");
press("Enter", { metaKey: true, target: overlay.querySelector("input") });
await settle(3);
assert.match(overlay.text, /title/i, "an empty title is refused in the form");
assert.equal(app.calls.length, callsBefore, "nothing was written");
const title = overlay.querySelector("input");
title.value = "A new idea";
title.dispatch("input");
press("Enter", { metaKey: true, target: title });
await settle(10);
const made = app.calls.at(-1);
assert.equal(made.method, "POST");
assert.match(made.url, /\/records$/);
assert.equal(made.body.title, "A new idea");
assert.equal(overlay.kids.length, 0, "the form closes");
assert.match(toasts.text, /Captured A new idea as raw/);

// drag and drop: a card dropped on a lane is the same move as L
press("1");
await settle(5);
const dt = { data: {}, setData(t, v) { this.data[t] = v; }, effectAllowed: "", dropEffect: "" };
const card = view.querySelector('.card[data-id="261101-aaaa"]');
const lane = view.querySelector(".col.s-ready");
card.dispatch("dragstart", { dataTransfer: dt });
assert.ok(card.classList.contains("dragging"));
assert.equal(dt.data["text/plain"], "261101-aaaa");
lane.querySelector(".col-body").dispatch("dragover", { dataTransfer: dt });
assert.ok(lane.classList.contains("drop-target"), "the lane under the card is marked");
const n0 = app.calls.length;
lane.querySelector(".col-body").dispatch("drop", { dataTransfer: dt });
await settle(10);
assert.equal(app.calls.length, n0 + 1);
assert.deepEqual(app.calls.at(-1).body.changes, { status: "ready" });
assert.equal(view.querySelectorAll(".dragging, .drop-target").length, 0, "the marks are cleared");

// a write that finds the record changed writes nothing more and says so
press("Escape");
press("1");
await settle(5);
staleNext = true;
press("L");
await settle(10);
assert.match(toasts.text, /changed since it was loaded/);

console.log("keys ok");
process.exit(0);
