// The close wizard with a new ADR whose link fails (work item 261029-nkqv):
// the wizard says the ADR is not linked and Retry writes the rest, without
// making the ADR a second time (run by TestDOMModulesUnderNode).
import assert from "node:assert/strict";
import { startApp, settle } from "./apphelper.mjs";

const json = (status, body) => ({ ok: status < 400, status, statusText: String(status), json: async () => body });
let failSet = 1;
const app = await startApp({
  onWrite: (u, init, records) => {
    const body = init.body ? JSON.parse(init.body) : {};
    if (u.endsWith("/preview")) {
      const rec = records.find((x) => x.id === body.writes.find((w) => w.id).id);
      return json(200, { ok: true, writes: body.writes.map((w) => ({ ok: true, operation: w.op, record: rec, diff: { path: rec.path, new_path: rec.path, hunks: [] } })) });
    }
    if (init.method === "POST" && u.endsWith("/records")) return json(201, { ok: true, record: { id: "261116-adr1", type: "adr", path: "docs/adr/261116-adr1-pick-x.md", etag: "a0", fields: { title: body.title, status: "proposed" } } });
    const id = decodeURIComponent((u.match(/records\/([^/]+)/) || [])[1] || "");
    const r = records.find((x) => x.id === id);
    if (init.method === "PATCH" && failSet) {
      failSet--;
      return json(500, { ok: false, error: { code: "internal", message: "disk full" } });
    }
    const fields = { ...r.fields, ...(body.changes || {}) };
    if (u.endsWith("/close")) fields.status = body.status;
    return json(200, { ok: true, record: { ...r, fields, etag: "w" + app.calls.length } });
  },
});
const { doc, press } = app;
const overlay = doc.getElementById("overlay-root");
const type = (el, v) => {
  el.value = v;
  el.dispatch("input");
};

// the active card: d opens the close wizard
press("l");
press("l");
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

// the write: new succeeds, the set that links it fails
press("Enter", { target: overlay.querySelector(".diff") || overlay });
await settle(10);
const posts = () => app.calls.filter((c) => c.method === "POST" && /\/records$/.test(c.url)).length;
assert.equal(posts(), 1, "the ADR was made");
assert.match(overlay.text, /disk full/);
assert.match(overlay.text, /261116-adr1 was created, but this item does not link to it yet/, "the wizard says the ADR is not linked");
const retry = overlay.querySelector("[data-act=retry]");
assert.ok(retry, "and offers a retry");

// retry: the link and the close are written, the ADR is not made again
retry.dispatch("click");
await settle(15);
assert.equal(posts(), 1, "no second ADR");
const patches = app.calls.filter((c) => c.method === "PATCH");
assert.equal(patches.length, 2, "the link was tried twice");
assert.deepEqual(patches.at(-1).body.changes.decisions, ["261116-adr1"], "the retry links the ADR that exists");
assert.match(app.calls.at(-1).url, /\/close$/, "then the close");
assert.equal(overlay.kids.length, 0, "the wizard closes");

console.log("orphan ok");
process.exit(0);
