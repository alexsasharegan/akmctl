#!/usr/bin/env python3
"""Fold block-scalar values (key: >, >-, |, |-) in SKILL.md frontmatter into one double-quoted line.

akmctl rule S2 reads a skill's frontmatter with the YAML-lite parser, which rejects block scalars.
Usage: fold_skill_frontmatter.py <SKILL.md> ...   (edits in place; prints the paths it changed)
The text is joined with single spaces (folded style), which is what `>` / `>-` mean to a YAML reader.
"""
import re, sys

for p in sys.argv[1:]:
    t = open(p, encoding="utf-8").read()
    end = t.index("\n---", 4)
    fm, rest = t[4:end], t[end:]
    out, lines, i = [], fm.split("\n"), 0
    changed = False
    while i < len(lines):
        m = re.match(r"^([A-Za-z_][\w.-]*):\s*([>|][+-]?)\s*$", lines[i])
        if not m:
            out.append(lines[i]); i += 1; continue
        j, parts = i + 1, []
        while j < len(lines) and (lines[j].startswith(" ") or lines[j].strip() == ""):
            parts.append(lines[j].strip()); j += 1
        text = " ".join(x for x in parts if x) if m.group(2).startswith(">") else "\n".join(parts)
        if "\n" in text:
            sys.exit(p + ": literal block with line breaks; fold by hand")
        out.append('%s: "%s"' % (m.group(1), text.replace("\\", "\\\\").replace('"', '\\"')))
        i, changed = j, True
    if changed:
        open(p, "w", encoding="utf-8").write("---\n" + "\n".join(out) + rest)
        print(p)
