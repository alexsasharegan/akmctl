#!/bin/sh
# Write the WinGet manifests for one release, from its SHA256SUMS.
#
#   scripts/winget-manifests.sh 0.3.3 dist/0.3.3/SHA256SUMS out/
#
# Output: three files in the output folder, named as winget-pkgs wants them:
# alexsasharegan.akmctl.yaml (version), .installer.yaml and .locale.en-US.yaml.
# It reads and writes files only; it submits nothing (work item 261108-fabf).
set -eu

version=${1:?usage: scripts/winget-manifests.sh <version> <SHA256SUMS> <outdir>}
sums=${2:?usage: scripts/winget-manifests.sh <version> <SHA256SUMS> <outdir>}
out=${3:?usage: scripts/winget-manifests.sh <version> <SHA256SUMS> <outdir>}
version=${version#v}
case $version in
*[!0-9A-Za-z.+-]* | '') echo "winget-manifests.sh: bad version: $version" >&2; exit 2 ;;
esac

id=alexsasharegan.akmctl
repo=alexsasharegan/akmctl
zip=akmctl_${version}_windows_amd64.zip
sha=$(awk -v f="$zip" '$2 == f || $2 == "*" f { print toupper($1) }' "$sums")
case $sha in
[0-9A-F][0-9A-F]*) ;;
*) echo "winget-manifests.sh: no checksum for $zip in $sums" >&2; exit 1 ;;
esac
[ ${#sha} -eq 64 ] || { echo "winget-manifests.sh: bad checksum for $zip" >&2; exit 1; }

mkdir -p "$out"
schema=1.6.0

cat > "$out/$id.yaml" <<YAML
# Created by scripts/winget-manifests.sh
PackageIdentifier: $id
PackageVersion: $version
DefaultLocale: en-US
ManifestType: version
ManifestVersion: $schema
YAML

cat > "$out/$id.installer.yaml" <<YAML
# Created by scripts/winget-manifests.sh
PackageIdentifier: $id
PackageVersion: $version
InstallerType: zip
NestedInstallerType: portable
NestedInstallerFiles:
  - RelativeFilePath: akmctl_${version}_windows_amd64\\akmctl.exe
    PortableCommandAlias: akmctl
Installers:
  - Architecture: x64
    InstallerUrl: https://github.com/$repo/releases/download/v$version/$zip
    InstallerSha256: $sha
ManifestType: installer
ManifestVersion: $schema
YAML

cat > "$out/$id.locale.en-US.yaml" <<YAML
# Created by scripts/winget-manifests.sh
PackageIdentifier: $id
PackageVersion: $version
PackageLocale: en-US
Publisher: alexsasharegan
PublisherUrl: https://github.com/alexsasharegan
PackageName: akmctl
PackageUrl: https://github.com/$repo
License: MIT
LicenseUrl: https://github.com/$repo/blob/v$version/LICENSE
ShortDescription: Command-line tool and API for agent knowledge management (AKM).
Moniker: akmctl
Tags:
  - agents
  - markdown
  - knowledge-management
ReleaseNotesUrl: https://github.com/$repo/releases/tag/v$version
ManifestType: defaultLocale
ManifestVersion: $schema
YAML

echo "wrote $out/$id.yaml, $id.installer.yaml and $id.locale.en-US.yaml"
