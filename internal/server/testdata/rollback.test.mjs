// The close wizard with a new ADR whose link fails (work items 261029-nkqv and
// 261029-xtqj): the move is one apply call, the server undoes the ADR it made,
// and the wizard says nothing was written and offers no retry of a half
// (run by TestDOMModulesUnderNode).
import assert from "node:assert/strict";
import { startApp, settle } from "./apphelper.mjs";

const json = (status, body) => ({ ok: status < 400, status, statusText: String(status), json: async () => body });
let fail = { write: 1, rolled_back: true };
const app = await startApp({
  onWrite: (u, init, records) => {
    const body = init.body ? JSON.parse(init.body) : {};
    const rec = records.find((x) => x.id === body.writes.find((w) => w.id).id);
    if (u.endsWith("/preview")) return json(200, { ok: true, writes: body.writes.map((w) => ({ ok: true, operation: w.op, record: rec, diff: { path: rec.path, new_path: rec.path, hunks: [] } })) });
    if (!u.endsWith("/apply")) throw new Error("a write went around apply: " + u);
    if (fail) {
      const details = fail;
      fail = null;
      return json(500, { ok: false, error: { code: "unreadable", message: "disk full", details } });
    }
    const adr = { id: "261116-adr1", type: "adr", path: "docs/adr/261116-adr1-pick-x.md", etag: "a0", fields: { title: "Pick X", status: "proposed" } };
    return json(200, { ok: true, writes: body.writes.map((w, n) => ({ ok: true, operation: w.op, record: w.op === "new" ? adr : { ...rec, etag: "w" + n, fields: { ...rec.fields, ...(w.op === "close" ? { status: w.status } : {}) } } })) });
  },
});
const { doc, press } = app;
const overlay = doc.getElementById("overlay-root");
const toasts = doc.getElementById("toasts");
const type = (el, v) => {
  el.value = v;
  el.dispatch("input");
};

async function toReview() {
  press("d");
  await settle(10);
  assert.match(overlay.text, /Close as done/);
  type(overlay.querySelector("textarea"), "It works.");
  press("Enter", { metaKey: true, target: overlay.querySelector("textarea") });
  await settle(3);
  const radio = overlay.querySelector("#wz-dec-new");
  radio.name = "wz-dec"; // the fake document does not reflect these attributes as properties
  radio.value = "new";
  radio.dispatch("change");
  await settle(3);
  type(overlay.querySelector("#wz-adr-title"), "Pick X");
  press("Enter", { target: overlay.querySelector("#wz-adr-title") });
  await settle(10);
  assert.match(overlay.text, /Review/);
}

// the active card: d opens the close wizard
press("l");
press("l");
await toReview();

// the write: one apply call; the link fails and the server undoes the ADR
press("Enter", { target: overlay.querySelector(".diff") || overlay });
await settle(10);
const applies = () => app.calls.filter((c) => /\/apply$/.test(c.url));
assert.equal(applies().length, 1, "the move is one call");
assert.deepEqual(applies()[0].body.writes.map((w) => w.op), ["new", "set", "close"]);
assert.deepEqual(applies()[0].body.writes[1].changes.decisions, ["$new"], "the server fills in the id of the ADR it makes");
assert.match(overlay.text, /disk full/);
assert.match(overlay.text, /Nothing was written/, "the wizard says the move was not made");
assert.doesNotMatch(overlay.text, /was created/, "and names no half-made ADR");
assert.equal(overlay.querySelector("[data-act=retry]"), null, "there is no half to retry");

// a rollback that fails names the files left behind
press("Escape");
await settle(3);
press("Escape");
await settle(3);
press("Escape");
await settle(3);
assert.equal(overlay.kids.length, 0, "the wizard is closed");
fail = { write: 2, rolled_back: false, rollback_errors: [{ path: "docs/adr/261116-adr1-pick-x.md", error: "permission denied" }] };
await toReview();
press("Enter", { target: overlay.querySelector(".diff") || overlay });
await settle(10);
assert.match(overlay.text, /could not be put back: docs\/adr\/261116-adr1-pick-x\.md/);
assert.doesNotMatch(overlay.text, /Nothing was written/);

// the same move again goes through whole
press("Escape");
await settle(3);
press("Escape");
await settle(3);
press("Escape");
await settle(3);
await toReview();
press("Enter", { target: overlay.querySelector(".diff") || overlay });
await settle(15);
assert.equal(applies().length, 3);
assert.equal(overlay.kids.length, 0, "the wizard closes");
assert.match(toasts.text, /Closed .* as done\. ADR 261116-adr1 is proposed/);

console.log("rollback ok");
process.exit(0);
