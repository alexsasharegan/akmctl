#!/usr/bin/env python3
"""Rewrite bare prose mentions of legacy ADR numbers ("ADR-019", "ADR 019") to the id form ("ADR 261027-ensf").

  adr_prose.py --root <git root> --adr-dir "<dir>" --include <prefix> [--include ...] [--report out.json] [--apply]

The number -> id map is read from the ADR files themselves (`adr_number` and `id` in the frontmatter), so it
runs after the migration. Not rewritten: link spans ([[...]], Markdown links), fenced code and inline code,
file-name forms (ADR-019.md, a path segment), files under --exclude (another ADR series with the same ADR-NNN look), numbers that are not in the map (e.g. ADR-042, an example),
and any file with Shavian text (counted and listed instead). Files outside --include are only counted.
Default is a dry run. Bytes outside the matched mentions are unchanged.
"""
import argparse, json, os, re, sys
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import akmlinks as L

MENTION = re.compile(rb"(?<![\w/])ADR[- ](\d{1,3})(?![\d\w]|\.md|/)")

def adr_map(root, adr_dir):
    m = {}
    d = os.path.join(root, adr_dir)
    for n in sorted(os.listdir(d)):
        if not n.endswith(".md") or n == "README.md":
            continue
        t = open(os.path.join(d, n), encoding="utf-8").read()
        end = t.find("\n---", 4)
        fm = t[4:end]
        num = re.search(r"^adr_number:\s*(\d+)\s*$", fm, re.M)
        id_ = re.search(r"^id:\s*\"?(\d{6}-[0-9a-hjkmnp-tv-z]{4})\"?\s*$", fm, re.M)
        if num and id_:
            m[int(num.group(1))] = id_.group(1)
    return m

def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--root", required=True)
    ap.add_argument("--adr-dir", required=True)
    ap.add_argument("--include", action="append", required=True)
    ap.add_argument("--exclude", action="append", default=[], help="path prefix whose mentions belong to another ADR series; only counted, never listed by file")
    ap.add_argument("--shavian-ok", action="append", default=[], help="path of a Shavian-containing file whose ADR mentions may be rewritten (the match is ASCII-only, so only the mention bytes change)")
    ap.add_argument("--report")
    ap.add_argument("--apply", action="store_true")
    a = ap.parse_args()
    amap = adr_map(a.root, a.adr_dir)
    rep = {"excluded_other_series": {}, "map_size": len(amap), "rewritten": {}, "outside_include": {}, "shavian_skipped": {}, "unmapped": {}, "in_code_or_link": {}}
    total = 0
    for p in L.list_md(a.root):
        fp = os.path.join(a.root, p)
        b = open(fp, "rb").read()
        if not MENTION.search(b):
            continue
        skip = [(l["start"], l["end"]) for l in L.find_links(b)] + L.code_spans(b)
        skip.sort()
        ex = next((x for x in a.exclude if p.startswith(x.rstrip("/") + "/")), None)
        if ex:
            rep["excluded_other_series"][ex] = rep["excluded_other_series"].get(ex, 0) + len(MENTION.findall(b))
            continue
        edits, unmapped, hidden = [], 0, 0
        for m in MENTION.finditer(b):
            if L.in_spans(m.start(), skip) or any(s <= m.start() < e for s, e in skip):
                hidden += 1
                continue
            n = int(m.group(1))
            if n not in amap:
                unmapped += 1
                continue
            edits.append((m.start(), m.end(), b"ADR " + amap[n].encode()))
        if hidden:
            rep["in_code_or_link"][p] = hidden
        if unmapped:
            rep["unmapped"][p] = unmapped
        if not edits:
            continue
        inc = any(p.startswith(x.rstrip("/") + "/") or p == x for x in a.include)
        if not inc:
            rep["outside_include"][p] = len(edits)
        elif L.has_shavian(b) and p not in a.shavian_ok:
            rep["shavian_skipped"][p] = len(edits)
        else:
            rep["rewritten"][p] = len(edits)
            total += len(edits)
            if a.apply:
                for s, e, new in sorted(edits, reverse=True):
                    b = b[:s] + new + b[e:]
                open(fp, "wb").write(b)
    rep["total_rewritten"] = total
    if a.report:
        json.dump(rep, open(a.report, "w", encoding="utf-8"), indent=1, ensure_ascii=False)
    print(("applied" if a.apply else "dry run"), "rewritten", total, "in", len(rep["rewritten"]), "files; outside include",
          sum(rep["outside_include"].values()), "in", len(rep["outside_include"]), "files; Shavian files skipped", len(rep["shavian_skipped"]),
          "; unmapped", sum(rep["unmapped"].values()), "; inside code/link", sum(rep["in_code_or_link"].values()))

main()
