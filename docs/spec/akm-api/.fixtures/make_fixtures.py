#!/usr/bin/env python3
"""Regenerate the fixture repositories and the expected responses (API spec section 9).

A throwaway reference, not an implementation of the tool: it knows only the two sample
repositories below. Run it, review the diff, commit. Python 3 standard library only.
"""
import datetime as dt, hashlib, json, os, re, shutil

HERE = os.path.dirname(os.path.abspath(__file__))
AS_OF = "2026-10-28"

LAYOUT = "layout: 0.3.0\nextension.adr-index: 0.3.0\nextension.work-items: 0.3.0\n"
AGENTS = """# Sample

A small sample repository for the AKM API fixtures, used by the maintainers of the sample tool.

## Commands

- Check the records: `akmctl check`

## Where things live

- Decisions: `docs/adr/`. Roadmap: `docs/ROADMAP.md`. Open work: `akmctl ls --status ready,active,blocked`.

## Rules

- Work only under a work item. Update its `## Next` and `## Log` with the code.

## Do not, without asking

- Edit the body of an accepted ADR.
"""

def work(id_, slug, title, status, created, updated, extra=None, sections=None, summary=None, done=False, closed=None):
    fm = [("id", id_), ("title", title)]
    if summary: fm.append(("summary", summary))
    fm.append(("status", status))
    for k, v in (extra or []): fm.append((k, v))
    fm += [("created", created), ("updated", updated)]
    if closed: fm.append(("closed", closed))
    body = "\n# %s\n\n" % title + "".join("## %s\n%s\n\n" % (h, t) for h, t in (sections or [])).rstrip("\n") + "\n"
    return ("docs/work/%s%s-%s.md" % ("done/" if done else "", id_, slug), fm, body)

def adr(id_, slug, title, summary, status, date, scope=None):
    fm = [("id", id_), ("title", title), ("summary", summary), ("status", status), ("date", date)]
    if scope: fm.append(("scope", scope))
    body = "\n# %s\n\n## Context\nWhy this came up.\n\n## Considered Options\n1. The chosen way.\n2. Do nothing.\n\n## Decision\nThe chosen way.\n\n## Consequences\nIt works.\n" % title
    return ("docs/adr/%s-%s.md" % (id_, slug), fm, body)

def fmt(v):
    if isinstance(v, list): return "[" + ", ".join(fmt(x) if " " in x else x for x in v) + "]" if False else "[" + ", ".join(('"%s"' % x if " " in x else x) for x in v) + "]"
    if isinstance(v, int): return str(v)
    if re.search(r"[:#]|^\s|\s$", v): return '"%s"' % v
    return v

def render(rec):
    path, fm, body = rec
    return path, "---\n" + "".join("%s: %s\n" % (k, fmt(v)) for k, v in fm) + "---\n" + body

GOAL = ("Goal", "Do the thing.\n\n- [ ] It is done.")
RECS = [
    work("261010-j9k0", "choose-a-licence", "Choose a licence", "done", "2026-10-10", "2026-10-12", done=True, closed="2026-10-12",
         sections=[GOAL, ("Outcome", "MIT, in `LICENSE`.")]),
    work("261015-c4d5", "machine-readable-listings", "Machine-readable listings", "active", "2026-10-15", "2026-10-26",
         extra=[("kind", "feature"), ("priority", 3)], summary="Make listings available as JSON for tools.",
         sections=[GOAL, ("Next", "Finish the children."), ("Log", "- 2026-10-15: opened.")]),
    work("261020-a2b3", "add-json-output-to-ls", "Add JSON output to ls", "active", "2026-10-20", "2026-10-27",
         extra=[("parent", "261015-c4d5"), ("priority", 2)], summary="Print ls results as JSON.",
         sections=[GOAL, ("Next", "Write the encoder test.\n\nSaved: 2026-10-27 18:00"), ("Log", "- 2026-10-20: started.")]),
    work("261022-e5f6", "document-the-exit-codes", "Document the exit codes", "ready", "2026-10-22", "2026-10-22",
         extra=[("parent", "261015-c4d5")], sections=[GOAL]),
    work("261024-g7h8", "publish-the-first-release", "Publish the first release", "blocked", "2026-10-24", "2026-10-25",
         extra=[("blocked_by", ["261020-a2b3", "legal review"]), ("priority", 1)],
         sections=[GOAL, ("Next", "Wait for the JSON output.")]),
    work("261026-m1n2", "explore-a-plugin-system", "Explore a plugin system", "raw", "2026-10-26", "2026-10-26",
         sections=[("Notes", "Maybe later.")]),
    adr("261012-p3q4", "store-records-as-markdown-files", "Store records as Markdown files",
        "Records are Markdown files with frontmatter; read this before changing the format.", "accepted", "2026-10-12", ["docs/"]),
    adr("261018-r5s6", "use-a-local-server-for-the-ui", "Use a local server for the UI",
        "The UI talks to a loopback server; read this before adding a network feature.", "proposed", "2026-10-18"),
]
ROADMAP = "# Roadmap\n\n## Now\n- [261015-c4d5] Machine-readable listings\n## Next\n- [261024-g7h8] First release\n## Later\n## Dropped\n"

def parse(text):
    m = re.match(r"---\n(.*?)\n---\n(.*)", text, re.S)
    fm, body = {}, m.group(2)
    for i, line in enumerate(m.group(1).split("\n")):
        k, _, v = line.partition(": ")
        if v.startswith("["):
            v = [x.strip().strip('"') for x in v[1:-1].split(",") if x.strip()]
        elif re.fullmatch(r"-?\d+", v): v = int(v)
        elif v in ("true", "false"): v = v == "true"
        else: v = v.strip('"')
        fm[k] = v
    return fm, body

def pad(rows):
    w = [max(len(r[i]) for r in rows) for i in range(len(rows[0]))]
    out = []
    for n, r in enumerate(rows):
        out.append("| " + " | ".join(c.ljust(w[i]) for i, c in enumerate(r)) + " |")
        if n == 0: out.append("| " + " | ".join("-" * w[i] for i in range(len(w))) + " |")
    return "\n".join(out)

def index_text(adrs):
    groups = [("In force", ["accepted"]), ("Proposed", ["proposed"]), ("Not in force", ["rejected", "deprecated", "superseded"])]
    out = "<!-- generated by akmctl index; do not edit -->\n"
    for name, sts in groups:
        rows = [["id", "title", "date", "scope", "summary", "supersedes / superseded_by"]]
        for path, fm in sorted(adrs, key=lambda a: a[1]["id"]):
            if fm["status"] in sts:
                rows.append(["[%s](%s)" % (fm["id"], os.path.basename(path)), fm["title"], fm["date"], ", ".join(fm.get("scope", [])),
                             fm["summary"], " / ".join(x for x in (fm.get("supersedes"), fm.get("superseded_by")) if x)])
        if len(rows) > 1: out += "\n## %s\n\n%s\n" % (name, pad(rows))
    return out

