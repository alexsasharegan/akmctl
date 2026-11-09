// The project filter in the header narrows the board and is remembered
// (run by TestDOMModulesUnderNode).
import assert from "node:assert/strict";
import { startApp, settle } from "./apphelper.mjs";

const app = await startApp({
  tweak: (records) => {
    // 261101-aaaa becomes a project with one item under it; 261101-cccc is outside every project
    records[1].fields.parent = "261101-aaaa";
  },
});
const { doc } = app;
const view = doc.getElementById("view");
const filter = doc.getElementById("project-filter");
assert.ok(filter, "the filter is in the header");
assert.equal(filter.parent.hidden, false, "it shows when the repository has a project");
const opts = filter.querySelectorAll("option");
assert.match(opts[0].text, /All projects/);
assert.ok(opts.some((o) => /T 261101-aaaa \(\d+\)/.test(o.text)), "the project is listed with its count");
assert.equal(view.querySelectorAll(".card").length, 3);

const proj = opts.find((o) => /T 261101-aaaa/.test(o.text));
filter.value = proj.getAttribute("value");
filter.dispatch("change");
await settle(5);
assert.match(app.hist.at(-1), /project=/, "the choice is in the address");
const ids = () => view.querySelectorAll(".card").map((c) => c.dataset.id).sort();
assert.ok(ids().includes("261101-bbbb") && !ids().includes("261101-cccc"), "the board shows the project only: " + ids());
assert.match(doc.getElementById("toasts").text, /Showing .*: \d+ items/);
assert.equal(localStorage.getItem("akm-project:/home/me/demo"), proj.getAttribute("value"), "remembered for the repository");

filter.value = "~";
filter.dispatch("change");
await settle(5);
assert.deepEqual(ids(), ["261101-cccc"], "No project shows the rest");

filter.value = "";
filter.dispatch("change");
await settle(5);
assert.match(doc.getElementById("toasts").text, /Showing all projects/);
assert.equal(localStorage.getItem("akm-project:/home/me/demo"), null);

console.log("project ok");
process.exit(0);
