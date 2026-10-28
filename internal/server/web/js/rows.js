// Keyboard navigation over table rows (tr.row), shared by the tables, the ADR
// list and the check view. One row has tabindex 0; j and k move it.

/**
 * wrap: the element that holds the rows (listeners go on it).
 * opts: open(id) for a row with a data-id, changed() after the position moves.
 */
export function rowNav(wrap, opts) {
  let rows = [];
  let at = -1;
  let pendingG = 0;

  const sync = () => rows.forEach((r, i) => (r.tabIndex = i === at ? 0 : -1));

  wrap.addEventListener("focusin", (e) => {
    const r = e.target.closest && e.target.closest("tr.row");
    const i = r ? rows.indexOf(r) : -1;
    if (i >= 0 && i !== at) {
      at = i;
      sync();
      opts.changed();
    }
  });
  wrap.addEventListener("click", (e) => {
    const r = e.target.closest && e.target.closest("tr.row");
    if (r && r.dataset.id && !e.target.closest("a")) opts.open(r.dataset.id);
  });

  const nav = {
    /** Call after the rows are drawn. */
    setRows(list) {
      rows = list;
      at = rows.length ? Math.min(Math.max(at, 0), rows.length - 1) : -1;
      sync();
    },
    index: () => at,
    count: () => rows.length,
    focus() {
      sync();
      if (rows[at]) {
        rows[at].focus({ preventScroll: true });
        rows[at].scrollIntoView({ block: "nearest" });
      }
      opts.changed();
    },
    /** Handle a key; true when used. */
    keys(e) {
      const k = e.key;
      if (k === "g") {
        const now = Date.now();
        if (now - pendingG < 900) {
          pendingG = 0;
          if (rows.length) {
            at = 0;
            nav.focus();
          }
        } else pendingG = now;
        return true;
      }
      pendingG = 0;
      if (!rows.length) return false;
      if (k === "j" || k === "ArrowDown") at = Math.min(rows.length - 1, at + 1);
      else if (k === "k" || k === "ArrowUp") at = Math.max(0, at - 1);
      else if (k === "G" || k === "End") at = rows.length - 1;
      else if (k === "Home") at = 0;
      else if (k === "o" || (k === "Enter" && !(e.target.closest && e.target.closest("a, button")))) {
        const id = rows[Math.max(at, 0)].dataset.id;
        if (id) opts.open(id);
        return true;
      } else return false;
      nav.focus();
      return true;
    },
  };
  return nav;
}
