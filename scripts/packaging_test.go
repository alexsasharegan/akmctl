// Package scripts_test runs the packaging scripts on a made-up release. The
// scripts read and write files only; they publish nothing.
package scripts_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const sum = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func sh(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the scripts are POSIX shell; the Windows job runs the WinGet one with bash in CI")
	}
	p, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("no sh")
	}
	return p
}

func sums(t *testing.T, lines ...string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "SHA256SUMS")
	if err := os.WriteFile(p, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

// TestWingetManifests is work item 261108-fabf: the three manifests carry the
// version, the URL of the Windows archive and its checksum, and a missing
// checksum is an error, not an empty field.
func TestWingetManifests(t *testing.T) {
	shell := sh(t)
	s := sums(t, sum+"  akmctl_9.9.9_linux_amd64.tar.gz", sum[:8]+strings.Repeat("a", 56)+"  akmctl_9.9.9_windows_amd64.zip")
	out := filepath.Join(t.TempDir(), "m")
	if b, err := exec.Command(shell, "winget-manifests.sh", "v9.9.9", s, out).CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, b)
	}
	read := func(name string) string {
		b, err := os.ReadFile(filepath.Join(out, "alexsasharegan.akmctl"+name))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	installer := read(".installer.yaml")
	for _, want := range []string{
		"PackageVersion: 9.9.9\n",
		"InstallerUrl: https://github.com/alexsasharegan/akmctl/releases/download/v9.9.9/akmctl_9.9.9_windows_amd64.zip\n",
		"InstallerSha256: 01234567" + strings.Repeat("A", 56) + "\n",
		"RelativeFilePath: akmctl_9.9.9_windows_amd64\\akmctl.exe\n",
		"NestedInstallerType: portable\n",
		"ManifestType: installer\n",
	} {
		if !strings.Contains(installer, want) {
			t.Errorf("installer manifest lacks %q:\n%s", want, installer)
		}
	}
	if v := read(".yaml"); !strings.Contains(v, "ManifestType: version\n") || !strings.Contains(v, "PackageVersion: 9.9.9\n") {
		t.Errorf("version manifest:\n%s", v)
	}
	if l := read(".locale.en-US.yaml"); !strings.Contains(l, "ManifestType: defaultLocale\n") || !strings.Contains(l, "LicenseUrl: https://github.com/alexsasharegan/akmctl/blob/v9.9.9/LICENSE\n") {
		t.Errorf("locale manifest:\n%s", l)
	}
	// no checksum for the Windows archive, and a bad version, are refused
	bad := sums(t, sum+"  akmctl_9.9.9_linux_amd64.tar.gz")
	if b, err := exec.Command(shell, "winget-manifests.sh", "9.9.9", bad, filepath.Join(t.TempDir(), "x")).CombinedOutput(); err == nil || !strings.Contains(string(b), "no checksum") {
		t.Errorf("missing checksum: %v %s", err, b)
	}
	if _, err := exec.Command(shell, "winget-manifests.sh", "9 9", s, filepath.Join(t.TempDir(), "y")).CombinedOutput(); err == nil {
		t.Error("bad version accepted")
	}
}
