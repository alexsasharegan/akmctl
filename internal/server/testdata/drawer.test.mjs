// The record drawer: the field editors and the claim bar, driven by clicks on
// the fake document (run by TestDOMModulesUnderNode).
import assert from "node:assert/strict";
import { startApp, settle } from "./apphelper.mjs";

const reply = (record, extra = {}) => ({ ok: true, status: 200, statusText: "OK", json: async () => ({ ok: true, record, ...extra }) });
const app = await startApp({ claims: true, identity: "me", onWrite: (u, init, records) => reply({ ...records[1], etag: "w1" }) });
const { doc } = app;
const drawer = doc.getElementById("drawer");
const toasts = doc.getElementById("toasts");

location.hash = "#/record/261101-bbbb";
await settle();
assert.match(drawer.text, /Claim as me/, "a ready item can be claimed");
assert.match(doc.getElementById("hintbar").text, /RECORD.*ggtopGbottom/, "the hint bar lists gg and G while the drawer is open");

// edit a field: Edit, change, Save
const rowOf = (name) => drawer.querySelector('tr[data-field="' + name + '"]');
const editBtn = (name) => rowOf(name).querySelector("button");
editBtn("priority").dispatch("click");
const input = rowOf("priority").querySelector("input");
assert.ok(input, "Edit swaps the value for a control");
input.value = "3";
input.dispatch("input");
rowOf("priority").querySelector("button.primary").dispatch("click");
await settle(10);
let w = app.calls.at(-1);
assert.equal(w.method, "PATCH");
assert.deepEqual(w.body.changes, { priority: 3 });
assert.match(w.headers["If-Match"], /"e2"/);

// Cancel puts the value back and writes nothing
await settle();
const n = app.calls.length;
editBtn("priority").dispatch("click");
rowOf("priority").querySelectorAll("button").find((b) => b.text === "Cancel").dispatch("click");
assert.equal(rowOf("priority").querySelector("input"), null);
assert.equal(app.calls.length, n);

// a refused value is told in the row and nothing is written
editBtn("priority").dispatch("click");
rowOf("priority").querySelector("input").value = "lots";
rowOf("priority").querySelector("button.primary").dispatch("click");
await settle(3);
assert.ok(rowOf("priority").querySelector(".fm-msg.err"), "a bad number is refused in the row");
assert.equal(app.calls.length, n);

// a select: open the list, pick an option, save
editBtn("kind").dispatch("click");
const sel = rowOf("kind").querySelector(".sel-btn");
assert.ok(sel, "kind is a select");
sel.dispatch("click");
assert.equal(rowOf("kind").querySelector(".sel-list").hidden, false, "the list opens");
const opt = rowOf("kind").querySelectorAll(".sel-opt").find((o) => o.text === "bug");
assert.ok(opt, "the kinds are offered");
opt.dispatch("click");
assert.match(sel.text, /bug/);
rowOf("kind").querySelector("button.primary").dispatch("click");
await settle(10);
assert.deepEqual(app.calls.at(-1).body.changes, { kind: "bug" });

// a summary: one line, a counter, saved as typed
editBtn("summary").dispatch("click");
const sum = rowOf("summary").querySelector("textarea");
assert.ok(sum, "summary is a text area");
sum.value = "A short summary";
sum.dispatch("input");
rowOf("summary").querySelector("button.primary").dispatch("click");
await settle(10);
assert.deepEqual(app.calls.at(-1).body.changes, { summary: "A short summary" });

// a picker: type to find an item, choose it
editBtn("parent").dispatch("click");
const pick = rowOf("parent").querySelector("input");
assert.ok(pick, "parent is a picker");
pick.value = "aaaa";
pick.dispatch("input");
await settle(3);
const hit = rowOf("parent").querySelectorAll("li").find((li) => /261101-aaaa/.test(li.text));
assert.ok(hit, "the picker lists the match");
hit.dispatch("click");
rowOf("parent").querySelector("button.primary").dispatch("click");
await settle(10);
assert.deepEqual(app.calls.at(-1).body.changes, { parent: "261101-aaaa" });

// the id is locked
assert.match(rowOf("id").text, /locked/);

// claim: the button sends the claim as the configured name
const claim = drawer.querySelectorAll("button").find((b) => /Claim as me/.test(b.text));
claim.dispatch("click");
await settle(10);
w = app.calls.at(-1);
assert.equal(w.method, "POST");
assert.match(w.url, /records\/261101-bbbb\/claim$/);
assert.equal(w.body.as, "me");
assert.match(toasts.text, /Claimed 261101-bbbb/);

// gg and G scroll the drawer to the top and the bottom (the board's keys and 900 ms window)
{
  const tall = () => {
    drawer.scrollHeight = 2000;
    drawer.clientHeight = 500;
    drawer.scrollTop = 700;
  };
  const key = (k, props) => app.press(k, { target: doc.body, ...props });
  const realNow = Date.now;
  let t = 1_000_000;
  Date.now = () => t;
  try {
    tall();
    key("G");
    assert.equal(drawer.scrollTop, 1500, "G scrolls to the bottom");

    tall();
    key("g");
    assert.equal(drawer.scrollTop, 700, "one g does nothing");
    t += 899;
    key("g");
    assert.equal(drawer.scrollTop, 0, "gg within 900 ms scrolls to the top");

    tall();
    t += 5000;
    key("g");
    t += 900;
    key("g");
    assert.equal(drawer.scrollTop, 700, "a second g after 900 ms is a first g");
    t += 100;
    key("g");
    assert.equal(drawer.scrollTop, 0, "and the one after it completes gg");

    tall();
    t += 5000;
    key("g");
    key("x");
    key("g");
    assert.equal(drawer.scrollTop, 700, "another key between the two g presses cancels the first");

    tall();
    t += 5000;
    key("g");
    key("g", { ctrlKey: true });
    key("g", { metaKey: true });
    key("g", { altKey: true });
    assert.equal(drawer.scrollTop, 700, "g with ctrl, meta or alt is ignored");
    key("G", { ctrlKey: true });
    key("G", { metaKey: true });
    key("G", { altKey: true });
    assert.equal(drawer.scrollTop, 700, "G with ctrl, meta or alt is ignored");
  } finally {
    Date.now = realNow;
  }
}

// the help lists both keys in its Record group
app.press("?");
assert.match(doc.getElementById("overlay-root").text, /scroll to the top[\s\S]*scroll to the bottom/, "the help lists gg and G for the record");

console.log("drawer ok");
process.exit(0);