def build(name, mutate=None):
    root = os.path.join(HERE, "repos", name)
    shutil.rmtree(root, ignore_errors=True)
    files = {"AGENTS.md": AGENTS, "docs/.layout": LAYOUT, "docs/ROADMAP.md": ROADMAP}
    for r in RECS:
        p, t = render(r); files[p] = t
    if mutate: mutate(files)
    adrs = [(p, parse(t)[0]) for p, t in files.items() if p.startswith("docs/adr/") and p.endswith(".md") and not p.endswith("README.md")]
    files["docs/adr/README.md"] = index_text(adrs)
    for p, t in files.items():
        fp = os.path.join(root, p); os.makedirs(os.path.dirname(fp), exist_ok=True)
        open(fp, "w", newline="\n").write(t)
    return files

def broken(files):
    k = "docs/work/261022-e5f6-document-the-exit-codes.md"
    files[k] = files[k].replace("parent: 261015-c4d5", "parent: 261099-aaaa")
    k = "docs/work/261020-a2b3-add-json-output-to-ls.md"
    files[k] = re.sub(r"## Next\n.*?\n\n(?=## Log)", "", files[k], flags=re.S)
    files["docs/ROADMAP.md"] = ROADMAP.replace("## Later\n", "## Later\n- [261099-bbbb] Ghost\n")

WORK_F = ["id", "title", "summary", "status", "kind", "parent", "blocked_by", "priority", "cycle", "keep", "created", "updated", "closed", "authored_by", "decisions"]
ADR_F = ["id", "title", "summary", "status", "date", "scope", "supersedes", "superseded_by", "authored_by"]

def load(files):
    recs = []
    claims = "extension.claims" in files.get("docs/.layout", "")
    for p, t in sorted(files.items()):
        if re.match(r"docs/(work(/done)?|adr)/\d{6}-", p):
            fm, body = parse(t)
            typ = "adr" if p.startswith("docs/adr/") else "work"
            fields = {k: fm.get(k) for k in (ADR_F if typ == "adr" else WORK_F + (["claimed_by", "claimed_at"] if claims else []))}
            for k in fm:
                if k not in fields: fields[k] = fm[k]
            d = fm["id"][:6]
            recs.append(dict(id=fm["id"], type=typ, path=p, id_date="20%s-%s-%s" % (d[:2], d[2:4], d[4:]),
                             etag=hashlib.sha256(t.encode()).hexdigest()[:16], fields=fields, _body=body, _text=t))
    return sorted(recs, key=lambda r: r["id"])

SHORT = {"work": ["id", "title", "summary", "status", "kind", "parent", "priority", "updated"],
         "adr": ["id", "title", "summary", "status", "date", "scope"]}

def dflt(r, extra=()):
    return SHORT[r["type"]] + (list(extra) if r["type"] == "work" else [])

def pub(r, fields=None, next_=False):
    o = {k: r[k] for k in ("id", "type", "path", "id_date", "etag")}
    o["fields"] = dict(r["fields"]) if fields is None else {k: r["fields"].get(k) for k in fields}
    if next_ and r["type"] == "work":
        m = re.search(r"## Next\n(.*?)(\n\n|\Z)", r["_body"], re.S)
        o["next_line"] = m.group(1).split("\n")[0] if m else None
    return o

def roadmap_rank(files):
    ids = re.findall(r"^- \[(\d{6}-\w{4})\]", files["docs/ROADMAP.md"].split("## Dropped")[0], re.M)
    return {i: n for n, i in enumerate(ids)}

def sort_recs(recs, keys, rank):
    for key in reversed(keys):
        desc = key.startswith("-"); k = key.lstrip("-")
        def val(r, k=k):
            if k == "roadmap": return rank.get(r["id"])
            return r["fields"].get(k) if k != "id" else r["id"]
        present = [r for r in recs if val(r) is not None]; absent = [r for r in recs if val(r) is None]
        present.sort(key=val, reverse=desc)
        recs = present + absent
    return recs

def commit(step, item, detail, *paths):
    """The commit object of ADR 261030-tbkh."""
    return dict(step=step, item=item, detail=detail, paths=list(paths),
                message="%s(%s): %s\n\nAKM-Step: %s\nAKM-Item: %s" % (item, step, detail, step, item))

def ok(op, req, **kw): return dict(api_version="1.0", operation=op, request=req, ok=True, **kw)
def err(op, req, code, msg, **d):
    e = dict(code=code, message=msg)
    if d: e["details"] = d
    return dict(api_version="1.0", operation=op, request=req, ok=False, error=e)

E5F6 = "docs/work/261022-e5f6-document-the-exit-codes.md"
A2B3 = "docs/work/261020-a2b3-add-json-output-to-ls.md"
A2B3_DONE = "docs/work/done/261020-a2b3-add-json-output-to-ls.md"
# the open item that names a2b3 in blocked_by: close reports it as a dependent
A2B3_DEPENDENTS = [dict(id="261024-g7h8", status="blocked", title="Publish the first release", path="docs/work/261024-g7h8-publish-the-first-release.md")]

def claims_layout(files):
    files["docs/.layout"] += "extension.claims: 0.3.0\n"
    files[A2B3] = files[A2B3].replace("updated: 2026-10-27\n---", "updated: 2026-10-27\nclaimed_by: agent-9\nclaimed_at: 2026-10-28T08:00\n---")

def closable(files):
    files[A2B3] += "\n## Outcome\nShipped.\n"


def mkdiff(path, new_path, old, new):
    """The diff member of a dry-run response: unified hunks, 3 lines of context."""
    def lines(t):
        ls = t.split("\n")
        if ls[-1] == "": ls.pop()
        else: ls[-1] += "\0"  # a missing final newline is a change
        return ls
    a, b = lines(old), lines(new)
    # longest common subsequence; a line that matches is taken at once, so a run of
    # added or removed lines sits at the latest position it can
    n, m = len(a), len(b)
    L = [[0] * (m + 1) for _ in range(n + 1)]
    for i in range(n - 1, -1, -1):
        for j in range(m - 1, -1, -1):
            L[i][j] = L[i + 1][j + 1] + 1 if a[i] == b[j] else max(L[i + 1][j], L[i][j + 1])
    ops, i, j = [], 0, 0
    while i < n or j < m:
        if i < n and j < m and a[i] == b[j]: ops.append(("ctx", i + 1, j + 1)); i += 1; j += 1
        elif i < n and (j == m or L[i + 1][j] >= L[i][j + 1]): ops.append(("rem", i + 1, None)); i += 1
        else: ops.append(("add", None, j + 1)); j += 1
    def line(o):
        op, x, y = o
        return dict(op=op, text=(a[x - 1] if x else b[y - 1]).rstrip("\0"), old=x, new=y)
    hunks, k = [], 0
    while k < len(ops):
        while k < len(ops) and ops[k][0] == "ctx": k += 1
        if k == len(ops): break
        start, last, e = max(0, k - 3), k, k
        while e < len(ops):
            if ops[e][0] != "ctx": last = e
            elif e - last > 6: break
            e += 1
        stop = min(len(ops), last + 4)
        out = [line(o) for o in ops[start:stop]]
        olds = [l["old"] for l in out if l["old"]]; news = [l["new"] for l in out if l["new"]]
        sp = lambda v: ("%d" % (v[0] if v else 0) if len(v) == 1 else "%d,%d" % (v[0] if v else 0, len(v)))
        hunks.append(dict(header="@@ -%s +%s @@" % (sp(olds), sp(news)), lines=out))
        k = stop
    return dict(path=path, new_path=new_path, hunks=hunks)

