package cli

import "runtime/debug"

func init() {
	// A development build names the commit it was built from (semver build metadata).
	if Version != "0.0.0-dev" {
		return
	}
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, s := range bi.Settings {
			if s.Key == "vcs.revision" && len(s.Value) >= 7 {
				Version += "+" + s.Value[:7]
			}
		}
	}
}
