// The move wizard beyond the close: a missing Goal and Next, the way back, a
// write that fails, and the identity step of a claim (run by TestDOMModulesUnderNode).
import assert from "node:assert/strict";
import { startApp, settle } from "./apphelper.mjs";

const json = (status, body) => ({ ok: status < 400, status, statusText: String(status), json: async () => body });
let failNext = null;
const app = await startApp({
  claims: true,
  // no sections at all: every move to ready or active has to ask
  bodyOf: (r) => "# " + r.fields.title + "\n",
  onWrite: (u, init, records) => {
    const body = init.body ? JSON.parse(init.body) : {};
    const id = decodeURIComponent((u.match(/records\/([^/]+)/) || [])[1] || "");
    const r = records.find((x) => x.id === id);
    if (failNext && !body.dry_run) {
      const f = failNext;
      failNext = null;
      return json(f.status, { ok: false, error: f.error });
    }
    if (body.dry_run) return json(200, { ok: true, dry_run: true, record: r, diff: { path: r.path, new_path: r.path, hunks: [{ old_start: 1, old_lines: 1, new_start: 1, new_lines: 2, lines: [{ op: " ", text: "# t" }, { op: "+", text: "## Goal" }] }] } });
    return json(200, { ok: true, record: { ...r, etag: "w" + app.calls.length, fields: { ...r.fields, ...(body.changes || {}) } } });
  },
});
const { doc, press } = app;
const overlay = doc.getElementById("overlay-root");
const toasts = doc.getElementById("toasts");
const type = (el, v) => {
  el.value = v;
  el.dispatch("input");
};

// L on the raw card (no Goal): the wizard asks for it, with a way back
press("L");
await settle(10);
assert.match(overlay.text, /Move to ready/);
assert.match(overlay.text, /What does done look like\?/);
assert.equal(app.calls.length, 0, "nothing is written while questions are open");
press("Escape");
await settle(3);
assert.equal(overlay.kids.length, 0, "Esc on the first step cancels");
assert.equal(app.calls.length, 0);

press("L");
await settle(10);
type(overlay.querySelector("textarea"), "Done means the tests pass.");
press("Enter", { metaKey: true, target: overlay.querySelector("textarea") });
await settle(10);
assert.match(overlay.text, /Review/);
assert.ok(overlay.querySelector(".diff") || /Goal/.test(overlay.text), "the review shows the diff of the write");
const dry = app.calls.filter((c) => c.body && c.body.dry_run);
assert.ok(dry.length >= 1, "the review is made of dry runs");
assert.equal(dry[0].method, "PUT");
assert.match(dry[0].url, /sections\/Goal$/);

// Esc goes back a step, not out
press("Escape", { target: overlay.querySelector(".diff") || overlay });
await settle(3);
assert.match(overlay.text, /What does done look like\?/, "Esc on the review goes back");
assert.equal(overlay.querySelector("textarea").value, "Done means the tests pass.", "the answer is kept");
press("Enter", { metaKey: true, target: overlay.querySelector("textarea") });
await settle(10);

// a failing write is shown in the wizard and the dialog stays open
failNext = { status: 500, error: { code: "internal", message: "disk full" } };
press("Enter", { target: overlay.querySelector(".diff") || overlay });
await settle(10);
assert.match(overlay.text, /disk full/, "the error is shown");
assert.equal(overlay.querySelector(".dialog.wizard") !== null, true, "the wizard stays open");
press("Escape");
await settle(3);
press("Escape");
await settle(3);
assert.equal(overlay.kids.length, 0);

// a stale write offers to reload
press("L");
await settle(10);
type(overlay.querySelector("textarea"), "Goal text.");
press("Enter", { metaKey: true, target: overlay.querySelector("textarea") });
await settle(10);
failNext = { status: 412, error: { code: "stale", message: "changed" } };
press("Enter", { target: overlay.querySelector(".diff") || overlay });
await settle(10);
assert.match(overlay.text, /changed|conflict|reload/i, "a conflict is told");
press("Escape");
await settle(3);
press("Escape");
await settle(3);

// claiming into active with no name asks for Goal, Next and a name
press("l"); // to 261101-bbbb in the ready lane
const n = app.calls.length;
press("L");
await settle(10);
assert.match(overlay.text, /Move to active/);
assert.match(overlay.text, /step 1 of 4/, "goal, next, identity and review");
type(overlay.querySelector("textarea"), "A goal.");
press("Enter", { metaKey: true, target: overlay.querySelector("textarea") });
await settle(3);
assert.match(overlay.text, /next step/i);
type(overlay.querySelector("textarea"), "Do the first thing.");
press("Enter", { metaKey: true, target: overlay.querySelector("textarea") });
await settle(3);
assert.match(overlay.text, /claim|name|identity/i, "the third question is who is claiming");
press("Enter", { target: overlay.querySelector("input") });
await settle(3);
assert.match(overlay.text, /Identity|name/i, "an empty name is refused");
type(overlay.querySelector("input"), "claude-test");
press("Enter", { target: overlay.querySelector("input") });
await settle(10);
assert.match(overlay.text, /Review/);
press("Enter", { target: overlay.querySelector(".diff") || overlay });
await settle(15);
const writes = app.calls.slice(n).filter((c) => !c.body || !c.body.dry_run).map((c) => c.method + " " + c.url.replace("/api/v1/records/261101-bbbb", ""));
assert.ok(writes.some((w) => /sections\/Goal/.test(w)) && writes.some((w) => /sections\/Next/.test(w)) && writes.some((w) => /\/claim$/.test(w)), "Goal, Next and the claim are written: " + writes.join(", "));
assert.equal(app.calls.findLast((c) => /\/claim$/.test(c.url)).body.as, "claude-test");
assert.equal(overlay.kids.length, 0);
assert.match(toasts.text, /Moved .* to active/);

console.log("wizard ok");
process.exit(0);
