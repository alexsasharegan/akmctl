#!/bin/sh
# Build the release archives for every supported platform.
#
#   scripts/release.sh 0.1.0
#
# Output: dist/<version>/akmctl_<version>_<os>_<arch>.tar.gz (.zip for windows)
# and dist/<version>/SHA256SUMS. Needs Go, tar, zip and shasum or sha256sum.
# Set SKIP_TESTS=1 to skip go test.
set -eu

version=${1:?usage: scripts/release.sh <version>, for example 0.1.0}
version=${version#v}
case $version in
*[!0-9A-Za-z.+-]* | '') echo "release.sh: bad version: $version" >&2; exit 2 ;;
esac

cd "$(dirname "$0")/.."
module=$(go list -m)
out=dist/$version
rm -rf "$out"
mkdir -p "$out"

if [ "${SKIP_TESTS:-}" != 1 ]; then
	go vet ./...
	go test ./...
fi

targets="darwin/arm64 darwin/amd64 linux/amd64 linux/arm64 windows/amd64"
for t in $targets; do
	os=${t%/*}
	arch=${t#*/}
	name=akmctl_${version}_${os}_${arch}
	stage=$out/.stage/$name
	mkdir -p "$stage"
	bin=akmctl
	[ "$os" = windows ] && bin=akmctl.exe
	CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -trimpath \
		-ldflags "-s -w -X $module/internal/cli.Version=$version" \
		-o "$stage/$bin" ./cmd/akmctl
	cp LICENSE README.md CHANGELOG.md "$stage/"
	if [ "$os" = windows ]; then
		(cd "$out/.stage" && zip -qr "../$name.zip" "$name")
	else
		tar -czf "$out/$name.tar.gz" -C "$out/.stage" "$name"
	fi
	echo "built $name"
done
rm -rf "$out/.stage"

cd "$out"
if command -v sha256sum >/dev/null 2>&1; then
	sha256sum akmctl_* > SHA256SUMS
else
	shasum -a 256 akmctl_* > SHA256SUMS
fi
cat SHA256SUMS
