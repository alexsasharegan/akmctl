package scripts_test

import (
	"os/exec"
	"strings"
	"testing"
)

// TestCIRetry is work item 261109-e540: a transient download error is retried
// and then skipped with a warning, and any other failure, a checksum mismatch
// included, fails at once.
func TestCIRetry(t *testing.T) {
	shell := sh(t)
	const prelude = `. ./ci-retry.sh
CI_RETRY_PAUSE=0
count=$(mktemp)
echo 0 >"$count"
try() { n=$(cat "$count"); n=$((n + 1)); echo "$n" >"$count"; }
`
	for _, c := range []struct {
		name, body string
		exit       int
		tries      string
		output     string
	}{
		{"a 504 twice, then success", `cmd() { try; if [ "$n" -lt 3 ]; then echo "curl: (22) The requested URL returned error: 504"; return 22; fi; echo fetched; }`, 0, "3", "fetched"},
		{"a 504 every time is a warning and exit 0", `cmd() { try; echo "HTTP 504: Gateway Timeout"; return 1; }`, 0, "3", "::warning title=download failed::"},
		{"a checksum mismatch fails at once", `cmd() { try; echo "Error: SHA256 mismatch"; return 4; }`, 4, "1", "SHA256 mismatch"},
		{"a missing command fails at once", `cmd() { try; echo "brew: command not found"; return 127; }`, 127, "1", "command not found"},
	} {
		t.Run(c.name, func(t *testing.T) {
			script := prelude + c.body + "\nstatus=0\n(ci_retry cmd || ci_skip_or_fail $?) || status=$?\necho \"tries=$(cat \"$count\")\"\nexit $status\n"
			cmd := exec.Command(shell, "-c", script)
			out, err := cmd.CombinedOutput()
			code := 0
			if ee, ok := err.(*exec.ExitError); ok {
				code = ee.ExitCode()
			} else if err != nil {
				t.Fatal(err)
			}
			if code != c.exit || !strings.Contains(string(out), "tries="+c.tries+"\n") || !strings.Contains(string(out), c.output) {
				t.Fatalf("exit %d (want %d), want tries=%s and %q in:\n%s", code, c.exit, c.tries, c.output, out)
			}
		})
	}
}
