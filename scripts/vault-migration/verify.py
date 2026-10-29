#!/usr/bin/env python3
"""Verify the link repair and the ADR prose rewrite against the tree before the migration.

  verify.py --repo <git root> --before <ref> --after <ref> --migration m.json --renames renames.json
            [--tool-dirs "<dir>" ...] [--out report.json] [--final <ref>]

  before / after   git refs (tag, commit) holding the trees to compare; they are exported with
                   `git archive` into temporary folders and removed afterwards.
  tool-dirs        folders whose notes the migration itself edits (tasks, ADRs): they get a weaker
                   check (b'). Everything else must equal its original after masking.
  final            optional ref of the end state, for check (c) only.

Checks
  (a) link multiset: (source, resolved target) before and after, with both sides mapped to final paths.
  (b) masked diff: link targets and ADR number mentions are masked in both texts; every file outside
      tool-dirs must then be byte-equal. (b') for tool-dirs notes: differences are allowed only in the
      frontmatter block, a new `## Goal` line, blank lines, and `## Done when` becoming `Done when:`.
  (c) no unresolved link added (source, raw target) in after (or final) beyond those already in before.
  (d) Shavian files changed: path and number of changed links.
  (e) changed files in 1 Calendar, 4 RESOURCES, 5 ARCHIVES.
Prints counts and paths only; never file content. Exit 1 when (a), (b), (b') or (c) fail.
"""
from collections import Counter
import argparse, difflib, json, os, re, shutil, subprocess, sys, tempfile
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import akmlinks as L

MENT = re.compile(rb"(?<![\w/])ADR[- ](\d{1,3}(?![\d\w]|\.md|/)|\d{6}-[0-9a-hjkmnp-tv-z]{4}\b)")

def export(repo, ref, dest):
    p = subprocess.Popen(["git", "-C", repo, "archive", ref], stdout=subprocess.PIPE)
    subprocess.run(["tar", "-x", "-C", dest], stdin=p.stdout, check=True)
    p.wait()

def masked(b, side):
    """Replace link targets and ADR mentions by placeholders; return masked bytes."""
    links = L.find_links(b)
    spans = L.code_spans(b)
    edits = []
    for lk in links:
        edits.append((lk["tstart"], lk["tend"], b"\x00T\x00"))
    skip = [(l["start"], l["end"]) for l in links] + spans
    for m in MENT.finditer(b):
        if any(s <= m.start() < e for s, e in skip):
            continue
        edits.append((m.start(), m.end(), b"\x00A\x00"))
    for s, e, new in sorted(edits, reverse=True):
        b = b[:s] + new + b[e:]
    return b

def link_records(root, paths, idx, pm):
    """Counter of (source, target-or-marker) for every link."""
    from collections import Counter
    c, unres = Counter(), Counter()
    for p in paths:
        b = open(os.path.join(root, p), "rb").read()
        src = pm.get(p, p)
        for lk in L.find_links(b):
            if lk["kind"] == "wiki":
                kind, res = idx.resolve_wiki(lk["target"].decode("utf-8", "replace"), p)
                if kind == "ok" or kind == "self":
                    c[(src, pm.get(res, res))] += 1
                elif kind == "ambiguous":
                    c[(src, "AMBIGUOUS:" + "|".join(sorted(pm.get(x, x) for x in res)))] += 1
                elif kind == "none":
                    unres[(src, norm_raw(res))] += 1
            else:
                res = idx.resolve_md(lk["target"].decode("utf-8", "replace"), p)
                if res:
                    c[(src, pm.get(res, res))] += 1
    return c, unres

def norm_raw(t):
    return os.path.basename(t.replace("\\", "/")).lower()

def fm_split(b):
    e = L.frontmatter_end(b)
    return b[:e], b[e:]

