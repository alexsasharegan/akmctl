// The workbench page end to end under a fake document: app.js reads the token
// from the address, loads the records from a fake API and draws the board
// (run by TestDOMModulesUnderNode).
import assert from "node:assert/strict";
import { startApp, settle } from "./apphelper.mjs";

const { doc, seen, hist, listeners } = await startApp();

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
assert.equal(doc.body.dataset.view, "adrs");
assert.equal(host.querySelectorAll("tr.row").length, 1);
assert.equal(doc.querySelector('#nav [data-nav="adrs"]').getAttribute("aria-current"), "page");
assert.equal(doc.querySelector('#nav [data-nav="board"]').getAttribute("aria-current"), null);

// a record route opens the drawer; the record text is shown as text
location.hash = "#/record/261101-bbbb";
await settle();
const drawer = doc.getElementById("drawer");
assert.ok(seen.some((u) => u.startsWith("/api/v1/records/261101-bbbb")), "the record was fetched");
assert.match(drawer.text, /261101-bbbb/);
assert.match(drawer.text, /locked: the id never changes/, "the fields table is drawn");
assert.match(drawer.text, /Do <b>it<\/b>\./, "the Goal section is drawn, its markup as text");
assert.match(drawer.text, /Start\./, "the Next section is drawn");
assert.match(drawer.text, /filed\./, "the Log is drawn");
assert.equal(drawer.querySelectorAll("b").length, 0, "markup in a record is text");

console.log("app ok");
process.exit(0);
