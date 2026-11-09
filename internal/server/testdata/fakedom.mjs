// A small fake DOM, enough for the render functions of web/js (run by
// TestDOMModulesUnderNode). It implements only what those modules call; a
// missing member throws, so a new DOM call in a module shows up as a test
// failure that says which member to add.
/** A removed element that had focus loses it, as in a browser. */
const dropFocus = (n) => {
  const a = globalThis.document && document.activeElement;
  for (let x = a; x; x = x.parent) if (x === n) return void (document.activeElement = null);
};

class FNode {
  constructor() {
    this.parent = null;
    this.kids = [];
  }
  get childNodes() {
    return this.kids;
  }
  get children() {
    return this.kids.filter((k) => k instanceof FElement);
  }
  _adopt(k) {
    if (k.parent) k.parent.kids.splice(k.parent.kids.indexOf(k), 1);
    k.parent = this;
    return k;
  }
  append(...ns) {
    for (const n of ns) this.kids.push(this._adopt(typeof n === "string" ? new FText(n) : n));
  }
  prepend(...ns) {
    const add = ns.map((n) => this._adopt(typeof n === "string" ? new FText(n) : n));
    this.kids.unshift(...add);
  }
  insertBefore(n, ref) {
    this._adopt(n);
    const i = ref ? this.kids.indexOf(ref) : -1;
    if (i < 0) this.kids.push(n);
    else this.kids.splice(i, 0, n);
    return n;
  }
  before(...ns) {
    for (const n of ns) this.parent.insertBefore(typeof n === "string" ? new FText(n) : n, this);
  }
  after(...ns) {
    const next = this.parent.kids[this.parent.kids.indexOf(this) + 1] || null;
    for (const n of ns) this.parent.insertBefore(typeof n === "string" ? new FText(n) : n, next);
  }
  remove() {
    dropFocus(this);
    if (this.parent) this.parent.kids.splice(this.parent.kids.indexOf(this), 1);
    this.parent = null;
  }
  replaceChildren(...ns) {
    for (const k of this.kids) {
      dropFocus(k);
      k.parent = null;
    }
    this.kids = [];
    this.append(...ns);
  }
}

class FText extends FNode {
  constructor(data) {
    super();
    this.data = data;
  }
  get textContent() {
    return this.data;
  }
}

const parseSimple = (sel) => {
  const out = { tag: null, id: null, classes: [], attrs: [], notAttrs: [], checked: false };
  const re = /^([a-zA-Z][\w-]*)|#([\w-]+)|\.([\w-]+)|\[([\w-]+)(?:=["']?([^\]"']*)["']?)?\]|:not\(\[([\w-]+)\]\)|:(checked)/y;
  let m;
  let i = 0;
  while (i < sel.length && (re.lastIndex = i, (m = re.exec(sel)))) {
    if (m[1]) out.tag = m[1].toLowerCase();
    else if (m[2]) out.id = m[2];
    else if (m[3]) out.classes.push(m[3]);
    else if (m[4]) out.attrs.push([m[4], m[5]]);
    else if (m[6]) out.notAttrs.push(m[6]);
    else out.checked = true;
    i = re.lastIndex;
  }
  if (i !== sel.length) throw new Error("fakedom: selector not supported: " + sel);
  return out;
};

