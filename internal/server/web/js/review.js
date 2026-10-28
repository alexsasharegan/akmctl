// The diff of a dry run, as rows and as nodes. Added and removed lines are told
// apart by a sign, a word for screen readers and the border style, not by
// colour alone. All text goes in as text nodes.
import { el } from "./util.js";

const SIGN = { add: "+", rem: "−", ctx: " ", hunk: "" };
const WORD = { add: "added", rem: "removed", ctx: "unchanged", hunk: "hunk header" };

/**
 * Rows for the `diff` member of a dry-run response ({ path, new_path, hunks }).
 * Returns { path, newPath, renamed, rows, added, removed }, where a row is
 * { kind: hunk|ctx|add|rem, old, new, sign, word, text }.
 */
export function diffRows(diff) {
  const d = diff || {};
  const rows = [];
  let added = 0;
  let removed = 0;
  for (const h of d.hunks || []) {
    rows.push({ kind: "hunk", old: null, new: null, sign: SIGN.hunk, word: WORD.hunk, text: h.header });
    for (const l of h.lines || []) {
      const kind = l.op === "add" ? "add" : l.op === "rem" ? "rem" : "ctx";
      if (kind === "add") added++;
      if (kind === "rem") removed++;
      rows.push({ kind, old: l.old, new: l.new, sign: SIGN[kind], word: WORD[kind], text: l.text });
    }
  }
  return { path: d.path || "", newPath: d.new_path || d.path || "", renamed: !!d.new_path && d.new_path !== d.path, rows, added, removed };
}

/** The diff as a region: file header, tally, a key to the signs, then the lines. */
export function renderDiff(diff) {
  const d = diffRows(diff);
  const lines = d.rows.map((r) =>
    el(
      "div",
      { class: "dl " + r.kind },
      el("span", { class: "ln", "aria-hidden": "true", text: r.old === null || r.old === undefined ? "" : String(r.old) }),
      el("span", { class: "ln", "aria-hidden": "true", text: r.new === null || r.new === undefined ? "" : String(r.new) }),
      el("span", { class: "sg", "aria-hidden": "true", text: r.sign }),
      el("span", { class: "tx" }, el("span", { class: "sr", text: r.word + ": " }), r.text === "" ? " " : r.text),
    ),
  );
  const empty = !d.rows.length;
  return el(
    "section",
    { class: "diff-block" },
    el(
      "div",
      { class: "diff-top" },
      el("span", { class: "file mono" }, d.path, d.renamed ? [" ", el("span", { class: "g", "aria-hidden": "true", text: "→" }), el("span", { class: "sr", text: " renamed to " }), " " + d.newPath] : null),
      el("span", { class: "tally" }, el("span", { class: "plus", text: "+" + d.added + " added" }), " ", el("span", { class: "minus", text: "−" + d.removed + " removed" })),
    ),
    el(
      "div",
      { class: "diff-key", "aria-hidden": "true" },
      el("span", {}, el("i", { class: "k add", text: "+" }), " added (solid bar)"),
      el("span", {}, el("i", { class: "k rem", text: "−" }), " removed (dashed bar)"),
      el("span", {}, el("i", { class: "k ctx", text: "·" }), " unchanged"),
    ),
    empty
      ? el("p", { class: "muted", text: "This change leaves the file as it is." })
      : el("div", { class: "diff", role: "region", "aria-label": "Unified diff of " + d.path, tabindex: "0" }, lines),
  );
}
