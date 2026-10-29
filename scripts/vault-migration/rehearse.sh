#!/bin/bash
# Dry-run harness: copies a vault repository and runs the scripted migration steps on the copy,
# with the link repair emulated by link_repair.py instead of Obsidian. Prints a timing per step.
#   rehearse.sh <source repo root> <destination dir (must not exist)> <akmctl binary>
# Never run it with the real repository as the destination. Needs git, rsync, python3.
set -euo pipefail
SRC=$(cd "$1" && pwd); DEST=$2; AKM=$3
HERE=$(cd "$(dirname "$0")" && pwd)
V="Deep Thought"
WORK=$(mktemp -d)
[ ! -e "$DEST" ] || { echo "destination exists: $DEST" >&2; exit 2; }
case "$DEST/" in "$SRC"/*) echo "destination is inside the source" >&2; exit 2;; esac
T0=$(python3 -c 'import time;print(time.time())'); LAST=$T0
step() { local now; now=$(python3 -c 'import time;print(time.time())'); printf '%-34s %6.2fs\n' "$1" "$(python3 -c "print($now-$LAST)")"; LAST=$now; }

rsync -a --exclude node_modules --exclude /.bak --exclude /.tmp "$SRC"/ "$DEST"/ ; step "copy"
cd "$DEST"
git config user.name dryrun; git config user.email dryrun@example.invalid; git config core.hooksPath /dev/null
git add -A; git diff --cached --quiet || git commit -qm "dry run baseline"; step "baseline commit"

mkdir -p docs
printf 'work_dir: %s/9 pkm/tasks\nadr_dir: %s/9 pkm/adr\nroadmap: %s/9 pkm/design/akm-roadmap.md\n' "$V" "$V" "$V" > docs/.layout
printf -- '---\ntype: roadmap\nstatus: active\ncreated: 2026-10-29\nupdated: 2026-10-29\nauthored_by: agent\n---\n# Roadmap\n\n## Now\n\n## Next\n\n## Later\n\n## Dropped\n' > "$V/9 pkm/design/akm-roadmap.md"
python3 "$HERE/make_migration.py" . > "$WORK"/rehearse-migration.json; step "manifest, roadmap, migration file"
git add -A; git commit -qm "dry run: manifest and roadmap"

"$AKM" migrate --mode report --migration "$WORK"/rehearse-migration.json --format json > "$WORK"/rehearse-report.json || true; step "migrate report"
"$AKM" migrate --mode plan --migration "$WORK"/rehearse-migration.json --format json > "$WORK"/rehearse-plan.json || true; step "migrate plan"
"$AKM" migrate --mode apply-safe --acknowledge-backup --migration "$WORK"/rehearse-migration.json --format json > "$WORK"/rehearse-apply.json; step "migrate apply-safe"
W="$WORK" python3 - <<'PY'
import json, os
W = os.environ["W"]
d = json.load(open(W + "/rehearse-apply.json"))
r = {s["id"]: s["renames"] for s in d["steps"] if s["kind"] == "mint_ids"}
json.dump(r["task-ids"] + r["adr-ids"], open(W + "/rehearse-renames.json", "w"), indent=1, ensure_ascii=False)
PY
python3 "$HERE/adr_slugs.py" --root . --adr-dir "$V/9 pkm/adr" --renames "$WORK"/rehearse-renames.json; step "renames list, ADR slugs"
git add -A; git commit -qm "migrate apply-safe result (tool steps)"
python3 "$HERE/link_repair.py" --root . --renames "$WORK"/rehearse-renames.json --report "$WORK"/rehearse-linkrep.json; step "link repair + renames (emulated)"
git add -A; git commit -qm "link repair and renames"
python3 "$HERE/adr_prose.py" --root . --adr-dir "$V/9 pkm/adr" --include "$V/9 pkm" --include "$V/2 PROJECTS" --include "$V/3 AREAS" --include "$V/0 Inbox" --include .agents --exclude "$V/3 AREAS/AMEX" --report "$WORK"/rehearse-prose.json --apply; step "ADR prose rewrite"
git add -A; git commit -qm "ADR prose mentions"
python3 "$HERE/apply_summaries.py" "$AKM" . "$V/9 pkm/adr" "$HERE/adr_summaries.json"; step "ADR summaries (akmctl set x21)"
"$AKM" index; step "akmctl index"
"$AKM" check --format json > "$WORK"/rehearse-check.json || true; step "akmctl check"
python3 "$HERE/verify.py" --repo . --before pre-migrate-0.0.0-to-1.0.0 --after HEAD --migration "$WORK"/rehearse-migration.json --renames "$WORK"/rehearse-renames.json --tool-dirs "$V/9 pkm/tasks/" "$V/9 pkm/adr/" --out "$WORK"/rehearse-verify.json || true; step "verify (a)-(e)"
printf 'total %6.2fs\n' "$(python3 -c "import time;print(time.time()-$T0)")"
