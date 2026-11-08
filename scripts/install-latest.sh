#!/bin/sh
# Install the newest commit of main whose CI passed as the global akmctl.
#
#   scripts/install-latest.sh            install it, or say it is already installed
#   scripts/install-latest.sh --dry-run  say what it would do, change nothing
#
# This is the maintainer's development channel (work item 261108-8m4a), not the
# release: users get tagged releases. "Passed" means the newest completed `ci`
# run of a push to main succeeded in every job (macOS, Ubuntu and Windows).
# The build is made by `go install ...@<commit>` from a clean checkout of that
# commit, never from a working tree, so it is never stamped +dirty. Needs gh
# (logged in) and go. Set AKMCTL_REPO to use another repository.
set -eu

dry=0
case ${1:-} in
--dry-run) dry=1 ;;
'') ;;
*) echo "usage: scripts/install-latest.sh [--dry-run]" >&2; exit 2 ;;
esac

repo=${AKMCTL_REPO:-alexsasharegan/akmctl}
need() { command -v "$1" >/dev/null 2>&1 || { echo "install-latest.sh: $1 is not installed" >&2; exit 1; }; }
need gh
need go

# Newest completed, successful run of a push to main with every job successful.
sha=
for id in $(gh run list --repo "$repo" --branch main --workflow ci.yml --event push --limit 20 \
	--json databaseId,status,conclusion -q '.[]|select(.status=="completed" and .conclusion=="success")|.databaseId'); do
	jobs=$(gh run view "$id" --repo "$repo" --json jobs -q '[.jobs[]|select(.conclusion=="success")]|length')
	total=$(gh run view "$id" --repo "$repo" --json jobs -q '.jobs|length')
	if [ "$total" -ge 3 ] && [ "$jobs" = "$total" ]; then
		sha=$(gh run view "$id" --repo "$repo" --json headSha -q .headSha)
		break
	fi
done
[ -n "$sha" ] || { echo "install-latest.sh: no green ci run of main found for $repo" >&2; exit 1; }
want=$(printf %s "$sha" | cut -c1-12)

# The version string ends in the 12-character commit; a +dirty build never matches.
have=
if command -v akmctl >/dev/null 2>&1; then
	have=$(akmctl version 2>/dev/null | sed -n 's/^akmctl [^ ]*-\([0-9a-f]\{12\}\)\( .*\)\{0,1\}$/\1/p')
fi
if [ "$have" = "$want" ]; then
	echo "akmctl is up to date at $want"
	exit 0
fi

if [ "$dry" = 1 ]; then
	echo "would install $repo at $want (installed: ${have:-none or not a clean main build})"
	exit 0
fi

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
(cd "$work" && GOPROXY=direct GOFLAGS=-mod=mod go install "github.com/$repo/cmd/akmctl@$sha")
echo "installed akmctl at $want: $(akmctl version)"