def set_section(text, name, body, append=False):
    """Reference for the section operation, for the sample files only: no code fences, LF."""
    lines = text.rstrip("\n").split("\n")
    heads = [i for i, l in enumerate(lines) if l.startswith("## ")]
    at = next((i for i in heads if lines[i][3:].strip().lower() == name.lower()), None)
    new = body.strip("\n").split("\n") if body.strip("\n") else []
    if at is None:
        log = next((i for i in heads if lines[i][3:].strip().lower() == "log"), None)
        block = ["## " + name[:1].upper() + name[1:]] + ([""] if new else []) + new
        if log is None: return "\n".join(lines + [""] + block) + "\n"
        return "\n".join(lines[:log] + block + [""] + lines[log:]) + "\n"
    end = next((i for i in heads if i > at), len(lines))
    cont = lines[at + 1:end]
    idx = [i for i, l in enumerate(cont) if l.strip()]
    lead, core, trail = (cont[:idx[0]], cont[idx[0]:idx[-1] + 1], cont[idx[-1] + 1:]) if idx else (cont, [], [])
    core = core + new if append else new
    if not trail and end < len(lines): trail = [""]
    return "\n".join(lines[:at + 1] + lead + core + trail + lines[end:]) + "\n"

def write_cases(cases, basic):
    clm = build("claims", claims_layout)
    clo = build("closable", closable)
    def after_rec(files, path, removed=()):
        fs = dict(files)
        for r in removed: fs.pop(r, None)
        return next(r for r in load(fs) if r["path"] == path)
    def case(op, n, name, repo, request, expected, after=None, removed=None, would_be=None):
        c = dict(name=name, repo=repo, request=request, expected=expected)
        if after: c["after"] = after
        if removed: c["removed"] = removed
        if would_be: c["would_be"] = would_be
        cases.append((op, n, c))
    # set
    req = dict(id="261022-e5f6", changes=dict(priority=4, summary="Document them."), as_of="2026-10-28")
    text = basic[E5F6].replace("updated: 2026-10-22\n---\n", "updated: 2026-10-28\npriority: 4\nsummary: Document them.\n---\n")
    files = dict(basic); files[E5F6] = text
    case("set", 1, "change two fields, refresh updated", "basic", req,
         ok("set", req, record=pub(after_rec(files, E5F6)), changed=["priority", "summary"]), after={E5F6: text})
    req = dict(id="261022-e5f6", changes=dict(status="blocked"))
    case("set", 2, "blocked needs blocked_by", "basic", req, err("set", req, "refused", "status blocked needs a non-empty blocked_by", reason="rule", rule="WK5"))
    cur = next(r for r in load(basic) if r["id"] == "261022-e5f6")
    req = dict(id="261022-e5f6", changes=dict(priority=1), if_match="0000000000000000")
    case("set", 3, "a stale etag writes nothing", "basic", req, err("set", req, "stale", "the record changed since it was read", etag=cur["etag"]))
    # close
    req = dict(id="261020-a2b3", status="done")
    case("close", 1, "an empty Outcome refuses", "basic", req, err("close", req, "refused", "261020-a2b3 has no ## Outcome; write it before closing", reason="empty_outcome"))
    req = dict(id="261020-a2b3", status="done", as_of="2026-10-28")
    text = clo[A2B3].replace("status: active\n", "status: done\n").replace("updated: 2026-10-27\n", "updated: 2026-10-28\nclosed: 2026-10-28\n")
    files = dict(clo); del files[A2B3]; files[A2B3_DONE] = text
    case("close", 2, "close moves the file to done/", "closable", req,
         ok("close", req, record=pub(after_rec(files, A2B3_DONE)), moved_from=A2B3, moved_to=A2B3_DONE, dependents=A2B3_DEPENDENTS, commit=commit("close", "261020-a2b3", "done", A2B3, A2B3_DONE)), after={A2B3_DONE: text}, removed=[A2B3])
    req = dict(id="261020-a2b3", status="dropped", as_of="2026-10-28")
    text = clo[A2B3].replace("status: active\n", "status: dropped\n").replace("updated: 2026-10-27\n", "updated: 2026-10-28\nclosed: 2026-10-28\n")
    files = dict(clo); del files[A2B3]; files[A2B3_DONE] = text
    case("close", 6, "dropping an item also writes closed", "closable", req,
         ok("close", req, record=pub(after_rec(files, A2B3_DONE)), moved_from=A2B3, moved_to=A2B3_DONE, dependents=A2B3_DEPENDENTS, commit=commit("drop", "261020-a2b3", "dropped", A2B3, A2B3_DONE)), after={A2B3_DONE: text}, removed=[A2B3])
    # cycle, keep and the managed closed field
    req = dict(id="261022-e5f6", changes=dict(cycle="2026-PI3-S4"), as_of="2026-10-28")
    text = basic[E5F6].replace("updated: 2026-10-22\n---\n", "updated: 2026-10-28\ncycle: 2026-PI3-S4\n---\n")
    files = dict(basic); files[E5F6] = text
    case("set", 5, "set the cycle label", "basic", req, ok("set", req, record=pub(after_rec(files, E5F6)), changed=["cycle"]), after={E5F6: text})
    req = dict(id="261020-a2b3", changes=dict(keep=True), as_of="2026-10-28")
    text = basic[A2B3].replace("updated: 2026-10-27\n---\n", "updated: 2026-10-28\nkeep: true\n---\n")
    files = dict(basic); files[A2B3] = text
    case("set", 6, "pin an item with keep", "basic", req, ok("set", req, record=pub(after_rec(files, A2B3)), changed=["keep"]), after={A2B3: text})
    req = dict(id="261022-e5f6", changes=dict(closed="2026-10-28"))
    case("set", 7, "closed is set only by close", "basic", req, err("set", req, "invalid_request", 'changes: "closed" cannot be set on a work (id, created, updated, closed, date and claim fields are managed by the tool)', parameter="changes"))
    req = dict(id="261026-m1n2", changes=dict(status="ready"), as_of="2026-10-28")
    M1N2 = "docs/work/261026-m1n2-explore-a-plugin-system.md"
    text = basic[M1N2].replace("status: raw\n", "status: ready\n").replace("updated: 2026-10-26\n", "updated: 2026-10-28\n")
    files = dict(basic); files[M1N2] = text
    case("set", 8, "a status change returns the commit object", "basic", req, ok("set", req, record=pub(after_rec(files, M1N2)), changed=["status"], commit=commit("ready", "261026-m1n2", "ready", M1N2)), after={M1N2: text})
    # claim, release (repo with the claims extension)
    req = dict(id="261022-e5f6", **{"as": "agent-1"}, as_of="2026-10-28T09:30")
    text = clm[E5F6].replace("status: ready\n", "status: active\n").replace("updated: 2026-10-22\n---\n", "updated: 2026-10-28\nclaimed_by: agent-1\nclaimed_at: 2026-10-28T09:30\n---\n")
    text = text.rstrip("\n") + "\n\n## Log\n- 2026-10-28: claimed (agent-1).\n"
    files = dict(clm); files[E5F6] = text
    case("claim", 1, "claim a ready item", "claims", req, ok("claim", req, record=pub(after_rec(files, E5F6)), previous_holder=None, commit=commit("claim", "261022-e5f6", "claimed by agent-1", E5F6)), after={E5F6: text})
    req = dict(id="261020-a2b3", **{"as": "agent-1"}, as_of="2026-10-28T09:30")
    case("claim", 2, "another holder's fresh claim refuses", "claims", req, err("claim", req, "refused", "261020-a2b3 is claimed by agent-9", reason="claimed"))
    req = dict(id="261020-a2b3", as_of="2026-10-28")
    text = clm[A2B3].replace("claimed_by: agent-9\nclaimed_at: 2026-10-28T08:00\n", "").replace("updated: 2026-10-27\n", "updated: 2026-10-28\n")
    text = text.rstrip("\n") + "\n- 2026-10-28: released (agent-9).\n"
    files = dict(clm); files[A2B3] = text
    case("release", 1, "release a claim", "claims", req, ok("release", req, record=pub(after_rec(files, A2B3)), commit=commit("release", "261020-a2b3", "released", A2B3)), after={A2B3: text})

    # dry_run: the would-be record carries the current etag, nothing is written
    req = dict(id="261022-e5f6", changes=dict(priority=4, summary="Document them."), dry_run=True, as_of="2026-10-28")
    text = basic[E5F6].replace("updated: 2026-10-22\n---\n", "updated: 2026-10-28\npriority: 4\nsummary: Document them.\n---\n")
    files = dict(basic); files[E5F6] = text
    rec = after_rec(files, E5F6); rec["etag"] = cur["etag"]
    case("set", 4, "dry run: the change as a diff, nothing written", "basic", req,
         ok("set", req, record=pub(rec), changed=["priority", "summary"], diff=mkdiff(E5F6, E5F6, basic[E5F6], text)), would_be={E5F6: text})
    req = dict(id="261020-a2b3", status="done", outcome="Shipped in 0.2.", dry_run=True, as_of="2026-10-28")
    text = set_section(basic[A2B3], "Outcome", "Shipped in 0.2.").replace("status: active\n", "status: done\n").replace("updated: 2026-10-27\n", "updated: 2026-10-28\nclosed: 2026-10-28\n")
    files = dict(basic); del files[A2B3]; files[A2B3_DONE] = text
    rec = after_rec(files, A2B3_DONE); rec["etag"] = next(r for r in load(basic) if r["path"] == A2B3)["etag"]
    case("close", 4, "dry run of a close with an outcome: the diff names the new path", "basic", req,
         ok("close", req, record=pub(rec), moved_from=A2B3, moved_to=A2B3_DONE, dependents=A2B3_DEPENDENTS, commit=commit("close", "261020-a2b3", "done", A2B3, A2B3_DONE), diff=mkdiff(A2B3, A2B3_DONE, basic[A2B3], text)), would_be={A2B3_DONE: text})
    # close with outcome
    req = dict(id="261020-a2b3", status="done", outcome="Shipped in 0.2.", as_of="2026-10-28")
    files = dict(basic); del files[A2B3]; files[A2B3_DONE] = text
    case("close", 3, "close writes the outcome in the same write", "basic", req,
         ok("close", req, record=pub(after_rec(files, A2B3_DONE)), moved_from=A2B3, moved_to=A2B3_DONE, dependents=A2B3_DEPENDENTS, commit=commit("close", "261020-a2b3", "done", A2B3, A2B3_DONE)), after={A2B3_DONE: text}, removed=[A2B3])
    vault_write_cases(case, basic)
    # section
    def sec(n, name, repo, req, newtext, path, secname, sectext, changed=True, dryrun=False):
        src = {"basic": basic}[repo]
        files = dict(src); files[path] = newtext
        rec = after_rec(files, path)
        if dryrun: rec["etag"] = next(r for r in load(src) if r["path"] == path)["etag"]
        extra = dict(diff=mkdiff(path, path, src[path], newtext)) if dryrun else {}
        case("section", n, name, repo, req,
             ok("section", req, record=pub(rec), section=dict(name=secname, text=sectext), changed=changed, **extra),
             after=None if dryrun else {path: newtext}, would_be={path: newtext} if dryrun else None)
    req = dict(id="261022-e5f6", name="Goal", text="List every exit code.\n\n- [ ] The table is written.", as_of="2026-10-28")
    t1 = set_section(basic[E5F6], "Goal", req["text"]).replace("updated: 2026-10-22\n", "updated: 2026-10-28\n")
    sec(1, "replace a section", "basic", req, t1, E5F6, "goal", req["text"])
    req = dict(id="261020-a2b3", name="log", text="- 2026-10-28: tests written.", mode="append", as_of="2026-10-28")
    t2 = set_section(basic[A2B3], "Log", req["text"], append=True).replace("updated: 2026-10-27\n", "updated: 2026-10-28\n")
    sec(2, "append to the Log section", "basic", req, t2, A2B3, "log", "- 2026-10-20: started.\n- 2026-10-28: tests written.")
    req = dict(id="261020-a2b3", name="Notes", text="Check the encoder against a large listing.", as_of="2026-10-28")
    t3 = set_section(basic[A2B3], "Notes", req["text"]).replace("updated: 2026-10-27\n", "updated: 2026-10-28\n")
    sec(3, "a missing section is created before the Log section", "basic", req, t3, A2B3, "notes", req["text"])
    req = dict(id="261020-a2b3", name="Next", text="Write the decoder test.", dry_run=True, as_of="2026-10-28")
    t4 = set_section(basic[A2B3], "Next", req["text"]).replace("updated: 2026-10-27\n", "updated: 2026-10-28\n")
    sec(4, "dry run: the replacement as a diff, nothing written", "basic", req, t4, A2B3, "next", req["text"], dryrun=True)
    req = dict(id="261020-a2b3", name="Log", text="- rewritten")
    case("section", 5, "the Log section cannot be replaced", "basic", req, err("section", req, "invalid_request", "the Log section is append-only: use mode append", parameter="mode"))
    req = dict(id="261010-j9k0", name="Notes", text="Late note.")
    case("section", 6, "a closed record is refused", "basic", req, err("section", req, "refused", "261010-j9k0 is closed; its sections are not edited", reason="closed"))

