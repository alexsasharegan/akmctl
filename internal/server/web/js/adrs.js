// The ADR list: every decision record, newest first.
import { el, text } from "./util.js";
import { adrRows, field } from "./logic.js";
import { rowNav } from "./rows.js";

const HINTS = [
  ["j k", "row"],
  ["g g", "first"],
  ["G", "last"],
  ["Enter", "open"],
  ["/", "search"],
  ["?", "help"],
];

/** ctx: records(), setPage(p), open(id). Returns { keys, focus, destroy }. */
export function mountAdrs(host, ctx) {
  const adrs = adrRows(ctx.records());
  const wrap = el("div", { class: "adrs" });
  host.replaceChildren(wrap);
  const nav = rowNav(wrap, { open: ctx.open, changed: () => ctx.setPage(describe()) });
  const describe = () => ({ mode: "ADRS", hints: HINTS, where: nav.count() ? "row " + (nav.index() + 1) + " of " + nav.count() : "no ADRs" });

  if (!adrs.length) {
    wrap.append(el("h1", { class: "pg-title", text: "ADRs" }), el("p", { class: "empty-col", text: "This repository has no ADRs." }));
  } else {
    const head = ["ID", "Status", "Title", "Date", "Scope", "Summary"];
    wrap.append(
      el("h1", { class: "pg-title" }, "ADRs ", el("span", { class: "count", text: "(" + adrs.length + ")" })),
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
            adrs.map((r) =>
              el(
                "tr",
                { class: "row", tabindex: "-1", "data-id": r.id },
                el("td", { class: "c-id" }, el("a", { class: "rowlink", href: "#/record/" + encodeURIComponent(r.id), text: r.id })),
                el("td", {}, el("span", { class: "stat st-" + field(r, "status"), text: text(field(r, "status")) })),
                el("td", { class: "c-name", text: text(field(r, "title")) }),
                el("td", { text: text(field(r, "date")) }),
                el("td", { class: "c-detail", text: text(field(r, "scope")) }),
                el("td", { class: "c-detail", text: text(field(r, "summary")) }),
              ),
            ),
          ),
        ),
      ),
    );
  }
  nav.setRows([...wrap.querySelectorAll("tr.row")]);
  ctx.setPage(describe());
  return {
    keys: (e) => nav.keys(e),
    focus: () => nav.focus(),
    destroy: () => host.replaceChildren(),
  };
}
