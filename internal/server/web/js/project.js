// The project filter: a select in the header that narrows the board, the
// tables and the palette to one project. The choice lives in the URL hash
// (?project=) and, per repository, in localStorage.
import { el } from "./util.js";
import { projectCounts, NO_PROJECT } from "./logic.js";

const key = (root) => "akm-project:" + (root || "");

/** The remembered selection for a repository, or null. */
export function recall(root) {
  try {
    const v = localStorage.getItem(key(root));
    return v === null || v === "" ? null : v;
  } catch (e) {
    return null;
  }
}

function remember(root, sel) {
  try {
    if (sel === null) localStorage.removeItem(key(root));
    else localStorage.setItem(key(root), sel);
  } catch (e) {
    /* kept in memory only */
  }
}

/**
 * Put the control into `host`. ctx: records() all loaded records, root() the
 * repository root, get() the current selection, onChange(sel). Returns
 * { refresh() }: call it after the records or the selection changed.
 */
export function mountProjectFilter(host, ctx) {
  const select = el("select", { id: "project-filter", "aria-label": "Filter by project" });
  const wrap = el("label", { class: "project-filter", hidden: "true" }, el("span", { class: "project-label", text: "Project" }), select);
  host.prepend(wrap);
  select.addEventListener("change", () => {
    const v = select.value === "" ? null : select.value;
    remember(ctx.root(), v);
    ctx.onChange(v);
  });
  function refresh() {
    const { list, none } = projectCounts(ctx.records());
    // no project field in this repository: nothing to filter on
    wrap.hidden = list.length === 0;
    const sel = ctx.get();
    const opts = [el("option", { value: "", text: "All projects" })];
    for (const p of list) opts.push(el("option", { value: p.target, text: p.target + " (" + p.count + ")" }));
    if (none > 0) opts.push(el("option", { value: NO_PROJECT, text: "No project (" + none + ")" }));
    // a selection with no records any more stays selectable, so it can be cleared
    if (sel !== null && sel !== NO_PROJECT && !list.some((p) => p.target === sel)) opts.push(el("option", { value: sel, text: sel + " (0)" }));
    select.replaceChildren(...opts);
    select.value = sel === null ? "" : sel;
  }
  return { refresh };
}
