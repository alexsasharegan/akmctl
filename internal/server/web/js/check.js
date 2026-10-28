// The check view: the findings of `check`, errors first. Severity is a word
// and a symbol as well as a colour.
import { el, text } from "./util.js";
import { orderFindings, SEVERITY } from "./logic.js";
import { describeError } from "./store.js";
import { rowNav } from "./rows.js";

const HINTS = [
  ["j k", "row"],
  ["g g", "first"],
  ["G", "last"],
  ["Enter", "open record"],
  ["r", "run again"],
  ["?", "help"],
];

const SYMBOL = { E: "✖", W: "▲" };

/** ctx: check(), byId(), setPage(p), open(id). Returns { keys, focus, destroy }. */
export function mountCheck(host, ctx) {
  let alive = true;
  const wrap = el("div", { class: "check-view" });
  host.replaceChildren(wrap);
  const nav = rowNav(wrap, { open: ctx.open, changed: () => ctx.setPage(describe()) });
  let summary = "Running check…";
  const describe = () => ({ mode: "CHECK", hints: HINTS, where: nav.count() && nav.index() >= 0 ? "finding " + (nav.index() + 1) + " of " + nav.count() : summary });
  wrap.append(el("h1", { class: "pg-title", text: "Check" }), el("p", { class: "muted", text: summary }));
  ctx.setPage(describe());

  ctx
    .check()
    .then((data) => {
      if (!alive) return;
      const byId = ctx.byId();
      const found = orderFindings(data.findings || []);
      summary = (data.clean ? "Clean" : "Not clean") + ": " + data.errors + " errors, " + data.warnings + " warnings";
      const state = data.reconcile && data.reconcile.state;
      const parts = [
        el("h1", { class: "pg-title", text: "Check" }),
        el("p", { class: "check-sum" }, el("strong", { text: data.clean ? "Clean. " : "Not clean. " }), data.errors + " errors, " + data.warnings + " warnings. Manifest: " + (state || "unknown") + "."),
      ];
      const steps = data.reconcile && data.reconcile.steps;
      if (steps && steps.length) parts.push(el("ul", { class: "steps" }, steps.map((s) => el("li", { text: s }))));
      if (!found.length) parts.push(el("p", { class: "empty-col", text: "No findings." }));
      else {
        const head = ["Severity", "Rule", "Where", "Finding", "Record"];
        parts.push(
          el(
            "div",
            { class: "table-wrap" },
            el(
              "table",
              { class: "grid" },
              el("thead", {}, el("tr", {}, head.map((h) => el("th", { scope: "col", text: h })))),
              el(
                "tbody",
                {},
                found.map((x) =>
                  el(
                    "tr",
                    { class: "row", tabindex: "-1", "data-id": x.id && byId.has(x.id) ? x.id : null },
                    el("td", {}, el("span", { class: "sev sev-" + x.severity }, el("span", { "aria-hidden": "true", text: SYMBOL[x.severity] || "?" }), " " + (SEVERITY[x.severity] || x.severity))),
                    el("td", { class: "mono", text: text(x.rule) }),
                    el("td", { class: "mono", text: x.file ? x.file + (x.line ? ":" + x.line : "") : "" }),
                    el("td", { class: "c-name", text: text(x.message) }),
                    el("td", {}, x.id && byId.has(x.id) ? el("a", { class: "rowlink", href: "#/record/" + encodeURIComponent(x.id), text: x.id }) : text(x.id)),
                  ),
                ),
              ),
            ),
          ),
        );
      }
      wrap.replaceChildren(...parts);
      nav.setRows([...wrap.querySelectorAll("tr.row")]);
      ctx.setPage(describe());
    })
    .catch((e) => {
      if (alive) wrap.replaceChildren(el("h1", { class: "pg-title", text: "Check" }), el("div", { class: "err", role: "alert", text: describeError(e) }));
    });

  return {
    keys: (e) => nav.keys(e),
    focus: () => nav.focus(),
    destroy() {
      alive = false;
      host.replaceChildren();
    },
  };
}
