// The child-count chip on the board card and in the drawer summary line
// (run by TestDOMModulesUnderNode).
import assert from "node:assert/strict";
import { startApp, settle } from "./apphelper.mjs";

const work = (id, status, f) => ({ id, type: "work", path: "docs/work/" + id + ".md", etag: "e" + id, fields: { id, title: "T " + id, status, summary: null, kind: null, parent: null, priority: null, created: "2026-10-01", updated: "2026-10-27", blocked_by: null, claimed_by: null, claimed_at: null, ...f } });
const app = await startApp({
  tweak: (records) => {
    // aaaa: three done, one ready, one active and one dropped child: "3/5"
    // bbbb: a chore with one done child: "1/1"
    // hhhh: a parent whose only child is dropped: no chip
    Object.assign(records[1].fields, { kind: "chore", priority: 1 });
    records.push(
      work("261101-d001", "done", { parent: "261101-aaaa" }),
      work("261101-d002", "done", { parent: "261101-aaaa" }),
      work("261101-d003", "done", { parent: "261101-aaaa" }),
      work("261101-d004", "ready", { parent: "261101-aaaa" }),
      work("261101-d005", "active", { parent: "261101-aaaa" }),
      work("261101-d006", "dropped", { parent: "261101-aaaa" }),
      work("261101-d007", "done", { parent: "261101-bbbb" }),
      work("261101-hhhh", "ready"),
      work("261101-d008", "dropped", { parent: "261101-hhhh" }),
    );
  },
});
const { doc } = app;
const view = doc.getElementById("view");
const drawer = doc.getElementById("drawer");
const card = (id) => view.querySelector('.card[data-id="' + id + '"]');
const chip = (c) => c.querySelector(".kids-chip");
const shown = (c) => c.querySelector('[aria-hidden="true"]').text;

// the board card
{
  const c = chip(card("261101-aaaa"));
  assert.ok(c, "a parent's card shows the chip");
  assert.equal(shown(c), "↳ 3/5", "arrow, then done over total; the dropped child is in neither");
  assert.equal(c.getAttribute("title"), "3 of 5 child items done");
  assert.equal(c.querySelector(".sr-only").text, "3 of 5 child items done", "the accessible name is the same words");
  assert.doesNotMatch(c.getAttribute("class"), /wait/, "the chip is not dashed");
  assert.equal(card("261101-cccc").querySelectorAll(".kids-chip").length, 0, "a card with no children has no chip");
  assert.equal(card("261101-hhhh").querySelectorAll(".kids-chip").length, 0, "nor one whose children are all dropped");
}

// a chore parent shows it, between the priority chip and the kind chip, and the kind chip is as before
{
  const b = card("261101-bbbb");
  const c = chip(b);
  assert.ok(c, "the kind does not matter");
  assert.equal(shown(c), "↳ 1/1");
  assert.equal(c.getAttribute("title"), "1 of 1 child item done", "singular for one child");
  const classes = b.querySelector(".card-head").children.map((n) => n.getAttribute("class"));
  const at = (re) => classes.findIndex((k) => re.test(k));
  assert.ok(at(/chip pri/) >= 0 && at(/chip pri/) < at(/kids-chip/) && at(/kids-chip/) < at(/chip kind/), "after the priority chip, before the kind chip: " + classes.join(" | "));
  const kind = b.querySelector(".chip.kind");
  assert.equal(kind.text, "chore");
  assert.equal(kind.getAttribute("class"), "chip kind", "the kind chip is unchanged");
}

// the drawer summary line
{
  location.hash = "#/record/261101-aaaa";
  await settle();
  const sum = drawer.querySelector(".dlg-sum");
  const c = sum.querySelector(".kids-chip");
  assert.ok(c, "the summary line shows the chip");
  assert.equal(shown(c), "↳ 3/5");
  assert.equal(c.getAttribute("title"), "3 of 5 child items done");
  assert.equal(c.tagName, "BUTTON", "it is a button, so the keyboard operates it");
  const pri = sum.children.findIndex((n) => /priority/.test(n.text));
  assert.ok(pri < 0 || pri < sum.children.indexOf(c), "after the priority chip when there is one");
  const list = drawer.querySelector("#drawer-kids");
  assert.ok(list && /Children \(6\)/.test(list.text), "the Children list is there, with the dropped child in it");
  c.dispatch("click");
  assert.equal(doc.activeElement, list, "the chip moves focus to the Children list");

  location.hash = "#/record/261101-bbbb";
  await settle();
  const own = drawer.querySelector(".dlg-sum");
  const kids = own.children;
  assert.ok(kids.findIndex((n) => /priority\s*1/.test(n.text)) < kids.indexOf(own.querySelector(".kids-chip")), "after the priority chip");

  location.hash = "#/record/261101-cccc";
  await settle();
  assert.equal(drawer.querySelectorAll(".kids-chip").length, 0, "no children, no chip");
  location.hash = "#/record/261101-hhhh";
  await settle();
  assert.equal(drawer.querySelectorAll(".kids-chip").length, 0, "all children dropped, no chip");
}

console.log("kids ok");
process.exit(0);
