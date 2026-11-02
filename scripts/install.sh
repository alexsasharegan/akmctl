#!/bin/sh
# Install the latest akmctl release.
#
#   curl -fsSL https://raw.githubusercontent.com/alexsasharegan/akmctl/main/scripts/install.sh | sh
#
# Picks the archive for this OS and CPU, checks it against SHA256SUMS and puts
# `akmctl` in $AKMCTL_INSTALL_DIR (default ~/.local/bin). Fetching with curl
# keeps macOS from quarantining the binary, as it does for a browser download.
#
#   AKMCTL_VERSION=0.3.3     install that version instead of the latest
#   AKMCTL_INSTALL_DIR=DIR   install into DIR
set -eu

repo=alexsasharegan/akmctl
dir=${AKMCTL_INSTALL_DIR:-$HOME/.local/bin}

die() {
	echo "install.sh: $*" >&2
	exit 1
}

case $(uname -s) in
Darwin) os=darwin ;;
Linux) os=linux ;;
*) die "unsupported system $(uname -s); on Windows download the .zip from https://github.com/$repo/releases" ;;
esac
case $(uname -m) in
arm64 | aarch64) arch=arm64 ;;
x86_64 | amd64) arch=amd64 ;;
*) die "unsupported CPU $(uname -m)" ;;
esac

version=${AKMCTL_VERSION:-}
if [ -z "$version" ]; then
	# The redirect of /releases/latest names the tag and costs no API quota.
	url=$(curl -fsSIL -o /dev/null -w '%{url_effective}' "https://github.com/$repo/releases/latest") ||
		die "could not look up the latest release"
	version=${url##*/}
fi
version=${version#v}
case $version in
*[!0-9A-Za-z.+-]* | '') die "bad version: $version" ;;
esac

name=akmctl_${version}_${os}_${arch}
base=https://github.com/$repo/releases/download/v$version
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "downloading akmctl $version for $os/$arch"
curl -fsSL -o "$tmp/$name.tar.gz" "$base/$name.tar.gz" || die "no archive $name.tar.gz in release v$version"
curl -fsSL -o "$tmp/SHA256SUMS" "$base/SHA256SUMS" || die "no SHA256SUMS in release v$version"

want=$(grep " $name.tar.gz\$" "$tmp/SHA256SUMS" | cut -d' ' -f1)
[ -n "$want" ] || die "SHA256SUMS has no line for $name.tar.gz"
if command -v sha256sum >/dev/null 2>&1; then
	got=$(sha256sum "$tmp/$name.tar.gz" | cut -d' ' -f1)
else
	got=$(shasum -a 256 "$tmp/$name.tar.gz" | cut -d' ' -f1)
fi
[ "$got" = "$want" ] || die "checksum mismatch for $name.tar.gz"

tar -xzf "$tmp/$name.tar.gz" -C "$tmp"
mkdir -p "$dir"
cp "$tmp/$name/akmctl" "$dir/akmctl.new"
chmod 755 "$dir/akmctl.new"
mv "$dir/akmctl.new" "$dir/akmctl"
echo "installed $dir/akmctl"

case :$PATH: in
*:"$dir":*) ;;
*) echo "add it to your PATH:  export PATH=\"$dir:\$PATH\"" ;;
esac