VW, VA, VR = "Deep Thought/9 pkm/tasks", "Deep Thought/9 pkm/adr", "Deep Thought/9 pkm/ROADMAP.md"

def relocate(p):
    for a, b in (("docs/work/", VW + "/"), ("docs/adr/", VA + "/")):
        if p.startswith(a): return b + p[len(a):]
    return VR if p == "docs/ROADMAP.md" else p

def build_vault(name):
    """The basic repository with its folders moved by the manifest, under a folder with spaces in its names."""
    root = os.path.join(HERE, "repos", name)
    shutil.rmtree(root, ignore_errors=True)
    basic = build("basic")
    files = {relocate(p): t for p, t in basic.items()}
    files["docs/.layout"] = LAYOUT + "work_dir: %s\nadr_dir: %s\nroadmap: %s\n" % (VW, VA, VR)
    for p, t in files.items():
        fp = os.path.join(root, p); os.makedirs(os.path.dirname(fp), exist_ok=True)
        open(fp, "w", newline="\n").write(t)
    return files

def vault_cases(case, basic, recs, rank):
    """Cases on the repository whose manifest sets work_dir, adr_dir and roadmap."""
    vault = build_vault("vault")
    vrecs = [dict(r, path=relocate(r["path"])) for r in recs]
    req = dict(type="work"); s = sort_recs([r for r in vrecs if r["type"] == "work"], ["id"], rank)
    case("ls", 6, "work items are read from the folder the manifest names", "vault", req, ok("ls", req, total=6, count=6, records=[pub(r, dflt(r)) for r in s]))
    case("check", 5, "a clean repository with moved folders", "vault", dict(as_of=AS_OF), ok("check", dict(as_of=AS_OF), clean=True, errors=0, warnings=0, count=0, findings=[], reconcile=dict(state="current")))
    req = dict(dry_run=True, include_content=True)
    case("index", 2, "the ADR index is written in the ADR folder", "vault", req, ok("index", req, path=VA + "/README.md", changed=False, written=False, content=vault[VA + "/README.md"]))
    repository = dict(name="vault", layout="0.3.0", extensions=["adr-index", "work-items"], reference_style="id", claim_stale_hours=24, work_dir=VW, adr_dir=VA, git_commits=False)
    case("version", 3, "the folders the manifest names are reported", "vault", {},
         ok("version", {}, tool=dict(name="akmctl", version="0.1.0"), api_versions=["1.0"], layout_versions=["0.3.0"], extensions=["claims"], repository=repository))


