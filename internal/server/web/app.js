"use strict";
// The workbench page. Everything from the repository is shown with textContent
// only, never as HTML.
(function () {
  const params = new URLSearchParams(location.hash.slice(1));
  let token = params.get("token") || sessionStorage.getItem("akm-token") || "";
  if (params.get("token")) {
    sessionStorage.setItem("akm-token", token);
    history.replaceState(null, "", location.pathname);
  }
  const OPEN = ["raw", "ready", "active", "blocked"];
  const CLOSED = ["done", "dropped"];
  const view = document.getElementById("view");
  const drawer = document.getElementById("drawer");
  const toastEl = document.getElementById("toast");
  const showClosed = document.getElementById("showClosed");
  let tab = "board";

  function el(tag, props, ...kids) {
    const n = document.createElement(tag);
    for (const [k, v] of Object.entries(props || {})) {
      if (k === "class") n.className = v;
      else if (k.startsWith("on")) n.addEventListener(k.slice(2), v);
      else if (v !== undefined && v !== null) n.setAttribute(k, v);
    }
    for (const kid of kids.flat()) if (kid !== null && kid !== undefined) n.append(kid);
    return n;
  }

  let toastTimer;
  function toast(msg, bad) {
    toastEl.textContent = msg;
    toastEl.className = bad ? "bad" : "";
    toastEl.hidden = false;
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => { toastEl.hidden = true; }, bad ? 8000 : 2500);
  }

  async function api(method, path, body, etag) {
    const headers = { Authorization: "Bearer " + token };
    if (body) headers["Content-Type"] = "application/json";
    if (etag) headers["If-Match"] = '"' + etag + '"';
    const res = await fetch("/api/v1/" + path, { method, headers, body: body ? JSON.stringify(body) : undefined });
    const data = await res.json();
    if (!data.ok) {
      const e = data.error || {};
      const why = e.details && (e.details.reason || e.details.rule) ? " (" + [e.details.reason, e.details.rule].filter(Boolean).join(", ") + ")" : "";
      throw new Error(e.code + ": " + e.message + why);
    }
    return data;
  }

  const f = (r, k) => (r.fields[k] === null || r.fields[k] === undefined ? "" : r.fields[k]);

  function card(r) {
    const meta = el("div", { class: "m" }, el("span", {}, r.id));
    if (f(r, "priority") !== "") meta.append(el("span", { class: "badge" }, "P" + f(r, "priority")));
    if (f(r, "parent")) meta.append(el("span", {}, "in " + f(r, "parent")));
    const by = r.fields.blocked_by;
    if (by && by.length) meta.append(el("span", {}, "waits on " + by.join(", ")));
    const c = el("div", { class: "card", draggable: "true", "data-id": r.id }, el("div", { class: "t" }, f(r, "title")), meta);
    if (f(r, "summary")) c.append(el("div", { class: "s" }, f(r, "summary")));
    c.addEventListener("click", () => openDrawer(r.id));
    c.addEventListener("dragstart", (e) => { e.dataTransfer.setData("text/plain", r.id); });
    return c;
  }

  async function board() {
    const statuses = OPEN.concat(showClosed.checked ? CLOSED : []);
    const data = await api("GET", "board?statuses=" + statuses.join(","));
    const cols = el("div", { class: "cols" });
    for (const col of data.columns) {
      const box = el("section", { class: "col", "data-status": col.status }, el("h2", {}, col.status + " (" + col.count + ")"));
      if (!col.records.length) box.append(el("div", { class: "empty" }, "nothing here"));
      for (const r of col.records) box.append(card(r));
      box.addEventListener("dragover", (e) => { e.preventDefault(); box.classList.add("over"); });
      box.addEventListener("dragleave", () => box.classList.remove("over"));
      box.addEventListener("drop", (e) => { e.preventDefault(); box.classList.remove("over"); move(e.dataTransfer.getData("text/plain"), col.status); });
      cols.append(box);
    }
    view.replaceChildren(cols);
  }

  async function move(id, status) {
    if (CLOSED.includes(status)) { toast("Open the item and close it there: closing needs an Outcome.", true); return; }
    try {
      const cur = await api("GET", "records/" + id);
      if (cur.record.fields.status === status) return;
      const changes = { status };
      if (status === "blocked" && !(cur.record.fields.blocked_by || []).length) {
        const why = prompt("Blocked on what? (a person, a role, or an id)");
        if (!why) return;
        changes.blocked_by = why.split(",").map((s) => s.trim()).filter(Boolean);
      }
      await api("PATCH", "records/" + id, { changes }, cur.record.etag);
      toast(id + " is now " + status);
      await render();
    } catch (e) { toast(e.message, true); }
  }

  async function adrs() {
    const data = await api("GET", "records?type=adr&sort=-id");
    const body = el("tbody");
    for (const r of data.records) {
      const row = el("tr", { class: "row" }, el("td", {}, r.id), el("td", {}, f(r, "title")), el("td", {}, f(r, "status")), el("td", {}, f(r, "date")), el("td", {}, f(r, "summary")));
      row.addEventListener("click", () => openDrawer(r.id));
      body.append(row);
    }
    view.replaceChildren(el("table", {}, el("thead", {}, el("tr", {}, ["id", "title", "status", "date", "summary"].map((h) => el("th", {}, h)))), body));
  }

  async function check() {
    const data = await api("GET", "check");
    const body = el("tbody");
    for (const x of data.findings) {
      body.append(el("tr", {}, el("td", { class: x.severity }, x.severity + " " + x.rule), el("td", {}, x.file ? x.file + (x.line ? ":" + x.line : "") : ""), el("td", {}, x.message)));
    }
    view.replaceChildren(
      el("p", {}, data.errors + " errors, " + data.warnings + " warnings. Manifest: " + data.reconcile.state + "."),
      data.findings.length ? el("table", {}, el("thead", {}, el("tr", {}, ["rule", "where", "finding"].map((h) => el("th", {}, h)))), body) : el("div", { class: "empty" }, "No findings."),
    );
  }

  async function openDrawer(id) {
    try {
      const data = await api("GET", "records/" + id);
      const r = data.record;
      const dl = el("dl");
      for (const [k, v] of Object.entries(r.fields)) {
        if (v === null) continue;
        dl.append(el("dt", {}, k), el("dd", {}, Array.isArray(v) ? v.join(", ") : String(v)));
      }
      const bar = el("div", { class: "bar" });
      if (r.type === "work") {
        const sel = el("select", {}, ["raw", "ready", "active", "blocked"].map((s) => el("option", { value: s, selected: s === r.fields.status ? "" : null }, s)));
        sel.addEventListener("change", () => move(r.id, sel.value).then(() => openDrawer(r.id)));
        bar.append(el("label", {}, "status ", sel));
        if (!["done", "dropped"].includes(r.fields.status)) {
          bar.append(
            el("button", { onclick: () => act("POST", "records/" + r.id + "/claim", { as: prompt("Claim as:", "workbench") }, r) }, "Claim"),
            el("button", { onclick: () => act("DELETE", "records/" + r.id + "/claim", null, r) }, "Release"),
            el("button", { onclick: () => act("POST", "records/" + r.id + "/close", { status: "done" }, r) }, "Close as done"),
            el("button", { onclick: () => act("POST", "records/" + r.id + "/close", { status: "dropped" }, r) }, "Drop"),
          );
        }
      }
      bar.append(el("button", { onclick: () => { drawer.hidden = true; } }, "Close panel"));
      const parts = [el("h2", {}, r.fields.title || r.id), el("div", { class: "m" }, r.path), bar, dl];
      for (const name of data.sections) {
        if (name === "log" && data.sections.length > 1) continue;
        const sec = await api("GET", "records/" + id + "?section=" + encodeURIComponent(name));
        parts.push(el("h3", {}, name), el("pre", {}, sec.section.text));
      }
      drawer.replaceChildren(...parts);
      drawer.hidden = false;
    } catch (e) { toast(e.message, true); }
  }

  async function act(method, path, body, r) {
    if (body && body.as === null) return;
    try {
      await api(method, path, body, r.etag);
      toast("Done");
      await render();
      await openDrawer(r.id);
    } catch (e) { toast(e.message, true); }
  }

  async function render() {
    if (!token) {
      view.replaceChildren(el("div", { class: "empty" }, "No token. Open the URL that akmctl ui printed."));
      return;
    }
    try {
      if (tab === "board") await board();
      else if (tab === "adrs") await adrs();
      else await check();
    } catch (e) { toast(e.message, true); }
  }

  document.getElementById("tabs").addEventListener("click", (e) => {
    const t = e.target.dataset && e.target.dataset.tab;
    if (!t) return;
    tab = t;
    for (const b of document.querySelectorAll("#tabs button")) b.classList.toggle("on", b.dataset.tab === t);
    render();
  });
  document.getElementById("reload").addEventListener("click", render);
  showClosed.addEventListener("change", render);
  render();
})();