def tool_note_ok(a, bb):
    """(b') compare two masked texts; True if only allowed differences."""
    fa, ra = fm_split(a)
    fb, rb = fm_split(bb)
    la, lb = ra.split(b"\n"), rb.split(b"\n")
    from collections import Counter as C
    gone = C(x.strip() for x in la if x.strip()) - C(x.strip() for x in lb if x.strip())
    came = C(x.strip() for x in lb if x.strip()) - C(x.strip() for x in la if x.strip())
    if set(gone) <= {b"## Done when"} and set(came) <= {b"## Goal", b"Done when:"} and sum(gone.values()) <= 1 and sum(came.values()) <= 2:
        return True, 0  # the Goal scaffold: same lines, a new heading, `## Done when` moved into the Goal
    sm = difflib.SequenceMatcher(None, la, lb, autojunk=False)
    for tag, i1, i2, j1, j2 in sm.get_opcodes():
        if tag == "equal":
            continue
        old, new = la[i1:i2], lb[j1:j2]
        if tag == "insert" and all(x.strip() in (b"", b"## Goal") for x in new):
            continue
        if tag == "replace" and [x.strip() for x in old] == [b"## Done when"] and all(x.strip() in (b"", b"Done when:") for x in new):
            continue
        if tag == "replace" and all(x.strip() in (b"", b"## Done when") for x in old) and all(x.strip() in (b"", b"Done when:") for x in new):
            continue
        if tag == "delete" and all(x.strip() == b"" for x in old):
            continue
        return False, i1 + 1
    return True, 0