def vault_write_cases(case, basic):
    # close moves the file to done/ inside the configured work folder
    req = dict(id="261020-a2b3", status="done", outcome="Shipped in 0.2.", as_of="2026-10-28")
    text = set_section(basic[A2B3], "Outcome", "Shipped in 0.2.").replace("status: active\n", "status: done\n").replace("updated: 2026-10-27\n", "updated: 2026-10-28\nclosed: 2026-10-28\n")
    files = dict(basic); del files[A2B3]; files[A2B3_DONE] = text
    rec = dict(next(r for r in load(files) if r["path"] == A2B3_DONE), path=relocate(A2B3_DONE))
    case("close", 5, "close moves the file to done/ under the configured work folder", "vault", req,
         ok("close", req, record=pub(rec), moved_from=relocate(A2B3), moved_to=relocate(A2B3_DONE), dependents=[dict(d, path=relocate(d["path"])) for d in A2B3_DEPENDENTS], commit=commit("close", "261020-a2b3", "done", relocate(A2B3), relocate(A2B3_DONE))), after={relocate(A2B3_DONE): text}, removed=[relocate(A2B3)])

ALPHABET = "0123456789abcdefghjkmnpqrstvwxyz"

def seeded_id(date, seed, name):
    """The id mint_ids draws for a seed: the first four bytes of sha256(seed NUL name NUL 0), each reduced to 5 bits."""
    h = hashlib.sha256((seed + "\0" + name + "\0" + "0").encode()).digest()
    return date + "-" + "".join(ALPHABET[b & 31] for b in h[:4])

LEGACY_TASK = "Deep Thought/9 pkm/tasks/2026-10-20 Write the docs.md"
LEGACY_ADR = "Deep Thought/9 pkm/adr/ADR-001 Use SQLite.md"

def build_legacy(name):
    """A notes repository with no layout history: a manifest that names the folders, notes without ids."""
    root = os.path.join(HERE, "repos", name)
    shutil.rmtree(root, ignore_errors=True)
    files = {
        "AGENTS.md": AGENTS,
        "docs/.layout": "work_dir: Deep Thought/9 pkm/tasks\nadr_dir: Deep Thought/9 pkm/adr\nroadmap: Deep Thought/9 pkm/ROADMAP.md\n",
        LEGACY_TASK: "---\ntype: task\nstatus: ready\ncreated_at: 2026-10-20\nupdated_at: 2026-10-21\n---\n# Write the docs\n\nWhy it matters.\n\n## Done when\n\n- [ ] It is written.\n",
        LEGACY_ADR: "---\nid: 1\ntype: adr\nstatus: accepted\ndate: 2026-10-01\n---\n# Use SQLite\n",
    }
    for p, t in files.items():
        fp = os.path.join(root, p); os.makedirs(os.path.dirname(fp), exist_ok=True)
        open(fp, "w", newline="\n").write(t)
    return files

