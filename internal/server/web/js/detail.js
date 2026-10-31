// The record drawer: the fields of one record, what points at it, and its
// sections read from the API one at a time. Body text is shown as text.
import { el, text } from "./util.js";
import { field, isId, childrenOf, obsidianUrl, wikilinkParts, wikilinkLabel } from "./logic.js";
import { show, findNote, describeError, state } from "./store.js";
import { claimControls } from "./flow.js";

// Sections that start collapsed (everything else starts open).
const COLLAPSED = new Set(["log"]);

// While refresh() draws the drawer again: what the reader had, by section name
// ({ open: Map name to bool, text: Map name to shown text }), so nothing jumps.
let keep = null;

/**
 * A wikilink value as a node. It is plain text at first; the API says which
 * note it names, and then it becomes a link that opens the note in Obsidian
 * (and, when the note is a record, a second link that opens it here), or is
 * marked unresolved and is not clickable.
 */
function wikilinkNode(w, byId) {
  const node = el("span", { class: "wikilink", text: wikilinkLabel(w), title: "[[" + w.target + "]]" });
  findNote(w.target)
    .then((hit) => {
      if (!node.isConnected) return;
      const url = hit ? obsidianUrl(state.config, hit.path) : "";
      if (!hit || !url) {
        node.classList.add("unresolved");
        node.title = "No note named " + w.target + " was found";
        return;
      }
      const link = el("a", { class: "wikilink", href: url, text: wikilinkLabel(w), title: "Open " + hit.path + " in Obsidian" });
      const nodes = [link];
      if (hit.id && byId.has(hit.id)) nodes.push(" ", el("a", { class: "rowlink wikilink-here", href: "#/record/" + encodeURIComponent(hit.id), title: "Open " + hit.id + " here", text: "here" }));
      node.replaceWith(...nodes);
    })
    .catch(() => {
      node.title = "The note could not be looked up";
    });
  return node;
}

/** A value as nodes: ids that are in the repository become links, wikilinks resolve through the API. */
export function valueNodes(v, byId) {
  const one = (x) => {
    const s = String(x);
    const w = wikilinkParts(s);
    if (w) return wikilinkNode(w, byId);
    return isId(s) && byId.has(s) ? el("a", { class: "rowlink", href: "#/record/" + encodeURIComponent(s), text: s }) : s;
  };
  if (Array.isArray(v)) return v.flatMap((x, i) => (i ? [", ", one(x)] : [one(x)]));
  return [one(v)];
}

/**
 * drawerEl: the <aside>; viewEl: the page behind it (made inert while open).
 * ctx: byId(), records(), close() leaves the record route.
 */
