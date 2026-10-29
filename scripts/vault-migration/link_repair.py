#!/usr/bin/env python3
"""Emulate Obsidian's link repair for a list of file renames, then perform the renames.

  link_repair.py --root <git root> --renames renames.json [--report out.json] [--no-rename]

renames.json: [{"from": "<repo-relative path>", "to": "<repo-relative path>"}, ...] (the `renames` the
`mint_ids` step prints). Generic: nothing here knows about tasks or ADRs.

What it does, as observed in the Obsidian pilot (decision log Y42):
  - every [[old name]] becomes [[new name]]; heading/block parts, aliases ([[old|text]]), embeds and
    quoted frontmatter links keep their form; in table cells the \\| escape stays;
  - Markdown links [text](old%20name.md) are rewritten the same way, keeping their folder part;
  - links inside fenced code blocks and inline code are not links and are left alone;
  - a file is rewritten only if a link changed. Bytes outside the link spans are not touched, so
    Shavian text or anything else passes through. Paths and counts only are reported.
Then each file is moved from `from` to `to` (os.rename; git sees a rename).
"""
import argparse, json, os, sys, re
from urllib.parse import quote
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import akmlinks as L

def new_wiki_target(target, new_path):
    """Replace the last segment of a wikilink target, keeping an explicit path prefix and .md ext."""
    t = target.decode("utf-8")
    ext = ".md" if t.lower().endswith(".md") else ""
    stem = t[:-3] if ext else t
    prefix = stem.rsplit("/", 1)[0] + "/" if "/" in stem else ""
    base = new_path.rsplit("/", 1)[-1][:-3]
    return (prefix + base + ext).encode("utf-8")

def new_md_dest(dest, new_path):
    d = dest.decode("utf-8")
    angle = d.startswith("<") and d.endswith(">")
    inner = d[1:-1] if angle else d
    encoded = "%20" in inner
    head, _, anchor = inner.partition("#")
    pre = head.rsplit("/", 1)[0] + "/" if "/" in head else ""
    base = new_path.rsplit("/", 1)[-1]
    nb = quote(base, safe="()'!*~._-") if encoded else base
    out = pre + nb + (("#" + anchor) if anchor else "")
    return (("<" + out + ">") if angle else out).encode("utf-8")

def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--root", required=True)
    ap.add_argument("--renames", required=True)
    ap.add_argument("--report")
    ap.add_argument("--no-rename", action="store_true", help="repair links but do not move the files")
    a = ap.parse_args()
    root = a.root
    renames = json.load(open(a.renames, encoding="utf-8"))
    rmap = {r["from"]: r["to"] for r in renames}
    paths = L.list_md(root)
    idx = L.Index(paths)
    for r in renames:
        if r["from"] not in idx.paths:
            sys.exit("rename source missing: " + r["from"])
        if os.path.exists(os.path.join(root, r["to"])):
            sys.exit("rename target exists: " + r["to"])
    changed, ambiguous, total = {}, [], 0
    for p in paths:
        fp = os.path.join(root, p)
        b = open(fp, "rb").read()
        edits = []
        for lk in L.find_links(b):
            if lk["kind"] == "wiki":
                kind, res = idx.resolve_wiki(lk["target"].decode("utf-8", "replace"), p)
                if kind == "ambiguous":
                    if any(x in rmap for x in res):
                        ambiguous.append({"file": p, "candidates": res})
                    continue
                if kind != "ok" or res not in rmap:
                    continue
                edits.append((lk["tstart"], lk["tend"], new_wiki_target(lk["target"], rmap[res])))
            else:
                res = idx.resolve_md(lk["target"].decode("utf-8", "replace"), p)
                if res is None or res not in rmap:
                    continue
                edits.append((lk["tstart"], lk["tend"], new_md_dest(lk["target"], rmap[res])))
        if edits:
            for s, e, new in sorted(edits, reverse=True):
                b = b[:s] + new + b[e:]
            tmp = fp + ".tmp-linkrepair"
            open(tmp, "wb").write(b)
            os.replace(tmp, fp)
            changed[p] = {"links": len(edits), "shavian": L.has_shavian(b)}
            total += len(edits)
    moved = 0
    if not a.no_rename:
        for r in renames:
            os.rename(os.path.join(root, r["from"]), os.path.join(root, r["to"]))
            moved += 1
    rep = {"renames": len(renames), "moved": moved, "files_changed": len(changed), "links_changed": total,
           "ambiguous": ambiguous, "changed": changed}
    if a.report:
        json.dump(rep, open(a.report, "w", encoding="utf-8"), indent=1, ensure_ascii=False)
    print(f"renames {len(renames)}, moved {moved}, files with links rewritten {len(changed)}, links rewritten {total}, "
          f"ambiguous {len(ambiguous)}, Shavian files touched {sum(1 for v in changed.values() if v['shavian'])}")

if __name__ == "__main__":
    main()
