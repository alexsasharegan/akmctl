#!/usr/bin/env python3
"""Generate the one-off AKM migration file (0.0.0 -> 1.0.0) for the vault.

Reads the task notes' `status` to emit one move_file step per done/dropped task.
Standard library only. Usage: make_migration.py <repo-root> > vault-0.0.0-to-1.0.0.json
"""
import json, re, sys, os

root = sys.argv[1]
V = "Deep Thought/9 pkm"
TASKS, ADR = f"{V}/tasks", f"{V}/adr"

def status_of(path):
    with open(path, encoding="utf-8") as f:
        text = f.read()
    if not text.startswith("---\n"):
        return None
    end = text.find("\n---", 4)
    m = re.search(r"^status:\s*\"?([a-z]+)\"?\s*$", text[4:end], re.M)
    return m.group(1) if m else None

steps = []
def add(id, kind, desc, **kw):
    s = {"id": id, "kind": kind, "class": "auto", "description": desc}
    s.update(kw)
    steps.append(s)

# 1. The ADR conventions leave the ADR folder, so adr/README.md can become the generated index.
add("conventions", "move_file", "ADR conventions note moves beside the design notes (adr/README.md becomes the generated index).",
    **{"from": f"{ADR}/README.md", "to": f"{V}/design/adr-conventions.md"})
add("candidates", "move_file", "ADR Candidates moves beside the conventions note.",
    **{"from": f"{ADR}/ADR Candidates.md", "to": f"{V}/design/ADR Candidates.md"})

# 2. Closed tasks go to done/ (WK9).
closed = []
for n in sorted(os.listdir(os.path.join(root, TASKS))):
    p = os.path.join(root, TASKS, n)
    if n.endswith(".md") and os.path.isfile(p) and status_of(p) in ("done", "dropped"):
        closed.append(n)
for i, n in enumerate(closed, 1):
    add(f"done-{i:03d}", "move_file", "Closed task moves to done/.",
        **{"from": f"{TASKS}/{n}", "to": f"{TASKS}/done/{n}"})

# 3. Tasks: fields, title, goal.
add("task-created", "rename_field", "created_at -> created", dir=TASKS, recursive=True, **{"from": "created_at", "to": "created"})
add("task-updated", "rename_field", "updated_at -> updated", dir=TASKS, recursive=True, **{"from": "updated_at", "to": "updated"})
add("task-title", "title_from_h1", "title from the first heading", dir=TASKS, recursive=True)
add("task-goal", "goal_scaffold", "mechanical ## Goal scaffold (review the open ones)", dir=TASKS, recursive=True)
add("task-ids", "mint_ids", "ids from created; files NOT renamed (Obsidian renames them)", dir=TASKS, recursive=True,
    date_field="created", rename=False, seed="vault-tasks-2026-10")

# 4. ADRs: the numeric id first, then ids, date, and the same date field names as tasks.
add("adr-number", "rename_field", "numeric id -> adr_number", dir=ADR, **{"from": "id", "to": "adr_number"})
add("adr-date", "copy_field", "date = created_at", dir=ADR, **{"from": "created_at", "to": "date"})
add("adr-ids", "mint_ids", "ids from created_at; files NOT renamed (Obsidian renames them)", dir=ADR,
    date_field="created_at", rename=False, seed="vault-adr-2026-10")
add("adr-created", "rename_field", "created_at -> created", dir=ADR, **{"from": "created_at", "to": "created"})
add("adr-updated", "rename_field", "updated_at -> updated", dir=ADR, **{"from": "updated_at", "to": "updated"})

json.dump({"from": "0.0.0", "to": "1.0.0", "steps": steps}, sys.stdout, indent=2, ensure_ascii=False)
print()
