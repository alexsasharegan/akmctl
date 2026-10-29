// The move wizard: one question per step, only for what the record still
// lacks, then a review of the exact diff of every write before any is made.
// A write that finds the record changed shows the conflict and writes nothing more.
import { el } from "./util.js";
import { api, ApiError, describeError } from "./store.js";
import { buildOps, requestFor, resolveOp, previewable, opLabel, validateStep, wizardKey } from "./moves.js";
import { slugify } from "./capture.js";
import { renderDiff } from "./review.js";

const STEP_NAMES = { goal: "Goal", next: "Next step", blocked_by: "Waiting on", outcome: "Outcome", decision: "Decision", identity: "Identity", unblock: "Waiting on", review: "Review" };

const TEXT = {
  goal: { q: "What does done look like?", help: "One or two lines. It is written under ## Goal; say how you will know it is finished.", rows: 3, field: "goal" },
  next: { q: "What is the next step?", help: "The first thing to do when you pick this up. It is written under ## Next.", rows: 3, field: "next" },
  blocked_by: { q: "What is it waiting on?", help: "One entry per line: a record id, or a few words about a person or a thing outside the repository.", rows: 3, field: "blocked_by" },
};


/**
 * deps: { openDialog, closeDialog, setDialogHints, MOD }.
 * spec: { rec (with etag), title, i (move input for buildOps), plan, adrs: [{ id, title, status }],
 *         outcomeDefault, onDone(result), onReload() }.
 */