def legacy_migration():
    T, A = "Deep Thought/9 pkm/tasks", "Deep Thought/9 pkm/adr"
    return {"from": "0.0.0", "to": "0.3.0", "steps": [
        {"id": "created", "kind": "rename_field", "class": "auto", "description": "created_at becomes created", "dir": T, "from": "created_at", "to": "created"},
        {"id": "updated", "kind": "rename_field", "class": "auto", "description": "updated_at becomes updated", "dir": T, "from": "updated_at", "to": "updated"},
        {"id": "title", "kind": "title_from_h1", "class": "auto", "description": "title from the first heading", "dir": T},
        {"id": "goal", "kind": "goal_scaffold", "class": "auto", "description": "add a Goal section", "dir": T, "statuses": ["ready", "active"]},
        {"id": "adr-number", "kind": "rename_field", "class": "auto", "description": "the legacy id becomes adr_number", "dir": A, "from": "id", "to": "adr_number"},
        {"id": "ids-tasks", "kind": "mint_ids", "class": "auto", "description": "ids for tasks", "dir": T, "rename": False, "seed": "vault-tasks"},
        {"id": "ids-adrs", "kind": "mint_ids", "class": "auto", "description": "ids for ADRs", "dir": A, "date_field": "date", "rename": False, "seed": "vault-adrs"},
    ]}

def legacy_cases(case):
    build_legacy("legacy")
    mig = legacy_migration()
    req = dict(mode="plan", migrations=[mig])
    t_id = seeded_id("261020", "vault-tasks", LEGACY_TASK)
    a_id = seeded_id("261001", "vault-adrs", LEGACY_ADR)
    t_to = "Deep Thought/9 pkm/tasks/%s-write-the-docs.md" % t_id
    a_to = "Deep Thought/9 pkm/adr/%s-adr-001-use-sqlite.md" % a_id
    def st(id_, kind, files, **kw): return dict(id=id_, kind=kind, **{"class": "auto"}, files=files, status="pending", **kw)
    steps = [
        st("created", "rename_field", [LEGACY_TASK]), st("updated", "rename_field", [LEGACY_TASK]), st("title", "title_from_h1", [LEGACY_TASK]),
        st("goal", "goal_scaffold", [LEGACY_TASK], notes=[LEGACY_TASK + ": added ## Goal (introduction moved under it, Done when kept as a list); review the text"]),
        st("adr-number", "rename_field", [LEGACY_ADR]),
        st("ids-tasks", "mint_ids", [LEGACY_TASK], notes=["%s -> %s (id %s, date from created); not renamed: rename the file as listed in renames" % (LEGACY_TASK, t_to, t_id)], renames=[dict(**{"from": LEGACY_TASK}, to=t_to)]),
        st("ids-adrs", "mint_ids", [LEGACY_ADR], notes=["%s -> %s (id %s, date from date); not renamed: rename the file as listed in renames" % (LEGACY_ADR, a_to, a_id)], renames=[dict(**{"from": LEGACY_ADR}, to=a_to)]),
    ]
    finding = dict(rule="S1", severity="W", file="docs/.layout", line=None, id=None, message="the manifest has no layout key", **{"class": "info"})
    case("migrate", 2, "a one-off migration file plans steps for a repository with no layout history", "legacy", req,
         ok("migrate", req, mode="plan", **{"from": "0.0.0"}, to="0.3.0", state="behind", steps=steps, findings=[finding], clean=False))
    req = dict(mode="plan")
    case("migrate", 3, "without the migrations parameter there is no path from 0.0.0", "legacy", req,
         ok("migrate", req, mode="plan", **{"from": "0.0.0"}, to="0.3.0", state="unsupported", steps=[], findings=[finding], clean=False))
    bad = dict(mig, steps=[dict(mig["steps"][0], kind="run_script")])
    req = dict(mode="plan", migrations=[bad])
    case("migrate", 4, "an unknown step kind in a migration file is rejected", "legacy", req,
         err("migrate", req, "invalid_request", "migrations[0]: step created: unknown kind \"run_script\"", parameter="migrations"))

# init: the files a run writes (paths relative to the target directory). The
# small files are compared exactly; the others (AGENTS.md, the skills, ...) are only
# required to exist, so the wording of a template is not part of the contract.
AKM_REFS = ["board", "check", "claim", "close", "commits", "config", "find", "index", "init", "ls", "migrate", "new",
    "release", "search", "section", "set", "show", "version", "when"]
INIT_COMMON = [".agents/skills/akm/SKILL.md"] + [".agents/skills/akm/references/%s.md" % r for r in AKM_REFS] + [
    ".agents/skills/akm-record-decision/SKILL.md", ".agents/skills/akm-record-decision/assets/adr.md",
    ".agents/skills/akm-track-work/SKILL.md", ".agents/skills/akm-track-work/assets/work-item.md",
    ".prettierignore", ".prettierrc.json", "AGENTS.md", "CLAUDE.md", "README.md"]
INIT_LAYOUT = "layout: 0.3.0\nextension.adr-index: 0.3.0\nextension.work-items: 0.3.0\nextension.claims: 0.3.0\nextension.safety: 0.3.0\n"

COMMIT_SKILL = ".agents/skills/akm-commit/SKILL.md"

def init_files(docs="docs", wikilinks=False, git_commits=False, tree=False):
    work, adr, road = docs + "/work", docs + "/adr", docs + "/ROADMAP.md"
    paths = INIT_COMMON + [docs + "/.akm-deny", docs + "/.layout", road, adr + "/README.md", work + "/.gitkeep", work + "/done/.gitkeep"]
    paths += ([COMMIT_SKILL] if git_commits else []) + ([".gitignore"] if tree else [])
    layout = INIT_LAYOUT
    if docs != "docs": layout += "work_dir: %s\nadr_dir: %s\nroadmap: %s\n" % (work, adr, road)
    if wikilinks: layout += "reference_style: wikilink\n"
    if git_commits: layout += "git_commits: true\n"
    exact = {"CLAUDE.md": "@AGENTS.md\n", docs + "/.layout": layout,
             ".prettierrc.json": '{ "embeddedLanguageFormatting": "off" }\n', ".prettierignore": adr + "/README.md\n",
             adr + "/README.md": "<!-- generated by akmctl index; do not edit -->\n", work + "/.gitkeep": "", work + "/done/.gitkeep": ""}
    if tree: exact[".gitignore"] = "# the advisory lock akmctl takes while it writes\n" + docs + "/.akm.lock\n"
    return sorted(paths), exact

