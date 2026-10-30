// The five tables. Rows are focusable; j and k move, Enter opens, h and l
// switch table. The "Recently closed" table asks the API for done and dropped
// work updated in the last 30 days.
import { el, ageText, text } from "./util.js";
import { TABLES, readyActive, openByParent, blockedRows, staleRows, closedRows, field, priorityOf, blockedBy, staleness } from "./logic.js";
import { describeError } from "./store.js";
import { rowNav } from "./rows.js";

export const CLOSED_DAYS = 30;

const statusCell = (r) => el("span", { class: "stat st-" + field(r, "status"), text: text(field(r, "status")) });
const idCell = (r) => el("a", { class: "rowlink", href: "#/record/" + encodeURIComponent(r.id), text: r.id });
const titleCell = (r) => text(field(r, "title")) || "(no title)";
const priCell = (r) => (priorityOf(r) > 0 ? "P" + priorityOf(r) : "");
const ageCell = (r, today) => {
  const { days, stale } = staleness(r, today);
  return days === null ? "" : el("span", { class: stale ? "age stale" : "age", text: (stale ? "stale " : "") + ageText(days) });
};

function columns(slug, today) {
  const id = { head: "ID", cls: "c-id", cell: idCell };
  const status = { head: "Status", cell: statusCell };
  const pri = { head: "Pri", cls: "num", cell: priCell };
  const title = { head: "Title", cls: "c-name", cell: titleCell };
  const kind = { head: "Kind", cell: (r) => text(field(r, "kind")) };
  const updated = { head: "Updated", cell: (r) => text(field(r, "updated")) };
  const age = { head: "Age", cell: (r) => ageCell(r, today) };
  const waits = { head: "Waits on", cls: "c-detail", cell: (r) => blockedBy(r).join(", ") };
  switch (slug) {
    case "ready-active":
      return [id, status, pri, title, kind, { head: "Parent", cell: (r) => text(field(r, "parent")) }, updated, age];
    case "by-parent":
      return [id, status, pri, title, kind, updated, age];
    case "blocked":
      return [id, status, pri, title, waits, updated, age];
    case "stale":
      return [id, status, title, updated, age];
    case "closed":
      return [id, status, title, kind, { head: "Closed", cell: (r) => text(field(r, "updated")) }];
  }
  return [];
}

function rowsFor(slug, records, closed, today) {
  switch (slug) {
    case "ready-active":
      return readyActive(records).map((r) => ({ r }));
    case "by-parent":
      return openByParent(records).flatMap((g) => [{ group: g }, ...g.records.map((r) => ({ r }))]);
    case "blocked":
      return blockedRows(records).map((r) => ({ r }));
    case "stale":
      return staleRows(records, today).map((x) => ({ r: x.record }));
    case "closed":
      return closedRows(closed || []).map((r) => ({ r }));
  }
  return [];
}

const EMPTY = {
  "ready-active": "Nothing is ready or active.",
  "by-parent": "No open work.",
  blocked: "Nothing is blocked.",
  stale: "No open work.",
  closed: "Nothing was closed in the last " + CLOSED_DAYS + " days.",
};

function buildTable(slug, rows, today) {
  const cols = columns(slug, today);
  const body = el("tbody");
  for (const row of rows) {
    if (row.group) {
      const g = row.group;
      body.append(
        el(
          "tr",
          { class: "group" },
          el(
            "th",
            { colspan: String(cols.length), scope: "colgroup" },
            g.parent ? (g.known ? el("a", { class: "rowlink", href: "#/record/" + encodeURIComponent(g.parent), text: g.parent }) : g.parent) : null,
            g.parent ? " " : null,
            g.title,
            el("span", { class: "count", text: " (" + g.records.length + ")" }),
          ),
        ),
      );
      continue;
    }
    body.append(el("tr", { class: "row", tabindex: "-1", "data-id": row.r.id }, cols.map((c) => el("td", { class: c.cls || "" }, c.cell(row.r)))));
  }
  return el("div", { class: "table-wrap" }, el("table", { class: "grid" }, el("thead", {}, el("tr", {}, cols.map((c) => el("th", { class: c.cls || "", scope: "col", text: c.head })))), body));
}

/** ctx: records(), today(), loadClosed(), setPage(p), open(id), go(slug). Returns { keys, focus, destroy }. */
export function mountTables(host, ctx, slug) {
  const today = ctx.today();
  const records = ctx.records();
  let alive = true;
  const tabs = el(
    "nav",
    { class: "tabs", "aria-label": "Tables" },
    TABLES.map((t) => el("a", { class: "tab", href: "#/tables/" + t.slug, "aria-current": t.slug === slug ? "page" : null }, t.title)),
  );
  const content = el("div", { class: "tables-body" });
  const wrap = el("div", { class: "tables" }, tabs, content);
  host.replaceChildren(wrap);
  const title = TABLES.find((t) => t.slug === slug).title;

  const nav = rowNav(wrap, { open: ctx.open, changed: () => ctx.setPage(describe()) });

  function describe() {
    const n = nav.count();
    return { mode: "TABLES", hints: HINTS, where: n && nav.index() >= 0 ? title + " \u00B7 row " + (nav.index() + 1) + " of " + n : title };
  }

  function show(closed) {
    const data = rowsFor(slug, records, closed, today);
    content.replaceChildren(data.some((d) => d.r) ? buildTable(slug, data, today) : el("p", { class: "empty-col", text: EMPTY[slug] }));
    nav.setRows([...content.querySelectorAll("tr.row")]);
    ctx.setPage(describe());
  }

  if (slug === "closed") {
    content.replaceChildren(el("p", { class: "empty-col", text: "Loading closed work…" }));
    ctx
      .loadClosed()
      .then((list) => alive && show(list))
      .catch((e) => {
        if (alive) content.replaceChildren(el("div", { class: "err", role: "alert", text: describeError(e) }));
      });
  } else show();
  ctx.setPage(describe());

  const idx = TABLES.findIndex((t) => t.slug === slug);
  return {
    keys(e) {
      const k = e.key;
      if (k === "h" || k === "ArrowLeft" || k === "l" || k === "ArrowRight") {
        const d = k === "h" || k === "ArrowLeft" ? -1 : 1;
        ctx.go(TABLES[(idx + d + TABLES.length) % TABLES.length].slug);
        return true;
      }
      return nav.keys(e);
    },
    focus: () => nav.focus(),
    destroy() {
      alive = false;
      host.replaceChildren();
    },
  };
}

const HINTS = [
  ["j k", "row"],
  ["h l", "table"],
  ["g g", "first"],
  ["G", "last"],
  ["Enter", "open"],
  ["/", "search"],
  ["?", "help"],
];
