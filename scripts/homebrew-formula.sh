#!/bin/sh
# Write the Homebrew formula for one release, from its SHA256SUMS.
#
#   scripts/homebrew-formula.sh 0.3.3 dist/0.3.3/SHA256SUMS akmctl.rb
#
# The formula installs the release archive for macOS and Linux (arm64 and
# amd64). It reads and writes files only; it publishes nothing (work item
# 261108-43gk).
set -eu

usage="usage: scripts/homebrew-formula.sh <version> <SHA256SUMS> <output file>"
version=${1:?$usage}
sums=${2:?$usage}
out=${3:?$usage}
version=${version#v}
case $version in
*[!0-9A-Za-z.+-]* | '') echo "homebrew-formula.sh: bad version: $version" >&2; exit 2 ;;
esac

repo=alexsasharegan/akmctl
sha() {
	name=akmctl_${version}_$1.tar.gz
	s=$(awk -v f="$name" '$2 == f || $2 == "*" f { print tolower($1) }' "$sums")
	if [ ${#s} -ne 64 ]; then
		echo "homebrew-formula.sh: no checksum for $name in $sums" >&2
		exit 1
	fi
	echo "$s"
}
url() { echo "https://github.com/$repo/releases/download/v$version/akmctl_${version}_$1.tar.gz"; }

darwin_arm64=$(sha darwin_arm64)
darwin_amd64=$(sha darwin_amd64)
linux_arm64=$(sha linux_arm64)
linux_amd64=$(sha linux_amd64)

mkdir -p "$(dirname "$out")"
cat > "$out" <<RUBY
# Created by scripts/homebrew-formula.sh
class Akmctl < Formula
  desc "Command-line tool and API for agent knowledge management (AKM)"
  homepage "https://github.com/$repo"
  version "$version"
  license "MIT"

  on_macos do
    on_arm do
      url "$(url darwin_arm64)"
      sha256 "$darwin_arm64"
    end
    on_intel do
      url "$(url darwin_amd64)"
      sha256 "$darwin_amd64"
    end
  end

  on_linux do
    on_arm do
      url "$(url linux_arm64)"
      sha256 "$linux_arm64"
    end
    on_intel do
      url "$(url linux_amd64)"
      sha256 "$linux_amd64"
    end
  end

  def install
    bin.install "akmctl"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/akmctl version")
  end
end
RUBY
echo "wrote $out"
