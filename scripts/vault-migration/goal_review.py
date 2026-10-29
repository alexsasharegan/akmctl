#!/usr/bin/env python3
"""List the open tasks (ready, active, blocked) whose ## Goal came from goal_scaffold, with simple flags.

  goal_review.py <task-dir> > table (Markdown)
Flags: EMPTY (no text under ## Goal), SHORT (under 80 characters), LONG (over 1500), NO-DONE-WHEN (no
"Done when" line or list in the Goal), GOAL-IS-LOG (text starts with a dated log line).
Reads note text only to measure it; nothing is quoted.
"""
import os, re, sys
d = sys.argv[1]
rows = []
for n in sorted(os.listdir(d)):
    if not n.endswith(".md"):
        continue
    t = open(os.path.join(d, n), encoding="utf-8").read()
    fm_end = t.find("\n---", 4)
    fm = t[4:fm_end]
    st = re.search(r"^status:\s*(\S+)", fm, re.M).group(1)
    if st not in ("ready", "active", "blocked"):
        continue
    title = re.search(r"^title:\s*(.*)$", fm, re.M).group(1).strip().strip('"')
    m = re.search(r"^## Goal\n(.*?)(?=^## |\Z)", t[fm_end:], re.M | re.S)
    g = m.group(1).strip() if m else ""
    flags = []
    if not g: flags.append("EMPTY")
    elif len(g) < 80: flags.append("SHORT")
    if len(g) > 1500: flags.append("LONG")
    if not re.search(r"done when", g, re.I): flags.append("NO-DONE-WHEN")
    if re.match(r"- \d{4}-\d\d-\d\d", g): flags.append("GOAL-IS-LOG")
    rows.append((n[:11], st, title, len(g), " ".join(flags) or "-"))
print("| id | status | title | Goal chars | flags |\n| --- | --- | --- | --- | --- |")
for r in rows:
    print("| %s | %s | %s | %d | %s |" % (r[0], r[1], r[2].replace("|", "\\|"), r[3], r[4]))
print("\n%d open tasks; %d with at least one flag other than NO-DONE-WHEN" % (len(rows), sum(1 for r in rows if set(r[4].split()) - {"NO-DONE-WHEN", "-"})), file=sys.stderr)