export class FElement extends FNode {
  constructor(tag) {
    super();
    this.tagName = tag.toUpperCase();
    this.attrs = new Map();
    this.listeners = new Map();
    this.dataset = {};
    this.style = { setProperty(k, v) { this[k] = v; }, removeProperty(k) { delete this[k]; } };
    this.hidden = false;
    this.focused = false;
  }
  get className() {
    return this.attrs.get("class") || "";
  }
  set className(v) {
    this.attrs.set("class", String(v));
  }
  get classList() {
    const self = this;
    const list = () => self.className.split(/\s+/).filter(Boolean);
    const api = {
      contains: (c) => list().includes(c),
      add: (...cs) => (self.className = [...new Set([...list(), ...cs])].join(" ")),
      remove: (...cs) => (self.className = list().filter((c) => !cs.includes(c)).join(" ")),
      toggle: (c, force) => {
        const on = force === undefined ? !list().includes(c) : force;
        on ? api.add(c) : api.remove(c);
        return on;
      },
    };
    return api;
  }
  get value() {
    if (this._value !== undefined) return this._value;
    return this.tagName === "TEXTAREA" ? this.textContent : this.attrs.get("value") || "";
  }
  set value(v) {
    this._value = String(v);
  }
  get id() {
    return this.attrs.get("id") || "";
  }
  get tabIndex() {
    return Number(this.attrs.get("tabindex") ?? -1);
  }
  set tabIndex(v) {
    this.attrs.set("tabindex", String(v));
  }
  get textContent() {
    return this.kids.map((k) => k.textContent).join("");
  }
  set textContent(v) {
    this.replaceChildren(String(v));
  }
  setAttribute(k, v) {
    this.attrs.set(k, String(v));
    if (k.startsWith("data-")) this.dataset[k.slice(5).replace(/-(\w)/g, (_, c) => c.toUpperCase())] = String(v);
  }
  getAttribute(k) {
    return this.attrs.has(k) ? this.attrs.get(k) : null;
  }
  removeAttribute(k) {
    this.attrs.delete(k);
  }
  hasAttribute(k) {
    return this.attrs.has(k);
  }
  addEventListener(type, fn) {
    if (!this.listeners.has(type)) this.listeners.set(type, []);
    this.listeners.get(type).push(fn);
  }
  removeEventListener(type, fn) {
    const l = this.listeners.get(type) || [];
    if (l.includes(fn)) l.splice(l.indexOf(fn), 1);
  }
  /** Fire an event at this node, then up through its parents (bubbling). */
  dispatch(type, props = {}) {
    const e = { type, target: this, key: "", shiftKey: false, defaultPrevented: false, preventDefault() { this.defaultPrevented = true; }, stopPropagation() {}, ...props };
    for (let n = this; n; n = n.parent) for (const fn of (n.listeners && n.listeners.get(type)) || []) fn(e);
    return e;
  }
  contains(n) {
    for (let x = n; x; x = x.parent) if (x === this) return true;
    return false;
  }
  focus() {
    this.focused = true;
    if (globalThis.document) document.activeElement = this;
  }
  blur() {
    this.focused = false;
  }
  scrollIntoView() {}
  getBoundingClientRect() {
    return { left: 0, top: 0, right: 0, bottom: 0, width: 0, height: 0 };
  }
  _matches(sel) {
    const s = parseSimple(sel.trim());
    if (s.tag && this.tagName.toLowerCase() !== s.tag) return false;
    if (s.id && this.id !== s.id) return false;
    if (!s.classes.every((c) => this.classList.contains(c))) return false;
    if (s.checked && !this.checked) return false;
    if (s.notAttrs.some((k) => this.attrs.has(k))) return false;
    return s.attrs.every(([k, v]) => this.attrs.has(k) && (v === undefined || this.attrs.get(k) === v));
  }
  /** Comma lists and descendant selectors ("#nav [data-nav]"); no child or sibling combinators. */
  matches(sel) {
    return sel.split(",").some((one) => {
      const parts = one.trim().split(/\s+/);
      if (!this._matches(parts.pop())) return false;
      let n = this.parent;
      while (parts.length) {
        while (n && !(n instanceof FElement && n._matches(parts.at(-1)))) n = n.parent;
        if (!n) return false;
        parts.pop();
        n = n.parent;
      }
      return true;
    });
  }
  closest(sel) {
    for (let n = this; n; n = n.parent) if (n instanceof FElement && n.matches(sel)) return n;
    return null;
  }
  querySelectorAll(sel) {
    const out = [];
    const walk = (n) => {
      for (const k of n.kids) {
        if (k instanceof FElement) {
          if (k.matches(sel)) out.push(k);
          walk(k);
        }
      }
    };
    walk(this);
    return out;
  }
  querySelector(sel) {
    return this.querySelectorAll(sel)[0] || null;
  }
  /** The visible text of the subtree, for assertions. */
  get text() {
    return this.textContent;
  }
}

