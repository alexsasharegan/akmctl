#!/usr/bin/env python3
"""Replace the slug part of the ADR renames with one made from each ADR's title.

mint_ids builds the slug from the old file name, which for ADR-019.md gives `adr-019`; a slug from the
title reads better (`akm-an-agent-knowledge-management`). Usage:
  adr_slugs.py --root <git root> --adr-dir "<repo-relative dir>" --renames renames.json [--max 40]
Rewrites renames.json in place, for the entries whose `to` is directly in --adr-dir.
"""
import argparse, json, os, re, sys, unicodedata

def slugify(title, n):
    s = unicodedata.normalize("NFKD", title).encode("ascii", "ignore").decode().lower()
    words = re.findall(r"[a-z0-9]+", s)
    out = ""
    for w in words:
        nxt = (out + "-" + w) if out else w
        if len(nxt) > n:
            break
        out = nxt
    stop = {"a", "an", "and", "the", "of", "to", "in", "for", "with", "as", "over", "on", "by", "after", "from", "no", "is", "at", "or", "then"}
    parts = out.split("-")
    while len(parts) > 1 and parts[-1] in stop:
        parts.pop()
    return "-".join(parts) or "untitled"

def title_of(path):
    text = open(path, encoding="utf-8").read()
    end = text.find("\n---", 4)
    m = re.search(r"^title:\s*(.*)$", text[4:end], re.M)
    t = m.group(1).strip() if m else ""
    if len(t) >= 2 and t[0] == t[-1] and t[0] in "\"'":
        t = t[1:-1].replace('\\"', '"').replace("''", "'")
    return t

ap = argparse.ArgumentParser()
ap.add_argument("--root", required=True)
ap.add_argument("--adr-dir", required=True)
ap.add_argument("--renames", required=True)
ap.add_argument("--max", type=int, default=40)
a = ap.parse_args()
r = json.load(open(a.renames, encoding="utf-8"))
d = a.adr_dir.rstrip("/")
for x in r:
    if x["to"].rsplit("/", 1)[0] == d:
        base = x["to"].rsplit("/", 1)[1]
        id_ = re.match(r"^(\d{6}-[0-9a-hjkmnp-tv-z]{4})-", base).group(1)
        t = title_of(os.path.join(a.root, x["from"]))
        x["to"] = f"{d}/{id_}-{slugify(t, a.max)}.md"
json.dump(r, open(a.renames, "w", encoding="utf-8"), indent=1, ensure_ascii=False)
