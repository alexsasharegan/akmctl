#!/usr/bin/env python3
"""Write the drafted ADR summaries with `akmctl set`. Usage: apply_summaries.py <akmctl> <repo> <adr-dir> adr_summaries.json"""
import json, os, re, subprocess, sys
akm, repo, adr_dir, summ = sys.argv[1:5]
s = json.load(open(summ, encoding="utf-8"))
for n in sorted(os.listdir(os.path.join(repo, adr_dir))):
    if not n.endswith(".md") or n == "README.md":
        continue
    t = open(os.path.join(repo, adr_dir, n), encoding="utf-8").read()
    fm = t[4:t.find("\n---", 4)]
    num = re.search(r"^adr_number:\s*(\d+)", fm, re.M)
    id_ = n[:11]
    if num and num.group(1) in s:
        r = subprocess.run([akm, "set", id_, "summary=" + s[num.group(1)], "--repo", repo], capture_output=True, text=True)
        if r.returncode:
            sys.exit(f"{id_}: {r.stderr or r.stdout}")
print("summaries written")