export function openWizard(deps, spec) {
  const { rec, i: move, plan } = spec;
  const to = move.to;
  const steps = plan.steps.concat("review");
  const W = {
    at: 0,
    err: "",
    phase: "edit", // edit | loading | review | writing | conflict | error
    a: { unblock: true, goal: "", next: "", blocked_by: "", outcome: spec.outcomeDefault || "", identity: spec.identity || "", decision: { mode: "none", id: "", title: "", slug: "" } },
    previews: [],
    takeover: false,
    fatal: null,
    seq: 0,
    written: 0,
    ops: [],
  };
  const root = el("div", { class: "dialog wizard", role: "dialog", "aria-modal": "true", "aria-labelledby": "wz-t" });
  const cur = () => steps[W.at];
  const MOD = deps.MOD;
  const ops = () => buildOps(move, W.a);

  // a textarea keeps Enter for a new line, so its steps advance on Mod+Enter
  const nextKey = (s) => (TEXT[s] || s === "outcome" ? MOD + "+Enter" : "Enter");

  // ---- hints ----

  function hints() {
    const s = cur();
    const back = W.at ? "back" : "cancel";
    let h;
    let mode = "STEP " + (W.at + 1) + " / " + steps.length;
    if (s !== "review") h = [[nextKey(s), "next"], ["Esc", back]];
    else {
      mode = "REVIEW";
      if (W.phase === "conflict") h = [["Enter", "reload the records"], ["Esc", "back"]];
      else if (W.phase === "review") h = [["Enter", "write"], ["Esc", back]];
      else h = [["Esc", back]];
    }
    deps.setDialogHints(mode, h);
  }

  // ---- bodies ----

  function textStep(s) {
    const t = TEXT[s];
    const id = "wz-" + s;
    return [
      el(
        "label",
        { class: "field", for: id },
        el("span", { class: "label", text: t.q }),
        el("span", { class: "help", text: t.help }),
        el("textarea", { id, rows: String(t.rows), "data-f": t.field, "aria-invalid": W.err ? "true" : null, "aria-describedby": W.err ? "wz-err" : null }, W.a[t.field]),
      ),
    ];
  }

  function outcomeStep() {
    const drop = to === "dropped";
    return [
      el(
        "label",
        { class: "field", for: "wz-outcome" },
        el("span", { class: "label", text: drop ? "Why was it dropped?" : "What was the outcome?" }),
        el("span", { class: "help", text: drop ? "The reason, in a line or two. It is written under ## Outcome." : "What changed and where, in one to three lines. Add the evidence for the done-when line. It is written under ## Outcome." }),
        el("textarea", { id: "wz-outcome", rows: "4", "data-f": "outcome", "aria-invalid": W.err ? "true" : null, "aria-describedby": W.err ? "wz-err" : null }, W.a.outcome),
      ),
    ];
  }

  function unblockStep() {
    return [
      el(
        "label",
        { class: "check" + (W.a.unblock ? " on" : "") },
        el("input", { type: "checkbox", id: "wz-unblock", checked: W.a.unblock ? true : null }),
        el("span", { class: "label", text: "Clear blocked_by?" }),
        el("span", { class: "help", text: "It no longer waits on what it listed. Untick to keep the list." }),
      ),
    ];
  }

  function identityStep() {
    return [
      el(
        "label",
        { class: "field", for: "wz-identity" },
        el("span", { class: "label", text: "Who is claiming it?" }),
        el("span", { class: "help", text: "This repository records claims. The name goes in claimed_by and the log, and is remembered in this browser." }),
        el("input", { id: "wz-identity", type: "text", autocomplete: "off", "data-f": "identity", value: W.a.identity }),
      ),
    ];
  }

  function decisionStep() {
    const d = W.a.decision;
    const radio = (mode, label, help) =>
      el(
        "label",
        { class: "check" + (d.mode === mode ? " on" : "") },
        el("input", { type: "radio", name: "wz-dec", id: "wz-dec-" + mode, value: mode, checked: d.mode === mode ? true : null }),
        el("span", { class: "box", "aria-hidden": "true", text: d.mode === mode ? "●" : "○" }),
        el("span", { class: "q" }, label, help ? el("span", { class: "help", text: " " + help }) : null),
      );
    const kids = [
      el(
        "fieldset",
        { class: "adr" },
        el("legend", { text: "Did this decide something?" }),
        radio("none", "No", "Nothing to record."),
        radio("link", "Yes, link an ADR that exists", "Adds its id to decisions."),
        radio("new", "Yes, create a new ADR", "Makes a proposed ADR and links it. Run akmctl index afterwards to refresh the index."),
      ),
    ];
    if (d.mode === "link") {
      kids.push(
        el(
          "label",
          { class: "field", for: "wz-adr" },
          el("span", { class: "label", text: "Which ADR?" }),
          el("select", { id: "wz-adr" }, el("option", { value: "", text: "(choose)" }), spec.adrs.map((a) => el("option", { value: a.id, selected: a.id === d.id ? true : null, text: a.id + " " + a.title + " (" + a.status + ")" }))),
        ),
      );
    }
    if (d.mode === "new") {
      kids.push(
        el(
          "label",
          { class: "field", for: "wz-adr-title" },
          el("span", { class: "label", text: "Title of the new ADR" }),
          el("span", { class: "help", text: d.slug ? "File name: " + d.slug : "The name is made from the title." }),
          el("input", { id: "wz-adr-title", type: "text", autocomplete: "off", value: d.title }),
        ),
      );
    }
    return kids;
  }

  function reviewBody() {
    switch (W.phase) {
      case "loading":
        return [el("p", { class: "muted", text: "Asking the server what would change…" })];
      case "writing":
        return [el("p", { class: "muted", text: "Writing…" })];
      case "conflict":
        return [
          el(
            "div",
            { class: "conflict", role: "alert" },
            el("p", { class: "eyebrow" }, el("span", { class: "g", "aria-hidden": "true", text: "⚠ " }), "Conflict"),
            el("h3", { text: "The record changed since you loaded it" }),
            el("p", { text: W.written ? W.written + " of " + W.ops.length + " writes were made before the change was found; the rest were not made." : "It was edited outside the workbench. Nothing was written." }),
            el("p", { text: "Reload the records to see its current text, then make the move again." }),
          ),
        ];
      case "error":
        return [
          el("div", { class: "err", role: "alert", text: W.fatal }),
          W.canTakeover ? el("p", { class: "help", text: "You can take the claim over: it is recorded in the log." }) : null,
          W.written ? el("p", { class: "help", text: W.written + " of " + W.ops.length + " writes were made before this." }) : null,
        ];
    }
    const parts = [el("p", { class: "muted", text: "This is exactly what each write will change. Nothing is written until you confirm." + (W.previews.length > 1 ? " Each diff is against the file as it is now; the writes are made in this order, so a later one builds on the earlier." : "") })];
    W.previews.forEach((p, n) => {
      parts.push(el("h3", { class: "sec-h" }, "Write " + (n + 1) + " of " + W.previews.length + ": " + opLabel(p.op, rec.id)));
      if (p.diff) parts.push(renderDiff(p.diff));
      else parts.push(el("p", { class: "note-card", text: p.note }));
    });
    return parts;
  }

  function body() {
    const s = cur();
    const err = W.err ? el("p", { id: "wz-err", class: "err", role: "alert" }, el("span", { class: "g", "aria-hidden": "true", text: "⚠ " }), W.err) : null;
    let kids;
    if (TEXT[s]) kids = textStep(s);
    else if (s === "outcome") kids = outcomeStep();
    else if (s === "identity") kids = identityStep();
    else if (s === "unblock") kids = unblockStep();
    else if (s === "decision") kids = decisionStep();
    else return reviewBody();
    return err ? kids.concat(err) : kids;
  }

  // ---- drawing ----

  function render(focusId) {
    const s = cur();
    const list = el(
      "ol",
      { class: "wsteps", "aria-label": "Progress" },
      steps.map((n, k) => el("li", { class: k < W.at ? "done" : k === W.at ? "now" : "", "aria-current": k === W.at ? "step" : null }, el("span", { class: "g", "aria-hidden": "true", text: k < W.at ? "✓" : k === W.at ? "●" : "○" }), " " + (k + 1) + ". " + STEP_NAMES[n])),
    );
    const heading = to === "done" ? "Close as done" : to === "dropped" ? "Drop" : "Move to " + to;
    const last = s === "review";
    const writeLabel = to === "done" ? "Write and close" : to === "dropped" ? "Write and drop" : to === "blocked" ? "Write and block" : to === "active" && move.claims ? "Write and claim" : "Write the move";
    const canWrite = W.phase === "review";
    const foot = [el("button", { type: "button", class: "btn quiet", "data-act": "back" }, (W.at === 0 ? "Cancel " : "Back "), el("kbd", { text: "Esc" }))];
    if (last && W.phase === "conflict") foot.push(el("button", { type: "button", class: "btn primary", "data-act": "reload" }, "Reload the records ", el("kbd", { text: "Enter" })));
    else if (last && W.canTakeover && W.phase === "error") foot.push(el("button", { type: "button", class: "btn primary", "data-act": "takeover" }, "Take over the claim"));
    else if (last) foot.push(el("button", { type: "button", class: "btn primary", "data-act": "write", disabled: canWrite ? null : true }, writeLabel + " ", el("kbd", { text: "Enter" })));
    else foot.push(el("button", { type: "button", class: "btn primary", "data-act": "next" }, "Next ", el("kbd", { text: nextKey(s) })));
    root.replaceChildren(
      el(
        "header",
        { class: "dlg-head" },
        el("div", {}, el("p", { class: "eyebrow", text: heading + " · step " + (W.at + 1) + " of " + steps.length }), el("h2", { id: "wz-t" }, el("span", { class: "mono", text: rec.id }), " " + spec.title), el("p", { class: "sub" }, move.from + " ", el("span", { class: "g", "aria-hidden": "true", text: "→" }), el("span", { class: "sr", text: " to" }), " " + to)),
        steps.length > 1 ? list : null,
      ),
      el("div", { class: "dlg-body" }, body()),
      el("footer", { class: "dlg-foot" }, foot),
    );
    hints();
    let f = focusId ? root.querySelector("#" + focusId) : null;
    if (!f) {
      if (last) f = root.querySelector(".diff") || root.querySelector(".conflict") || root.querySelector("[data-act=write]:not([disabled])") || root.querySelector("[data-act=back]");
      else f = root.querySelector("textarea, input[type=text], input:checked, select");
    }
    if (f) f.focus({ preventScroll: true });
  }

  // ---- steps ----

  function next() {
    if (W.phase !== "edit" && cur() !== "review") return;
    if (cur() === "review") return; // Enter writes; Ctrl+Enter never does
    W.err = validateStep(cur(), W.a, to);
    if (W.err) {
      const dm = W.a.decision.mode;
      render(cur() === "decision" ? (dm === "link" ? "wz-adr" : dm === "new" ? "wz-adr-title" : null) : "wz-" + cur());
      return;
    }
    W.at++;
    if (cur() === "review") loadReview();
    else render();
  }

  function back() {
    W.err = "";
    W.seq++;
    if (W.at === 0) {
      deps.closeDialog();
      return;
    }
    W.phase = "edit";
    W.fatal = null;
    W.canTakeover = false;
    W.at--;
    render();
  }

  // ---- review (dry runs) ----

  async function loadReview() {
    const mine = ++W.seq;
    W.phase = "loading";
    W.previews = [];
    W.fatal = null;
    W.canTakeover = false;
    W.ops = ops();
    render();
    const previews = [];
    try {
      for (const op of W.ops) {
        if (!previewable(op, move)) {
          previews.push({ op, note: op.kind === "new" ? "A new ADR file is created when you write. It has no diff to show: its status is proposed, and it is linked from this record's decisions." : op.kind === "claim" ? "The claim is written after the status change: it sets claimed_by and claimed_at and adds a log line." : "This write names the new ADR, which does not exist yet, so it is shown when written." });
          continue;
        }
        const q = requestFor(op, rec.id, { dry_run: true, takeover: W.takeover });
        const r = await api(q.method, q.path, { body: q.body, etag: rec.etag });
        if (mine !== W.seq) return;
        previews.push({ op, diff: r.diff });
      }
      if (mine !== W.seq) return;
      W.previews = previews;
      W.phase = "review";
    } catch (e) {
      if (mine !== W.seq) return;
      fail(e);
    }
    render();
  }

  function fail(e) {
    if (e instanceof ApiError && e.code === "stale") {
      W.phase = "conflict";
      return;
    }
    W.phase = "error";
    W.fatal = describeError(e);
    W.canTakeover = e instanceof ApiError && e.code === "refused" && e.details && e.details.reason === "claimed";
  }

  // ---- writes ----

  async function write() {
    if (W.phase !== "review") return;
    W.phase = "writing";
    render();
    let etag = rec.etag;
    let newId = null;
    W.written = 0;
    let last = null;
    try {
      for (const raw of W.ops) {
        const op = resolveOp(raw, newId);
        const q = requestFor(op, rec.id, { takeover: W.takeover });
        const r = await api(q.method, q.path, { body: q.body, etag: op.kind === "new" ? undefined : etag });
        if (op.kind === "new") newId = r.record.id;
        else etag = r.record.etag;
        last = r;
        W.written++;
      }
    } catch (e) {
      fail(e);
      render();
      return;
    }
    deps.closeDialog();
    spec.onDone({ ops: W.ops, record: last && last.record, newAdr: newId, identity: W.a.identity });
  }

  // ---- events ----

  root.addEventListener("click", (e) => {
    const b = e.target.closest("[data-act]");
    if (!b) return;
    const a = b.dataset.act;
    if (a === "back") back();
    else if (a === "next") next();
    else if (a === "write") write();
    else if (a === "reload") reload();
    else if (a === "takeover") {
      W.takeover = true;
      loadReview();
    }
  });
  root.addEventListener("input", (e) => {
    const t = e.target;
    if (t.dataset && t.dataset.f) W.a[t.dataset.f] = t.value;
    else if (t.id === "wz-adr") W.a.decision.id = t.value;
    else if (t.id === "wz-adr-title") {
      W.a.decision.title = t.value;
      W.a.decision.slug = slugify(t.value);
      const help = root.querySelector("#wz-adr-title").parentNode.querySelector(".help");
      if (help) help.textContent = W.a.decision.slug ? "File name: " + W.a.decision.slug : "The name is made from the title.";
    }
  });
  root.addEventListener("change", (e) => {
    if (e.target.id === "wz-unblock") W.a.unblock = e.target.checked;
    else if (e.target.name === "wz-dec") {
      W.a.decision.mode = e.target.value;
      W.err = "";
      render("wz-dec-" + e.target.value);
    }
  });

  function reload() {
    deps.closeDialog();
    spec.onReload();
  }

  function onKey(e) {
    const act = wizardKey({ key: e.key, mod: e.metaKey || e.ctrlKey, tag: e.target.tagName, type: e.target.type, phase: W.phase, step: cur() });
    if (!act) return;
    e.preventDefault();
    if (act === "back") back();
    else if (act === "next") next();
    else if (act === "reload") reload();
    else write();
  }

  deps.openDialog(root, { mode: "STEP", hints: [], onKey });
  if (cur() === "review") loadReview();
  else render();
}