def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--repo", required=True)
    ap.add_argument("--before", required=True)
    ap.add_argument("--after", required=True)
    ap.add_argument("--final")
    ap.add_argument("--migration", required=True)
    ap.add_argument("--renames", required=True)
    ap.add_argument("--tool-dirs", nargs="*", default=[])
    ap.add_argument("--exempt", help="file with one repo-relative path per line: notes edited by hand after the scripted steps (reported, not checked in (b); their link differences are listed apart)")
    ap.add_argument("--out")
    a = ap.parse_args()
    exempt = set(l.rstrip("\n") for l in open(a.exempt, encoding="utf-8")) if a.exempt else set()

    mig = json.load(open(a.migration, encoding="utf-8"))
    mig = mig[0] if isinstance(mig, list) else mig
    moves = {s["from"]: s["to"] for s in mig["steps"] if s["kind"] == "move_file"}
    renames = {r["from"]: r["to"] for r in json.load(open(a.renames, encoding="utf-8"))}
    # original path -> path after the migration: move first, then rename
    def path_map(paths):
        pm = {}
        for p in paths:
            q = moves.get(p, p)
            q = renames.get(q, q)
            if q != p:
                pm[p] = q
        return pm

    tmp = tempfile.mkdtemp(prefix="akm-verify-")
    try:
        dirs = {}
        for name, ref in (("before", a.before), ("after", a.after)) + ((("final", a.final),) if a.final else ()):
            d = os.path.join(tmp, name)
            os.makedirs(d)
            export(a.repo, ref, d)
            dirs[name] = d

        def mdfiles(d):
            out = []
            for root, ds, fs in os.walk(d):
                ds[:] = [x for x in ds if x != ".git"]
                for f in fs:
                    if f.lower().endswith(".md"):
                        out.append(os.path.relpath(os.path.join(root, f), d).replace(os.sep, "/"))
            return sorted(out)

        pb, pa = mdfiles(dirs["before"]), mdfiles(dirs["after"])
        pm = path_map(pb)
        ib, ia = L.Index(pb), L.Index(pa)
        cb, ub = link_records(dirs["before"], pb, ib, pm)
        ca, ua = link_records(dirs["after"], pa, ia, {})
        res = {"files_before": len(pb), "files_after": len(pa)}
        # (a)
        only_b, only_a = cb - ca, ca - cb
        res["a_links_before"], res["a_links_after"] = sum(cb.values()), sum(ca.values())
        res["a_only_before"] = [[s, t, n] for (s, t), n in sorted(only_b.items()) if s not in exempt]
        res["a_only_after"] = [[s, t, n] for (s, t), n in sorted(only_a.items()) if s not in exempt]
        res["a_exempt_sources"] = sorted({s for (s, t) in list(only_b) + list(only_a) if s in exempt})
        res["a_ok"] = not res["a_only_before"] and not res["a_only_after"]
        # (c)
        newun = Counter({k: v for k, v in (ua - ub).items() if k[0] not in exempt})
        res["c_unresolved_before"], res["c_unresolved_after"] = sum(ub.values()), sum(ua.values())
        res["c_added"] = [[s, t, n] for (s, t), n in sorted(newun.items())]
        res["c_ok"] = not newun
        if a.final:
            pf = mdfiles(dirs["final"])
            cf, uf = link_records(dirs["final"], pf, L.Index(pf), {})
            nf = uf - ub
            res["c_final_added"] = [[s, t, n] for (s, t), n in sorted(nf.items())]
        # (b)
        tool_dirs = [d.rstrip("/") + "/" for d in a.tool_dirs]
        bad, badtool, checked, checked_tool, changed = [], [], 0, 0, []
        for p in pb:
            q = pm.get(p, p)
            if q not in ia.paths:
                res.setdefault("missing_after", []).append(p)
                continue
            b0 = open(os.path.join(dirs["before"], p), "rb").read()
            b1 = open(os.path.join(dirs["after"], q), "rb").read()
            if b0 != b1:
                changed.append(q)
            if q in exempt or p in exempt:
                res.setdefault("b_exempt", []).append(q)
                continue
            is_tool = any(q.startswith(d) for d in tool_dirs) or any(p.startswith(d) for d in tool_dirs)
            m0, m1 = masked(b0, "before"), masked(b1, "after")
            if is_tool:
                checked_tool += 1
                ok, line = tool_note_ok(m0, m1)
                if not ok:
                    badtool.append([q, line])
            else:
                checked += 1
                if m0 != m1:
                    l0, l1 = m0.split(b"\n"), m1.split(b"\n")
                    line = next((i + 1 for i, (x, y) in enumerate(zip(l0, l1)) if x != y), min(len(l0), len(l1)) + 1)
                    bad.append([q, line])
        res["b_checked_strict"], res["b_checked_tool_notes"] = checked, checked_tool
        res["b_bad"], res["b_bad_tool"] = bad, badtool
        res["b_ok"] = not bad and not badtool
        res["new_files"] = sorted(set(pa) - {pm.get(p, p) for p in pb})
        # (d), (e): per changed file, count of rewritten links = links whose target text differs
        rep = {}
        for q in changed:
            p = next((k for k, v in pm.items() if v == q), q)
            if p not in ib.paths:
                continue
            l0 = [(lk["target"]) for lk in L.find_links(open(os.path.join(dirs["before"], p), "rb").read())]
            b1 = open(os.path.join(dirs["after"], q), "rb").read()
            l1 = [(lk["target"]) for lk in L.find_links(b1)]
            n = sum(1 for x, y in zip(l0, l1) if x != y) if len(l0) == len(l1) else -1
            rep[q] = {"links": n, "shavian": L.has_shavian(b1)}
        res["changed_files"] = len(changed)
        res["d_shavian_changed"] = {k: v["links"] for k, v in sorted(rep.items()) if v["shavian"]}
        res["e_calendar_resources_archives"] = {k: v["links"] for k, v in sorted(rep.items())
                                                 if k.startswith(("Deep Thought/1 Calendar/", "Deep Thought/4 RESOURCES/", "Deep Thought/5 ARCHIVES/"))}
        res["amex_area_changed"] = sum(1 for k in rep if "/3 AREAS/AMEX/" in k)
    finally:
        shutil.rmtree(tmp, ignore_errors=True)
    if a.out:
        json.dump(res, open(a.out, "w", encoding="utf-8"), indent=1, ensure_ascii=False)
    print("(a) link multiset equal: %s (before %d, after %d; only-before %d, only-after %d)" % (res["a_ok"], res["a_links_before"], res["a_links_after"], len(res["a_only_before"]), len(res["a_only_after"])))
    print("(b) masked diff ok: %s (strict %d files, %d bad; tool notes %d, %d bad)" % (res["b_ok"], res["b_checked_strict"], len(bad), res["b_checked_tool_notes"], len(badtool)))
    print("(c) no unresolved link added: %s (before %d, after %d, added %d)" % (res["c_ok"], res["c_unresolved_before"], res["c_unresolved_after"], len(res["c_added"])))
    print("(d) Shavian files changed: %d (links changed %d)" % (len(res["d_shavian_changed"]), sum(max(0, v) for v in res["d_shavian_changed"].values())))
    print("(e) changed in 1 Calendar / 4 RESOURCES / 5 ARCHIVES: %d files" % len(res["e_calendar_resources_archives"]))
    print("changed files total %d; new files %d; missing after %d" % (res["changed_files"], len(res["new_files"]), len(res.get("missing_after", []))))
    sys.exit(0 if (res["a_ok"] and res["b_ok"] and res["c_ok"]) else 1)

main()