/** A media query list that never changes. */
export const mql = (matches) => ({ matches, addEventListener() {}, removeEventListener() {} });

/** Install a fake document (and the storage the modules use) on globalThis. */
export function installDOM() {
  const html = new FElement("html");
  const body = new FElement("body");
  html.append(body);
  const doc = {
    documentElement: html,
    body,
    title: "",
    activeElement: null,
    createElement: (t) => new FElement(t),
    createTextNode: (d) => new FText(d),
    getElementById: (id) => html.querySelector("#" + id),
    querySelector: (s) => html.querySelector(s),
    querySelectorAll: (s) => html.querySelectorAll(s),
    listeners: new Map(),
    addEventListener(t, fn) {
      if (!this.listeners.has(t)) this.listeners.set(t, []);
      this.listeners.get(t).push(fn);
    },
    removeEventListener(t, fn) {
      const l = this.listeners.get(t) || [];
      if (l.includes(fn)) l.splice(l.indexOf(fn), 1);
    },
    /** Fire a document-level event (a key press). */
    dispatch(type, props = {}) {
      const e = { type, key: "", target: body, metaKey: false, ctrlKey: false, altKey: false, shiftKey: false, defaultPrevented: false, preventDefault() { this.defaultPrevented = true; }, stopPropagation() {}, ...props };
      for (const fn of [...(this.listeners.get(type) || [])]) fn(e);
      return e;
    },
  };
  const store = () => {
    const m = new Map();
    return { getItem: (k) => (m.has(k) ? m.get(k) : null), setItem: (k, v) => m.set(k, String(v)), removeItem: (k) => m.delete(k) };
  };
  globalThis.window = globalThis;
  globalThis.innerWidth = 1280;
  globalThis.innerHeight = 800;
  globalThis.CSS = { escape: (s) => String(s).replace(/[^\w-]/g, (c) => "\\" + c) };
  globalThis.document = doc;
  globalThis.localStorage = store();
  globalThis.sessionStorage = store();
  globalThis.matchMedia = () => mql(false);
  return doc;
}

/** A host element attached to the fake document. */
export function mountPoint(id = "view") {
  const h = new FElement("main");
  h.setAttribute("id", id);
  document.body.append(h);
  return h;
}

/** The page skeleton of index.html: the elements app.js looks up. */
export function installPage() {
  const doc = installDOM();
  const add = (parent, tag, attrs = {}) => {
    const n = doc.createElement(tag);
    for (const [k, v] of Object.entries(attrs)) n.setAttribute(k, v);
    parent.append(n);
    return n;
  };
  add(doc.body, "a", { class: "skip", href: "#view" });
  const app = add(doc.body, "div", { class: "app" });
  const top = add(app, "header", { class: "top" });
  add(top, "span", { id: "repo", class: "repo" });
  const nav = add(top, "nav", { id: "nav" });
  for (const v of ["board", "tables", "adrs", "check"]) add(nav, "a", { class: "nav-item", href: "#/" + v, "data-nav": v });
  const actions = add(top, "div", { class: "top-actions" });
  add(actions, "button", { id: "palette-btn" });
  add(actions, "button", { id: "theme-btn" });
  add(app, "main", { id: "view", tabindex: "-1" });
  add(app, "footer", { id: "hintbar" });
  add(doc.body, "aside", { id: "drawer", class: "drawer" });
  add(doc.body, "div", { id: "overlay-root" });
  add(doc.body, "div", { id: "toasts" });
  return doc;
}