def init_cases(cases, basic):
    def case(n, name, repo, req, expected, after=None, after_paths=None, tree=False):
        c = dict(name=name, repo=repo, request=req, expected=expected)
        if after: c["after"] = after
        if after_paths: c["after_paths"] = after_paths
        if tree: c["git_work_tree"] = True
        cases.append(("init", n, c))
    def run(n, name, dirname, docs="docs", wikilinks=False, existing=(), git_commits=False, forced_off=False, tree=False):
        paths, exact = init_files(docs, wikilinks, git_commits and not forced_off, tree)
        req = dict(dir=dirname) if dirname else dict()
        if docs != "docs": req["docs_dir"] = docs
        if wikilinks: req["wikilinks"] = True
        if git_commits: req["git_commits"] = True
        written = [p for p in paths if p not in existing]
        skipped = [p for p in paths if p in existing]
        pre = "" if dirname in (None, ".") else dirname + "/"
        case(n, name, "basic", req, ok("init", req, docs_dir=docs, wikilinks=wikilinks, dir=dirname or ".", layout="0.3.0", written=written, skipped=skipped, git_commits=git_commits and not forced_off, git_commits_forced_off=forced_off, git="none"),
             after={pre + p: exact[p] for p in written if p in exact}, after_paths=[pre + p for p in written if p not in exact], tree=tree)
    run(1, "defaults in a new directory: records under docs/, no .gitignore outside a work tree", "new-project")
    run(2, "docs_dir and wikilinks: records and manifest under knowledge/, wikilink references", "vault", docs="knowledge", wikilinks=True)
    on_disk = [p for p in basic if p in init_files()[0]]
    run(3, "a file that exists is skipped, never overwritten", None, existing=on_disk)
    req = dict(docs_dir="a/b")
    cases.append(("init", 4, dict(name="a docs_dir that is not a single folder name is rejected", repo="basic", request=req,
        expected=err("init", req, "invalid_request", 'docs_dir "a/b" must be a single folder name, with no path separator', parameter="docs_dir"))))
    req = dict(dir="../elsewhere")
    cases.append(("init", 5, dict(name="a directory outside the repository is rejected", repo="basic", request=req,
        expected=err("init", req, "invalid_request", 'dir "../elsewhere" must not contain ".."', parameter="dir"))))
    run(6, "git_commits writes the manifest key and installs the optional skill akm-commit; inside a work tree .gitignore is written", "with-commits", git_commits=True, tree=True)
    run(7, "git_commits stays off outside a git work tree, and the response says so", "no-tree", git_commits=True, forced_off=True)
    run(8, "in a work tree the .gitignore names the lock under the docs dir", "tree-project", docs="knowledge", tree=True)
    req = dict(profile="code")
    cases.append(("init", 9, dict(name="the retired profile parameter is rejected", repo="basic", request=req,
        expected=err("init", req, "invalid_request", 'unknown parameter "profile"', parameter="profile"))))

# config: the git_commits setting. The skill's text is not part of the contract (a case
# lists the skill's files only), and no case removes an installed skill, because whether it
# differs from the shipped copy depends on that text.
def config_setting(value):
    return dict(key="git_commits", value=value, paths=[".agents/skills/akm-commit"],
        description="Install the akm-commit skill, which tells an agent to commit after each logical step. Off: the skill is not installed.")

def config_cases(cases, basic):
    def case(n, name, repo, req, expected, after=None, after_paths=None, tree=True):
        c = dict(name=name, repo=repo, request=req, expected=expected)
        if after: c["after"] = after
        if after_paths: c["after_paths"] = after_paths
        if tree: c["git_work_tree"] = True
        cases.append(("config", n, c))
    def done(req, value, applied, plan, settings_value):
        return ok("config", req, key="git_commits", value=value, changed=bool(plan["write"] or plan["remove"] or plan["modify"]), applied=applied,
                  needs_confirmation=False, plan=plan, edits=[], settings=[config_setting(settings_value)])
    none = dict(write=[], remove=[], modify=[])
    on_plan = dict(write=[COMMIT_SKILL], remove=[], modify=["docs/.layout"])
    case(1, "list the settings", "basic", {}, ok("config", {}, settings=[config_setting(False)]))
    req = dict(key="git_commits")
    case(2, "read one setting", "commits", req, ok("config", req, settings=[config_setting(True)]))
    req = dict(key="git_commits", value=True)
    case(3, "on: the key is written and the skill installed", "basic", req, done(req, True, True, on_plan, True),
         after={"docs/.layout": basic["docs/.layout"] + "git_commits: true\n"}, after_paths=[COMMIT_SKILL])
    req = dict(key="git_commits", value=True, dry_run=True)
    case(4, "dry run: the plan, nothing written", "basic", req, done(req, True, False, on_plan, False))
    req = dict(key="git_commits", value=False)
    case(5, "off with no skill directory: only the key changes, no confirmation", "commits", req,
         done(req, False, True, dict(write=[], remove=[], modify=["docs/.layout"]), False), after={"docs/.layout": basic["docs/.layout"] + "git_commits: false\n"})
    case(6, "off when it is already off changes nothing", "basic", req, done(req, False, False, none, False))
    req = dict(key="git_commits", value=True)
    case(7, "a path that is not the tool's own is a conflict and nothing is written", "conflict", req,
         err("config", req, "refused", "not overwriting what is not the tool's own: .agents/skills/akm-commit", reason="conflict", paths=[".agents/skills/akm-commit"]))
    req = dict(key="color", value=True)
    case(8, "an unknown key is rejected", "basic", req, err("config", req, "invalid_request", 'key: unknown value "color"', parameter="key"))
    req = dict(key="git_commits", value=True)
    case(10, "on outside a git work tree is refused", "basic", req, err("config", req, "refused", "this repository is not inside a git work tree: run git init first; git commits need one", reason="not_a_git_repo"), tree=False)
    req = dict(value=True)
    case(9, "a value needs a key", "basic", req, err("config", req, "invalid_request", "value needs a key", parameter="key"))

