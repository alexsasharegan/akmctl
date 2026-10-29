#!/usr/bin/env python3
"""Rename the files of a renames list through the Obsidian CLI, so Obsidian repairs the links.

  obsidian_renames.py --renames renames.json [--vault "Deep Thought"] [--vault-prefix "Deep Thought/"] [--dry-run] [--log run.log]

For each {"from", "to"} it runs
  obsidian --vault="<vault>" rename file="<old name without .md>" name="<new file name>.md"
and checks on disk that `from` is gone and `to` exists before going on. It stops at the first failure
(nothing is rolled back; the entries done so far are in the log, and a rerun skips entries whose `to`
exists and `from` does not). The git tree must be committed before this runs, so `git diff` shows exactly
what Obsidian changed. Obsidian must be running. NOT used in the dry run (it needs the real vault).
Paths in renames.json are repository-relative; the vault folder is the first path segment (--vault-prefix).
Standard library only.
"""
import argparse, json, os, subprocess, sys, time

ap = argparse.ArgumentParser()
ap.add_argument("--renames", required=True)
ap.add_argument("--root", default=".", help="repository root (default: current directory)")
ap.add_argument("--vault", default="Deep Thought")
ap.add_argument("--dry-run", action="store_true")
ap.add_argument("--log")
ap.add_argument("--pause", type=float, default=0.3, help="seconds to wait after each rename")
a = ap.parse_args()
renames = json.load(open(a.renames, encoding="utf-8"))
log = open(a.log, "a", encoding="utf-8") if a.log else None
done = skipped = 0
for i, r in enumerate(renames, 1):
    old, new = os.path.join(a.root, r["from"]), os.path.join(a.root, r["to"])
    if not os.path.exists(old) and os.path.exists(new):
        skipped += 1
        continue
    stem = os.path.basename(r["from"])[:-3]
    cmd = ["obsidian", "--vault=" + a.vault, "rename", "file=" + stem, "name=" + os.path.basename(r["to"])]
    if a.dry_run:
        print(" ".join(repr(c) if " " in c else c for c in cmd))
        continue
    p = subprocess.run(cmd, capture_output=True, text=True)
    time.sleep(a.pause)
    ok = os.path.exists(new) and not os.path.exists(old)
    line = f"{i}/{len(renames)} rc={p.returncode} ok={ok} {r['from']} -> {r['to']}"
    print(line)
    if log:
        log.write(line + "\n" + p.stdout + p.stderr + "\n"); log.flush()
    if p.returncode or not ok:
        sys.exit(f"stopped at {i}: {p.stdout.strip()} {p.stderr.strip()}")
    done += 1
print(f"renamed {done}, already done {skipped}, total {len(renames)}")
