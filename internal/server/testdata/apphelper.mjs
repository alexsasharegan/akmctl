// The workbench page with a fake API behind it, for app.test.mjs and
// keys.test.mjs (run by TestDOMModulesUnderNode).
import { installPage } from "./fakedom.mjs";

export const sleep = (ms = 5) => new Promise((r) => setTimeout(r, ms));
/** Let pending promises and timers of the page settle. */
export const settle = async (n = 20) => {
  for (let i = 0; i < n; i++) await sleep();
};

export async function startApp({ hash = "#token=tok123", onWrite, claims = false, identity = "", bodyOf, tweak } = {}) {
  const doc = installPage();

  const work = (id, status, f) => ({ id, status, title: "T " + id, summary: null, kind: null, parent: null, priority: null, created: "2026-10-01", updated: "2026-10-27", blocked_by: null, claimed_by: null, claimed_at: null, ...f });
  const records = [
    { id: "261101-aaaa", type: "work", path: "docs/work/261101-aaaa.md", etag: "e1", fields: work("261101-aaaa", "raw") },
    { id: "261101-bbbb", type: "work", path: "docs/work/261101-bbbb.md", etag: "e2", fields: work("261101-bbbb", "ready", { priority: 2 }) },
    { id: "261101-cccc", type: "work", path: "docs/work/261101-cccc.md", etag: "e3", fields: work("261101-cccc", "active") },
    { id: "261029-0001", type: "adr", path: "docs/adr/261029-0001.md", etag: "e4", fields: { id: "261029-0001", title: "An ADR", status: "accepted", date: "2026-10-01" } },
  ];
  
  if (tweak) tweak(records);
  const seen = [];
  let auth = "";
  globalThis.fetch = async (url, init = {}) => {
    seen.push(String(url));
    auth = (init.headers && init.headers.Authorization) || "";
    const u = new URL(url, "http://127.0.0.1");
    const ok = (body) => ({ ok: true, status: 200, statusText: "OK", json: async () => ({ ok: true, ...body }) });
    if (auth !== "Bearer tok123") return { ok: false, status: 401, statusText: "Unauthorized", json: async () => ({ ok: false, error: { code: "unauthorized", message: "no valid token" } }) };
    if (u.pathname === "/ui/config") return ok({ root: "/home/me/demo", identity });
    if (u.pathname === "/api/v1/version") return ok({ work_dir: "docs/work", adr_dir: "docs/adr", repository: { extensions: claims ? ["claims"] : [] } });
    if (u.pathname === "/api/v1/records") {
      const type = u.searchParams.get("type");
      return ok({ records: records.filter((r) => !type || r.type === type) });
    }
    const one = u.pathname.match(/^\/api\/v1\/records\/([^/]+)$/);
    if (one) {
      const r = records.find((x) => x.id === decodeURIComponent(one[1]));
      if (!r) return { ok: false, status: 404, statusText: "Not Found", json: async () => ({ ok: false, error: { code: "not_found", message: "no such record" } }) };
      const body = bodyOf ? bodyOf(r) : "# " + r.fields.title + "\n\n## Goal\n\nDo <b>it</b>.\n\n## Next\n\nStart.\n\n## Log\n\n- 2026-10-27: filed.\n";
      const names = [...body.matchAll(/^## (.+)$/gm)].map((m) => m[1]);
      const want = u.searchParams.get("section");
      if (want) {
        const m = body.match(new RegExp("^## " + want + "\\n([\\s\\S]*?)(?=^## |$(?![\\s\\S]))", "m"));
        return ok({ record: { ...r, fields: { ...r.fields } }, section: { name: want, text: m ? m[1] : "" } });
      }
      return ok({ record: { ...r, fields: { ...r.fields } }, body, sections: names });
    }
    return ok({});
  };
  const hist = [];
  globalThis.history = { replaceState: (_s, _t, url) => hist.push(url) };
  let current = hash;
  globalThis.location = {
    get hash() {
      return current;
    },
    // a browser fires hashchange when the fragment changes
    set hash(v) {
      current = v;
      if (listeners.hashchange) listeners.hashchange();
    },
    pathname: "/",
    search: "",
  };
  const listeners = {};
  globalThis.window.addEventListener = (t, fn) => (listeners[t] = fn);
  const calls = [];
  const realFetch = globalThis.fetch;
  globalThis.fetch = async (url, init = {}) => {
    if (init.method && init.method !== "GET") {
      calls.push({ method: init.method, url: String(url), body: init.body ? JSON.parse(init.body) : undefined, headers: init.headers });
      if (onWrite) return onWrite(String(url), init, records);
    }
    return realFetch(url, init);
  };
  await import("./app.js");
  await settle();
  return { doc, seen, hist, listeners, records, calls, press: (key, props = {}) => doc.dispatch("keydown", { key, shiftKey: key.length === 1 && key !== key.toLowerCase() && key !== "?" && key !== "/", target: doc.activeElement || doc.body, ...props }) };
}