def main():
    basic = build("basic"); bro = build("broken", broken)
    recs = load(basic); rank = roadmap_rank(basic); cases = []
    def case(op, n, name, repo, request, expected): cases.append((op, n, dict(name=name, repo=repo, request=request, expected=expected)))
    # ls
    req = dict(type="work"); s = sort_recs([r for r in recs if r["type"] == "work"], ["id"], rank)
    case("ls", 1, "all work items, default order", "basic", req, ok("ls", req, total=6, count=6, records=[pub(r, dflt(r)) for r in s]))
    req = dict(status=["ready", "active"], sort=["roadmap", "-priority", "id"], next=True, fields=["title", "status", "priority"])
    s = sort_recs([r for r in recs if r["type"] == "work" and r["fields"]["status"] in ("ready", "active")], req["sort"], rank)
    case("ls", 2, "open work, roadmap then priority, with next lines", "basic", req, ok("ls", req, total=3, count=3, records=[pub(r, req["fields"], True) for r in s]))
    req = dict(type="adr", status=["accepted"])
    s = [r for r in recs if r["type"] == "adr" and r["fields"]["status"] == "accepted"]
    case("ls", 3, "accepted ADRs", "basic", req, ok("ls", req, total=1, count=1, records=[pub(r, dflt(r)) for r in s]))
    req = dict(parent="261015-c4d5", count_only=True)
    case("ls", 4, "count the children of a parent", "basic", req, ok("ls", req, total=2, count=0, records=[]))
    req = dict(status=["finished"])
    case("ls", 5, "an unknown status is rejected", "basic", req, err("ls", req, "invalid_request", "status: unknown value \"finished\"", parameter="status"))
    # search
    req = dict(text="json")
    hits = []
    for r in recs:
        m = []
        for f in ("title", "summary", "scope"):
            v = r["fields"].get(f)
            if v and "json" in (v if isinstance(v, str) else " ".join(v)).lower():
                lines = r["_text"].split("\n")
                ln = next(i for i, l in enumerate(lines, 1) if l.startswith(f + ": "))
                m.append(dict(field=f, line=ln, snippet=lines[ln - 1]))
        if m: hits.append((len(m), r, m))
    hits.sort(key=lambda h: (-h[0], h[1]["id"]))
    case("search", 1, "title and summary search", "basic", req, ok("search", req, total=len(hits), count=len(hits), records=[dict(pub(r, dflt(r)), matches=m) for _, r, m in hits]))
    # show
    r = next(r for r in recs if r["id"] == "261020-a2b3"); req = dict(id=r["id"], section="next")
    secs = [h.lower() for h in re.findall(r"^## (.+)$", r["_body"], re.M)]
    nxt = re.search(r"## Next\n(.*?)\n\n## ", r["_body"], re.S).group(1)
    case("show", 1, "one section of a work item", "basic", req, ok("show", req, record=pub(r), sections=secs, section=dict(name="next", text=nxt)))
    req = dict(id="261099-zzzz")
    case("show", 2, "a missing id", "basic", req, err("show", req, "not_found", "no record with id 261099-zzzz"))
    # board
    cols = []
    for st in ("active", "blocked", "ready"):
        s = sort_recs([r for r in recs if r["type"] == "work" and r["fields"]["status"] == st], ["roadmap", "-priority", "id"], rank)
        cols.append(dict(status=st, count=len(s), records=[pub(r, dflt(r, ["blocked_by", "cycle", "keep"])) for r in s]))
    tot = {st: sum(1 for r in recs if r["type"] == "work" and r["fields"]["status"] == st) for st in ("raw", "ready", "active", "blocked", "done", "dropped")}
    case("board", 1, "default columns", "basic", {}, ok("board", {}, columns=cols, totals=tot))
    # check
    case("check", 1, "a clean repository", "basic", dict(as_of=AS_OF), ok("check", dict(as_of=AS_OF), clean=True, errors=0, warnings=0, count=0, findings=[], reconcile=dict(state="current")))
    def line_of(files, path, prefix): return next(i for i, l in enumerate(files[path].split("\n"), 1) if l.startswith(prefix))
    f = []
    p = "docs/ROADMAP.md"; f.append(dict(rule="RM1", severity="E", file=p, line=line_of(bro, p, "- [261099-bbbb]"), id=None, message="roadmap entry 261099-bbbb does not resolve to a work item", **{"class": None}))
    p = "docs/work/261020-a2b3-add-json-output-to-ls.md"; f.append(dict(rule="WK6", severity="E", file=p, line=None, id="261020-a2b3", message="status active requires a non-empty ## Next", **{"class": None}))
    p = "docs/work/261022-e5f6-document-the-exit-codes.md"; f.append(dict(rule="WK4", severity="E", file=p, line=line_of(bro, p, "parent:"), id="261022-e5f6", message="parent 261099-aaaa does not resolve to an existing item", **{"class": None}))
    f.sort(key=lambda x: (x["file"], x["line"] or 0, x["rule"]))
    case("check", 2, "three errors", "broken", dict(as_of=AS_OF), ok("check", dict(as_of=AS_OF), clean=False, errors=3, warnings=0, count=3, findings=f, reconcile=dict(state="current")))
    case("check", 3, "one rule family", "broken", dict(rules=["RM"], as_of=AS_OF), ok("check", dict(rules=["RM"], as_of=AS_OF), clean=False, errors=1, warnings=0, count=1, findings=[x for x in f if x["rule"].startswith("RM")], reconcile=dict(state="current")))
    req = dict(target="0.3.0", as_of=AS_OF)
    f4 = [dict(x, **{"class": "manual"}) for x in f]
    case("check", 4, "findings are classified for a target version", "broken", req, ok("check", req, clean=False, errors=3, warnings=0, count=3, findings=f4, reconcile=dict(state="current")))
    req = {}
    case("migrate", 1, "a repository on the current layout has nothing to migrate", "basic", req, ok("migrate", req, mode="report", **{"from": "0.3.0"}, to="0.3.0", state="current", steps=[], findings=[], clean=True))
    # find, when
    r = next(r for r in recs if r["id"] == "261022-e5f6"); req = dict(id=r["id"])
    case("find", 1, "a path from an id", "basic", req, ok("find", req, id=r["id"], type="work", path=r["path"]))
    req = dict(id="261022-e5f6", as_of=AS_OF)
    case("when", 1, "the date of an id", "basic", req, ok("when", req, id="261022-e5f6", date="2026-10-22", weekday=dt.date(2026, 10, 22).strftime("%A"), age_days=(dt.date(2026, 10, 28) - dt.date(2026, 10, 22)).days))
    req = dict(dry_run=True, include_content=True)
    case("index", 1, "the ADR index, dry run", "basic", req, ok("index", req, path="docs/adr/README.md", changed=False, written=False, content=basic["docs/adr/README.md"]))
    for n, repo, files in ((1, "basic", basic), (2, "claims", build("claims", claims_layout))):
        exts = sorted(l.split(":")[0][len("extension."):] for l in files["docs/.layout"].split("\n") if l.startswith("extension."))
        repository = dict(name=repo, layout="0.3.0", extensions=exts, reference_style="id", claim_stale_hours=24, work_dir="docs/work", adr_dir="docs/adr", git_commits=False)
        case("version", n, "the repository this server or tool is working on" if n == 1 else "enabled extensions are listed by name", repo, {},
             ok("version", {}, tool=dict(name="akmctl", version="0.1.0"), api_versions=["1.0"], layout_versions=["0.3.0"], extensions=["claims"], repository=repository))
    vault_cases(case, basic, recs, rank)
    legacy_cases(case)
    write_cases(cases, basic)
    init_cases(cases, basic)
    build("commits", lambda f: f.update({"docs/.layout": f["docs/.layout"] + "git_commits: true\n"}))
    build("conflict", lambda f: f.update({".agents/skills/akm-commit": "a file where the skill directory would go\n"}))
    config_cases(cases, basic)
    shutil.rmtree(os.path.join(HERE, "cases"), ignore_errors=True); os.makedirs(os.path.join(HERE, "cases"))
    for op, n, c in cases:
        json.dump(c, open(os.path.join(HERE, "cases", "%s-%d.json" % (op, n)), "w"), indent=2); open(os.path.join(HERE, "cases", "%s-%d.json" % (op, n)), "a").write("\n")
    print(len(cases), "cases")

main()
