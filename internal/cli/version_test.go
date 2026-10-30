package cli

import (
	"runtime/debug"
	"testing"
)

func TestResolveVersion(t *testing.T) {
	rev := []debug.BuildSetting{{Key: "vcs.revision", Value: "0123456789abcdef"}}
	cases := []struct {
		name, v string
		bi      debug.BuildInfo
		want    string
	}{
		{"linked", "0.2.0", debug.BuildInfo{Main: debug.Module{Version: "v9.9.9"}}, "0.2.0"},
		{"go install", "0.0.0-dev", debug.BuildInfo{Main: debug.Module{Version: "v0.2.1"}}, "0.2.1"},
		{"dev build", "0.0.0-dev", debug.BuildInfo{Main: debug.Module{Version: "(devel)"}, Settings: rev}, "0.0.0-dev+0123456"},
		{"no info", "0.0.0-dev", debug.BuildInfo{}, "0.0.0-dev"},
	}
	for _, c := range cases {
		if got := resolveVersion(c.v, &c.bi); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}