export function createDrawer(drawerEl, viewEl, ctx) {
  let current = null;
  let seq = 0;
  let returnFocus = null;
  const scrim = el("div", { class: "scrim drawer-scrim", hidden: true, onclick: () => ctx.close() });
  drawerEl.before(scrim);

  function sectionBlock(id, name, mySeq) {
    const body = el("pre", { class: "body" });
    const box = el("details", { class: "sec" }, el("summary", { text: name }), body);
    const shown = keep && keep.text.get(name);
    let loaded = false;
    const load = async () => {
      if (loaded) return;
      loaded = true;
      body.textContent = shown || "Loading…";
      try {
        const data = await show(id, name);
        if (mySeq === seq) body.textContent = data.section.text.replace(/^\n+|\s+$/g, "") || "(empty)";
      } catch (e) {
        loaded = false;
        if (mySeq === seq) body.textContent = describeError(e);
      }
    };
    box.addEventListener("toggle", () => box.open && load());
    if (keep && keep.open.has(name) ? keep.open.get(name) : !COLLAPSED.has(name)) {
      box.open = true;
      load();
    }
    return box;
  }

  function render(r, sections, mySeq) {
    const byId = ctx.byId();
    const dl = el("dl", { class: "fields" });
    for (const [k, v] of Object.entries(r.fields)) {
      if (v === null || v === undefined || k === "id" || k === "title" || (Array.isArray(v) && !v.length)) continue;
      dl.append(el("dt", { text: k }), el("dd", {}, valueNodes(v, byId)));
    }
    const kids = r.type === "work" ? childrenOf(ctx.records(), r.id) : [];
    const parts = [
      el(
        "header",
        { class: "dlg-head" },
        el(
          "div",
          {},
          el("p", { class: "eyebrow" }, r.type === "adr" ? "ADR " : "Work item ", el("span", { class: "mono", text: r.id })),
          el("h2", { id: "drawer-t", tabindex: "-1", text: text(field(r, "title")) || r.id }),
          el("p", { class: "sub mono", text: r.path }),
        ),
        obsidianUrl(state.config, r.path) ? el("a", { class: "btn quiet", id: "drawer-obsidian", href: obsidianUrl(state.config, r.path) }, "Open in Obsidian ", el("kbd", { text: "o" })) : null,
        el("button", { type: "button", class: "btn quiet", id: "drawer-close", onclick: () => ctx.close() }, "Close ", el("kbd", { text: "Esc" })),
      ),
      el("div", { class: "dlg-body" }, dl),
    ];
    const body = parts[1];
    const claimSlot = el("div", { class: "claim-slot" });
    body.prepend(claimSlot);
    claimControls(r, () => self && self.refresh()).then((n) => n && claimSlot.append(n)).catch(() => {});
    if (kids.length) {
      body.append(
        el("h3", { class: "sec-h", text: "Children (" + kids.length + ")" }),
        el(
          "ul",
          { class: "kids" },
          kids.map((c) => el("li", {}, el("a", { class: "rowlink", href: "#/record/" + encodeURIComponent(c.id), text: c.id }), " ", el("span", { class: "stat", text: text(field(c, "status")) }), " " + text(field(c, "title")))),
        ),
      );
    }
    for (const name of sections) body.append(sectionBlock(r.id, name, mySeq));
    if (!sections.length) body.append(el("p", { class: "muted", text: "This record has no ## sections." }));
    drawerEl.replaceChildren(...parts);
    drawerEl.setAttribute("aria-labelledby", "drawer-t");
  }

  let self = null;

  function reveal() {
    scrim.hidden = false;
    drawerEl.hidden = false;
    viewEl.inert = true;
  }

  self = {
    isOpen: () => current !== null,
    id: () => current,
    async open(id) {
      const mySeq = ++seq;
      if (current === null) returnFocus = document.activeElement;
      current = id;
      reveal();
      drawerEl.replaceChildren(el("p", { class: "muted pad", text: "Loading " + id + "…" }));
      try {
        const data = await show(id);
        if (mySeq !== seq) return;
        render(data.record, data.sections || [], mySeq);
        const h = drawerEl.querySelector("#drawer-t");
        if (h) h.focus({ preventScroll: true });
        drawerEl.scrollTop = 0;
      } catch (e) {
        if (mySeq !== seq) return;
        drawerEl.replaceChildren(
          el("header", { class: "dlg-head" }, el("h2", { id: "drawer-t", tabindex: "-1", text: id }), el("button", { type: "button", class: "btn quiet", onclick: () => ctx.close() }, "Close ", el("kbd", { text: "Esc" }))),
          el("div", { class: "dlg-body" }, el("div", { class: "err", role: "alert", text: describeError(e) })),
        );
        drawerEl.querySelector("#drawer-t").focus({ preventScroll: true });
      }
    },
    /** Draw the open record again from the file, keeping scroll, open sections and focus. For the live board. */
    async refresh() {
      if (current === null) return;
      const id = current;
      const mySeq = ++seq;
      const top = drawerEl.scrollTop;
      const secs = [...drawerEl.querySelectorAll("details.sec")];
      const name = (d) => d.querySelector("summary").textContent;
      const a = document.activeElement;
      const focusKey = a && drawerEl.contains(a) ? (a.matches("summary") ? "sec:" + name(a.parentNode) : a.id) : "";
      let data;
      try {
        data = await show(id);
      } catch (e) {
        if (e && e.code === "not_found" && mySeq === seq && current === id) {
          drawerEl.replaceChildren(el("header", { class: "dlg-head" }, el("h2", { id: "drawer-t", tabindex: "-1", text: id }), el("button", { type: "button", class: "btn quiet", id: "drawer-close", onclick: () => ctx.close() }, "Close ", el("kbd", { text: "Esc" }))), el("p", { class: "muted pad", text: "This record was removed." }));
        }
        return; // otherwise keep what is shown; the next change tries again
      }
      if (mySeq !== seq || current !== id) return;
      keep = { open: new Map(secs.map((d) => [name(d), d.open])), text: new Map(secs.map((d) => [name(d), d.querySelector(".body").textContent])) };
      try {
        render(data.record, data.sections || [], mySeq);
      } finally {
        keep = null;
      }
      drawerEl.scrollTop = top;
      const target = focusKey.startsWith("sec:") ? [...drawerEl.querySelectorAll("details.sec")].find((d) => "sec:" + name(d) === focusKey)?.querySelector("summary") : focusKey && drawerEl.querySelector("#" + CSS.escape(focusKey));
      if (target) target.focus({ preventScroll: true });
    },
    /** Hide the drawer and give focus back to where it came from. */
    close() {
      if (current === null) return;
      seq++;
      current = null;
      drawerEl.hidden = true;
      scrim.hidden = true;
      viewEl.inert = false;
      drawerEl.replaceChildren();
      const back = returnFocus;
      returnFocus = null;
      return back;
    },
    /** Keys while the drawer is open: Esc closes, j and k scroll. */
    keys(e) {
      if (e.key === "Escape") {
        ctx.close();
        return true;
      }
      if (e.key === "j" || e.key === "ArrowDown") {
        drawerEl.scrollBy(0, 80);
        return true;
      }
      if (e.key === "k" || e.key === "ArrowUp") {
        drawerEl.scrollBy(0, -80);
        return true;
      }
      return false;
    },
  };
  return self;
}
