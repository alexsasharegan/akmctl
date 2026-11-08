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

	"github.com/alexsasharegan/akmctl/internal/testutil"
)

const sum = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func sh(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the scripts are POSIX shell; the Windows job runs the WinGet one with bash in CI")
	}
	return testutil.RequireTool(t, "sh")
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

// TestHomebrewFormula is work item 261108-43gk: the formula names the archive
// and checksum of each of the four macOS and Linux targets, and a missing
// checksum is an error.
func TestHomebrewFormula(t *testing.T) {
	shell := sh(t)
	var lines []string
	for i, t := range []string{"darwin_arm64", "darwin_amd64", "linux_arm64", "linux_amd64"} {
		lines = append(lines, strings.Repeat(string(rune('a'+i)), 64)+"  akmctl_9.9.9_"+t+".tar.gz")
	}
	out := filepath.Join(t.TempDir(), "f", "akmctl.rb")
	if b, err := exec.Command(shell, "homebrew-formula.sh", "v9.9.9", sums(t, lines...), out).CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, b)
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	f := string(b)
	for i, target := range []string{"darwin_arm64", "darwin_amd64", "linux_arm64", "linux_amd64"} {
		url := "url \"https://github.com/alexsasharegan/akmctl/releases/download/v9.9.9/akmctl_9.9.9_" + target + ".tar.gz\"\n"
		sum := "sha256 \"" + strings.Repeat(string(rune('a'+i)), 64) + "\"\n"
		at := strings.Index(f, url)
		if at < 0 || !strings.HasPrefix(strings.TrimLeft(f[at+len(url):], " "), sum) {
			t.Errorf("%s: url and sha256 not together in:\n%s", target, f)
		}
	}
	for _, want := range []string{"class Akmctl < Formula", "version \"9.9.9\"", "bin.install \"akmctl\"", "license \"MIT\""} {
		if !strings.Contains(f, want) {
			t.Errorf("formula lacks %q", want)
		}
	}
	if b, err := exec.Command(shell, "homebrew-formula.sh", "9.9.9", sums(t, lines[:3]...), filepath.Join(t.TempDir(), "x.rb")).CombinedOutput(); err == nil || !strings.Contains(string(b), "no checksum") {
		t.Errorf("missing checksum: %v %s", err, b)
	}
}

// fakeTools writes gh, go and akmctl stand-ins into a folder for
// install-latest.sh: gh lists the runs in FAKE_RUNS and answers per-run queries
// from FAKE_SHA_<id>, FAKE_OK_<id> and FAKE_TOTAL_<id>; go records its arguments
// in $FAKE_LOG; akmctl prints FAKE_VERSION.
func fakeTools(t *testing.T) (dir, log string) {
	t.Helper()
	dir = t.TempDir()
	log = filepath.Join(t.TempDir(), "go.log")
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("#!/bin/sh\n"+body), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write("gh", `
if [ "$1 $2" = "run list" ]; then printf '%s\n' $FAKE_RUNS; exit 0; fi
id=$3
for last; do :; done
case $last in
*headSha*) eval "echo \$FAKE_SHA_$id" ;;
*select*) eval "echo \$FAKE_OK_$id" ;;
*) eval "echo \$FAKE_TOTAL_$id" ;;
esac
`)
	write("go", `echo "$GOPROXY $GOFLAGS $*" >> "$FAKE_LOG"`)
	write("akmctl", `[ -n "$FAKE_VERSION" ] && echo "$FAKE_VERSION"; exit 0`)
	return dir, log
}

func runInstallLatest(t *testing.T, dir, log string, env []string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command("sh", append([]string{"install-latest.sh"}, args...)...)
	cmd.Env = append([]string{"PATH=" + dir + ":/usr/bin:/bin", "FAKE_LOG=" + log, "HOME=" + t.TempDir()}, env...)
	b, err := cmd.CombinedOutput()
	return string(b), err
}

// TestInstallLatest is work item 261108-8m4a: the script installs the newest
// main commit whose every job passed, by commit, and does nothing when that
// commit is installed.
func TestInstallLatest(t *testing.T) {
	sh(t)
	dir, log := fakeTools(t)
	sha1 := "1111111111111111111111111111111111111111" // newest run, but a job failed
	sha2 := "2222222222222222222222222222222222222222" // newest fully green
	env := []string{
		"FAKE_RUNS=101 102 103",
		"FAKE_SHA_101=" + sha1, "FAKE_OK_101=2", "FAKE_TOTAL_101=3",
		"FAKE_SHA_102=" + sha2, "FAKE_OK_102=3", "FAKE_TOTAL_102=3",
		"FAKE_SHA_103=3333333333333333333333333333333333333333", "FAKE_OK_103=3", "FAKE_TOTAL_103=3",
	}
	logged := func() string { b, _ := os.ReadFile(log); return string(b) }

	// an older or dirty build is replaced by the newest green commit, by its full hash
	out, err := runInstallLatest(t, dir, log, append(env, "FAKE_VERSION=akmctl 0.3.4-0.20261108000000-aaaaaaaaaaaa+dirty (API 1.0)"))
	if err != nil || !strings.HasPrefix(logged(), "direct -mod=mod install github.com/alexsasharegan/akmctl/cmd/akmctl@"+sha2) {
		t.Fatalf("install: %v\n%s\nlog: %s", err, out, logged())
	}
	if strings.Contains(logged(), sha1) {
		t.Errorf("installed the commit with a failed job: %s", logged())
	}

	// the same commit is left alone
	os.Remove(log)
	out, err = runInstallLatest(t, dir, log, append(env, "FAKE_VERSION=akmctl 0.3.4-0.20261108000000-222222222222 (API 1.0)"))
	if err != nil || !strings.Contains(out, "up to date at 222222222222") || logged() != "" {
		t.Errorf("up to date: %v\n%s\nlog: %s", err, out, logged())
	}

	// --dry-run changes nothing
	out, err = runInstallLatest(t, dir, log, append(env, "FAKE_VERSION=akmctl 0.3.4-0.20261108000000-aaaaaaaaaaaa (API 1.0)"), "--dry-run")
	if err != nil || !strings.Contains(out, "would install") || logged() != "" {
		t.Errorf("dry run: %v\n%s\nlog: %s", err, out, logged())
	}

	// no fully green run is an error and installs nothing
	out, err = runInstallLatest(t, dir, log, []string{"FAKE_RUNS=101", "FAKE_SHA_101=" + sha1, "FAKE_OK_101=2", "FAKE_TOTAL_101=3"})
	if err == nil || !strings.Contains(out, "no green ci run") || logged() != "" {
		t.Errorf("no green run: %v\n%s\nlog: %s", err, out, logged())
	}
	if _, err := runInstallLatest(t, dir, log, nil, "--bad"); err == nil {
		t.Error("a bad flag was accepted")
	}
}
